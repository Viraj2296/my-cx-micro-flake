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

	TraceabilityRecordTrailTable = "traceability_record_trail"

	TraceabilityEmailTemplateTable      = "traceability_email_template"
	TraceabilityEmailTemplateFieldTable = "traceability_email_template_field"
	TraceabilityComponentTable          = "traceability_component"
	TraceabilityOrdersTable             = "traceability_order"

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

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	CommonFieldLastUpdatedBy = "lastUpdatedBy"
	CommonFieldLastUpdatedAt = "lastUpdatedAt"

	MouldTestRequestScheduledTemplateType = 1
	MouldTestRequestSubmittedTemplateType = 2
	MouldTestRequestApprovedTemplateType  = 3
	MouldShotCountTemplateType            = 4

	ProjectID  = "906d0fd569404c59956503985b330132"
	ModuleName = "traceability"
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
