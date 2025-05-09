package models

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type SparePartInventoryEmailEscalationInfo struct {
	SparePartId     int    `json:"sparePartId"`
	EmailRecipients []int  `json:"emailRecipients"`
	ObjectStatus    string `json:"objectStatus"`
}

func (v *SparePartInventoryEmailEscalationInfo) Serialised() []byte {
	marshal, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	} else {
		return marshal
	}
}

func GetSparePartInventoryEmailEscalationInfo(objectInfo datatypes.JSON) *SparePartInventoryEmailEscalationInfo {
	var SparePartEmailEscalationInfo SparePartInventoryEmailEscalationInfo
	err := json.Unmarshal(objectInfo, &SparePartEmailEscalationInfo)
	if err == nil {
		return &SparePartEmailEscalationInfo
	}

	return nil
}
