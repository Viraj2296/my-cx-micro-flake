package database

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/moulds/const_util"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

var emptyObject interface{}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.MouldsRecordTrailTable:
		var dbObjects []MouldsRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldsRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldsRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.MouldTestRequestTable:
		var dbObjects []MouldTestRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldTestRequest{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldTestRequest{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.MouldModificationHistoryTable:
		var dbObjects []MouldModificationHistory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldModificationHistory{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldModificationHistory{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
		return nil, errors.New(const_util.GetUnknownObjectType)
	}
}
func CreateRecordTrail(database *gorm.DB, objectInterface MouldsRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == const_util.MouldMasterTable:
		object := MouldMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldStatusTable:
		object := MouldStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldCategoryTable:
		object := MouldCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldSubCategoryTable:
		object := MouldSubCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldTestRequestTable:
		object := MouldTestRequest{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.PartMasterTable:
		object := PartMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldEmailTemplateTable:
		object := MouldEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldTestStatusTable:
		object := MouldTestStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldManualShotCountTable:
		object := MouldManualShotCount{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldMachineMasterTable:
		object := MouldMachineMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldSettingTable:
		object := MouldSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.MouldShoutCountViewTable:
		object := MouldShotCountView{MouldId: objectInterface.Id, ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.MouldId
	case table == const_util.MouldModificationHistoryTable:
		object := MouldModificationHistory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(const_util.CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == const_util.MouldMasterTable:
		dbObject := MouldMaster{Id: recordId}
		err = database.Debug().Model(&MouldMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldStatusTable:
		dbObject := MouldStatus{Id: recordId}
		err = database.Debug().Model(&MouldStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldTestStatusTable:
		dbObject := MouldTestStatus{Id: recordId}
		err = database.Debug().Model(&MouldTestStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldCategoryTable:
		dbObject := MouldCategory{Id: recordId}
		err = database.Debug().Model(&MouldCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldSubCategoryTable:
		dbObject := MouldSubCategory{Id: recordId}
		err = database.Debug().Model(&MouldSubCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldComponentTable:
		dbObject := MouldComponent{Id: recordId}
		err = database.Debug().Model(&MouldComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldTestRequestTable:
		dbObject := MouldTestRequest{Id: recordId}
		err = database.Debug().Model(&MouldTestRequest{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.PartMasterTable:
		dbObject := PartMaster{Id: recordId}
		err = database.Debug().Model(&PartMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldEmailTemplateTable:
		dbObject := MouldEmailTemplate{Id: recordId}
		err = database.Debug().Model(&MouldEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldEmailTemplateFieldTable:
		dbObject := MouldEmailTemplateField{Id: recordId}
		err = database.Debug().Model(&MouldEmailTemplateField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldManualShotCountTable:
		dbObject := MouldManualShotCount{Id: recordId}
		err = database.Debug().Model(&MouldManualShotCount{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldMachineMasterTable:
		dbObject := MouldMachineMaster{Id: recordId}
		err = database.Debug().Model(&MouldMachineMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldSettingTable:
		dbObject := MouldSetting{Id: recordId}
		err = database.Debug().Model(&MouldSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.MouldShoutCountViewTable:
		dbObject := MouldShotCountView{MouldId: recordId}
		err = database.Debug().Model(&MouldShotCountView{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.MouldId, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(const_util.GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.MouldMasterTable:

		var dbObjects []MouldMaster
		if len(objectCount) > 0 {
			err = database.Model(&MouldMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldCategoryTable:
		var dbObjects []MouldCategory
		if len(objectCount) > 0 {
			err = database.Model(&MouldCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldSubCategoryTable:
		var dbObjects []MouldSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&MouldSubCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldSubCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldStatusTable:
		var dbObjects []MouldStatus
		if len(objectCount) > 0 {
			err = database.Model(&MouldStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldTestStatusTable:
		var dbObjects []MouldTestStatus
		if len(objectCount) > 0 {
			err = database.Model(&MouldTestStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldTestStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldComponentTable:
		var dbObjects []MouldComponent
		if len(objectCount) > 0 {
			err = database.Model(&MouldComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldTestRequestTable:
		var dbObjects []MouldTestRequest
		if len(objectCount) > 0 {
			err = database.Model(&MouldTestRequest{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldTestRequest{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.PartMasterTable:
		var dbObjects []PartMaster
		if len(objectCount) > 0 {
			err = database.Model(&PartMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&PartMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldEmailTemplateTable:
		var dbObjects []MouldEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&MouldEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldManualShotCountTable:
		var dbObjects []MouldManualShotCount
		if len(objectCount) > 0 {
			err = database.Model(&MouldManualShotCount{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldManualShotCount{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldMachineMasterTable:
		var dbObjects []MouldMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineMaster{}).Find(&dbObjects).Error
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
		return nil, errors.New(const_util.GetUnknownObjectType)
	}

}

func GetCount(database *gorm.DB, table string, condition ...string) int64 {
	var iNumberOfRecords int64
	switch {
	case table == const_util.MouldMasterTable:
		if len(condition) > 0 {
			database.Model(&MouldMaster{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&MouldMaster{}).Count(&iNumberOfRecords)
		}

	case table == const_util.MouldCategoryTable:
		if len(condition) > 0 {
			database.Model(&MouldCategory{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&MouldCategory{}).Count(&iNumberOfRecords)
		}

	case table == const_util.MouldSubCategoryTable:
		if len(condition) > 0 {
			database.Model(&MouldSubCategory{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&MouldSubCategory{}).Count(&iNumberOfRecords)
		}
	case table == const_util.MouldStatusTable:
		if len(condition) > 0 {
			database.Model(&MouldStatus{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&MouldStatus{}).Count(&iNumberOfRecords)
		}
	case table == const_util.MouldTestStatusTable:
		if len(condition) > 0 {
			database.Model(&MouldTestStatus{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&MouldTestStatus{}).Count(&iNumberOfRecords)
		}
	case table == const_util.MouldTestRequestTable:
		if len(condition) > 0 {
			database.Debug().Model(&MouldTestRequest{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Debug().Model(&MouldTestRequest{}).Count(&iNumberOfRecords)
		}
	case table == const_util.PartMasterTable:
		if len(condition) > 0 {
			database.Debug().Model(&PartMaster{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Debug().Model(&PartMaster{}).Count(&iNumberOfRecords)
		}
	case table == const_util.MouldEmailTemplateTable:
		if len(condition) > 0 {
			database.Debug().Model(&MouldEmailTemplate{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Debug().Model(&MouldEmailTemplate{}).Count(&iNumberOfRecords)
		}
	case table == const_util.MouldManualShotCountTable:
		if len(condition) > 0 {
			database.Debug().Model(&MouldManualShotCount{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Debug().Model(&MouldManualShotCount{}).Count(&iNumberOfRecords)
		}
	case table == const_util.MouldMachineMasterTable:
		if len(condition) > 0 {
			database.Debug().Model(&MouldMachineMaster{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Debug().Model(&MouldMachineMaster{}).Count(&iNumberOfRecords)
		}
	default:
		return 0
	}
	return iNumberOfRecords
}
func GetConditionalObjects(database *gorm.DB, table string, condition string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.MouldMasterTable:
		var dbObjects []MouldMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.MouldCategoryTable:
		var dbObjects []MouldCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.MouldSubCategoryTable:
		var dbObjects []MouldSubCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldSubCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldSubCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.MouldStatusTable:
		var dbObjects []MouldStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldTestStatusTable:
		var dbObjects []MouldTestStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldTestStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldTestStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldTestRequestTable:
		var dbObjects []MouldTestRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldTestRequest{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldTestRequest{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.PartMasterTable:
		var dbObjects []PartMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&PartMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&PartMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldEmailTemplateTable:
		var dbObjects []MouldEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldComponentTable:
		var dbObjects []MouldComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldManualShotCountTable:
		var dbObjects []MouldManualShotCount
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldManualShotCount{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldManualShotCount{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldMachineMasterTable:
		var dbObjects []MouldMachineMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMachineMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldSettingTable:
		var dbObjects []MouldSetting
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldSetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldSetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.MouldShoutCountViewTable:
		var dbObjects []MouldShotCountView
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldShotCountView{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldShotCountView{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.MouldId, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	default:
		return nil, errors.New(const_util.GetUnknownObjectType)
	}
}

func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	switch {
	case table == const_util.MouldMasterTable:
		err = database.Debug().Model(&MouldMaster{}).Delete(&MouldMaster{Id: objectInterface.Id}).Error
	case table == const_util.MouldStatusTable:
		err = database.Debug().Model(&MouldStatus{}).Delete(&MouldStatus{Id: objectInterface.Id}).Error
	case table == const_util.MouldCategoryTable:
		err = database.Debug().Model(&MouldCategory{}).Delete(&MouldCategory{Id: objectInterface.Id}).Error
	case table == const_util.MouldSubCategoryTable:
		err = database.Debug().Model(&MouldSubCategory{}).Delete(&MouldSubCategory{Id: objectInterface.Id}).Error
	case table == const_util.MouldTestRequestTable:
		err = database.Debug().Model(&MouldTestRequest{}).Delete(&MouldTestRequest{Id: objectInterface.Id}).Error
	case table == const_util.PartMasterTable:
		err = database.Debug().Model(&PartMaster{}).Delete(&PartMaster{Id: objectInterface.Id}).Error
	case table == const_util.MouldEmailTemplateTable:
		err = database.Debug().Model(&MouldEmailTemplate{}).Delete(&MouldEmailTemplate{Id: objectInterface.Id}).Error
	case table == const_util.MouldManualShotCountTable:
		err = database.Debug().Model(&MouldManualShotCount{}).Delete(&MouldManualShotCount{Id: objectInterface.Id}).Error
	case table == const_util.MouldMachineMasterTable:
		err = database.Debug().Model(&MouldMachineMaster{}).Delete(&MouldMachineMaster{Id: objectInterface.Id}).Error
	default:
		return errors.New(const_util.GetUnknownObjectType)
	}

	return err

}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == const_util.MouldMasterTable:
		err = database.Debug().Model(&MouldMaster{}).Count(&numberOfRecords).Error
	case table == const_util.MouldStatusTable:
		err = database.Debug().Model(&MouldStatus{}).Count(&numberOfRecords).Error
	case table == const_util.MouldCategoryTable:
		err = database.Debug().Model(&MouldCategory{}).Count(&numberOfRecords).Error
	case table == const_util.MouldSubCategoryTable:
		err = database.Debug().Model(&MouldSubCategory{}).Count(&numberOfRecords).Error
	case table == const_util.MouldTestRequestTable:
		err = database.Debug().Model(&MouldTestRequest{}).Count(&numberOfRecords).Error
	case table == const_util.PartMasterTable:
		err = database.Debug().Model(&PartMaster{}).Count(&numberOfRecords).Error
	case table == const_util.MouldEmailTemplateTable:
		err = database.Debug().Model(&MouldEmailTemplate{}).Count(&numberOfRecords).Error
	case table == const_util.MouldManualShotCountTable:
		err = database.Debug().Model(&MouldManualShotCount{}).Count(&numberOfRecords).Error
	case table == const_util.MouldMachineMasterTable:
		err = database.Debug().Model(&MouldMachineMaster{}).Count(&numberOfRecords).Error
	default:
		err = errors.New(const_util.UpdateUnknownObjectType)
	}

	if err != nil {
		return -1
	}
	return numberOfRecords
}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == const_util.MouldMasterTable:
		err = database.Debug().Model(&MouldMaster{}).Take(&MouldMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldStatusTable:
		err = database.Debug().Model(&MouldStatus{}).Take(&MouldStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldTestStatusTable:
		err = database.Debug().Model(&MouldTestStatus{}).Take(&MouldTestStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldCategoryTable:
		err = database.Debug().Model(&MouldCategory{}).Take(&MouldCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldSubCategoryTable:
		err = database.Debug().Model(&MouldSubCategory{}).Take(&MouldSubCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldTestRequestTable:
		err = database.Debug().Model(&MouldTestRequest{}).Take(&MouldTestRequest{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.PartMasterTable:
		err = database.Debug().Model(&PartMaster{}).Take(&PartMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldEmailTemplateTable:
		err = database.Debug().Model(&MouldEmailTemplate{}).Take(&MouldEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldComponentTable:
		err = database.Debug().Model(&MouldComponent{}).Take(&MouldComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldManualShotCountTable:
		err = database.Debug().Model(&MouldManualShotCount{}).Take(&MouldManualShotCount{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldMachineMasterTable:
		err = database.Debug().Model(&MouldMachineMaster{}).Take(&MouldMachineMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldShoutCountViewTable:
		err = database.Debug().Model(&MouldShotCountView{}).Take(&MouldShotCountView{MouldId: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.MouldSettingTable:
		err = database.Debug().Model(&MouldSetting{}).Take(&MouldSetting{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(const_util.UpdateUnknownObjectType)
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
	case table == const_util.MouldMasterTable:
		err = database.Debug().Model(&MouldMaster{}).Take(&MouldMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldStatusTable:
		err = database.Debug().Model(&MouldStatus{}).Take(&MouldStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldTestStatusTable:
		err = database.Debug().Model(&MouldTestStatus{}).Take(&MouldTestStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldCategoryTable:
		err = database.Debug().Model(&MouldCategory{}).Take(&MouldCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldSubCategoryTable:
		err = database.Debug().Model(&MouldSubCategory{}).Take(&MouldSubCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldTestRequestTable:
		err = database.Debug().Model(&MouldTestRequest{}).Take(&MouldTestRequest{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.PartMasterTable:
		err = database.Debug().Model(&PartMaster{}).Take(&PartMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldEmailTemplateTable:
		err = database.Debug().Model(&MouldEmailTemplate{}).Take(&MouldEmailTemplate{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldManualShotCountTable:
		err = database.Debug().Model(&MouldManualShotCount{}).Take(&MouldManualShotCount{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.MouldMachineMasterTable:
		err = database.Debug().Model(&MouldMachineMaster{}).Take(&MouldMachineMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(const_util.UpdateUnknownObjectType)
	}

	return err
}

func GetMouldShoutCountByID(database *gorm.DB, mouldID int) (*MouldShotCountView, error) {
	var mouldShoutCount MouldShotCountView

	err := database.Debug().
		Where("mould_id = ?", mouldID).
		First(&mouldShoutCount).Error

	if err != nil {
		return nil, err
	}

	return &mouldShoutCount, nil
}
