package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"

	"github.com/gin-gonic/gin"
)

func (v *SchedulerService) handleComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "get_events_by_resources" {
		v.WorkflowActionHandler.GetEventsByResources(ctx)
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}
