package database

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/it/handler/const_util"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

var emptyObject interface{}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.ITServiceRecordTrailTable:
		var dbObjects []ITServiceRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.ITServiceRequestTable:
		var dbObjects []ITServiceRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceRequest{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceRequest{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface ITServiceRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}
func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == const_util.ITServiceComponentTable:
		object := ITServiceComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceRequestStatusTable:
		object := ITServiceRequestStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceRequestCategoryTable:
		object := ITServiceRequestCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceRequestSubCategoryTable:
		object := ITServiceRequestSubCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceRequestTable:
		object := ITServiceRequest{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceEmailTemplateTable:
		object := ITServiceEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.IITServiceCategoryTemplateTable:
		object := ITServiceCategoryTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceSAPChangeReasonsTable:
		object := ITServiceSAPChangeReason{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceSAPAuthorizationFunctionsTable:
		object := ITServiceSAPAuthorizationFunction{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == const_util.ITServiceWorkflowEngineTable:
		object := ITServiceWorkflowEngine{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(const_util.CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == const_util.ITServiceRequestTable:
		dbObject := ITServiceRequest{Id: recordId}
		err = database.Debug().Model(&ITServiceRequest{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceRequestStatusTable:
		dbObject := ITServiceRequestStatus{Id: recordId}
		err = database.Debug().Model(&ITServiceRequestStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceRequestCategoryTable:
		dbObject := ITServiceRequestCategory{Id: recordId}
		err = database.Debug().Model(&ITServiceRequestCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceRequestSubCategoryTable:
		dbObject := ITServiceRequestSubCategory{Id: recordId}
		err = database.Debug().Model(&ITServiceRequestSubCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceComponentTable:
		dbObject := ITServiceComponent{Id: recordId}
		err = database.Debug().Model(&ITServiceComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceEmailTemplateTable:
		dbObject := ITServiceEmailTemplate{Id: recordId}
		err = database.Debug().Model(&ITServiceEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceEmailTemplateFieldTable:
		dbObject := ITServiceEmailTemplateField{Id: recordId}
		err = database.Debug().Model(&ITServiceEmailTemplateField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.IITServiceCategoryTemplateTable:
		dbObject := ITServiceCategoryTemplate{Id: recordId}
		err = database.Debug().Model(&ITServiceCategoryTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceSAPChangeReasonsTable:
		dbObject := ITServiceSAPChangeReason{Id: recordId}
		err = database.Debug().Model(&ITServiceSAPChangeReason{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceSAPAuthorizationFunctionsTable:
		dbObject := ITServiceSAPAuthorizationFunction{Id: recordId}
		err = database.Debug().Model(&ITServiceSAPAuthorizationFunction{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == const_util.ITServiceWorkflowEngineTable:
		dbObject := ITServiceWorkflowEngine{Id: recordId}
		err = database.Debug().Model(&ITServiceWorkflowEngine{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(const_util.GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == const_util.ITServiceComponentTable:

		var dbObjects []ITServiceComponent
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == const_util.ITServiceRequestTable:
		var dbObjects []ITServiceRequest
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceRequest{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequest{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceRequestCategoryTable:
		var dbObjects []ITServiceRequestCategory
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceRequestCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequestCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceEmailTemplateTable:
		var dbObjects []ITServiceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceEmailTemplateFieldTable:
		var dbObjects []ITServiceEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceEmailTemplateField{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceEmailTemplateField{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceRequestSubCategoryTable:
		var dbObjects []ITServiceRequestSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceRequestSubCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequestSubCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceRequestStatusTable:
		var dbObjects []ITServiceRequestStatus
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceRequestStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequestStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.IITServiceCategoryTemplateTable:
		var dbObjects []ITServiceCategoryTemplate
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceCategoryTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceCategoryTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceSAPChangeReasonsTable:
		var dbObjects []ITServiceSAPChangeReason
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceSAPChangeReason{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceSAPChangeReason{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceSAPAuthorizationFunctionsTable:
		var dbObjects []ITServiceSAPAuthorizationFunction
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceSAPAuthorizationFunction{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceSAPAuthorizationFunction{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceWorkflowEngineTable:
		var dbObjects []ITServiceWorkflowEngine
		if len(objectCount) > 0 {
			err = database.Model(&ITServiceWorkflowEngine{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceWorkflowEngine{}).Find(&dbObjects).Error
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
	case table == const_util.ITServiceRequestStatusTable:
		var dbObjects []ITServiceRequestStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceRequestStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequestStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == const_util.ITServiceRequestTable:
		var dbObjects []ITServiceRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceRequest{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequest{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceRequestCategoryTable:
		var dbObjects []ITServiceRequestCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceRequestCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequestCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceRequestSubCategoryTable:
		var dbObjects []ITServiceRequestSubCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceRequestSubCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ITServiceRequestSubCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceEmailTemplateTable:
		var dbObjects []ITServiceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceComponentTable:
		var dbObjects []ITServiceComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.IITServiceCategoryTemplateTable:
		var dbObjects []ITServiceCategoryTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceCategoryTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceCategoryTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceSAPChangeReasonsTable:
		var dbObjects []ITServiceSAPChangeReason
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceSAPChangeReason{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceSAPChangeReason{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceSAPAuthorizationFunctionsTable:
		var dbObjects []ITServiceSAPAuthorizationFunction
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceSAPAuthorizationFunction{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceSAPAuthorizationFunction{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == const_util.ITServiceWorkflowEngineTable:
		var dbObjects []ITServiceWorkflowEngine
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ITServiceWorkflowEngine{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ITServiceWorkflowEngine{}).Where(condition).Find(&dbObjects).Error
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
	case table == const_util.ITServiceRequestStatusTable:
		err = database.Debug().Model(&ITServiceRequestStatus{}).Delete(&ITServiceRequestStatus{Id: objectInterface.Id}).Error
	case table == const_util.ITServiceRequestCategoryTable:
		err = database.Debug().Model(&ITServiceRequestCategory{}).Delete(&ITServiceRequestCategory{Id: objectInterface.Id}).Error
	case table == const_util.ITServiceRequestSubCategoryTable:
		err = database.Debug().Model(&ITServiceRequestSubCategory{}).Delete(&ITServiceRequestSubCategory{Id: objectInterface.Id}).Error
	case table == const_util.ITServiceRequestTable:
		err = database.Debug().Model(&ITServiceRequest{}).Delete(&ITServiceRequest{Id: objectInterface.Id}).Error
	case table == const_util.ITServiceWorkflowEngineTable:
		err = database.Debug().Model(&ITServiceWorkflowEngine{}).Delete(&ITServiceWorkflowEngine{Id: objectInterface.Id}).Error
	case table == const_util.IITServiceCategoryTemplateTable:
		err = database.Debug().Model(&ITServiceCategoryTemplate{}).Delete(&ITServiceCategoryTemplate{Id: objectInterface.Id}).Error
	default:
		return errors.New(const_util.GetUnknownObjectType)
	}

	return err

}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == const_util.ITServiceRequestStatusTable:
		err = database.Debug().Model(&ITServiceRequestStatus{}).Count(&numberOfRecords).Error
	case table == const_util.ITServiceRequestCategoryTable:
		err = database.Debug().Model(&ITServiceRequestCategory{}).Count(&numberOfRecords).Error
	case table == const_util.ITServiceRequestSubCategoryTable:
		err = database.Debug().Model(&ITServiceRequestSubCategory{}).Count(&numberOfRecords).Error
	case table == const_util.ITServiceRequestTable:
		err = database.Debug().Model(&ITServiceRequest{}).Count(&numberOfRecords).Error
	case table == const_util.IITServiceCategoryTemplateTable:
		err = database.Debug().Model(&ITServiceCategoryTemplate{}).Count(&numberOfRecords).Error
	case table == const_util.ITServiceSAPChangeReasonsTable:
		err = database.Debug().Model(&ITServiceSAPChangeReason{}).Count(&numberOfRecords).Error
	case table == const_util.ITServiceSAPAuthorizationFunctionsTable:
		err = database.Debug().Model(&ITServiceSAPAuthorizationFunction{}).Count(&numberOfRecords).Error
	case table == const_util.ITServiceWorkflowEngineTable:
		err = database.Debug().Model(&ITServiceWorkflowEngine{}).Count(&numberOfRecords).Error
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
	case table == const_util.ITServiceRequestStatusTable:
		err = database.Debug().Model(&ITServiceRequestStatus{}).Take(&ITServiceRequestStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceRequestCategoryTable:
		err = database.Debug().Model(&ITServiceRequestCategory{}).Take(&ITServiceRequestCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceRequestSubCategoryTable:
		err = database.Debug().Model(&ITServiceRequestSubCategory{}).Take(&ITServiceRequestSubCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceRequestTable:
		err = database.Debug().Model(&ITServiceRequest{}).Take(&ITServiceRequest{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceEmailTemplateTable:
		err = database.Debug().Model(&ITServiceEmailTemplate{}).Take(&ITServiceEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceComponentTable:
		err = database.Debug().Model(&ITServiceComponent{}).Take(&ITServiceComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.IITServiceCategoryTemplateTable:
		err = database.Debug().Model(&ITServiceCategoryTemplate{}).Take(&ITServiceCategoryTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceSAPChangeReasonsTable:
		err = database.Debug().Model(&ITServiceSAPChangeReason{}).Take(&ITServiceSAPChangeReason{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceSAPAuthorizationFunctionsTable:
		err = database.Debug().Model(&ITServiceSAPAuthorizationFunction{}).Take(&ITServiceSAPAuthorizationFunction{Id: recordId}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceWorkflowEngineTable:
		err = database.Debug().Model(&ITServiceWorkflowEngine{}).Take(&ITServiceWorkflowEngine{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == const_util.ITServiceRequestStatusTable:
		err = database.Debug().Model(&ITServiceRequestStatus{}).Take(&ITServiceRequestStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceRequestCategoryTable:
		err = database.Debug().Model(&ITServiceRequestCategory{}).Take(&ITServiceRequestCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceRequestSubCategoryTable:
		err = database.Debug().Model(&ITServiceRequestSubCategory{}).Take(&ITServiceRequestSubCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.ITServiceRequestTable:
		err = database.Debug().Model(&ITServiceRequest{}).Take(&ITServiceRequest{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == const_util.IITServiceCategoryTemplateTable:
		err = database.Debug().Model(&ITServiceCategoryTemplate{}).Take(&ITServiceCategoryTemplate{Id: objectInterface.Id}).UpdateColumns(updateObject).Error

	default:
		err = errors.New(const_util.UpdateUnknownObjectType)
	}

	return err
}

func CountByCondition(database *gorm.DB, table, condition string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == const_util.ITServiceRequestTable:
		err = database.Debug().Model(&ITServiceRequest{}).Where(condition).Count(&numberOfRecords).Error
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
