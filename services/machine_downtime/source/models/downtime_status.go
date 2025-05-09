package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type MachineDowntimeStatusInfo struct {
	Status        string `json:"status"`
	Description   string `json:"description"`
	ColorCode     string `json:"colorCode"`
	CreatedAt     string `json:"createdAt"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	CreatedBy     int    `json:"createdBy"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

func GetMachineDowntimeStatusInfo(objectInfo datatypes.JSON) *MachineDowntimeStatusInfo {
	var machineDowntimeStatusInfo MachineDowntimeStatusInfo
	err := json.Unmarshal(objectInfo, &machineDowntimeStatusInfo)
	if err == nil {
		return &machineDowntimeStatusInfo
	}
	return nil
}
