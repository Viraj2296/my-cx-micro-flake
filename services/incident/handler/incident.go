package handler

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type IncidentRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type IncidentComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentSafetyCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type IncidentQualityCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type IncidentDeliveryCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type IncidentInventoryCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type IncidentProductivityCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentInventory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentDelivery struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentSafety struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type IncidentQuality struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentProductivity struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentTarget struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type IncidentInfo struct {
	Year          int       `json:"year"`
	Month         int       `json:"month"`
	Day           int       `json:"day"`
	Target        int       `json:"target"`
	IsNoWorkDay   bool      `json:"isNoWorkDay"`
	Department    int       `json:"department"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
	ObjectStatus  string    `json:"objectStatus"`
	Incidents     Incidents `json:"incidents"`
}

type Incidents struct {
	Id        int  `json:"id"`
	Indicator bool `json:"indicator"`
}

type IncidentCategoryInfo struct {
	Name          string    `json:"name"`
	Site          int       `json:"site"`
	ColorCode     string    `json:"colorCode"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
	ObjectStatus  string    `json:"objectStatus"`
	Department    int       `json:"department"`
}

func (v *IncidentSafetyCategory) getSafetyCategoryInfo() *IncidentCategoryInfo {
	incidentSafetyCategoryInfo := IncidentCategoryInfo{}
	json.Unmarshal(v.ObjectInfo, &incidentSafetyCategoryInfo)
	return &incidentSafetyCategoryInfo
}
func (v *IncidentQualityCategory) getQualityCategoryInfo() *IncidentCategoryInfo {
	incidentSafetyCategoryInfo := IncidentCategoryInfo{}
	json.Unmarshal(v.ObjectInfo, &incidentSafetyCategoryInfo)
	return &incidentSafetyCategoryInfo
}
func (v *IncidentDeliveryCategory) getDeliveryCategoryInfo() *IncidentCategoryInfo {
	incidentSafetyCategoryInfo := IncidentCategoryInfo{}
	json.Unmarshal(v.ObjectInfo, &incidentSafetyCategoryInfo)
	return &incidentSafetyCategoryInfo
}
func (v *IncidentInventoryCategory) getInventoryCategoryInfo() *IncidentCategoryInfo {
	incidentSafetyCategoryInfo := IncidentCategoryInfo{}
	json.Unmarshal(v.ObjectInfo, &incidentSafetyCategoryInfo)
	return &incidentSafetyCategoryInfo
}
func (v *IncidentProductivityCategory) getProductivityCategoryInfo() *IncidentCategoryInfo {
	incidentSafetyCategoryInfo := IncidentCategoryInfo{}
	json.Unmarshal(v.ObjectInfo, &incidentSafetyCategoryInfo)
	return &incidentSafetyCategoryInfo
}
