package models

import (
	"encoding/json"

	"go.cerex.io/transcendflow/component"
	"gorm.io/datatypes"
)

type MachineDowntimeInfo struct {
	JobReferenceId string                    `json:"jobReferenceId"`
	CheckInTime    string                    `json:"checkInTime"`
	CheckInDate    string                    `json:"checkInDate"`
	CheckOutDate   string                    `json:"checkOutDate"`
	CheckOutTime   string                    `json:"checkOutTime"`
	CheckInUserId  int                       `json:"checkInUserId"`
	CheckOutUserId int                       `json:"checkOutUserId"`
	MachineId      int                       `json:"machineId"`
	CreatedAt      string                    `json:"createdAt"`
	LastUpdatedAt  string                    `json:"lastUpdatedAt"`
	CreatedBy      int                       `json:"createdBy"`
	LastUpdatedBy  int                       `json:"lastUpdatedBy"`
	CanCheckOut    bool                      `json:"canCheckOut"`
	CanCheckIn     bool                      `json:"canCheckIn"`
	CanCancel      bool                      `json:"canCancel"`
	FaultType      int                       `json:"faultType"`
	FaultCode      int                       `json:"faultCode"`
	Status         int                       `json:"status"`
	AssignedUserId int                       `json:"assignedUserId"`
	Remarks        string                    `json:"remarks"`
	ActionRemarks  []component.ActionRemarks `json:"actionRemarks"`
	ObjectStatus   string                    `json:"objectStatus"`
}

func GetMachineDowntimeInfo(objectInfo datatypes.JSON) *MachineDowntimeInfo {
	var downtimeInfo MachineDowntimeInfo
	err := json.Unmarshal(objectInfo, &downtimeInfo)
	if err == nil {
		return &downtimeInfo
	}

	return nil
}

func (v *MachineDowntimeInfo) Serialised() ([]byte, error) {
	return json.Marshal(v)
}
