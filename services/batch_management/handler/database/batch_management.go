package database

import (
	"encoding/json"
	"gorm.io/datatypes"
	"time"
)

type MouldOverview struct {
	SnapshotTime datatypes.Time `json:"snapshotTime"`
}

type BatchManagementRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type BatchManagementComponent struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type BatchManagementRawMaterial struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type BatchManagementMaterialType struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type BatchManagementMould struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type BatchManagementPrinter struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type RawMaterialBatchInfo struct {
	Name              string      `json:"name"`
	Label             string      `json:"label"`
	BatchId           string      `json:"batchId"`
	CanPrint          bool        `json:"canPrint"`
	Location          int         `json:"location"`
	VendorId          int         `json:"vendorId"`
	CreatedAt         time.Time   `json:"createdAt"`
	CreatedBy         int         `json:"createdBy"`
	LabelImage        string      `json:"labelImage"`
	Description       string      `json:"description"`
	LabelStatus       int         `json:"labelStatus"`
	VendorLotNo       string      `json:"vendorLotNo"`
	MaterialType      int         `json:"materialType"`
	ObjectStatus      string      `json:"objectStatus"`
	LastUpdatedAt     time.Time   `json:"lastUpdatedAt"`
	LastUpdatedBy     int         `json:"lastUpdatedBy"`
	MaterialBatchNo   string      `json:"materialBatchNo"`
	GoodReceiveNumber string      `json:"goodReceiveNumber"`
	LabelPrintMessage interface{} `json:"labelPrintMessage"`
}

func GeRawMaterialBatchInfo(serialisedData datatypes.JSON) *RawMaterialBatchInfo {
	rawMaterialBatchInfo := RawMaterialBatchInfo{}
	err := json.Unmarshal(serialisedData, &rawMaterialBatchInfo)
	if err != nil {
		return &RawMaterialBatchInfo{}
	}
	return &rawMaterialBatchInfo
}

func (v *RawMaterialBatchInfo) Serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON{}
	}
	return serialisedData

}

type MouldBatchInfo struct {
	Label             string      `json:"label"`
	MouldId           int         `json:"mouldId"`
	CanPrint          bool        `json:"canPrint"`
	Location          int         `json:"location"`
	StopTime          time.Time   `json:"stopTime"`
	CreatedAt         time.Time   `json:"createdAt"`
	CreatedBy         int         `json:"createdBy"`
	MachineId         int         `json:"machineId"`
	LabelImage        string      `json:"labelImage"`
	NoOfLabels        interface{} `json:"noOfLabels"`
	OperatorId        int         `json:"operatorId"`
	LabelStatus       int         `json:"labelStatus"`
	MouldBatchId      string      `json:"mouldBatchId"`
	ObjectStatus      string      `json:"objectStatus"`
	LastUpdatedAt     time.Time   `json:"lastUpdatedAt"`
	LastUpdatedBy     int         `json:"lastUpdatedBy"`
	RawMaterialId     int         `json:"rawMaterialId"`
	ScheduleEventId   int         `json:"scheduleEventId"`
	MouldBatchNumber  string      `json:"mouldBatchNumber"`
	LabelPrintMessage interface{} `json:"labelPrintMessage"`
}

func GeRawMouldBatchInfo(serialisedData datatypes.JSON) *MouldBatchInfo {
	mouldBatchInfo := MouldBatchInfo{}
	err := json.Unmarshal(serialisedData, &mouldBatchInfo)
	if err != nil {
		return &MouldBatchInfo{}
	}
	return &mouldBatchInfo
}

func (v *MouldBatchInfo) Serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON{}
	}
	return serialisedData

}

type PrinterInfo struct {
	Name             string `json:"name"`
	Message          string `json:"message"`
	Location         int    `json:"location"`
	CreatedAt        string `json:"createdAt"`
	CreatedBy        int    `json:"createdBy"`
	NetworkIP        string `json:"networkIP"`
	Description      string `json:"description"`
	NetworkPort      int    `json:"networkPort"`
	ObjectStatus     string `json:"objectStatus"`
	LastUpdatedAt    string `json:"lastUpdatedAt"`
	LastUpdatedBy    int    `json:"lastUpdatedBy"`
	ConnectionStatus string `json:"connectionStatus"`
	CoordinateX      int    `json:"coordinateX"`
	CoordinateY      int    `json:"coordinateY"`
	ModuleSize       int    `json:"moduleSize"`
	Magnification    int    `json:"magnification"`
}

func GetPrinterInfo(serialisedData datatypes.JSON) *PrinterInfo {
	printerInfo := PrinterInfo{}
	err := json.Unmarshal(serialisedData, &printerInfo)
	if err != nil {
		return &PrinterInfo{}
	}
	return &printerInfo
}

func (v *PrinterInfo) Serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON{}
	}
	return serialisedData

}
