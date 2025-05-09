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

func (v *ActionService) HandleUserSubmit(ctx *gin.Context) {
	v.Logger.Info("handle user submit is received")
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
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
		if !v.EmailHandler.IsEmailTemplateExist(dbConnection, const_util.UserSubmitRequestEmailTemplateType) {
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
		serviceRequestInfo["canAddCustomExecution"] = true

		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "USER SUBMITTED",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Request successfully submitted. Please wait while it is being processed by subsequent layers.",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks

		templateId := util.InterfaceToInt(serviceRequestInfo["templateFields"])
		_, serviceRequestTemplateGeneralObject := database.Get(dbConnection, const_util.FacilityServiceCategoryTemplateTable, templateId)
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
		if hodRoutingOption == const_util.UserTechnicianUserWorkflowEngine {
			v.Logger.Info("user submit, HOD routing option", zap.Any("routing_option", hodRoutingOption), zap.Any("category_id", categoryId))
			listOfWorkflowUsers := database.GetWorkFlowUsers(dbConnection, hodRoutingOption, basicUserInfo.UserId, categoryId, const_util.TechniciansFromWorkflow1)

			if len(listOfWorkflowUsers) == 0 {
				v.Logger.Error("no workflow users configured to route the request")
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "No Actions Allowed",
						Description: "Please consider configuring users for notification",
					})
				return
			}
			for _, workflowUser := range listOfWorkflowUsers {
				v.Logger.Info("generating email during user submit", zap.Any("user_id", workflowUser))
				v.EmailHandler.EmailGenerator(dbConnection, const_util.UserSubmitRequestRoutingOneEmailTemplateType, workflowUser, const_util.FacilityServiceMyRequestComponent, recordId)
			}
			serviceRequestInfo["canExecPartyComplete"] = false
			serviceRequestInfo["isAssignable"] = true
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowTechnician
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingTechnicianParty

		} else {
			if util.InterfaceToString(serviceRequestInfo["hodEmail"]) != "" {
				v.Logger.Info("generating email during user submit", zap.Any("hod_email", util.InterfaceToString(serviceRequestInfo["hodEmail"])))
				v.EmailHandler.EmailGeneratorByEmail(dbConnection, const_util.UserSubmitRequestEmailTemplateType, util.InterfaceToString(serviceRequestInfo["hodEmail"]), const_util.FacilityServiceMyRequestComponent, recordId)
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
		err := database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
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
			err, categoryInterface := database.Get(dbConnection, const_util.FacilityServiceCategoryTemplateTable, categoryId)
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

		// Now notifiy the technical team if the request is marked as security related.
		isSafetyRelated := util.InterfaceToIntArray(serviceRequestInfo["isSafetyRelated"])
		if len(isSafetyRelated) > 0 {
			v.Logger.Info("request is marked as safety related, so sending email to safety team")
			err, adminInterface := database.Get(dbConnection, const_util.FacilityServiceAdminSettingTable, 1)
			if err == nil {
				adminSettingInfo := database.GetFacilityServiceAdminSettingInfo(adminInterface.ObjectInfo)
				if len(adminSettingInfo.SecurityTeams) > 0 {
					v.Logger.Info("security team is configured, sending email now")
					for _, userId := range adminSettingInfo.SecurityTeams {
						v.EmailHandler.EmailGenerator(dbConnection, const_util.UsedSubmissionAdminSettingEmailTemplateType, userId, const_util.FacilityServiceMyRequestComponent, recordId)
					}
				} else {
					v.Logger.Error("no security team configued to send an email")
				}
			} else {
				v.Logger.Error("no admin setting defined")
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
