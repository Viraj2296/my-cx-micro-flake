package database

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type LabourManagementRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type LabourManagementComponent struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type LabourManagementShiftMaster struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type LabourManagementShiftProduction struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type LabourManagementAttendance struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type LabourManagementShiftStatus struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type LabourManagementShiftTemplate struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type LabourManagementShiftProductionAttendance struct {
	ShiftProductionId int       `json:"shift_production_id" gorm:"primary_key;not_null"`
	ShiftAttendanceId int       `json:"shift_attendance_id" gorm:"primary_key;not_null"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         int       `json:"created_by"`
}

// ShiftTemplateInfo The purpose of the shift template is keep the different template for each sites
type ShiftTemplateInfo struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	ShiftStartTime string `json:"shiftStartTime"`
	ShiftPeriod    int    `json:"shiftPeriod"`
	LastUpdatedAt  string `json:"lastUpdatedAt"`
	LastUpdatedBy  int    `json:"lastUpdatedBy"`
	CreatedAt      string `json:"createdAt"`
	CreatedBy      int    `json:"createdBy"`
	ObjectStatus   string `json:"objectStatus"`
	SiteId         int    `json:"siteId"`
}

func GetSShiftTemplateInfo(serialisedData datatypes.JSON) *ShiftTemplateInfo {
	shiftTemplateInfo := ShiftTemplateInfo{}
	err := json.Unmarshal(serialisedData, &shiftTemplateInfo)
	if err != nil {
		return &ShiftTemplateInfo{}
	}
	return &shiftTemplateInfo
}

type ShiftMasterProductionInfo struct {
	ShiftId                    int             `json:"shiftId"`
	ScheduledEventId           int             `json:"scheduledEventId"`
	ActualManpower             int             `json:"actualManpower"`         // this is based on scanning using mobile app
	ActualManHourPartTimer     int             `json:"actualManHourPartTimer"` // this is the calculated valued based on scanning, we can use the job role type contract or not
	ShiftTargetOutputPartTimer int             `json:"shiftTargetOutputPartTimer"`
	ShiftActualOutputPartTimer int             `json:"shiftActualOutputPartTimer"`
	ShiftActualOutput          int             `json:"shiftActualOutput"`
	Remarks                    string          `json:"remarks"`
	LastUpdatedAt              string          `json:"lastUpdatedAt"`
	LastUpdatedBy              int             `json:"lastUpdatedBy"`
	CreatedAt                  string          `json:"createdAt"`
	CreatedBy                  int             `json:"createdBy"`
	MachineId                  int             `json:"machineId"`
	ActionRemarks              []ActionRemarks `json:"actionRemarks"`
	CanUpdateManpower          bool            `json:"canUpdateManpower"`
}

func GetShiftMasterProductionInfo(serialisedData datatypes.JSON) *ShiftMasterProductionInfo {
	shiftMasterProductionInfo := ShiftMasterProductionInfo{}
	err := json.Unmarshal(serialisedData, &shiftMasterProductionInfo)
	if err != nil {
		return &ShiftMasterProductionInfo{}
	}
	return &shiftMasterProductionInfo
}

func (v *ShiftMasterProductionInfo) Serialised() datatypes.JSON {
	serialised, _ := json.Marshal(v)
	return serialised
}

type ShiftMasterInfo struct {
	SiteId                int             `json:"siteId"`
	ShiftReferenceId      string          `json:"shiftReferenceId"`
	CreatedAt             string          `json:"createdAt"`
	CreatedBy             int             `json:"createdBy"`
	CanCheckIn            bool            `json:"canCheckIn"`
	ShiftStatus           int             `json:"shiftStatus"`
	CanShiftStop          bool            `json:"canShiftStop"`
	DepartmentId          int             `json:"departmentId"`
	ObjectStatus          string          `json:"objectStatus"`
	ShiftEndDate          string          `json:"shiftEndDate"`
	ShiftEndTime          string          `json:"shiftEndTime"`
	CanShiftStart         bool            `json:"canShiftStart"`
	LastUpdatedAt         string          `json:"lastUpdatedAt"`
	LastUpdatedBy         int             `json:"lastUpdatedBy"`
	ActualManPower        int             `json:"actualManPower"`
	ShiftStartDate        string          `json:"shiftStartDate"`
	ShiftStartTime        string          `json:"shiftStartTime"`
	ShiftSupervisor       int             `json:"shiftSupervisor"`
	ScheduledOrderEvents  []int           `json:"scheduledOrderEvents"` // when creating the shift, users can select the list of lines this shift assigned it.
	IsSupervisorCheckedIn bool            `json:"isSupervisorCheckedIn"`
	ShiftTemplateId       int             `json:"shiftTemplateId"`
	ActionRemarks         []ActionRemarks `json:"actionRemarks"`
	CanRollBack           bool            `json:"canRollBack"`
}

func (v *ShiftMasterInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetShiftMasterInfo(serialisedData datatypes.JSON) *ShiftMasterInfo {
	shiftMaster := ShiftMasterInfo{}
	err := json.Unmarshal(serialisedData, &shiftMaster)
	if err != nil {
		return &ShiftMasterInfo{}
	}
	return &shiftMaster
}

type AttendanceInfo struct {
	ShiftResourceId    int    `json:"shiftResourceId"`
	CheckInDate        string `json:"checkInDate"`
	CheckInTime        string `json:"checkInTime"`
	UserResourceId     int    `json:"userResourceId"`
	CheckOutDate       string `json:"checkOutDate"`
	CheckOutTime       string `json:"checkOutTime"`
	CreatedAt          string `json:"createdAt"`
	LastUpdatedAt      string `json:"lastUpdatedAt"`
	LastUpdatedBy      int    `json:"lastUpdatedBy"`
	ManufacturingLines []int  `json:"manufacturingLines"` // we may have multiple lines this means employees responsible for multiple lines, and line supervisors or leads
	CreatedBy          int    `json:"createdBy"`
}

func GetAttendanceInfo(serialisedData datatypes.JSON) *AttendanceInfo {
	attendanceInfo := AttendanceInfo{}
	err := json.Unmarshal(serialisedData, &attendanceInfo)
	if err != nil {
		return &AttendanceInfo{}
	}
	return &attendanceInfo
}

func (v *AttendanceInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

type ShiftStatusInfo struct {
	Status        string    `json:"status"`
	ColorCode     string    `json:"colorCode"`
	CreatedAt     time.Time `json:"createdAt"`
	Description   string    `json:"description"`
	ObjectStatus  string    `json:"objectStatus"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func (v *ShiftStatusInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetShiftStatusInfo(serialisedData datatypes.JSON) *ShiftStatusInfo {
	shiftStatusInfo := ShiftStatusInfo{}
	err := json.Unmarshal(serialisedData, &shiftStatusInfo)
	if err != nil {
		return &ShiftStatusInfo{}
	}
	return &shiftStatusInfo
}

type AssemblyMachineLineInfo struct {
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	ObjectStatus  string    `json:"objectStatus"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func GetAssemblyMachineLineInfo(serialisedData datatypes.JSON) *AssemblyMachineLineInfo {
	assemblyMachineLineInfo := AssemblyMachineLineInfo{}
	err := json.Unmarshal(serialisedData, &assemblyMachineLineInfo)
	if err != nil {
		return &AssemblyMachineLineInfo{}
	}
	return &assemblyMachineLineInfo
}

type LabourManagementSettingInfo struct {
	LabourManagementSummaryJobRoles    []int  `json:"labourManagementSummaryJobRoles"`
	ShiftSupervisorJobRoleId           int    `json:"shiftSupervisorJobRoleId"`
	ShiftOperatorJobRoleId             int    `json:"shiftOperatorJobRoleId"`
	ContractEmployeeTypeId             int    `json:"contractEmployeeTypeId"`
	LineDisplayRoles                   []int  `json:"lineDisplayRoles"`
	ShiftCreationRoles                 []int  `json:"shiftCreationRoles"`
	ShiftAutoStopTime                  string `json:"shiftAutoStopTime"`
	ShitAutoStopEmailNotificationUsers []int  `json:"shitAutoStopEmailNotificationUsers"`
}

func GetLabourManagementSettingInfo(serialisedData datatypes.JSON) *LabourManagementSettingInfo {
	labourManagementSettingInfo := LabourManagementSettingInfo{}
	err := json.Unmarshal(serialisedData, &labourManagementSettingInfo)
	if err != nil {
		return &LabourManagementSettingInfo{}
	}
	return &labourManagementSettingInfo
}

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}

type ActualShiftValueCache struct {
	TimeInterval string `json:"timeInterval"`
	ActualValue  int    `json:"actualValue"`
}

type RefreshExportRequest struct {
	SelectedId []string `json:"selectEdId"`
}
