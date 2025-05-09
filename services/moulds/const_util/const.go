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

	MouldsRecordTrailTable = "moulds_record_trail"

	MouldComponentTable               = "mould_component"
	MouldMasterTable                  = "mould_master"
	MouldCategoryTable                = "mould_category"
	MouldSubCategoryTable             = "mould_sub_category"
	MouldStatusTable                  = "mould_status"
	MouldTestStatusTable              = "mould_test_status"
	MouldTestRequestTable             = "mould_test_request"
	MouldTestRequestComponent         = "mould_test_request"
	PartMasterTable                   = "part_master"
	MouldSettingTable                 = "mould_setting"
	MouldMasterComponent              = "mould_master"
	MouldRequestTestComponent         = "mould_test_request"
	MouldStatusComponent              = "mould_status"
	MouldCategoryComponent            = "mould_category"
	MouldSubCategoryComponent         = "mould_sub_category"
	MyMouldTestRequestComponent       = "my_mould_test_request"
	PartMasterComponent               = "part_master"
	MouldManualShotCountTable         = "mould_manual_shot_count"
	MouldMachineMasterTable           = "mould_machine_master"
	MouldManualShotCountComponent     = "mould_manual_shot_count"
	MouldSettingComponent             = "mould_setting"
	MouldShoutCountViewTable          = "mould_shot_count_view"
	MouldShoutCountViewComponent      = "mould_shot_count_view"
	MouldModificationHistoryTable     = "mould_modification_history"
	MouldModificationHistoryComponent = "mould_modification_history"

	MouldEmailTemplateTable      = "mould_email_template"
	MouldEmailTemplateFieldTable = "mould_email_template_field"

	ErrorGettingObjectsInformation          = 6000
	FailedToDownloadTheImportFileUrl        = 6001
	UnableToReadCSVFile                     = 6002
	ParsingCSVFileFailed                    = 6003
	SchemaIsNotMatchedWithUploadedCSV       = 6004
	ErrorUpdatingObjectInformation          = 6005
	ErrorRemovingObjectInformation          = 6006
	ErrorGettingIndividualObjectInformation = 6007
	ErrorCreatingObjectInformation          = 6008

	FieldValidationFailed = 6010

	UnableToCreateExportFile = 6009
	ErrorGettingActionFields = 6056
	InvalidMouldComponent    = 6010
	InvalidScheduleStatus    = 6054

	MouldTestWorkFlowTooling           = 1
	MouldTestWorkFlowPlanner           = 2
	MouldTestWorkFlowProcessDepartment = 3
	MouldTestWorkFlowQuality           = 4

	//Mould status const
	MouldStatusPrototyping      = 1
	MouldStatusQualification    = 2
	MouldStatusCustomerApproval = 3
	MouldStatusActive           = 4
	MouldStatusMaintenance      = 5
	MouldStatusRepair           = 6
	MouldStatusScrap            = 7
	MouldStatusWaiver           = 8
	MouldStatusInterimApproval  = 9
	MouldStatusExport           = 10

	MouldTestRequestActionCreated        = "Created"
	MouldTestRequestActionScheduled      = "Scheduled"
	MouldTestRequestActionTestInProgress = "Test In-Progress"
	MouldTestRequestActionPassed         = "Passed"
	MouldTestRequestActionFailed         = "Failed"
	MouldTestRequestActionApproved       = "Approved"

	InactiveMould                     = "Inactive mould"
	InvalidSourceError                = "Invalid Source"
	InactiveMouldDescription          = "System couldn't able to create mould test request"
	InactiveMouldWorkOrderDescription = "System couldn't able to create work order due to invalid mould status"
	InvalidMouldStatusDescription     = "Getting mould status from the system failed, check mould statuses are defined"
	DuplicateRecordFound              = "Error in Releasing the Mould Test Request"
	ErrorCannotModifyData             = "Error in changing modification"

	ActionUpdateSchedulerEvent      = "update_scheduler_event"
	ActionReleaseMouldTestRequest   = "release"
	ActionUnReleaseMouldTestRequest = "hold"
	ActionCustomerConfirm           = "customer_confirm"

	ActionCheckInMouldTestRequest  = "checkIn"
	ActionCompleteMouldTestRequest = "complete"
	ActionApproveMouldTestRequest  = "approve_test_request"
	ActionModifyMouldMasterRequest = "modify"

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	MouldMasterFieldMouldStatus           = "mouldStatus"
	MouldMasterFieldCanSubmitTestRequest  = "canSubmitTestRequest"
	MouldMasterFieldCanCreateWorkOrder    = "canCreateWorkOrder"
	MouldMasterFieldCanCustomerApprove    = "canCustomerApprove"
	MouldMasterFieldShotCount             = "shotCount"
	MouldMasterFieldToolLife              = "toolLife"
	MouldMasterFieldIsNotificationSend    = "isNotificationSend"
	MouldMasterFieldMouldLifeNotification = "mouldLifeNotification"
	ShotCountRatio                        = 0.75

	CommonFieldLastUpdatedBy = "lastUpdatedBy"
	CommonFieldLastUpdatedAt = "lastUpdatedAt"

	MouldTestRequestScheduledTemplateType = 1
	MouldTestRequestSubmittedTemplateType = 2
	MouldTestRequestApprovedTemplateType  = 3
	MouldShotCountTemplateType            = 4
	DuplicateValueFound                   = 800
	CannotModifyData                      = 801

	QualificationStatus = 2

	ProjectID  = "906d0fd569404c59956503985b330132"
	ModuleName = "moulds"
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
