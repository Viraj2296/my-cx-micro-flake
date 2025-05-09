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

func (v *ActionService) HandleITManagerApprove(ctx *gin.Context) {

	recordId, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canITManagerApprove"]) {
		// dont' allow if ack email not configured or route email configured
		if !v.EmailHandler.IsEmailTemplateExist(v.Database, const_util.SapApproveEmailTemplateType) {
			v.Logger.Error("handle IT manager approve has failed due to invalid email template type", zap.Any("record_id", recordId))
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
			v.Logger.Error("handle IT manager approve has failed due to getting category has failed", zap.Any("category_id", categoryId))
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

		listOfWorkflowUsers := database.GetWorkFlowUsers(v.Database, workflowEngineId, basicUserInfo.UserId, categoryId, const_util.UserHodExecutionHod)

		if !(len(listOfWorkflowUsers) > 0) {
			v.Logger.Error("handle IT manager approve has failed due to empty workflow users", zap.Any("record_id", recordId))
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
		serviceRequestInfo["canITManagerApprove"] = false
		serviceRequestInfo["canExecPartyComplete"] = false
		serviceRequestInfo["isAssignable"] = true
		serviceRequestInfo["canEdit"] = true
		serviceRequestInfo["serviceStatus"] = const_util.WorkFlowExecutionParty
		serviceRequestInfo["actionStatus"] = const_util.ActionPendingExecParty
		// send the email about ack saying, you request is under review

		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "APPROVED BY IT MANAGER",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Great news! Your request has been approved by the IT Manager. We are now processing it and will notify you once complete.",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(v.Database, const_util.ITServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle IT manager approve has failed due to update resource failed", zap.String("error", err.Error()))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		for _, workflowUser := range listOfWorkflowUsers {
			v.Logger.Info("function, handleITManagerApprove", zap.Any("workflowuser", workflowUser))
			err := v.EmailHandler.EmailGenerator(v.Database, const_util.ITManagerApprovalEmailTemplateType, workflowUser, const_util.ITServiceMyITManagementRequestComponent, recordId)
			if err != nil {
				v.Logger.Info("function, handleITManagerApprove, error generating email", zap.String("error", err.Error()))
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
