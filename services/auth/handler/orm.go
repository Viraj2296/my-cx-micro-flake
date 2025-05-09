package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"strconv"
)

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == AuthRecordTrailTable:
		var dbObjects []AuthRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AuthRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&AuthRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err

	default:
		return nil, errors.New(GetUnknownObjectType)
	}
}
func CreateRecordTrail(database *gorm.DB, objectInterface AuthRecordTrail) (error, int) {
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

func (as *AuthService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := as.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := as.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			err, listOfObjects := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(listOfObjects)
					for _, referenceObject := range listOfObjects {
						as.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (as *AuthService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := as.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := as.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			err, listOfObjects := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(listOfObjects) > 0 {
					for _, referenceObject := range listOfObjects {
						err := ArchiveObject(dbConnection, referenceTable, referenceObject)
						if err == nil {
							err := as.CreateUserRecordMessage(ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
							if err == nil {
								as.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
							}
						}

					}
				}
			}

		}
	}
}
