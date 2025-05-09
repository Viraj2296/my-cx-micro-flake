package model

import (
	"encoding/json"
	"gorm.io/datatypes"
)

// DeviceUserInfo are holding information about how the devices can authenticate with MES, only configured devices are allowed
type DeviceUserInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	MacAddress    string `json:"macAddress"`
	CreatedBy     int    `json:"createdBy"`
	CreatedAt     string `json:"createdAt"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	ObjectStatus  string `json:"objectStatus"`
	IPAddress     string `json:"ipAddress"`
}

func GetDeviceInfo(serialisedData datatypes.JSON) *DeviceUserInfo {
	deviceUserInfo := DeviceUserInfo{}
	err := json.Unmarshal(serialisedData, &deviceUserInfo)
	if err != nil {
		return &DeviceUserInfo{}
	}
	return &deviceUserInfo
}

func (v *DeviceUserInfo) serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON{}
	}
	return serialisedData

}
