package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type UpdateScheduleEvent struct {
	ShiftId                     int `json:"shiftId"`
	SchedulerEventId            int `json:"schedulerEventId"`
	ShiftTargetOutputPartTimers int `json:"shiftTargetOutputPartTimers"`
	ShiftActualOutputPartTimers int `json:"shiftActualOutputPartTimers"`
}

func (v *ActionService) UpdateShiftEvents(ctx *gin.Context) {
	v.Logger.Info("update product complete count received")

	updateScheduleEvent := UpdateScheduleEvent{}
	if err := ctx.ShouldBindBodyWith(&updateScheduleEvent, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("sending error failed", zap.Error(err))
		}
		return
	}

	// get the scheduled order events to this shift, and check the requested event id is present
	err, shiftMasterInterface := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, updateScheduleEvent.ShiftId)

	if err != nil {
		response.DispatchDetailedError(ctx, common.ObjectNotFound,
			&response.DetailedError{
				Header:      const_util.GetError(common.ObjectNotFoundError).Error(),
				Description: "Request shift is not available in the system.",
			})
		return
	}
	shiftInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
	if shiftInfo.ShiftStatus == const_util.ShiftStatusCompleted {
		response.DispatchDetailedError(ctx, common.ObjectNotFound,
			&response.DetailedError{
				Header:      "Invalid Operation",
				Description: "This shift has already been completed, updating shift related scheduled event is not possible now.",
			})
		return
	}

	var condition = " object_info->>'$.shiftId' = " + strconv.Itoa(updateScheduleEvent.ShiftId) + " AND object_info->>'$.scheduledEventId' = " + strconv.Itoa(updateScheduleEvent.SchedulerEventId)
	err, listOfShiftProduction := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, condition)
	if len(listOfShiftProduction) == 0 {
		response.DispatchDetailedError(ctx, common.ObjectNotFound,
			&response.DetailedError{
				Header:      "Invalid Shift Production Order",
				Description: "System couldn't able to find the scheduled order event related to this shift, please report to system administrator",
			})
		return
	}

	shiftMasterProductionInfo := database.GetShiftMasterProductionInfo((listOfShiftProduction)[0].ObjectInfo)
	shiftMasterProductionInfo.ShiftActualOutputPartTimer = updateScheduleEvent.ShiftActualOutputPartTimers
	shiftMasterProductionInfo.ShiftTargetOutputPartTimer = updateScheduleEvent.ShiftTargetOutputPartTimers
	var updatingFields = make(map[string]interface{})
	updatingFields["object_info"] = shiftMasterProductionInfo.Serialised()
	err = database.Update(v.Database, const_util.LabourManagementShiftProductionTable, (listOfShiftProduction)[0].Id, updatingFields)
	if err != nil {
		v.Logger.Error("error updating shift production", zap.Any("error", err.Error()))
		response.SendInternalSystemError(ctx)
		return
	}

	v.Logger.Info("successfully updated the master info")
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the scheduled order event",
	})

}
