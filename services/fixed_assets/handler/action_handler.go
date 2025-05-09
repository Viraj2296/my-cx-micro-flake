package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"

	"github.com/gin-gonic/gin"
)

func (v *FixedAssetService) recordPOSTActionHandler(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == "submit" {
		v.handleSubmit(ctx)
		return

	} else if actionName == "hod_approve" {
		v.handleHODApprove(ctx)
		return

	} else if actionName == "hod_reject" {
		v.handleHODReject(ctx)
		return

	} else if actionName == "ceo_approve" {
		v.handleCEOApprove(ctx)
		return

	} else if actionName == "ceo_reject" {
		v.handleCEOReject(ctx)
		return

	} else if actionName == "notify_vendor" {
		v.handleNotifyVendor(ctx)
		return

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
