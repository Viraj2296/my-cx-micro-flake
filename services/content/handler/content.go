package handler

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type SharePermission struct {
	SharedWith int `json:"sharedWith"`
	RoleId     int `json:"roleId"`
}

type ContentRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}
type ContenetMasterInfo struct {
	Url              string            `json:"url"`
	Name             string            `json:"name"`
	Size             string            `json:"size"`
	IsFile           bool              `json:"isFile"`
	ChainReference   string            `json:"chainReference"`
	Path             string            `json:"path"`
	MIMEType         string            `json:"mimeType"`
	CreatedAt        string            `json:"createdAt"`
	CreatedBy        int               `json:"createdBy"`
	LastUpdatedAt    string            `json:"lastUpdatedAt"`
	LastUpdatedBy    int               `json:"lastUpdatedBy"`
	ObjectStatus     string            `json:"objectStatus"`
	IsShared         bool              `json:"isShared"`
	Share            []SharePermission `json:"share"`
	Tags             []string          `json:"tags"`
	FavoriteList     []int             `json:"favoriteList"`
	FileTypeIcon     string            `json:"fileTypeIcon"`
	FilePreviewImage string            `json:"filePreviewImage"`
	Description      string            `json:"description"`
}

type ContentComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (contentInfo *ContenetMasterInfo) Serialize() []byte {
	rawData, _ := json.Marshal(contentInfo)
	return rawData
}

type ContentMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (cm *ContentMaster) GetContentInfo() *ContenetMasterInfo {
	contentInfo := ContenetMasterInfo{}
	json.Unmarshal(cm.ObjectInfo, &contentInfo)
	return &contentInfo
}

type GroupByCardView struct {
	GroupByField string                   `json:"groupByField"`
	Cards        []map[string]interface{} `json:"cards"`
}
