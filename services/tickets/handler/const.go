package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	TicketsServiceRecordTrailTable = "tickets_service_record_trail"

	TicketsServiceComponentTable               = "tickets_service_component"
	TicketsServiceRequestTable                 = "tickets_service_request"
	TicketsServiceRequestCategoryTable         = "tickets_service_request_category"
	TicketsServiceRequestSubCategoryTable      = "tickets_service_request_sub_category"
	TicketsServiceMyRequestComponent           = "tickets_service_my_request"
	TicketsServiceMyDepartmentRequestComponent = "tickets_service_my_department_request"
	TicketsServiceMyReviewRequestComponent     = "tickets_service_my_review_request"
	TicketsServiceMyExecutionRequestComponent  = "tickets_service_my_execution_request"

	TicketsServiceRequestComponent        = "tickets_service_request"
	TicketsServiceRequestStatusTable      = "tickets_service_request_status"
	TicketsServiceEmailTemplateTable      = "tickets_service_email_template"
	TicketsServiceEmailTemplateFieldTable = "tickets_service_email_template_field"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	WorkFlowUser           = 1
	WorkFlowHOD            = 2
	WorkFlowReviewParty    = 3
	WorkFlowExecutionParty = 4

	ActionPendingSubmission     = "Pending Submission"
	ActionPendingHoD            = "Pending HoD"
	ActionReturnedByHoD         = "Returned By HoD"
	ActionRejectedByHoD         = "Rejected By HoD"
	ActionPendingITReview       = "Pending IT Review"
	ActionRejectedByITReview    = "Rejected By IT Review"
	ActionPendingExecParty      = "Pending Execution Party"
	ActionPendingAcknowledgment = "Pending Acknowledgement"
	ActionCaseReOpened          = "Reopened"
	ActionClosed                = "Closed"

	ActionAPIHoDApprove       = "hod_approve"
	ActionAPIHoDReturn        = "hod_return"
	ActionAPIHoDReject        = "hod_reject"
	ActionUserSubmit          = "submit"
	ActionUserAcknowledgement = "acknowledgement"
	ActionAPIITApprove        = "it_approve"
	ActionAPIITReject         = "it_reject"

	ActionAPIExecutionPartyDeliver = "execution_party_deliver"

	ModuleName = "facility"

	ProjectID = "906d0fd569404c59956503985b330132"

	/*
		A1	Pending HOD
		A2	Returned by HOD
		A3	Rejected by HOD
		B1	Pending IT Review
		B2	Rejected by IT Review
		C1	Pending Exec Party
		C2	Case re-assigned
		C3	Pending acknowledgement
		C4	Case reopened
	*/

	UserSubmitRequestEmailTemplateType       = 1
	HeadOfDepartmentApproveEmailTemplateType = 2
	HeadOfDepartmentReturnEmailTemplateType  = 3
	HeadOfDepartmentRejectEmailTemplateType  = 4
	ITApproveEmailTemplateType               = 5
	ITRejectEmailTemplateType                = 6
	ExecutionPartyDeliverTemplateType        = 7
	UserAcknowledgeEmailTemplateType         = 8
	SubmitForExecutionEmailTemplateType      = 9
)

func getError(errorString string) error {
	return errors.New(errorString)
}

func sendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}
func sendArchiveFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Archived Failed",
			Description: "Internal system error during archive process. This is normally happen when the system is not configured properly. Please report to system administrator",
		})
	return
}
