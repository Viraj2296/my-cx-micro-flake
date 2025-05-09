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

func (v *ActionService) HandleUserSubmit(ctx *gin.Context) {
	v.Logger.Info("handle user submit is received")
	recordId, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canUserSubmit"]) {

		if !isDoHConfigured(basicUserInfo.UserId) {
			v.Logger.Error("handle user submit failed due to no head of department assigned", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Head Of Department Assigned",
					Description: "Sorry, There is no head of department configured to process your request, please ask system admin to assign you in to corresponding department",
				})
			return
		}
		if !v.EmailHandler.IsEmailTemplateExist(v.Database, const_util.UserSubmitRequestEmailTemplateType) {
			v.Logger.Error("handle user submit failed due to submit request email template not configured", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Sorry, no email configuration is done yet, please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		if existingLevelCounter, ok := serviceRequestInfo["levelCounter"]; ok {
			var intLevelCounter = util.InterfaceToInt(existingLevelCounter)
			serviceRequestInfo["levelCounter"] = intLevelCounter + 1
		}
		serviceRequestInfo["canEdit"] = false
		serviceRequestInfo["canUserSubmit"] = false

		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "USER SUBMITTED",
			UserId:        basicUserInfo.UserId,
			Remarks:       "User has successfully submitted the request. Please wait until it is processed by other layers",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks

		templateId := util.InterfaceToInt(serviceRequestInfo["templateFields"])
		_, serviceRequestTemplateGeneralObject := database.Get(v.Database, const_util.IITServiceCategoryTemplateTable, templateId)
		categoryTemplateInfo := make(map[string]interface{})
		json.Unmarshal(serviceRequestTemplateGeneralObject.ObjectInfo, &categoryTemplateInfo)

		hodRoutingOption := util.InterfaceToInt(categoryTemplateInfo["hodRoutingOption"])

		if util.InterfaceToInt(categoryTemplateInfo["hodRoutingOption"]) == 0 {
			v.Logger.Error("handle user submit failed due to hodRoutingOption not configured", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Please configure the email routing option for HOD in template",
				})
			return
		}

		categoryId := util.InterfaceToInt(serviceRequestInfo["categoryId"])
		var isHodRouting bool
		if hodRoutingOption == const_util.HODRoutingWorkflowEngine {
			v.Logger.Info("user submit, HOD routing option", zap.Any("routing_option", hodRoutingOption), zap.Any("category_id", categoryId))
			listOfWorkflowUsers := database.GetWorkFlowUsers(v.Database, hodRoutingOption, basicUserInfo.UserId, categoryId, const_util.UserExecutionReviewer)

			if len(listOfWorkflowUsers) == 0 {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "No Actions Allowed",
						Description: "Please consider configuring users for notification",
					})
				return
			}
			for _, workflowUser := range listOfWorkflowUsers {
				v.Logger.Info("generating email during user submit", zap.Any("user_id", workflowUser))
				v.EmailHandler.EmailGenerator(v.Database, const_util.UserSubmitExecutionRequestEmailTemplateType, workflowUser, const_util.ITServiceMyRequestComponent, recordId)
			}
			serviceRequestInfo["canExecPartyComplete"] = false
			serviceRequestInfo["canAddCustomExecution"] = true
			serviceRequestInfo["isAssignable"] = true
			serviceRequestInfo["canAddCustomExecution"] = true
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowExecutionParty
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingExecParty

		} else {
			if util.InterfaceToString(serviceRequestInfo["hodEmail"]) != "" {
				v.Logger.Info("generating email during user submit", zap.Any("hod_email", util.InterfaceToString(serviceRequestInfo["hodEmail"])))
				v.EmailHandler.EmailGeneratorByEmail(v.Database, const_util.UserSubmitRequestEmailTemplateType, util.InterfaceToString(serviceRequestInfo["hodEmail"]), const_util.ITServiceMyRequestComponent, recordId)
			} else {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception [Invalid HOD Email]",
						Description: "Please configure the Head Of Department email address in the request",
					})
				return
			}
			serviceRequestInfo["canHODReturn"] = true
			serviceRequestInfo["canHODReject"] = true
			serviceRequestInfo["canHODApprove"] = true
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowHOD
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingHoD
			isHodRouting = true
		}

		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := database.Update(v.Database, const_util.ITServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		// only do it for hod routing now
		if isHodRouting {
			// get the category name
			err, categoryInterface := database.Get(v.Database, const_util.IITServiceCategoryTemplateTable, categoryId)
			if err == nil {
				// if the error is nill only create the system notification, otherwise log
				var categoryFields = make(map[string]interface{})
				json.Unmarshal(categoryInterface.ObjectInfo, &categoryFields)
				if categoryInterface, ok := categoryFields["name"]; ok {
					if serviceRequestInterface, ok := serviceRequestInfo["serviceRequestId"]; ok {
						var categoryName = util.InterfaceToString(categoryInterface)
						var serviceRequestId = util.InterfaceToString(serviceRequestInterface)
						var hodEmail = util.InterfaceToString(serviceRequestInfo["hodEmail"])
						v.EmailHandler.CreateMyDepartmentSystemNotification(const_util.ProjectID, hodEmail, categoryName, serviceRequestId, recordId)
					}
				}
			}
		}
		v.Logger.Info("handle user submit is successfully processed", zap.Any("record_id", recordId))
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
