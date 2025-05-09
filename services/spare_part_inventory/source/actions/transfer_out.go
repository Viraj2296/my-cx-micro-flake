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

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"
)

func (v *Actions) TransferOut(ctx *gin.Context) {
	v.Logger.Info("transfer out request received")
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
	transferPayload := dto.TransferRequest{}
	if err := ctx.ShouldBindBodyWith(&transferPayload, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}

	err, sparePartInventoryMasterInfo := models.GetSparePartInventoryMasterInfo(objectInterface.ObjectInfo)

	if err != nil {
		v.Logger.Error("cannot unmarshal", logging.String("error", err.Error()))
		error_util.SendUnmarshlingFailed(ctx)
		return
	}

	if sparePartInventoryMasterInfo.OnHandQty < 1 || sparePartInventoryMasterInfo.OnHandQty < transferPayload.Quantity {
		v.Logger.Error("insufficient quantity at source")
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Insufficient quantity at source",
				Description: "The spare part has insufficient quantity at the source location.",
			})
		return

	}

	sparePartInventoryMasterInfo.OnHandQty -= transferPayload.Quantity

	var serialisedData = sparePartInventoryMasterInfo.Serialised()
	err = v.Repository.UpdateResource(targetTable, recordId, serialisedData, userId)

	if err != nil {
		v.Logger.Error("error check in error", zap.Error(err))
		response.SendInternalSystemError(ctx)
		return
	}

	if err := v.logTransaction(userId, *sparePartInventoryMasterInfo, transferPayload.Quantity, transferPayload.DestinationLocation, transferPayload.SourceLocation, transferPayload.ServiceNotification); err != nil {
		v.Logger.Error("error updating destination record", zap.String("error", err.Error()))
		error_util.SendResourceUpdateFailed(ctx)
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully completed the transfer out action",
	})

}

func (v *Actions) logTransaction(userId int, objectInfoMap models.SparePartInventoryMasterInfo, qty int, locationId int, sourceLocation int, serviceNotificationNumber string) error {

	transactionIn := models.SparePartInventoryTransactionInfo{
		PartNumber:                objectInfoMap.SparePartNumber,
		PartDescription:           objectInfoMap.SparePartDescription,
		DestinationLocation:       locationId,
		Transaction:               consts.TransactionStatusIn,
		Qty:                       qty,
		SourceLocation:            sourceLocation,
		ServiceNotificationNumber: serviceNotificationNumber,
	}
	transactionOut := models.SparePartInventoryTransactionInfo{
		PartNumber:                objectInfoMap.SparePartNumber,
		PartDescription:           objectInfoMap.SparePartDescription,
		DestinationLocation:       locationId,
		Transaction:               consts.TransactionStatusOut,
		Qty:                       qty,
		SourceLocation:            sourceLocation,
		ServiceNotificationNumber: serviceNotificationNumber,
	}

	var transactionsList []models.SparePartInventoryTransactionInfo
	transactionsList = append(transactionsList, transactionIn, transactionOut)
	var targetTable = v.ComponentManager.GetTargetTable(consts.SparePartInventorTransactionComponent)
	for _, transaction := range transactionsList {
		var serialisedData = transaction.Serialised()
		err, resourceId := v.Repository.CreateResource(targetTable, serialisedData, userId)
		if err == nil {
			v.Logger.Info("new spare part transaction is created", zap.Any("record_id", resourceId))
		} else {
			v.Logger.Error("error creating spare part transaction record", zap.Any("error", err.Error()))
		}
	}
	return nil
}
