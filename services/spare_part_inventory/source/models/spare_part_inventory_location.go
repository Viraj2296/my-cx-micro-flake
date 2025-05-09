package models

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type SparePartInventoryLocationInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetSparePartInventoryLocationInfo(serialisedData datatypes.JSON) *SparePartInventoryLocationInfo {
	sparePartInventoryLocationInfo := SparePartInventoryLocationInfo{}
	err := json.Unmarshal(serialisedData, &sparePartInventoryLocationInfo)
	if err != nil {

		return &SparePartInventoryLocationInfo{}
	}

	return &sparePartInventoryLocationInfo
}
