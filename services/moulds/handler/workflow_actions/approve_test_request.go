package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) ApproveTestRequest(ctx *gin.Context) {
	recordId := util.GetRecordId(ctx)
	v.Logger.Info("handing approve test request ", zap.Int("recordId", recordId))
	err, mouldTestRequestGeneralObject := database.Get(v.Database, const_util.MouldTestRequestTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	userId := common.GetUserId(ctx)
	mouldTestRequest := database.MouldTestRequest{ObjectInfo: mouldTestRequestGeneralObject.ObjectInfo}
	mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()

	var actionStatus string
	var actionRemark string

	if mouldTestRequestInfo.MouldTestStatus == const_util.MouldTestWorkFlowProcessDepartment {
		mouldTestRequestInfo.MouldTestStatus = const_util.MouldTestWorkFlowProcessDepartment
		actionRemark = "Approved by process engineer with status  "
	} else if mouldTestRequestInfo.MouldTestStatus == const_util.MouldTestWorkFlowQuality {
		mouldTestRequestInfo.MouldTestStatus = const_util.MouldTestWorkFlowTooling
		actionRemark = "Approved by quality engineer with status  "

	}
	mouldTestRequestInfo.ApprovedBy = userId

	v.Logger.Info("setting mould action status", zap.Any("action_status", mouldTestRequestInfo.ActionStatus))
	existingActionRemarks := mouldTestRequestInfo.ActionRemarks
	existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
		Status:        actionStatus,
		UserId:        userId,
		Remarks:       "Mould Test Request has been approved, It is now ready go for production [" + actionRemark + "]",
		ProcessedTime: GetTimeDifference(mouldTestRequestInfo.CreatedAt),
	})
	mouldTestRequestInfo.CanApprove = false
	mouldTestRequestInfo.ActionRemarks = existingActionRemarks
	mouldTestRequestInfo.ActionStatus = const_util.MouldTestRequestActionApproved
	paramErr := machineService.UpdateMachineParamEditStatus(const_util.ProjectID, false, mouldTestRequestInfo.MachineParamId)

	var serialisedData = mouldTestRequestInfo.DatabaseSerialize(userId)
	v.Logger.Info("updating approve mould test request", zap.Any("approve_request", string(mouldTestRequestInfo.Serialize())))
	err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, serialisedData)
	if err != nil {
		v.Logger.Error("error updating mould test request table", zap.Error(err))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal system error",
				Description: "Internal system error starting test request, please summit error code to system admin",
			})
		return
	}
	if paramErr != nil {
		v.Logger.Error("requested machine param resource not found", zap.Error(err))
		response.SendResourceNotFound(ctx)
		return
	}

	// Update the mould Customer Approval Status
	err, generalObject := database.Get(v.Database, const_util.MouldMasterTable, mouldTestRequestInfo.MouldId)

	if err != nil {
		v.Logger.Error("error getting mould master table", zap.Error(err))
	}
	updatingMouldMasterData := make(map[string]interface{})

	var objectFields = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &objectFields)
	objectFields[const_util.MouldMasterFieldMouldStatus] = const_util.MouldStatusActive

	objectFields[const_util.CommonFieldLastUpdatedBy] = userId
	objectFields[const_util.CommonFieldLastUpdatedAt] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	serialisedMouldMaster, _ := json.Marshal(objectFields)
	updatingMouldMasterData["object_info"] = serialisedMouldMaster

	err = database.Update(v.Database, const_util.MouldMasterTable, mouldTestRequestInfo.MouldId, updatingMouldMasterData)
	if err != nil {
		v.Logger.Error("error updating mould master data", zap.Error(err))
	}
	err = v.EmailHandler.EmailGenerator(v.Database, const_util.MouldTestRequestApprovedTemplateType, mouldTestRequestInfo.TestedBy, const_util.MouldTestRequestComponent, recordId)
	if err != nil {
		v.Logger.Error("error sending email", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Test is successfully completed",
		Code:    0,
	})
}
