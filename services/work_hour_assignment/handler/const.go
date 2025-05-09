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

	WorkHourAssignmentRecordTrailTable = "work_hour_assignment_record_trail"

	WorkHourAssignmentComponentTable               = "work_hour_assignment_component"
	WorkHourAssignmentRequestTable                 = "work_hour_assignment_request"
	WorkHourAssignmentTasksTable                   = "work_hour_assignment_task"
	WorkHourAssignmentMyRequestComponent           = "work_hour_assignment_my_request"
	WorkHourAssignmentMyDepartmentRequestComponent = "work_hour_assignment_my_department_request"

	WorkHourAssignmentJRMasterTable = "work_hour_assignment_master_jr"
	WorkHourAssignmentTLMasterTable = "work_hour_assignment_master_tl"
	WorkHourAssignmentMRMasterTable = "work_hour_assignment_master_mr"

	WorkHourAssignmentRequestComponent        = "work_hour_assignment_request"
	WorkHourAssignmentRequestStatusTable      = "work_hour_assignment_request_status"
	WorkHourAssignmentEmailTemplateTable      = "work_hour_assignment_email_template"
	WorkHourAssignmentEmailTemplateFieldTable = "work_hour_assignment_email_template_field"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	WorkFlowUser = 1
	WorkFlowHOD  = 2

	ActionPendingSubmission = "DRAFT"
	ActionPendingHoD        = "PENDING BY HOD"
	ActionReturnedByHoD     = "RETURNED BY HOD"
	ActionRejectedByHoD     = "REJECTED BY HOD"
	ActionApprovedByHoD     = "APPROVED BY HOD"

	ActionAPIHoDApprove       = "hod_approve"
	ActionAPIHoDReturn        = "hod_return"
	ActionAPIHoDReject        = "hod_reject"
	ActionUserSubmit          = "submit"
	ActionUserAcknowledgement = "acknowledgement"
	ActionAPIITApprove        = "it_approve"
	ActionAPIITReject         = "it_reject"

	ActionAPIExecutionPartyDeliver = "execution_party_deliver"

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
	ModuleName                               = "work_hour_assignment"
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
