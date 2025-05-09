package handler

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/datatypes"
)

type KanbanTaskStatusRequest struct {
	TaskStatus int `json:"taskStatus"`
}

type ToolingComponent struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingProgram struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingStatus struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingLocation struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingMouldVendor struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingProjectTaskCheckList struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingProject struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingGatingType struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingHotRunnerBrand struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingRunnerType struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingHotRunnerConnectorType struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingHotRunnerController struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingHRController struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingHotRunnerControllerBrand struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingCoolingFitting struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingProjectInfo struct {
	ProjectReferenceId    string          `json:"projectReferenceId"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	ToolerId              int             `json:"toolerId"`
	MaterialId            int             `json:"materialId"`
	PartId                int             `json:"partId"`
	Customer              string          `json:"customer"`
	ProductionOrderNumber string          `json:"productionOrderNumber"`
	AmountForecast        float64         `json:"amountForecast"`
	AllocatedBudget       float64         `json:"allocatedBudget"`
	PlannedStartDate      string          `json:"plannedStartDate"`
	PlannedEndDate        string          `json:"plannedEndDate"`
	ActualStartDate       string          `json:"actualStartDate"`
	ActualEndDate         string          `json:"actualEndDate"`
	Status                int             `json:"status"`
	TotalDuration         int             `json:"totalDuration"`
	Remarks               string          `json:"remarks"`
	CreatedAt             string          `json:"createdAt"`
	LastUpdatedAt         string          `json:"lastUpdatedAt"`
	CanCancel             bool            `json:"canCancel"`
	CanClose              bool            `json:"canClose"`
	CanApprove            bool            `json:"canApprove"`
	CompletedPercentage   int             `json:"completedPercentage"`
	LastYearAmountUsed    float64         `json:"lastYearAmountUsed"`
	TotalAmountUsed       float64         `json:"totalAmountUsed"`
	TotalForcastAmount    float64         `json:"totalForcastAmount"`
	CompletionBF          int             `json:"completionBF"`
	Supervisor            int             `json:"supervisor"`
	ActionRemarks         []ActionRemarks `json:"actionRemarks"`
	ObjectStatus          string          `json:"objectStatus"`
	IsKickOff             bool            `json:"isKickOff"`
}

func (v *ToolingProject) getToolingProjectInfo() *ToolingProjectInfo {
	itServiceRequestInfo := ToolingProjectInfo{}
	json.Unmarshal(v.ObjectInfo, &itServiceRequestInfo)
	return &itServiceRequestInfo
}

type ToolingProjectTask struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingProjectStatus struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingTaskInfo struct {
	// we need to have reference to project
	ProjectId        int             `json:"projectId"`
	TaskReferenceId  string          `json:"taskReferenceId"`
	Description      string          `json:"description"`
	AssignedUserId   int             `json:"assignedUserId"`
	Name             string          `json:"name"`
	DepartmentId     int             `json:"departmentId"`
	PlannedStartDate string          `json:"plannedStartDate"`
	PlannedEndDate   string          `json:"plannedEndDate"`
	ActualStartDate  string          `json:"actualStartDate"`
	ActualEndDate    string          `json:"actualEndDate"`
	Duration         int             `json:"duration"`
	Status           int             `json:"status"`
	ActionRemarks    []ActionRemarks `json:"actionRemarks"`
	CreatedAt        string          `json:"createdAt"`
	CreatedBy        int             `json:"createdBy"`
	LastUpdatedAt    string          `json:"lastUpdatedAt"`
	LastUpdatedBy    int             `json:"lastUpdatedBy"`
	CanCheckIn       bool            `json:"canCheckIn"`
	CanCheckOut      bool            `json:"canCheckOut"`
	CanApprove       bool            `json:"canApprove"`
	CanReject        bool            `json:"canReject"`
	Remark           string          `json:"remark"`
	ObjectStatus     string          `json:"objectStatus"`
	TargetDate       string          `json:"targetDate"`
}

func (v *ToolingTaskInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (v *ToolingProjectInfo) DatabaseSerialize() map[string]interface{} {
	updatingData := make(map[string]interface{})
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (v *ToolingProjectInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *ToolingTaskInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *ToolingProjectTask) getToolingTaskInfo() *ToolingTaskInfo {
	toolingTaskInfo := ToolingTaskInfo{}
	json.Unmarshal(v.ObjectInfo, &toolingTaskInfo)
	return &toolingTaskInfo
}

type ToolingProjectTaskStatus struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type TaskStatusInfo struct {
	Status        string `json:"status"`
	ColorCode     string `json:"colorCode"`
	CreatedAt     string `json:"createdAt"`
	CreatedBy     int    `json:"createdBy"`
	Preference    int    `json:"preference"`
	Description   string `json:"description"`
	ObjectStatus  string `json:"objectStatus"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

func (v *ToolingProjectTaskStatus) getToolingProjectTaskStatusInfo() *TaskStatusInfo {
	taskStatusInfo := TaskStatusInfo{}
	json.Unmarshal(v.ObjectInfo, &taskStatusInfo)
	return &taskStatusInfo
}

type ToolingProjectSprint struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (v *ToolingProjectSprint) getToolingProjectSprintInfo() *ToolingProjectSprintInfo {
	toolingProjectSprintInfo := ToolingProjectSprintInfo{}
	json.Unmarshal(v.ObjectInfo, &toolingProjectSprintInfo)
	return &toolingProjectSprintInfo
}

type ToolingRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type ToolingEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}

type ToolingProjectSprintInfo struct {
	ProjectId            int     `json:"projectId"`
	SprintName           string  `json:"sprintName"`
	Description          string  `json:"description"`
	ObjectStatus         string  `json:"objectStatus"`
	StartDate            string  `json:"startDate"`
	Status               int     `json:"status"`
	EndDate              string  `json:"endDate"`
	ActualAmount         float64 `json:"actualAmount"`
	ForcastAmount        float64 `json:"forcastAmount"`
	ListOfAssignedTasks  []int   `json:"listOfAssignedTasks"`
	ListOfCompletedTasks []int   `json:"listOfCompletedTasks"`
	LastUpdatedAt        string  `json:"lastUpdatedAt"`
	CreatedAt            string  `json:"createdAt"`
	CanActivate          bool    `json:"canActivate"`
	CanComplete          bool    `json:"canComplete"`
}

func (v *ToolingProjectSprintInfo) DatabaseSerialize() map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (v *ToolingProjectSprintInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}
