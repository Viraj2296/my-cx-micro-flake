package models

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type SparePartInventoryTransactionInfo struct {
	PartNumber                string `json:"part_number"`
	PartDescription           string `json:"part_description"`
	DestinationLocation       int    `json:"destination_location"`
	Transaction               string `json:"transaction"`
	Qty                       int    `json:"qty"`
	SourceLocation            int    `json:"source_location"`
	ServiceNotificationNumber string `json:"service_notification_number"`
}

func GetSparePartInventoryTransactionInfo(serialisedData datatypes.JSON) *SparePartInventoryTransactionInfo {
	sparePartInventoryTransactionInfo := SparePartInventoryTransactionInfo{}
	err := json.Unmarshal(serialisedData, &sparePartInventoryTransactionInfo)
	if err != nil {

		return &SparePartInventoryTransactionInfo{}
	}

	return &sparePartInventoryTransactionInfo
}
func (v *SparePartInventoryTransactionInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}
