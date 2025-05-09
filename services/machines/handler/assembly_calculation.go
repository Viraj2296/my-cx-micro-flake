package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

/*
   run_time = planned production time - stop time
            = 420 min - 60 = 373 minutes

good count = Total count - reject count
           =  19271 - 423 = 18,848

availability = run time / planned production time
= 373 / 420 = 88 %

performance = ideal cycle time * total count / run time
            = 1 seconds * 19271/ 373*60 = 86.11

Quality = good count / total count
          18848/ 19271 = 97 %
*/
/*
  Daily Planned Quantity (scheduled quantity from machine time line event) : 15026
Overall Planned Quantity : 200000  - done
Daily completed Quantity : Actual   - ok
Overall Completed Quantity : 16000 - done
Daily Reject Quantity :  reject quantity  - 90    - ok
Overall rejected Quanatity : Accumuldated one  - yesterday (200) + current rejected so far (90) = total = 290  - ok
Daily completed percentage : Actual/Daily planned :  300/15026 * 100  - done
Overall completed percentage : Completed/ overall planned quantity :16900/200000 * 100 = ? - ok
*/

func (ms *MachineService) MakeAssemblyMachineStatisticsCalculation() {
	//Get all data from Machine timeline event table
	dbConnection := ms.BaseService.ServiceDatabases[ProjectID]

	//Check whether we are getting messages or not. If yes set it as Live
	ms.updateAssemblyMachineConnectStatus(dbConnection)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	_, listOfEvents := productionOrderInterface.GetAssemblyScheduledEvents(ProjectID)
	ms.BaseService.Logger.Info("events processing", zap.Any("number_of_events", len(*listOfEvents)))

	// Get production order complete status id from preference level
	orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(ProjectID, ScheduleStatusPreferenceSeven)

	for _, eventTimeLine := range *listOfEvents {
		//Decode Machine timeline event info
		var etlInfo ScheduledOrderEventInfo
		err := json.Unmarshal(eventTimeLine.ObjectInfo, &etlInfo)
		if err == nil {

			machineStatus, newMachineId := getAssemblyMachineStatus(ms, dbConnection, etlInfo.MachineId)
			// if the machine is not defined in the configuration, then we don't use it for calculation
			assemblyMachineMasterInfo := GetAssemblyMachineFromId(dbConnection, etlInfo.MachineId)
			if assemblyMachineMasterInfo != nil {
				isLineContainArea := ms.isAreaAvailableInLine(assemblyMachineMasterInfo.MessageFlag, assemblyMachineMasterInfo.Area)
				if !isLineContainArea {
					ms.BaseService.Logger.Warn("machine is not configured to do the calculation, skipping it", zap.Int("machine_id", etlInfo.MachineId))
					continue
				}
			} else {
				ms.BaseService.Logger.Error("invalid machine ID skipping", zap.Int("machine_id", etlInfo.MachineId))
				continue
			}

			//Read schedule start date and end date
			startDate := convertStringToDateTime(etlInfo.StartDate)
			endDate := convertStringToDateTime(etlInfo.EndDate)
			if startDate.Error == nil || endDate.Error == nil {
				//Check end date is overdue or not
				//Adding singapore time in utc time
				// event time is not passed in under
				// Event status in completed then we don't do calculation
				if time.Now().Before(endDate.DateTime) && etlInfo.EventStatus != orderStatusId {
					startEpoch := startDate.DateTimeEpochMilli
					endEpoch := endDate.DateTimeEpochMilli
					topicName, messageFlag := getAssemblyTopicName(etlInfo.MachineId, ms, ProjectID)

					hmiInfoList := getHMIInfo(ms, ProjectID, etlInfo.MachineId, eventTimeLine.Id, AssemblyMachineHmiTable)
					ms.BaseService.Logger.Info("HMI info list", zap.Any("machine_id:", etlInfo.MachineId), zap.Any("list:", hmiInfoList))
					if len(hmiInfoList) > 0 {
						//Get machine reset assembly object
						var startHmi = getResetHmiFlag(hmiInfoList)
						if startHmi.EventId != eventTimeLine.Id {
							startHmi = hmiInfoList[0]
						}

						// getting messages based on HMI first created time and current time for processing
						// we needed only the in_timestamp message greater
						messageQuery := "select * from message where body->>'$.in_timestamp' >= " + strconv.FormatInt(convertStringToDateTime(startHmi.Created).DateTimeEpochMilli, 10) + " and body->>'$.in_timestamp' <= " + strconv.FormatInt(time.Now().UnixMilli(), 10) + " and topic='" + topicName + "' order by body->>'$.in_timestamp' desc"
						ms.BaseService.Logger.Info("running query to select the messages", zap.Any("query", messageQuery))
						var messages []Message
						dbConnection.Raw(messageQuery).Scan(&messages)
						//messageOffSet = getCycleCountOffset(messages)
						//ms.BaseService.Logger.Infow("message offset", "offset", messageOffSet)
						// Fetch data based on the schedule from message table

						if len(messages) > 0 {

							err, productionOrderObject := productionOrderInterface.GetAssemblyProductionOrderInfo(ProjectID, etlInfo.EventSourceId, etlInfo.MachineId)
							if err != nil {
								ms.BaseService.Logger.Error("error getting production order info", zap.Any("event_source_id", etlInfo.EventSourceId), zap.Any("machine_id", etlInfo.MachineId))
								continue
							}
							productionOrderInfo := GetProductionOrderInfo(productionOrderObject.ObjectInfo)
							ms.BaseService.Logger.Info("production order:", zap.Any("machine_id:", etlInfo.MachineId), zap.Any("production order info ", productionOrderInfo))

							// accumulatedCycleCount := getAccumulatedCycleCount(messages)
							// how much manufactured so far from PUMAS

							cycleCount, remark := getAssemblyCycleCount(messages, messageFlag)
							ms.BaseService.Logger.Info("cycle count", zap.Any("machine_id:", etlInfo.MachineId), zap.Any("cycleCount", cycleCount))

							// If the manual reset happened, then we need substract the cycle count from off set
							var resetCycleCount = getResetOffsetAssemblyCycleCount(dbConnection, hmiInfoList, eventTimeLine.Id, messageFlag, topicName)
							if resetCycleCount != 0 {
								cycleCount = cycleCount - resetCycleCount
							}
							//=========================================================================================================================
							productCount := cycleCount
							ms.BaseService.Logger.Info("product count ", zap.Any("machine_id:", etlInfo.MachineId), zap.Any("product count:", productCount), zap.Any("offset value for reset:", resetCycleCount))

							//TODO, get the cavity using interface function
							// actualProduced :=cycleCount * noOfCavity

							//Get all hmis for given machine and event

							//based on UI or HMI
							dailyRejectedCount := getTotalRejectQtyForEvent(hmiInfoList)

							overallRejectedCount := getAssemblyOverallRejectedQuantity(ms, dbConnection, etlInfo.EventSourceId)
							//getOverallRejectedQuantity(ms, dbConnection, etlInfo.ProductionOrder, eventTimeLines)
							goodCount := getGoodCount(productCount, dailyRejectedCount)
							machineDownTime := getAssemblyDowntime(hmiInfoList, dbConnection, etlInfo.MachineId, startEpoch, endEpoch)

							plannedProductionTime := getPlannedProductionTime(hmiInfoList) // this is in milliseconds
							// plannedProductionTime := endEpoch - startEpoch
							availability := getAvailability(plannedProductionTime, machineDownTime)

							//capacity := getCapacity(productionOrder.CycleTime, startEpoch, machineDownTime)
							//performance := getPerformance(goodCount, capacity)

							runTime := plannedProductionTime - machineDownTime

							performance := (float64(productionOrderInfo.CycleTime) * float64(productCount)) / float64(runTime/1000)
							var quality float64
							quality = 0.0
							if productCount > 0 {
								quality = float64(goodCount) / float64(productCount)
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

							_, listOfEventForProduction := productionOrderInterface.GetAssemblyScheduledEventByProductionId(ProjectID, etlInfo.EventSourceId)

							actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, etlInfo.MachineId, eventTimeLine.Id, AssemblyMachineHmiTable)
							completedValue := getAssemblyCompletedValue(eventTimeLine.Id, productCount, etlInfo.ProductionOrder, etlInfo.MachineId, listOfEventForProduction, ms, ProjectID)
							progressPercentage := getProgressPercentage(productCount, etlInfo.ScheduledQty)

							ms.BaseService.Logger.Info("assembly calculation summary", zap.Any("machine_id", etlInfo.MachineId), zap.Any("event_id", eventTimeLine.Id),
								zap.Any("planned_production_time", plannedProductionTime), zap.Any("machine_down_time", machineDownTime), zap.Any("performance", performance),
								zap.Any("availability", availability), zap.Any("completedValue", completedValue), zap.Any("progressPercentage", progressPercentage),
								zap.Any("actualStartTime", actualStartTime), zap.Any("actualEndTime", actualEndTime), zap.Any("quality", quality))
							statsInfo := MachineStatisticsInfo{
								EventId:                    eventTimeLine.Id,
								ProductionOrderId:          etlInfo.EventSourceId,
								CurrentStatus:              machineStatus,
								PartId:                     productionOrderInfo.PartNumber,
								ScheduleStartTime:          etlInfo.StartDate,
								ScheduleEndTime:            etlInfo.EndDate,
								ActualStartTime:            actualStartTime,
								ActualEndTime:              actualEndTime,
								EstimatedEndTime:           getEstimatedEndTime(float32(productionOrderInfo.CycleTime), productCount, actualStartTime, etlInfo.ScheduledQty),
								Oee:                        int(oee * 10000),
								Availability:               int(availability * 10000),
								Performance:                int(performance * 10000),
								Quality:                    int(quality * 10000),
								PlannedQuality:             productionOrderInfo.OrderQty,
								DailyPlannedQty:            etlInfo.ScheduledQty,
								Completed:                  completedValue,
								Rejects:                    dailyRejectedCount,
								OverallRejectedQty:         overallRejectedCount,
								CompletedPercentage:        getDailyCompletedPercentage(etlInfo.ScheduledQty, productCount),
								OverallCompletedPercentage: getOverallCompletedPercentage(productionOrderInfo.OrderQty, completedValue),
								Actual:                     productCount,
								ProgressPercentage:         progressPercentage,
								DownTime:                   int(machineDownTime),
								WarningMessage:             make([]string, 0),
								Remark:                     remark,
							}
							statsInfoJson, _ := json.Marshal(statsInfo)
							machineStatistics := AssemblyMachineStatistics{
								MachineId: etlInfo.MachineId,
								TS:        time.Now().Unix(),
								StatsInfo: statsInfoJson,
							}
							err = dbConnection.Create(&machineStatistics).Error
							if err != nil {
								ms.BaseService.Logger.Error("inserting stats has failed", zap.String("error", err.Error()))
							}
							updateAssemblyMachineTimeLine(ProjectID, eventTimeLine.Id, productCount, progressPercentage, dailyRejectedCount, int(oee*10000))
							setHelpStopInMachineStop(messages, messageFlag, etlInfo, eventTimeLine.Id, hmiInfoList, dbConnection, ms)

							//	Need to update view columns
							initialCycleCount, cycleCountMessage := getFirstAssemblyCycleCount(messages, messageFlag)
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

							err = ms.ViewManager.CreateOrUpdateAssemblyView(etlInfo.MachineId, updates)
							if err != nil {
								ms.BaseService.Logger.Error("error updating assembly view", zap.String("error", err.Error()))
							}
						}
					} else {
						continue
					}
				} else {
					hmiInfoList := getHMIInfo(ms, ProjectID, etlInfo.MachineId, eventTimeLine.Id, AssemblyMachineHmiTable)

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
						hmiSettingInfo := getAssemblyHmiSetting(dbConnection, ms, etlInfo.MachineId)
						ms.BaseService.Logger.Info("hmi setting info", zap.Any("hmi_setting_info", hmiSettingInfo), zap.Any("endDate:", endDate))
						//End date should be in seconds
						appendWarningAssemblyMessage(ms, dbConnection, hmiSettingInfo.WarningMessageGenerationPeriod, endDate.DateTimeEpoch, eventTimeLine.Id, etlInfo.MachineId, etlInfo.Name, newMachineId, hmiSettingInfo.WarningTargetEmailId)
						addForceStopAssemblyScheduler(ProjectID, eventTimeLine.Id)
						//If user not stop HMI data, Then program automatically stopped hmi
						// Condition
						// Based on hmi stop configuration hmi will be stopped
						automaticAssemblyHmiStop(hmiSettingInfo.HmiAutoStopPeriod, lastHMIInfoResult, dbConnection, ms, endDate.DateTimeEpoch)
					}

				}
			}

		}
	}
}

func getResetOffsetAssemblyCycleCount(dbConnection *gorm.DB, hmiInfoList []MachineHMIInfoResult, eventId int, cycleCountKey string, topicName string) int {
	var cycleCount = 0
	var resetHmi = getResetHmiFlag(hmiInfoList)
	if resetHmi.EventId != eventId {
		return cycleCount
	}
	var startHmi = hmiInfoList[0]
	var messageQuery = "select * from message where body->>'$.in_timestamp' >= " + strconv.FormatInt(convertStringToDateTime(startHmi.Created).DateTimeEpochMilli, 10) + " and body->>'$.in_timestamp' <= " + strconv.FormatInt(convertStringToDateTime(resetHmi.Created).DateTimeEpochMilli, 10) + " and topic='" + topicName + "'"
	var messages []Message
	dbConnection.Raw(messageQuery).Scan(&messages)
	for _, message := range messages {
		messageBody := make(map[string]interface{})
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if val, ok := messageBody[cycleCountKey]; ok {
			cycleCount = util.InterfaceToInt(val)
			break
		} else {
			continue
		}
	}

	return cycleCount
}

func setHelpStopInMachineStop(messages []Message, messageFlag string,
	etlInfo ScheduledOrderEventInfo,
	eventId int,
	hmiInfoList []MachineHMIInfoResult,
	dbConnection *gorm.DB, ms *MachineService) {

	messageKey := messageFlag + "_orderNum"
	stnMessageKey := messageFlag + "_stnCall"
	orderMessageBody := getOrderNameInMessages(messages, messageKey)

	// This message in reverse order(asc)
	stnCallMessages := getStnCallMessages(messages, stnMessageKey)

	if len(orderMessageBody) != 0 {
		scheduledOrderName := util.InterfaceToString(orderMessageBody[messageKey])

		for _, helpMsg := range stnCallMessages {
			stnCall := util.InterfaceToString(helpMsg[stnMessageKey])
			insertedTs := util.InterfaceToInt(helpMsg["in_timestamp"])
			containsOne := strings.Contains(stnCall, "1")

			if containsOne {
				if scheduledOrderName == etlInfo.Name {
					if !ifAlreadyStopInserted(hmiInfoList, int64(insertedTs), HelpStopReasonId) {
						hmiInfo := MachineHMIInfoResult{ReasonId: 1, MachineId: etlInfo.MachineId, EventId: eventId, HmiStatus: "stopped"}
						updateOverDueAssemblyHmiInfo(hmiInfo, dbConnection, ms, "HMI was stopped because operator clicked help button")
					}
				}

			}
		}
	}
}

func ifAlreadyStopInserted(hmiInfoList []MachineHMIInfoResult, insertedTs int64, reasonId int) bool {
	for _, hmiInfo := range hmiInfoList {
		createdTime := util.ConvertStringToDateTime(hmiInfo.Created)
		if hmiInfo.HmiStatus == "stopped" && createdTime.DateTimeEpochMilli >= insertedTs && hmiInfo.ReasonId == reasonId {
			return true
		}
	}

	return false
}

func getStnCallMessages(messages []Message, stnMessageKey string) []map[string]interface{} {
	helpMessageBody := make([]map[string]interface{}, 0)

	// Loop over the slice in reverse
	for i := len(messages) - 1; i >= 0; i-- {
		msgBody := make(map[string]interface{})
		json.Unmarshal(messages[i].Body, &msgBody)

		if _, ok := msgBody[stnMessageKey]; ok {
			helpMessageBody = append(helpMessageBody, msgBody)
		}
	}

	return helpMessageBody
}

func getOrderNameInMessages(messages []Message, messageKey string) map[string]interface{} {
	orderMessageBody := make(map[string]interface{})
	for _, msg := range messages {
		msgBody := make(map[string]interface{})
		json.Unmarshal(msg.Body, &msgBody)

		if _, ok := msgBody[messageKey]; ok {
			return msgBody
		}
	}

	return orderMessageBody
}

func automaticAssemblyHmiStop(hmiAutoStopPeriod string, info MachineHMIInfoResult, dbConnection *gorm.DB, ms *MachineService, scheduleStopTime int64) {
	//Schedule stop time in seconds
	duration, _ := time.ParseDuration(hmiAutoStopPeriod)
	durationInSeconds := duration.Seconds()
	if time.Now().UTC().Unix() > (scheduleStopTime + int64(durationInSeconds)) {
		fmt.Println("is it exceeeded: ", time.Now().UTC().Unix(), " scheduleStopTime :", scheduleStopTime, " duration seconds :", int64(durationInSeconds))
		updateOverDueAssemblyHmiInfo(info, dbConnection, ms, "HMI was stopped by program")
	}
}

func getAssemblyHmiSetting(dbConnection *gorm.DB, ms *MachineService, machineId int) *HmiSettingInfo {
	err, generalObject := Get(dbConnection, AssemblyMachineHmiSettingTable, machineId)
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

// If the user not stopped the hmi info, program automatically is inserted stopped hmi
// This function only update created , status and remark attributes
func updateOverDueAssemblyHmiInfo(info MachineHMIInfoResult, dbConnection *gorm.DB, ms *MachineService, remark string) {
	objectInfo := MachineHMIInfo{
		CreatedAt: time.Now().Format("2006-01-02T15:04:05.000Z"),
		EventId:   info.EventId,
		Operator:  info.Operator,
		MachineId: info.MachineId,
		Status:    info.Status,
		HMIStatus: "stopped",
		Remark:    remark,
	}
	hmiObjectInfo, err := json.Marshal(objectInfo)

	if err != nil {
		ms.BaseService.Logger.Error("error in marshalling hmi info", zap.String("error", err.Error()))
	}

	dbError, _ := Create(dbConnection, AssemblyMachineHmiTable, component.GeneralObject{ObjectInfo: hmiObjectInfo})

	if dbError != nil {
		ms.BaseService.Logger.Error("error in updating hmi info", zap.String("error", err.Error()))
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateAssemblyOrderPreferenceLevel(ProjectID, 2, info.EventId, ScheduleStatusPreferenceSix)

}

// Update machine master whether we receive message or not. Live or Waiting for feed
func (ms *MachineService) updateAssemblyMachineConnectStatus(dbConnection *gorm.DB) {
	// year, month, day := time.Now().Date()
	timeNow := time.Now().UTC().UnixMilli()
	generalObjects, _ := GetObjects(dbConnection, AssemblyMachineMasterTable)
	for _, assemblyMachine := range *generalObjects {
		assemblyMachineMasterInfo := AssemblyMachineMasterInfo{}
		_, originalMachineObject := Get(dbConnection, AssemblyMachineMasterTable, assemblyMachine.Id)
		err := json.Unmarshal(originalMachineObject.ObjectInfo, &assemblyMachineMasterInfo)
		hmiSettingInfo := getAssemblyHmiSetting(dbConnection, ms, assemblyMachine.Id)

		if hmiSettingInfo == nil {
			continue
		}
		isLineContainArea := ms.isAreaAvailableInLine(assemblyMachineMasterInfo.MessageFlag, assemblyMachineMasterInfo.Area)

		machineLiveDetectionInterval := hmiSettingInfo.MachineLiveDetectionInterval
		duration, _ := time.ParseDuration(machineLiveDetectionInterval)
		durationInSeconds := duration.Milliseconds()

		liveDetectionPeriod := timeNow - durationInSeconds

		if err != nil {
			ms.BaseService.Logger.Error("error in unmarshalling machine master", zap.String("error", err.Error()))
			continue
		}
		//Get past 15 seconds message
		cycleCountMessageQuery := "select * from message where ts > " + strconv.FormatInt(liveDetectionPeriod, 10) + " and topic = '" + "machines/" + assemblyMachineMasterInfo.Level + "' order by ts desc"
		ms.BaseService.Logger.Info("processing the machine status", zap.String("machine_id", assemblyMachineMasterInfo.NewMachineId), zap.String("query", cycleCountMessageQuery))
		var messagesCycleCount []Message
		dbConnection.Raw(cycleCountMessageQuery).Scan(&messagesCycleCount)

		connectionStatusCycleCount, currentCycleCount := findAssemblyCycleCountIncrement(messagesCycleCount, assemblyMachineMasterInfo.MessageFlag)
		delayTime, delayStatus, findMsg := findDelayMessage(messagesCycleCount)
		if findMsg {
			assemblyMachineMasterInfo.DelayPeriod = delayTime
			assemblyMachineMasterInfo.DelayStatus = delayStatus
		}

		assemblyMachineMasterInfo.CurrentCycleCount = currentCycleCount
		ms.BaseService.Logger.Info("calculation processing", zap.Any("AssemblyMessageFlag", assemblyMachineMasterInfo.MessageFlag))
		ms.BaseService.Logger.Info("calculation processing", zap.Any("CurrentCycleCount", assemblyMachineMasterInfo.CurrentCycleCount), zap.Any("connectionStatusCycleCount", connectionStatusCycleCount), zap.Any("isLineContainArea", isLineContainArea))
		var machineConnectStatus = 2
		if connectionStatusCycleCount {
			ms.BaseService.Logger.Info("connectionStatusCycleCount", zap.Any("connectionStatusCycleCount", connectionStatusCycleCount))
			assemblyMachineMasterInfo.MachineConnectStatus = 1
			machineConnectStatus = 1
			if !isLineContainArea {
				assemblyMachineMasterInfo.MachineConnectStatus = 2
				machineConnectStatus = 2
			}

		} else {
			assemblyMachineMasterInfo.MachineConnectStatus = 2
			machineConnectStatus = 2
		}
		assemblyMachineMasterInfo.LastUpdatedMachineLiveStatus = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

		updateQuery := "UPDATE assembly_machine_master SET object_info = JSON_SET(object_info, '$.machineConnectStatus'," + strconv.Itoa(assemblyMachineMasterInfo.MachineConnectStatus) + ", '$.delayPeriod'," + strconv.FormatInt(assemblyMachineMasterInfo.DelayPeriod, 10) + ", '$.delayStatus','" + assemblyMachineMasterInfo.DelayStatus + "', '$.lastUpdatedMachineLiveStatus','" + assemblyMachineMasterInfo.LastUpdatedMachineLiveStatus + "', '$.currentCycleCount'," + strconv.Itoa(assemblyMachineMasterInfo.CurrentCycleCount) + ") WHERE id = ?"
		result := dbConnection.Exec(updateQuery, assemblyMachine.Id)
		if result.Error != nil {
			ms.BaseService.Logger.Error("error updating machine master", zap.String("error", err.Error()))
		}
		updates := map[string]interface{}{
			"delay_period":           strconv.FormatInt(delayTime, 10),
			"delay_status":           delayStatus,
			"current_cycle_count":    currentCycleCount,
			"machine_connect_status": machineConnectStatus,
		}
		err = ms.ViewManager.CreateOrUpdateAssemblyView(assemblyMachine.Id, updates)
		if err != nil {
			ms.BaseService.Logger.Error("error updating assembly view", zap.String("error", err.Error()))
		}
	}

}

func getResetHmiFlag(listOfMachineHmi []MachineHMIInfoResult) MachineHMIInfoResult {
	// Method 1: Using a for loop with index
	for i := len(listOfMachineHmi) - 1; i >= 0; i-- {
		var machineHmi = listOfMachineHmi[i]
		if machineHmi.Remark == "reset" {
			return machineHmi
		}
	}
	return MachineHMIInfoResult{}
}

func (ms *MachineService) isAreaAvailableInLine(MachineShortCode string, area string) bool {

	for _, configurations := range ms.AssemblyMachineConfiguration {
		if configurations.MessageFlag == MachineShortCode {
			for _, areaConfig := range configurations.Area {
				if area == areaConfig {
					return true
				}
			}
		}
	}

	return false
}

func findAssemblyCycleCountIncrement(messageList []Message, cycleCountKey string) (bool, int) {
	machineStatus := false
	firstCycleCountFound := false
	firstCycleCountValue := 0

	for _, msg := range messageList {
		msgBody := make(map[string]interface{}, 0)

		json.Unmarshal(msg.Body, &msgBody)

		if valBody, ok := msgBody[cycleCountKey]; ok {
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
func updateAssemblyMachineTimeLine(projectId string, eventId int, completedQty int, progressPercent float64, rejectedQty int, oee int) {

	var updatingFields = make(map[string]interface{})
	updatingFields["completedQty"] = completedQty
	updatingFields["percentDone"] = progressPercent
	updatingFields["rejectedQty"] = rejectedQty
	updatingFields["oee"] = oee
	serializedObject, _ := json.Marshal(updatingFields)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, eventId, serializedObject)

}

func addForceStopAssemblyScheduler(projectId string, eventId int) {

	var updatingFields = make(map[string]interface{})
	updatingFields["canForceStop"] = true

	serializedObject, _ := json.Marshal(updatingFields)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, eventId, serializedObject)

}

func getAssemblyTopicName(machineId int, ms *MachineService, projectId string) (string, string) {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]
	err, machineMaster := Get(dbConnection, AssemblyMachineMasterTable, machineId)

	if err != nil {
		ms.BaseService.Logger.Error("error in fetching machine master", zap.String("error", err.Error()))
		return "", ""
	}

	var assemblyMachineInfo AssemblyMachineMasterInfo
	errMachineInfo := json.Unmarshal(machineMaster.ObjectInfo, &assemblyMachineInfo)
	if errMachineInfo != nil {
		return "", ""
	}
	return "machines/" + assemblyMachineInfo.Level, assemblyMachineInfo.MessageFlag
}

// Machine status is decided by hmi or cycle count of message
func getAssemblyMachineStatus(ms *MachineService, dbConnection *gorm.DB, machineId int) (string, string) {
	err, machineMaster := Get(dbConnection, AssemblyMachineMasterTable, machineId)
	var machineInfo AssemblyMachineMasterInfo

	err = json.Unmarshal(machineMaster.ObjectInfo, &machineInfo)

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

func getFirstAssemblyCycleCount(messages []Message, cycleCountKey string) (int, map[string]interface{}) {
	var cycleCount int
	var cycleCountMessage map[string]interface{}

	// This is for finding first message
	for i := len(messages) - 1; i >= 0; i-- {
		var message = messages[i]
		messageBody := make(map[string]interface{})
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if val, ok := messageBody[cycleCountKey]; ok {
			cycleCount = util.InterfaceToInt(val)
			cycleCountMessage = messageBody
			break
		} else {
			continue
		}
	}

	// This is for finding last message
	for _, message := range messages {
		messageBody := make(map[string]interface{})
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if _, ok := messageBody[cycleCountKey]; ok {
			cycleCountMessage = messageBody
			break
		} else {
			continue
		}
	}

	return cycleCount, cycleCountMessage
}

func getAssemblyCycleCount(messages []Message, cycleCountKey string) (int, string) {
	var cycleCountArray []int
	var isZeroDetected bool
	isZeroDetected = false
	for _, message := range messages {
		messageBody := make(map[string]interface{})
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if val, ok := messageBody[cycleCountKey]; ok {
			cycleCount := util.InterfaceToInt(val)
			cycleCountArray = append(cycleCountArray, cycleCount)
			if cycleCount == 0 {
				isZeroDetected = true
			}
		} else {
			continue
		}
	}
	if !isZeroDetected { // no zero value detected, last value number is the cycle count
		if len(cycleCountArray) == 0 {
			return 0, ""
		} else {
			return cycleCountArray[0], ""
		}
	}
	// we need to check whether the part id is changing between 0 cycle count, if it is changing, then we shouldn't add between 0
	//TODO
	cycleCount := sumBetweenZero(cycleCountArray)
	var remark string
	if cycleCount < 0 {
		cycleCount = 0
		remark = "Cycle count has negative value"
	}
	return cycleCount, remark
}

func sumBetweenZero(numbers []int) int {
	var segmentMax int
	var maxSum int
	for _, num := range numbers {
		if num == 0 {
			if segmentMax != 0 {
				maxSum += segmentMax
			}
			segmentMax = 0
		} else if num > segmentMax {
			segmentMax = num
		}
	}
	if segmentMax != 0 {
		maxSum += segmentMax
	}
	return maxSum
}

func appendWarningAssemblyMessage(ms *MachineService, dbConnection *gorm.DB, warningMessageGenrationPeriod string, scheduleStopTime int64, eventId int, machineId int, eventName string, machineName string, targetEmailId string) {
	// Schedule stop time in seconds
	statisticsQuery := "select * from assembly_machine_statistics where stats_info ->> '$.eventId' = " + strconv.Itoa(eventId) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"
	var machineStatics AssemblyMachineStatistics
	var machineStatsInfo MachineStatisticsInfo
	dbConnection.Raw(statisticsQuery).Scan(&machineStatics)

	_ = json.Unmarshal(machineStatics.StatsInfo, &machineStatsInfo)

	duration, _ := time.ParseDuration(warningMessageGenrationPeriod)
	durationInSeconds := duration.Seconds()

	// Find how many times wrning message should have been added
	timeDiff := time.Now().UTC().Unix() - scheduleStopTime
	noOfWaringMessage := timeDiff / int64(durationInSeconds)

	statsInfoWarnMessageLength := len(machineStatsInfo.WarningMessage)

	if int(noOfWaringMessage) != statsInfoWarnMessageLength {
		getCurrentSingaporeTime := util.GetZoneCurrentTime("Asia/Singapore")
		//Append warning message
		warningMessage := "Schedule is overrunning, do you want proceed stop? " + getCurrentSingaporeTime
		machineStatsInfo.WarningMessage = append(machineStatsInfo.WarningMessage, warningMessage)

		machineStstsInfo, _ := json.Marshal(machineStatsInfo)

		//Update last calculated machine statics
		dbError := dbConnection.Model(&AssemblyMachineStatistics{}).Where("ts = ?", machineStatics.TS).Where("machine_id = ?", machineStatics.MachineId).Update("stats_info", machineStstsInfo).Error
		if dbError != nil {
			ms.BaseService.Logger.Error("error in updating machine statitics with warning message", zap.String("error", dbError.Error()))
		}

		if targetEmailId != "" {
			ms.emailGenerator(eventName, machineName, targetEmailId, make([]string, 0))
		}
	}
}

func getAssemblyDowntime(hmiInfoList []MachineHMIInfoResult, dbConnection *gorm.DB, machineId int, startEpochMilli int64, endEpochMilli int64) int64 {
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

	machineLiveUnplannedDown := getAssemblyMachineLiveUnplanndedDownTime(dbConnection, machineId, startEpochMilli, endEpochMilli)
	return downTime + machineLiveUnplannedDown
}

func getAssemblyMachineLiveUnplanndedDownTime(dbConnection *gorm.DB, machineId int, startEpochMilli int64, endEpochMilli int64) int64 {
	//Get machine statistics based on the machine id and epochs
	// machine status == !Live
	// list stop_start list
	//While loop
	//	if row satisfy the condition machine status
	//  	machine status == Live
	//      append stop_start_list
	startTimeInSeconds := int64(startEpochMilli / 1000)
	endTimeInSeconds := int64(endEpochMilli / 1000)
	statsQuery := "select * from assembly_machine_statistics where machine_id=" + strconv.Itoa(machineId) + " and ts >= " + strconv.FormatInt(startTimeInSeconds, 10) + " and ts <= " + strconv.FormatInt(endTimeInSeconds, 10) + " order by ts asc"
	var machineStats []AssemblyMachineStatistics
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
			err = json.Unmarshal([]byte(machineStats[index-1].StatsInfo), &machineStatsInfo)

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

func getAssemblyCompletedValue(currentEtlId int, currentActual int, productionOrder string, machineId int, timeLineEvents *[]component.GeneralObject, ms *MachineService, projectId string) int {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]
	var completedValue int = currentActual
	var etlInfo ScheduledOrderEventInfo
	var statisticsResult MachineStatistics
	var statsInfo MachineStatisticsInfo
	for _, eventTimeLine := range *timeLineEvents {

		if eventTimeLine.Id == currentEtlId {
			continue
		}

		err := json.Unmarshal([]byte(eventTimeLine.ObjectInfo), &etlInfo)
		if err != nil {
			continue
		}

		startDate := convertStringToDateTime(etlInfo.StartDate)
		endDate := convertStringToDateTime(etlInfo.EndDate)

		if startDate.Error != nil || endDate.Error != nil {
			continue
		}
		if etlInfo.ProductionOrder == productionOrder {

			statQuery := "select * from assembly_machine_statistics where ts >= " + strconv.FormatInt(startDate.DateTimeEpoch, 10) + " and ts <= " + strconv.FormatInt(endDate.DateTimeEpoch, 10) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"

			dbConnection.Raw(statQuery).Scan(&statisticsResult)

			errStats := json.Unmarshal([]byte(statisticsResult.StatsInfo), &statsInfo)
			if errStats != nil {
				continue
			}

			completedValue += statsInfo.Actual

		}

	}

	return completedValue
}

func getAssemblyOverallRejectedQuantity(ms *MachineService, dbConnection *gorm.DB, productionOrderId int) int {

	overallRejectedQty := 0
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	listOfScheduledOrderEvents := productionOrderInterface.GetChildEventsOfAssemblyProductionOrder(ProjectID, productionOrderId)
	for _, eventId := range listOfScheduledOrderEvents {
		//Decode Machine timeline event info

		conditionString := " object_info ->> '$.eventId' = " + strconv.Itoa(eventId)
		hmiObject, _ := GetConditionalObjects(dbConnection, AssemblyMachineHmiTable, conditionString)

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
