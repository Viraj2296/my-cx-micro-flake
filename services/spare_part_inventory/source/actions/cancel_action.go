package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"net/http"

	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func (v *Actions) CancelAction(ctx *gin.Context) {
	v.Logger.Info("Cancel request received")

	var recordId = header_parser.GetRecordId(ctx)
	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)
	var userId = header_parser.GetUserId(ctx)
	err, objectInterface := v.Repository.GetResource(targetTable, recordId)

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

	sparePartInventoryRequestInfo.RequestStatus = "CANCELLED"

	var serialisedData = sparePartInventoryRequestInfo.Serialised()
	err = v.Repository.UpdateResource(targetTable, recordId, serialisedData, userId)

	if err != nil {
		v.Logger.Error("error check in error", zap.Error(err))
		response.SendInternalSystemError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully Cancel the Request",
	})
}
