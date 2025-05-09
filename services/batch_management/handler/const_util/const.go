package const_util

import (
	"errors"
)

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	BatchManagementRecordTrailTable = "batch_management_record_trail"

	BatchManagementEmailTemplateTable      = "batch_management_email_template"
	BatchManagementEmailTemplateFieldTable = "batch_management_email_template_field"
	BatchManagementComponentTable          = "batch_management_component"
	BatchManagementRawMaterialTable        = "batch_management_raw_material"
	BatchManagementMouldTable              = "batch_management_mould"
	BatchManagementPrinterTable            = "batch_management_printer"
	ManufacturingVendorMasterTable         = "manufacturing_vendor_master"
	BatchManagementRawMaterialType         = "batch_management_material_type"

	ComponentBatchManagementRawMaterial = "batch_management_raw_material"
	RasinType                           = 1

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
	ModuleName = "moulds"
)

func GetError(errorString string) error {
	return errors.New(errorString)
}
