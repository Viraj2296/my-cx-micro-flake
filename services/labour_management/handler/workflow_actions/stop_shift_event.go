package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StopSchedulerEvent struct {
	EventId int `json:"eventId"`
}

func (v *ActionService) HandleStopShiftEvent(ctx *gin.Context) {
	v.Logger.Info("handle stop shift event is received")
	stopSchedulerEvent := StopSchedulerEvent{}
	if err := ctx.ShouldBindBodyWith(&stopSchedulerEvent, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, scheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, stopSchedulerEvent.EventId)
	if err != nil {
		v.Logger.Error("error getting scheduler order information", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InterModuleCommunicationProblem,
			&response.DetailedError{
				Header:      "Error Getting Scheduler Event",
				Description: "Internal system error during update [" + err.Error() + "]",
			})
		return
	}
	scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(scheduledOrderInterface.ObjectInfo)
	err, productionOrderData := productionOrderInterface.GetAssemblyProductionOrderById(scheduledOrderInfo.EventSourceId)
	if err != nil {
		v.Logger.Error("error getting production order information", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.UnmarshalingError,
			&response.DetailedError{
				Header:      "Error Decoding Assembly Order",
				Description: "Internal system error during update [" + err.Error() + "]",
			})
		return
	}
	var productionOrderInfo = GetAssemblyProductionOrderInfo(productionOrderData.ObjectInfo)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	// get the machine details
	err, machineGeneralComponent := machineService.GetAssemblyMachineInfoById(productionOrderInfo.MachineId)
	if err != nil {
		v.Logger.Error("error getting machine  information", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InterModuleCommunicationProblem,
			&response.DetailedError{
				Header:      "Error Getting Machine Information",
				Description: "Internal system error during update [" + err.Error() + "]",
			})
		return
	}
	machineInfo := GetAssemblyMachineMasterInfo(machineGeneralComponent.ObjectInfo)

	recordId, dbConnection, userInfo, shiftMasterInfo := v.getBasicInfo(ctx)
	if shiftMasterInfo.CanShiftStop {

		if shiftMasterInfo.ShiftStatus == const_util.ShiftStatusActive {

			productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
			orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(const_util.ProjectID, const_util.ScheduleStatusPreferenceSix)
			var updatingData = make(map[string]interface{})
			updatingData["eventStatus"] = orderStatusId
			serializedObject, _ := json.Marshal(updatingData)
			v.Logger.Info("stopping an individual scheduler event", zap.Int("event_id", stopSchedulerEvent.EventId), zap.Any("event_object", string(serializedObject)))
			err := productionOrderInterface.UpdateAssemblyScheduledOrderFields(const_util.ProjectID, stopSchedulerEvent.EventId, serializedObject)
			if err != nil {
				v.Logger.Error("error updating scheduler order event", zap.Error(err))
				response.DispatchDetailedError(ctx, common.InterModuleCommunicationProblem,
					&response.DetailedError{
						Header:      "Error Update Event Status",
						Description: "Internal system error during update [" + err.Error() + "]",
					})
				return
			}
			// get the line for this event

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
				attendanceInfo := database.GetAttendanceInfo(attendanceInterface.ObjectInfo)
				if util.HasInt(machineInfo.AssemblyLineOption, attendanceInfo.ManufacturingLines) {

					// if the user is checked in for multiple lines, then we need to create the separate attendance record, and update timing.
					if len(attendanceInfo.ManufacturingLines) > 1 {

						var modifiedLines = make([]int, 0)
						modifiedLines = append(modifiedLines, machineInfo.AssemblyLineOption)
						newAttendanceRecord := database.AttendanceInfo{
							ShiftResourceId:    attendanceInfo.ShiftResourceId,
							CheckInDate:        util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
							CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
							LastUpdatedAt:      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
							LastUpdatedBy:      userInfo.UserId,
							CheckInTime:        util.GetCurrentTimeInSingapore("15:04:05"),
							CreatedBy:          userInfo.UserId,
							UserResourceId:     attendanceInfo.UserResourceId,
							CheckOutDate:       util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
							CheckOutTime:       util.GetCurrentTimeInSingapore("15:04:05"),
							ManufacturingLines: modifiedLines,
						}
						generalObject := component.GeneralObject{ObjectInfo: newAttendanceRecord.Serialised()}
						err, attendanceRecordId := database.CreateFromGeneralObject(dbConnection, const_util.LabourManagementAttendanceTable, generalObject)
						if err != nil {
							v.Logger.Error("error creating attendance record for checkout", zap.Error(err))
							return
						}
						v.Logger.Info("new attendance record is created ", zap.Int("record_id", attendanceRecordId))

					} else {
						attendanceInfo.LastUpdatedBy = userInfo.UserId
						attendanceInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
						attendanceInfo.CheckOutDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
						attendanceInfo.CheckOutTime = util.GetCurrentTimeInSingapore("15:04:05")
						updateObject := make(map[string]interface{})
						updateObject["object_info"] = attendanceInfo.Serialised()
						err := database.Update(dbConnection, const_util.LabourManagementAttendanceTable, attendanceInterface.Id, updateObject)
						if err != nil {
							// ideally all the shift employees should check out in complete shift, but log the error
							v.Logger.Error("error updating checkout time", zap.String("error", err.Error()))
						} else {
							v.Logger.Info("successfully checked to user id", zap.Int("user_id", attendanceInterface.Id))
						}
					}

				} else {
					v.Logger.Warn("line is not found in the checked in lines", zap.Int("assembly_line", machineInfo.AssemblyLineOption))
				}

			}

			//Stop machine hmi
			err = v.createAssemblyHmiEntry(stopSchedulerEvent.EventId, userInfo.UserId, "stopped")
			if err != nil {
				v.Logger.Error("error in creating stop hmi through labour management", zap.Error(err))
			}

			v.Logger.Info("handle Stop Shift event is successfully processed", zap.Any("record_id", recordId))
			ctx.JSON(http.StatusOK, response.GeneralResponse{
				Code:    0,
				Message: "Your shift has been stopped successfully",
			})

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
