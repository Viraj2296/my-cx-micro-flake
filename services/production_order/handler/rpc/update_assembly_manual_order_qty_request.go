package rpc

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type UpdateAssemblyManualOrderQuantityRequest struct {
	EventId           int `json:"eventId"`
	CompletedQuantity int `json:"completedQuantity"`
	RequestBy         int `json:"requestBy"`
}

func GetUpdateAssemblyManualOrderQuantityRequest(serialisedData datatypes.JSON) (error, *UpdateAssemblyManualOrderQuantityRequest) {
	updateAssemblyManualOrderQuantityRequest := UpdateAssemblyManualOrderQuantityRequest{}
	err := json.Unmarshal(serialisedData, &updateAssemblyManualOrderQuantityRequest)
	if err != nil {
		return err, nil
	}
	return nil, &updateAssemblyManualOrderQuantityRequest
}
