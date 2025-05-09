package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (v *ActionService) HandleUserAcknowledgement(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canUserAcknowledge"]) {
		if !v.EmailHandler.IsEmailTemplateExist(dbConnection, const_util.UserAcknowledgeEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}

		categoryId := util.InterfaceToInt(serviceRequestInfo["categoryId"])
		err, categoryObject := database.Get(dbConnection, const_util.FacilityServiceRequestCategoryTable, categoryId)

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

		err, categoryTemplateObject := database.Get(dbConnection, const_util.FacilityServiceCategoryTemplateTable, categoryTemplateIdId)
		categoryTemplateInfo := make(map[string]interface{})
		json.Unmarshal(categoryTemplateObject.ObjectInfo, &categoryTemplateInfo)
		//workflowEngineId := util.InterfaceToInt(categoryTemplateInfo["hodRoutingOption"])

		var assignedExecutionParty = util.InterfaceToInt(serviceRequestInfo["assignedExecutionParty"])
		//var entryIndex int
		//if workflowEngineId == const_util.HodWorkFlow {
		//	entryIndex = 3
		//} else if workflowEngineId == const_util.ExecutionWorkFlow {
		//	entryIndex = 2
		//} else {
		//	entryIndex = 3
		//}
		//
		//listOfWorkflowUsers := database.GetWorkFlowUsers(dbConnection, workflowEngineId, basicUserInfo.UserId, categoryId, entryIndex)

		if !(assignedExecutionParty > 0) {
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

		err, workflowEngineObject := database.Get(dbConnection, const_util.FacilityServiceWorkflowEngineTable, 3)
		if err == nil {
			workFlowEngineInfo := database.GetWorkFlowEngineInfo(workflowEngineObject.ObjectInfo)
			var categoryIdInt = util.InterfaceToInt(categoryId)
			if util.HasInt(categoryIdInt, workFlowEngineInfo.ListOfTemplates) {
				var existingActionStatus = util.InterfaceToString(serviceRequestInfo["actionStatus"])
				if existingActionStatus == const_util.ActionPendingTesting {
					serviceRequestInfo["actionStatus"] = const_util.ActionTested
					serviceRequestInfo["serviceStatus"] = const_util.WorkFlowTechnician
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
		}
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle execution party deliver has failed due to update resource failed", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		//v.Logger.Info("execution users", zap.Any("users", listOfWorkflowUsers))

		v.EmailHandler.EmailGenerator(dbConnection, emailTemplate, assignedExecutionParty, const_util.FacilityServiceMyRequestComponent, recordId)

		//send an additional email notification to all the security team members configured.

		//isSafetyRelated := util.InterfaceToIntArray(serviceRequestInfo["isSafetyRelated"])
		//if len(isSafetyRelated) > 0 {
		//	v.Logger.Info("request is marked as safety related, so sending email to safety team once it is acknowledged by the user")
		//	err, adminInterface := database.Get(dbConnection, const_util.FacilityServiceAdminSettingTable, 1)
		//	if err == nil {
		//		adminSettingInfo := database.GetFacilityServiceAdminSettingInfo(adminInterface.ObjectInfo)
		//		if len(adminSettingInfo.SecurityTeams) > 0 {
		//			v.Logger.Info("security team is configured, sending email now as completion note")
		//			for _, userId := range adminSettingInfo.SecurityTeams {
		//				v.EmailHandler.EmailGenerator(dbConnection, const_util.UserAcknowledgeEmailTemplateType, userId, const_util.FacilityServiceMyRequestComponent, recordId)
		//			}
		//		} else {
		//			v.Logger.Error("no security team configured to send an email")
		//		}
		//	} else {
		//		v.Logger.Error("no admin setting defined")
		//	}
		//}

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
