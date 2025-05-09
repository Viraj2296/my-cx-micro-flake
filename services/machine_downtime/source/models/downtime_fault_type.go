package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type MachineDowntimeFaultTypeInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	CreatedAt     string `json:"createdAt"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	CreatedBy     int    `json:"createdBy"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

func GetMachineDowntimeFaultTypeInfo(objectInfo datatypes.JSON) *MachineDowntimeFaultTypeInfo {
	var machineDowntimeFaultTypeInfo MachineDowntimeFaultTypeInfo
	err := json.Unmarshal(objectInfo, &machineDowntimeFaultTypeInfo)
	if err == nil {
		return &machineDowntimeFaultTypeInfo
	}
	return nil
}

func (v *MachineDowntimeFaultTypeInfo) Serialised() ([]byte, error) {
	return json.Marshal(v)
}
