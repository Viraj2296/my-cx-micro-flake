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

func (v *LabourManagementService) recordPOSTActionHandler(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	if actionName == "update_complete_count" {
		v.WorkflowActionHandler.UpdateProductCompleteCount(ctx)
		return
	} else if actionName == "update_manpower" {
		v.WorkflowActionHandler.HandleUpdateManpower(ctx)
		return
	} else if actionName == "start_shift" {
		v.WorkflowActionHandler.HandleStartShift(ctx)
		return
	} else if actionName == "get_user_lines" {
		v.BaseService.Logger.Info("getting user lines")
		v.WorkflowActionHandler.GetUserLines(ctx)
		return
	} else if actionName == "get_shift_lines" {
		v.BaseService.Logger.Info("getting shift lines")
		v.WorkflowActionHandler.GetShiftLines(ctx)
		return
	} else if actionName == "stop_shift" {
		v.WorkflowActionHandler.HandleStopShift(ctx)
		return
	} else if actionName == "stop_shift_event" {
		v.WorkflowActionHandler.HandleStopShiftEvent(ctx)
		return
	} else if actionName == "rollback_shift" {
		v.WorkflowActionHandler.RollBackShift(ctx)
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}

func (v *LabourManagementService) handleComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "check_in" {
		v.WorkflowActionHandler.HandleCheckIn(ctx)
		return
	} else if actionName == "validate_check_in" {
		v.WorkflowActionHandler.ValidateCheckIn(ctx)
		return
	} else if actionName == "check_out" {
		v.WorkflowActionHandler.HandleCheckOut(ctx)
		return
	} else if actionName == "update_shift_events" {
		v.BaseService.Logger.Info("handling update shift events ")
		v.WorkflowActionHandler.UpdateShiftEvents(ctx)
		return
	} else if actionName == "create_shift" {
		v.WorkflowActionHandler.CreateShift(ctx)
		return
	} else if actionName == "get_shift_scheduler_events" {
		v.BaseService.Logger.Info("handling get shift scheduler events ")
		v.WorkflowActionHandler.GetShiftSchedulerEvents(ctx)
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}
func (v *LabourManagementService) handleModulePOSTAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	if actionName == "line_display" {
		v.BaseService.Logger.Info("handling line display")
		v.WorkflowActionHandler.GetLineDisplay(ctx)
	} else if actionName == "shift_production_tv_view" {
		v.BaseService.Logger.Info("handling shift_production_tv_view")
		v.WorkflowActionHandler.GetShiftProductionTvView(ctx)
	} else if actionName == "assembly_shift_summary_display" {
		v.BaseService.Logger.Info("handling assembly shift summary display")
		v.WorkflowActionHandler.GetaAssemblyShiftSummaryDisplay(ctx)
	} else if actionName == "assembly_shift_summary_display_v2" {
		v.BaseService.Logger.Info("handling assembly_shift_summary_display_v2")
		v.WorkflowActionHandler.GetaAssemblyShiftSummaryDisplayV2(ctx)
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}
func (v *LabourManagementService) recordGetActionHandler(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "get_shift_history_details" {
		v.BaseService.Logger.Info("handling get shift history details")
		v.WorkflowActionHandler.GetHistoryShiftDetails(ctx)
	} else if actionName == "get_active_shift_details" {
		v.BaseService.Logger.Info("handling get active shift  details for mobile")
		v.WorkflowActionHandler.GetActiveShiftDetails(ctx)
	} else if actionName == "get_shift_templates" {
		v.BaseService.Logger.Info("handling get shift templates")
		v.WorkflowActionHandler.GetShiftTemplates(ctx)
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}
