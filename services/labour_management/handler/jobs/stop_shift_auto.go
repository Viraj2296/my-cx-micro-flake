package jobs

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (v *JobService) StopShiftAuto() {
	if v.PoolingInterval == 0 || v.PoolingInterval < 0 {
		v.PoolingInterval = 30
	}
	v.Logger.Info("machine stats pooling is starting up....", zap.Int("pooling_interval", v.PoolingInterval))
	var duration = time.Duration(v.PoolingInterval) * time.Second
	var condition = " object_info->>'$.shiftStatus' = 1"
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(const_util.ProjectID, const_util.ScheduleStatusPreferenceSix)

	for {

		err, listOfShift := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, condition)
		if err != nil {
			v.Logger.Error("error fetching list of shift", zap.Error(err))
		}

		for _, shift := range listOfShift {
			currentTime := time.Now().UTC()
			shiftInfo := database.GetShiftMasterInfo(shift.ObjectInfo)
			shiftEndTime := shiftInfo.ShiftEndTime
			endDateTime := shiftInfo.ShiftEndDate + "T" + shiftEndTime + ".000Z"
			// conditionTemplate := "TIME(JSON_UNQUOTE(JSON_EXTRACT(object_info, '$.shiftStartTime'))) >= '" + shiftStartTime + "'"
			// err, shiftTemplate := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, condition)
			// if err != nil || len(shiftTemplate) == 0 {
			// 	v.Logger.Error("error fetching list of shift", zap.Error(err))
			// }
			// templateObjectInfo := database.GetSShiftTemplateInfo(shiftTemplate[0].ObjectInfo)
			// templateShiftStartTime := templateObjectInfo.ShiftPeriod
			err, shiftSetting := database.Get(v.Database, const_util.LabourManagementSettingTable, 1)
			if err != nil {
				v.Logger.Error("error fetching shift setting", zap.Error(err))
			}
			shiftSettingInfo := database.GetLabourManagementSettingInfo(shiftSetting.ObjectInfo)
			endTimeDateValue := util.ConvertSingaporeTimeToUTC(endDateTime)
			endDateTimeStr, _ := time.Parse(const_util.ISOTimeLayout, endTimeDateValue)

			durationLimit, _ := time.ParseDuration(shiftSettingInfo.ShiftAutoStopTime)
			durationSinceEnd := currentTime.Sub(endDateTimeStr)
			if durationSinceEnd > durationLimit {
				var condition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(shift.Id)
				err, listOfAttendance := database.GetConditionalObjects(v.Database, const_util.LabourManagementAttendanceTable, condition)
				if err != nil {
					v.Logger.Error("error getting attendance", zap.String("error", err.Error()))

				}
				v.Logger.Info("checking out attendance size", zap.Int("size", len(listOfAttendance)))
				// now all the attendance, update the checkout date, and time
				for _, attendanceInterface := range listOfAttendance {
					attendance := database.GetAttendanceInfo(attendanceInterface.ObjectInfo)
					if attendance.CheckOutDate == "" {
						attendance.CheckOutDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
						attendance.CheckOutTime = util.GetCurrentTimeInSingapore("15:04:05")
						updateObject := make(map[string]interface{})
						updateObject["object_info"] = attendance.Serialised()
						err := database.Update(v.Database, const_util.LabourManagementAttendanceTable, attendanceInterface.Id, updateObject)
						if err != nil {
							// ideally all the shift employees should check out in complete shift, but log the error
							v.Logger.Error("error updating checkout time", zap.String("error", err.Error()))
						} else {
							v.Logger.Info("successfully checked to user id", zap.Int("user_id", attendanceInterface.Id))
						}
					} else {
						v.Logger.Warn("employee has already checkout, so skipping it", zap.Any("attendance", attendance))
					}

				}
				shiftInfo.ShiftStatus = const_util.ShiftStatusCompleted
				shiftInfo.CanShiftStop = false
				shiftInfo.CanShiftStart = false
				shiftInfo.CanCheckIn = false
				var updateObject = make(map[string]interface{})
				updateObject["object_info"] = shiftInfo.Serialised()
				err = database.Update(v.Database, const_util.LabourManagementShiftMasterTable, shift.Id, updateObject)
				if err != nil {
					v.Logger.Error("error updating shift master as complete", zap.String("error", err.Error()))

				}

				v.Logger.Info("handle stop shift is successfully processed", zap.Any("record_id", shift.Id))

				for _, eventId := range shiftInfo.ScheduledOrderEvents {
					var updatingData = make(map[string]interface{})
					updatingData["eventStatus"] = orderStatusId
					serializedObject, _ := json.Marshal(updatingData)
					v.Logger.Info("stopping the shift, stop  an individual scheduler event", zap.Int("event_id", eventId), zap.Any("event_object", string(serializedObject)))
					err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(const_util.ProjectID, eventId, serializedObject)
					if err != nil {
						v.Logger.Error("error updating scheduler order event in the stop shift", zap.Error(err))
						// we can not send an error inside loop, stay silent
					}
					err = v.createAssemblyHmiEntry(eventId, 0, "stopped")
					if err != nil {
						v.Logger.Error("error in creating stop hmi through labour management", zap.Error(err))
					}
				}

				// update the manpower flag to shift production now
				var shiftProductionCondition = " object_info->>'$.shiftId' = " + strconv.Itoa(shift.Id)
				err, shiftProdInterfaceObjects := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftProductionCondition)
				if err == nil {
					for _, shiftProdInterfaceObject := range shiftProdInterfaceObjects {
						shiftProductionInfo := database.GetShiftMasterProductionInfo(shiftProdInterfaceObject.ObjectInfo)
						shiftProductionInfo.CanUpdateManpower = false
						var updateShiftProdObject = make(map[string]interface{})
						updateShiftProdObject["object_info"] = shiftProductionInfo.Serialised()
						err = database.Update(v.Database, const_util.LabourManagementShiftProductionTable, shiftProdInterfaceObject.Id, updateShiftProdObject)
						if err != nil {
							v.Logger.Error("error updating scheduler production canUpdateManpower flag", zap.Error(err))
						}
					}
				}
				targetUsers := shiftSettingInfo.ShitAutoStopEmailNotificationUsers
				for _, userId := range targetUsers {
					userInfo := authService.GetUserInfoById(userId)
					var emailTemplate = v.EscalationEmailTemplate
					emailTemplate = strings.Replace(emailTemplate, "[USER]", userInfo.Username, 1)
					emailTemplate = strings.Replace(emailTemplate, "[SHIFTREFERENCE]", shiftInfo.ShiftReferenceId, 1)
					emailTemplate = strings.Replace(emailTemplate, "[DATEVALUE]", currentTime.Format(const_util.ISOTimeLayout), 1)
					v.sendEmail(userInfo.Email, emailTemplate)
				}

			}

		}

		time.Sleep(duration)
	}
}
func (v *JobService) sendEmail(targetEmail string, emailContent string) {
	var emailMessages []common.Message

	emailMessage := common.Message{
		To:          []string{targetEmail},
		SingleEmail: false,
		Subject:     "System Alert: Shift Automatically Stopped by MES",
		Body: map[string]string{
			"text/html": emailContent,
		},
		Info:          "",
		EmbeddedFiles: nil,
		AttachedFiles: nil,
	}

	emailMessages = append(emailMessages, emailMessage)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	err := notificationService.CreateMessages("906d0fd569404c59956503985b330132", emailMessages)
	if err != nil {
		v.Logger.Error("error creating notification messages", zap.Error(err))
	} else {
		v.Logger.Info("notification messages successfully created")
	}

}
func (v *JobService) createAssemblyHmiEntry(eventId, userId int, hmiStatus string) error {
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	err, assemblyScheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, eventId)
	if err != nil {
		return err
	}
	var scheduleOrderInfo = make(map[string]interface{})
	json.Unmarshal(assemblyScheduledOrderInterface.ObjectInfo, &scheduleOrderInfo)

	var machineId = util.InterfaceToInt(scheduleOrderInfo["machineId"])
	var timeNowStr = util.GetCurrentTime(const_util.ISOTimeLayout)
	var machineHmiInfo = map[string]interface{}{
		"eventId":   eventId,
		"createdAt": timeNowStr,
		"createdBy": userId,
		"hmiStatus": hmiStatus,
		"machineId": machineId,
	}
	err = machineService.CreateAssemblyHmiEntry(const_util.ProjectID, machineHmiInfo)
	return err
}
