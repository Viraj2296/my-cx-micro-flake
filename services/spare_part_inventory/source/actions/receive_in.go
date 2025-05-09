package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/spare_part_inventory/source/dto"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"net/http"

	"go.cerex.io/transcendflow/logging"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (v *Actions) InventoryReceiveIn(ctx *gin.Context) {
	v.Logger.Info("handle check in received")
	var userId = header_parser.GetRecordId(ctx)
	inventoryReceiveInRequest := dto.InventoryReceiveInputRequest{}
	if err := ctx.ShouldBindBodyWith(&inventoryReceiveInRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", logging.Error(err))
		}
		return
	}
	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)

	err, objectInterface := v.Repository.GetResource(targetTable, inventoryReceiveInRequest.ResourceId)

	if err != nil {
		v.Logger.Error("error check in error", logging.Error(err))
		response.SendInternalSystemError(ctx)
	}

	if !common.ValidateObjectStatus(objectInterface.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Resource Status",
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}
	err, inventoryMasterInfo := models.GetSparePartInventoryMasterInfo(objectInterface.ObjectInfo)
	if err != nil {
		v.Logger.Error("error getting records", logging.String("error", err.Error()))
		error_util.SendUnmarshlingFailed(ctx)
		return
	}

	inventoryMasterInfo.OnHandQty += inventoryReceiveInRequest.Quantity

	var serialisedData = inventoryMasterInfo.Serialised()

	err = v.Repository.UpdateResource(targetTable, inventoryReceiveInRequest.ResourceId, serialisedData, userId)
	if err != nil {
		v.Logger.Error("error updating inventory master data", logging.Error(err))
		response.SendInternalSystemError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the Inventory receive action",
	})
}
