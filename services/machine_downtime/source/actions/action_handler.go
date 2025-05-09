package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *Actions) PostModuleAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	if actionName == "machine_downtime_display" {
		v.Logger.Info("handling machine_downtime_display")
		v.GetMachineDowntimeDisplay(ctx)
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}
func (v *Actions) PostComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	v.Logger.Info("handling component action", zap.String("action_name", actionName))
	if actionName == "check_in" {
		v.HandleCheckIn(ctx)
		return
	} else if actionName == "checkin_faults" {
		v.Logger.Info("handling checkin_faults")
		v.HandleCheckInFaults(ctx)
		return
	} else if actionName == "checkout_faults" {
		v.Logger.Info("handling checkout_faults")
		v.HandleCheckOutFaults(ctx)
		return
	} else if actionName == "view_notification" {
		v.HandleViewNotifications(ctx)
		return
	} else {
		v.Logger.Warn("invalid action", zap.String("action_name", actionName))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}

func (v *Actions) PostComponentResourceAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	if actionName == "check_out" {
		v.Logger.Info("handling checkout")
		v.HandleCheckout(ctx)
		return
	} else if actionName == "cancel_job" {
		v.Logger.Info("handling cancel_job")
		v.HandleCancelJob(ctx)
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

	if actionName == "fault_info" {
		v.HandleGetFaultInfo(ctx)
		v.Logger.Info("handling fault info")

	} else if actionName == "notification_history" {
		v.Logger.Info("handling notification history")
		v.HandleGetNotificationHistory(ctx)
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
	actionName := util.GetActionName(ctx)

	if actionName == "view_job_detail" {
		v.HandleJobInfoById(ctx)
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
