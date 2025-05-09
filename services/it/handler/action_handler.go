package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"

	"github.com/gin-gonic/gin"
)

type ServiceRequestAction struct {
	Remark string `json:"remark"`
}

func (v *ITService) recordPOSTActionHandler(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == const_util.ActionAPIHoDApprove {
		v.WorkfowActionHandler.HandleDepartmentHeadApprove(ctx)
		return
	} else if actionName == const_util.ActionAPIHoDReturn {
		v.WorkfowActionHandler.HandleDepartmentHeadReturn(ctx)
		return
	} else if actionName == const_util.ActionAPIHoDReject {
		v.WorkfowActionHandler.HandleDepartmentHeadReject(ctx)
		return
	} else if actionName == const_util.ActionUserSubmit {
		v.WorkfowActionHandler.HandleUserSubmit(ctx)
		return
	} else if actionName == const_util.ActionUserAcknowledgement {
		v.WorkfowActionHandler.HandleUserAcknowledgement(ctx)
		return
	} else if actionName == const_util.ActionAPIAssignMyself {
		v.WorkfowActionHandler.HandleExecutionAssignedMyself(ctx)
		return
	} else if actionName == const_util.ActionAPIExecutionPartyDeliver {
		v.WorkfowActionHandler.HandleExecutionPartyDeliver(ctx)
		return
	} else if actionName == const_util.ActionAPISapApprove {
		v.WorkfowActionHandler.HandleSapApprove(ctx)
		return
	} else if actionName == const_util.ActionAPISapReject {
		v.WorkfowActionHandler.HandleSapReject(ctx)
		return
	} else if actionName == const_util.ActionAPIITManagerApprove {
		v.WorkfowActionHandler.HandleITManagerApprove(ctx)
		return
	} else if actionName == const_util.ActionAPIITManagerReject {
		v.WorkfowActionHandler.HandleITManagerReject(ctx)
		return
	} else if actionName == const_util.ActionAPIRollBack {
		v.handleRollBack(ctx)
		return
	} else if actionName == const_util.ActionAPICancel {
		v.WorkfowActionHandler.HandleCancelRequest(ctx)
		return
	} else if actionName == const_util.ActionReassignExecution {
		v.WorkfowActionHandler.HandleReassignExecution(ctx)
		return
	} else if actionName == const_util.ActionTransferExecution {
		v.WorkfowActionHandler.HandleTransferAssignedExecutionParty(ctx)
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
