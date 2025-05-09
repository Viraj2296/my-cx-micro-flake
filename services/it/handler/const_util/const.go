package const_util

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

	ITServiceRecordTrailTable = "it_service_record_trail"

	ITServiceComponentTable                 = "it_service_component"
	ITServiceRequestTable                   = "it_service_request"
	ITServiceRequestCategoryTable           = "it_service_request_category"
	ITServiceRequestSubCategoryTable        = "it_service_request_sub_category"
	ITServiceSAPChangeReasonsTable          = "it_service_sap_change_reason"
	ITServiceSAPAuthorizationFunctionsTable = "it_service_sap_authorization_function"
	ITServiceWorkflowEngineTable            = "it_service_workflow_engine"

	ITServiceMyRequestComponent             = "it_service_my_request"
	ITServiceMyDepartmentRequestComponent   = "it_service_my_department_request"
	ITServiceMyReviewRequestComponent       = "it_service_my_review_request"
	ITServiceMyExecutionRequestComponent    = "it_service_my_execution_request"
	ITServiceMySapRequestComponent          = "it_service_my_sap_request"
	ITServiceAllSAPRequestComponent         = "it_service_all_sap_request"
	ITServiceMyITManagementRequestComponent = "it_service_my_it_manager_request"

	ITServiceRequestComponent        = "it_service_request"
	ITServiceRequestStatusTable      = "it_service_request_status"
	ITServiceEmailTemplateTable      = "it_service_email_template"
	ITServiceEmailTemplateFieldTable = "it_service_email_template_field"
	IITServiceCategoryTemplateTable  = "it_service_category_template"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	WorkFlowUser           = 1
	WorkFlowHOD            = 2
	WorkFlowSapManager     = 3
	WorkFlowExecutionParty = 4
	WorkFlowITManager      = 5

	ActionPendingSubmission     = "Pending Submission"
	ActionPendingHoD            = "Pending HoD"
	ActionReturnedByHoD         = "Returned By HoD"
	ActionRejectedByHoD         = "Rejected By HoD"
	ActionPendingSapManager     = "Pending Sap Manager"
	ActionRejectedBySap         = "Rejected By Sap Manager"
	ActionPendingExecParty      = "Pending Execution Party"
	ActionPendingAcknowledgment = "Pending Acknowledgement"
	ActionCaseReOpened          = "Reopened"
	ActionRejected              = "Rejected By IT Manager"
	ActionClosed                = "Closed"
	ActionPendingTesting        = "Pending Testing"
	ActionCanceled              = "Canceled"
	ActionTested                = "Tested"
	ActionPendingITManager      = "Pending IT Manager"

	ActionAPIHoDApprove       = "hod_approve"
	ActionAPIHoDReturn        = "hod_return"
	ActionAPIHoDReject        = "hod_reject"
	ActionUserSubmit          = "submit"
	ActionAPIAssignMyself     = "assign_myself"
	ActionUserAcknowledgement = "acknowledgement"
	ActionAPISapApprove       = "sap_approve"
	ActionAPISapReject        = "sap_reject"
	ActionAPICancel           = "cancel"
	ActionAPIITManagerApprove = "it_manager_approve"
	ActionAPIITManagerReject  = "it_manager_reject"
	ActionReassignExecution   = "reassign_execution"
	ActionTransferExecution   = "transfer_execution"

	ActionAPIExecutionPartyDeliver = "execution_party_deliver"

	SAPManagerWorkFlowEngine = 3
	ITManagerWorkFlowEngine  = 4
	HODRoutingWorkflowEngine = 2
	HodRoutingManual         = 1

	ActionAPIRollBack = "roll_back"

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

	UserSubmitRequestEmailTemplateType            = 1
	HeadOfDepartmentApproveEmailTemplateType      = 2
	HeadOfDepartmentReturnEmailTemplateType       = 3
	HeadOfDepartmentRejectEmailTemplateType       = 4
	ITApproveEmailTemplateType                    = 5
	ITRejectEmailTemplateType                     = 6
	ExecutionPartyDeliverTemplateType             = 7
	UserAcknowledgeEmailTemplateType              = 8
	SubmitForExecutionEmailTemplateType           = 9
	UserSubmitExecutionRequestEmailTemplateType   = 10
	SapApproveEmailTemplateType                   = 11
	HeadOfDepartmentApproveToSapEmailTemplateType = 12
	UserTestedEmailTemplateType                   = 13
	HeadOfDepartmentApproveToITEmailTemplateType  = 14
	ITManagerApprovalEmailTemplateType            = 15
	ITManagerApprovalRejectEmailTemplateType      = 16
	SapRejectToHODEmailTemplateType               = 17
	SapRejectToUserEmailTemplateType              = 18
	AssignExecutionerEmailTemplateType            = 19

	ModuleName           = "it"
	InvalidSchema        = 5009
	ItExecutionPartyRole = 13

	UserHodExecutionUser = 1
	UserHodExecutionHod  = 2

	UserHodExecutionReviewer = 3

	UserExecutionUser     = 1
	UserExecutionReviewer = 2
	UserExecutionUserAck  = 2

	UserHodSapManager = 2
	UserHODITManager  = 2

	HodWorkFlow        = 1
	ExecutionWorkFlow  = 2
	SapManagerWorkFlow = 3
	ITManagerWorkFlow  = 4

	HodWorkFlowExecutionEntry        = 3
	ExecutionWorkFlowExecutionEntry  = 2
	SapManagerWorkFlowExecutionEntry = 3
	ITManagerWorkFlowExecutionEntry

	InvalidPayload       = 5010
	ErrorConvertingField = 5009
)

func GetError(errorString string) error {
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
