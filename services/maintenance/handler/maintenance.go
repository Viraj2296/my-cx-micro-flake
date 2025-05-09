package handler

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type KanbanTaskStatusRequest struct {
	TaskStatus int `json:"taskStatus"`
}

type MaintenanceRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type MaintenanceOverview struct {
	SnapshotTime datatypes.Time `json:"snapshotTime"`
	//MaintenanceSnapshot
}

type MaintenanceComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenanceWorkOrder struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenanceCorrectiveWorkOrder struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldMaintenanceCorrectiveWorkOrder struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldMaintenanceCorrectiveWorkOrderTask struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MaintenanceWorkOrderCorrectiveTask struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenanceFaultCode struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenancePreventiveWorkOrderStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MouldMaintenancePreventiveWorkOrder struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldCorrectiveMaintenanceJrOption struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *MouldCorrectiveMaintenanceJrOption) getWorkOrderTaskInfo() *WorkOrderInfo {
	workOrderInfo := WorkOrderInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workOrderInfo)
	return &workOrderInfo
}

type WorkOrderInfo struct {
	Name                        string          `json:"name"`
	Labels                      []string        `json:"labels"`
	AssetId                     int             `json:"assetId"`
	CreatedAt                   string          `json:"createdAt"`
	CreatedBy                   int             `json:"createdBy"`
	Attachment                  string          `json:"attachment"`
	Confidence                  int             `json:"confidence"`
	ModuleName                  string          `json:"moduleName"`
	Description                 string          `json:"description"`
	Supervisors                 []int           `json:"supervisors"`
	WorkOrderReferenceId        string          `json:"workOrderReferenceId"`
	ObjectStatus                string          `json:"objectStatus"`
	HasAttachment               bool            `json:"hasAttachment"`
	LastUpdatedAt               string          `json:"lastUpdatedAt"`
	LastUpdatedBy               int             `json:"lastUpdatedBy"`
	RemainderCron               string          `json:"remainderCron"`
	WorkOrderType               string          `json:"workOrderType"`
	WorkOrderStatus             int             `json:"workOrderStatus"`
	DetailedDescription         string          `json:"detailedDescription"`
	WorkOrderRepetitiveCron     string          `json:"workOrderRepetitiveCron"`
	WorkOrderScheduledEndDate   string          `json:"workOrderScheduledEndDate"`
	WorkOrderActualEndDate      string          `json:"workOrderActualEndDate"`
	WorkOrderScheduledStartDate string          `json:"workOrderScheduledStartDate"`
	IsRepetitive                bool            `json:"isRepetitive"`
	IsParent                    bool            `json:"isParent"`
	LastTimeWorkOrderCreation   int64           `json:"lastTimeWorkOrderCreation"`
	CanComplete                 bool            `json:"canComplete"`
	CanUpdate                   bool            `json:"canUpdate"`
	CanRelease                  bool            `json:"canRelease"`
	CanForceStop                bool            `json:"canForceStop"`
	CanUnRelease                bool            `json:"canUnRelease"`
	IsRemainder                 bool            `json:"isRemainder"`
	IsFirstTime                 bool            `json:"isFirstTime"`
	CardImage                   string          `json:"cardImage"`
	CanContinue                 int             `json:"canContinue"`
	RemainderDate               string          `json:"remainderDate"`
	RemainderEndDate            string          `json:"remainderEndDate"`
	RepeatInterval              int             `json:"repeatInterval"`
	RepeatFrequency             int             `json:"repeatFrequency"`
	EmailLastSendDate           string          `json:"emailLastSendDate"`
	ActionRemarks               []ActionRemarks `json:"actionRemarks"`
	MouldId                     int             `json:"mouldId"`
	JrOptionId                  string          `json:"jrOptionId"`
	DownTimeHours               float64         `json:"downTimeHours"`
}

func (v *WorkOrderInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *WorkOrderInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (mwo *MaintenanceWorkOrder) getWorkOrderInfo() *WorkOrderInfo {
	workOrderInfo := WorkOrderInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workOrderInfo)
	return &workOrderInfo
}

func (mwo *MouldMaintenancePreventiveWorkOrder) getWorkOrderInfo() *WorkOrderInfo {
	workOrderInfo := WorkOrderInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workOrderInfo)
	return &workOrderInfo
}

type MaintenanceWorkOrderTask struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *MaintenanceWorkOrderTask) getWorkOrderTaskInfo() *WorkOrderTaskInfo {
	workOrderInfo := WorkOrderTaskInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workOrderInfo)
	return &workOrderInfo
}

type MouldMaintenancePreventiveWorkOrderTask struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *MouldMaintenancePreventiveWorkOrderTask) getWorkOrderTaskInfo() *WorkOrderTaskInfo {
	workOrderInfo := WorkOrderTaskInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workOrderInfo)
	return &workOrderInfo
}

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}

type WorkOrderTaskInfo struct {
	Priority             string          `json:"priority"`
	TaskDate             string          `json:"taskDate"`
	CheckInDate          string          `json:"checkInDate"`
	CheckOutDate         string          `json:"checkOutDate"`
	CanApprove           bool            `json:"canApprove"`
	CanReject            bool            `json:"canReject"`
	CanCheckIn           bool            `json:"canCheckIn"`
	CanCheckOut          bool            `json:"canCheckOut"`
	TaskName             string          `json:"taskName"`
	CreatedAt            string          `json:"createdAt"`
	CreatedBy            int             `json:"createdBy"`
	Attachment           string          `json:"attachment"`
	TaskStatus           int             `json:"taskStatus"`
	WorkOrderId          int             `json:"workOrderId"`
	ObjectStatus         string          `json:"objectStatus"`
	HasAttachment        bool            `json:"hasAttachment"`
	LastUpdatedAt        string          `json:"lastUpdatedAt"`
	LastUpdatedBy        int             `json:"lastUpdatedBy"`
	AssignedUserId       int             `json:"assignedUserId"`
	WorkOrderTaskId      string          `json:"workOrderTaskId"`
	ShortDescription     string          `json:"shortDescription"`
	DetailDescription    string          `json:"detailDescription"`
	EstimatedTaskEndDate string          `json:"estimatedTaskEndDate"`
	IsOrderReleased      bool            `json:"isOrderReleased"`
	CanUpdate            bool            `json:"canUpdate"`
	ActionRemarks        []ActionRemarks `json:"actionRemarks"`
	CanComplete          bool            `json:"canComplete"`
	Remarks              string          `json:"remarks"`
	Remark               string          `json:"remark"`
}

func (mwo *MaintenanceWorkOrderStatus) getMaintenanceWorkOrderStatusInfo() *MaintenanceWorkOrderStatusInfo {
	maintenanceWorkOrderStatusInfo := MaintenanceWorkOrderStatusInfo{}
	json.Unmarshal(mwo.ObjectInfo, &maintenanceWorkOrderStatusInfo)
	return &maintenanceWorkOrderStatusInfo
}

func (v *WorkOrderTaskInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *WorkOrderTaskInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

type MaintenanceWorkOrderStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenanceWorkOrderStatusInfo struct {
	Status                     string `json:"status"`
	ColorCode                  string `json:"colorCode"`
	CreatedAt                  int    `json:"createdAt"`
	Description                string `json:"description"`
	ObjectStatus               string `json:"objectStatus"`
	LastUpdatedAt              string `json:"lastUpdatedAt"`
	LastUpdatedBy              int    `json:"lastUpdatedBy"`
	EmailTemplateId            int    `json:"emailTemplateId"`
	NotificationUserList       []int  `json:"notificationUserList"`
	IsEmailNotificationEnabled bool   `json:"isEmailNotificationEnabled"`
}

type MaintenanceWorkOrderTaskStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *MaintenanceWorkOrderTaskStatus) getMaintenanceWorkOrderTaskStatusInfo() *MaintenanceWorkOrderTaskStatusInfo {
	maintenanceWorkOrderTaskStatusInfo := MaintenanceWorkOrderTaskStatusInfo{}
	json.Unmarshal(mwo.ObjectInfo, &maintenanceWorkOrderTaskStatusInfo)
	return &maintenanceWorkOrderTaskStatusInfo
}

type MaintenanceWorkOrderTaskStatusInfo struct {
	Status        string `json:"status"`
	ColorCode     string `json:"colorCode"`
	CreatedAt     int    `json:"createdAt"`
	Description   string `json:"description"`
	ObjectStatus  string `json:"objectStatus"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

type MaintenanceEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenanceEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaintenanceStatusOverview struct {
	GroupByField string      `json:"groupByField"`
	Cards        interface{} `json:"cards"`
}

type MachineMaintenanceSetting struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MachineMaintenanceSettingInfo struct {
	CreatedAt                       time.Time `json:"createdAt"`
	CreatedBy                       int       `json:"createdBy"`
	ObjectStatus                    string    `json:"objectStatus"`
	LastUpdatedAt                   time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy                   int       `json:"lastUpdatedBy"`
	AllowedPreventiveWorkOrderGroup []int     `json:"allowedPreventiveWorkOrderGroup"`
	AllowedCorrectiveWorkOrderGroup []int     `json:"allowedCorrectiveWorkOrderGroup"`
}

func (v *MachineMaintenanceSettingInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetMachineMaintenanceSettingInfo(serialisedData datatypes.JSON) *MachineMaintenanceSettingInfo {
	machineMaintenanceSettingInfo := MachineMaintenanceSettingInfo{}
	err := json.Unmarshal(serialisedData, &machineMaintenanceSettingInfo)
	if err != nil {
		return &MachineMaintenanceSettingInfo{}
	}
	return &machineMaintenanceSettingInfo
}

type MouldMaintenanceSetting struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MouldMaintenanceSettingInfo struct {
	CreatedAt                       time.Time `json:"createdAt"`
	CreatedBy                       int       `json:"createdBy"`
	ObjectStatus                    string    `json:"objectStatus"`
	LastUpdatedAt                   time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy                   int       `json:"lastUpdatedBy"`
	AllowedPreventiveWorkOrderGroup []int     `json:"allowedPreventiveWorkOrderGroup"`
	AllowedCorrectiveWorkOrderGroup []int     `json:"allowedCorrectiveWorkOrderGroup"`
}

func (v *MouldMaintenanceSettingInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetMouldMaintenanceSettingInfo(serialisedData datatypes.JSON) *MouldMaintenanceSettingInfo {
	mouldMaintenanceSettingInfo := MouldMaintenanceSettingInfo{}
	err := json.Unmarshal(serialisedData, &mouldMaintenanceSettingInfo)
	if err != nil {
		return &MouldMaintenanceSettingInfo{}
	}
	return &mouldMaintenanceSettingInfo
}
