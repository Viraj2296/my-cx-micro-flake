package handler

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/datatypes"
)

type ProductionOrderOverview struct {
	SnapshotTime datatypes.Time `json:"snapshotTime"`
}
type AssemblyManualOrderCompletedQuantityHistory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type ProductionOrderRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type ProductionOrderMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyProductionOrder struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (v *ProductionOrderMaster) getProductionOrderInfo() *ProductionOrderInfo {
	productionOrderInfo := ProductionOrderInfo{}
	json.Unmarshal(v.ObjectInfo, &productionOrderInfo)
	return &productionOrderInfo
}

type ScheduledOrderEvent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyScheduledOrderEvent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MaterialMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (v *ScheduledOrderEvent) getScheduledOrderEventInfo() *ScheduledOrderEventInfo {
	scheduledOrderEventInfo := ScheduledOrderEventInfo{}
	json.Unmarshal(v.ObjectInfo, &scheduledOrderEventInfo)
	return &scheduledOrderEventInfo
}

type ProductionOrderComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ProductionOrderStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ProductionOderStatusInfo struct {
	Status       string `json:"status"`
	ColorCode    string `json:"colorCode"`
	Preference   int    `json:"preference"`
	Description  string `json:"description"`
	ObjectStatus string `json:"objectStatus"`
}

func (pos *ProductionOrderStatus) getProductionOrderStatusInfo() *ProductionOrderStatusInfo {
	productionOrderStatusInfo := ProductionOrderStatusInfo{}
	json.Unmarshal(pos.ObjectInfo, &productionOrderStatusInfo)
	return &productionOrderStatusInfo
}

type ProductionOrderStatusInfo struct {
	Status       string `json:"status"`
	ColorCode    string `json:"colorCode"`
	Preference   int    `json:"preference"`
	Description  string `json:"description"`
	ObjectStatus string `json:"objectStatus"`
}

type ProductionOrderInfo struct {
	OrderStatus                      int         `json:"orderStatus"`
	Balance                          int         `json:"balance"`
	MouldId                          int         `json:"mouldId"`
	ProdQty                          int         `json:"prodQty"`
	Remarks                          string      `json:"remarks"`
	BaseUnit                         string      `json:"baseUnit"`
	OrderQty                         int         `json:"orderQty"`
	PartName                         string      `json:"partName"`
	CycleTime                        float64     `json:"cycleTime"`
	DailyRate                        int         `json:"dailyRate"`
	MachineId                        int         `json:"machineId"`
	ProdOrder                        string      `json:"prodOrder"`
	PartNumber                       int         `json:"partNumber"`
	QtyNeeded1                       int         `json:"qtyNeeded1"`
	QtyNeeded2                       int         `json:"qtyNeeded2"`
	WorkCenter                       string      `json:"workCenter"`
	MaterialUsed1                    interface{} `json:"materialUsed1"`
	MaterialUsed2                    string      `json:"materialUsed2"`
	RecycleMaterial                  string      `json:"recycleMaterial"`
	BomTextPercentageOfMaterialUsage float64     `json:"bomTextPercentageOfMaterialUsage"`
	RemainingScheduledQty            int         `json:"remainingScheduledQty"`
	ObjectStatus                     string      `json:"objectStatus"`
	CreatedAt                        string      `json:"createdAt"`
	CreatedBy                        int         `json:"createdBy"`
	LastUpdatedAt                    string      `json:"lastUpdatedAt"`
	LastUpdatedBy                    int         `json:"lastUpdatedBy"`
}

type AssemblyProductionOrderInfo struct {
	OrderStatus                      int     `json:"orderStatus"`
	Balance                          int     `json:"balance"`
	MouldId                          int     `json:"mouldId"`
	ProdQty                          int     `json:"prodQty"`
	Remarks                          string  `json:"remarks"`
	OrderQty                         int     `json:"orderQty"`
	PartName                         string  `json:"partName"`
	CycleTime                        float64 `json:"cycleTime"`
	DailyRate                        int     `json:"dailyRate"`
	MachineId                        int     `json:"machineId"`
	ProdOrder                        string  `json:"prodOrder"`
	PartNumber                       int     `json:"partNumber"`
	WorkCenter                       string  `json:"workCenter"`
	RecycleMaterial                  string  `json:"recycleMaterial"`
	BomTextPercentageOfMaterialUsage float64 `json:"bomTextPercentageOfMaterialUsage"`
	RemainingScheduledQty            int     `json:"remainingScheduledQty"`
	ObjectStatus                     string  `json:"objectStatus"`
	CreatedAt                        string  `json:"createdAt"`
	CreatedBy                        int     `json:"createdBy"`
	LastUpdatedAt                    string  `json:"lastUpdatedAt"`
	LastUpdatedBy                    int     `json:"lastUpdatedBy"`
}

type ScheduledOrderEventInfo struct {
	Name                 string  `json:"name"`
	EventSourceId        int     `json:"eventSourceId"`
	EndDate              string  `json:"endDate"`
	IconCls              string  `json:"iconCls"`
	Draggable            bool    `json:"draggable"`
	EventType            string  `json:"eventType"`
	EventStatus          int     `json:"eventStatus"`
	StartDate            string  `json:"startDate"`
	EventColor           string  `json:"eventColor"`
	ProductionOrder      string  `json:"productionOrder"`
	MachineId            int     `json:"machineId"`
	ScheduledQty         int     `json:"scheduledQty"`
	CompletedQty         int     `json:"completedQty"`
	PercentDone          float64 `json:"percentDone"`
	RejectedQty          int     `json:"rejectedQty"`
	ObjectStatus         string  `json:"objectStatus"`
	CreatedBy            int     `json:"createdBy"`
	CreatedAt            string  `json:"createdAt"`
	LastUpdatedAt        string  `json:"lastUpdatedAt"`
	LastUpdatedBy        int     `json:"lastUpdatedBy"`
	Oee                  int     `json:"oee"`
	EnableCustomCavity   bool    `json:"enableCustomCavity"`
	MouldId              int     `json:"mouldId"`
	CustomCavity         int     `json:"customCavity"`
	MouldUp              string  `json:"mouldUp"`
	MouldDown            string  `json:"mouldDown"`
	IsAbortEnabled       bool    `json:"isAbortEnabled"`
	CanComplete          bool    `json:"canComplete"`
	IsUpdate             bool    `json:"isUpdate"`
	CanForceStop         bool    `json:"canForceStop"`
	IsRecoverySchedule   bool    `json:"isRecoverySchedule"`
	RecoveryScheduleId   int     `json:"recoveryScheduleId"`
	MouldBatchResourceId int     `json:"mouldBatchResourceId"`
	PlannedManpower      int     `json:"plannedManpower"` // this was introduced to set the default manpower when creating the schedule
}

type ToolingScheduledOrderEventInfo struct {
	Name            string  `json:"name"`
	EventSourceId   int     `json:"eventSourceId"`
	EndDate         string  `json:"endDate"`
	IconCls         string  `json:"iconCls"`
	Draggable       bool    `json:"draggable"`
	EventType       string  `json:"eventType"`
	EventStatus     int     `json:"eventStatus"`
	StartDate       string  `json:"startDate"`
	EventColor      string  `json:"eventColor"`
	ProductionOrder string  `json:"productionOrder"`
	MachineId       int     `json:"machineId"`
	CompletedQty    string  `json:"completedQty"`
	PercentDone     float64 `json:"percentDone"`
	RejectedQty     int     `json:"rejectedQty"`
	ObjectStatus    string  `json:"objectStatus"`
	CreatedBy       int     `json:"createdBy"`
	CreatedAt       string  `json:"createdAt"`
	LastUpdatedAt   string  `json:"lastUpdatedAt"`
	LastUpdatedBy   int     `json:"lastUpdatedBy"`
	Oee             int     `json:"oee"`
	PartId          int     `json:"partId"`
	IsAbortEnabled  bool    `json:"isAbortEnabled"`
	CanComplete     bool    `json:"canComplete"`
	IsUpdate        bool    `json:"isUpdate"`
	SetupTime       string  `json:"setupTime"`
	CanSetupTime    bool    `json:"canSetupTime"`
	CanForceStop    bool    `json:"canForceStop"`
}

func (v *ToolingScheduledOrderEvent) getToolingScheduledOrderEventInfo() *ToolingScheduledOrderEventInfo {
	scheduledOrderEventInfo := ToolingScheduledOrderEventInfo{}
	json.Unmarshal(v.ObjectInfo, &scheduledOrderEventInfo)
	return &scheduledOrderEventInfo
}

func (v *ScheduledOrderEventInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *ToolingScheduledOrderEventInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *ProductionOrderInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *ProductionOrderInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (v *ScheduledOrderEventInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (v *ToolingScheduledOrderEventInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

type ToolingBomMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingPartMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingOrderMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingScheduledOrderEvent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
