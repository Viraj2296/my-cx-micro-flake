package models

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type SignalHistoryInfo struct {
	CreatedAt               int64 `json:"createdAt"`
	SignalMachine           int   `json:"signalMachine"`
	IsNotificationProcessed bool  `json:"isNotificationProcessed"` // if the notification is processed, we don't need to load it next time
}

func (v *SignalHistoryInfo) Serialised() []byte {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return serialisedData
}
func GetSignalHistoryInfo(objectInfo datatypes.JSON) *SignalHistoryInfo {
	var signalHistoryInfo SignalHistoryInfo
	err := json.Unmarshal(objectInfo, &signalHistoryInfo)
	if err == nil {
		return &signalHistoryInfo
	}
	return nil
}
