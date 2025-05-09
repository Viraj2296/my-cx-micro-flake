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

	ToolingComponentTable         = "tooling_component"
	ToolingProjectTable           = "tooling_project"
	ToolingProjectComponent       = "tooling_project"
	ToolingProjectTaskTable       = "tooling_project_task"
	ToolingProjectTaskComponent   = "tooling_project_task"
	ToolingProjectStatusTable     = "tooling_project_status"
	ToolingProjectTaskStatusTable = "tooling_project_task_status"
	ToolingProjectSprintTable     = "tooling_project_sprint"

	ToolingGatingTypeTab                 = "tooling_gating_type"
	ToolingHotRunnerBrandTable           = "tooling_hot_runner_brand"
	ToolingRunnerTypeTab                 = "tooling_runner_type"
	ToolingHotRunnerConnectorTypeTable   = "tooling_hot_runner_connector_type"
	ToolingHotRunnerControllerTable      = "tooling_hot_runner_controller"
	ToolingHRControllerTable             = "tooling_hr_controller"
	ToolingHotRunnerControllerBrandTable = "tooling_hot_runner_controller_brand"
	ToolingCoolingFittingTable           = "tooling_cooling_fitting"

	ToolingProgramTable         = "tooling_program"
	ToolingStatusTable          = "tooling_status"
	ToolingLocationTable        = "tooling_location"
	ToolingMouldVendorTable     = "tooling_mould_vendor"
	ToolingProjectTaskListTable = "tooling_project_task_check_list"

	ToolingProjectMasterComponent     = "tooling_project_master"
	ToolingTaskComponent              = "tooling_task"
	ToolingProjectStatusComponent     = "tooling_project_status"
	ToolingProjectTaskStatusComponent = "tooling_project_task_status"
	ToolingProjectSprintComponent     = "tooling_project_sprint"
	ToolingProjectMyTaskComponent     = "tooling_project_my_task"

	ToolingEmailTemplateTable      = "tooling_email_template"
	ToolingEmailTemplateFieldTable = "tooling_email_template_field"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008
	ErrorCreatingSchedule                   = 5009
	AlreadyScheduled                        = 5010

	InvalidEventId = 6058

	FieldValidationFailed = 6010

	ToolingRecordTrailTable = "tooling_record_trail"

	ToolingProjectTaskAssignmentEmailTemplateType = 1
	ToolingProjectTaskStatusEmailTemplateType     = 2

	ProjectTaskToDO       = 1
	ProjectTaskReDO       = 2
	ProjectTaskInProgress = 3
	ProjectTaskDone       = 4
	ProjectTaskApproved   = 5

	ProjectDRAFT       = 1
	ProjectCREATED     = 2
	ProjectDFM         = 3
	ProjectDESIGN      = 4
	ProjectFABRICATION = 5
	ProjectAPPROVAL    = 6
	ProjectCLOSED      = 7
	ProjectCANCELLED   = 8
	ProjectPENDING     = 9

	ProjectCreated = 1

	ProjectSprintInProgress = 2
	ProjectSprintCompleted  = 3

	ProjectID = "906d0fd569404c59956503985b330132"

	ISOTimeLayout = "2006-01-02T15:04:05.000Z"

	ModuleName = "tooling"
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
