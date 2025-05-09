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

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == ToolingProjectTable:
		object := ToolingProject{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingProjectTaskTable:
		object := ToolingProjectTask{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingProjectStatusTable:
		object := ToolingProjectStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingProjectTaskStatusTable:
		object := ToolingProjectTaskStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingEmailTemplateTable:
		object := ToolingEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingProjectSprintTable:
		object := ToolingProjectSprint{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingProgramTable:
		object := ToolingProgram{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingStatusTable:
		object := ToolingStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingLocationTable:
		object := ToolingLocation{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingMouldVendorTable:
		object := ToolingMouldVendor{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingProjectTaskListTable:
		object := ToolingProjectTaskCheckList{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingGatingTypeTab:
		object := ToolingGatingType{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingHotRunnerBrandTable:
		object := ToolingHotRunnerBrand{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingRunnerTypeTab:
		object := ToolingRunnerType{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingHotRunnerConnectorTypeTable:
		object := ToolingHotRunnerConnectorType{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingHotRunnerControllerTable:
		object := ToolingHotRunnerController{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingHRControllerTable:
		object := ToolingHRController{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingHotRunnerControllerBrandTable:
		object := ToolingHotRunnerControllerBrand{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingCoolingFittingTable:
		object := ToolingCoolingFitting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}

}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == ToolingProjectTable:
		dbObject := ToolingProject{Id: recordId}
		err = database.Debug().Model(&ToolingProject{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingProjectTaskTable:
		dbObject := ToolingProjectTask{Id: recordId}
		err = database.Debug().Model(&ToolingProjectTask{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingProjectStatusTable:
		dbObject := ToolingProjectStatus{Id: recordId}
		err = database.Debug().Model(&ToolingProjectStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingProjectTaskStatusTable:
		dbObject := ToolingProjectTaskStatus{Id: recordId}
		err = database.Debug().Model(&ToolingProjectTaskStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingEmailTemplateTable:
		dbObject := ToolingEmailTemplate{Id: recordId}
		err = database.Debug().Model(&ToolingEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingProjectSprintTable:
		dbObject := ToolingProjectSprint{Id: recordId}
		err = database.Debug().Model(&ToolingProjectSprint{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingProgramTable:
		dbObject := ToolingProgram{Id: recordId}
		err = database.Debug().Model(&ToolingProgram{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingStatusTable:
		dbObject := ToolingStatus{Id: recordId}
		err = database.Debug().Model(&ToolingStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingLocationTable:
		dbObject := ToolingLocation{Id: recordId}
		err = database.Debug().Model(&ToolingLocation{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingMouldVendorTable:
		dbObject := ToolingMouldVendor{Id: recordId}
		err = database.Debug().Model(&ToolingMouldVendor{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingProjectTaskListTable:
		dbObject := ToolingProjectTaskCheckList{Id: recordId}
		err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingGatingTypeTab:
		dbObject := ToolingGatingType{Id: recordId}
		err = database.Debug().Model(&ToolingGatingType{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingHotRunnerBrandTable:
		dbObject := ToolingHotRunnerBrand{Id: recordId}
		err = database.Debug().Model(&ToolingHotRunnerBrand{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingRunnerTypeTab:
		dbObject := ToolingRunnerType{Id: recordId}
		err = database.Debug().Model(&ToolingRunnerType{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingHotRunnerConnectorTypeTable:
		dbObject := ToolingHotRunnerConnectorType{Id: recordId}
		err = database.Debug().Model(&ToolingHotRunnerConnectorType{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingHotRunnerControllerTable:
		dbObject := ToolingHotRunnerController{Id: recordId}
		err = database.Debug().Model(&ToolingHotRunnerController{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingHRControllerTable:
		dbObject := ToolingHRController{Id: recordId}
		err = database.Debug().Model(&ToolingHRController{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingHotRunnerControllerBrandTable:
		dbObject := ToolingHotRunnerControllerBrand{Id: recordId}
		err = database.Debug().Model(&ToolingHotRunnerControllerBrand{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingCoolingFittingTable:
		dbObject := ToolingCoolingFitting{Id: recordId}
		err = database.Debug().Model(&ToolingCoolingFitting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}
}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == ToolingComponentTable:

		var dbObjects []ToolingProject
		if len(objectCount) > 0 {
			err = database.Model(&ToolingComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectTable:

		var dbObjects []ToolingProject
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProject{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProject{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectTaskTable:

		var dbObjects []ToolingProjectTask
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProjectTask{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProjectTask{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectStatusTable:

		var dbObjects []ToolingProjectStatus
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProjectStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProjectStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectTaskStatusTable:

		var dbObjects []ToolingProjectTaskStatus
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProjectTaskStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProjectTaskStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingEmailTemplateTable:
		var dbObjects []ToolingEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&ToolingEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingEmailTemplateFieldTable:
		var dbObjects []ToolingEmailTemplateField
		if len(objectCount) > 0 {
			err = database.Model(&ToolingEmailTemplateField{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingEmailTemplateField{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectSprintTable:
		var dbObjects []ToolingProjectSprint
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProjectSprint{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProjectSprint{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProgramTable:
		var dbObjects []ToolingProgram
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProgram{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProgram{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingStatusTable:
		var dbObjects []ToolingStatus
		if len(objectCount) > 0 {
			err = database.Model(&ToolingStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingLocationTable:
		var dbObjects []ToolingLocation
		if len(objectCount) > 0 {
			err = database.Model(&ToolingLocation{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingLocation{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMouldVendorTable:
		var dbObjects []ToolingMouldVendor
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMouldVendor{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMouldVendor{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectTaskListTable:
		var dbObjects []ToolingProjectTaskCheckList
		if len(objectCount) > 0 {
			err = database.Model(&ToolingProjectTaskCheckList{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingProjectTaskCheckList{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingGatingTypeTab:
		var dbObjects []ToolingGatingType
		if len(objectCount) > 0 {
			err = database.Model(&ToolingGatingType{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingGatingType{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerBrandTable:
		var dbObjects []ToolingHotRunnerBrand
		if len(objectCount) > 0 {
			err = database.Model(&ToolingHotRunnerBrand{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingHotRunnerBrand{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingRunnerTypeTab:
		var dbObjects []ToolingRunnerType
		if len(objectCount) > 0 {
			err = database.Model(&ToolingRunnerType{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingRunnerType{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerConnectorTypeTable:
		var dbObjects []ToolingHotRunnerConnectorType
		if len(objectCount) > 0 {
			err = database.Model(&ToolingHotRunnerConnectorType{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingHotRunnerConnectorType{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerControllerTable:
		var dbObjects []ToolingHotRunnerController
		if len(objectCount) > 0 {
			err = database.Model(&ToolingHotRunnerController{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingHotRunnerController{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHRControllerTable:
		var dbObjects []ToolingHRController
		if len(objectCount) > 0 {
			err = database.Model(&ToolingHRController{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingHRController{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerControllerBrandTable:
		var dbObjects []ToolingHotRunnerControllerBrand
		if len(objectCount) > 0 {
			err = database.Model(&ToolingHotRunnerControllerBrand{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingHotRunnerControllerBrand{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingCoolingFittingTable:
		var dbObjects []ToolingCoolingFitting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingCoolingFitting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingCoolingFitting{}).Find(&dbObjects).Error
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
	case table == ToolingProjectTable:
		var dbObjects []ToolingProject
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProject{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProject{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingProjectTaskTable:
		var dbObjects []ToolingProjectTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProjectTask{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProjectTask{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingProjectStatusTable:
		var dbObjects []ToolingProjectStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProjectStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProjectStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingProjectTaskStatusTable:
		var dbObjects []ToolingProjectTaskStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProjectTaskStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProjectTaskStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingEmailTemplateTable:
		var dbObjects []ToolingEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectSprintTable:
		var dbObjects []ToolingProjectSprint
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProjectSprint{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProjectSprint{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingComponentTable:
		var dbObjects []ToolingComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProgramTable:
		var dbObjects []ToolingProgram
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProgram{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProgram{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingStatusTable:
		var dbObjects []ToolingStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingLocationTable:
		var dbObjects []ToolingLocation
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingLocation{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingLocation{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMouldVendorTable:
		var dbObjects []ToolingMouldVendor
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingMouldVendor{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingMouldVendor{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingProjectTaskListTable:
		var dbObjects []ToolingProjectTaskCheckList
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingGatingTypeTab:
		var dbObjects []ToolingGatingType
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingGatingType{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingGatingType{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerBrandTable:
		var dbObjects []ToolingHotRunnerBrand
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingHotRunnerBrand{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingHotRunnerBrand{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingRunnerTypeTab:
		var dbObjects []ToolingRunnerType
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingRunnerType{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingRunnerType{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerConnectorTypeTable:
		var dbObjects []ToolingHotRunnerConnectorType
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingHotRunnerConnectorType{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingHotRunnerConnectorType{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerControllerTable:
		var dbObjects []ToolingHotRunnerController
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingHotRunnerController{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingHotRunnerController{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHRControllerTable:
		var dbObjects []ToolingHRController
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingHRController{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingHRController{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingHotRunnerControllerBrandTable:
		var dbObjects []ToolingHotRunnerControllerBrand
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingHotRunnerControllerBrand{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingHotRunnerControllerBrand{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingCoolingFittingTable:
		var dbObjects []ToolingCoolingFitting
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingCoolingFitting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingCoolingFitting{}).Where(condition).Find(&dbObjects).Error
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
	case table == ToolingProjectTable:
		err = database.Debug().Model(&ToolingProject{}).Delete(&ToolingProject{Id: objectInterface.Id}).Error
	case table == ToolingProjectTaskTable:
		err = database.Debug().Model(&ToolingProjectTask{}).Delete(&ToolingProjectTask{Id: objectInterface.Id}).Error
	case table == ToolingProjectStatusTable:
		err = database.Debug().Model(&ToolingProjectStatus{}).Delete(&ToolingProjectStatus{Id: objectInterface.Id}).Error
	case table == ToolingProjectTaskStatusTable:
		err = database.Debug().Model(&ToolingProjectTaskStatus{}).Delete(&ToolingProjectTaskStatus{Id: objectInterface.Id}).Error
	case table == ToolingProjectSprintTable:
		err = database.Debug().Model(&ToolingProjectSprint{}).Delete(&ToolingProjectSprint{Id: objectInterface.Id}).Error
	case table == ToolingProgramTable:
		err = database.Debug().Model(&ToolingProgram{}).Delete(&ToolingProgram{Id: objectInterface.Id}).Error
	case table == ToolingStatusTable:
		err = database.Debug().Model(&ToolingStatus{}).Delete(&ToolingStatus{Id: objectInterface.Id}).Error
	case table == ToolingLocationTable:
		err = database.Debug().Model(&ToolingLocation{}).Delete(&ToolingLocation{Id: objectInterface.Id}).Error
	case table == ToolingMouldVendorTable:
		err = database.Debug().Model(&ToolingMouldVendor{}).Delete(&ToolingMouldVendor{Id: objectInterface.Id}).Error
	case table == ToolingProjectTaskListTable:
		err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Delete(&ToolingProjectTaskCheckList{Id: objectInterface.Id}).Error
	case table == ToolingCoolingFittingTable:
		err = database.Debug().Model(&ToolingCoolingFitting{}).Delete(&ToolingCoolingFitting{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err
}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == ToolingProjectTable:
		err = database.Debug().Model(&ToolingProject{}).Take(&ToolingProject{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingProjectTaskTable:
		err = database.Debug().Model(&ToolingProjectTask{}).Take(&ToolingProjectTask{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingProjectStatusTable:
		err = database.Debug().Model(&ToolingProjectStatus{}).Take(&ToolingProjectStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingProjectTaskStatusTable:
		err = database.Debug().Model(&ToolingProjectTaskStatus{}).Take(&ToolingProjectTaskStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingEmailTemplateTable:
		err = database.Debug().Model(&ToolingEmailTemplate{}).Take(&ToolingEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingProjectSprintTable:
		err = database.Debug().Model(&ToolingProjectSprint{}).Take(&ToolingProjectSprint{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingComponentTable:
		err = database.Debug().Model(&ToolingComponent{}).Take(&ToolingComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingProgramTable:
		err = database.Debug().Model(&ToolingProgram{}).Take(&ToolingProgram{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingStatusTable:
		err = database.Debug().Model(&ToolingStatus{}).Take(&ToolingStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingLocationTable:
		err = database.Debug().Model(&ToolingLocation{}).Take(&ToolingLocation{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingMouldVendorTable:
		err = database.Debug().Model(&ToolingMouldVendor{}).Take(&ToolingMouldVendor{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingProjectTaskListTable:
		err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Take(&ToolingProjectTaskCheckList{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingGatingTypeTab:
		err = database.Debug().Model(&ToolingGatingType{}).Take(&ToolingGatingType{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerBrandTable:
		err = database.Debug().Model(&ToolingHotRunnerBrand{}).Take(&ToolingHotRunnerBrand{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingRunnerTypeTab:
		err = database.Debug().Model(&ToolingRunnerType{}).Take(&ToolingRunnerType{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerConnectorTypeTable:
		err = database.Debug().Model(&ToolingHotRunnerConnectorType{}).Take(&ToolingHotRunnerConnectorType{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerControllerTable:
		err = database.Debug().Model(&ToolingHotRunnerController{}).Take(&ToolingHotRunnerController{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingHRControllerTable:
		err = database.Debug().Model(&ToolingHRController{}).Take(&ToolingHRController{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerControllerBrandTable:
		err = database.Debug().Model(&ToolingHotRunnerControllerBrand{}).Take(&ToolingHotRunnerControllerBrand{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingCoolingFittingTable:
		err = database.Debug().Model(&ToolingCoolingFitting{}).Take(&ToolingCoolingFitting{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == ToolingProjectTable:
		err = database.Debug().Model(&ToolingProject{}).Take(&ToolingProject{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingProjectTaskTable:
		err = database.Debug().Model(&ToolingProjectTask{}).Take(&ToolingProjectTask{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingProjectStatusTable:
		err = database.Debug().Model(&ToolingProjectStatus{}).Take(&ToolingProjectStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingProjectTaskStatusTable:
		err = database.Debug().Model(&ToolingProjectTaskStatus{}).Take(&ToolingProjectTaskStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingProjectSprintTable:
		err = database.Debug().Model(&ToolingProjectSprint{}).Take(&ToolingProjectSprint{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingProgramTable:
		err = database.Debug().Model(&ToolingProgram{}).Take(&ToolingProgram{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingStatusTable:
		err = database.Debug().Model(&ToolingStatus{}).Take(&ToolingStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingLocationTable:
		err = database.Debug().Model(&ToolingLocation{}).Take(&ToolingLocation{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingMouldVendorTable:
		err = database.Debug().Model(&ToolingMouldVendor{}).Take(&ToolingMouldVendor{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingProjectTaskListTable:
		err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Take(&ToolingProjectTaskCheckList{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingGatingTypeTab:
		err = database.Debug().Model(&ToolingGatingType{}).Take(&ToolingGatingType{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerBrandTable:
		err = database.Debug().Model(&ToolingHotRunnerBrand{}).Take(&ToolingHotRunnerBrand{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingRunnerTypeTab:
		err = database.Debug().Model(&ToolingRunnerType{}).Take(&ToolingRunnerType{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerConnectorTypeTable:
		err = database.Debug().Model(&ToolingHotRunnerConnectorType{}).Take(&ToolingHotRunnerConnectorType{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerControllerTable:
		err = database.Debug().Model(&ToolingHotRunnerController{}).Take(&ToolingHotRunnerController{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingHRControllerTable:
		err = database.Debug().Model(&ToolingHRController{}).Take(&ToolingHRController{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingHotRunnerControllerBrandTable:
		err = database.Debug().Model(&ToolingHotRunnerControllerBrand{}).Take(&ToolingHotRunnerControllerBrand{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingCoolingFittingTable:
		err = database.Debug().Model(&ToolingCoolingFitting{}).Take(&ToolingCoolingFitting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}
	return err
}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == ToolingProjectTable:
		err = database.Debug().Model(&ToolingProject{}).Count(&numberOfRecords).Error
	case table == ToolingProjectTaskTable:
		err = database.Debug().Model(&ToolingProjectTask{}).Count(&numberOfRecords).Error
	case table == ToolingProjectStatusTable:
		err = database.Debug().Model(&ToolingProjectStatus{}).Count(&numberOfRecords).Error
	case table == ToolingProjectTaskStatusTable:
		err = database.Debug().Model(&ToolingProjectTaskStatus{}).Count(&numberOfRecords).Error
	case table == ToolingProgramTable:
		err = database.Debug().Model(&ToolingProgram{}).Count(&numberOfRecords).Error
	case table == ToolingStatusTable:
		err = database.Debug().Model(&ToolingStatus{}).Count(&numberOfRecords).Error
	case table == ToolingLocationTable:
		err = database.Debug().Model(&ToolingLocation{}).Count(&numberOfRecords).Error
	case table == ToolingMouldVendorTable:
		err = database.Debug().Model(&ToolingMouldVendor{}).Count(&numberOfRecords).Error
	case table == ToolingProjectTaskListTable:
		err = database.Debug().Model(&ToolingProjectTaskCheckList{}).Count(&numberOfRecords).Error
	case table == ToolingGatingTypeTab:
		err = database.Debug().Model(&ToolingGatingType{}).Count(&numberOfRecords).Error
	case table == ToolingHotRunnerBrandTable:
		err = database.Debug().Model(&ToolingHotRunnerBrand{}).Count(&numberOfRecords).Error
	case table == ToolingRunnerTypeTab:
		err = database.Debug().Model(&ToolingRunnerType{}).Count(&numberOfRecords).Error
	case table == ToolingHotRunnerConnectorTypeTable:
		err = database.Debug().Model(&ToolingHotRunnerConnectorType{}).Count(&numberOfRecords).Error
	case table == ToolingHotRunnerControllerTable:
		err = database.Debug().Model(&ToolingHotRunnerController{}).Count(&numberOfRecords).Error
	case table == ToolingHRControllerTable:
		err = database.Debug().Model(&ToolingHRController{}).Count(&numberOfRecords).Error
	case table == ToolingHotRunnerControllerBrandTable:
		err = database.Debug().Model(&ToolingHotRunnerControllerBrand{}).Count(&numberOfRecords).Error
	case table == ToolingCoolingFittingTable:
		err = database.Debug().Model(&ToolingCoolingFitting{}).Count(&numberOfRecords).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	if err != nil {
		return -1
	}
	return numberOfRecords
}

func GetCount(database *gorm.DB, table string, condition ...string) int64 {
	var iNumberOfRecords int64
	switch {
	case table == ToolingProjectTable:
		if len(condition) > 0 {
			database.Model(&ToolingProject{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&ToolingProject{}).Count(&iNumberOfRecords)
		}

	case table == ToolingProjectTaskTable:
		if len(condition) > 0 {
			database.Model(&ToolingProjectTask{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&ToolingProjectTask{}).Count(&iNumberOfRecords)
		}

	case table == ToolingProjectStatusTable:
		if len(condition) > 0 {
			database.Model(&ToolingProjectStatus{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&ToolingProjectStatus{}).Count(&iNumberOfRecords)
		}
	case table == ToolingProjectTaskStatusTable:
		if len(condition) > 0 {
			database.Model(&ToolingProjectTaskStatus{}).Where(condition[0]).Count(&iNumberOfRecords)
		} else {
			database.Model(&ToolingProjectTaskStatus{}).Count(&iNumberOfRecords)
		}
	default:
		return 0
	}
	return iNumberOfRecords

}

func CreateRecordTrail(database *gorm.DB, objectInterface ToolingRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == ToolingRecordTrailTable:
		var dbObjects []ToolingRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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

func (v *ToolingService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *ToolingService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
