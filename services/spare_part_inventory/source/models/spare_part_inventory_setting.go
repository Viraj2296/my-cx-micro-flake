package models

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type SparePartInventorySettingInfo struct {
	IsEnableEscalation               bool  `json:"isEnableEscalation"`
	InventoryStockLimitAlertUsers    []int `json:"inventoryStockLimitAlertUsers"`
	EscalationWaitingPeriod          int   `json:"escalationWaitingPeriod"`
	InventoryThresholdQuantity       int   `json:"inventoryThresholdQuantity"`
	InitialWaitingPeriod             int   `json:"initialWaitingPeriod"`
	InventoryStockAlertEmailTemplate int   `json:"inventoryStockAlertEmailTemplate"` // based on selected ID, we can load the template and send the alert email
}

func GetSparePartInventorySettingInfo(serialisedData datatypes.JSON) *SparePartInventorySettingInfo {
	sparePartInventorySettingInfo := SparePartInventorySettingInfo{}
	err := json.Unmarshal(serialisedData, &sparePartInventorySettingInfo)
	if err != nil {

		return &SparePartInventorySettingInfo{}
	}

	return &sparePartInventorySettingInfo
}
