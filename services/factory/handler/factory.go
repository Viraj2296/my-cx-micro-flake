package handler

import (
	"encoding/json"
	"gorm.io/datatypes"
	"time"
)

type FactoryRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type FactoryComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactorySite struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryBuilding struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryGeneralFacility struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryCustomerAsset struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type SiteInfo struct {
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	Country       string    `json:"country"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func (v *FactorySite) getFactorySiteInfo() *SiteInfo {
	siteInfo := SiteInfo{}
	json.Unmarshal(v.ObjectInfo, &siteInfo)
	return &siteInfo
}

type FactoryPlant struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryLocation struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryDepartment struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryUnit struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryLevel struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FactoryArea struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type DepartmentInfo struct {
	Name          string    `json:"name"`
	Plant         int       `json:"plant"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	ObjectStatus  string    `json:"objectStatus"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func (v *FactoryDepartment) getDepartmentInfo() *DepartmentInfo {
	departmentInfo := DepartmentInfo{}
	json.Unmarshal(v.ObjectInfo, &departmentInfo)
	return &departmentInfo
}

type FactoryDepartmentSection struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type SectionInfo struct {
	Name          string    `json:"name"`
	Department    int       `json:"department"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	ObjectStatus  string    `json:"objectStatus"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func (v *FactoryDepartmentSection) getSectionInfo() *SectionInfo {
	sectionInfo := SectionInfo{}
	json.Unmarshal(v.ObjectInfo, &sectionInfo)
	return &sectionInfo
}
