package models

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type SparePartInventoryMasterInfo struct {
	Brand                string `json:"brand"`
	MinQty               int    `json:"minQty"`
	Customer             string `json:"customer"`
	Location             int    `json:"location"`
	CreatedAt            string `json:"createdAt"`
	OnHandQty            int    `json:"onHandQty"`
	MachineIds           []int  `json:"machineIds"`
	EquipmentType        string `json:"equipmentType"`
	RepairRemarks        string `json:"repairRemarks"`
	SparePartImage       string `json:"sparePartImage"`
	SparePartNumber      string `json:"sparePartNumber"`
	SparePartDescription string `json:"sparePartDescription"`
}

func GetSparePartInventoryMasterInfo(serialisedData datatypes.JSON) (error, *SparePartInventoryMasterInfo) {
	var sparePartInventoryMasterInfo SparePartInventoryMasterInfo
	err := json.Unmarshal(serialisedData, &sparePartInventoryMasterInfo)

	if err == nil {

		return nil, &sparePartInventoryMasterInfo
	}

	return err, &sparePartInventoryMasterInfo
}
func (v *SparePartInventoryMasterInfo) Serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("marshalling error", err.Error())
		return []byte{}
	}
	return serialisedData
}
