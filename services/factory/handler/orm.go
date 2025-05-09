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
	case table == FactoryRecordTrailTable:
		var dbObjects []FactoryRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FactoryRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface FactoryRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == FactoryComponentTable:
		object := FactoryComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryLocationTable:
		object := FactoryLocation{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactorySiteTable:
		object := FactorySite{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryPlantTable:
		object := FactoryPlant{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryDepartmentTable:
		object := FactoryDepartment{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryDepartmentSectionTable:
		object := FactoryDepartmentSection{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryGeneralFacilityTable:
		object := FactoryGeneralFacility{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryCustomerAssetTable:
		object := FactoryCustomerAsset{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryBuildingTable:
		object := FactoryBuilding{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryUnitTable:
		object := FactoryUnit{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryLevelTable:
		object := FactoryLevel{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FactoryAreaTable:
		object := FactoryArea{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == FactorySiteTable:
		dbObject := FactorySite{Id: recordId}
		err = database.Debug().Model(&FactorySite{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryLocationTable:
		dbObject := FactoryLocation{Id: recordId}
		err = database.Debug().Model(&FactoryLocation{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryPlantTable:
		dbObject := FactoryPlant{Id: recordId}
		err = database.Debug().Model(&FactoryPlant{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryComponentTable:
		dbObject := FactoryComponent{Id: recordId}
		err = database.Debug().Model(&FactoryComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryDepartmentTable:
		dbObject := FactoryDepartment{Id: recordId}
		err = database.Debug().Model(&FactoryDepartment{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryDepartmentSectionTable:
		dbObject := FactoryDepartmentSection{Id: recordId}
		err = database.Debug().Model(&FactoryDepartmentSection{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryGeneralFacilityTable:
		dbObject := FactoryGeneralFacility{Id: recordId}
		err = database.Debug().Model(&FactoryGeneralFacility{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryCustomerAssetTable:
		dbObject := FactoryCustomerAsset{Id: recordId}
		err = database.Debug().Model(&FactoryCustomerAsset{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryBuildingTable:
		dbObject := FactoryBuilding{Id: recordId}
		err = database.Debug().Model(&FactoryBuilding{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryUnitTable:
		dbObject := FactoryUnit{Id: recordId}
		err = database.Debug().Model(&FactoryUnit{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryLevelTable:
		dbObject := FactoryLevel{Id: recordId}
		err = database.Debug().Model(&FactoryLevel{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FactoryAreaTable:
		dbObject := FactoryArea{Id: recordId}
		err = database.Debug().Model(&FactoryArea{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == FactoryComponentTable:

		var dbObjects []FactoryComponent
		if len(objectCount) > 0 {
			err = database.Model(&FactoryComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == FactorySiteTable:
		var dbObjects []FactorySite
		if len(objectCount) > 0 {
			err = database.Model(&FactorySite{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactorySite{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryDepartmentTable:
		var dbObjects []FactoryDepartment
		if len(objectCount) > 0 {
			err = database.Model(&FactoryDepartment{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryDepartment{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryPlantTable:
		var dbObjects []FactoryPlant
		if len(objectCount) > 0 {
			err = database.Model(&FactoryPlant{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryPlant{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryLocationTable:
		var dbObjects []FactoryLocation
		if len(objectCount) > 0 {
			err = database.Model(&FactoryLocation{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryLocation{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryDepartmentSectionTable:
		var dbObjects []FactoryDepartmentSection
		if len(objectCount) > 0 {
			err = database.Model(&FactoryDepartmentSection{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryDepartmentSection{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryGeneralFacilityTable:
		var dbObjects []FactoryGeneralFacility
		if len(objectCount) > 0 {
			err = database.Model(&FactoryGeneralFacility{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryGeneralFacility{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryCustomerAssetTable:
		var dbObjects []FactoryCustomerAsset
		if len(objectCount) > 0 {
			err = database.Model(&FactoryCustomerAsset{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryCustomerAsset{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryBuildingTable:
		var dbObjects []FactoryBuilding
		if len(objectCount) > 0 {
			err = database.Model(&FactoryBuilding{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryBuilding{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryUnitTable:
		var dbObjects []FactoryUnit
		if len(objectCount) > 0 {
			err = database.Model(&FactoryUnit{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryUnit{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryLevelTable:
		var dbObjects []FactoryLevel
		if len(objectCount) > 0 {
			err = database.Model(&FactoryLevel{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryLevel{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryAreaTable:
		var dbObjects []FactoryArea
		if len(objectCount) > 0 {
			err = database.Model(&FactoryArea{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryArea{}).Find(&dbObjects).Error
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
	case table == FactoryLocationTable:
		var dbObjects []FactoryLocation
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryLocation{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryLocation{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FactoryDepartmentTable:
		var dbObjects []FactoryDepartment
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryDepartment{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryDepartment{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FactorySiteTable:
		var dbObjects []FactorySite
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactorySite{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactorySite{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryPlantTable:
		var dbObjects []FactoryPlant
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryPlant{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryPlant{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryComponentTable:
		var dbObjects []FactoryComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryDepartmentSectionTable:
		var dbObjects []FactoryDepartmentSection
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryDepartmentSection{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryDepartmentSection{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryGeneralFacilityTable:
		var dbObjects []FactoryGeneralFacility
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryGeneralFacility{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryGeneralFacility{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryCustomerAssetTable:
		var dbObjects []FactoryCustomerAsset
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryCustomerAsset{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryCustomerAsset{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryBuildingTable:
		var dbObjects []FactoryBuilding
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryBuilding{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryBuilding{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryUnitTable:
		var dbObjects []FactoryUnit
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryUnit{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryUnit{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryLevelTable:
		var dbObjects []FactoryLevel
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryLevel{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryLevel{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FactoryAreaTable:
		var dbObjects []FactoryArea
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FactoryArea{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FactoryArea{}).Where(condition).Find(&dbObjects).Error
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
	case table == FactoryLocationTable:
		err = database.Debug().Model(&FactoryLocation{}).Delete(&FactoryLocation{Id: objectInterface.Id}).Error
	case table == FactorySiteTable:
		err = database.Debug().Model(&FactorySite{}).Delete(&FactorySite{Id: objectInterface.Id}).Error
	case table == FactoryPlantTable:
		err = database.Debug().Model(&FactoryPlant{}).Delete(&FactoryPlant{Id: objectInterface.Id}).Error
	case table == FactoryDepartmentTable:
		err = database.Debug().Model(&FactoryDepartment{}).Delete(&FactoryDepartment{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == FactoryLocationTable:
		err = database.Debug().Model(&FactoryLocation{}).Count(&numberOfRecords).Error
	case table == FactorySiteTable:
		err = database.Debug().Model(&FactorySite{}).Count(&numberOfRecords).Error
	case table == FactoryPlantTable:
		err = database.Debug().Model(&FactoryPlant{}).Count(&numberOfRecords).Error
	case table == FactoryDepartmentTable:
		err = database.Debug().Model(&FactoryDepartment{}).Count(&numberOfRecords).Error
	case table == FactoryDepartmentSectionTable:
		err = database.Debug().Model(&FactoryDepartmentSection{}).Count(&numberOfRecords).Error
	case table == FactoryUnitTable:
		err = database.Debug().Model(&FactoryUnit{}).Count(&numberOfRecords).Error
	case table == FactoryLevelTable:
		err = database.Debug().Model(&FactoryLevel{}).Count(&numberOfRecords).Error
	case table == FactoryAreaTable:
		err = database.Debug().Model(&FactoryArea{}).Count(&numberOfRecords).Error
	default:
		return -1
	}
	if err != nil {
		return -1
	}
	return numberOfRecords
}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == FactoryLocationTable:
		err = database.Debug().Model(&FactoryLocation{}).Take(&FactoryLocation{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactorySiteTable:
		err = database.Debug().Model(&FactorySite{}).Take(&FactorySite{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryPlantTable:
		err = database.Debug().Model(&FactoryPlant{}).Take(&FactoryPlant{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryDepartmentTable:
		err = database.Debug().Model(&FactoryDepartment{}).Take(&FactoryDepartment{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryComponentTable:
		err = database.Debug().Model(&FactoryComponent{}).Take(&FactoryComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryDepartmentSectionTable:
		err = database.Debug().Model(&FactoryDepartmentSection{}).Take(&FactoryDepartmentSection{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryGeneralFacilityTable:
		err = database.Debug().Model(&FactoryGeneralFacility{}).Take(&FactoryGeneralFacility{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryCustomerAssetTable:
		err = database.Debug().Model(&FactoryCustomerAsset{}).Take(&FactoryCustomerAsset{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryBuildingTable:
		err = database.Debug().Model(&FactoryBuilding{}).Take(&FactoryBuilding{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryUnitTable:
		err = database.Debug().Model(&FactoryUnit{}).Take(&FactoryUnit{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryLevelTable:
		err = database.Debug().Model(&FactoryLevel{}).Take(&FactoryLevel{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FactoryAreaTable:
		err = database.Debug().Model(&FactoryArea{}).Take(&FactoryArea{Id: recordId}).UpdateColumns(updateObject).Error

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
	case table == FactoryLocationTable:
		err = database.Debug().Model(&FactoryLocation{}).Take(&FactoryLocation{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FactorySiteTable:
		err = database.Debug().Model(&FactorySite{}).Take(&FactorySite{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FactoryPlantTable:
		err = database.Debug().Model(&FactoryPlant{}).Take(&FactoryPlant{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FactoryDepartmentTable:
		err = database.Debug().Model(&FactoryDepartment{}).Take(&FactoryDepartment{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FactoryDepartmentSectionTable:
		err = database.Debug().Model(&FactoryDepartmentSection{}).Take(&FactoryDepartmentSection{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *FactoryService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *FactoryService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
