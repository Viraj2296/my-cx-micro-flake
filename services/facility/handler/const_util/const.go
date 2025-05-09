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

	FacilityServiceRecordTrailTable = "facility_service_record_trail"

	FacilityServiceComponentTable          = "facility_service_component"
	FacilityServiceRequestTable            = "facility_service_request"
	FacilityServiceRequestCategoryTable    = "facility_service_request_category"
	FacilityServiceRequestSubCategoryTable = "facility_service_request_sub_category"
	FacilityServiceWorkflowEngineTable     = "facility_service_workflow_engine"
	FacilityServiceAdminSettingTable       = "facility_service_admin_setting"

	FacilityServiceMyRequestComponent                = "facility_service_my_request"
	FacilityServiceMyDepartmentRequestComponent      = "facility_service_my_department_request"
	FacilityServiceMyReviewRequestComponent          = "facility_service_my_review_request"
	FacilityServiceMyExecutionRequestComponent       = "facility_service_my_execution_request"
	FacilityServiceMyEHSManagerRequestComponent      = "facility_service_my_ehs_manager_request"
	FacilityServiceMyTechnicianRequestComponent      = "facility_service_my_execution_request"
	FacilityServiceMyFacilityManagerRequestComponent = "facility_service_my_facility_manager_request"
	FacilityServiceRequestComponent                  = "facility_service_request"
	FacilityServiceAdminSettingComponent             = "facility_service_admin_setting"

	FacilityServiceRequestStatusTable        = "facility_service_request_status"
	FacilityServiceEmailTemplateTable        = "facility_service_email_template"
	FacilityServiceEmailTemplateFieldTable   = "facility_service_email_template_field"
	FacilityServiceCategoryTemplateTable     = "facility_service_category_template"
	FacilityServiceAllSafetyRequestComponent = "facility_service_all_safety_request"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	WorkFlowUser            = 1
	WorkFlowHOD             = 2
	WorkFlowFacilityManager = 3
	WorkFlowEHSManager      = 4
	WorkFlowTechnician      = 5

	ActionPendingSubmission         = "Pending Submission"
	ActionPendingHoD                = "Pending HoD"
	ActionReturnedByHoD             = "Returned By HoD"
	ActionRejectedByHoD             = "Rejected By HoD"
	ActionRejectedByFacilityManager = "Rejected By Facility Manager"
	ActionPendingEHSManager         = "Pending EHS Manager"
	ActionPendingFacilityManager    = "Pending Facility Manager"
	ActionPendingExecParty          = "Pending Execution Party"
	ActionPendingTechnicianParty    = "Pending Technician Party"
	ActionPendingAcknowledgment     = "Pending Acknowledgement"
	ActionCaseReOpened              = "Reopened"
	ActionRejected                  = "Rejected By EHS Manager"
	ActionClosed                    = "Closed"
	ActionPendingTesting            = "Pending Testing"
	ActionCanceled                  = "Cancelled"

	ActionTested                    = "Tested"

	ActionAPIHoDApprove            = "hod_approve"
	ActionAPIHoDReturn             = "hod_return"
	ActionAPIHoDReject             = "hod_reject"
	ActionUserSubmit               = "submit"
	ActionAPIAssignMyself          = "assign_myself"
	ActionUserAcknowledgement      = "acknowledgement"
	ActionAPISapApprove            = "sap_approve"
	ActionAPICancel                = "cancel"
	ActionAPEHSManagerApprove      = "ehs_manager_approve"
	ActionAPEHSManagerReject       = "ehs_manager_reject"
	ActionReassignExecution        = "reassign_execution"
	ActionFacilityManagerApprove   = "facility_manager_approve"
	ActionFacilityManagerReject    = "facility_manager_reject"
	ActionAPIExecutionPartyDeliver = "execution_party_deliver"

	// workflow engine const
	UserTechnicianUserWorkflowEngine                             = 1
	UserHoDFacilityManagerEHSManagerTechnicianUserWorkflowEngine = 2

	ActionAPIRollBack = "roll_back"

	ProjectID = "906d0fd569404c59956503985b330132"

	RemarkTypeRichTextEditor = "richTextEditor"

	UserSubmitRequestEmailTemplateType       = 1
	HeadOfDepartmentApproveEmailTemplateType = 2
	HeadOfDepartmentReturnEmailTemplateType  = 3
	HeadOfDepartmentRejectEmailTemplateType  = 4
	FacilityManagerApproveEmailTemplateType  = 5
	FacilityCanceledRequestEmailTemplateType = 18

	EHSManagerApproveEmailTemplateType                    = 7
	EHSManagerRejectEmailTemplateType                     = 8
	ExecutionPartyAssignTemplateType                      = 9
	ExecutionPartyDeliverTemplateType                     = 10
	UserAcknowledgeEmailTemplateType                      = 9
	FacilityManagerApprovalForTechnicianEmailTemplateType = 11
	FacilityManagerApprovalForEHSManagerEmailTemplateType = 12
	FacilityManagerRejectEmailTemplateType                = 13
	FacilityManagerRejectToHODEmailTemplateType           = 14
	UserSubmitRequestRoutingOneEmailTemplateType          = 15
	UsedSubmissionAdminSettingEmailTemplateType           = 16
	ExecutionAdminSettingEmailTemplateType                = 17
	AssignExecutionerEmailTemplateType                    = 18
	ModuleName                                            = "facility"
	InvalidSchema                                         = 5009
	ItExecutionPartyRole                                  = 13

	UserHodExecutionUser = 1
	UserHodExecutionHod  = 2

	UserHodExecutionReviewer = 3

	UserExecutionUser     = 1
	UserExecutionReviewer = 2
	UserExecutionUserAck  = 2

	TechniciansFromWorkflow1 = 2

	UserHodSapManager = 2
	UserHODITManager  = 2
	FacilityManager   = 6

	HodWorkFlow        = 1
	ExecutionWorkFlow  = 2
	EhsManagerWorkFlow = 3
	TechnicianWorkFlow = 4

	HodWorkFlowExecutionEntry       = 3
	ExecutionWorkFlowExecutionEntry = 2

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
