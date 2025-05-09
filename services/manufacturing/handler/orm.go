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

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == ManufacturingRecordTrailTable:
		var dbObjects []ManufacturingRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ManufacturingRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ManufacturingRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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

func CreateRecordTrail(database *gorm.DB, objectInterface ManufacturingRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}
func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == ManufacturingMouldingPartTable:
		object := ManufacturingMouldingPart{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ManufacturingAssemblyPartTable:
		object := ManufacturingAssemblyPart{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ManufacturingVendorMasterTable:
		object := ManufacturingVendorMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == ManufacturingMouldingPartTable:
		dbObject := ManufacturingMouldingPart{Id: recordId}
		err = database.Debug().Model(&ManufacturingMouldingPart{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ManufacturingAssemblyPartTable:
		dbObject := ManufacturingAssemblyPart{Id: recordId}
		err = database.Debug().Model(&ManufacturingAssemblyPart{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ManufacturingVendorMasterTable:
		dbObject := ManufacturingVendorMaster{Id: recordId}
		err = database.Debug().Model(&ManufacturingVendorMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}
}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == ManufacturingComponentTable:

		var dbObjects []ManufacturingComponent
		if len(objectCount) > 0 {
			err = database.Model(&ManufacturingComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ManufacturingMouldingPartTable:

		var dbObjects []ManufacturingMouldingPart
		if len(objectCount) > 0 {
			err = database.Model(&ManufacturingMouldingPart{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingMouldingPart{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ManufacturingAssemblyPartTable:
		var dbObjects []ManufacturingAssemblyPart
		if len(objectCount) > 0 {
			err = database.Model(&ManufacturingAssemblyPart{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingAssemblyPart{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ManufacturingVendorMasterTable:
		var dbObjects []ManufacturingVendorMaster
		if len(objectCount) > 0 {
			err = database.Model(&ManufacturingVendorMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingVendorMaster{}).Find(&dbObjects).Error
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
	case table == ManufacturingMouldingPartTable:
		var dbObjects []ManufacturingMouldingPart
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ManufacturingMouldingPart{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingMouldingPart{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ManufacturingAssemblyPartTable:
		var dbObjects []ManufacturingAssemblyPart
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ManufacturingAssemblyPart{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingAssemblyPart{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ManufacturingVendorMasterTable:
		var dbObjects []ManufacturingVendorMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ManufacturingVendorMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ManufacturingVendorMaster{}).Where(condition).Find(&dbObjects).Error
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

func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	switch {
	case table == ManufacturingMouldingPartTable:
		err = database.Debug().Model(&ManufacturingMouldingPart{}).Delete(&ManufacturingMouldingPart{Id: objectInterface.Id}).Error
	case table == ManufacturingAssemblyPartTable:
		err = database.Debug().Model(&ManufacturingAssemblyPart{}).Delete(&ManufacturingAssemblyPart{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == ManufacturingMouldingPartTable:
		err = database.Debug().Model(&ManufacturingMouldingPart{}).Count(&numberOfRecords).Error
	case table == ManufacturingAssemblyPartTable:
		err = database.Debug().Model(&ManufacturingAssemblyPart{}).Count(&numberOfRecords).Error
	case table == ManufacturingVendorMasterTable:
		err = database.Debug().Model(&ManufacturingVendorMaster{}).Count(&numberOfRecords).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	if err != nil {
		return -1
	} else {
		return numberOfRecords
	}
}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == ManufacturingMouldingPartTable:
		err = database.Debug().Model(&ManufacturingMouldingPart{}).Take(&ManufacturingMouldingPart{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ManufacturingAssemblyPartTable:
		err = database.Debug().Model(&ManufacturingAssemblyPart{}).Take(&ManufacturingAssemblyPart{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ManufacturingVendorMasterTable:
		err = database.Debug().Model(&ManufacturingVendorMaster{}).Take(&ManufacturingVendorMaster{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
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
	case table == ManufacturingMouldingPartTable:
		err = database.Debug().Model(&ManufacturingMouldingPart{}).Take(&ManufacturingMouldingPart{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ManufacturingAssemblyPartTable:
		err = database.Debug().Model(&ManufacturingAssemblyPart{}).Take(&ManufacturingAssemblyPart{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *ManufacturingService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *ManufacturingService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
