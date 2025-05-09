package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

type UpdateProductComplete struct {
	ActualOutput     int    `json:"actualOutput"`
	Remarks          string `json:"remarks"`
	SchedulerEventId int    `json:"schedulerEventId"`
}

func (v *ActionService) UpdateProductCompleteCount(ctx *gin.Context) {
	v.Logger.Info("update product complete count received")
	shiftMasterId := util.GetRecordId(ctx)
	var userId = header_parser.GetUserId(ctx)
	// get the scheduled order events to this shift, and check the requested event id is present
	err, shiftMasterInterface := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, shiftMasterId)

	if err != nil {
		response.DispatchDetailedError(ctx, common.ObjectNotFound,
			&response.DetailedError{
				Header:      const_util.GetError(common.ObjectNotFoundError).Error(),
				Description: "Request shift is not available in the system.",
			})
		return
	}
	updateProductComplete := UpdateProductComplete{}
	if err := ctx.ShouldBindBodyWith(&updateProductComplete, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	shiftInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
	if util.IsElementExistIntArray(shiftInfo.ScheduledOrderEvents, updateProductComplete.SchedulerEventId) {
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
		v.Logger.Info("received update product complete", zap.Any("body", updateProductComplete))
		err, scheduledEventInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, updateProductComplete.SchedulerEventId)
		if err != nil {
			v.Logger.Error("error getting production order information", zap.String("error", err.Error()))

		}
		var scheduledEventInfo = GetAssemblyScheduledOrderEventInfo(scheduledEventInterface.ObjectInfo)
		var scheduledQty = scheduledEventInfo.ScheduledQty
		var progressDone = 0.0
		if updateProductComplete.ActualOutput < scheduledQty {
			progressDone = float64((updateProductComplete.ActualOutput / scheduledQty) * 100.0)
		}

		// call the assembly module and update the fields
		var updatingFields = make(map[string]interface{})
		updatingFields["completedQty"] = updateProductComplete.ActualOutput
		updatingFields["percentDone"] = progressDone
		updatingFields["remarks"] = updateProductComplete.Remarks
		serializedObject, _ := json.Marshal(updatingFields)
		err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(const_util.ProjectID, updateProductComplete.SchedulerEventId, serializedObject)
		if err != nil {
			v.Logger.Error("error updating assemble scheduler order fields", zap.String("error", err.Error()))
		}
		type updateAssemblyManualOrderQuantityRequest struct {
			EventId           int `json:"eventId"`
			CompletedQuantity int `json:"completedQuantity"`
			RequestBy         int `json:"requestBy"`
		}
		var manualUpdateQuantityRequest updateAssemblyManualOrderQuantityRequest
		manualUpdateQuantityRequest.RequestBy = userId
		manualUpdateQuantityRequest.CompletedQuantity = updateProductComplete.ActualOutput
		manualUpdateQuantityRequest.EventId = updateProductComplete.SchedulerEventId
		serialisedUpdatedQuantity, _ := json.Marshal(manualUpdateQuantityRequest)

		_, err = productionOrderInterface.UpdateAssemblyManualOrderCompletedQuantity(serialisedUpdatedQuantity)
		if err != nil {
			v.Logger.Error("error updating product completed quantity", zap.String("error", err.Error()))
		}

		v.Logger.Info("successfully updated the event")
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Successfully updated the scheduled order event",
		})
	} else {
		v.Logger.Error("requested event is not present in the scheduled events, check the shift master", zap.Any("event", updateProductComplete.SchedulerEventId))
		response.DispatchDetailedError(ctx, common.ObjectNotFound,
			&response.DetailedError{
				Header:      const_util.GetError(common.ObjectNotFoundError).Error(),
				Description: "Invalid scheduled order event, Please report this error to system admin",
			})
		return
	}

}
