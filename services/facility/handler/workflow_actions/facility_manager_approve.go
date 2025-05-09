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

func (v *ActionService) HandleFacilityManagerApprove(ctx *gin.Context) {
	v.Logger.Info("handle head of department approve request is received")
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canFacilityManagerApprove"]) {
		if !v.EmailHandler.IsEmailTemplateExist(dbConnection, const_util.FacilityManagerApproveEmailTemplateType) {
			v.Logger.Error("No head of facility manager approve email template configured, user_id", zap.Any("user_id", basicUserInfo.UserId))
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

		err, categoryTemplateObject := database.Get(dbConnection, const_util.FacilityServiceCategoryTemplateTable, categoryTemplateIdId)
		categoryTemplateInfo := make(map[string]interface{})
		json.Unmarshal(categoryTemplateObject.ObjectInfo, &categoryTemplateInfo)
		workflowEngineId := util.InterfaceToInt(categoryTemplateInfo["hodRoutingOption"])

		isSafetyRelated := util.InterfaceToIntArray(serviceRequestInfo["isSafetyRelated"])
		var listOfWorkflowUsers []int

		if len(isSafetyRelated) > 0 {
			listOfWorkflowUsers = database.GetWorkFlowUsers(dbConnection, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.EhsManagerWorkFlow)
		} else {
			listOfWorkflowUsers = database.GetWorkFlowUsers(dbConnection, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.TechnicianWorkFlow)
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
		serviceRequestInfo["canHODReject"] = false
		serviceRequestInfo["canFacilityManagerApprove"] = false
		serviceRequestInfo["canFacilityManagerReject"] = false

		serviceRequestInfo["canEdit"] = false

		if len(isSafetyRelated) > 0 {
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowEHSManager
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingEHSManager
			serviceRequestInfo["canEHSManagerApprove"] = true
			serviceRequestInfo["canEHSManagerReject"] = true
		} else {
			serviceRequestInfo["serviceStatus"] = const_util.WorkFlowTechnician
			serviceRequestInfo["actionStatus"] = const_util.ActionPendingTechnicianParty
			serviceRequestInfo["isAssignable"] = true
		}

		// Check if the service request is not safety-related and set canAddCustomExecution to true if so
		if value, ok := serviceRequestInfo["isSafetyRelated"]; ok {
			if !util.InterfaceToBool(value) {
				serviceRequestInfo["canAddCustomExecution"] = true
			}
		}

		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		if len(isSafetyRelated) > 0 {
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "APPROVED BY FACILITY MANAGER",
				UserId:        basicUserInfo.UserId,
				Remarks:       "Great, your request has been approved by Facility Manager. Please wait until EHS Manager is approving your request",
				ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
			})
		} else {
			existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
				Status:        "APPROVED BY FACILITY MANAGER",
				UserId:        basicUserInfo.UserId,
				Remarks:       "Great, your request has been approved by Facility Manager. Please wait until Technician proceed wit the request",
				ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
			})
		}

		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		if len(isSafetyRelated) > 0 {

			for _, workflowUser := range listOfWorkflowUsers {
				err = v.EmailHandler.EmailGenerator(dbConnection, const_util.FacilityManagerApprovalForEHSManagerEmailTemplateType, workflowUser, const_util.FacilityServiceMyDepartmentRequestComponent, recordId)
				if err != nil {
					v.Logger.Error("error generating email during head of department approval", zap.String("error", err.Error()))
				}
			}
		} else {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.EmailHandler.EmailGenerator(dbConnection, const_util.FacilityManagerApprovalForTechnicianEmailTemplateType, workflowUser, const_util.FacilityServiceMyDepartmentRequestComponent, recordId)
				if err != nil {
					v.Logger.Error("error generating email during head of department approval", zap.String("error", err.Error()))
				}
			}
		}

		v.Logger.Info("Head of department approval has been successfully completed, record_id", zap.Any("record_id", recordId))
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})

	} else {
		v.Logger.Error("Head of department approval has failed due to canFacilityManagerApprove is false, record_id", zap.Any("record_id", recordId))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
