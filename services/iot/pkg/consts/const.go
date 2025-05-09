package consts

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

	MessageTable                 = "message"
	IOTRecordTrailTable          = "iot_record_trail"
	IOTComponentTable            = "iot_component"
	IOTNodeTable                 = "iot_node"
	IOTDriverTable               = "iot_driver"
	IOTDataSourceMasterTable     = "iot_datasource_master"
	IOTDriverStatisticsTable     = "iot_driver_statistics"
	IoTDriverComponent           = "iot_driver"
	IOTDriverStatisticsComponent = "iot_driver_statistics"

	IOTEmailTemplateTable      = "iot_email_template"
	IOTEmailTemplateFieldTable = "iot_email_template_field"

	IOTDriverDynamicFieldTable       = "iot_driver_dynamic_field"
	IOTDriverFieldConfigurationTable = "iot_driver_dynamic_field_configuration"

	ErrorGettingObjectsInformation          = 6000
	FailedToDownloadTheImportFileUrl        = 6001
	UnableToReadCSVFile                     = 6002
	ParsingCSVFileFailed                    = 6003
	SchemaIsNotMatchedWithUploadedCSV       = 6004
	ErrorUpdatingObjectInformation          = 6005
	ErrorRemovingObjectInformation          = 6006
	ErrorGettingIndividualObjectInformation = 6007
	ErrorCreatingObjectInformation          = 6008

	ProjectID  = "906d0fd569404c59956503985b330132"
	TimeLayout = "2006-01-02T15:04:05.000Z"
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
