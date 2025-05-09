package handler

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

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
	EnableCustomCavity   bool    `json:"enableCustomCavity"`
	MouldId              int     `json:"mouldId"`
	CustomCavity         int     `json:"customCavity"`
	MouldUp              string  `json:"mouldUp"`
	MouldDown            string  `json:"mouldDown"`
	IsAbortEnabled       bool    `json:"isAbortEnabled"`
	CanComplete          bool    `json:"canComplete"`
	IsUpdate             bool    `json:"isUpdate"`
	CanForceStop         bool    `json:"canForceStop"`
	MouldBatchResourceId int     `json:"mouldBatchResourceId"`
}

func GetScheduledOrderEventInfo(eventObject datatypes.JSON) *ScheduledOrderEventInfo {
	scheduledOrderEventInfo := ScheduledOrderEventInfo{}
	json.Unmarshal(eventObject, &scheduledOrderEventInfo)
	return &scheduledOrderEventInfo
}

type ProductionOrderInfo struct {
	Balance                          int            `json:"balance"`
	MouldId                          int            `json:"mouldId"`
	ProdQty                          int            `json:"prodQty"`
	Remarks                          string         `json:"remarks"`
	BaseUnit                         string         `json:"baseUnit"`
	OrderQty                         int            `json:"orderQty"`
	PartName                         string         `json:"partName"`
	CycleTime                        float32        `json:"cycleTime"`
	DailyRate                        int            `json:"dailyRate"`
	MachineId                        int            `json:"machineId"`
	ProdOrder                        string         `json:"prodOrder"`
	PartNumber                       int            `json:"partNumber"`
	QtyNeeded1                       int            `json:"qtyNeeded1"`
	QtyNeeded2                       int            `json:"qtyNeeded2"`
	WorkCenter                       string         `json:"workCenter"`
	OrderStatus                      int            `json:"orderStatus"`
	MaterialUsed1                    int            `json:"materialUsed1"`
	MaterialUsed2                    string         `json:"materialUsed2"`
	RecycleMaterial                  string         `json:"recycleMaterial"`
	RemainingScheduledQty            int            `json:"remainingScheduledQty"`
	BomTextPercentageOfMaterialUsage int            `json:"bomTextPercentageOfMaterialUsage"`
	NoOfCustomMouldCavity            int            `json:"noOfCustomMouldCavity"`
	ObjectInfo                       datatypes.JSON `json:"objectInfo"`
}

func GetProductionOrderInfo(eventObject datatypes.JSON) *ProductionOrderInfo {
	productionOrderInfo := ProductionOrderInfo{}
	json.Unmarshal(eventObject, &productionOrderInfo)
	return &productionOrderInfo
}

type PartInfo struct {
	Image          string    `json:"image"`
	PartNumber     string    `json:"partNumber"`
	Description    string    `json:"description"`
	LastUpdatedAt  time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy  int       `json:"lastUpdatedBy"`
	IsBatchManaged bool      `json:"isBatchManaged"`
}

func GetPartInfo(eventObject datatypes.JSON) *PartInfo {
	partInfo := PartInfo{}
	json.Unmarshal(eventObject, &partInfo)
	return &partInfo
}

type ProductionOrderStatusInfo struct {
	Status       string `json:"status"`
	ColorCode    string `json:"colorCode"`
	Preference   int    `json:"preference"`
	Description  string `json:"description"`
	ObjectStatus string `json:"objectStatus"`
}

func GetProductionOrderStatusInfo(eventObject datatypes.JSON) *ProductionOrderStatusInfo {
	productionOrderStatusInfo := ProductionOrderStatusInfo{}
	json.Unmarshal(eventObject, &productionOrderStatusInfo)
	return &productionOrderStatusInfo
}
