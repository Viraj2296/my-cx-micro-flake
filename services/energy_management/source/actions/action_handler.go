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
	if actionName == "energy_usage_display" {
		v.Logger.Info("handling energy_usage_display")
		v.GetMachineEnergyUsageDisplay(ctx)
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

	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Invalid Action",
			Description: "Your action can not be performed against this request due to sequence validation",
		})
}

func (v *Actions) PostComponentResourceAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	v.Logger.Info("handling component action", zap.String("action_name", actionName))

	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Invalid Action",
			Description: "Your action can not be performed against this request due to sequence validation",
		})
}

func (v *Actions) GetModuleAction(ctx *gin.Context) {
}
func (v *Actions) GetComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	v.Logger.Info("handling component action", zap.String("action_name", actionName))
	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Invalid Action",
			Description: "Your action can not be performed against this request due to sequence validation",
		})
}

func (v *Actions) GetComponentResourceAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	v.Logger.Info("handling component action", zap.String("action_name", actionName))
	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Invalid Action",
			Description: "Your action can not be performed against this request due to sequence validation",
		})
}
