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

func (v *ActionService) HandleEHSManagerApprove(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canEHSManagerApprove"]) {
		// dont' allow if ack email not configured or route email configured
		if !v.EmailHandler.IsEmailTemplateExist(dbConnection, const_util.EHSManagerApproveEmailTemplateType) {
			v.Logger.Error("handle EHS manager approve has failed due to invalid email template type", zap.Any("record_id", recordId))
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
			v.Logger.Error("handle EHS manager approve has failed due to getting category has failed", zap.Any("category_id", categoryId))
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
		workflowEngineId := util.InterfaceToInt(categoryTemplateInfo["hodRoutingOption"])

		listOfWorkflowUsers := database.GetWorkFlowUsers(dbConnection, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.TechnicianWorkFlow)

		if !(len(listOfWorkflowUsers) > 0) {
			v.Logger.Error("handle EHS manager approve has failed due to empty workflow users", zap.Any("record_id", recordId))
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
		serviceRequestInfo["canEHSManagerApprove"] = false
		serviceRequestInfo["canEHSManagerReject"] = false
		serviceRequestInfo["canExecPartyComplete"] = false
		serviceRequestInfo["isAssignable"] = true
		serviceRequestInfo["canEdit"] = true
		serviceRequestInfo["serviceStatus"] = const_util.WorkFlowTechnician
		serviceRequestInfo["actionStatus"] = const_util.ActionPendingTechnicianParty
		// send the email about ack saying, you request is under review

		// Check if the service request is safety-related and set canAddCustomExecution to true if so
		if value, ok := serviceRequestInfo["isSafetyRelated"]; ok {
			if util.InterfaceToBool(value) {
				serviceRequestInfo["canAddCustomExecution"] = true
			}
		}

		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "APPROVED BY EHS MANAGER",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Great news! Your request has been approved by the EHS Manager. We are now processing it and will notify you once complete.",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle EHS manager approve has failed due to update resource failed", zap.String("error", err.Error()))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		for _, workflowUser := range listOfWorkflowUsers {
			v.Logger.Info("sending email to workflow users", zap.Any("workflow_user", workflowUser))
			err := v.EmailHandler.EmailGenerator(dbConnection, const_util.EHSManagerApproveEmailTemplateType, workflowUser, const_util.FacilityServiceMyEHSManagerRequestComponent, recordId)
			if err != nil {
				v.Logger.Info("function, canEHSManagerApprove, error generating email", zap.String("error", err.Error()))
			}

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
