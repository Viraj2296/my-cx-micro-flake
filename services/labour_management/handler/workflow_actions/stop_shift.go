package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) HandleStopShift(ctx *gin.Context) {
	v.Logger.Info("handle stop shift is received")
	recordId, dbConnection, userInfo, shiftMasterInfo := v.getBasicInfo(ctx)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(const_util.ProjectID, const_util.ScheduleStatusPreferenceSix)
	if shiftMasterInfo.CanShiftStop {

		if shiftMasterInfo.ShiftStatus == const_util.ShiftStatusActive {
			var condition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(recordId)
			err, listOfAttendance := database.GetConditionalObjects(dbConnection, const_util.LabourManagementAttendanceTable, condition)
			if err != nil {
				v.Logger.Error("error getting attendance", zap.String("error", err.Error()))
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
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
					err := database.Update(dbConnection, const_util.LabourManagementAttendanceTable, attendanceInterface.Id, updateObject)
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
			// update the shift status as active
			shiftMasterInfo.ShiftStatus = const_util.ShiftStatusCompleted
			shiftMasterInfo.CanShiftStop = false
			shiftMasterInfo.CanShiftStart = false
			shiftMasterInfo.CanCheckIn = false
			var updateObject = make(map[string]interface{})
			updateObject["object_info"] = shiftMasterInfo.Serialised()
			err = database.Update(dbConnection, const_util.LabourManagementShiftMasterTable, recordId, updateObject)
			if err != nil {
				v.Logger.Error("error updating shift master as complete", zap.String("error", err.Error()))
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

			v.Logger.Info("handle stop shift is successfully processed", zap.Any("record_id", recordId))
			// now update all the events as completed
			for _, eventId := range shiftMasterInfo.ScheduledOrderEvents {
				var updatingData = make(map[string]interface{})
				updatingData["eventStatus"] = orderStatusId
				serializedObject, _ := json.Marshal(updatingData)
				v.Logger.Info("stopping the shift, stop  an individual scheduler event", zap.Int("event_id", eventId), zap.Any("event_object", string(serializedObject)))
				err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(const_util.ProjectID, eventId, serializedObject)
				if err != nil {
					v.Logger.Error("error updating scheduler order event in the stop shift", zap.Error(err))
					// we can not send an error inside loop, stay silent
				}

				// Stop machine hmi
				err = v.createAssemblyHmiEntry(eventId, userInfo.UserId, "stopped")
				if err != nil {
					v.Logger.Error("error in creating stop hmi through labour management", zap.Error(err))
				}
			}

			// update the manpower flag to shift production now
			var shiftProductionCondition = " object_info->>'$.shiftId' = " + strconv.Itoa(recordId)
			err, shiftProdInterfaceObjects := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftProductionCondition)
			if err == nil {
				for _, shiftProdInterfaceObject := range shiftProdInterfaceObjects {
					shiftProductionInfo := database.GetShiftMasterProductionInfo(shiftProdInterfaceObject.ObjectInfo)
					shiftProductionInfo.CanUpdateManpower = false
					var updateShiftProdObject = make(map[string]interface{})
					updateShiftProdObject["object_info"] = shiftProductionInfo.Serialised()
					err = database.Update(dbConnection, const_util.LabourManagementShiftProductionTable, shiftProdInterfaceObject.Id, updateShiftProdObject)
					if err != nil {
						v.Logger.Error("error updating scheduler production canUpdateManpower flag", zap.Error(err))
					}
				}
			}
			ctx.JSON(http.StatusOK, response.GeneralResponse{
				Code:    0,
				Message: "Your shift has been stopped successfully",
			})

		} else if shiftMasterInfo.ShiftStatus == const_util.ShiftStatusCompleted {
			endDateTime := shiftMasterInfo.ShiftEndDate + "T" + shiftMasterInfo.ShiftEndTime + ".000Z"
			existingActionRemarks := shiftMasterInfo.ActionRemarks
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "SHIFT HAS BEEN STOPPED BY SYSTEM",
				UserId:        userInfo.UserId,
				Remarks:       "Shift has been stopped by system at " + endDateTime,
				ProcessedTime: getTimeDifference(util.InterfaceToString(shiftMasterInfo.CreatedAt)),
			})
			shiftMasterInfo.ActionRemarks = existingActionRemarks
			var updateObject = make(map[string]interface{})
			updateObject["object_info"] = shiftMasterInfo.Serialised()
			err := database.Update(dbConnection, const_util.LabourManagementShiftMasterTable, recordId, updateObject)
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}
		} else {
			v.Logger.Error("shift is not in active status to complete, shift status", zap.Any("status", shiftMasterInfo.ShiftStatus))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Shift Status",
					Description: "This operation is invalid, you can not complete this shift now, please report this to system admin",
				})
		}

	} else {
		v.Logger.Error("CanShiftStop is not true", zap.Any("status", shiftMasterInfo.ShiftStatus))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "This operation is invalid, you can not stop the shift at this movement, Please report this to system admin",
			})
	}

}
