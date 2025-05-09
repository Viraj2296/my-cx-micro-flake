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

	MaintenanceRecordTrailTable = "maintenance_record_trail"

	MaintenanceComponentTable                 = "maintenance_component"
	MaintenanceWorkOrderTable                 = "maintenance_work_order"
	MaintenanceWorkOrderTaskTable             = "maintenance_work_order_task"
	MaintenanceWorkOrderStatusTable           = "maintenance_work_order_status"
	MaintenancePreventiveWorkOrderStatusTable = "maintenance_preventive_work_order_status"
	MaintenanceWorkOrderTaskStatusTable       = "maintenance_work_order_task_status"
	MaintenanceFaultCodeTable                 = "maintenance_fault_code"

	MaintenanceWorkOrderTaskComponent             = "maintenance_work_order_task"
	MaintenanceWorkOrderMyTaskComponent           = "my_maintenance_work_order_task"
	MaintenanceWorkOrderMyCorrectiveTaskComponent = "my_maintenance_corrective_work_order_task"
	MaintenanceWorkOrderComponent                 = "maintenance_work_order"
	MaintenanceCorrectiveWorkOrderComponent       = "maintenance_corrective_work_order"
	MaintenanceWorkOrderCorrectiveTaskComponent   = "maintenance_work_order_corrective_task"

	MaintenanceCorrectiveWorkOrderTable     = "maintenance_corrective_work_order"
	MaintenanceWorkOrderCorrectiveTaskTable = "maintenance_work_order_corrective_task"

	MouldMaintenancePreventiveWorkOrderTable     = "mould_maintenance_preventive_work_order"
	MouldMaintenancePreventiveWorkOrderTaskTable = "mould_maintenance_preventive_work_order_task"
	MouldMaintenanceCorrectiveWorkOrderTable     = "mould_maintenance_corrective_work_order"
	MouldMaintenanceCorrectiveWorkOrderTaskTable = "mould_maintenance_corrective_work_order_task"

	MouldMaintenancePreventiveWorkOrderComponent     = "mould_maintenance_preventive_work_order"
	MouldMaintenancePreventiveWorkOrderTaskComponent = "mould_maintenance_preventive_work_order_task"
	MouldMaintenanceCorrectiveWorkOrderComponent     = "mould_maintenance_corrective_work_order"
	MouldMaintenanceCorrectiveWorkOrderTaskComponent = "mould_maintenance_corrective_work_order_task"

	MyMouldMaintenancePreventiveWorkOrderTaskComponent = "my_mould_maintenance_preventive_work_order_task"
	MyMouldMaintenanceCorrectiveWorkOrderTaskComponent = "my_mould_maintenance_corrective_work_order_task"

	MouldCorrectiveMaintenanceJrOptionTable     = "mould_corrective_maintenance_jr_option"
	MouldCorrectiveMaintenanceJrOptionComponent = "mould_corrective_maintenance_jr_option"

	MachineMaintenanceSettingTable     = "machine_maintenance_setting"
	MachineMaintenanceSettingComponent = "machine_maintenance_setting"
	MouldMaintenanceSettingTable       = "mould_maintenance_setting"
	MouldMaintenanceSettingComponent   = "mould_maintenance_setting"

	PreventiveWorkOrderTaskAssignmentNotification      = 1
	WorkOrderTaskAssignmentNotification                = 6
	WorkOrderSupervisorNotication                      = 3
	CorrectiveWorkOrderSupervisorNotification          = 4
	CorrectiveWorkOrderCompletionNotification          = 5
	WorkOrderCompletionNotification                    = 8
	MouldWorkOrderTaskAssignmentNotification           = 9
	MouldPreventiveWorkOrderTaskAssignmentNotification = 10
	MouldCorrectiveWorkOrderAssignmentNotification     = 11
	MouldPreventiveWorkOrderAssignmentNotification     = 12
	MouldCorrectiveWorkOrderCompletionNotification     = 13
	MouldPreventiveWorkOrderCompletionNotification     = 14
	CorrectiveWorkOrderCreateNotification              = 15
	MouldCorrectiveWorkOrderCreateNotification         = 16

	MaintenanceEmailTemplateTable      = "maintenance_email_template"
	MaintenanceEmailTemplateFieldTable = "maintenance_email_template_field"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ProductionContinued    = 1
	ProductionNotContinued = 2

	InvalidScheduleStatus = 6054
	InvalidSourceError    = "Invalid Source"

	ISOTimeLayout                 = "2006-01-02T15:04:05.000Z"
	ScheduleStatusPreferenceThree = 3
	ScheduleStatusPreferenceFour  = 4

	WorkOrderCreated    = 1
	WorkOrderScheduled  = 2
	WorkOrderInProgress = 3
	WorkOrderDone       = 4

	WorkOrderTaskToDo       = 1
	WorkOrderTaskReDo       = 2
	WorkOrderTaskInProgress = 3
	WorkOrderTaskDone       = 4
	WorkOrderTaskApproved   = 5

	// this notification is generated if the status changes
	WorkOrderNotificationEmailTemplateType = 1
	// this will be used to notify users that they have assigned to task
	WorkOrderTaskNotificationEmailTemplateType = 2
	WorkOrderTaskStatusEmailTemplateType       = 3

	AssetClassMoulds           = "moulds"
	AssetClassMachines         = "machines"
	AssetClassAssemblyMachines = "assemblyMachines"
	AssetClassToolingMachines  = "toolingMachines"

	WorkOrderTypePreventive = "Preventive"
	WorkOrderTypeCorrective = "Corrective"

	ModuleName = "maintenance"

	ProjectID = "906d0fd569404c59956503985b330132"

	DailyRepeat   = 1
	WeeklyRepeat  = 2
	MonthlyRepeat = 3

	EmailNotificationPreventiveOrder = 7
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
