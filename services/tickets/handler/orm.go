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
	case table == TicketsServiceRecordTrailTable:
		var dbObjects []TicketsServiceRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&TicketsServiceRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface TicketsServiceRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}
func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == TicketsServiceComponentTable:
		object := TicketsServiceComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == TicketsServiceRequestStatusTable:
		object := TicketsServiceRequestStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == TicketsServiceRequestCategoryTable:
		object := TicketsServiceRequestCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == TicketsServiceRequestSubCategoryTable:
		object := TicketsServiceRequestSubCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == TicketsServiceRequestTable:
		object := TicketsServiceRequest{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == TicketsServiceEmailTemplateTable:
		object := TicketsServiceEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id

	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == TicketsServiceRequestTable:
		dbObject := TicketsServiceRequest{Id: recordId}
		err = database.Debug().Model(&TicketsServiceRequest{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == TicketsServiceRequestStatusTable:
		dbObject := TicketsServiceRequestStatus{Id: recordId}
		err = database.Debug().Model(&TicketsServiceRequestStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == TicketsServiceRequestCategoryTable:
		dbObject := TicketsServiceRequestCategory{Id: recordId}
		err = database.Debug().Model(&TicketsServiceRequestCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == TicketsServiceRequestSubCategoryTable:
		dbObject := TicketsServiceRequestSubCategory{Id: recordId}
		err = database.Debug().Model(&TicketsServiceRequestSubCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == TicketsServiceComponentTable:
		dbObject := TicketsServiceComponent{Id: recordId}
		err = database.Debug().Model(&TicketsServiceComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == TicketsServiceEmailTemplateTable:
		dbObject := TicketsServiceEmailTemplate{Id: recordId}
		err = database.Debug().Model(&TicketsServiceEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == TicketsServiceComponentTable:

		var dbObjects []TicketsServiceComponent
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == TicketsServiceRequestTable:
		var dbObjects []TicketsServiceRequest
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceRequest{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequest{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceRequestCategoryTable:
		var dbObjects []TicketsServiceRequestCategory
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceRequestCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequestCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceEmailTemplateTable:
		var dbObjects []TicketsServiceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceEmailTemplateFieldTable:
		var dbObjects []TicketsServiceEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceEmailTemplateField{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceEmailTemplateField{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceRequestSubCategoryTable:
		var dbObjects []TicketsServiceRequestSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceRequestSubCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequestSubCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceRequestStatusTable:
		var dbObjects []TicketsServiceRequestStatus
		if len(objectCount) > 0 {
			err = database.Model(&TicketsServiceRequestStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequestStatus{}).Find(&dbObjects).Error
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
	case table == TicketsServiceRequestStatusTable:
		var dbObjects []TicketsServiceRequestStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceRequestStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequestStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == TicketsServiceRequestTable:
		var dbObjects []TicketsServiceRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceRequest{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequest{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceRequestCategoryTable:
		var dbObjects []TicketsServiceRequestCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceRequestCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequestCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceRequestSubCategoryTable:
		var dbObjects []TicketsServiceRequestSubCategory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceRequestSubCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&TicketsServiceRequestSubCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceEmailTemplateTable:
		var dbObjects []TicketsServiceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&TicketsServiceEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == TicketsServiceComponentTable:
		var dbObjects []TicketsServiceComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&TicketsServiceComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&TicketsServiceComponent{}).Where(condition).Find(&dbObjects).Error
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
	case table == TicketsServiceRequestStatusTable:
		err = database.Debug().Model(&TicketsServiceRequestStatus{}).Delete(&TicketsServiceRequestStatus{Id: objectInterface.Id}).Error
	case table == TicketsServiceRequestCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestCategory{}).Delete(&TicketsServiceRequestCategory{Id: objectInterface.Id}).Error
	case table == TicketsServiceRequestSubCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestSubCategory{}).Delete(&TicketsServiceRequestSubCategory{Id: objectInterface.Id}).Error
	case table == TicketsServiceRequestTable:
		err = database.Debug().Model(&TicketsServiceRequest{}).Delete(&TicketsServiceRequest{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == TicketsServiceRequestStatusTable:
		err = database.Debug().Model(&TicketsServiceRequestStatus{}).Count(&numberOfRecords).Error
	case table == TicketsServiceRequestCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestCategory{}).Count(&numberOfRecords).Error
	case table == TicketsServiceRequestSubCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestSubCategory{}).Count(&numberOfRecords).Error
	case table == TicketsServiceRequestTable:
		err = database.Debug().Model(&TicketsServiceRequest{}).Count(&numberOfRecords).Error
	case table == TicketsServiceEmailTemplateTable:
		err = database.Debug().Model(&TicketsServiceEmailTemplate{}).Count(&numberOfRecords).Error

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
	case table == TicketsServiceRequestStatusTable:
		err = database.Debug().Model(&TicketsServiceRequestStatus{}).Take(&TicketsServiceRequestStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == TicketsServiceRequestCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestCategory{}).Take(&TicketsServiceRequestCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == TicketsServiceRequestSubCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestSubCategory{}).Take(&TicketsServiceRequestSubCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == TicketsServiceRequestTable:
		err = database.Debug().Model(&TicketsServiceRequest{}).Take(&TicketsServiceRequest{Id: recordId}).UpdateColumns(updateObject).Error
	case table == TicketsServiceEmailTemplateTable:
		err = database.Debug().Model(&TicketsServiceEmailTemplate{}).Take(&TicketsServiceEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == TicketsServiceComponentTable:
		err = database.Debug().Model(&TicketsServiceComponent{}).Take(&TicketsServiceComponent{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == TicketsServiceRequestStatusTable:
		err = database.Debug().Model(&TicketsServiceRequestStatus{}).Take(&TicketsServiceRequestStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == TicketsServiceRequestCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestCategory{}).Take(&TicketsServiceRequestCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == TicketsServiceRequestSubCategoryTable:
		err = database.Debug().Model(&TicketsServiceRequestSubCategory{}).Take(&TicketsServiceRequestSubCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == TicketsServiceRequestTable:
		err = database.Debug().Model(&TicketsServiceRequest{}).Take(&TicketsServiceRequest{Id: objectInterface.Id}).UpdateColumns(updateObject).Error

	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (ts *TicketsService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := ts.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := ts.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(*listOfObjects)
					for _, referenceObject := range *listOfObjects {
						ts.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (ts *TicketsService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := ts.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := ts.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					for _, referenceObject := range *listOfObjects {
						fmt.Println("referenceTable : ", referenceTable, " id :", referenceObject)
						ArchiveObject(dbConnection, referenceTable, referenceObject)
						ts.CreateUserRecordMessage(ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
						ts.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
					}
				}
			}

		}
	}
}
