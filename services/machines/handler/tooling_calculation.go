package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func MakeToolingMachineStatisticsCalculation(ms *MachineService) {
	//Get all data from Machine timeline event table
	dbConnection := ms.BaseService.ServiceDatabases[ProjectID]

	//Check whether we are getting messages or not. If yes set it as Live
	updateToolingMachineConnectStatus(ms, dbConnection)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	_, listOfEvents := productionOrderInterface.GetToolingScheduledEvents(ProjectID)
	ms.BaseService.Logger.Info("events processingss", zap.Any("number_of_events", len(*listOfEvents)))

	// Get production order complete status id from preference level
	orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(ProjectID, ScheduleStatusPreferenceSeven)

	for _, eventTimeLine := range *listOfEvents {
		//Decode Machine timeline event info
		etlInfo := make(map[string]interface{})
		err := json.Unmarshal(eventTimeLine.ObjectInfo, &etlInfo)
		if err == nil {
			toolingMachineId := util.InterfaceToInt(etlInfo["machineId"])
			scheduleStartDate := util.InterfaceToString(etlInfo["startDate"])
			scheduleName := util.InterfaceToString(etlInfo["name"])
			scheduleEndDate := util.InterfaceToString(etlInfo["endDate"])
			scheduleEventStatus := util.InterfaceToInt(etlInfo["eventStatus"])
			scheduleEventSourceId := util.InterfaceToInt(etlInfo["eventSourceId"])
			scheduleProductionOrder := util.InterfaceToString(etlInfo["productionOrder"])
			schedulePartId := util.InterfaceToInt(etlInfo["partId"])

			machineStatus, newMachineId := getToolingMachineStatus(ms, dbConnection, toolingMachineId)
			//Read schedule start date and end date
			startDate := convertStringToDateTime(scheduleStartDate)
			endDate := convertStringToDateTime(scheduleEndDate)
			if startDate.Error == nil || endDate.Error == nil {

				//Check end date is overdue or not
				//Adding singapore time in utc time
				// event time is not passed in under
				// Event status in completed then we don't do calculation
				if time.Now().Before(endDate.DateTime) && scheduleEventStatus != orderStatusId {

					startEpoch := startDate.DateTimeEpochMilli
					endEpoch := endDate.DateTimeEpochMilli
					topicName := getToolingTopicName(toolingMachineId, ms, ProjectID)

					hmiInfoList := getHMIInfo(ms, ProjectID, toolingMachineId, eventTimeLine.Id, ToolingMachineHmiTable)
					ms.BaseService.Logger.Info("Tooling HMI info list", zap.Any("machine_id:", toolingMachineId), zap.Any("list", hmiInfoList))
					if len(hmiInfoList) > 0 {
						// getting messages based on HMI first created time and current time for processing
						// we needed only the in_timestamp message greater
						messageQuery := "select * from message where body->>'$.in_timestamp' >= " + strconv.FormatInt(convertStringToDateTime(hmiInfoList[0].Created).DateTimeEpochMilli, 10) + " and body->>'$.in_timestamp' <= " + strconv.FormatInt(time.Now().UnixMilli(), 10) + " and topic='" + topicName + "' order by body->>'$.in_timestamp' desc"
						ms.BaseService.Logger.Info("tooling running query to select the messages", zap.Any("query", messageQuery))
						var messages []Message
						dbConnection.Raw(messageQuery).Scan(&messages)
						//messageOffSet = getCycleCountOffset(messages)
						//ms.BaseService.Logger.Info("message offset", "offset", messageOffSet)
						// Fetch data based on the schedule from message table
						if len(messages) > 0 || topicName == "machines/L1_Makino_F3Graphite" {

							err, productionOrderObject := productionOrderInterface.GetToolingProductionOrderInfo(ProjectID, scheduleEventSourceId, toolingMachineId)
							if err != nil {
								ms.BaseService.Logger.Error("error getting tooling production order info", zap.Any("event_source_id", scheduleEventSourceId), zap.Any("machine_id", toolingMachineId))
								continue
							}

							_, toolingPartMaster := productionOrderInterface.GetToolingPartById(ProjectID, schedulePartId)
							toolingPartMasterInfo := make(map[string]interface{})
							json.Unmarshal(toolingPartMaster.ObjectInfo, &toolingPartMasterInfo)

							//totalDuration := util.InterfaceToFloat(partInfo["day"])*24 + util.InterfaceToFloat(partInfo["hour"]) + util.InterfaceToFloat(partInfo["minute"])/60
							totalDuration := endDate.DateTimeEpoch - startDate.DateTimeEpoch
							cycleTime := int(totalDuration)
							//cycleTime := util.InterfaceToInt(toolingPartMasterInfo["day"])*86400 + util.InterfaceToInt(toolingPartMasterInfo["hour"])*3600 + util.InterfaceToInt(toolingPartMasterInfo["minute"])*60
							scheduleScheduledQty := cycleTime

							productionOrderInfo := make(map[string]interface{})
							json.Unmarshal(productionOrderObject.ObjectInfo, &productionOrderInfo)

							//productionOrderQty := util.InterfaceToInt(productionOrderInfo["day"])*86400 + util.InterfaceToInt(productionOrderInfo["hour"])*3600 + util.InterfaceToInt(productionOrderInfo["minute"])*60
							partNo := util.InterfaceToIntArray(productionOrderInfo["partNo"])
							productionOrderQty := getOverallPlannedQty(ProjectID, partNo, productionOrderInterface)
							ms.BaseService.Logger.Info("tooling production order:", zap.Any("machine_id:", toolingMachineId), zap.Any("production order info ", productionOrderInfo))

							actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, toolingMachineId, eventTimeLine.Id, ToolingMachineHmiTable)
							// accumulatedCycleCount := getAccumulatedCycleCount(messages)
							// how much manufactured so far from PUMAS
							var cycleCount int
							var cycleCountOneMinBefore int
							var remark string
							if topicName == "machines/L1_Makino_F3Graphite" {
								cycleCount = getGraphiteMachineCycleCount(actualStartTime)
								cycleCountOneMinBefore = cycleTime - 60
							} else {
								cycleCount, remark = getToolingCycleCount(messages)
								cycleCountOneMinBefore = getToolingCycleCountBeforeOneMin(messages)
							}

							ms.BaseService.Logger.Info("tooling cycle count", zap.Any("machine_id:", toolingMachineId), zap.Any("cycleCount", cycleCount), zap.Any("cycleCountOneMinBefore", cycleCountOneMinBefore))

							//based on UI or HMI
							dailyRejectedCount := getTotalRejectQtyForEvent(hmiInfoList)

							overallRejectedCount := getToolingOverallRejectedQuantity(ms, dbConnection, scheduleEventSourceId)
							//getOverallRejectedQuantity(ms, dbConnection, etlInfo.ProductionOrder, eventTimeLines)
							goodCount := getGoodCount(cycleCount, dailyRejectedCount)
							machineDownTime := getToolingDowntime(hmiInfoList, dbConnection, toolingMachineId, startEpoch, endEpoch)

							plannedProductionTime := getPlannedProductionTime(hmiInfoList) // this is in milliseconds
							// plannedProductionTime := endEpoch - startEpoch
							availability := getAvailability(plannedProductionTime, machineDownTime)

							runTime := plannedProductionTime - machineDownTime

							performance := (float64(cycleTime) * float64(cycleCount)) / float64(runTime/1000)
							var quality float64
							quality = 0.0
							if cycleCount > 0 {
								quality = float64(goodCount) / float64(cycleCount)
							}

							if performance > float64(1) {
								performance = float64(1)
							}

							if availability > float64(1) {
								availability = float64(1)
							}

							if quality > float64(1) {
								quality = float64(1)
							}

							oee := availability * performance * quality // 0.7479

							if oee > float64(1) {
								oee = float64(1)
							}
							_, scheduleOrderNeighbours := productionOrderInterface.GetToolingScheduledEventsByProductionId(ProjectID, scheduleEventSourceId)
							completedValue := getToolingCompletedValue(eventTimeLine.Id, cycleCount, scheduleProductionOrder, toolingMachineId, scheduleOrderNeighbours, ms, ProjectID)
							progressPercentage := getProgressPercentage(cycleCount, scheduleScheduledQty)
							ms.BaseService.Logger.Info("summary", zap.Any("machine_id", toolingMachineId),
								zap.Any("planned_production_time", plannedProductionTime), zap.Any("machine_down_time", machineDownTime), zap.Any("performance", performance),
								zap.Any("availability", availability), zap.Any("completedValue", completedValue), zap.Any("progressPercentage", progressPercentage),
								zap.Any("actualStartTime", actualStartTime), zap.Any("actualEndTime", actualEndTime), zap.Any("quality", quality))
							statsInfo := MachineStatisticsInfo{
								EventId:                    eventTimeLine.Id,
								ProductionOrderId:          scheduleEventSourceId,
								CurrentStatus:              machineStatus,
								PartId:                     schedulePartId,
								ScheduleStartTime:          scheduleStartDate,
								ScheduleEndTime:            scheduleEndDate,
								ActualStartTime:            actualStartTime,
								ActualEndTime:              actualEndTime,
								EstimatedEndTime:           getEstimatedEndTime(float32(cycleTime), cycleCount, actualStartTime, scheduleScheduledQty),
								Oee:                        int(oee * 10000),
								Availability:               int(availability * 10000),
								Performance:                int(performance * 10000),
								Quality:                    int(quality * 10000),
								PlannedQuality:             productionOrderQty,
								DailyPlannedQty:            scheduleScheduledQty,
								Completed:                  completedValue,
								Rejects:                    dailyRejectedCount,
								OverallRejectedQty:         overallRejectedCount,
								CompletedPercentage:        getDailyCompletedPercentage(scheduleScheduledQty, cycleCount),
								OverallCompletedPercentage: getOverallCompletedPercentage(productionOrderQty, completedValue),
								Actual:                     cycleCount,
								ProgressPercentage:         progressPercentage,
								DownTime:                   int(machineDownTime),
								ExpectedProductQyt:         getExpectedProductQty(cycleCount, cycleCountOneMinBefore, actualStartTime, endDate.DateTimeEpoch),
								WarningMessage:             make([]string, 0),
								Remark:                     remark,
							}
							statsInfoJson, _ := json.Marshal(statsInfo)
							machineStatistics := ToolingMachineStatistics{
								MachineId: toolingMachineId,
								TS:        time.Now().Unix(),
								StatsInfo: statsInfoJson,
							}
							//ms.BaseService.Logger.Info("machineStatistics", "stats", string(statsInfoJson))
							err = dbConnection.Create(&machineStatistics).Error
							if err != nil {
								ms.BaseService.Logger.Error("inserting stats has failed", zap.String("error", err.Error()))
							}
							updateToolingMachineTimeLine(ProjectID, eventTimeLine.Id, cycleCount, progressPercentage, dailyRejectedCount, int(oee*10000))

							if topicName == "machines/L1_Makino_F3Graphite" {
								var cycleCountMessage = getFirstGraphiteMachineCycleCount(actualStartTime)
								serializeMsg, err := json.Marshal(cycleCountMessage)
								var updates = make(map[string]interface{})
								if err != nil {
									updates = map[string]interface{}{
										"starting_cycle_count": 0,
									}
								} else {
									updates = map[string]interface{}{
										"starting_cycle_count": 0,
										"message_cycle_count":  serializeMsg,
									}
								}

								err = ms.ViewManager.CreateOrUpdateToolingView(toolingMachineId, updates)
								if err != nil {
									ms.BaseService.Logger.Error("error updating tooling view", zap.String("error", err.Error()))
								}
							} else {
								var initialCycleCount, cycleCountMessage = getFirstToolingCycleCount(messages)
								serializeMsg, err := json.Marshal(cycleCountMessage)
								var updates = make(map[string]interface{})
								if err != nil {
									updates = map[string]interface{}{
										"starting_cycle_count": initialCycleCount,
									}
								} else {
									updates = map[string]interface{}{
										"starting_cycle_count": initialCycleCount,
										"message_cycle_count":  serializeMsg,
									}
								}

								err = ms.ViewManager.CreateOrUpdateToolingView(toolingMachineId, updates)
								if err != nil {
									ms.BaseService.Logger.Error("error updating tooling view", zap.String("error", err.Error()))
								}
							}

						}
					} else {
						continue
					}
				} else {
					hmiInfoList := getHMIInfo(ms, ProjectID, toolingMachineId, eventTimeLine.Id, ToolingMachineHmiTable)

					if len(hmiInfoList) == 0 {
						continue
					}

					ms.BaseService.Logger.Info("machine is running out of time:", zap.Any("hmi_list", hmiInfoList))
					var lastHMIInfoResult MachineHMIInfoResult
					for index := len(hmiInfoList) - 1; index >= 0; index-- {
						if hmiInfoList[index].HmiStatus == "" {
							continue
						} else {
							lastHMIInfoResult = hmiInfoList[index]
							break
						}

					}
					ms.BaseService.Logger.Info("machine is running out of time:", zap.Any("last_hmi_info_result", lastHMIInfoResult))
					if lastHMIInfoResult.HmiStatus != "stopped" {
						//If the schedule is overdue but but still hmi isn't stopped the we have to append warning message
						hmiSettingInfo := getHmiSetting(dbConnection, ms, toolingMachineId)
						ms.BaseService.Logger.Info("hmi setting info", zap.Any("hmi_setting_info", hmiSettingInfo), zap.Any("endDate:", endDate))
						//End date should be in seconds
						appendToolingWarningMessage(ms, dbConnection, hmiSettingInfo.WarningMessageGenerationPeriod, endDate.DateTimeEpoch, eventTimeLine.Id, toolingMachineId, scheduleName, newMachineId, hmiSettingInfo.WarningTargetEmailId)
						addToolingForceStopFlag(ProjectID, eventTimeLine.Id)
						//If user not stop HMI data, Then program automatically stopped hmi
						// Condition
						// Based on hmi stop configuration hmi will be stopped
						automaticToolingHmiStop(hmiSettingInfo.HmiAutoStopPeriod, lastHMIInfoResult, dbConnection, ms, endDate.DateTimeEpoch)
					}

				}
			}

		}
	}
}

func getOverallPlannedQty(projectId string, partNoList []int, productionOrderInterface common.ProductionOrderInterface) int {
	var productionOrderQty int
	for _, partId := range partNoList {
		_, generalObject := productionOrderInterface.GetToolingPartById(projectId, partId)
		partInfo := make(map[string]interface{})
		json.Unmarshal(generalObject.ObjectInfo, &partInfo)

		productionOrderQty = productionOrderQty + util.InterfaceToInt(partInfo["day"])*86400 + util.InterfaceToInt(partInfo["hour"])*3600 + util.InterfaceToInt(partInfo["minute"])*60

	}

	return productionOrderQty
}

func getFirstGraphiteMachineCycleCount(actualStartTime string) map[string]interface{} {
	timeNow := time.Now().UTC().Unix()

	var cycleCountMessage = map[string]interface{}{
		"actualStartTime": actualStartTime,
		"currentTime":     timeNow,
	}

	return cycleCountMessage
}

func getGraphiteMachineCycleCount(actualStartTime string) int {
	timeNow := time.Now().UTC().Unix()
	startTimeObject := util.ConvertStringToDateTime(actualStartTime)

	cycleCount := timeNow - startTimeObject.DateTimeEpoch

	return int(cycleCount)
}

func automaticToolingHmiStop(hmiAutoStopPeriod string, info MachineHMIInfoResult, dbConnection *gorm.DB, ms *MachineService, scheduleStopTime int64) {
	//Schedule stop time in seconds
	duration, _ := time.ParseDuration(hmiAutoStopPeriod)
	durationInSeconds := duration.Seconds()
	fmt.Println("durationInSeconds: ", durationInSeconds)
	if time.Now().UTC().Unix() > (scheduleStopTime + int64(durationInSeconds)) {
		fmt.Println("is it exceeeded: ", time.Now().UTC().Unix(), " scheduleStopTime :", scheduleStopTime, " duration seconds :", int64(durationInSeconds))
		updateOverDueToolingHmiInfo(info, dbConnection, ms)
	}
}

func appendToolingWarningMessage(ms *MachineService, dbConnection *gorm.DB, warningMessageGenrationPeriod string, scheduleStopTime int64, eventId int, machineId int, eventName string, machineName string, targetEmailId string) {
	// Schedule stop time in seconds
	statisticsQuery := "select * from tooling_machine_statistics where stats_info ->> '$.eventId' = " + strconv.Itoa(eventId) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"
	var machineStatics ToolingMachineStatistics
	var machineStatsInfo MachineStatisticsInfo

	dbConnection.Raw(statisticsQuery).Scan(&machineStatics)

	_ = json.Unmarshal(machineStatics.StatsInfo, &machineStatsInfo)

	duration, _ := time.ParseDuration(warningMessageGenrationPeriod)
	durationInSeconds := duration.Seconds()

	// Find how many times wrning message should have been added
	timeDiff := time.Now().UTC().Unix() - scheduleStopTime
	noOfWaringMessage := timeDiff / int64(durationInSeconds)

	statsInfoarnMessageLength := len(machineStatsInfo.WarningMessage)

	if int(noOfWaringMessage) != statsInfoarnMessageLength {
		getCurrentSingaporeTime := util.GetZoneCurrentTime("Asia/Singapore")
		//Append warning message
		warningMessage := "Schedule is overrunning, do you want proceed stop? " + getCurrentSingaporeTime
		machineStatsInfo.WarningMessage = append(machineStatsInfo.WarningMessage, warningMessage)

		updatedMachineStatsInfo, _ := json.Marshal(machineStatsInfo)

		//Update last calculated machine statics
		dbError := dbConnection.Model(&ToolingMachineStatistics{}).Where("ts = ?", machineStatics.TS).Where("machine_id = ?", machineStatics.MachineId).Update("stats_info", updatedMachineStatsInfo).Error
		if dbError != nil {
			ms.BaseService.Logger.Error("error in updating machine statics with warning message", zap.String("error", dbError.Error()))
		}

		if targetEmailId != "" {
			ms.emailGenerator(eventName, machineName, targetEmailId, make([]string, 0))
		}

	}

}

func getToolingHmiSetting(dbConnection *gorm.DB, ms *MachineService, machineId int) *HmiSettingInfo {
	err, generalObject := Get(dbConnection, ToolingMachineHmiSettingTable, machineId)
	var hmiSettingInfo HmiSettingInfo
	if err != nil {
		ms.BaseService.Logger.Error("error getting machine hmi setting", zap.String("error", err.Error()))
		return nil
	}

	if generalObject.Id == 0 {
		ms.BaseService.Logger.Error("error no machine hmi setting is found", zap.String("error", err.Error()))
		return nil
	}

	err = json.Unmarshal(generalObject.ObjectInfo, &hmiSettingInfo)
	if err != nil {
		ms.BaseService.Logger.Info("fetched object", zap.Any("machine_id", machineId), zap.Any("object", string(generalObject.ObjectInfo)))
		ms.BaseService.Logger.Error("error getting machine hmi setting", zap.String("error", err.Error()))
		return nil
	}

	return &hmiSettingInfo
}

// If the user not stopped thselect * from fuyu_mes.tooling_machine_hmi where object_info->>'$.machineId' = 1 and object_info->>'$.eventId' = 74 order by object_info->>'$.createdAt') asc
// hmi info, program automatically is inserted stopped hmi
// This function only update created , status and remark attributes
func updateOverDueToolingHmiInfo(info MachineHMIInfoResult, dbConnection *gorm.DB, ms *MachineService) {
	objectInfo := MachineHMIInfo{
		CreatedAt: time.Now().Format("2006-01-02T15:04:05.000Z"),
		EventId:   info.EventId,
		Operator:  info.Operator,
		MachineId: info.MachineId,
		Status:    info.Status,
		HMIStatus: "stopped",
		Remark:    "HMI was stopped by program",
	}

	hmiObjectInfo, err := json.Marshal(objectInfo)

	if err != nil {
		ms.BaseService.Logger.Error("error in marshalling hmi info", zap.String("error", err.Error()))
	}

	dbError, _ := Create(dbConnection, ToolingMachineHmiTable, component.GeneralObject{ObjectInfo: hmiObjectInfo})

	if dbError != nil {
		ms.BaseService.Logger.Error("error in updating hmi info", zap.String("error", err.Error()))
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateToolingOrderPreferenceLevel(ProjectID, 2, info.EventId, ScheduleStatusPreferenceSix)

}

func getAllToolingMachines(ms *MachineService, dbConnection *gorm.DB) []ToolingMachineMaster {
	machineQuery := "select * from tooling_machine_master"
	var machineMasters []ToolingMachineMaster
	dbConnection.Raw(machineQuery).Scan(&machineMasters)

	return machineMasters
}

func searchOPCode(messageList []Message) (bool, bool) {
	foundOPCode := false
	opCodeValue := false

	for _, msg := range messageList {
		msgBody := make(map[string]interface{}, 0)

		json.Unmarshal(msg.Body, &msgBody)
		if valBody, ok := msgBody["OP"]; ok {
			opCodeValue = util.InterfaceToBool(valBody)

			foundOPCode = true
			break
		}

	}

	return foundOPCode, opCodeValue
}

// Update machine master whether we receive message or not. Live or Waiting for feed
func updateToolingMachineConnectStatus(ms *MachineService, dbConnection *gorm.DB) {
	// year, month, day := time.Now().Date()
	timeNow := time.Now().UTC().UnixMilli()
	machineMasters := getAllToolingMachines(ms, dbConnection)
	for _, machine := range machineMasters {
		machineMasterInfo := ToolingMachineMasterInfo{}
		_, originalMachineObject := Get(dbConnection, ToolingMachineMasterTable, machine.Id)
		err := json.Unmarshal(originalMachineObject.ObjectInfo, &machineMasterInfo)
		hmiSettingInfo := getToolingHmiSetting(dbConnection, ms, machine.Id)

		if hmiSettingInfo == nil {
			continue
		}
		machineLiveDetectionInterval := hmiSettingInfo.MachineLiveDetectionInterval
		duration, _ := time.ParseDuration(machineLiveDetectionInterval)
		durationInSeconds := duration.Milliseconds()

		liveDetectionPeriod := timeNow - durationInSeconds

		if err != nil {
			ms.BaseService.Logger.Error("error in unmarshelling machine master", zap.String("error", err.Error()))
			continue
		}
		//Get past 15 seconds message
		cycleCountMessageQuery := "select * from message where ts > " + strconv.FormatInt(liveDetectionPeriod, 10) + " and topic = '" + "machines/" + machineMasterInfo.TopicFlag + "' order by ts desc"
		var messagesCycleCount []Message
		dbConnection.Raw(cycleCountMessageQuery).Scan(&messagesCycleCount)

		var connectionStatusCycleCount bool
		var currentCycleCount int

		if machineMasterInfo.TopicFlag == "L1_Makino_F3Graphite" {
			isOpCodeFound, opCodeValue := searchOPCode(messagesCycleCount)
			if isOpCodeFound {
				if opCodeValue {
					connectionStatusCycleCount = true
				} else {
					connectionStatusCycleCount = false
				}
			} else {
				if machineMasterInfo.MachineConnectStatus == 1 {
					connectionStatusCycleCount = true
				} else {
					connectionStatusCycleCount = false
				}
			}
		} else {
			connectionStatusCycleCount, currentCycleCount = findToolingCycleCountIncrement(messagesCycleCount)
		}

		delayTime, delayStatus, findMsg := findDelayMessage(messagesCycleCount)

		if findMsg {
			machineMasterInfo.DelayPeriod = delayTime
			machineMasterInfo.DelayStatus = delayStatus
		}

		machineMasterInfo.CurrentCycleCount = currentCycleCount

		fmt.Println("[message_query_for_cycle_count]:", cycleCountMessageQuery)

		ms.BaseService.Logger.Info("tooling calculation processing", zap.Any("machineMasterInfo", machineMasterInfo.NewMachineId))
		ms.BaseService.Logger.Info("tooling calculation processing", zap.Any("CurrentCycleCount", machineMasterInfo.CurrentCycleCount))
		if connectionStatusCycleCount {
			machineMasterInfo.MachineConnectStatus = 1
		} else {
			machineMasterInfo.MachineConnectStatus = 2
		}
		machineMasterInfo.LastUpdatedMachineLiveStatus = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		//machineObjectInfo, errMachines := json.Marshal(machineMasterInfo)
		//
		//if errMachines != nil {
		//	ms.BaseService.Logger.Error("error in marshelling machine master", "error", err.Error())
		//	continue
		//}
		//dbError := dbConnection.Model(&ToolingMachineMaster{}).Where("id = ?", machine.Id).Update("object_info", machineObjectInfo).Error
		updateQuery := "UPDATE tooling_machine_master SET object_info = JSON_SET(object_info, '$.machineConnectStatus'," + strconv.Itoa(machineMasterInfo.MachineConnectStatus) + ", '$.delayPeriod'," + strconv.FormatInt(machineMasterInfo.DelayPeriod, 10) + ", '$.delayStatus','" + machineMasterInfo.DelayStatus + "', '$.lastUpdatedMachineLiveStatus','" + machineMasterInfo.LastUpdatedMachineLiveStatus + "') WHERE id = ?"
		result := dbConnection.Exec(updateQuery, machine.Id)
		if result.Error != nil {
			ms.BaseService.Logger.Error("error updating machine master", zap.String("error", err.Error()))
		}

		//	Updating tooling machine view
		updates := map[string]interface{}{
			"delay_period":           strconv.FormatInt(delayTime, 10),
			"delay_status":           delayStatus,
			"current_cycle_count":    currentCycleCount,
			"machine_connect_status": machineMasterInfo.MachineConnectStatus,
		}
		err = ms.ViewManager.CreateOrUpdateToolingView(machine.Id, updates)
		if err != nil {
			ms.BaseService.Logger.Error("error updating tooling view", zap.String("error", err.Error()))
		}

	}

}

func findToolingCycleCountIncrement(messageList []Message) (bool, int) {
	machineStatus := false
	firstCycleCountFound := false
	firstCycleCountValue := 0

	for _, msg := range messageList {
		msgBody := make(map[string]interface{}, 0)

		json.Unmarshal(msg.Body, &msgBody)
		if valBody, ok := msgBody["MachiningTime"]; ok {
			cycleCount := util.InterfaceToInt(valBody)

			if firstCycleCountFound {
				if firstCycleCountValue > cycleCount {
					machineStatus = true
					break
				}
			}

			if !firstCycleCountFound {
				firstCycleCountFound = true
				firstCycleCountValue = cycleCount
			}

		}

	}

	return machineStatus, firstCycleCountValue
}

// Update given Machine time line with completed qty, progress percent and rejected quantity
func updateToolingMachineTimeLine(projectId string, eventId int, completedQty int, progressPercent float64, rejectedQty int, oee int) {
	stringTimeFormat := strconv.Itoa(completedQty) + "s"
	durationInTime, _ := time.ParseDuration(stringTimeFormat)

	var updatingFields = make(map[string]interface{})
	updatingFields["completedQty"] = fmt.Sprintf("%.2f", durationInTime.Hours()) + " hours"
	updatingFields["percentDone"] = progressPercent
	updatingFields["rejectedQty"] = rejectedQty
	updatingFields["oee"] = oee
	serializedObject, _ := json.Marshal(updatingFields)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateToolingScheduledOrderFields(projectId, eventId, serializedObject)

}

// Update given Machine time line with completed qty, progress percent and rejected quantity
func addToolingForceStopFlag(projectId string, eventId int) {

	var updatingFields = make(map[string]interface{})
	updatingFields["canForceStop"] = true

	serializedObject, _ := json.Marshal(updatingFields)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateToolingScheduledOrderFields(projectId, eventId, serializedObject)

}

func getToolingTopicName(machineId int, ms *MachineService, projectId string) string {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]
	machineQuery := "select * from tooling_machine_master where id = " + strconv.Itoa(machineId)
	var machineMaster ToolingMachineMaster
	err := dbConnection.Raw(machineQuery).Scan(&machineMaster).Error

	if err != nil {
		ms.BaseService.Logger.Error("error in fetching machine master", zap.String("error", err.Error()))
		return ""
	}

	var machineInfo ToolingMachineMasterInfo
	errMachineInfo := json.Unmarshal(machineMaster.ObjectInfo, &machineInfo)
	if errMachineInfo != nil {
		return ""
	}
	return "machines/" + machineInfo.TopicFlag
}

// Machine status is decided by hmi or cycle count of message
func getToolingMachineStatus(ms *MachineService, dbConnection *gorm.DB, machineId int) (string, string) {
	machineQuery := "select * from tooling_machine_master where id=" + strconv.Itoa(machineId)
	var machineMaster ToolingMachineMaster
	var machineInfo ToolingMachineMasterInfo
	dbConnection.Raw(machineQuery).Scan(&machineMaster)

	err := json.Unmarshal(machineMaster.ObjectInfo, &machineInfo)

	if err != nil {
		ms.BaseService.Logger.Error("error in unmarshalling machine master", zap.String("error", err.Error()))
		return "", ""
	}

	var machineConnectStatus map[string]interface{}
	_, connectStatusObject := Get(dbConnection, MachineConnectStatusTable, machineInfo.MachineConnectStatus)
	json.Unmarshal(connectStatusObject.ObjectInfo, &machineConnectStatus)

	statusName := util.InterfaceToString(machineConnectStatus["status"])

	return statusName, machineInfo.NewMachineId
}

func getToolingDowntime(hmiInfoList []MachineHMIInfoResult, dbConnection *gorm.DB, machineId int, startEpochMilli int64, endEpochMilli int64) int64 {
	//DownTime in Milliseconds
	resultLength := len(hmiInfoList)
	var downTime int64 = 0

	for index, hmi := range hmiInfoList {
		if hmi.Status == "stopped" {
			if (index + 1) < resultLength {
				downTime = downTime + (hmiInfoList[index+1].TS.DateTimeEpochMilli - hmi.TS.DateTimeEpochMilli)
			} else {
				downTime = downTime + (time.Now().UTC().UnixMilli() - hmi.TS.DateTimeEpochMilli)
			}

		}
	}

	machineLiveUnplannedDown := getToolingMachineLiveUnplanndedDownTime(dbConnection, machineId, startEpochMilli, endEpochMilli)
	return downTime + machineLiveUnplannedDown
}

// Should be in Milliseconds
func getToolingMachineLiveUnplanndedDownTime(dbConnection *gorm.DB, machineId int, startEpochMilli int64, endEpochMilli int64) int64 {
	//Get machine statistics based on the machine id and epochs
	// machine status == !Live
	// list stop_start list
	//While loop
	//	if row satisfy the condition machine status
	//  	machine status == Live
	//      append stop_start_list
	startTimeInSeconds := int64(startEpochMilli / 1000)
	endTimeInSeconds := int64(endEpochMilli / 1000)
	statsQuery := "select * from tooling_machine_statistics where machine_id=" + strconv.Itoa(machineId) + " and ts >= " + strconv.FormatInt(startTimeInSeconds, 10) + " and ts <= " + strconv.FormatInt(endTimeInSeconds, 10) + " order by ts asc"
	var machineStats []ToolingMachineStatistics
	var machineStatsInfo MachineStatisticsInfo
	dbConnection.Raw(statsQuery).Scan(&machineStats)
	fmt.Println("statsQuery: ", statsQuery)
	fmt.Println("result size:", len(machineStats))
	var machineConnectStatus map[string]interface{}
	_, connectStatusObject := Get(dbConnection, MachineConnectStatusTable, machineConnectStatusWaitingForFeed)
	json.Unmarshal(connectStatusObject.ObjectInfo, &machineConnectStatus)

	searchingStatus := util.InterfaceToString(machineConnectStatus["status"])
	stopStartTS := make([]int64, 0)
	var machineDownTime int64
	machineDownTime = 0
	// Live WF WF Live WF Live live live live WF WF
	// t1    2 3   4    5  6    7    8    9   10 11
	// 4- t1 + 6 -4 +  11 -10
	// Live Live Live WF WF WF WF WF Live
	//                 t1            -
	//Create start and stop list by iterating through the machine statistics
	for index, statistics := range machineStats {
		err := json.Unmarshal([]byte(statistics.StatsInfo), &machineStatsInfo)

		if err != nil {
			continue
		}

		if machineStatsInfo.CurrentStatus == searchingStatus {
			stopStartTS = append(stopStartTS, statistics.TS)
			if index == 0 {
				// if the first is waiting for feed, then we need to skip it
				continue
			}
			err = json.Unmarshal(machineStats[index-1].StatsInfo, &machineStatsInfo)

			if err != nil {
				continue
			}
			searchingStatus = machineStatsInfo.CurrentStatus
		}
	}

	lenOfStopTS := len(stopStartTS)
	if lenOfStopTS == 0 {
		return machineDownTime
	}
	// L  L  L W W W W L L L
	// [ 5    10   30  ]
	// 10-5  + unix -30
	if lenOfStopTS%2 != 0 {
		// we got the odd values, so lets make it even for loop
		stopStartTS = append(stopStartTS, time.Now().Unix())
	}

	for index := 0; index < lenOfStopTS; index += 2 {
		machineDownTime += stopStartTS[index+1] - stopStartTS[index]
	}

	return machineDownTime * 1000

}

// TODO based on all the events belongs to that proudction order A = childs
func getToolingOverallRejectedQuantity(ms *MachineService, dbConnection *gorm.DB, productionOrderId int) int {

	overallRejectedQty := 0
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	listOfScheduledOrderEvents := productionOrderInterface.GetChildEventsOfToolingProductionOrder(ProjectID, productionOrderId)
	for _, eventId := range listOfScheduledOrderEvents {
		//Decode Machine timeline event info

		conditionString := " object_info ->> '$.eventId' = " + strconv.Itoa(eventId)
		hmiObject, _ := GetConditionalObjects(dbConnection, ToolingMachineHmiTable, conditionString)

		for _, machineHmi := range *hmiObject {
			machineHMI := MachineHMI{ObjectInfo: machineHmi.ObjectInfo}

			rejectedQty := machineHMI.getMachineHMIInfo().RejectedQuantity
			if rejectedQty != nil {
				overallRejectedQty += *rejectedQty
			}

		}

	}

	return overallRejectedQty

}

func getFirstToolingCycleCount(messages []Message) (int, ToolingMessage) {
	var cycleCount int
	var cycleCountMsg ToolingMessage
	// we need to basically count total numbers excluding 0
	//var cycleCountArray []int
	for i := len(messages) - 1; i >= 0; i-- {
		var message = messages[i]
		var messageBody = ToolingMessage{}
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil || messageBody.MachiningTime == nil {
			continue
		} else {
			cycleCount = *messageBody.MachiningTime
			break
		}
	}

	// This is for finding last message
	for _, message := range messages {
		var messageBody = ToolingMessage{}
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil || messageBody.MachiningTime == nil {
			continue
		} else {
			cycleCountMsg = messageBody
			break
		}
	}

	return cycleCount, cycleCountMsg
}

func getToolingCycleCount(messages []Message) (int, string) {
	// first get the last Machining Time - start Machining Time
	var messageBody ToolingMessage

	getLastMachiningTimeFlag := false
	var lastMachiningTime int
	var firstMachiningTime int
	// we need to basically count total numbers excluding 0
	//var cycleCountArray []int
	for _, message := range messages {
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil || messageBody.MachiningTime == nil {
			continue
		}

		if !getLastMachiningTimeFlag {
			lastMachiningTime = *messageBody.MachiningTime
			getLastMachiningTimeFlag = true
		}

		firstMachiningTime = *messageBody.MachiningTime
	}

	machiningTime := lastMachiningTime - firstMachiningTime
	remark := ""
	if machiningTime < 0 {
		machiningTime = 0
		remark = "Cycle count has negative value"
	}

	return machiningTime, remark
}

func getToolingCycleCountBeforeOneMin(messages []Message) int {
	beforeOneMinTime := time.Now().UTC().UnixMilli() - 60000

	breakPointer := len(messages) - 1
	for index, message := range messages {
		messageBody := ToolingMessage{}
		err := json.Unmarshal(message.Body, &messageBody)

		if err == nil && messageBody.InTimestamp != nil {
			if int64(*messageBody.InTimestamp) < beforeOneMinTime {
				breakPointer = index
				break
			}
		}
	}

	oneMinBeforeMessages := messages[breakPointer:]
	cycleCount, _ := getToolingCycleCount(oneMinBeforeMessages)

	return cycleCount

}

func getToolingCompletedValue(currentEtlId int, currentActual int, productionOrder string, machineId int, timeLineEvents *[]component.GeneralObject, ms *MachineService, projectId string) int {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]
	completedValue := currentActual

	for _, eventTimeLine := range *timeLineEvents {

		if eventTimeLine.Id == currentEtlId {
			continue
		}
		etlInfo := make(map[string]interface{})
		err := json.Unmarshal(eventTimeLine.ObjectInfo, &etlInfo)
		if err != nil {
			continue
		}

		startDate := convertStringToDateTime(util.InterfaceToString(etlInfo["startDate"]))
		endDate := convertStringToDateTime(util.InterfaceToString(etlInfo["endDate"]))

		if startDate.Error != nil || endDate.Error != nil {
			continue
		}
		if util.InterfaceToString(etlInfo["productionOrder"]) == productionOrder {
			statisticsResult := ToolingMachineStatistics{}
			statsInfo := MachineStatisticsInfo{}
			statQuery := "select * from tooling_machine_statistics where stats_info ->> '$.eventId'= " + strconv.Itoa(eventTimeLine.Id) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"

			dbConnection.Raw(statQuery).Scan(&statisticsResult)

			errStats := json.Unmarshal(statisticsResult.StatsInfo, &statsInfo)
			if errStats != nil {
				continue
			}
			completedValue += statsInfo.Actual

		}

	}

	return completedValue
}

type ToolingMessage struct {
	AlarmMsg      *[]string `json:"AlarmMsg"`
	WarningMsg    *[]int    `json:"WarningMsg"`
	ActFeedRate   *float32  `json:"ActFeedRate"`
	AutoRunTime   *int      `json:"AutoRunTime"`
	InTimestamp   *int      `json:"in_timestamp"`
	MachiningTime *int      `json:"MachiningTime"`
	RemainingTime *int      `json:"RemainingTime"`
}
