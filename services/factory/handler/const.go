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

	FactoryRecordTrailTable = "factory_record_trail"

	FactoryComponentTable         = "it_service_component"
	FactorySiteTable              = "factory_site"
	FactoryLocationTable          = "factory_location"
	FactoryPlantTable             = "factory_plant"
	FactoryDepartmentTable        = "factory_department"
	FactoryDepartmentSectionTable = "factory_department_section"
	FactoryGeneralFacilityTable   = "factory_general_facility"
	FactoryCustomerAssetTable     = "factory_customer_asset"
	FactoryBuildingTable          = "factory_building"
	FactoryUnitTable              = "factory_unit"
	FactoryLevelTable             = "factory_level"
	FactoryAreaTable              = "factory_area"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	ProjectID = "906d0fd569404c59956503985b330132"

	ScheduleStatusPreferenceThree = 3
	ScheduleStatusPreferenceFour  = 4

	ModuleName = "factory"
)

type InvalidRequest struct {
	Message string `json:"message"`
}

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
