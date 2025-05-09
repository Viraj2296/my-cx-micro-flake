package database

import (
	"gorm.io/datatypes"
	"time"
)

type MachinesRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type AssemblyMachineMaster struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

/*
CREATE TABLE assembly_machine_view (
    id INT NOT NULL,  -- This will act as the primary key and also the foreign key from assembly_machine_master
    machine_connect_status INT,
    current_cycle_count INT,
    delay_status VARCHAR(255),
    delay_period BIGINT,
    last_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),  -- Primary key on id
    FOREIGN KEY (id) REFERENCES assembly_machine_master(id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO `assembly_machine_view` (`machine_master_id`, `machine_connect_status`, `current_cycle_count`, `delay_status`, `delay_period`)
SELECT `id`, 0, 0, 'none', 0 FROM `machine_master`;
*/
// AssemblyMachineView view of any tables should be lined via foreign key
type AssemblyMachineView struct {
	Id                    int                   `gorm:"column:id; primary_key;not_null" json:"machineConnectStatus"`
	MachineConnectStatus  int                   `gorm:"column:machine_connect_status" json:"machineConnectStatus"`
	CurrentCycleCount     int                   `gorm:"column:current_cycle_count" json:"currentCycleCount"`
	DelayStatus           string                `gorm:"column:delay_status" json:"delayStatus"`
	DelayPeriod           int64                 `gorm:"column:delay_period" json:"delayPeriod"`
	LastUpdatedAt         time.Time             `gorm:"column:last_updated_at;autoUpdateTime" json:"lastUpdatedAt"`
	CreatedAt             time.Time             `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	AssemblyMachineMaster AssemblyMachineMaster `gorm:"foreignKey:Id"`
	StartingCycleCount    int                   `gorm:"column:starting_cycle_count" json:"starting_cycle_count"`
	MessageCycleCount     datatypes.JSON        `gorm:"column:message_cycle_count" json:"message_cycle_count"`
}

type MouldingMachineView struct {
	Id                    int            `gorm:"column:id; primary_key;not_null" json:"machineConnectStatus"`
	MachineConnectStatus  int            `gorm:"column:machine_connect_status" json:"machineConnectStatus"`
	CurrentCycleCount     int            `gorm:"column:current_cycle_count" json:"currentCycleCount"`
	DelayStatus           string         `gorm:"column:delay_status" json:"delayStatus"`
	DelayPeriod           int64          `gorm:"column:delay_period" json:"delayPeriod"`
	LastUpdatedAt         time.Time      `gorm:"column:last_updated_at;autoUpdateTime" json:"lastUpdatedAt"`
	CreatedAt             time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	AssemblyMachineMaster MachineMaster  `gorm:"foreignKey:Id"`
	StartingCycleCount    int            `gorm:"column:starting_cycle_count" json:"starting_cycle_count"`
	MessageCycleCount     datatypes.JSON `gorm:"column:message_cycle_count" json:"message_cycle_count"`
}

type ToolingMachineView struct {
	Id                    int                  `gorm:"column:id; primary_key;not_null" json:"machineConnectStatus"`
	MachineConnectStatus  int                  `gorm:"column:machine_connect_status" json:"machineConnectStatus"`
	CurrentCycleCount     int                  `gorm:"column:current_cycle_count" json:"currentCycleCount"`
	DelayStatus           string               `gorm:"column:delay_status" json:"delayStatus"`
	DelayPeriod           int64                `gorm:"column:delay_period" json:"delayPeriod"`
	LastUpdatedAt         time.Time            `gorm:"column:last_updated_at;autoUpdateTime" json:"lastUpdatedAt"`
	CreatedAt             time.Time            `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	AssemblyMachineMaster ToolingMachineMaster `gorm:"foreignKey:Id"`
	StartingCycleCount    int                  `gorm:"column:starting_cycle_count" json:"starting_cycle_count"`
	MessageCycleCount     datatypes.JSON       `gorm:"column:message_cycle_count" json:"message_cycle_count"`
}

type AssemblyMessageBody struct {
	L01OpId                 string  `json:"L01_opId"`
	L02OpId                 string  `json:"L02_opId"`
	L03OpId                 string  `json:"L03_opId"`
	L04OpId                 string  `json:"L04_opId"`
	L05OpId                 string  `json:"L05_opId"`
	L06OpId                 string  `json:"L06_opId"`
	L01Count                int     `json:"L01_count"`
	L02Count                int     `json:"L02_count"`
	L03Count                int     `json:"L03_count"`
	L04Count                int     `json:"L04_count"`
	L05Count                int     `json:"L05_count"`
	L06Count                int     `json:"L06_count"`
	L01StnCall              *string `json:"L01_stnCall"`
	L02StnCall              *string `json:"L02_stnCall"`
	L03StnCall              *string `json:"L03_stnCall"`
	L04StnCall              *string `json:"L04_stnCall"`
	L05StnCall              *string `json:"L05_stnCall"`
	L06StnCall              *string `json:"L06_stnCall"`
	L01ModelNum             string  `json:"L01_modelNum"`
	L01OrderNum             string  `json:"L01_orderNum"`
	L02ModelNum             string  `json:"L02_modelNum"`
	L02OrderNum             string  `json:"L02_orderNum"`
	L03ModelNum             string  `json:"L03_modelNum"`
	L03OrderNum             string  `json:"L03_orderNum"`
	L04ModelNum             string  `json:"L04_modelNum"`
	L04OrderNum             string  `json:"L04_orderNum"`
	L05ModelNum             string  `json:"L05_modelNum"`
	L05OrderNum             string  `json:"L05_orderNum"`
	L06ModelNum             string  `json:"L06_modelNum"`
	L06OrderNum             string  `json:"L06_orderNum"`
	InTimestamp             int64   `json:"in_timestamp"`
	MachineConnectionStatus int     `json:"machine_connection_status"`
}

type MachineMaster struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingMachineMaster struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
