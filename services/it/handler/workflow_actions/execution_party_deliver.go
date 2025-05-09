package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

func (v *ActionService) HandleExecutionPartyDeliver(ctx *gin.Context) {

	var deliverRemarksFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&deliverRemarksFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	deliverRemark := util.InterfaceToString(deliverRemarksFields["remark"])

	recordId, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canExecPartyComplete"]) {
		if !v.EmailHandler.IsEmailTemplateExist(v.Database, const_util.ExecutionPartyDeliverTemplateType) {
			v.Logger.Error("handle execution party deliver has failed due to invalid email template type", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}

		if basicUserInfo.UserId != util.InterfaceToInt(serviceRequestInfo["assignedExecutionParty"]) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "You aren't allowed to make action to this request",
				})
			return
		}

		var actionStatus = const_util.ActionPendingAcknowledgment
		if categoryId, ok := serviceRequestInfo["categoryId"]; ok {
			err, workflowEngineObject := database.Get(v.Database, const_util.ITServiceWorkflowEngineTable, 3)
			if err == nil {
				workflowEngine := database.ITServiceWorkflowEngine{ObjectInfo: workflowEngineObject.ObjectInfo}
				var categoryIdInt = util.InterfaceToInt(categoryId)
				if util.HasInt(categoryIdInt, workflowEngine.GetWorkFlowEngineInfo().ListOfTemplates) {
					var existingActionStatus = util.InterfaceToString(serviceRequestInfo["actionStatus"])
					if existingActionStatus == const_util.ActionTested {
						actionStatus = const_util.ActionPendingAcknowledgment
					} else {
						actionStatus = const_util.ActionPendingTesting
					}
				} else {
					actionStatus = const_util.ActionPendingAcknowledgment
				}
			} else {
				actionStatus = const_util.ActionPendingAcknowledgment
			}
		} else {
			actionStatus = const_util.ActionPendingAcknowledgment
		}

		if existingLevelCounter, ok := serviceRequestInfo["levelCounter"]; ok {
			var intLevelCounter = util.InterfaceToInt(existingLevelCounter)
			serviceRequestInfo["levelCounter"] = intLevelCounter + 1
		}
		serviceRequestInfo["canHODApprove"] = false
		serviceRequestInfo["canHODReturn"] = false
		serviceRequestInfo["canHODReject"] = false

		serviceRequestInfo["canUserSubmit"] = false
		serviceRequestInfo["canCancel"] = false
		serviceRequestInfo["canUserAcknowledge"] = true
		serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
		serviceRequestInfo["actionStatus"] = actionStatus
		serviceRequestInfo["canExecPartyComplete"] = false
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		if actionStatus == const_util.ActionPendingAcknowledgment {
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "EXECUTED",
				UserId:        basicUserInfo.UserId,
				Remarks:       deliverRemark + " Great, your request has been executed successfully, please go and acknowledge or re-open if modification needed",
				ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
			})
		}
		if actionStatus == const_util.ActionPendingTesting {
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "EXECUTED",
				UserId:        basicUserInfo.UserId,
				Remarks:       deliverRemark + " Great, your request has been executed successfully, Please do the testing before proceed further",
				ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
			})
		}

		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := database.Update(v.Database, const_util.ITServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle execution party deliver has failed due to update resource", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		v.EmailHandler.EmailGenerator(v.Database, const_util.ExecutionPartyDeliverTemplateType, util.InterfaceToInt(serviceRequestInfo["createdBy"]), const_util.ITServiceMyExecutionRequestComponent, recordId)

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
