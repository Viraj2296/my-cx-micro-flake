package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
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
func (v *WorkHourAssignmentService) getBasicInfo(ctx *gin.Context) (int, *gorm.DB, common.UserBasicInfo, *WorkHourAssignmentRequestInfo) {
	projectId := util.GetProjectId(ctx)

	recordId := util.GetRecordId(ctx)
	userId := common.GetUserId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, serviceRequestGeneralObject := Get(dbConnection, WorkHourAssignmentRequestTable, recordId)
	serviceRequestInfo := WorkHourAssignmentRequestInfo{}
	json.Unmarshal(serviceRequestGeneralObject.ObjectInfo, &serviceRequestInfo)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)

	return recordId, dbConnection, basicUserInfo, &serviceRequestInfo
}

func (v *WorkHourAssignmentService) getStatusName(dbConnection *gorm.DB, serviceStatusId int) string {
	_, requestStatus := Get(dbConnection, WorkHourAssignmentRequestStatusTable, serviceStatusId)
	var requestStatusInfo WorkHourAssignmentRequestStatusInfo
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

func isDoSConfigured(userId int) bool {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	listOfHeads := authService.GetHeadOfSections(userId)
	if len(listOfHeads) == 0 {
		return false
	}
	return true
}

func (v *WorkHourAssignmentService) handleUserSubmit(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if serviceRequestInfo.CanUserSubmit {
		var listOfHeads []common.UserBasicInfo
		// first check is there any section defined for that department which that user belongs,
		factoryService := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
		departmentList := basicUserInfo.Department

		var sectionList []int

		for _, departmentId := range departmentList {

			section, _ := factoryService.GetSections(departmentId)

			for _, sectionOption := range section {
				sectionFound := false
				for _, existingSection := range sectionList {
					if sectionOption == existingSection {
						sectionFound = true
					}
				}

				if !sectionFound {
					sectionList = append(sectionList, sectionOption)
				}
			}
		}

		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

		if len(sectionList) > 0 {
			// if sections defined (else case), then we need to check any section head is defined under this department section which that user belongs
			if !isDoSConfigured(basicUserInfo.UserId) {
				// if the section is not defined , then throw error, "No section head is assigned"
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "No Head Of Section Assigned",
						Description: "Sorry, There is no head of section configured to process your request, please ask system admin to assign you in to corresponding department",
					})
				return
			}
			listOfHeads = authService.GetHeadOfSections(basicUserInfo.UserId)
		} else {
			// if not defined, then route to head of department , the following checkc is needed before route
			if !isDoHConfigured(basicUserInfo.UserId) {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "No Head Of Department Assigned",
						Description: "Sorry, There is no head of department configured to process your request, please ask system admin to assign you in to corresponding department",
					})
				return
			}
			listOfHeads = authService.GetHeadOfDepartments(basicUserInfo.UserId)
		}

		if !v.isEmailTemplateExist(dbConnection, UserSubmitRequestEmailTemplateType) {
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
		serviceRequestInfo.AssignmentStatus = WorkFlowHOD
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
		err := Update(dbConnection, WorkHourAssignmentRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		// authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		// userId := common.GetUserId(ctx)
		// headOfDepartments := authService.GetHeadOfDepartments(userId)
		v.BaseService.Logger.Info("selected hod", zap.Any("head", listOfHeads))
		for _, hodUserId := range listOfHeads {
			v.emailGenerator(dbConnection, UserSubmitRequestEmailTemplateType, hodUserId.UserId, WorkHourAssignmentMyRequestComponent, recordId)
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

func (v *WorkHourAssignmentService) handleDepartmentHeadApprove(ctx *gin.Context) {

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if serviceRequestInfo.CanHODApprove {
		// dont' allow if ack email not configured or route email configured
		if !v.isEmailTemplateExist(dbConnection, HeadOfDepartmentApproveEmailTemplateType) {
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
		serviceRequestInfo.CanEdit = false
		serviceRequestInfo.AssignmentStatus = WorkFlowUser
		serviceRequestInfo.ActionStatus = ActionApprovedByHoD
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
		err := Update(dbConnection, WorkHourAssignmentRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		v.emailGenerator(dbConnection, HeadOfDepartmentApproveEmailTemplateType, serviceRequestInfo.CreatedBy, WorkHourAssignmentMyDepartmentRequestComponent, recordId)
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

func (v *WorkHourAssignmentService) handleDepartmentHeadReturn(ctx *gin.Context) {

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	returnRemark := util.InterfaceToString(returnFields["remark"])
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if serviceRequestInfo.CanHODReturn {
		if !v.isEmailTemplateExist(dbConnection, HeadOfDepartmentReturnEmailTemplateType) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		serviceRequestInfo.CanEdit = true
		serviceRequestInfo.CanHODReturn = false
		serviceRequestInfo.CanHODReject = false
		serviceRequestInfo.CanHODApprove = false
		serviceRequestInfo.CanUserSubmit = true
		serviceRequestInfo.AssignmentStatus = WorkFlowUser
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
		err := Update(dbConnection, WorkHourAssignmentRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		v.emailGenerator(dbConnection, HeadOfDepartmentReturnEmailTemplateType, serviceRequestInfo.CreatedBy, WorkHourAssignmentMyDepartmentRequestComponent, recordId)
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

func (v *WorkHourAssignmentService) handleDepartmentHeadReject(ctx *gin.Context) {

	var rejectionRemarkFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&rejectionRemarkFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	rejectionRemarks := util.InterfaceToString(rejectionRemarkFields["remark"])
	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if serviceRequestInfo.CanHODReject {
		if !v.isEmailTemplateExist(dbConnection, HeadOfDepartmentRejectEmailTemplateType) {
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
		serviceRequestInfo.AssignmentStatus = WorkFlowUser
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
		err := Update(dbConnection, WorkHourAssignmentRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		v.emailGenerator(dbConnection, HeadOfDepartmentRejectEmailTemplateType, serviceRequestInfo.CreatedBy, WorkHourAssignmentMyDepartmentRequestComponent, recordId)

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
