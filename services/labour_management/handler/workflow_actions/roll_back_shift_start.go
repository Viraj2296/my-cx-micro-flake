package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (v *ActionService) RollBackShift(ctx *gin.Context) {
	// make all the scheduled events in to scheduld stage again
	// make all the HMI entry status archived
	// make the shift master canStart true, canCheckIn True again
	v.Logger.Info("handle start shift is received")
	recordId, dbConnection, userInfo, shiftMasterInfo := v.getBasicInfo(ctx)
	// Order status id is confirmed
	var orderStatusId = 4

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		v.Logger.Error("invalid payload, return error now", zap.String("error", err.Error()))
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	shiftStart, err := time.Parse("2006-01-02", shiftMasterInfo.ShiftStartDate)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Input",
				Description: "Invalid shift start date",
			})
		return
	}

	// Get today's date
	today := time.Now()

	// Check if the input date matches today's date
	isToday := shiftStart.Year() == today.Year() &&
		shiftStart.Month() == today.Month() &&
		shiftStart.Day() == today.Day()

	if !isToday {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "User can only roll back today shift",
			})
		return
	}

	if shiftMasterInfo.CanRollBack {
		// update the shift status as active
		shiftMasterInfo.ShiftStatus = const_util.ShiftStatusPending
		shiftMasterInfo.CanShiftStop = false
		shiftMasterInfo.CanShiftStart = true
		shiftMasterInfo.CanCheckIn = true
		shiftMasterInfo.LastUpdatedBy = userInfo.UserId
		//shiftMasterInfo.CanRollBack = false
		shiftMasterInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		existingActionRemarks := shiftMasterInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "SHIFT MASTER ROLL BACK",
			UserId:        userInfo.UserId,
			Remarks:       util.InterfaceToString(returnFields["remark"]),
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
		var productionOrderInterface = common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
		var machineService = common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
		// update all the events as
		for _, eventId := range shiftMasterInfo.ScheduledOrderEvents {
			err = productionOrderInterface.UpdateAssemblyOrderPreferenceLevel(const_util.ProjectID, userInfo.UserId, eventId, orderStatusId)

			if err != nil {
				v.Logger.Error("error updating assembly schedule order status", zap.Error(err))
				continue
			}

			//	update hmi
			err = machineService.ArchivedAssemblyHmiEntry(const_util.ProjectID, eventId)
			if err != nil {
				v.Logger.Error("error in archiving hmi through labour management", zap.Error(err))
			}
		}
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

	v.Logger.Info("handle rool back shift master is successfully processed", zap.Any("record_id", recordId))
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Your shift has been rolled back successfully",
	})

}
