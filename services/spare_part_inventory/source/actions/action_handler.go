package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"

	"github.com/gin-gonic/gin"
)

func (v *Actions) PostModuleAction(ctx *gin.Context) {
}
func (v *Actions) PostComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "receive_in" {
		v.InventoryReceiveIn(ctx)
		return
	} else if actionName == "get_part_info" {
		v.GetPartInventoryInfo(ctx)
		v.Logger.Info("handling fault info")
		return
	} else if actionName == "get_spare_part_request" {
		v.GetSparePartRequest(ctx)
		v.Logger.Info("handling get spare part request")
		return
	} else if actionName == "create_spare_request" {
		v.CreateSparePartRequest(ctx)
		v.Logger.Info("handling create spare request")
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}

func (v *Actions) PostComponentResourceAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	if actionName == "transfer_out" {
		v.TransferOut(ctx)
		return
	} else if actionName == "cancel" {
		v.CancelAction(ctx)
		v.Logger.Info("handling fault info")
		return
	} else if actionName == "approve" {
		v.ApproveAction(ctx)
		v.Logger.Info("handling fault info")
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}

func (v *Actions) GetModuleAction(ctx *gin.Context) {
}
func (v *Actions) GetComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "requested_parts" {
		v.RequestedSpareParts(ctx)
		v.Logger.Info("handling requested spare parts")
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}

func (v *Actions) GetComponentResourceAction(ctx *gin.Context) {

}
