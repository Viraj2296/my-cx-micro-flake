package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (v *ActionService) HandleUserAcknowledgement(ctx *gin.Context) {

	recordId, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canUserAcknowledge"]) {
		if !v.EmailHandler.IsEmailTemplateExist(v.Database, const_util.UserAcknowledgeEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}

		categoryId := util.InterfaceToInt(serviceRequestInfo["categoryId"])
		err, categoryObject := database.Get(v.Database, const_util.ITServiceRequestCategoryTable, categoryId)

		if err != nil {
			v.Logger.Error("handle user ack has failed due to getting resource has failed", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring category in the it service request",
				})
			return
		}

		categoryInfo := make(map[string]interface{})
		json.Unmarshal(categoryObject.ObjectInfo, &categoryInfo)
		categoryTemplateIdId := util.InterfaceToInt(categoryInfo["categoryTemplate"])

		err, categoryTemplateObject := database.Get(v.Database, const_util.IITServiceCategoryTemplateTable, categoryTemplateIdId)
		categoryTemplateInfo := make(map[string]interface{})
		json.Unmarshal(categoryTemplateObject.ObjectInfo, &categoryTemplateInfo)
		workflowEngineId := util.InterfaceToInt(categoryTemplateInfo["hodRoutingOption"])

		var entryIndex int
		if workflowEngineId == const_util.HodWorkFlow {
			entryIndex = 3
		} else if workflowEngineId == const_util.ExecutionWorkFlow {
			entryIndex = 2
		} else {
			entryIndex = 3
		}

		listOfWorkflowUsers := database.GetWorkFlowUsers(v.Database, workflowEngineId, basicUserInfo.UserId, categoryId, entryIndex)

		if !(len(listOfWorkflowUsers) > 0) {
			v.Logger.Error("handle execution party deliver has failed due workflow users are empty", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Email Ids aren't configured",
					Description: "Email can't be sent to respective persons. Please configure the email ids.",
				})
			return
		}
		if existingLevelCounter, ok := serviceRequestInfo["levelCounter"]; ok {
			var intLevelCounter = util.InterfaceToInt(existingLevelCounter)
			serviceRequestInfo["levelCounter"] = intLevelCounter + 1
		}
		serviceRequestInfo["canUserSubmit"] = false
		serviceRequestInfo["canUserAcknowledge"] = false
		serviceRequestInfo["canHODReject"] = false
		serviceRequestInfo["canHODApprove"] = false
		serviceRequestInfo["canHODReturn"] = false

		err, workflowEngineObject := database.Get(v.Database, const_util.ITServiceWorkflowEngineTable, 3)
		if err == nil {
			workflowEngine := database.ITServiceWorkflowEngine{ObjectInfo: workflowEngineObject.ObjectInfo}
			var categoryIdInt = util.InterfaceToInt(categoryId)
			if util.HasInt(categoryIdInt, workflowEngine.GetWorkFlowEngineInfo().ListOfTemplates) {
				var existingActionStatus = util.InterfaceToString(serviceRequestInfo["actionStatus"])
				if existingActionStatus == const_util.ActionPendingTesting {
					serviceRequestInfo["actionStatus"] = const_util.ActionTested
					serviceRequestInfo["serviceStatus"] = const_util.WorkFlowExecutionParty
					serviceRequestInfo["canExecPartyComplete"] = true
				} else {
					serviceRequestInfo["canEdit"] = false
					serviceRequestInfo["actionStatus"] = const_util.ActionClosed
					serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
				}

			} else {
				serviceRequestInfo["actionStatus"] = const_util.ActionClosed
				serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
			}
		} else {
			serviceRequestInfo["actionStatus"] = const_util.ActionClosed
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
		}

		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		var emailTemplate int
		if util.InterfaceToString(serviceRequestInfo["actionStatus"]) == const_util.ActionClosed {
			emailTemplate = const_util.UserAcknowledgeEmailTemplateType
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "ACCEPTED BY USER",
				UserId:        basicUserInfo.UserId,
				Remarks:       "Thank you for accepting your request deliver",
				ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
			})
		} else {
			emailTemplate = const_util.UserTestedEmailTemplateType
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "TESTED",
				UserId:        basicUserInfo.UserId,
				Remarks:       "Tested by user",
				ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
			})
		}

		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(v.Database, const_util.ITServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle execution party deliver has failed due to update resource failed", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		v.Logger.Info("execution users", zap.Any("users", listOfWorkflowUsers))
		for _, workflowUser := range listOfWorkflowUsers {
			v.EmailHandler.EmailGenerator(v.Database, emailTemplate, workflowUser, const_util.ITServiceMyRequestComponent, recordId)
		}

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
