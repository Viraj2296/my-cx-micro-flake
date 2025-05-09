package rpc

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type AssemblyManualOrderHistoryByEventRequest struct {
	EventId int `json:"eventId"`
}

type AssemblyManualOrderHistoryByEventResponse struct {
	CreatedAt         string `json:"createdAt"`
	CompletedQuantity int    `json:"completedQuantity"`
}

func GetAssemblyManualOrderHistoryByEventRequest(serialisedData datatypes.JSON) (error, *AssemblyManualOrderHistoryByEventRequest) {
	assemblyManualOrderHistoryByEventRequest := AssemblyManualOrderHistoryByEventRequest{}
	err := json.Unmarshal(serialisedData, &assemblyManualOrderHistoryByEventRequest)
	if err != nil {
		return err, nil
	}
	return nil, &assemblyManualOrderHistoryByEventRequest
}
