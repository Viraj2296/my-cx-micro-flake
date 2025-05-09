package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}
func (ts *TicketsService) getBasicInfo(ctx *gin.Context) (int, *gorm.DB, common.UserBasicInfo, *TicketsServiceRequestInfo) {
	projectId := util.GetProjectId(ctx)

	recordId := util.GetRecordId(ctx)
	userId := common.GetUserId(ctx)
	dbConnection := ts.BaseService.ServiceDatabases[projectId]
	_, serviceRequestGeneralObject := Get(dbConnection, TicketsServiceRequestTable, recordId)
	serviceRequestInfo := TicketsServiceRequestInfo{}
	json.Unmarshal(serviceRequestGeneralObject.ObjectInfo, &serviceRequestInfo)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)

	return recordId, dbConnection, basicUserInfo, &serviceRequestInfo
}

func (ts *TicketsService) getStatusName(dbConnection *gorm.DB, serviceStatusId int) string {
	_, requestStatus := Get(dbConnection, TicketsServiceRequestStatusTable, serviceStatusId)
	var requestStatusInfo TicketsServiceRequestStatusInfo
	json.Unmarshal(requestStatus.ObjectInfo, &requestStatusInfo)
	return requestStatusInfo.Status
}

func isDoHConfigured(userId int) bool {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	listOfHeads := authService.GetHeadOfDepartments(userId)
	if len(listOfHeads) == 0 {
		return false
	}
	return true
}

func (ts *TicketsService) handleUserSubmit(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanUserSubmit {

		if !isDoHConfigured(basicUserInfo.UserId) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Head Of Department Assigned",
					Description: "Sorry, There is no head of department configured to process your request, please ask system admin to assign you in to corresponding department",
				})
			return
		}
		if !ts.isEmailTemplateExist(dbConnection, UserSubmitRequestEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Sorry, no email configuration is done yet, please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanEdit = false
		serviceRequestInfo.CanUserSubmit = false
		serviceRequestInfo.CanHODReturn = true
		serviceRequestInfo.CanHODReject = true
		serviceRequestInfo.CanHODApprove = true
		serviceRequestInfo.ServiceStatus = WorkFlowHOD
		serviceRequestInfo.ActionStatus = ActionPendingHoD
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "USER SUBMITTED",
			UserId:        basicUserInfo.UserId,
			Remarks:       "User has successfully submitted the request. Please wait until it is processed by other layers",
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userId := common.GetUserId(ctx)
		headOfDepartments := authService.GetHeadOfDepartments(userId)
		for _, hodUserId := range headOfDepartments {
			ts.emailGenerator(dbConnection, UserSubmitRequestEmailTemplateType, hodUserId.UserId, TicketsServiceMyRequestComponent, recordId)
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

func (ts *TicketsService) handleDepartmentHeadApprove(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanHODApprove {
		// dont' allow if ack email not configured or route email configured
		if !ts.isEmailTemplateExist(dbConnection, HeadOfDepartmentApproveEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanHODReject = false
		serviceRequestInfo.CanITApprove = true
		serviceRequestInfo.CanEdit = false
		serviceRequestInfo.CanITReject = true
		serviceRequestInfo.ServiceStatus = WorkFlowReviewParty
		serviceRequestInfo.ActionStatus = ActionPendingITReview
		// send the email about ack saying, you request is under review

		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "APPROVED BY HOD",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Great, your request has been approved by Head Of Department. Please wait until IT is approving your request",
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		listOfWorkflowUsers := ts.getTargetWorkflowGroup(dbConnection, recordId, WorkFlowReviewParty)
		for _, workflowUser := range listOfWorkflowUsers {
			ts.emailGenerator(dbConnection, HeadOfDepartmentApproveEmailTemplateType, workflowUser, TicketsServiceMyDepartmentRequestComponent, recordId)
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

func (ts *TicketsService) handleDepartmentHeadReturn(ctx *gin.Context) {

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	returnRemark := util.InterfaceToString(returnFields["remark"])
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanHODReturn {
		if !ts.isEmailTemplateExist(dbConnection, HeadOfDepartmentReturnEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanEdit = true
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanUserSubmit = true
		serviceRequestInfo.ServiceStatus = WorkFlowUser
		serviceRequestInfo.ActionStatus = ActionReturnedByHoD
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "RETURNED BY HOD",
			UserId:        basicUserInfo.UserId,
			Remarks:       returnRemark,
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		ts.emailGenerator(dbConnection, HeadOfDepartmentReturnEmailTemplateType, serviceRequestInfo.CreatedBy, TicketsServiceMyDepartmentRequestComponent, recordId)
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

func (ts *TicketsService) handleDepartmentHeadReject(ctx *gin.Context) {

	var rejectionRemarkFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&rejectionRemarkFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	rejectionRemarks := util.InterfaceToString(rejectionRemarkFields["remark"])
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanHODReject {
		if !ts.isEmailTemplateExist(dbConnection, HeadOfDepartmentRejectEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanEdit = false
		serviceRequestInfo.CanHODReject = false
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.ServiceStatus = WorkFlowUser
		serviceRequestInfo.ActionStatus = ActionRejectedByHoD
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "REJECTED BY HOD",
			UserId:        basicUserInfo.UserId,
			Remarks:       rejectionRemarks,
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		ts.emailGenerator(dbConnection, HeadOfDepartmentRejectEmailTemplateType, serviceRequestInfo.CreatedBy, TicketsServiceMyDepartmentRequestComponent, recordId)

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

func (ts *TicketsService) handleReviewPartyApprove(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanITApprove {
		if !ts.isEmailTemplateExist(dbConnection, ITApproveEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanExecPartyComplete = true

		serviceRequestInfo.CanITApprove = false
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanHODReject = false
		serviceRequestInfo.CanUserSubmit = false
		serviceRequestInfo.CanUserAcknowledge = false
		serviceRequestInfo.CanITReject = false
		serviceRequestInfo.ServiceStatus = WorkFlowExecutionParty
		serviceRequestInfo.ActionStatus = ActionPendingExecParty
		// send the email about ack saying, you request is under review

		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "APPROVED BY IT",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Great, your request has been approved by IT, Wait until it is get executed",
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		listOfWorkflowUsers := ts.getTargetWorkflowGroup(dbConnection, recordId, WorkFlowExecutionParty)
		for _, workflowUser := range listOfWorkflowUsers {
			ts.emailGenerator(dbConnection, ITApproveEmailTemplateType, workflowUser, TicketsServiceMyReviewRequestComponent, recordId)
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

func (ts *TicketsService) handleReviewPartyReject(ctx *gin.Context) {

	var rejectionRemarkFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&rejectionRemarkFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	rejectionRemarks := util.InterfaceToString(rejectionRemarkFields["remark"])
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanITReject {
		if !ts.isEmailTemplateExist(dbConnection, ITRejectEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanHODReject = false
		serviceRequestInfo.CanUserSubmit = false
		serviceRequestInfo.CanUserAcknowledge = false
		serviceRequestInfo.CanITApprove = false
		serviceRequestInfo.CanITReject = false
		serviceRequestInfo.ServiceStatus = WorkFlowUser
		serviceRequestInfo.ActionStatus = ActionRejectedByITReview
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "REJECTED BY IT",
			UserId:        basicUserInfo.UserId,
			Remarks:       rejectionRemarks,
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		ts.emailGenerator(dbConnection, ITRejectEmailTemplateType, serviceRequestInfo.CreatedBy, TicketsServiceMyReviewRequestComponent, recordId)

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

func (ts *TicketsService) handleExecutionPartyDeliver(ctx *gin.Context) {

	var deliverRemarksFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&deliverRemarksFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	deliverRemark := util.InterfaceToString(deliverRemarksFields["remark"])

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanExecPartyComplete {
		if !ts.isEmailTemplateExist(dbConnection, ExecutionPartyDeliverTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanHODReject = false

		serviceRequestInfo.CanUserSubmit = false
		serviceRequestInfo.CanUserAcknowledge = true
		serviceRequestInfo.CanITApprove = false
		serviceRequestInfo.CanITReject = false
		serviceRequestInfo.ServiceStatus = WorkFlowUser
		serviceRequestInfo.ActionStatus = ActionPendingAcknowledgment
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "EXECUTED",
			UserId:        basicUserInfo.UserId,
			Remarks:       deliverRemark + " Great, your request has been executed successfully, please go and acknowledge or re-open if modification needed",
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		ts.emailGenerator(dbConnection, ExecutionPartyDeliverTemplateType, serviceRequestInfo.CreatedBy, TicketsServiceMyExecutionRequestComponent, recordId)

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

func (ts *TicketsService) handleUserAcknowledgement(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := ts.getBasicInfo(ctx)
	if serviceRequestInfo.CanUserAcknowledge {
		if !ts.isEmailTemplateExist(dbConnection, UserAcknowledgeEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanUserSubmit = false
		serviceRequestInfo.CanUserAcknowledge = false
		serviceRequestInfo.CanHODReject = false
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanITApprove = false
		serviceRequestInfo.CanITReject = false
		serviceRequestInfo.ServiceStatus = WorkFlowUser
		serviceRequestInfo.ActionStatus = ActionClosed
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "ACCEPTED BY USER",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Thank you for accepting your request deliver",
			ProcessedTime: getTimeDifference(serviceRequestInfo.CreatedAt),
		})
		serviceRequestInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, TicketsServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		listOfExecutionUsers := ts.getTargetWorkflowGroup(dbConnection, recordId, WorkFlowExecutionParty)
		ts.BaseService.Logger.Info("execution users", zap.Any("users", listOfExecutionUsers))
		for _, executionUserId := range listOfExecutionUsers {
			ts.emailGenerator(dbConnection, UserAcknowledgeEmailTemplateType, executionUserId, TicketsServiceRequestComponent, recordId)
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
