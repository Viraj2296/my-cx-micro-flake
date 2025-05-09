package handler

import (
	"cx-micro-flake/pkg/common/component"
	"encoding/json"

	"gorm.io/datatypes"
)

type Message struct {
	Id    int `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	Topic string
	Body  datatypes.JSON
	Ts    int64
}
type IOTRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}
type IOTComponent struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTNode struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTDriver struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

// IOTDriverStatistics store the driver stats information, record id is the driver id
type IOTDriverStatistics struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	RecordId   int            `json:"recordId" gorm:"primary_key;not_null"`
	Timestamp  int64          `json:"timestamp"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTStaticsMessage struct {
	Topic string         `json:"topic"`
	Id    int            `json:"id"`
	Ts    int64          `json:"ts"`
	Body  datatypes.JSON `json:"body"`
}

type IOTDataSourceMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTDriverDynamicField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IOTDriverDynamicFieldConfiguration struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (v *IOTDriverDynamicFieldConfiguration) getIoTDriverDynamicFieldInfo() *IoTDriverDynamicFieldInfo {
	iotDriverDynamicFieldInfo := IoTDriverDynamicFieldInfo{}
	json.Unmarshal(v.ObjectInfo, &iotDriverDynamicFieldInfo)
	return &iotDriverDynamicFieldInfo
}

type DynamicFields struct {
	Property   string `json:"property"`
	Type       string `json:"type"`
	GridSystem string `json:"gridSystem"`
	Label      string `json:"label"`
}

type IoTDriverDynamicFieldInfo struct {
	ComponentName string                   `json:"componentName"`
	Id            int                      `json:"id"`
	DynamicFields []component.RecordSchema `json:"dynamicFields"`
}

type DatasourceMasterInfo struct {
	Type               string `json:"type"`
	Category           string `json:"category"`
	ObjectStatus       string `json:"objectStatus"`
	DatasourceImageUrl string `json:"datasourceImageUrl"`
	ConnectionField    string `json:"connectionField"`
	LongDescription    string `json:"longDescription"`
	ShortDescription   string `json:"shortDescription"`
	DisplayName        string `json:"displayName"`
}
type DriverInfo struct {
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	Image                  string `json:"image"`
	Status                 string `json:"status"`
	Protocol               string `json:"protocol"`
	Description            string `json:"description"`
	ConnectionParam        string `json:"connectionParam"`
	EmailTemplateId        int    `json:"emailTemplateId"`
	EmailTriggerTime       string `json:"emailTriggerTime"`
	NotificationGroup      []int  `json:"notificationGroup"`
	TcpSourceConfiguration struct {
		Ip                 string `json:"ip"`
		Port               int    `json:"port"`
		Decoder            string `json:"decoder"`
		AuthEnabled        bool   `json:"authEnabled"`
		CompressionEnabled bool   `json:"compressionEnabled"`
	} `json:"tcpSourceConfiguration"`
}

func (v *IOTDataSourceMaster) getDatasourceMasterInfo() *DatasourceMasterInfo {
	datasourceMasterInfo := DatasourceMasterInfo{}
	json.Unmarshal(v.ObjectInfo, &datasourceMasterInfo)
	return &datasourceMasterInfo
}

type MachineData struct {
	CycleCount              int   `json:"cycle_count"`
	Auto                    int   `json:"auto"`
	MaintenanceDoorAlarm    int   `json:"maintenance_door_alarm"`
	SafetyDoorAlarm         int   `json:"safety_door_alarm"`
	EstopAlarm              int   `json:"estop_alarm"`
	TowerLightRed           int   `json:"tower_light_red"`
	MachineConnectionStatus int   `json:"machine_connection_status"`
	InTimestamp             int64 `json:"in_timestamp"`
}
