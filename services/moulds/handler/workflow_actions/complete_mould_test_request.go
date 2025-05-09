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
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

func (v *ActionService) CompleteMouldTestRequest(ctx *gin.Context) {
	recordId := util.GetRecordId(ctx)
	v.Logger.Info("processing complete mould test request", zap.Int("record_id", recordId))
	err, eventObject := database.Get(v.Database, const_util.MouldTestRequestTable, recordId)

	if err != nil {
		v.Logger.Error("requested resource not found", zap.Error(err))
		response.SendResourceNotFound(ctx)
		return
	}

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	completeRemark := util.InterfaceToString(returnFields["remark"])
	testResultStatus := util.InterfaceToString(returnFields["testStatus"])

	userId := common.GetUserId(ctx)
	mouldTestRequest := database.MouldTestRequest{ObjectInfo: eventObject.ObjectInfo}
	mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()

	mouldTestRequestInfo.MouldTestStatus = const_util.MouldTestWorkFlowQuality
	var actionStatus = "FAILED"
	if testResultStatus == "passed" {
		actionStatus = "PASSED"
		mouldTestRequestInfo.ActionStatus = const_util.MouldTestRequestActionPassed
	} else {
		actionStatus = "FAILED"
		mouldTestRequestInfo.ActionStatus = const_util.MouldTestRequestActionFailed
	}

	mouldTestRequestInfo.CanCheckOut = false
	mouldTestRequestInfo.CanContinueTest = false
	mouldTestRequestInfo.CanComplete = false
	mouldTestRequestInfo.IsUpdate = false
	mouldTestRequestInfo.Draggable = false
	mouldTestRequestInfo.CanApprove = true

	existingActionRemarks := mouldTestRequestInfo.ActionRemarks

	existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
		Status:        "TEST COMPLETED [" + actionStatus + "]",
		UserId:        userId,
		Remarks:       "The mould test request has been successfully completed [" + completeRemark + "]",
		ProcessedTime: GetTimeDifference(mouldTestRequestInfo.CreatedAt),
	})

	mouldTestRequestInfo.ActionRemarks = existingActionRemarks
	var serialisedData = mouldTestRequestInfo.DatabaseSerialize(userId)
	v.Logger.Info("updating complete mould test request", zap.Any("approve_request", string(mouldTestRequestInfo.Serialize())))
	err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, serialisedData)
	if err != nil {
		v.Logger.Error("error updating mould test request", zap.Error(err))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal system error",
				Description: "Sorry, your action couldn't able to perform now due to internal system error. Please report this error to system administrator.",
			})
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
	objectFields[const_util.MouldMasterFieldMouldStatus] = const_util.MouldStatusCustomerApproval

	objectFields[const_util.CommonFieldLastUpdatedBy] = userId
	objectFields[const_util.CommonFieldLastUpdatedAt] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	serialisedMouldMaster, _ := json.Marshal(objectFields)
	updatingMouldMasterData["object_info"] = serialisedMouldMaster

	err = database.Update(v.Database, const_util.MouldMasterTable, mouldTestRequestInfo.MouldId, updatingMouldMasterData)
	if err != nil {
		v.Logger.Error("error updating mould master data", zap.Error(err))
	}

	err = v.EmailHandler.EmailGenerator(v.Database, const_util.MouldTestRequestSubmittedTemplateType, mouldTestRequestInfo.CreatedBy, const_util.MouldTestRequestComponent, recordId)
	if err == nil {
		v.Logger.Error("email generation failed", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Test is successfully completed.",
		Code:    0,
	})
}
