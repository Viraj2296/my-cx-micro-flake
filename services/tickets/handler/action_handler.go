package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"github.com/gin-gonic/gin"
)

type ServiceRequestAction struct {
	Remark string `json:"remark"`
}

func (ts *TicketsService) recordPOSTActionHandler(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == ActionAPIHoDApprove {
		ts.handleDepartmentHeadApprove(ctx)
		return
	} else if actionName == ActionAPIHoDReturn {
		ts.handleDepartmentHeadReturn(ctx)
		return
	} else if actionName == ActionAPIHoDReject {
		ts.handleDepartmentHeadReject(ctx)
		return
	} else if actionName == ActionUserSubmit {
		ts.handleUserSubmit(ctx)
		return
	} else if actionName == ActionUserAcknowledgement {
		ts.handleUserAcknowledgement(ctx)
		return
	} else if actionName == ActionAPIITApprove {
		ts.handleReviewPartyApprove(ctx)
		return
	} else if actionName == ActionAPIITReject {
		ts.handleReviewPartyReject(ctx)
		return
	} else if actionName == ActionAPIExecutionPartyDeliver {
		ts.handleExecutionPartyDeliver(ctx)
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
