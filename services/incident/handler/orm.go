package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

var emptyObject interface{}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == IncidentRecordTrailTable:
		var dbObjects []IncidentRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentInventoryTable:
		var dbObjects []IncidentInventory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentInventory{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentInventory{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentDeliveryTable:
		var dbObjects []IncidentDelivery
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentDelivery{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentDelivery{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentSafetyTable:
		var dbObjects []IncidentSafety
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentSafety{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentSafety{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentQualityTable:
		var dbObjects []IncidentQuality
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentQuality{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentQuality{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentProductivityTable:
		var dbObjects []IncidentProductivity
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentProductivity{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentProductivity{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface IncidentRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == IncidentComponentTable:
		object := IncidentComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentSafetyCategoryTable:
		object := IncidentSafetyCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentQualityCategoryTable:
		object := IncidentQualityCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentDeliveryCategoryTable:
		object := IncidentDeliveryCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentInventoryCategoryTable:
		object := IncidentInventoryCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentProductivityCategoryTable:
		object := IncidentProductivityCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentInventoryTable:
		object := IncidentInventory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentDeliveryTable:
		object := IncidentDelivery{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentSafetyTable:
		object := IncidentSafety{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentQualityTable:
		object := IncidentQuality{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentProductivityTable:
		object := IncidentProductivity{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == IncidentTargetTable:
		object := IncidentTarget{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {

	case table == IncidentComponentTable:
		dbObject := IncidentComponent{Id: recordId}
		err = database.Debug().Model(&IncidentComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentSafetyTable:
		dbObject := IncidentSafety{Id: recordId}
		err = database.Debug().Model(&IncidentSafety{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentQualityTable:
		dbObject := IncidentQuality{Id: recordId}
		err = database.Debug().Model(&IncidentQuality{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentDeliveryTable:
		dbObject := IncidentDelivery{Id: recordId}
		err = database.Debug().Model(&IncidentDelivery{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentInventoryTable:
		dbObject := IncidentInventory{Id: recordId}
		err = database.Debug().Model(&IncidentInventory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentProductivityTable:
		dbObject := IncidentProductivity{Id: recordId}
		err = database.Debug().Model(&IncidentProductivity{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentQualityCategoryTable:
		dbObject := IncidentQualityCategory{Id: recordId}
		err = database.Debug().Model(&IncidentQualityCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentDeliveryCategoryTable:
		dbObject := IncidentDeliveryCategory{Id: recordId}
		err = database.Debug().Model(&IncidentDeliveryCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentInventoryCategoryTable:
		dbObject := IncidentInventoryCategory{Id: recordId}
		err = database.Debug().Model(&IncidentInventoryCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentProductivityCategoryTable:
		dbObject := IncidentProductivityCategory{Id: recordId}
		err = database.Debug().Model(&IncidentProductivityCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentSafetyCategoryTable:
		dbObject := IncidentSafetyCategory{Id: recordId}
		err = database.Debug().Model(&IncidentSafetyCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == IncidentTargetTable:
		dbObject := IncidentTarget{Id: recordId}
		err = database.Debug().Model(&IncidentTarget{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == IncidentComponentTable:

		var dbObjects []IncidentComponent
		if len(objectCount) > 0 {
			err = database.Model(&IncidentComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentSafetyCategoryTable:

		var dbObjects []IncidentSafetyCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentSafetyCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentSafetyCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentQualityCategoryTable:

		var dbObjects []IncidentQualityCategory
		if len(objectCount) > 0 {
			err = database.Model(&IncidentQualityCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentQualityCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentDeliveryCategoryTable:

		var dbObjects []IncidentDeliveryCategory
		if len(objectCount) > 0 {
			err = database.Model(&IncidentDeliveryCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentDeliveryCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentInventoryCategoryTable:

		var dbObjects []IncidentInventoryCategory
		if len(objectCount) > 0 {
			err = database.Model(&IncidentInventoryCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentInventoryCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentProductivityCategoryTable:

		var dbObjects []IncidentProductivityCategory
		if len(objectCount) > 0 {
			err = database.Model(&IncidentProductivityCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentProductivityCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentSafetyTable:

		var dbObjects []IncidentSafety
		if len(objectCount) > 0 {
			err = database.Model(&IncidentSafety{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentSafety{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentQualityTable:

		var dbObjects []IncidentQuality
		if len(objectCount) > 0 {
			err = database.Model(&IncidentQuality{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentQuality{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentDeliveryTable:

		var dbObjects []IncidentDelivery
		if len(objectCount) > 0 {
			err = database.Model(&IncidentDelivery{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentDelivery{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentInventoryTable:

		var dbObjects []IncidentInventory
		if len(objectCount) > 0 {
			err = database.Model(&IncidentInventory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentInventory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentProductivityTable:

		var dbObjects []IncidentProductivity
		if len(objectCount) > 0 {
			err = database.Model(&IncidentProductivity{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentProductivity{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentTargetTable:

		var dbObjects []IncidentTarget
		if len(objectCount) > 0 {
			err = database.Model(&IncidentTarget{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentTarget{}).Find(&dbObjects).Error
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
	case table == IncidentSafetyCategoryTable:
		var dbObjects []IncidentSafetyCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentSafetyCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentSafetyCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentQualityCategoryTable:
		var dbObjects []IncidentQualityCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentQualityCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentQualityCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == IncidentDeliveryCategoryTable:
		var dbObjects []IncidentDeliveryCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentDeliveryCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentDeliveryCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentInventoryCategoryTable:
		var dbObjects []IncidentInventoryCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentInventoryCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentInventoryCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentProductivityCategoryTable:
		var dbObjects []IncidentProductivityCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentProductivityCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentProductivityCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentSafetyTable:
		var dbObjects []IncidentSafety
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentSafety{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&IncidentSafety{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentQualityTable:
		var dbObjects []IncidentQuality
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentQuality{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentQuality{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentDeliveryTable:
		var dbObjects []IncidentDelivery
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentDelivery{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentDelivery{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentInventoryTable:
		var dbObjects []IncidentInventory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentInventory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentInventory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == IncidentProductivityTable:
		var dbObjects []IncidentProductivity
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentProductivity{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentProductivity{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentTargetTable:
		var dbObjects []IncidentTarget
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentTarget{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentTarget{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == IncidentComponentTable:
		var dbObjects []IncidentComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&IncidentComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&IncidentComponent{}).Where(condition).Find(&dbObjects).Error
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

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == IncidentQualityCategoryTable:
		err = database.Debug().Model(&IncidentQualityCategory{}).Count(&numberOfRecords).Error
	case table == IncidentSafetyCategoryTable:
		err = database.Debug().Model(&IncidentSafetyCategory{}).Count(&numberOfRecords).Error
	case table == IncidentDeliveryCategoryTable:
		err = database.Debug().Model(&IncidentDeliveryCategory{}).Count(&numberOfRecords).Error
	case table == IncidentInventoryCategoryTable:
		err = database.Debug().Model(&IncidentInventoryCategory{}).Count(&numberOfRecords).Error
	case table == IncidentProductivityCategoryTable:
		err = database.Debug().Model(&IncidentProductivityCategory{}).Count(&numberOfRecords).Error
	case table == IncidentTargetTable:
		err = database.Debug().Model(&IncidentTarget{}).Count(&numberOfRecords).Error

	default:
		return -1
	}
	if err != nil {
		return -1
	}
	return numberOfRecords
}

func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error

	return err

}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == IncidentDeliveryTable:
		err = database.Debug().Model(&IncidentDelivery{}).Take(&IncidentDelivery{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentProductivityTable:
		err = database.Debug().Model(&IncidentProductivity{}).Take(&IncidentProductivity{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentInventoryTable:
		err = database.Debug().Model(&IncidentInventory{}).Take(&IncidentInventory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentSafetyTable:
		err = database.Debug().Model(&IncidentSafety{}).Take(&IncidentSafety{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentQualityTable:
		err = database.Debug().Model(&IncidentQuality{}).Take(&IncidentQuality{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentSafetyCategoryTable:
		err = database.Debug().Model(&IncidentSafetyCategory{}).Take(&IncidentSafetyCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentProductivityCategoryTable:
		err = database.Debug().Model(&IncidentProductivityCategory{}).Take(&IncidentProductivityCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentInventoryCategoryTable:
		err = database.Debug().Model(&IncidentInventoryCategory{}).Take(&IncidentInventoryCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentDeliveryCategoryTable:
		err = database.Debug().Model(&IncidentDeliveryCategory{}).Take(&IncidentDeliveryCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentQualityCategoryTable:
		err = database.Debug().Model(&IncidentQualityCategory{}).Take(&IncidentQualityCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentTargetTable:
		err = database.Debug().Model(&IncidentTarget{}).Take(&IncidentTarget{Id: recordId}).UpdateColumns(updateObject).Error
	case table == IncidentComponentTable:
		err = database.Debug().Model(&IncidentComponent{}).Take(&IncidentComponent{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func ArchiveReferenceObjects(database *gorm.DB, table string, referenceField string, id int) error {
	var err error

	conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$." + referenceField + "\")) =" + strconv.Itoa(id)
	listOfObjects, err := GetConditionalObjects(database, table, conditionString)
	for _, object := range *listOfObjects {
		updateObject := make(map[string]interface{})
		var objectFields map[string]interface{}
		json.Unmarshal(object.ObjectInfo, &objectFields)
		objectFields["objectStatus"] = "Archived"
		objectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		serializedObject, _ := json.Marshal(objectFields)
		updateObject["object_info"] = serializedObject
		switch {
		case table == IncidentSafetyCategoryTable:
			err = database.Debug().Model(&IncidentSafetyCategory{}).Take(&IncidentSafetyCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == IncidentQualityCategoryTable:
			err = database.Debug().Model(&IncidentQualityCategory{}).Take(&IncidentQualityCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == IncidentDeliveryCategoryTable:
			err = database.Debug().Model(&IncidentDeliveryCategory{}).Take(&IncidentDeliveryCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == IncidentInventoryCategoryTable:
			err = database.Debug().Model(&IncidentInventoryCategory{}).Take(&IncidentInventoryCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == IncidentProductivityCategoryTable:
			err = database.Debug().Model(&IncidentProductivityCategory{}).Take(&IncidentProductivityCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == IncidentTargetTable:
			err = database.Debug().Model(&IncidentTarget{}).Take(&IncidentTarget{Id: object.Id}).UpdateColumns(updateObject).Error
		default:
			err = errors.New(UpdateUnknownObjectType)
		}
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
	case table == IncidentSafetyCategoryTable:
		err = database.Debug().Model(&IncidentSafetyCategory{}).Take(&IncidentSafetyCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == IncidentQualityCategoryTable:
		err = database.Debug().Model(&IncidentQualityCategory{}).Take(&IncidentQualityCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == IncidentDeliveryCategoryTable:
		err = database.Debug().Model(&IncidentDeliveryCategory{}).Take(&IncidentDeliveryCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == IncidentInventoryCategoryTable:
		err = database.Debug().Model(&IncidentInventoryCategory{}).Take(&IncidentInventoryCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == IncidentProductivityCategoryTable:
		err = database.Debug().Model(&IncidentProductivityCategory{}).Take(&IncidentProductivityCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == IncidentTargetTable:
		err = database.Debug().Model(&IncidentTarget{}).Take(&IncidentTarget{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *IncidentService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *IncidentService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
