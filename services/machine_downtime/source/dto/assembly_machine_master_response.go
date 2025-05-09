package dto

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type AssemblyMachineMasterInfoResponse struct {
	NewMachineId                 string `json:"newMachineId"`
	MessageFlag                  string `json:"messageFlag"`
	Area                         string `json:"area"`
	Level                        string `json:"level"`
	HelpButtonStationNo          string `json:"helpButtonStationNo"`
	Department                   int    `json:"department"`
	CreatedBy                    int    `json:"createdBy"`
	CreatedAt                    string `json:"createdAt"`
	ObjectStatus                 string `json:"objectStatus"`
	MachineImage                 string `json:"machineImage"`
	Description                  string `json:"description"`
	MachineStatus                int    `json:"machineStatus"`
	MachineConnectStatus         int    `json:"machineConnectStatus"`
	LastUpdatedAt                string `json:"lastUpdatedAt"`
	DelayStatus                  string `json:"delayStatus"`
	DelayPeriod                  int64  `json:"delayPeriod"`
	CurrentCycleCount            int    `json:"currentCycleCount"`
	CanCreateWorkOrder           bool   `json:"canCreateWorkOrder"`
	LastUpdatedMachineLiveStatus string `json:"lastUpdatedMachineLiveStatus"`
	AssemblyLineOption           int    `json:"assemblyLineOption"`
	Model                        string `json:"model"` // This was added as part of labour management module
	IsMESDriverConfigured        bool   `json:"isMESDriverConfigured"`
	IsEnabled                    bool   `json:"isEnabled"`
	AutoRejectHistoryGeneration  bool   `json:"AutoRejectHistoryGeneration"`
	LineMappingMessageFlag       string `json:"lineMappingMessageFlag"` // L01, L02, L03, L04, L05, L06
	EquipmentId                  string `json:"equipmentId"`            // This was added to check-in the machine using the mobile using qr code
	EquipmentName                string `json:"equipmentName"`          // This was added as part of machine downtime module
}

func GetAssemblyMachineMasterInfo(objectInfo datatypes.JSON) (AssemblyMachineMasterInfoResponse, error) {
	var assemblyMachineMasterInfo AssemblyMachineMasterInfoResponse
	err := json.Unmarshal(objectInfo, &assemblyMachineMasterInfo)
	return assemblyMachineMasterInfo, err
}
