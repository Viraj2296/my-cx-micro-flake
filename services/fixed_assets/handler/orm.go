package handler

import (
	"cx-micro-flake/pkg/common/component"
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
	case table == FixedAssetRecordTrailTable:
		var dbObjects []FixedAssetRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FixedAssetRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface FixedAssetRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == FixedAssetsComponentTable:
		object := FixedAssetsComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetMasterTable:
		object := FixedAssetMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetContractTable:
		object := FixedAssetContract{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetStatusTable:
		object := FixedAssetStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetContractStatusTable:
		object := FixedAssetContractStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetDisposalTable:
		object := FixedAssetDisposal{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetTransferTable:
		object := FixedAssetTransfer{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetLocationTable:
		object := FixedAssetLocation{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetCategoryTable:
		object := FixedAssetCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssetClassTable:
		object := FixedAssetClass{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetEmailTemplateTable:
		object := FixedAssetEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetEmailTemplateFieldTable:
		object := FixedAssetEmailTemplateField{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetSetupMasterTable:
		object := FixedAssetSetupMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetTransferStatusTable:
		object := FixedAssetTransferStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == FixedAssetDisposalStatusTable:
		object := FixedAssetDisposalStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == FixedAssetMasterTable:
		dbObject := FixedAssetMaster{Id: recordId}
		err = database.Debug().Model(&FixedAssetMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetStatusTable:
		dbObject := FixedAssetStatus{Id: recordId}
		err = database.Debug().Model(&FixedAssetStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetCategoryTable:
		dbObject := FixedAssetCategory{Id: recordId}
		err = database.Debug().Model(&FixedAssetCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetSubCategoryTable:
		dbObject := FixedAssetSubCategory{Id: recordId}
		err = database.Debug().Model(&FixedAssetSubCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetClassTable:
		dbObject := FixedAssetClass{Id: recordId}
		err = database.Debug().Model(&FixedAssetClass{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetContractTable:
		dbObject := FixedAssetContract{Id: recordId}
		err = database.Debug().Model(&FixedAssetContract{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetContractStatusTable:
		dbObject := FixedAssetContractStatus{Id: recordId}
		err = database.Debug().Model(&FixedAssetContractStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetsComponentTable:
		dbObject := FixedAssetsComponent{Id: recordId}
		err = database.Debug().Model(&FixedAssetsComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetDisposalTable:
		dbObject := FixedAssetDisposal{Id: recordId}
		err = database.Debug().Model(&FixedAssetDisposal{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetTransferTable:
		dbObject := FixedAssetTransfer{Id: recordId}
		err = database.Debug().Model(&FixedAssetTransfer{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetLocationTable:
		dbObject := FixedAssetLocation{Id: recordId}
		err = database.Debug().Model(&FixedAssetLocation{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssetClassTable:
		dbObject := FixedAssetClass{Id: recordId}
		err = database.Debug().Model(&FixedAssetClass{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetDynamicFieldTable:
		dbObject := FixedAssetDynamicField{Id: recordId}
		err = database.Debug().Model(&FixedAssetDynamicField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetDynamicFieldConfigurationTable:
		dbObject := FixedAssetDynamicFieldConfiguration{Id: recordId}
		err = database.Debug().Model(&FixedAssetDynamicFieldConfiguration{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetEmailTemplateTable:
		dbObject := FixedAssetEmailTemplate{Id: recordId}
		err = database.Debug().Model(&FixedAssetEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetEmailTemplateFieldTable:
		dbObject := FixedAssetEmailTemplateField{Id: recordId}
		err = database.Debug().Model(&FixedAssetEmailTemplateField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetSetupMasterTable:
		dbObject := FixedAssetSetupMaster{Id: recordId}
		err = database.Debug().Model(&FixedAssetSetupMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetTransferStatusTable:
		dbObject := FixedAssetTransferStatus{Id: recordId}
		err = database.Debug().Model(&FixedAssetTransferStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == FixedAssetDisposalStatusTable:
		dbObject := FixedAssetDisposalStatus{Id: recordId}
		err = database.Debug().Model(&FixedAssetDisposalStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == FixedAssetsComponentTable:

		var dbObjects []FixedAssetsComponent
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetsComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetsComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == FixedAssetMasterTable:
		var dbObjects []FixedAssetMaster
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetContractTable:
		var dbObjects []FixedAssetContract
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetContract{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetContract{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetContractStatusTable:
		var dbObjects []FixedAssetContractStatus
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetContractStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetContractStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetStatusTable:
		var dbObjects []FixedAssetStatus
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetCategoryTable:
		var dbObjects []FixedAssetCategory
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetSubCategoryTable:
		var dbObjects []FixedAssetSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetSubCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetSubCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == FixedAssetDisposalTable:
		var dbObjects []FixedAssetDisposal
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetDisposal{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FixedAssetDisposal{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetTransferTable:
		var dbObjects []FixedAssetTransfer
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetTransfer{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetTransfer{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetLocationTable:
		var dbObjects []FixedAssetLocation
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetLocation{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetLocation{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetClassTable:
		var dbObjects []FixedAssetClass
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetClass{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetClass{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetSetupMasterTable:
		var dbObjects []FixedAssetSetupMaster
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetSetupMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetSetupMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetTransferStatusTable:
		var dbObjects []FixedAssetTransferStatus
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetTransferStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetTransferStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == FixedAssetDisposalStatusTable:
		var dbObjects []FixedAssetDisposalStatus
		if len(objectCount) > 0 {
			err = database.Model(&FixedAssetDisposalStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetDisposalStatus{}).Find(&dbObjects).Error
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
	case table == FixedAssetMasterTable:
		var dbObjects []FixedAssetMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetContractTable:
		var dbObjects []FixedAssetContract
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetContract{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetContract{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetContractStatusTable:
		var dbObjects []FixedAssetContractStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetContractStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetContractStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetStatusTable:
		var dbObjects []FixedAssetStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetCategoryTable:
		var dbObjects []FixedAssetCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetSubCategoryTable:
		var dbObjects []FixedAssetSubCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetSubCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetSubCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetTransferTable:
		var dbObjects []FixedAssetTransfer
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetTransfer{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetTransfer{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetDisposalTable:
		var dbObjects []FixedAssetDisposal
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetDisposal{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetDisposal{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetLocationTable:
		var dbObjects []FixedAssetLocation
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetLocation{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetLocation{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetDynamicFieldTable:
		var dbObjects []FixedAssetDynamicField
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetDynamicField{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetDynamicField{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetClassTable:
		var dbObjects []FixedAssetClass
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetClass{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetClass{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetsComponentTable:
		var dbObjects []FixedAssetsComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetsComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetsComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetEmailTemplateTable:
		var dbObjects []FixedAssetEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetEmailTemplateFieldTable:
		var dbObjects []FixedAssetEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetEmailTemplateField{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetEmailTemplateField{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetSetupMasterTable:
		var dbObjects []FixedAssetSetupMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetSetupMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetSetupMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetTransferStatusTable:
		var dbObjects []FixedAssetTransferStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetTransferStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetTransferStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == FixedAssetDisposalStatusTable:
		var dbObjects []FixedAssetDisposalStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FixedAssetDisposalStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FixedAssetDisposalStatus{}).Where(condition).Find(&dbObjects).Error
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
	case table == FixedAssetMasterTable:
		err = database.Debug().Model(&FixedAssetMaster{}).Delete(&FixedAssetMaster{Id: objectInterface.Id}).Error
	case table == FixedAssetContractTable:
		err = database.Debug().Model(&FixedAssetContract{}).Delete(&FixedAssetContract{Id: objectInterface.Id}).Error
	case table == FixedAssetContractStatusTable:
		err = database.Debug().Model(&FixedAssetContractStatus{}).Delete(&FixedAssetContractStatus{Id: objectInterface.Id}).Error
	case table == FixedAssetStatusTable:
		err = database.Debug().Model(&FixedAssetStatus{}).Delete(&FixedAssetStatus{Id: objectInterface.Id}).Error
	case table == FixedAssetTransferStatusTable:
		err = database.Debug().Model(&FixedAssetTransferStatus{}).Delete(&FixedAssetTransferStatus{Id: objectInterface.Id}).Error
	case table == FixedAssetSetupMasterTable:
		err = database.Debug().Model(&FixedAssetSetupMaster{}).Delete(&FixedAssetSetupMaster{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == FixedAssetMasterTable:
		err = database.Debug().Model(&FixedAssetMaster{}).Count(&numberOfRecords).Error
	case table == FixedAssetContractTable:
		err = database.Debug().Model(&FixedAssetContract{}).Count(&numberOfRecords).Error
	case table == FixedAssetStatusTable:
		err = database.Debug().Model(&FixedAssetStatus{}).Count(&numberOfRecords).Error
	case table == FixedAssetCategoryTable:
		err = database.Debug().Model(&FixedAssetCategory{}).Count(&numberOfRecords).Error
	case table == FixedAssetSubCategoryTable:
		err = database.Debug().Model(&FixedAssetSubCategory{}).Count(&numberOfRecords).Error
	case table == FixedAssetClassTable:
		err = database.Debug().Model(&FixedAssetClass{}).Count(&numberOfRecords).Error
	case table == FixedAssetDisposalTable:
		err = database.Debug().Model(&FixedAssetDisposal{}).Count(&numberOfRecords).Error
	case table == FixedAssetTransferTable:
		err = database.Debug().Model(&FixedAssetTransfer{}).Count(&numberOfRecords).Error
	case table == FixedAssetLocationTable:
		err = database.Debug().Model(&FixedAssetLocation{}).Count(&numberOfRecords).Error
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
	case table == FixedAssetMasterTable:
		err = database.Debug().Model(&FixedAssetMaster{}).Take(&FixedAssetMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetContractTable:
		err = database.Debug().Model(&FixedAssetContract{}).Take(&FixedAssetContract{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetCategoryTable:
		err = database.Debug().Model(&FixedAssetCategory{}).Take(&FixedAssetCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetStatusTable:
		err = database.Debug().Model(&FixedAssetStatus{}).Take(&FixedAssetStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetSubCategoryTable:
		err = database.Debug().Model(&FixedAssetSubCategory{}).Take(&FixedAssetSubCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetClassTable:
		err = database.Debug().Model(&FixedAssetClass{}).Take(&FixedAssetClass{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetDisposalTable:
		err = database.Debug().Model(&FixedAssetDisposal{}).Take(&FixedAssetDisposal{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetTransferTable:
		err = database.Debug().Model(&FixedAssetTransfer{}).Take(&FixedAssetTransfer{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetLocationTable:
		err = database.Debug().Model(&FixedAssetLocation{}).Take(&FixedAssetLocation{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetsComponentTable:
		err = database.Debug().Model(&FixedAssetsComponent{}).Take(&FixedAssetsComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == FixedAssetEmailTemplateTable:
		err = database.Debug().Model(&FixedAssetEmailTemplate{}).Take(&FixedAssetEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == FixedAssetMasterTable:
		err = database.Debug().Model(&FixedAssetMaster{}).Take(&FixedAssetMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetCategoryTable:
		err = database.Debug().Model(&FixedAssetCategory{}).Take(&FixedAssetCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetContractTable:
		err = database.Debug().Model(&FixedAssetContract{}).Take(&FixedAssetContract{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetStatusTable:
		err = database.Debug().Model(&FixedAssetStatus{}).Take(&FixedAssetStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetSubCategoryTable:
		err = database.Debug().Model(&FixedAssetSubCategory{}).Take(&FixedAssetSubCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetClassTable:
		err = database.Debug().Model(&FixedAssetClass{}).Take(&FixedAssetClass{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetDisposalTable:
		err = database.Debug().Model(&FixedAssetDisposal{}).Take(&FixedAssetDisposal{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetTransferTable:
		err = database.Debug().Model(&FixedAssetTransfer{}).Take(&FixedAssetTransfer{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == FixedAssetLocationTable:
		err = database.Debug().Model(&FixedAssetLocation{}).Take(&FixedAssetLocation{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *FixedAssetService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *FixedAssetService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
