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

	MachinesRecordTrailTable = "machines_record_trail"

	MachineComponentTable              = "machine_component"
	MachineMasterTable                 = "machine_master"
	MachineCategoryTable               = "machine_category"
	MachineSubCategoryTable            = "machine_sub_category"
	MachineParameterTable              = "machine_parameter"
	MachineStatusTable                 = "machine_status"
	MachineSchedulerTable              = "machine_scheduler"
	MachineStatisticsTable             = "machine_statistics"
	MachineHMITable                    = "machine_hmi"
	MachineHMIStopReasonTable          = "machine_hmi_stop_reason"
	MachineDisplaySettingTable         = "machine_display_setting"
	MachineHMISettingSettingTable      = "machine_hmi_setting"
	MachineModuleSettingTable          = "machine_module_setting"
	MachineConnectStatusTable          = "machine_connect_status"
	AssemblyMachineDisplaySettingTable = "assembly_machine_display_setting"
	AssemblyMachineLineTable           = "assembly_machine_lines"

	AssemblyMachineMasterTable        = "assembly_machine_master"
	AssemblyMachineHmiTable           = "assembly_machine_hmi"
	AssemblyMachineHmiSettingTable    = "assembly_machine_hmi_setting"
	AssemblyMachineModuleSettingTable = "assembly_machine_module_setting"
	AssemblyMachineStatisticsTable    = "assembly_machine_statistics"

	MachineEmailTemplateTable     = "machine_email_template"
	MachineEmailMasterFieldsTable = "machine_email_master_fields"

	ToolingMachineMasterTable         = "tooling_machine_master"
	ToolingMachineHmiTable            = "tooling_machine_hmi"
	ToolingMachineHmiSettingTable     = "tooling_machine_hmi_setting"
	ToolingMachineDisplaySettingTable = "tooling_machine_display_setting"
	ToolingMachineStatisticsTable     = "tooling_machine_statistics"
	AssemblyLineTypeTable             = "assembly_line_type"

	AssemblyEquipmentTypeMasterTable = "assembly_equipment_type_master"
	AssemblyEquipmentNameTable       = "assembly_equipment_name"
	MachineWidgetTable               = "machine_widget"
	MachineDashboardTable            = "machine_dashboard"

	MachineMasterComponent         = "machine_master"
	MachineParamComponent          = "machine_param"
	MachineStatusComponent         = "machine_status"
	MachineCategoryComponent       = "machine_category"
	MachineSubCategoryComponent    = "machine_sub_category"
	MachineTimelineEventComponent  = "machine_timeline_event"
	MachineHMIComponent            = "machine_hmi"
	MachineDisplaySettingComponent = "machine_display_setting"
	AssemblyMachineMasterComponent = "assembly_machine_master"
	AssemblyMachineHmiComponent    = "assembly_machine_hmi"
	ToolingMachineHmiComponent     = "tooling_machine_hmi"
	ToolingMachineMasterComponent  = "tooling_machine_master"
	MachineFilterTable             = "machine_filter"

	MachineHMIRejectedComponent         = "machine_hmi_rejected"
	AssemblyMachineHMIRejectedComponent = "assembly_machine_hmi_rejected"
	ToolingMachineHMIRejectedComponent  = "tooling_machine_hmi_rejected"

	MouldMachineBrandTable     = "mould_machine_brand"
	MouldMachineBrandComponent = "mould_machine_brand"

	MouldMachineSettingComponent = "mould_machine_setting"
	MouldMachineSettingTable     = "mould_machine_setting"

	FactoryBuildingComponent = "factory_building"
	HelpStopReasonId         = 1

	InvalidSourceError = "Invalid Source"

	InvalidSchedulePosition = "Invalid Schedule Position"
	InvalidComponent        = 6010

	InvalidScheduleStatus    = 6054
	ErrorGettingActionFields = 6056
	QueryExecutionFailed     = 6057
	InvalidEventId           = 6058
	InvalidMachineId         = 6059

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	DecodingFailed       = 6070
	InvalidMachineStatus = 6071

	ScheduleStatusPreferenceTwo   = 2
	ScheduleStatusPreferenceThree = 3
	ScheduleStatusPreferenceFour  = 4
	ScheduleStatusPreferenceFive  = 5
	ScheduleStatusPreferenceSix   = 6
	ScheduleStatusPreferenceSeven = 7
	ScheduleStatusPreferenceEight = 8

	MachineStatusActive      = 1
	MachineStatusMaintenance = 2
	MachineStatusRepair      = 3
	MachineStatusInactive    = 4

	WorkOrderDoneStatus = 4

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	ProjectID                          = "906d0fd569404c59956503985b330132"
	TimeLayout                         = "2006-01-02T15:04:05.000Z"
	machineConnectStatusLive           = 1
	machineConnectStatusWaitingForFeed = 2
	machineConnectStatusMaintenance    = 3

	ModuleName = "machines"

	Moulding = "moulding"
	Assembly = "assembly"
	Tooling  = "tooling"
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
