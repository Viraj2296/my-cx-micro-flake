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

	IncidentRecordTrailTable = "incident_record_trail"

	IncidentComponentTable = "incident_component"

	IncidentSafetyCategoryTable       = "incident_safety_category"
	IncidentQualityCategoryTable      = "incident_quality_category"
	IncidentDeliveryCategoryTable     = "incident_delivery_category"
	IncidentInventoryCategoryTable    = "incident_inventory_category"
	IncidentProductivityCategoryTable = "incident_productivity_category"

	IncidentTargetTable = "incident_target"

	IncidentSafetyTable       = "incident_safety"
	IncidentQualityTable      = "incident_quality"
	IncidentDeliveryTable     = "incident_delivery"
	IncidentInventoryTable    = "incident_inventory"
	IncidentProductivityTable = "incident_productivity"

	IncidentSafetyComponent       = "incident_safety"
	IncidentQualityComponent      = "incident_quality"
	IncidentDeliveryComponent     = "incident_delivery"
	IncidentInventoryComponent    = "incident_inventory"
	IncidentProductivityComponent = "incident_productivity"

	IncidentSafetyCategoryComponent       = "incident_safety_category"
	IncidentQualityCategoryComponent      = "incident_quality_category"
	IncidentDeliveryCategoryComponent     = "incident_delivery_category"
	IncidentInventoryCategoryComponent    = "incident_inventory_category"
	IncidentProductivityCategoryComponent = "incident_productivity_category"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	ScheduleStatusPreferenceThree = 3
	ScheduleStatusPreferenceFour  = 4

	ModuleName   = "incident"
	ProjectID    = "906d0fd569404c59956503985b330132"
	Safety       = 1
	Quality      = 2
	Delivery     = 3
	Inventory    = 4
	Productivity = 5
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
