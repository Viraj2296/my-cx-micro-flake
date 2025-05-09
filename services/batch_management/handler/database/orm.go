package database

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/gorm"
)

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (error, []component.GeneralObject) {
	var err error

	var dbObjects []component.GeneralObject
	if len(objectCount) > 0 {
		err = database.Table(table).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
	} else {
		err = database.Table(table).Order(orderBy).Where(condition).Find(&dbObjects).Error
	}
	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}
}
func CreateRecordTrail(database *gorm.DB, objectInterface BatchManagementRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func CreateGeneralRecord(database *gorm.DB, table string, objectFields map[string]interface{}) error {
	err := database.Table(table).Create(objectFields).Error
	return err
}

func CreateFromGeneralObject(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	err := database.Table(table).Debug().Create(&objectInterface).Error
	if err != nil {
		return err, -1
	} else {
		return nil, objectInterface.Id
	}
}
func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	generalObject := component.GeneralObject{Id: recordId}
	err := database.Table(table).Find(&generalObject).Error
	if err != nil {
		return err, generalObject
	} else {
		return nil, generalObject
	}
}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	err = database.Table(table).Count(&numberOfRecords).Error

	if err != nil {
		return -1
	}
	return numberOfRecords
}

func CountByCondition(database *gorm.DB, table string, condition string) int64 {
	var err error
	var numberOfRecords int64
	err = database.Table(table).Where(condition).Count(&numberOfRecords).Error

	if err != nil {
		return -1
	}
	return numberOfRecords
}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (error, []component.GeneralObject) {
	var err error

	var dbObjects []component.GeneralObject
	if len(objectCount) > 0 {
		err = database.Table(table).Limit(objectCount[0]).Find(&dbObjects).Error
	} else {
		err = database.Table(table).Find(&dbObjects).Error
	}

	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}

}

func GetConditionalObjects(database *gorm.DB, table string, condition string, objectCount ...int) (error, []component.GeneralObject) {
	var err error
	var dbObjects []component.GeneralObject
	if len(objectCount) > 0 {
		err = database.Table(table).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
	} else {
		err = database.Table(table).Where(condition).Find(&dbObjects).Error
	}
	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}
}
func Delete(database *gorm.DB, table string, recordId int) error {
	var err error
	generalComponent := component.GeneralObject{Id: recordId}
	err = database.Table(table).Delete(&generalComponent).Error
	return err
}

func Archive(database *gorm.DB, table string, recordId int, lastUpdatedBy int) error {
	var err error
	existingObject := component.GeneralObject{Id: recordId}
	err = database.Table(table).Take(&existingObject).Error
	if err == nil {
		var objectFields = make(map[string]interface{})
		err := json.Unmarshal(existingObject.ObjectInfo, &objectFields)
		if err != nil {
			return err
		}
		objectFields["objectStatus"] = component.ObjectStatusArchived
		objectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		err = database.Table(table).Take(&existingObject).UpdateColumns(objectFields).Error
		return err
	}
	return err
}
func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	existingObject := component.GeneralObject{Id: recordId}
	err = database.Table(table).Take(&existingObject).UpdateColumns(updateObject).Error
	return err
}
func ArchiveObject(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	existingObject := component.GeneralObject{Id: objectInterface.Id}
	var objectFields map[string]interface{}
	err = json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
	if err != nil {
		return err
	}
	objectFields["objectStatus"] = "Archived"
	objectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	err = database.Table(table).Take(&existingObject).UpdateColumns(objectFields).Error
	return err
}
