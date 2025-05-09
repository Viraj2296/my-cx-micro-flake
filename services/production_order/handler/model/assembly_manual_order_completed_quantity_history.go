package model

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type AssemblyManualOrderCompletedQuantityHistoryInfo struct {
	LastUpdatedAt     string `json:"lastUpdatedAt"`
	LastUpdatedBy     int    `json:"lastUpdatedBy"`
	CreatedAt         string `json:"createdAt"`
	CreatedBy         int    `json:"createdBy"`
	ObjectStatus      string `json:"objectStatus"`
	EventId           int    `json:"eventId"` // eventID is only enough get all as it is linked to everything
	CompletedQuantity int    `json:"completedQuantity"`
}

func GetAssemblyManualOrderCompletedQuantityHistoryInfo(serialisedData datatypes.JSON) (error, *AssemblyManualOrderCompletedQuantityHistoryInfo) {
	assemblyManualOrderCompletedQuantityHistoryInfo := AssemblyManualOrderCompletedQuantityHistoryInfo{}
	err := json.Unmarshal(serialisedData, &assemblyManualOrderCompletedQuantityHistoryInfo)
	if err != nil {
		return err, nil
	}
	return nil, &assemblyManualOrderCompletedQuantityHistoryInfo
}

func (v *AssemblyManualOrderCompletedQuantityHistoryInfo) Serialised() ([]byte, error) {
	return json.Marshal(v)
}
