package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type MachineMasterInfo struct {
	Site                         string      `json:"site"`
	Brand                        int         `json:"brand"`
	Model                        string      `json:"model"`
	Plant                        string      `json:"plant"`
	Tonnage                      string      `json:"tonnage"`
	Category                     int         `json:"category"`
	Location                     string      `json:"location"`
	Supplier                     string      `json:"supplier"`
	Department                   int         `json:"department"`
	SubCategory                  int         `json:"subCategory"`
	MachineImage                 string      `json:"machineImage"`
	NewMachineId                 string      `json:"newMachineId"`
	OldMachineId                 string      `json:"oldMachineId"`
	SerialNumber                 string      `json:"serialNumber"`
	MachineStatus                int         `json:"machineStatus"`
	CommissionedDate             interface{} `json:"commissionedDate"`
	FixedAssetNumber             string      `json:"fixedAssetNumber"`
	MachineDescription           string      `json:"machineDescription"`
	MachineConnectStatus         int         `json:"machineConnectStatus"`
	ObjectStatus                 string      `json:"objectStatus"`
	CreatedBy                    int         `json:"createdBy"`
	CreatedAt                    string      `json:"createdAt"`
	LastUpdatedBy                int         `json:"lastUpdatedBy"`
	LastUpdatedAt                string      `json:"lastUpdatedAt"`
	LastUpdatedMachineLiveStatus string      `json:"lastUpdatedMachineLiveStatus"`
	CanCreateWorkOrder           bool        `json:"canCreateWorkOrder"`
	CanCreateCorrectiveWorkOrder bool        `json:"canCreateCorrectiveWorkOrder"`
	DelayStatus                  string      `json:"delayStatus"`
	DelayPeriod                  int64       `json:"delayPeriod"`
	CurrentCycleCount            int         `json:"currentCycleCount"`
	EnableProductionOrder        bool        `json:"enableProductionOrder"`
}

func GetMachineMasterInfo(objectInfo datatypes.JSON) *MachineMasterInfo {
	machineMasterInfo := MachineMasterInfo{}
	err := json.Unmarshal(objectInfo, &machineMasterInfo)
	if err != nil {
		return nil
	}
	return &machineMasterInfo
}
