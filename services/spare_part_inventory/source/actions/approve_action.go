package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/error_util"
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

func (v *Actions) ApproveAction(ctx *gin.Context) {
	v.Logger.Info("Cancel request received")

	var recordId = header_parser.GetRecordId(ctx)
	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)
	var userId = header_parser.GetUserId(ctx)
	err, objectInterface := v.Repository.GetResource(targetTable, recordId)

	approvePartPayload := dto.ApprovePartList{}
	if err := ctx.ShouldBindBodyWith(&approvePartPayload, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Spare Part Id",
				Description: "The Spare Part Id is not available in the system, record might have archived or removed from system",
			})
		return
	}

	err, sparePartInventoryRequestInfo := models.GetSparePartInventoryRepairRequestInfo(objectInterface.ObjectInfo)
	if err != nil {
		v.Logger.Error("cannot unmarshal", logging.String("error", err.Error()))
		error_util.SendUnmarshlingFailed(ctx)
		return
	}
	for _, spareMaster := range approvePartPayload.SparePartList {
		err, objectInterface := v.Repository.GetResource(consts.SparePartInventoryMasterComponent, spareMaster.SparePartId)
		if err != nil {
			v.Logger.Error("error check in error", logging.Error(err))
		}
		err, inventoryMasterInfo := models.GetSparePartInventoryMasterInfo(objectInterface.ObjectInfo)
		if err != nil {
			v.Logger.Error("error getting records", logging.String("error", err.Error()))
		}
		for _, dbPart := range sparePartInventoryRequestInfo.SpareParts {
			if dbPart.SparePartId == spareMaster.SparePartId {
				var different = 0
				if spareMaster.Quantity >= dbPart.Quantity {
					different = spareMaster.Quantity - dbPart.Quantity
					inventoryMasterInfo.OnHandQty -= different
				} else {
					different = dbPart.Quantity - spareMaster.Quantity
					inventoryMasterInfo.OnHandQty += different
				}
			}
		}

		// inventoryMasterInfo.OnHandQty -= spareMaster.Quantity

		serializedObject := inventoryMasterInfo.Serialised()
		err = v.Repository.UpdateResource(consts.SparePartInventoryMasterComponent, spareMaster.SparePartId, serializedObject, userId)

		if err != nil {
			v.Logger.Error("error check in error", zap.Error(err))
		}
	}

	sparePartInventoryRequestInfo.RequestStatus = "APPROVED"
	sparePartInventoryRequestInfo.SpareParts = approvePartPayload.SparePartList

	var serialisedData = sparePartInventoryRequestInfo.Serialised()
	err = v.Repository.UpdateResource(targetTable, recordId, serialisedData, userId)

	if err != nil {
		v.Logger.Error("error check in error", zap.Error(err))
		response.SendInternalSystemError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully Approve the Request",
	})
}
