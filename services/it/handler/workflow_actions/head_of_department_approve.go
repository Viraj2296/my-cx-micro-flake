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

func (v *ActionService) HandleDepartmentHeadApprove(ctx *gin.Context) {
	v.Logger.Info("handle head of department approve request is received")
	recordId, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canHODApprove"]) {
		if !v.EmailHandler.IsEmailTemplateExist(v.Database, const_util.HeadOfDepartmentApproveEmailTemplateType) {
			v.Logger.Error("No head of department approve email template configured, user_id", zap.Any("user_id", basicUserInfo.UserId))
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
			v.Logger.Error("No action allowed for head of department approve,consider adding category, resource_id", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring category in the service request",
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

		var listOfWorkflowUsers []int
		if workflowEngineId == const_util.SAPManagerWorkFlowEngine {
			v.Logger.Info("getting SAP workflow users for", zap.Any("user_id", basicUserInfo.UserId), zap.Any("category_id", categoryId))
			serviceRequestInfo["canEdit"] = true
			listOfWorkflowUsers = database.GetWorkFlowUsers(v.Database, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.UserHodSapManager)
		} else if workflowEngineId == const_util.ITManagerWorkFlowEngine {
			v.Logger.Info("getting IT Manager workflow users for", zap.Any("user_id", basicUserInfo.UserId), zap.Any("category_id", categoryId))
			serviceRequestInfo["canEdit"] = true
			listOfWorkflowUsers = database.GetWorkFlowUsers(v.Database, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.UserHODITManager)
		} else {
			v.Logger.Info("getting general workflow users for", zap.Any("user_id", basicUserInfo.UserId), zap.Any("category_id", categoryId))
			listOfWorkflowUsers = database.GetWorkFlowUsers(v.Database, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.UserHodExecutionReviewer)
		}

		if !(len(listOfWorkflowUsers) > 0) {
			v.Logger.Error("no workflow users configured for head of department approval, resource_id", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, const_util.ErrorGettingObjectsInformation,
				&response.DetailedError{
					Header:      "Email Ids aren't configured",
					Description: "No routing email is configured by admin, please contact admin to configure the routing emails before proceed",
				})
			return
		}

		if existingLevelCounter, ok := serviceRequestInfo["levelCounter"]; ok {
			var intLevelCounter = util.InterfaceToInt(existingLevelCounter)
			serviceRequestInfo["levelCounter"] = intLevelCounter + 1
		}
		serviceRequestInfo["canHODApprove"] = false
		serviceRequestInfo["canHODReturn"] = false
		serviceRequestInfo["canHODReject"] = false
		serviceRequestInfo["canExecPartyComplete"] = false

		serviceRequestInfo["canEdit"] = false

		// if this is SAP related one, then assign canAddCustomExecution true
		//TODO don't hard code the SAP, move this configuration to that category group can be configued based on deployment environment
		listOfSAPRequestIdInterface, err := database.GetConditionalObjects(v.Database, const_util.ITServiceRequestCategoryTable, " object_info->'$.categoryGroup' = 'SAP'")
		var listOfSAPRequestCategory []int
		for _, SAPRequestInterface := range *listOfSAPRequestIdInterface {
			listOfSAPRequestCategory = append(listOfSAPRequestCategory, SAPRequestInterface.Id)
		}
		if util.HasInt(categoryId, listOfSAPRequestCategory) {
			serviceRequestInfo["canAddCustomExecution"] = true
			v.Logger.Info("this is the SAP request, so adding canAddCustomExecution", zap.Any("request_id", recordId))
		} else {
			v.Logger.Info("this is not the SAP request, so not adding canAddCustomExecution", zap.Any("request_id", recordId))
		}

		if workflowEngineId == const_util.SAPManagerWorkFlowEngine {
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowSapManager
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingSapManager
			serviceRequestInfo["canSapApprove"] = true
			serviceRequestInfo["canSapReject"] = true
			serviceRequestInfo["canEdit"] = true
		} else if workflowEngineId == const_util.ITManagerWorkFlowEngine {
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowITManager
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingITManager
			serviceRequestInfo["canITManagerApprove"] = true
			serviceRequestInfo["canITManagerReject"] = true
			serviceRequestInfo["canEdit"] = true
		} else {
			serviceRequestInfo["isAssignable"] = true
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowExecutionParty
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingExecParty
		}

		// send the email about ack saying, you request is under review

		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "APPROVED BY HOD",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Great, your request has been approved by Head Of Department. Please wait until IT is approving your request",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(v.Database, const_util.ITServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		for _, workflowUser := range listOfWorkflowUsers {
			var err error
			if workflowEngineId == const_util.SAPManagerWorkFlowEngine {
				serviceRequestInfo["canEdit"] = true
				err = v.EmailHandler.EmailGenerator(v.Database, const_util.HeadOfDepartmentApproveToSapEmailTemplateType, workflowUser, const_util.ITServiceMyDepartmentRequestComponent, recordId)
			} else if workflowEngineId == const_util.ITManagerWorkFlowEngine {
				serviceRequestInfo["canEdit"] = true
				err = v.EmailHandler.EmailGenerator(v.Database, const_util.HeadOfDepartmentApproveToITEmailTemplateType, workflowUser, const_util.ITServiceMyDepartmentRequestComponent, recordId)
			} else {
				err = v.EmailHandler.EmailGenerator(v.Database, const_util.HeadOfDepartmentApproveEmailTemplateType, workflowUser, const_util.ITServiceMyDepartmentRequestComponent, recordId)
			}
			if err != nil {
				v.Logger.Error("error generating email during head of department approval", zap.String("error", err.Error()))
			}

		}

		v.Logger.Info("Head of department approval has been successfully completed, record_id", zap.Any("record_id", recordId))
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})

	} else {
		v.Logger.Error("Head of department approval has failed due to canHODApprove is false, record_id", zap.Any("record_id", recordId))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
