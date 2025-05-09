package handler

import (
	"cx-micro-flake/pkg/orm"
	"encoding/json"
	"gorm.io/datatypes"
)

type Project struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ProjectRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

func (v *Project) getProjectInfo() *ProjectInfo {
	projectInfo := ProjectInfo{}
	json.Unmarshal(v.ObjectInfo, &projectInfo)
	return &projectInfo
}

type ProjectInfo struct {
	Name               string             `json:"name"`
	Created            string             `json:"created"`
	Template           string             `json:"template"`
	AvatarUrl          string             `json:"avatarUrl"`
	CreatedBy          string             `json:"createdBy"`
	IsDefault          bool               `json:"isDefault"`
	Description        string             `json:"description"`
	LastUpdatedAt      string             `json:"lastUpdatedAt"`
	LastUpdatedBy      string             `json:"lastUpdatedBy"`
	ProjectDatasource  orm.DatabaseConfig `json:"projectDatasource"`
	ProjectReferenceId string             `json:"projectReferenceId"`
}

func (ai *ProjectInfo) Serialize() []byte {
	rawData, _ := json.Marshal(ai)
	return rawData
}
