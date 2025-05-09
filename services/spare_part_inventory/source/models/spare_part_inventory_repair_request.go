package models

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type SparePartList struct {
	Quantity    int `json:"quantity"`
	SparePartId int `json:"sparePartId"`
}

type SparePartInventoryRepairRequestInfo struct {
	IsNeedSparePart bool            `json:"isNeedSparePart"`
	SpareParts      []SparePartList `json:"spareParts"`
	JobId           string          `json:"jobId"`
	RequestStatus   string          `json:"requestStatus"` //CREATED, APPROVED, CANCELLED
	MachineId       int             `json:"machineId"`
}

func GetSparePartInventoryRepairRequestInfo(serialisedData datatypes.JSON) (error, *SparePartInventoryRepairRequestInfo) {
	var sparePartInventoryRepairRequestInfo SparePartInventoryRepairRequestInfo
	err := json.Unmarshal(serialisedData, &sparePartInventoryRepairRequestInfo)

	if err == nil {

		return nil, &sparePartInventoryRepairRequestInfo
	}

	return err, &sparePartInventoryRepairRequestInfo
}
func (v *SparePartInventoryRepairRequestInfo) Serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("marshalling error", err.Error())
		return []byte{}
	}
	return serialisedData
}
