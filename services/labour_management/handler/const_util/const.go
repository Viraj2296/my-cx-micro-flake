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

	LabourManagementRecordTrailTable = "labour_management_record_trail"

	LabourManagementEmailTemplateTable      = "labour_management_email_template"
	LabourManagementEmailTemplateFieldTable = "labour_management_email_template_field"
	LabourManagementComponentTable          = "labour_management_component"
	LabourManagementShiftMasterTable        = "labour_management_shift_master"
	LabourManagementAttendanceTable         = "labour_management_attendance"
	LabourManagementShiftStatusTable        = "labour_management_shift_status"
	LabourManagementShiftProductionTable    = "labour_management_shift_production"
	LabourManagementShiftTemplateTable      = "labour_management_shift_template"
	LabourManagementSettingTable            = "labour_management_setting"
	LabourManagementShiftProductionAttendanceTable = "labour_management_shift_production_attendance"

	LabourManagementSettingComponent        = "labour_management_setting"
	ErrorGettingObjectsInformation          = 6000
	FailedToDownloadTheImportFileUrl        = 6001
	UnableToReadCSVFile                     = 6002
	ParsingCSVFileFailed                    = 6003
	SchemaIsNotMatchedWithUploadedCSV       = 6004
	ErrorUpdatingObjectInformation          = 6005
	ErrorRemovingObjectInformation          = 6006
	ErrorGettingIndividualObjectInformation = 6007
	ErrorCreatingObjectInformation          = 6008
	ScheduleStatusPreferenceFive            = 5
	ScheduleStatusPreferenceSix             = 6
	FieldValidationFailed                   = 6010

	UnableToCreateExportFile = 6009
	ErrorGettingActionFields = 6056
	InvalidMouldComponent    = 6010
	InvalidScheduleStatus    = 6054

	CreatedStatus  = 1
	PrintingStatus = 2
	PrintedStatus  = 3

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	CommonFieldLastUpdatedBy = "lastUpdatedBy"
	CommonFieldLastUpdatedAt = "lastUpdatedAt"

	MouldTestRequestScheduledTemplateType = 1
	MouldTestRequestSubmittedTemplateType = 2
	MouldTestRequestApprovedTemplateType  = 3
	MouldShotCountTemplateType            = 4

	ProjectID  = "906d0fd569404c59956503985b330132"
	ModuleName = "labour_management"

	ShiftStatusActive    = 1
	ShiftStatusCompleted = 2
	ShiftStatusPending   = 3
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
