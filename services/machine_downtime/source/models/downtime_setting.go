package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type MachineDowntimeSettingInfo struct {
	IsEnableEscalation       bool  `json:"isEnableEscalation"`
	PrimaryEmailRecipients   []int `json:"primaryEmailRecipients"`
	EscalationWaitingPeriod  int   `json:"escalationWaitingPeriod"`
	InitialWaitingPeriod     int   `json:"initialWaitingPeriod"`
	PushNotificationJobRoles []int `json:"pushNotificationJobRoles"`
	JobAssigningUsers        []int `json:"jobAssigningUsers"`
}

func GetMachineDowntimeSettingInfo(serialisedData datatypes.JSON) *MachineDowntimeSettingInfo {
	labourManagementSettingInfo := MachineDowntimeSettingInfo{}
	err := json.Unmarshal(serialisedData, &labourManagementSettingInfo)
	if err != nil {
		return &MachineDowntimeSettingInfo{}
	}
	return &labourManagementSettingInfo
}
