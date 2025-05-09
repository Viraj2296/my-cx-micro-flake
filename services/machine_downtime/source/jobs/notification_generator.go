package jobs

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"go.cerex.io/transcendflow/component"
	"go.cerex.io/transcendflow/const_util"
	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/orm"
	"go.cerex.io/transcendflow/util"
	"go.uber.org/zap"
)

func (v *Jobs) GetRandomJobUser() int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	if len(v.JobAssigningUsers) > 0 {
		randomIndex := r.Intn(len(v.JobAssigningUsers))
		selectedUser := v.JobAssigningUsers[randomIndex]
		return selectedUser
	} else {
		return -1
	}
}

type SignalMachine struct {
	Id                  int   `json:"id"`
	SignalGeneratedTime int64 `json:"signalGeneratedTime"`
}

func (v *Jobs) GenerateNotificationMessage() {
	// Get necessary services
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	var duration = time.Duration(v.DowntimeConfig.JobServiceConfig.MachineHelpSignalPollingInterval) * time.Second
	for {
		// make sure we sleep in the earliest to handle continue
		time.Sleep(duration)
		v.Logger.Info("Checking for signals to generate message", logging.String("last_processed_time", util.ConvertEpochToDateTime(int(v.LastProcessedMachineHelpSignalTime))), logging.Int64("last_processed_time", v.LastProcessedMachineHelpSignalTime))
		err, listOfSignalledMachines := machineService.GetListOfMachinesNeededHelp(v.LastProcessedMachineHelpSignalTime)
		if err != nil {
			v.Logger.Error("Error getting machine help signal", zap.Error(err))
			continue
		}
		v.Logger.Info("list of signaled machines ", zap.Any("machines", listOfSignalledMachines))

		for _, machineDetailInterface := range listOfSignalledMachines {
			var isCreateFaultFailed = false
			signalMachine := SignalMachine{}
			err := json.Unmarshal(machineDetailInterface, &signalMachine)
			if err == nil {
				err, assemblyMachineObject := machineService.GetAssemblyMachineInfoById(signalMachine.Id)
				if err == nil {
					assemblyMachineMasterInfo, _ := dto.GetAssemblyMachineMasterInfo(assemblyMachineObject.ObjectInfo)
					v.Logger.Info("Help signalled machine master", logging.Any("assembly_machine_master", assemblyMachineMasterInfo), logging.Any("job_assigning_user", v.JobAssigningUsers))

					var referenceData = make(map[string]interface{})
					// initially assign the job to no body, assigning would happen during the check-in process
					var downtimeId = v.CreateFault(signalMachine.Id, 0)
					// now sendd the notification to  all the job assigning users
					for _, jobAssignUserId := range v.JobAssigningUsers {
						assignedUserDeviceToken := authService.GetUserOneSignalSubscriptionIds(jobAssignUserId)

						if downtimeId == -1 {
							v.Logger.Error("error creating fault")
							isCreateFaultFailed = true
							continue
						}
						v.Logger.Info("created a fault ", zap.Int("downtime_id", downtimeId))
						referenceData["downtimeId"] = downtimeId
						pushNotificationMessage := common.PushNotificationMessage{
							IncludePlayerIDs: assignedUserDeviceToken,
							Headings:         map[string]string{"en": "Assembly Machine is Down"},
							Contents:         map[string]string{"en": "Machine [" + assemblyMachineMasterInfo.Description + "] is requesting help. Please attend it, and check in"}, //Todo this newmachineid is still not update as null value in the push notification message, need to check this
							Data:             nil,
							DeliveryStatus:   "pending",
							RetryCount:       0,
							ReferenceData:    referenceData,
						}

						pushNotificationData, err := json.Marshal(pushNotificationMessage)
						if err != nil {
							v.Logger.Error("Failed to marshal push notification message", zap.Error(err))
							continue
						}
						notificationId, err := notificationService.CreatePushNotification(consts.ProjectID, pushNotificationData)
						if err != nil {
							v.Logger.Error("error creating push notification for general users", zap.Error(err))
						} else {
							//add the notification id to the user using auth srvice
							v.Logger.Info("Push notification has been successfully created", zap.Any("push_notification_id", notificationId))
							err = authService.AddPushNotificationIds(jobAssignUserId, notificationId)
							if err != nil {
								v.Logger.Error("Error adding push notification id to user", zap.Error(err))
							} else {
								v.Logger.Info("Successfully added notification ID to user profile")
							}

						}
					}

				} else {
					v.Logger.Error("error getting assembly machine info", zap.Error(err))
				}
				if !isCreateFaultFailed {

					signalHistoryInfo := models.SignalHistoryInfo{
						CreatedAt:               signalMachine.SignalGeneratedTime,
						SignalMachine:           signalMachine.Id,
						IsNotificationProcessed: false,
					}
					v.Logger.Info("signal history info", zap.Any("signal_history_info", signalHistoryInfo.Serialised()))
					err, recordId := orm.CreateFromResource(v.Database, consts.MachineDownTimeSignalHistoryTable, 1, signalHistoryInfo.Serialised())

					if err != nil {
						v.Logger.Error("error creates downtime signal history", zap.Error(err))
					} else {
						v.Logger.Info("Successfully created downtime signal history", zap.Any("record", recordId))
						v.LastProcessedMachineHelpSignalTime = signalMachine.SignalGeneratedTime
					}
				}

			} else {
				v.Logger.Error("error un-marshalling machine detail", zap.Error(err))
			}
		}

	}
}

func (v *Jobs) CreateFault(machineId int, randomJobAssignUserId int) int {
	machineDowntimeInfo := models.MachineDowntimeInfo{}
	machineDowntimeInfo.MachineId = machineId
	machineDowntimeInfo.CheckInDate = ""
	machineDowntimeInfo.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	machineDowntimeInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	machineDowntimeInfo.CheckInTime = ""
	machineDowntimeInfo.CreatedBy = 1
	machineDowntimeInfo.LastUpdatedBy = 1
	machineDowntimeInfo.CheckInUserId = 0
	machineDowntimeInfo.CanCheckIn = true
	machineDowntimeInfo.CanCheckOut = false
	machineDowntimeInfo.AssignedUserId = randomJobAssignUserId
	machineDowntimeInfo.Status = consts.DowntimeStatus_Fault_Reportd
	machineDowntimeInfo.Remarks = ""
	machineDowntimeInfo.CanCancel = true
	machineDowntimeInfo.JobReferenceId = "DT"
	machineDowntimeInfo.CheckOutDate = ""
	machineDowntimeInfo.CheckOutTime = ""
	machineDowntimeInfo.ObjectStatus = const_util.ObjectStatusActive
	machineDowntimeInfo.ActionRemarks = append(machineDowntimeInfo.ActionRemarks, component.ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(util.ISOTimeLayout),
		Status:        "Fault created by MES",
		UserId:        1,
		Remarks:       "MES has received the signal to create the fault",
		ProcessedTime: util.GetTimeDifference(util.InterfaceToString(machineDowntimeInfo.CreatedAt)),
	})
	serialisedData, _ := machineDowntimeInfo.Serialised()
	err, recordId := orm.CreateFromResource(v.Database, consts.MachineDownTimeMasterTable, 1, serialisedData)
	if err != nil {
		v.Logger.Error("create from general object error", zap.Error(err))
		return -1
	}

	// make the machine into maintenance mode.
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	err = machineService.MoveAssemblyMachineToMaintenance(consts.ProjectID, machineId)
	if err != nil {
		v.Logger.Error("move assembly machine to maintenance error", zap.Error(err))
	}

	jobNumber := getJobNumber(recordId)
	err, downtimeMasterObject := orm.Get(v.Database, consts.MachineDownTimeMasterTable, recordId)
	if err == nil {
		downtimeObjectInfo := models.GetMachineDowntimeInfo(downtimeMasterObject.ObjectInfo)
		if downtimeObjectInfo != nil {
			downtimeObjectInfo.JobReferenceId = jobNumber
			if serializedData, err := downtimeObjectInfo.Serialised(); err == nil {
				err := orm.UpdateSerialisedResourceFromId(v.Database, consts.MachineDownTimeMasterTable, recordId, 1, serializedData)
				if err != nil {
					v.Logger.Error("error updating job reference ID", zap.Error(err))
				}
			} else {
				v.Logger.Error("error serialising data ", zap.Error(err))
			}
		}

	} else {
		v.Logger.Error("error getting job reference ID", zap.Error(err))
	}

	return recordId
}
func getJobNumber(num int) string {
	return "DT" + fmt.Sprintf("%08d", num)
}
