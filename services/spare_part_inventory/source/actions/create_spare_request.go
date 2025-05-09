package actions

import (
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/spare_part_inventory/source/consts"
	"cx-micro-flake/services/spare_part_inventory/source/dto"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"net/http"

	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (v *Actions) CreateSparePartRequest(ctx *gin.Context) {
	v.Logger.Info("handle check in received")
	var userId = header_parser.GetRecordId(ctx)
	inventoryReceiveInRequest := dto.CreateSparePartRequest{}
	if err := ctx.ShouldBindBodyWith(&inventoryReceiveInRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", logging.Error(err))
		}
		return
	}

	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)
	for _, spareMaster := range inventoryReceiveInRequest.SpareParts {
		err, objectInterface := v.Repository.GetResource(consts.SparePartInventoryMasterComponent, spareMaster.SparePartId)
		if err != nil {
			v.Logger.Error("error check in error", logging.Error(err))
		}
		err, inventoryMasterInfo := models.GetSparePartInventoryMasterInfo(objectInterface.ObjectInfo)
		if err != nil {
			v.Logger.Error("error getting records", logging.String("error", err.Error()))
		}

		inventoryMasterInfo.OnHandQty -= spareMaster.Quantity

		serializedObject := inventoryMasterInfo.Serialised()
		err = v.Repository.UpdateResource(consts.SparePartInventoryMasterComponent, spareMaster.SparePartId, serializedObject, userId)

		if err != nil {
			v.Logger.Error("error check in error", zap.Error(err))
		}
	}

	spareRequest := models.SparePartInventoryRepairRequestInfo{}
	spareRequest.JobId = inventoryReceiveInRequest.JobId
	spareRequest.SpareParts = inventoryReceiveInRequest.SpareParts
	spareRequest.IsNeedSparePart = inventoryReceiveInRequest.IsNeedSparePart
	spareRequest.RequestStatus = "CREATED"
	spareRequest.MachineId = inventoryReceiveInRequest.MachineId

	var serialisedData = spareRequest.Serialised()

	err, _ := v.Repository.CreateResource(targetTable, serialisedData, userId)
	if err != nil {
		v.Logger.Error("error updating inventory master data", logging.Error(err))
		response.SendInternalSystemError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully Created the Inventory Request action",
	})
}
