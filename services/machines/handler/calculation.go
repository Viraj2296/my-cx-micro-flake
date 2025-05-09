package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (v *MachineService) OnTimer() {

	c := cron.New()

	moulding := getMachineConfiguration(v.CalculationConfig.MachineConfiguration, Moulding)
	assembly := getMachineConfiguration(v.CalculationConfig.MachineConfiguration, Assembly)
	tooling := getMachineConfiguration(v.CalculationConfig.MachineConfiguration, Tooling)

	if moulding.EnableCalculation {
		c.AddFunc(moulding.CalculationInterval, func() { MakeMachineStatisticsCalculation(v) })
	}

	if assembly.EnableCalculation {
		c.AddFunc(assembly.CalculationInterval, func() { v.MakeAssemblyMachineStatisticsCalculation() })
	}

	if tooling.EnableCalculation {
		c.AddFunc(tooling.CalculationInterval, func() { MakeToolingMachineStatisticsCalculation(v) })
	}

	c.Start()
}

func getMachineConfiguration(configuration []MachineConfig, machineType string) MachineConfig {
	machineConfig := MachineConfig{CalculationInterval: "@every 300s"}
	for _, config := range configuration {
		if config.MachineType == machineType {
			return config
		}
	}
	return machineConfig
}

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

func MakeMachineStatisticsCalculation(ms *MachineService) {
	//Get all data from Machine timeline event table
	dbConnection := ms.BaseService.ServiceDatabases[ProjectID]
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	//Check whether we are getting messages or not. If yes set it as Live
	updateMachineConnectStatus(ms, dbConnection, productionOrderInterface)

	mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, listOfEvents := productionOrderInterface.GetScheduledEvents(ProjectID)
	ms.BaseService.Logger.Info("events processing", zap.Any("number_of_events", len(*listOfEvents)))

	// Get production order complete status id from preference level
	orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(ProjectID, ScheduleStatusPreferenceSeven)

	for _, eventTimeLine := range *listOfEvents {
		//Decode Machine timeline event info
		var etlInfo ScheduledOrderEventInfo
		err := json.Unmarshal(eventTimeLine.ObjectInfo, &etlInfo)
		if err == nil {
			machineStatus, newMachineId := getMachineStatus(ms, dbConnection, etlInfo.MachineId)
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
					topicName := getTopicName(etlInfo.MachineId, ms, ProjectID)

					hmiInfoList := getHMIInfo(ms, ProjectID, etlInfo.MachineId, eventTimeLine.Id, MachineHMITable)
					ms.BaseService.Logger.Info("HMI info list", zap.Any("machine_id:", etlInfo.MachineId), zap.Any("list:", hmiInfoList))
					if len(hmiInfoList) > 0 {
						// getting messages based on HMI first created time and current time for processing
						// we needed only the in_timestamp message greater
						messageQuery := "select * from message where body->>'$.in_timestamp' >= " + strconv.FormatInt(convertStringToDateTime(hmiInfoList[0].Created).DateTimeEpochMilli, 10) + " and body->>'$.in_timestamp' <= " + strconv.FormatInt(time.Now().UnixMilli(), 10) + " and topic='" + topicName + "' order by body->>'$.in_timestamp' desc"
						ms.BaseService.Logger.Info("running query to select the messages", zap.Any("query", messageQuery))
						var messages []Message
						dbConnection.Raw(messageQuery).Scan(&messages)
						//messageOffSet = getCycleCountOffset(messages)
						//ms.BaseService.Logger.Infow("message offset", "offset", messageOffSet)
						// Fetch data based on the schedule from message table
						if len(messages) > 0 {

							err, productionOrderObject := productionOrderInterface.GetMachineProductionOrderInfo(ProjectID, etlInfo.EventSourceId, etlInfo.MachineId)
							if err != nil {
								ms.BaseService.Logger.Error("error getting production order info", zap.Any("event_source_id", etlInfo.EventSourceId), zap.Any("machine_id", etlInfo.MachineId))
								continue
							}
							productionOrderInfo := GetProductionOrderInfo(productionOrderObject.ObjectInfo)
							ms.BaseService.Logger.Info("production order", zap.Any("machine_id", etlInfo.MachineId), zap.Any("production order info ", productionOrderInfo), zap.Any("production_order_id", productionOrderObject.Id))

							// accumulatedCycleCount := getAccumulatedCycleCount(messages)
							// how much manufactured so far from PUMAS
							cycleCount, remark := getCycleCount(messages)

							cycleCountOneMinBefore := getCycleCountBeforeOneMin(messages)
							ms.BaseService.Logger.Info("cycle count", zap.Any("machine_id:", etlInfo.MachineId), zap.Any("cycleCount", cycleCount), zap.Any("cycleCountOneMinBefore", cycleCountOneMinBefore))

							productCount := ms.getProductCount(cycleCount, ProjectID, etlInfo, *productionOrderInfo)
							ms.BaseService.Logger.Info("product count", zap.Any("machine_id", etlInfo.MachineId), zap.Any("product count", productCount))

							//TODO, get the cavity using interface function
							// actualProduced :=cycleCount * noOfCavity

							//Get all hmis for given machine and event

							//based on UI or HMI
							dailyRejectedCount := getTotalRejectQtyForEvent(hmiInfoList)

							overallRejectedCount := getOverallRejectedQuantity(ms, dbConnection, etlInfo.EventSourceId)
							//getOverallRejectedQuantity(ms, dbConnection, etlInfo.ProductionOrder, eventTimeLines)
							goodCount := getGoodCount(productCount, dailyRejectedCount)
							machineDownTime := ms.getDowntime(hmiInfoList, dbConnection, etlInfo.MachineId, startEpoch, endEpoch)

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

							//_, listOfEventForProduction := productionOrderInterface.GetScheduledEventByProductionId(ProjectID, etlInfo.EventSourceId)

							actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, etlInfo.MachineId, eventTimeLine.Id, MachineHMITable)
							// completedValue := getCompletedValue(eventTimeLine.Id, productCount, etlInfo.ProductionOrder, etlInfo.MachineId, listOfEventForProduction, ms, ProjectID)
							completedValue := productionOrderInterface.GetCompletedQuantity(ProjectID, productionOrderObject.Id)
							progressPercentage := getProgressPercentage(productCount, etlInfo.ScheduledQty)

							ms.BaseService.Logger.Info("summary", zap.Any("machine_id", etlInfo.MachineId), zap.Any("event_id", eventTimeLine.Id), zap.Any("mould_id", etlInfo.MouldId),
								zap.Any("planned_production_time", plannedProductionTime), zap.Any("machine_down_time", machineDownTime), zap.Any("performance", performance),
								zap.Any("availability", availability), zap.Any("completedValue", completedValue), zap.Any("progressPercentage", progressPercentage),
								zap.Any("actualStartTime", actualStartTime), zap.Any("actualEndTime", actualEndTime), zap.Any("quality", quality), zap.Any("production_order_id", productionOrderObject.Id))
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
								ExpectedProductQyt:         getExpectedProductQty(cycleCount, cycleCountOneMinBefore, actualStartTime, endDate.DateTimeEpoch),
								WarningMessage:             make([]string, 0),
								Remark:                     remark,
							}
							statsInfoJson, _ := json.Marshal(statsInfo)
							machineStatistics := MachineStatistics{
								MachineId: etlInfo.MachineId,
								TS:        time.Now().Unix(),
								StatsInfo: statsInfoJson,
							}
							//ms.BaseService.Logger.Infow("machineStatistics", "stats", string(statsInfoJson))
							err = dbConnection.Create(&machineStatistics).Error
							if err != nil {
								ms.BaseService.Logger.Error("inserting stats has failed", zap.String("error", err.Error()))
							}
							updateMachineTimeLine(ProjectID, eventTimeLine.Id, productCount, progressPercentage, dailyRejectedCount, int(oee*10000))

							var key = strconv.Itoa(etlInfo.MouldId) + "_" + strconv.Itoa(eventTimeLine.Id)
							var deltaValue = 0
							if existingShotCount, ok := ms.MouldShotCountCache[key]; ok {
								deltaValue = cycleCount - existingShotCount
								ms.MouldShotCountCache[key] = cycleCount
							} else {
								ms.MouldShotCountCache[key] = cycleCount
								deltaValue = cycleCount
							}
							err = mouldInterface.UpdateShotCount(ProjectID, etlInfo.MouldId, deltaValue)
							if err != nil {
								ms.BaseService.Logger.Error("error updating shot count to mould ID", zap.Error(err), zap.Any("mould_id", etlInfo.MouldId))
							}

							//if productionOrder.OrderStatus == ScheduleStatusPreferenceTwo || productionOrder.OrderStatus == ScheduleStatusPreferenceSeven {
							//	continue
							//}
							//updateProductionOrderProdQty(ms, dbConnection, productionOrder, prodId, completedValue)
							// updateAllMachineTimeLineCompletedQty(ms, dbConnection, etlInfo.ProductionOrder, completedValue)

							//	Need to update view columns
							initialCycleCount, cycleCountMessage := getFirstCycleCount(messages)
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
							err = ms.ViewManager.CreateOrUpdateMouldingView(etlInfo.MachineId, updates)
							if err != nil {
								ms.BaseService.Logger.Error("error updating Moulding view", zap.String("error", err.Error()))
							}

						}
					} else {
						continue
					}
				} else {
					hmiInfoList := getHMIInfo(ms, ProjectID, etlInfo.MachineId, eventTimeLine.Id, MachineHMITable)

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

						// lastHmiStatus := hmiInfoList[len(hmiInfoList)-1]

						// if lastHmiStatus.Status != "stopped" {
						//If the schedule is overdue but but still hmi isn't stopped the we have to append warning message
						// hmiSettingInfo := getHmiSetting(dbConnection, ms, etlInfo.MachineId)
						//End date should be in seconds
						// appendWarningMessage(ms, dbConnection, hmiSettingInfo.WarningMessageGenerationPeriod, endDate.DateTimeEpoch, etlInfo.MachineId)

						//If user not stop HMI data, Then program automatically stopped hmi
						// Condition
						// Based on hmi stop configuration hmi will be stopped
						// automaticHmiStop(hmiSettingInfo.HmiAutoStopPeriod, lastHmiStatus, dbConnection, ms, endDate.DateTimeEpoch)
						// } else {

						// }

					}
					ms.BaseService.Logger.Info("machine is running out of time:", zap.Any("last_hmi_info_result", lastHMIInfoResult))
					if lastHMIInfoResult.HmiStatus != "stopped" {
						//If the schedule is overdue but but still hmi isn't stopped the we have to append warning message
						hmiSettingInfo := getHmiSetting(dbConnection, ms, etlInfo.MachineId)
						ms.BaseService.Logger.Info("hmi setting info", zap.Any("hmi_setting_info", hmiSettingInfo), zap.Any("endDate:", endDate))
						//End date should be in seconds
						appendWarningMessage(ms, dbConnection, hmiSettingInfo.WarningMessageGenerationPeriod, endDate.DateTimeEpoch, eventTimeLine.Id, etlInfo.MachineId, etlInfo.Name, newMachineId, hmiSettingInfo.WarningTargetEmailId)
						addForceStopFlag(ProjectID, eventTimeLine.Id)
						//If user not stop HMI data, Then program automatically stopped hmi
						// Condition
						// Based on hmi stop configuration hmi will be stopped
						ms.automaticHmiStop(hmiSettingInfo.HmiAutoStopPeriod, lastHMIInfoResult, dbConnection, ms, endDate.DateTimeEpoch)
					}

				}
			}

		}
	}
}

func (v *MachineService) automaticHmiStop(hmiAutoStopPeriod string, info MachineHMIInfoResult, dbConnection *gorm.DB, ms *MachineService, scheduleStopTime int64) {
	//Schedule stop time in seconds
	duration, _ := time.ParseDuration(hmiAutoStopPeriod)
	durationInSeconds := duration.Seconds()
	v.BaseService.Logger.Info("performing automatic HMI stop ", zap.Any("info", info))
	if time.Now().UTC().Unix() > (scheduleStopTime + int64(durationInSeconds)) {
		v.BaseService.Logger.Info("Time is exceeded ", zap.Any("now", time.Now().UTC().Unix()), zap.Any("scheduleStopTime", scheduleStopTime), zap.Any("duration seconds", int64(durationInSeconds)))
		updateOverDueHmiInfo(info, dbConnection, ms)
	}
}

func appendWarningMessage(ms *MachineService, dbConnection *gorm.DB, warningMessageGenrationPeriod string, scheduleStopTime int64, eventId int, machineId int, eventName string, machineName string, targetEmailId string) {
	// Schedule stop time in seconds
	statisticsQuery := "select * from machine_statistics where stats_info ->> '$.eventId' = " + strconv.Itoa(eventId) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"
	var machineStatics MachineStatistics
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

		machineStatsInfo, _ := json.Marshal(machineStatsInfo)

		//Update last calculated machine statics
		dbError := dbConnection.Model(&MachineStatistics{}).Where("ts = ?", machineStatics.TS).Where("machine_id = ?", machineStatics.MachineId).Update("stats_info", machineStatsInfo).Error

		if dbError != nil {
			ms.BaseService.Logger.Error("error in updating machine statics with warning message", zap.String("error", dbError.Error()))
		}

		if targetEmailId != "" {
			ms.emailGenerator(eventName, machineName, targetEmailId, make([]string, 0))
		}

	}

}

func getHmiSetting(dbConnection *gorm.DB, ms *MachineService, machineId int) *HmiSettingInfo {
	err, generalObject := Get(dbConnection, MachineHMISettingSettingTable, machineId)
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

func getPlannedProductionTime(hmiInfoList []MachineHMIInfoResult) int64 {
	if len(hmiInfoList) == 0 {
		return 0
	}
	createdTime := convertStringToDateTime(hmiInfoList[0].Created)
	timeNow := time.Now().UTC().UnixMilli()
	return timeNow - createdTime.DateTimeEpochMilli
}

// If the user not stopped the hmi info, program automatically is inserted stopped hmi
// This function only update created , status and remark attributes
func updateOverDueHmiInfo(info MachineHMIInfoResult, dbConnection *gorm.DB, ms *MachineService) {
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

	dbError, _ := Create(dbConnection, MachineHMITable, component.GeneralObject{ObjectInfo: hmiObjectInfo})

	if dbError != nil {
		ms.BaseService.Logger.Error("error in updating hmi info", zap.String("error", err.Error()))
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	productionOrderInterface.UpdateOrderPreferenceLevel(ProjectID, 2, info.EventId, ScheduleStatusPreferenceSix)

}

func getAllMachines(ms *MachineService, dbConnection *gorm.DB) []MachineMaster {
	machineQuery := "select * from machine_master"
	var machineMasters []MachineMaster
	dbConnection.Raw(machineQuery).Scan(&machineMasters)

	return machineMasters
}

// //Update machine master whether we receive message or not. Live or Waiting for feed
// func updateMachineConnectStatus(ms *MachineService, dbConnection *gorm.DB) {
// 	timeNow := time.Now().UTC().UnixMilli()
// 	machineMasters := getAllMachines(ms, dbConnection)
// 	var machineMasterInfo MachineMasterInfo
// 	for _, machine := range machineMasters {
// 		err := json.Unmarshal([]byte(machine.ObjectInfo), &machineMasterInfo)
// 		hmiSettingInfo := getHmiSetting(dbConnection, ms, machine.Id)

// 		if hmiSettingInfo == nil {
// 			continue
// 		}
// 		machineLiveDetectionInterval := hmiSettingInfo.MachineLiveDetectionInterval
// 		duration, _ := time.ParseDuration(machineLiveDetectionInterval)
// 		durationInSeconds := duration.Milliseconds()

// 		liveDetectionPeriod := timeNow - durationInSeconds

// 		if err != nil {
// 			ms.BaseService.Logger.Error("error in unmarshelling machine master", zap.String("error", err.Error()))
// 			continue
// 		}
// 		//Get past 15 seconds message
// 		messageQuery := "select * from message where ts > " + strconv.FormatInt(liveDetectionPeriod, 10) + " and topic = '" + "machines/" + machineMasterInfo.NewMachineId + "'"
// 		var messages []Message
// 		dbConnection.Raw(messageQuery).Scan(&messages)

// 		if len(messages) > 0 {
// 			if machineMasterInfo.MachineConnectStatus != "Live" {
// 				generateSystemNotification(ms, ProjectID, "Machine "+machineMasterInfo.NewMachineId+" is started", machine.Id)
// 			}
// 			machineMasterInfo.MachineConnectStatus = "Live"

// 		} else {
// 			if machineMasterInfo.MachineConnectStatus != "Waiting For Feed" {
// 				generateSystemNotification(ms, ProjectID, "Machine "+machineMasterInfo.NewMachineId+" is stopped", machine.Id)
// 			}
// 			machineMasterInfo.MachineConnectStatus = "Waiting For Feed"
// 		}
// 		machineMasterInfo.LastUpdatedMachineLiveStatus = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
// 		machineObjectInfo, errMachines := json.Marshal(machineMasterInfo)

// 		if errMachines != nil {
// 			ms.BaseService.Logger.Error("error in marshelling machine master", zap.String("error", err.Error()))
// 			continue
// 		}
// 		dbError := dbConnection.Model(&MachineMaster{}).Where("id = ?", machine.Id).Update("object_info", machineObjectInfo).Error

// 		if dbError != nil {
// 			ms.BaseService.Logger.Error("error updating machine master", zap.String("error", err.Error()))
// 		}

// 	}

// }

// Update machine master whether we receive message or not. Live or Waiting for feed
func updateMachineConnectStatus(ms *MachineService, dbConnection *gorm.DB, productionInterface common.ProductionOrderInterface) {
	// year, month, day := time.Now().Date()
	timeNow := time.Now().UTC().UnixMilli()
	machineMasters := getAllMachines(ms, dbConnection)

	for _, machine := range machineMasters {
		machineMasterInfo := MachineMasterInfo{}
		// Fetching data from Machine
		//_, originalMachineObject := Get(dbConnection, MachineMasterTable, machine.Id)
		err := json.Unmarshal(machine.ObjectInfo, &machineMasterInfo)
		hmiSettingInfo := getHmiSetting(dbConnection, ms, machine.Id)

		if hmiSettingInfo == nil {
			ms.BaseService.Logger.Warn("hmi setting is null for machine ID ", zap.Int("machine_id", machine.Id))
			continue
		}
		machineLiveDetectionInterval := hmiSettingInfo.MachineLiveDetectionInterval
		duration, _ := time.ParseDuration(machineLiveDetectionInterval)
		durationInSeconds := duration.Milliseconds()

		liveDetectionPeriod := timeNow - durationInSeconds

		if err != nil {
			ms.BaseService.Logger.Error("error in unmarshalling machine master", zap.String("error", err.Error()), zap.Int("machine_id", machine.Id))
			continue
		}
		//Get past 15 seconds message
		cycleCountMessageQuery := "select * from message where ts > " + strconv.FormatInt(liveDetectionPeriod, 10) + " and topic = '" + "machines/" + machineMasterInfo.NewMachineId + "' order by ts desc"
		var messagesCycleCount []Message
		dbConnection.Raw(cycleCountMessageQuery).Scan(&messagesCycleCount)

		connectionStatusCycleCount, currentCycleCount := findCycleCountIncrement(messagesCycleCount)
		delayTime, delayStatus, findMsg := findDelayMessage(messagesCycleCount)

		if findMsg {
			machineMasterInfo.DelayPeriod = delayTime
			machineMasterInfo.DelayStatus = delayStatus
		}

		machineMasterInfo.CurrentCycleCount = currentCycleCount

		// theTime := time.Date(year, month, day, 00, 01, 00, 00, time.UTC)
		//Get past 15 seconds message
		// messageQuery := "select * from message where ts > " + strconv.FormatInt(theTime.UnixMilli(), 10) + " and topic = '" + "machines/" + machineMasterInfo.NewMachineId + "' order by ts desc"
		// var messages []Message
		// dbConnection.Raw(messageQuery).Scan(&messages)

		// Searching machine connection status attribute in the messages
		// connectionStatus := searchMachineStatus(messages)

		// fmt.Println("[message_query]:", messageQuery)
		ms.BaseService.Logger.Info("calculation processing", zap.Any("machineMasterInfo", machineMasterInfo.NewMachineId))
		ms.BaseService.Logger.Info("calculation processing", zap.Any("CurrentCycleCount", machineMasterInfo.CurrentCycleCount))

		if connectionStatusCycleCount {
			//if machineMasterInfo.MachineConnectStatus != "Live" {
			//	generateSystemNotification(ms, ProjectID, "Machine "+machineMasterInfo.NewMachineId+" is started", machine.Id)
			//}
			machineMasterInfo.MachineConnectStatus = 1
		} else {
			//if machineMasterInfo.MachineConnectStatus != "Waiting For Feed" {
			//	generateSystemNotification(ms, ProjectID, "Machine "+machineMasterInfo.NewMachineId+" is stopped", machine.Id)
			//}
			machineMasterInfo.MachineConnectStatus = 2
		}

		machineMasterInfo.LastUpdatedMachineLiveStatus = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

		// Added corrective maintenance flag
		// this code is commented as client requested to enable the canCreateCorrectiveWorkOrder always true
		/*
			_, currentSchedule := productionInterface.GetCurrentScheduledEvent(ProjectID, machine.Id)

			if machineMasterInfo.EnableProductionOrder {
				if currentSchedule.Id == 0 {
					machineMasterInfo.CanCreateCorrectiveWorkOrder = false
				} else {
					machineMasterInfo.CanCreateCorrectiveWorkOrder = true
				}
			} else {
				machineMasterInfo.CanCreateCorrectiveWorkOrder = true
			}
		*/
		machineMasterInfo.CanCreateCorrectiveWorkOrder = true
		//machineObjectInfo, errMachines := json.Marshal(machineMasterInfo)

		// fmt.Println("[machineMasterInfo_for_cycle_count]:", machineMasterInfo.DelayPeriod)

		//if errMachines != nil {
		//	ms.BaseService.Logger.Error("error in marshelling machine master", zap.String("error", err.Error()))
		//	continue
		//}
		//dbError := dbConnection.Model(&MachineMaster{}).Where("id = ?", machine.Id).Update("object_info", machineObjectInfo).Error

		_, originalMachineObject := Get(dbConnection, MachineMasterTable, machine.Id)
		recentMachineMasterInfo := MachineMasterInfo{}
		json.Unmarshal(originalMachineObject.ObjectInfo, &recentMachineMasterInfo)

		if recentMachineMasterInfo.MachineStatus != MachineStatusMaintenance {
			updateQuery := "UPDATE machine_master SET object_info = JSON_SET(object_info, '$.machineConnectStatus'," + strconv.Itoa(machineMasterInfo.MachineConnectStatus) + ", '$.delayPeriod'," + strconv.FormatInt(machineMasterInfo.DelayPeriod, 10) + ", '$.delayStatus','" + machineMasterInfo.DelayStatus + "', '$.lastUpdatedMachineLiveStatus','" + machineMasterInfo.LastUpdatedMachineLiveStatus + "', '$.canCreateCorrectiveWorkOrder'," + strconv.FormatBool(machineMasterInfo.CanCreateCorrectiveWorkOrder) + ", '$.currentCycleCount'," + strconv.Itoa(machineMasterInfo.CurrentCycleCount) + ") WHERE id = ?"
			result := dbConnection.Exec(updateQuery, machine.Id)
			if result.Error != nil {
				ms.BaseService.Logger.Error("error updating machine master", zap.String("error", err.Error()))
			}
		}

		//	Adding data into moulding machine view
		updates := map[string]interface{}{
			"delay_period":           strconv.FormatInt(delayTime, 10),
			"delay_status":           delayStatus,
			"current_cycle_count":    currentCycleCount,
			"machine_connect_status": machineMasterInfo.MachineConnectStatus,
		}
		err = ms.ViewManager.CreateOrUpdateMouldingView(machine.Id, updates)
		if err != nil {
			ms.BaseService.Logger.Error("error updating moulding view", zap.String("error", err.Error()))
		}

	}

}

// Assumed message is in decending order
func findDelayMessage(messageList []Message) (int64, string, bool) {
	isFindMsg := false
	if len(messageList) != 0 {
		isFindMsg = true
		// timeNow := util.NowAsUnixMilli()

		msgBody := make(map[string]interface{}, 0)
		json.Unmarshal(messageList[0].Body, &msgBody)
		dbTs := messageList[0].TS

		var delay int64
		if _, ok := msgBody["in_timestamp"]; ok {
			delay = int64(math.Abs(float64(dbTs - int64(util.InterfaceToInt(msgBody["in_timestamp"])))))
		}

		if delay > 10000 {
			return delay / 1000, "Delayed Feed", isFindMsg
		}
		return 0, "On-Time", isFindMsg
	}
	return 0, "Delayed Feed", isFindMsg
}

func findCycleCountIncrement(messageList []Message) (bool, int) {
	machineStatus := false
	firstCycleCountFound := false
	firstCycleCountValue := 0

	for _, msg := range messageList {
		msgBody := make(map[string]interface{}, 0)

		json.Unmarshal(msg.Body, &msgBody)
		if valBody, ok := msgBody["cycle_count"]; ok {
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

// func searchMachineStatus(messageList []Message) bool {
// 	machineStatus := false

// 	for _, msg := range messageList {
// 		msgBody := make(map[string]interface{}, 0)

// 		json.Unmarshal(msg.Body, &msgBody)
// 		if valBody, ok := msgBody["machine_connection_status"]; ok {
// 			machineConnectionStatus := util.InterfaceToInt(valBody)
// 			if machineConnectionStatus == 1 {
// 				machineStatus = true
// 			}
// 			break
// 		}
// 	}
// 	return machineStatus
// }

func generateSystemNotification(ms *MachineService, projectId string, notificationDescription string, machineId int) {
	notificationHeader := notificationDescription
	ms.createSystemNotification(projectId, notificationHeader, notificationDescription, machineId)
}

func getProgressPercentage(completedQty int, scheduledQty int) float64 {
	fmt.Println("Progress percent===========================", completedQty)
	if scheduledQty != 0 {
		return math.Round((float64(completedQty) / float64(scheduledQty)) * 100)
	}
	return 0.0
}

////Update all the machine time line for given production order with completed quantity
//func updateAllMachineTimeLineCompletedQty(ms *MachineService, dbConnection *gorm.DB, productionOrderId string, completedQty int) {
//	eventQuery := "select * from machine_timeline_event where JSON_EXTRACT(object_info, \"$.productionOrder\") = '" + productionOrderId + "'"
//	var eventTimeLines []MachineTimelineEvent
//	dbConnection.Raw(eventQuery).Scan(&eventTimeLines)
//	var etlInfo MachineTimelineInfo
//	for _, eventTimeLine := range eventTimeLines {
//		err := json.Unmarshal([]byte(eventTimeLine.ObjectInfo), &etlInfo)
//		if err != nil {
//			ms.BaseService.Logger.Error("error unmarshelling machine time line", zap.String("error", err.Error()))
//			continue
//		}
//		startDate := convertStringToDateTime(etlInfo.StartDate)
//		if startDate.DateTimeEpoch > time.Now().Unix() {
//			etlInfo.CompletedQty = completedQty
//			etlObjectInfo, errEtl := json.Marshal(etlInfo)
//
//			if errEtl != nil {
//				continue
//			}
//
//			dbError := dbConnection.Model(&MachineTimelineEvent{}).Where("id = ?", eventTimeLine.Id).Update("object_info", etlObjectInfo).Error
//
//			if dbError != nil {
//				ms.BaseService.Logger.Error("error updateing machine time line", zap.String("error", err.Error()))
//			}
//
//		}
//	}
//}

// Update given Machine time line with completed qty, progress percent and rejected quantity
func updateMachineTimeLine(projectId string, eventId int, completedQty int, progressPercent float64, rejectedQty int, oee int) {

	var updatingFields = make(map[string]interface{})
	updatingFields["completedQty"] = completedQty
	updatingFields["percentDone"] = progressPercent
	updatingFields["rejectedQty"] = rejectedQty
	updatingFields["oee"] = oee
	serializedObject, _ := json.Marshal(updatingFields)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err := productionOrderInterface.UpdateScheduledOrderFields(projectId, eventId, serializedObject)
	if err != nil {
		fmt.Println("error updating scheduler order fields")
	}

}

// Update given Machine time line with completed qty, progress percent and rejected quantity
func addForceStopFlag(projectId string, eventId int) {

	var updatingFields = make(map[string]interface{})
	updatingFields["canForceStop"] = true

	serializedObject, _ := json.Marshal(updatingFields)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err := productionOrderInterface.UpdateScheduledOrderFields(projectId, eventId, serializedObject)
	if err != nil {
		fmt.Println("error updating scheduler order fields from  addForceStopFlag")
	}

}

// Get hmi info for given machine id and event id
func getHMIInfo(ms *MachineService, projectId string, machineId int, eventId int, targetTable string) []MachineHMIInfoResult {
	//Fetch hmi info which are in time line schedule
	dbConnection := ms.BaseService.ServiceDatabases[projectId]
	hmiQuery := "select * from " + targetTable + " where object_info->>'$.machineId' = " + strconv.Itoa(machineId) + " and object_info->>'$.eventId' = " + strconv.Itoa(eventId) + " order by object_info->>'$.createdAt' asc"
	var hmiResult []MachineHMI

	var hmiResponseList []MachineHMIInfoResult

	dbConnection.Raw(hmiQuery).Scan(&hmiResult)

	if len(hmiResult) == 0 {
		return make([]MachineHMIInfoResult, 0)
	}

	for _, hmi := range hmiResult {
		hmiInfo := MachineHMIInfo{}
		err := json.Unmarshal(hmi.ObjectInfo, &hmiInfo)

		if err != nil {
			ms.BaseService.Logger.Error("error updating machine hmi info", zap.String("error", err.Error()))
			continue
		}
		TS := convertStringToDateTime(hmiInfo.CreatedAt)

		if TS.Error != nil {
			ms.BaseService.Logger.Error("error in converting hmi date", zap.String("error", TS.Error.Error()))
			continue
		}

		hmiResponseList = append(hmiResponseList, MachineHMIInfoResult{
			Created:        hmiInfo.CreatedAt,
			EventId:        hmiInfo.EventId,
			Operator:       hmiInfo.Operator,
			ReasonId:       hmiInfo.ReasonId,
			MachineId:      hmiInfo.MachineId,
			Status:         hmiInfo.Status,
			HmiStatus:      hmiInfo.HMIStatus,
			RejectQuantity: hmiInfo.RejectedQuantity,
			Id:             hmi.Id,
			TS:             TS,
			Remark:         hmiInfo.Remark,
		})

	}

	return hmiResponseList
}

func getTopicName(machineId int, ms *MachineService, projectId string) string {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]
	machineQuery := "select * from machine_master where id = " + strconv.Itoa(machineId)
	var machineMaster MachineMaster
	err := dbConnection.Raw(machineQuery).Scan(&machineMaster).Error

	if err != nil {
		ms.BaseService.Logger.Error("error in fetching machine master", zap.String("error", err.Error()))
		return ""
	}

	var machineInfo MachineMasterInfo
	errMachineInfo := json.Unmarshal([]byte(machineMaster.ObjectInfo), &machineInfo)
	if errMachineInfo != nil {
		return ""
	}
	return "machines/" + machineInfo.NewMachineId
}

// Machine status is decided by hmi or cycle count of message
func getMachineStatus(ms *MachineService, dbConnection *gorm.DB, machineId int) (string, string) {
	machineQuery := "select * from machine_master where id=" + strconv.Itoa(machineId)
	var machineMaster MachineMaster
	var machineInfo MachineMasterInfo
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

func getDailyCompletedPercentage(PlannedQuality int, Completed int) float32 {
	if Completed > PlannedQuality {
		return 100.0
	}
	if PlannedQuality > 0 {
		return (float32(Completed) / float32(PlannedQuality)) * 100
	}
	return 0.0
}

// we send 100 incase we exceed a more than 100 %
func getOverallCompletedPercentage(ProductionQty int, OverallCompleted int) float32 {
	if OverallCompleted > ProductionQty {
		return 100.0
	}
	if ProductionQty > 0 {
		return (float32(OverallCompleted) / float32(ProductionQty)) * 100
	}
	return 0.0
}

func getAvailability(plannedProductionTime, downTime int64) float64 {
	runTime := plannedProductionTime - downTime
	if plannedProductionTime > 0 {
		return float64(runTime) / float64(plannedProductionTime)
	}
	return 0
}

func getGoodCount(productCount int, rejected int) int {
	goodCount := productCount - rejected

	if goodCount < 0 {
		goodCount = 0
	}
	return goodCount
}

func (v *MachineService) getDowntime(hmiInfoList []MachineHMIInfoResult, dbConnection *gorm.DB, machineId int, startEpochMilli int64, endEpochMilli int64) int64 {
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

	machineLiveUnplannedDown := v.getMachineLiveUnplannedDownTime(dbConnection, machineId, startEpochMilli, endEpochMilli)
	return downTime + machineLiveUnplannedDown
}

// Should be in Milliseconds
func (v *MachineService) getMachineLiveUnplannedDownTime(dbConnection *gorm.DB, machineId int, startEpochMilli int64, endEpochMilli int64) int64 {
	//Get machine statistics based on the machine id and epochs
	// machine status == !Live
	// list stop_start list
	//While loop
	//	if row satisfy the condition machine status
	//  	machine status == Live
	//      append stop_start_list
	startTimeInSeconds := int64(startEpochMilli / 1000)
	endTimeInSeconds := int64(endEpochMilli / 1000)
	statsQuery := "select * from machine_statistics where machine_id=" + strconv.Itoa(machineId) + " and ts >= " + strconv.FormatInt(startTimeInSeconds, 10) + " and ts <= " + strconv.FormatInt(endTimeInSeconds, 10) + " order by ts asc"
	var machineStats []MachineStatistics
	var machineStatsInfo MachineStatisticsInfo
	dbConnection.Raw(statsQuery).Scan(&machineStats)
	v.BaseService.Logger.Info("get machine line unplanned down time, stats query", zap.Any("statsQuery", statsQuery), zap.Any("status length ", len(machineStats)))
	//searchingStatus := machineConnectStatusWaitingForFeed
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

// TODO based on all the events belongs to that proudction order A = childs
func getOverallRejectedQuantity(ms *MachineService, dbConnection *gorm.DB, productionOrderId int) int {

	overallRejectedQty := 0
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	listOfScheduledOrderEvents := productionOrderInterface.GetChildEventsOfProductionOrder(ProjectID, productionOrderId)
	for _, eventId := range listOfScheduledOrderEvents {
		//Decode Machine timeline event info

		conditionString := " object_info ->> '$.eventId' = " + strconv.Itoa(eventId)
		hmiObject, err := GetConditionalObjects(dbConnection, MachineHMITable, conditionString)
		if err == nil {
			for _, machineHmi := range *hmiObject {
				machineHMI := MachineHMI{ObjectInfo: machineHmi.ObjectInfo}

				rejectedQty := machineHMI.getMachineHMIInfo().RejectedQuantity
				if rejectedQty != nil {
					overallRejectedQty += *rejectedQty
				}
			}
		} else {
			fmt.Println("[getOverallRejectedQuantity] error getting data from MachineHMITable ")
		}

	}

	return overallRejectedQty

}

func getTotalRejectQtyForEvent(hmiInfoList []MachineHMIInfoResult) int {
	overallRejectedCount := 0

	for _, hmi := range hmiInfoList {
		if hmi.RejectQuantity != nil {
			overallRejectedCount = overallRejectedCount + *hmi.RejectQuantity
		}

	}
	return overallRejectedCount
}

func convertStringToDateTime(dateTimeString string) DateTimeInfo {
	layout := "2006-01-02T15:04:05.000Z"
	dateTime, err := time.Parse(layout, dateTimeString)
	if err != nil {
		return DateTimeInfo{Error: err}
	}
	return DateTimeInfo{DateTimeString: dateTimeString, DateTime: dateTime, DateTimeEpochMilli: dateTime.UTC().UnixMilli(), DateTimeEpoch: dateTime.UTC().Unix(), Error: nil}
}

//func getCycleCount(messages []Message) int {
//	// first get the maximum number stop and again do that until hit the maximum
//	var messageBody MachineDataInfo
//
//	// Append 0 data if it doesn't in messages
//	messages = messagePreprocessing(messages)
//
//	var cycleCountArray []int
//	for _, message := range messages {
//		err := json.Unmarshal(message.Body, &messageBody)
//
//		if err != nil || messageBody.CycleCount == nil {
//			continue
//		}
//		cycleCountArray = append(cycleCountArray, *messageBody.CycleCount)
//	}
//	if len(cycleCountArray) < 2 {
//		return 0
//	}
//
//	var countCache = make(map[int][]int)
//	var cacheIndex int
//	//index 0 , 1 , 2 , 3 , 4, 5 , 6 , 7
//	// values 34, 33, 32,31,30, 0, 98, 97
//	for _, currentValue := range cycleCountArray {
//		if currentValue == 0 {
//			cacheIndex = cacheIndex + 1
//
//		} else {
//			countCache[cacheIndex] = append(countCache[cacheIndex], currentValue)
//
//		}
//
//	}
//
//	if len(countCache) == 1 {
//		countArray := countCache[0]
//		if len(countArray) > 0 {
//			return countArray[0] - countArray[len(countArray)-1]
//		} else {
//			return 0
//		}
//
//	} else {
//		var cycleCount int
//		var offset int
//		for _, countArray := range countCache {
//			cycleCount += countArray[0]
//			offset = countArray[len(countArray)-1]
//		}
//		return cycleCount - offset
//	}
//}

func calculateTotalDifference(arr []int) int {
	if len(arr) < 2 {
		return 0
	}

	// Initialize variables
	totalDifference := 0

	for i := 1; i < len(arr); i++ {
		diff := arr[i-1] - arr[i]

		// Treat negative differences as 1
		if diff < 0 {
			diff = 1
		}

		// Add the difference to the total
		totalDifference += diff
	}

	return totalDifference
}

func getFirstCycleCount(messages []Message) (int, MachineDataInfo) {

	var cycleCount int
	var cycleCountMessage = MachineDataInfo{}

	for i := len(messages) - 1; i >= 0; i-- {
		var message = messages[i]
		var messageBody = MachineDataInfo{}
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if messageBody.CycleCount != nil && *messageBody.CycleCount >= 0 {
			cycleCount = *messageBody.CycleCount
			break
		}
	}

	// This is for finding last message
	for _, message := range messages {
		var messageBody = MachineDataInfo{}
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if messageBody.CycleCount != nil && *messageBody.CycleCount >= 0 {
			cycleCountMessage = messageBody
			break
		}
	}

	return cycleCount, cycleCountMessage
}

func getCycleCount(messages []Message) (int, string) {
	// first get the maximum number stop and again do that until hit the maximum
	var messageBody MachineDataInfo
	// we need to basically count total numbers excluding 0
	var cycleCountArray []int
	var totalCycleCount = 0

	for _, message := range messages {
		err := json.Unmarshal(message.Body, &messageBody)
		if err != nil {
			continue
		}
		if messageBody.CycleCount != nil && *messageBody.CycleCount >= 0 {
			cycleCountArray = append(cycleCountArray, *messageBody.CycleCount)
		}
	}
	totalCycleCount = calculateTotalDifference(cycleCountArray)
	var remark string
	if totalCycleCount < 0 {
		totalCycleCount = 0
		remark = "Cycle count has negative value"
	}
	return totalCycleCount, remark
}

func getCycleCountBeforeOneMin(messages []Message) int {
	beforeOneMinTime := time.Now().UTC().UnixMilli() - 60000

	breakPointer := len(messages) - 1
	for index, message := range messages {
		messageBody := MachineDataInfo{}
		err := json.Unmarshal(message.Body, &messageBody)

		if err == nil {
			if int64(*messageBody.InTimestamp) < beforeOneMinTime {
				breakPointer = index
				break
			}
		}
	}

	oneMinBeforeMessages := messages[breakPointer:]
	cycleCount, _ := getCycleCount(oneMinBeforeMessages)

	return cycleCount

}

func messagePreprocessing(messages []Message) []Message {
	zeroInsertIndex := make([]map[string]interface{}, 0)
	for index, message := range messages {
		messageBody := MachineDataInfo{}
		err := json.Unmarshal(message.Body, &messageBody)

		if err != nil || messageBody.CycleCount == nil {
			continue
		}

		if *messageBody.CycleCount == 1 {
			previousMessageBody := MachineDataInfo{}
			if index != 0 {
				err = json.Unmarshal(messages[index-1].Body, &previousMessageBody)

				if err == nil {
					if previousMessageBody.CycleCount != nil {
						if *previousMessageBody.CycleCount != 0 {
							data := map[string]interface{}{"index": index, "in_timestamp": *previousMessageBody.InTimestamp, "topic": messages[index-1].Topic, "ts": message.TS}
							zeroInsertIndex = append(zeroInsertIndex, data)
						}
					}

				}
			}
		}
	}

	if len(zeroInsertIndex) == 0 {
		return messages
	}

	for _, zeroIndex := range zeroInsertIndex {
		messageBody := make(map[string]interface{})
		messageBody["cycle_count"] = 0
		messageBody["in_timestamp"] = zeroIndex["in_timestamp"]

		jsonBody, _ := json.Marshal(messageBody)
		msg := Message{Body: jsonBody, Topic: util.InterfaceToString(zeroIndex["topic"]), TS: int64(util.InterfaceToInt(zeroIndex["ts"]))}

		messages = insertSlice(messages, util.InterfaceToInt(zeroIndex["index"]), msg)
	}

	return messages
}

func insertSlice(slice []Message, index int, data Message) []Message {
	// Create a new slice with a length greater than the original slice
	newSlice := make([]Message, len(slice)+1)

	// Copy the elements before the insertion point
	copy(newSlice[:index], slice[:index])

	// Insert the new data
	newSlice[index] = data

	// Copy the elements after the insertion point
	copy(newSlice[index+1:], slice[index:])

	return newSlice
}

func (v *MachineService) getProductCount(cycleCount int, projectId string, scheduledOrderEventInfo ScheduledOrderEventInfo, productionOrderMaster ProductionOrderInfo) int {
	// if productionOrderInfo.NoOfCustomMouldCavity > 0 {
	// 	return cycleCount * productionOrderInfo.NoOfCustomMouldCavity
	// }
	if scheduledOrderEventInfo.CustomCavity > 0 {
		return cycleCount * scheduledOrderEventInfo.CustomCavity
	}

	// Get default mould id from mould info
	if scheduledOrderEventInfo.MouldId != 0 {
		mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		_, mouldGeneralObj := mouldService.GetMouldInfoById(projectId, scheduledOrderEventInfo.MouldId)
		mouldInfo := make(map[string]interface{})
		json.Unmarshal(mouldGeneralObj.ObjectInfo, &mouldInfo)

		ca := util.InterfaceToInt(mouldInfo["noOfCav"])
		productCount := cycleCount * ca
		v.BaseService.Logger.Info("getting product count ", zap.Any("mould_id", scheduledOrderEventInfo.MouldId), zap.Any(" loaded cavity ", ca), zap.Any("product count", productCount))
		return productCount
	} else {
		v.BaseService.Logger.Info("getting product count, mould id 0 ", zap.Any("mould_id", scheduledOrderEventInfo.MouldId))
	}

	// mouldRecordId := productionOrderInfo.MouldId
	// mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	// _, noOfCavity := mouldInterface.GetNoOfCavity(projectId, mouldRecordId)
	noOfCavity := 1

	return cycleCount * noOfCavity
}

// func getAccumulatedCycleCount(messages []Message) int {
// 	var messageBody MachineDataInfo
// 	var cycleCount int = 0
// 	for _, message := range messages {
// 		err := json.Unmarshal([]byte(message.Body), &messageBody)

//			if err != nil || messageBody.CycleCount == nil {
//				continue
//			}
//			cycleCount += *messageBody.CycleCount
//		}
//		return cycleCount
//	}
//
// 15026 - 6000
// --------------  , 15026 /5
//
//	6000
//
// ((current time x 60 x 60)/ current quantity) X (daily quantity – current quantity)
func getEstimatedEndTime(cycleTime float32, productCount int, actualStartTime string, dailyScheduledQuantity int) string {
	if cycleTime == 0 {
		cycleTime = 1
	}
	actualStartTimeTS := convertStringToDateTime(actualStartTime)
	currentTime := time.Now().UTC().Unix()
	timeTakenToComplete := currentTime - actualStartTimeTS.DateTimeEpoch
	var timeNeededToComplete int64
	if productCount > 0 {
		timeNeededToComplete = int64(((float32(dailyScheduledQuantity) - float32(productCount)) / float32(productCount)) * float32(timeTakenToComplete))
	}
	fmt.Println("dailyScheduledQuantity: ", dailyScheduledQuantity, "productCount : ", productCount, "timeTakenToComplete:", timeTakenToComplete)
	fmt.Println("currentTime : ", currentTime, "timeTakenToComplete: ", timeTakenToComplete, "timeNeededToComplete : ", timeNeededToComplete)
	estimatedEndTime := currentTime + timeNeededToComplete
	fmt.Println("estimatedEndTime : ", estimatedEndTime)
	return time.Unix(estimatedEndTime, 0).Format("2006-01-02T15:04:05.000Z")
}

func getExpectedProductQty(currentCycleCount int, oneMinBeforeCycleCount int, actualStartTime string, scheduleEndTime int64) int {
	actualStartTimeTS := convertStringToDateTime(actualStartTime)
	var expectedProductQty int
	currentCycleCountInLastMin := currentCycleCount - oneMinBeforeCycleCount
	expectedProductQty = int(scheduleEndTime-actualStartTimeTS.DateTimeEpoch) * currentCycleCountInLastMin
	return expectedProductQty
}

func getCompletedValue(currentEtlId int, currentActual int, productionOrder string, machineId int, timeLineEvents *[]component.GeneralObject, ms *MachineService, projectId string) int {
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

			statQuery := "select * from machine_statistics where ts >= " + strconv.FormatInt(startDate.DateTimeEpoch, 10) + " and ts <= " + strconv.FormatInt(endDate.DateTimeEpoch, 10) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"

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

// Actual start time based on hmi start and end time
func getActualStartStopTime(dbConnection *gorm.DB, machineId int, eventId int, targetTable string) (string, string) {
	hmiQuery := "select * from " + targetTable + " where JSON_EXTRACT(object_info, \"$.machineId\")=" + strconv.Itoa(machineId) + " and JSON_EXTRACT(object_info, \"$.eventId\")=" + strconv.Itoa(eventId) + " order by JSON_EXTRACT(object_info, \"$.createdAt\") asc"
	var machineHMIGeneralObjects []component.GeneralObject

	dbConnection.Raw(hmiQuery).Scan(&machineHMIGeneralObjects)

	if len(machineHMIGeneralObjects) == 0 {
		return "", ""
	}
	var firstStartedRecordFound bool
	firstStartedRecordFound = false
	var actualStartTime, actualEndTime string
	for _, machineHMIGeneralObject := range machineHMIGeneralObjects {
		machineHMIInfo := MachineHMIInfo{}
		json.Unmarshal(machineHMIGeneralObject.ObjectInfo, &machineHMIInfo)
		if machineHMIInfo.HMIStatus == "started" && !firstStartedRecordFound {
			// this is the first start
			firstStartedRecordFound = true
			actualStartTime = machineHMIInfo.CreatedAt

		}
		if machineHMIInfo.HMIStatus == "stopped" {
			// this is the first start
			actualEndTime = machineHMIInfo.CreatedAt

		}
	}

	return actualStartTime, actualEndTime

}

func getHighestCycleCount(messages []Message) int {
	//Messages are in descending order
	//Three use cases
	// One is there is no 0 in message then it returns 0
	// Second One 0 in messages then pick max value and return it
	// Third multiple 0s in message then accumulate the may cycle count
	var messageBody *MachineDataInfo
	var accumulatedCycleCount int
	for index, message := range messages {
		err := json.Unmarshal(message.Body, &messageBody)

		if err != nil {
			continue
		}

		if messageBody != nil && messageBody.CycleCount != nil {
			if *messageBody.CycleCount == 0 {
				err = json.Unmarshal(messages[index+1].Body, &messageBody)
				if err != nil {
					continue
				}
				//Check nil for index+1 message data
				if messageBody != nil && messageBody.CycleCount != nil {
					accumulatedCycleCount += *messageBody.CycleCount
				}

			}
		}

	}
	return accumulatedCycleCount
}

type MaxMinCycleCount struct {
	MaxValue *int
	MinValue *int
}

type MachineHMIInfoResult struct {
	Created                 string
	EventId                 int
	Operator                int
	ReasonId                int
	MachineId               int
	ProductionOrderSourceId int
	Status                  string
	RejectQuantity          *int
	HmiStatus               string
	Id                      int
	TS                      DateTimeInfo
	Remark                  string
}

type DateTimeInfo struct {
	DateTimeString     string
	DateTime           time.Time
	DateTimeEpochMilli int64
	DateTimeEpoch      int64
	Error              error
}

type MachineDataInfo struct {
	CycleCount              *int `json:"cycle_count"`
	Auto                    *int `json:"auto"`
	MaintenanceDoorAlarm    *int `json:"maintenance_door_alarm"`
	SafetyDoorAlarm         *int `json:"safety_door_alarm"`
	EstopAlarm              *int `json:"estop_alarm"`
	TowerLightRed           *int `json:"tower_light_red"`
	MachineConnectionStatus *int `json:"machine_connection_status"`
	InTimestamp             *int `json:"in_timestamp"`
}
