package model

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type AssemblyMachineHelpSignalViewInfo struct {
	CreatedAt             string `json:"createdAt"`
	LastUpdatedAt         string `json:"lastUpdatedAt"`
	HelpButtonPressedTime int    `json:"helpButtonPressedTime"`
}

func GetAssemblyMachineHelpSignalViewInfo(serialisedData datatypes.JSON) *AssemblyMachineHelpSignalViewInfo {
	assemblyMachineHelpSignalViewInfo := AssemblyMachineHelpSignalViewInfo{}
	err := json.Unmarshal(serialisedData, &assemblyMachineHelpSignalViewInfo)
	if err != nil {
		return nil
	}
	return &assemblyMachineHelpSignalViewInfo
}

func (v *AssemblyMachineHelpSignalViewInfo) Serialised() ([]byte, error) {
	return json.Marshal(v)
}
