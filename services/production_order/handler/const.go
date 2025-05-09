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

	ProductionOrderRecordTrailTable = "production_order_record_trail"

	ProductionOrderComponentTable = "production_order_component"
	ProductionOrderMasterTable    = "production_order_master"
	ScheduledOrderEventTable      = "scheduled_order_event"
	MaterialMasterTable           = "material_master"
	ProductionOrderStatusTable    = "production_order_status"

	AssemblyProductionOrderTable     = "assembly_production_order"
	AssemblyScheduledOrderEventTable = "assembly_scheduled_order_event"

	ScheduledOrderEventComponent         = "scheduled_order_event"
	ProductionOrderMasterComponent       = "production_order_master"
	AssemblyScheduledOrderEventComponent = "assembly_scheduled_order_event"
	ToolingOderMasterComponent           = "tooling_order_master"
	ToolingPartMasterComponent           = "tooling_part_master"
	ToolingScheduledOrderEventComponent  = "tooling_scheduled_order_event"

	InvalidScheduleEventStatusError = "Invalid Event Status"

	ToolingBomMasterTable           = "tooling_bom_master"
	ToolingOrderMasterTable         = "tooling_order_master"
	ToolingPartMasterTable          = "tooling_part_master"
	ToolingScheduledOrderEventTable = "tooling_scheduled_order_event"

	ScheduleEventNotFound = "Requested Event Not Found"
	ObjectNotFound        = "Requested Object Not Found"
	FieldNotExist         = "Field Not Found"
	DeleteEventFailed     = "Deleting Event Failed"

	ScheduleEventNotFoundDescription      = "System couldn't able to the requested event information, Is this a bot sending request?"
	FieldNotFoundDescription              = "Internal system error, Requested field is not available in existing object"
	ScheduleIsAlreadyReleaseDescription   = "This schedule is already release, please make it hold before proceed remove"
	InvalidProductionOrderDescription     = "Invalid production order, relationship missing between scheduled orders and parent order"
	DeleteScheduledEventFailedDescription = "Internal system error removing event corresponding to request scheduled order"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008
	ErrorCreatingSchedule                   = 5009
	AlreadyScheduled                        = 5010

	FieldNotFound     = 5010
	InvalidFieldValue = 5011

	FieldValidationFailed = 5012

	ScheduleStatusPreferenceOne   = 1
	ScheduleStatusPreferenceTwo   = 2
	ScheduleStatusPreferenceThree = 3
	ScheduleStatusPreferenceFour  = 4
	ScheduleStatusPreferenceFive  = 5
	ScheduleStatusPreferenceSix   = 6
	ScheduleStatusPreferenceSeven = 7
	ScheduleStatusPreferenceEight = 8

	InvalidSchedulePosition = "Invalid Schedule Position"
	InvalidScheduleStatus   = 6054
	InvalidSourceError      = "Invalid Source"

	OperatingZone = "Asia/Singapore"

	InvalidTimeRange               = "Invalid Time Range"
	PartiallyScheduled             = "Quantity Partially Scheduled"
	BalanceQuantityGreaterThanZero = "Automatic scheduled is not enabled orders where balance quantity is > 0"

	TimeRangeErrorDescription = "End time is greater than start time, please check the time range"

	TimeLayout = "2006-01-02T15:04:05.000Z"

	ModuleName = "production_order"

	ProjectID = "906d0fd569404c59956503985b330132"
)

func getError(errorString string) error {
	return errors.New(errorString)
}

func getDetailedError(header string, description string) *response.DetailedError {
	detailedError := response.DetailedError{
		Header:      header,
		Description: description,
	}
	return &detailedError
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
