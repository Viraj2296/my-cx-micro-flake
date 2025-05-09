package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type MachineDowntimeEmailEscalationInfo struct {
	DowntimeId      int    `json:"downtimeId"`
	CreatedAt       string `json:"createdAt"`
	CreatedBy       int    `json:"createdBy"`
	EmailRecipients []int  `json:"emailRecipients"`
	ObjectStatus    string `json:"objectStatus"`
}

func GetMachineDowntimeEmailEscalationInfo(objectInfo datatypes.JSON) *MachineDowntimeEmailEscalationInfo {
	var machineDowntimeEmailEscalationInfo MachineDowntimeEmailEscalationInfo
	err := json.Unmarshal(objectInfo, &machineDowntimeEmailEscalationInfo)
	if err == nil {
		return &machineDowntimeEmailEscalationInfo
	}

	return nil
}

func (v *MachineDowntimeEmailEscalationInfo) Serialised() []byte {
	marshal, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	} else {
		return marshal
	}
}
