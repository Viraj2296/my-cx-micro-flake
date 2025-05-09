package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type MachineDowntimeFaultCodeInfo struct {
	Name          string `json:"name"`
	FaultType     int    `json:"faultType"`
	Description   string `json:"description"`
	CreatedAt     string `json:"createdAt"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	CreatedBy     int    `json:"createdBy"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

func GetMachineDowntimeFaultCodeInfo(objectInfo datatypes.JSON) *MachineDowntimeFaultCodeInfo {
	var machineDowntimeFaultCodeInfo MachineDowntimeFaultCodeInfo
	err := json.Unmarshal(objectInfo, &machineDowntimeFaultCodeInfo)
	if err == nil {
		return &machineDowntimeFaultCodeInfo
	}
	return nil
}

func (v *MachineDowntimeFaultCodeInfo) Serialised() ([]byte, error) {
	return json.Marshal(v)
}
