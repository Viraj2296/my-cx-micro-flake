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
	case table == WorkHourAssignmentRecordTrailTable:
		var dbObjects []WorkHourAssignmentRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&WorkHourAssignmentRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface WorkHourAssignmentRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == WorkHourAssignmentComponentTable:
		object := WorkHourAssignmentComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentRequestTable:
		object := WorkHourAssignmentRequest{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentTasksTable:
		object := WorkHourAssignmentTask{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentJRMasterTable:
		object := WorkHourAssignmentMasterJR{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentTLMasterTable:
		object := WorkHourAssignmentMasterTL{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentMRMasterTable:
		object := WorkHourAssignmentMasterMR{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentRequestStatusTable:
		object := WorkHourAssignmentRequestStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentEmailTemplateTable:
		object := WorkHourAssignmentEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		object := WorkHourAssignmentEmailTemplateField{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == WorkHourAssignmentComponentTable:
		dbObject := WorkHourAssignmentComponent{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentRequestTable:
		dbObject := WorkHourAssignmentRequest{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentRequest{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentTasksTable:
		dbObject := WorkHourAssignmentTask{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentTask{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentJRMasterTable:
		dbObject := WorkHourAssignmentMasterJR{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentMasterJR{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentTLMasterTable:
		dbObject := WorkHourAssignmentMasterTL{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentMasterTL{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentMRMasterTable:
		dbObject := WorkHourAssignmentMasterMR{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentMasterMR{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentRequestStatusTable:
		dbObject := WorkHourAssignmentRequestStatus{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentRequestStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentEmailTemplateTable:
		dbObject := WorkHourAssignmentEmailTemplate{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		dbObject := WorkHourAssignmentEmailTemplateField{Id: recordId}
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplateField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == WorkHourAssignmentComponentTable:

		var dbObjects []WorkHourAssignmentComponent
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentRequestTable:
		var dbObjects []WorkHourAssignmentRequest
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentRequest{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentRequest{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentTasksTable:
		var dbObjects []WorkHourAssignmentTask
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentTask{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentTask{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentJRMasterTable:
		var dbObjects []WorkHourAssignmentMasterJR
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentMasterJR{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentMasterJR{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentTLMasterTable:
		var dbObjects []WorkHourAssignmentMasterTL
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentMasterTL{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentMasterTL{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentMRMasterTable:
		var dbObjects []WorkHourAssignmentMasterMR
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentMasterMR{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentMasterMR{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentRequestStatusTable:
		var dbObjects []WorkHourAssignmentRequestStatus
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentRequestStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentRequestStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentEmailTemplateTable:
		var dbObjects []WorkHourAssignmentEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		var dbObjects []WorkHourAssignmentEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Model(&WorkHourAssignmentEmailTemplateField{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentEmailTemplateField{}).Find(&dbObjects).Error
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
	case table == WorkHourAssignmentComponentTable:
		var dbObjects []WorkHourAssignmentComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == WorkHourAssignmentRequestTable:
		var dbObjects []WorkHourAssignmentRequest
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentRequest{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentRequest{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == WorkHourAssignmentTasksTable:
		var dbObjects []WorkHourAssignmentTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentTask{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentTask{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == WorkHourAssignmentJRMasterTable:
		var dbObjects []WorkHourAssignmentMasterJR
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentMasterJR{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentMasterJR{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == WorkHourAssignmentTLMasterTable:
		var dbObjects []WorkHourAssignmentMasterTL
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentMasterTL{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentMasterTL{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentMRMasterTable:
		var dbObjects []WorkHourAssignmentMasterMR
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentMasterMR{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentMasterMR{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentRequestStatusTable:
		var dbObjects []WorkHourAssignmentRequestStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentRequestStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentRequestStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentEmailTemplateTable:
		var dbObjects []WorkHourAssignmentEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		var dbObjects []WorkHourAssignmentEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Debug().Model(&WorkHourAssignmentEmailTemplateField{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&WorkHourAssignmentEmailTemplateField{}).Where(condition).Find(&dbObjects).Error
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
	case table == WorkHourAssignmentComponentTable:
		err = database.Debug().Model(&WorkHourAssignmentComponent{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentRequestTable:
		err = database.Debug().Model(&WorkHourAssignmentRequest{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentTasksTable:
		err = database.Debug().Model(&WorkHourAssignmentTask{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentJRMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterJR{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentTLMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterTL{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentMRMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterMR{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentRequestStatusTable:
		err = database.Debug().Model(&WorkHourAssignmentRequestStatus{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentEmailTemplateTable:
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplate{}).Count(&numberOfRecords).Error
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplateField{}).Count(&numberOfRecords).Error
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
	case table == WorkHourAssignmentComponentTable:
		err = database.Debug().Model(&WorkHourAssignmentComponent{}).Take(&WorkHourAssignmentComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentRequestTable:
		err = database.Debug().Model(&WorkHourAssignmentRequest{}).Take(&WorkHourAssignmentRequest{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentTasksTable:
		err = database.Debug().Model(&WorkHourAssignmentTask{}).Take(&WorkHourAssignmentTask{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentJRMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterJR{}).Take(&WorkHourAssignmentMasterJR{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentTLMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterTL{}).Take(&WorkHourAssignmentMasterTL{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentMRMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterMR{}).Take(&WorkHourAssignmentMasterMR{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentRequestStatusTable:
		err = database.Debug().Model(&WorkHourAssignmentRequestStatus{}).Take(&WorkHourAssignmentRequestStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentEmailTemplateTable:
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplate{}).Take(&WorkHourAssignmentEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplateField{}).Take(&WorkHourAssignmentEmailTemplateField{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == WorkHourAssignmentComponentTable:
		err = database.Debug().Model(&WorkHourAssignmentComponent{}).Take(&WorkHourAssignmentComponent{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentRequestTable:
		err = database.Debug().Model(&WorkHourAssignmentRequest{}).Take(&WorkHourAssignmentRequest{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentTasksTable:
		err = database.Debug().Model(&WorkHourAssignmentTask{}).Take(&WorkHourAssignmentTask{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentJRMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterJR{}).Take(&WorkHourAssignmentMasterJR{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentTLMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterTL{}).Take(&WorkHourAssignmentMasterTL{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentMRMasterTable:
		err = database.Debug().Model(&WorkHourAssignmentMasterMR{}).Take(&WorkHourAssignmentMasterMR{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentRequestStatusTable:
		err = database.Debug().Model(&WorkHourAssignmentRequestStatus{}).Take(&WorkHourAssignmentRequestStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentEmailTemplateTable:
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplate{}).Take(&WorkHourAssignmentEmailTemplate{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == WorkHourAssignmentEmailTemplateFieldTable:
		err = database.Debug().Model(&WorkHourAssignmentEmailTemplateField{}).Take(&WorkHourAssignmentEmailTemplateField{Id: objectInterface.Id}).UpdateColumns(updateObject).Error

	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *WorkHourAssignmentService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *WorkHourAssignmentService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
