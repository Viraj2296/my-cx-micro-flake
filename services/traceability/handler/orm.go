package handler

import (
	"cx-micro-flake/pkg/common/component"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

var emptyObject interface{}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == TraceabilityComponentTable:
		var dbObjects []TraceabilityComponent
		if len(objectCount) > 0 {
			err = database.Model(&TraceabilityComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TraceabilityComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TraceabilityOrdersTable:
		var dbObjects []TraceabilityOrder
		if len(objectCount) > 0 {
			err = database.Model(&TraceabilityOrder{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TraceabilityOrder{}).Find(&dbObjects).Error
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

func GetConditionalObjects(database *gorm.DB, table string, condition string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error
	switch {

	case table == TraceabilityComponentTable:
		var dbObjects []TraceabilityComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TraceabilityComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TraceabilityComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == TraceabilityOrdersTable:
		var dbObjects []TraceabilityOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TraceabilityOrder{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TraceabilityOrder{}).Where(condition).Find(&dbObjects).Error
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

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == TraceabilityComponentTable:
		err = database.Debug().Model(&TraceabilityComponent{}).Take(&TraceabilityComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == TraceabilityOrdersTable:
		err = database.Debug().Model(&TraceabilityOrder{}).Take(&TraceabilityOrder{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == TraceabilityComponentTable:
		object := TraceabilityComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == TraceabilityOrdersTable:
		object := TraceabilityOrder{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id

	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == TraceabilityRecordTrailTable:
		var dbObjects []TraceabilityRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TraceabilityRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&TraceabilityRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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

func CreateRecordTrail(database *gorm.DB, objectInterface TraceabilityRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == TraceabilityComponentTable:
		dbObject := TraceabilityComponent{Id: recordId}
		err = database.Debug().Model(&TraceabilityComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == TraceabilityOrdersTable:
		dbObject := TraceabilityOrder{Id: recordId}
		err = database.Debug().Model(&TraceabilityOrder{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {

	case table == TraceabilityOrdersTable:
		err = database.Debug().Model(&TraceabilityOrder{}).Count(&numberOfRecords).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	if err != nil {
		return -1
	} else {
		return numberOfRecords
	}
}

func (v *TraceabilityService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(*listOfObjects)
					for _, referenceObject := range *listOfObjects {
						v.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (v *TraceabilityService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					for _, referenceObject := range *listOfObjects {
						fmt.Println("referenceTable : ", referenceTable, " id :", referenceObject)
						ArchiveObject(dbConnection, referenceTable, referenceObject)
						v.CreateUserRecordMessage(ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
						v.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
					}
				}
			}

		}
	}
}
func ArchiveObject(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	updateObject := make(map[string]interface{})
	var objectFields map[string]interface{}
	json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
	objectFields["objectStatus"] = "Archived"
	serializedObject, _ := json.Marshal(objectFields)
	updateObject["object_info"] = serializedObject

	switch {
	case table == TraceabilityOrdersTable:
		err = database.Debug().Model(&TraceabilityOrder{}).Take(&TraceabilityOrder{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}
