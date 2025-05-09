package model

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type AssemblyHelpSignalProcessedTimeInfo struct {
	ProcessedMessageTs int    `json:"processedMessageTs"`
	LastUpdatedAt      string `json:"lastUpdatedAt"`
	CreatedAt          string `json:"createdAt"`
}

func GetAssemblyHelpSignalProcessedTimeInfo(serialisedData datatypes.JSON) *AssemblyHelpSignalProcessedTimeInfo {
	assemblyMachineHelpSignalViewInfo := AssemblyHelpSignalProcessedTimeInfo{}
	err := json.Unmarshal(serialisedData, &assemblyMachineHelpSignalViewInfo)
	if err != nil {
		return nil
	}
	return &assemblyMachineHelpSignalViewInfo
}

func (v *AssemblyHelpSignalProcessedTimeInfo) Serialised() ([]byte, error) {
	return json.Marshal(v)
}
