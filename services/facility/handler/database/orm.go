package database

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/facility/handler/const_util"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

var emptyObject interface{}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.FacilityServiceRecordTrailTable:
		var dbObjects []FacilityServiceRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.FacilityServiceRequestTable:
		var dbObjects []FacilityServiceRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceRequest{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceRequest{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface FacilityServiceRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}
func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == const_util.FacilityServiceComponentTable:
		object := FacilityServiceComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceRequestStatusTable:
		object := FacilityServiceRequestStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceRequestCategoryTable:
		object := FacilityServiceRequestCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceWorkflowEngineTable:
		object := FacilityServiceWorkflowEngine{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		object := FacilityServiceRequestSubCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceRequestTable:
		object := FacilityServiceRequest{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceEmailTemplateTable:
		object := FacilityServiceEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceCategoryTemplateTable:
		object := FacilityServiceCategoryTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.FacilityServiceAdminSettingTable:
		object := FacilityServiceAdminSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id

	default:
		return errors.New(const_util.CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == const_util.FacilityServiceRequestTable:
		dbObject := FacilityServiceRequest{Id: recordId}
		err = database.Debug().Model(&FacilityServiceRequest{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceRequestStatusTable:
		dbObject := FacilityServiceRequestStatus{Id: recordId}
		err = database.Debug().Model(&FacilityServiceRequestStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceRequestCategoryTable:
		dbObject := FacilityServiceRequestCategory{Id: recordId}
		err = database.Debug().Model(&FacilityServiceRequestCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceWorkflowEngineTable:
		dbObject := FacilityServiceWorkflowEngine{Id: recordId}
		err = database.Debug().Model(&FacilityServiceWorkflowEngine{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		dbObject := FacilityServiceRequestSubCategory{Id: recordId}
		err = database.Debug().Model(&FacilityServiceRequestSubCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceComponentTable:
		dbObject := FacilityServiceComponent{Id: recordId}
		err = database.Debug().Model(&FacilityServiceComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceEmailTemplateTable:
		dbObject := FacilityServiceEmailTemplate{Id: recordId}
		err = database.Debug().Model(&FacilityServiceEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceCategoryTemplateTable:
		dbObject := FacilityServiceCategoryTemplate{Id: recordId}
		err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceEmailTemplateFieldTable:
		dbObject := FacilityServiceEmailTemplateField{Id: recordId}
		err = database.Debug().Model(&FacilityServiceEmailTemplateField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceCategoryTemplateTable:
		dbObject := FacilityServiceCategoryTemplate{Id: recordId}
		err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.FacilityServiceAdminSettingTable:
		dbObject := FacilityServiceAdminSetting{Id: recordId}
		err = database.Debug().Model(&FacilityServiceAdminSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(const_util.GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.FacilityServiceComponentTable:

		var dbObjects []FacilityServiceComponent
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == const_util.FacilityServiceRequestTable:
		var dbObjects []FacilityServiceRequest
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceRequest{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequest{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceRequestCategoryTable:
		var dbObjects []FacilityServiceRequestCategory
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceRequestCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequestCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceWorkflowEngineTable:
		var dbObjects []FacilityServiceWorkflowEngine
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceWorkflowEngine{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceWorkflowEngine{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceEmailTemplateTable:
		var dbObjects []FacilityServiceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceCategoryTemplateTable:
		var dbObjects []FacilityServiceCategoryTemplate
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceCategoryTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceCategoryTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceEmailTemplateFieldTable:
		var dbObjects []FacilityServiceEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceEmailTemplateField{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceEmailTemplateField{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		var dbObjects []FacilityServiceRequestSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceRequestSubCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequestSubCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceRequestStatusTable:
		var dbObjects []FacilityServiceRequestStatus
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceRequestStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequestStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceCategoryTemplateTable:
		var dbObjects []FacilityServiceCategoryTemplate
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceCategoryTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceCategoryTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceAdminSettingTable:
		var dbObjects []FacilityServiceAdminSetting
		if len(objectCount) > 0 {
			err = database.Model(&FacilityServiceAdminSetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceAdminSetting{}).Find(&dbObjects).Error
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

func GetConditionalObjects(database *gorm.DB, table string, condition string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.FacilityServiceRequestStatusTable:
		var dbObjects []FacilityServiceRequestStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceRequestStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequestStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.FacilityServiceRequestTable:
		var dbObjects []FacilityServiceRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceRequest{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequest{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceRequestCategoryTable:
		var dbObjects []FacilityServiceRequestCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceRequestCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequestCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceWorkflowEngineTable:
		var dbObjects []FacilityServiceWorkflowEngine
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceWorkflowEngine{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceWorkflowEngine{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		var dbObjects []FacilityServiceRequestSubCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceRequestSubCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&FacilityServiceRequestSubCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceEmailTemplateTable:
		var dbObjects []FacilityServiceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceCategoryTemplateTable:
		var dbObjects []FacilityServiceCategoryTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceComponentTable:
		var dbObjects []FacilityServiceComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceCategoryTemplateTable:
		var dbObjects []FacilityServiceCategoryTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.FacilityServiceAdminSettingTable:
		var dbObjects []FacilityServiceAdminSetting
		if len(objectCount) > 0 {
			err = database.Debug().Model(&FacilityServiceAdminSetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&FacilityServiceAdminSetting{}).Where(condition).Find(&dbObjects).Error
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

func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	switch {
	case table == const_util.FacilityServiceRequestStatusTable:
		err = database.Debug().Model(&FacilityServiceRequestStatus{}).Delete(&FacilityServiceRequestStatus{Id: objectInterface.Id}).Error
	case table == const_util.FacilityServiceRequestCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestCategory{}).Delete(&FacilityServiceRequestCategory{Id: objectInterface.Id}).Error
	case table == const_util.FacilityServiceWorkflowEngineTable:
		err = database.Debug().Model(&FacilityServiceWorkflowEngine{}).Delete(&FacilityServiceWorkflowEngine{Id: objectInterface.Id}).Error
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestSubCategory{}).Delete(&FacilityServiceRequestSubCategory{Id: objectInterface.Id}).Error
	case table == const_util.FacilityServiceRequestTable:
		err = database.Debug().Model(&FacilityServiceRequest{}).Delete(&FacilityServiceRequest{Id: objectInterface.Id}).Error
	case table == const_util.FacilityServiceAdminSettingTable:
		err = database.Debug().Model(&FacilityServiceAdminSetting{}).Delete(&FacilityServiceAdminSetting{Id: objectInterface.Id}).Error
	default:
		return errors.New(const_util.GetUnknownObjectType)
	}

	return err

}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == const_util.FacilityServiceRequestStatusTable:
		err = database.Debug().Model(&FacilityServiceRequestStatus{}).Count(&numberOfRecords).Error
	case table == const_util.FacilityServiceRequestCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestCategory{}).Count(&numberOfRecords).Error
	case table == const_util.FacilityServiceWorkflowEngineTable:
		err = database.Debug().Model(&FacilityServiceWorkflowEngine{}).Count(&numberOfRecords).Error
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestSubCategory{}).Count(&numberOfRecords).Error
	case table == const_util.FacilityServiceRequestTable:
		err = database.Debug().Model(&FacilityServiceRequest{}).Count(&numberOfRecords).Error

	case table == const_util.FacilityServiceCategoryTemplateTable:
		err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Count(&numberOfRecords).Error
	case table == const_util.FacilityServiceAdminSettingTable:
		err = database.Debug().Model(&FacilityServiceAdminSetting{}).Count(&numberOfRecords).Error

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
	case table == const_util.FacilityServiceRequestStatusTable:
		err = database.Debug().Model(&FacilityServiceRequestStatus{}).Take(&FacilityServiceRequestStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceRequestCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestCategory{}).Take(&FacilityServiceRequestCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceWorkflowEngineTable:
		err = database.Debug().Model(&FacilityServiceWorkflowEngine{}).Take(&FacilityServiceWorkflowEngine{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestSubCategory{}).Take(&FacilityServiceRequestSubCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceRequestTable:
		err = database.Debug().Model(&FacilityServiceRequest{}).Take(&FacilityServiceRequest{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceEmailTemplateTable:
		err = database.Debug().Model(&FacilityServiceEmailTemplate{}).Take(&FacilityServiceEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceCategoryTemplateTable:
		err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Take(&FacilityServiceCategoryTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceComponentTable:
		err = database.Debug().Model(&FacilityServiceComponent{}).Take(&FacilityServiceComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceCategoryTemplateTable:
		err = database.Debug().Model(&FacilityServiceCategoryTemplate{}).Take(&FacilityServiceCategoryTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceAdminSettingTable:
		err = database.Debug().Model(&FacilityServiceAdminSetting{}).Take(&FacilityServiceAdminSetting{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == const_util.FacilityServiceRequestStatusTable:
		err = database.Debug().Model(&FacilityServiceRequestStatus{}).Take(&FacilityServiceRequestStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceRequestCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestCategory{}).Take(&FacilityServiceRequestCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceWorkflowEngineTable:
		err = database.Debug().Model(&FacilityServiceWorkflowEngine{}).Take(&FacilityServiceWorkflowEngine{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceRequestSubCategoryTable:
		err = database.Debug().Model(&FacilityServiceRequestSubCategory{}).Take(&FacilityServiceRequestSubCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceRequestTable:
		err = database.Debug().Model(&FacilityServiceRequest{}).Take(&FacilityServiceRequest{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.FacilityServiceAdminSettingTable:
		err = database.Debug().Model(&FacilityServiceAdminSetting{}).Take(&FacilityServiceAdminSetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(const_util.UpdateUnknownObjectType)
	}

	return err
}

func CountByCondition(database *gorm.DB, table, condition string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == const_util.FacilityServiceRequestTable:
		err = database.Debug().Model(&FacilityServiceRequest{}).Where(condition).Count(&numberOfRecords).Error
	default:
		return -1
	}

	if err != nil {
		return -1
	}
	return numberOfRecords
}

func GetLastObjects(database *gorm.DB, table string, condition string) (error, []component.GeneralObject) {
	var err error
	var dbObjects []component.GeneralObject

	err = database.Table(table).Where(condition).Last(&dbObjects).Error

	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}
}
