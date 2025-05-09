package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServiceRequestAction struct {
	Remark string `json:"remark"`
}

func (v *MouldService) recordPUTActionHandler(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == const_util.ActionUpdateSchedulerEvent {
		v.WorkflowActionHandler.UpdateSchedulerEvent(ctx)
		return
	} else {
		v.BaseService.Logger.Error("invalid PUT action received from client", zap.Any("action_name", actionName))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}

func (v *MouldService) recordPOSTActionHandler(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == const_util.ActionReleaseMouldTestRequest {
		v.WorkflowActionHandler.ReleaseMouldTestRequest(ctx)
		return
	} else if actionName == const_util.ActionCustomerConfirm {
		v.WorkflowActionHandler.CustomerConfirm(ctx)
		return
	} else if actionName == const_util.ActionUnReleaseMouldTestRequest {
		v.WorkflowActionHandler.UnReleaseMouldTestRequest(ctx)
		return
	} else if actionName == const_util.ActionCheckInMouldTestRequest {
		v.WorkflowActionHandler.CheckInRequest(ctx)
		return
	} else if actionName == const_util.ActionCompleteMouldTestRequest {
		v.WorkflowActionHandler.CompleteMouldTestRequest(ctx)
		return
	} else if actionName == const_util.ActionApproveMouldTestRequest {
		v.WorkflowActionHandler.ApproveTestRequest(ctx)
		return
	} else if actionName == const_util.ActionUpdateSchedulerEvent {
		v.WorkflowActionHandler.UpdateSchedulerEvent(ctx)
		return
	} else if actionName == const_util.ActionModifyMouldMasterRequest {
		v.WorkflowActionHandler.AddModificationCount(ctx)
		return
	} else {
		v.BaseService.Logger.Error("invalid POST action received from client", zap.Any("action_name", actionName))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
