package jobs

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/machines/handler/const_util"
	"cx-micro-flake/services/machines/handler/database"
	"cx-micro-flake/services/machines/handler/model"
	"encoding/json"
	"strconv"
	"time"

	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/util"
	"go.uber.org/zap"
)

func (v *JobService) GenerateHelpSignalHistory() {
	tDuration, err := time.ParseDuration(v.JobConfig.AssemblyHelpSignalGenerationPeriod)
	if err != nil {
		tDuration, err = time.ParseDuration("30s")
		if err != nil {
			v.Logger.Error("error parsing pooling interval, using default value", zap.String("error", err.Error()))
			tDuration = 10 * time.Second
		}
	}
	for {

		time.Sleep(tDuration)
		v.Logger.Info("generating help signal start, passing timestamp", logging.Any("ts", v.LastProcessedHelpSignalTs), logging.String("converted_time", util.ConvertEpochToDateTime(v.LastProcessedHelpSignalTs)))
		var recordCount = database.Count(v.Database, const_util.AssemblyHelpSignalProcessedTimeTable)
		var isProcessedTimeRecordAvailable = false
		if recordCount > 0 {
			isProcessedTimeRecordAvailable = true
		}
		rows, err := v.Database.Raw(const_util.SelectAssemblyHistoryMessage, v.LastProcessedHelpSignalTs).Rows()
		if err != nil {
			v.Logger.Error("generate help signal history", zap.Error(err))
			continue
		}
		for rows.Next() {
			var bodyData []byte
			var messageTs int
			if err := rows.Scan(&bodyData, &messageTs); err != nil {
				v.Logger.Error("Failed to scan message body", zap.Error(err))
				continue
			}

			var assemblyMessage database.AssemblyMessageBody
			if err := json.Unmarshal(bodyData, &assemblyMessage); err != nil {
				v.Logger.Error("Failed to unmarshal JSON into AssemblyMessageBody", zap.Error(err))
				continue
			}
			v.Logger.Info("processing message for signal history", zap.Any("assembly_message", assemblyMessage))
			v.LastProcessedHelpSignalTs = messageTs

			// insert the message processed time
			if !isProcessedTimeRecordAvailable {
				var assemblyHelpSignalProcessedTimeInfo = model.AssemblyHelpSignalProcessedTimeInfo{
					ProcessedMessageTs: messageTs,
					LastUpdatedAt:      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
					CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
				}
				serialisedData, _ := assemblyHelpSignalProcessedTimeInfo.Serialised()
				err, _ := database.CreateFromGeneralObject(v.Database, const_util.AssemblyHelpSignalProcessedTimeTable, component.GeneralObject{ObjectInfo: serialisedData})
				if err != nil {
					v.Logger.Error("failed to create assembly help signal processed time data", zap.Error(err))
				} else {
					isProcessedTimeRecordAvailable = true
				}
			} else {
				// update it
				_, assemblyHelpSignalObject := database.Get(v.Database, const_util.AssemblyHelpSignalProcessedTimeTable, 1)
				var assemblyHelpSignalProcessedTimeInfo = model.GetAssemblyHelpSignalProcessedTimeInfo(assemblyHelpSignalObject.ObjectInfo)
				assemblyHelpSignalProcessedTimeInfo.ProcessedMessageTs = messageTs
				assemblyHelpSignalProcessedTimeInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
				serialisedData, _ := assemblyHelpSignalProcessedTimeInfo.Serialised()
				var updatingData = make(map[string]interface{})
				updatingData["object_info"] = serialisedData
				err := database.Update(v.Database, const_util.AssemblyHelpSignalProcessedTimeTable, 1, updatingData)
				if err != nil {
					v.Logger.Error("Failed to update assemblyHelpSignalProcessedTime", zap.Error(err))
				} else {
					v.Logger.Info("successfully updated the assembly signal processed time", zap.Any("assembly_message", string(serialisedData)))
				}
			}

			// Define station call mappings for each LXXStnCall and LXXOpId pair
			stationCalls := map[string]*string{
				"L01": assemblyMessage.L01StnCall,
				"L02": assemblyMessage.L02StnCall,
				"L03": assemblyMessage.L03StnCall,
				"L04": assemblyMessage.L04StnCall,
				"L05": assemblyMessage.L05StnCall,
				"L06": assemblyMessage.L06StnCall,
			}

			// Loop through each station call to identify active help buttons
			for lineFlag, stnCall := range stationCalls {
				if stnCall != nil {
					var stationCallString = *stnCall
					v.Logger.Info("Station Call", zap.String(lineFlag+"_StationCall", stationCallString))
					// Find the index positions where the call signal is "1"
					for indexPosition, char := range stationCallString {
						var helpButtonNumber = indexPosition + 1
						if char == '1' {
							var machineId int
							err := v.Database.Raw(const_util.SelectAssemblyMasterFromMessageFlag, lineFlag, helpButtonNumber).Scan(&machineId).Error
							if err != nil {
								v.Logger.Error("Failed to retrieve machine info", zap.Error(err))
								continue
							}
							if machineId == 0 {
								v.Logger.Warn("invalid machine ID during 1 flag checks, so skipping..", logging.String("line_flag", lineFlag), logging.Int("help_button_number", helpButtonNumber))
								continue
							}
							v.Logger.Info("Processing signal for machine id", zap.Int("machine_id", machineId))
							var numberOfRecords = database.CountByCondition(v.Database, const_util.AssemblyMachineHelpSignalViewTable, " id = "+strconv.Itoa(machineId))
							if numberOfRecords == 0 {
								// insert record
								var assemblyMachineHelpSignalViewInfo = model.AssemblyMachineHelpSignalViewInfo{
									CreatedAt:             util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
									LastUpdatedAt:         util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
									HelpButtonPressedTime: messageTs,
								}
								serialisedData, err := assemblyMachineHelpSignalViewInfo.Serialised()
								if err != nil {
									v.Logger.Error("error serialising machine help signal info", logging.Error(err))
								} else {
									err, _ = database.CreateFromGeneralObject(v.Database, const_util.AssemblyMachineHelpSignalViewTable, component.GeneralObject{Id: machineId, ObjectInfo: serialisedData})
									if err != nil {
										v.Logger.Error("error creating help signal history for machine ID ", logging.Error(err), logging.Int("machine_id", machineId))
									} else {
										v.Logger.Info("successfully created the help signal history for machine ID ", logging.Int("machine_id", machineId))
									}
								}

							} else {
								// update it
								err, c := database.Get(v.Database, const_util.AssemblyMachineHelpSignalViewTable, machineId)
								if err != nil {
									v.Logger.Error("error getting machine ID to update", logging.Error(err))
								} else {
									var assemblyMachineSignalViewInfo = model.GetAssemblyMachineHelpSignalViewInfo(c.ObjectInfo)
									assemblyMachineSignalViewInfo.HelpButtonPressedTime = messageTs
									assemblyMachineSignalViewInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
									serialisedData, err := assemblyMachineSignalViewInfo.Serialised()
									if err != nil {
										v.Logger.Error("error serialising machine help signal info", logging.Error(err))
									} else {
										var updatingData = make(map[string]interface{})
										updatingData["object_info"] = serialisedData
										err := database.Update(v.Database, const_util.AssemblyMachineHelpSignalViewTable, machineId, updatingData)
										if err != nil {
											v.Logger.Error("error updating help signal history for machine ID ", logging.Error(err), logging.Int("machine_id", machineId))
										} else {
											v.Logger.Info("successfully updated the help signal history for machine ID ", logging.Int("machine_id", machineId))
										}
									}
								}
							}
						}

					}
				}
			}
		}

	}
}
