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
	case table == MachinesRecordTrailTable:
		var dbObjects []MachinesRecordTrail
		if len(objectCount) > 0 {
			err = database.Model(&MachinesRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachinesRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineHMITable:
		var dbObjects []MachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMI{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMI{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AssemblyMachineHmiTable:
		var dbObjects []AssemblyMachineHmi
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineHmi{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineHmi{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingMachineHmiTable:
		var dbObjects []ToolingMachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHMI{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineHMI{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MouldMachineBrandTable:
		var dbObjects []MouldMachineBrand
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineBrand{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineBrand{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface MachinesRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func CreateWithId(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == MachineDisplaySettingTable:
		object := MachineDisplaySetting{Id: objectInterface.Id, ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineHMISettingSettingTable:
		object := MachineHMISetting{Id: objectInterface.Id, ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == MachineMasterTable:
		object := MachineMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineStatusTable:
		object := MachineStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineCategoryTable:
		object := MachineCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineSubCategoryTable:
		object := MachineSubCategory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineParameterTable:
		object := MachineParameter{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineSchedulerTable:
		object := MachineScheduler{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineHMITable:
		object := MachineHMI{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineHMIStopReasonTable:
		object := MachineHMIStopReason{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineDisplaySettingTable:
		object := MachineDisplaySetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineModuleSettingTable:
		object := MachineModuleSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineDashboardTable:
		object := MachineDashboard{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineWidgetTable:
		object := MachineWidget{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineEmailTemplateTable:
		object := MachineEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineEmailMasterFieldsTable:
		object := MachineEmailMasterFields{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyMachineMasterTable:
		object := AssemblyMachineMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyMachineHmiTable:
		object := AssemblyMachineHmi{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingMachineMasterTable:
		object := ToolingMachineMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingMachineHmiTable:
		object := ToolingMachineHMI{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingMachineHmiSettingTable:
		object := ToolingMachineHmiSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingMachineDisplaySettingTable:
		object := ToolingMachineDisplaySetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyMachineLineTable:
		object := AssemblyMachineLines{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineConnectStatusTable:
		object := MachineConnectStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyLineTypeTable:
		object := AssemblyLineType{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyEquipmentNameTable:
		object := AssemblyEquipmentName{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyEquipmentTypeMasterTable:
		object := AssemblyEquipmentTypeMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineFilterTable:
		object := MachineFilter{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMachineBrandTable:
		object := MouldMachineBrand{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMachineSettingTable:
		object := MouldMachineSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == MachineMasterTable:
		dbObject := MachineMaster{Id: recordId}
		err = database.Model(&MachineMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineStatusTable:
		dbObject := MachineStatus{Id: recordId}
		err = database.Model(&MachineStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineCategoryTable:
		dbObject := MachineCategory{Id: recordId}
		err = database.Model(&MachineCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineSubCategoryTable:
		dbObject := MachineSubCategory{Id: recordId}
		err = database.Model(&MachineSubCategory{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineComponentTable:
		dbObject := MachineComponent{Id: recordId}
		err = database.Model(&MachineComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineParameterTable:
		dbObject := MachineParameter{Id: recordId}
		err = database.Model(&MachineParameter{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineSchedulerTable:
		dbObject := MachineScheduler{Id: recordId}
		err = database.Model(&MachineScheduler{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineHMITable:
		dbObject := MachineHMI{Id: recordId}
		err = database.Model(&MachineHMI{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineHMIStopReasonTable:
		dbObject := MachineHMIStopReason{Id: recordId}
		err = database.Model(&MachineHMIStopReason{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineDisplaySettingTable:
		dbObject := MachineDisplaySetting{Id: recordId}
		err = database.Model(&MachineDisplaySetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineHMISettingSettingTable:
		dbObject := MachineHMISetting{Id: recordId}
		err = database.Model(&MachineHMISetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineModuleSettingTable:
		dbObject := MachineModuleSetting{Id: recordId}
		err = database.Model(&MachineModuleSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineDashboardTable:
		dbObject := MachineDashboard{Id: recordId}
		err = database.Model(&MachineDashboard{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineWidgetTable:
		dbObject := MachineWidget{Id: recordId}
		err = database.Model(&MachineWidget{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineEmailTemplateTable:
		dbObject := MachineEmailTemplate{Id: recordId}
		err = database.Model(&MachineEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineEmailMasterFieldsTable:
		dbObject := MachineEmailMasterFields{Id: recordId}
		err = database.Model(&MachineEmailMasterFields{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyMachineMasterTable:
		dbObject := AssemblyMachineMaster{Id: recordId}
		err = database.Model(&AssemblyMachineMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyMachineHmiTable:
		dbObject := AssemblyMachineHmi{Id: recordId}
		err = database.Model(&AssemblyMachineHmi{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyMachineHmiSettingTable:
		dbObject := AssemblyMachineHmiSetting{Id: recordId}
		err = database.Model(&AssemblyMachineHmiSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyMachineDisplaySettingTable:
		dbObject := AssemblyMachineDisplaySetting{Id: recordId}
		err = database.Model(&AssemblyMachineDisplaySetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingMachineMasterTable:
		dbObject := ToolingMachineMaster{Id: recordId}
		err = database.Model(&ToolingMachineMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingMachineHmiTable:
		dbObject := ToolingMachineHMI{Id: recordId}
		err = database.Model(&ToolingMachineHMI{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingMachineDisplaySettingTable:
		dbObject := ToolingMachineDisplaySetting{Id: recordId}
		err = database.Model(&ToolingMachineDisplaySetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingMachineHmiSettingTable:
		dbObject := ToolingMachineHmiSetting{Id: recordId}
		err = database.Model(&ToolingMachineHmiSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyMachineLineTable:
		dbObject := AssemblyMachineLines{Id: recordId}
		err = database.Model(&AssemblyMachineLines{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineConnectStatusTable:
		dbObject := MachineConnectStatus{Id: recordId}
		err = database.Model(&MachineConnectStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyLineTypeTable:
		dbObject := AssemblyLineType{Id: recordId}
		err = database.Model(&AssemblyLineType{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyEquipmentTypeMasterTable:
		dbObject := AssemblyEquipmentTypeMaster{Id: recordId}
		err = database.Model(&AssemblyEquipmentTypeMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyEquipmentNameTable:
		dbObject := AssemblyEquipmentName{Id: recordId}
		err = database.Model(&AssemblyEquipmentName{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMachineBrandTable:
		dbObject := MouldMachineBrand{Id: recordId}
		err = database.Model(&MouldMachineBrand{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineFilterTable:
		dbObject := MachineFilter{Id: recordId}
		err = database.Model(&MachineFilter{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMachineSettingTable:
		dbObject := MouldMachineSetting{Id: recordId}
		err = database.Model(&MouldMachineSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == MachineMasterTable:

		var dbObjects []MachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&MachineMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineParameterTable:
		var dbObjects []MachineParameter
		if len(objectCount) > 0 {
			err = database.Model(&MachineParameter{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineParameter{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineCategoryTable:
		var dbObjects []MachineCategory
		if len(objectCount) > 0 {
			err = database.Model(&MachineCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineSubCategoryTable:
		var dbObjects []MachineSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&MachineSubCategory{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineSubCategory{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineStatusTable:
		var dbObjects []MachineStatus
		if len(objectCount) > 0 {
			err = database.Model(&MachineStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineComponentTable:
		var dbObjects []MachineComponent
		if len(objectCount) > 0 {
			err = database.Model(&MachineComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineSchedulerTable:
		var dbObjects []MachineScheduler
		if len(objectCount) > 0 {
			err = database.Model(&MachineScheduler{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineScheduler{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMITable:
		var dbObjects []MachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMI{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMI{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMIStopReasonTable:
		var dbObjects []MachineHMIStopReason
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMIStopReason{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMIStopReason{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineDisplaySettingTable:
		var dbObjects []MachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineDisplaySetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineDisplaySetting{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMISettingSettingTable:
		var dbObjects []MachineHMISetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMISetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMISetting{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineModuleSettingTable:
		var dbObjects []MachineModuleSetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineModuleSetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineModuleSetting{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineDashboardTable:
		var dbObjects []MachineDashboard
		if len(objectCount) > 0 {
			err = database.Model(&MachineDashboard{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineDashboard{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineWidgetTable:
		var dbObjects []MachineWidget
		if len(objectCount) > 0 {
			err = database.Model(&MachineWidget{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineWidget{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineEmailTemplateTable:
		var dbObjects []MachineEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&MachineEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineEmailMasterFieldsTable:
		var dbObjects []MachineEmailMasterFields
		if len(objectCount) > 0 {
			err = database.Model(&MachineEmailMasterFields{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineEmailMasterFields{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineMasterTable:
		var dbObjects []AssemblyMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineHmiTable:
		var dbObjects []AssemblyMachineHmi
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineHmi{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineHmi{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineDisplaySettingTable:
		var dbObjects []AssemblyMachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineDisplaySetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineDisplaySetting{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineMasterTable:
		var dbObjects []ToolingMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineHmiTable:
		var dbObjects []ToolingMachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHMI{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineHMI{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineHmiSettingTable:
		var dbObjects []ToolingMachineHmiSetting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHmiSetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineHmiSetting{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineDisplaySettingTable:
		var dbObjects []ToolingMachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineDisplaySetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineDisplaySetting{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineLineTable:
		var dbObjects []AssemblyMachineLines
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineLines{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineLines{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineConnectStatusTable:
		var dbObjects []MachineConnectStatus
		if len(objectCount) > 0 {
			err = database.Model(&MachineConnectStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineConnectStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyLineTypeTable:
		var dbObjects []AssemblyLineType
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyLineType{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyLineType{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyEquipmentNameTable:
		var dbObjects []AssemblyEquipmentName
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyEquipmentName{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyEquipmentName{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyEquipmentTypeMasterTable:
		var dbObjects []AssemblyEquipmentTypeMaster
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyEquipmentTypeMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyEquipmentTypeMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineFilterTable:
		var dbObjects []MachineFilter
		if len(objectCount) > 0 {
			err = database.Model(&MachineFilter{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineFilter{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMachineBrandTable:
		var dbObjects []MouldMachineBrand
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineBrand{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineBrand{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMachineSettingTable:
		var dbObjects []MouldMachineSetting
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineSetting{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineSetting{}).Find(&dbObjects).Error
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
	case table == MachineMasterTable:
		var dbObjects []MachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&MachineMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineParameterTable:
		var dbObjects []MachineParameter
		if len(objectCount) > 0 {
			err = database.Model(&MachineParameter{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineParameter{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineCategoryTable:
		var dbObjects []MachineCategory
		if len(objectCount) > 0 {
			err = database.Model(&MachineCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineSubCategoryTable:
		var dbObjects []MachineSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&MachineSubCategory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineSubCategory{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineStatusTable:
		var dbObjects []MachineStatus
		if len(objectCount) > 0 {
			err = database.Model(&MachineStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineSchedulerTable:
		var dbObjects []MachineScheduler
		if len(objectCount) > 0 {
			err = database.Model(&MachineScheduler{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineScheduler{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMIStopReasonTable:
		var dbObjects []MachineHMIStopReason
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMIStopReason{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMIStopReason{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineDisplaySettingTable:
		var dbObjects []MachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineDisplaySetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineDisplaySetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMITable:
		var dbObjects []MachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMI{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMI{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMISettingSettingTable:
		var dbObjects []MachineHMISetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMISetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineHMISetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineModuleSettingTable:
		var dbObjects []MachineModuleSetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineModuleSetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineModuleSetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineWidgetTable:
		var dbObjects []MachineWidget
		if len(objectCount) > 0 {
			err = database.Model(&MachineWidget{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineWidget{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineDashboardTable:
		var dbObjects []MachineDashboard
		if len(objectCount) > 0 {
			err = database.Model(&MachineDashboard{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineDashboard{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineEmailTemplateTable:
		var dbObjects []MachineEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&MachineEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineEmailMasterFieldsTable:
		var dbObjects []MachineEmailMasterFields
		if len(objectCount) > 0 {
			err = database.Model(&MachineEmailMasterFields{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineEmailMasterFields{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineComponentTable:
		var dbObjects []MachineComponent
		if len(objectCount) > 0 {
			err = database.Model(&MachineComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineMasterTable:
		var dbObjects []AssemblyMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineHmiTable:
		var dbObjects []AssemblyMachineHmi
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineHmi{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineHmi{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineDisplaySettingTable:
		var dbObjects []AssemblyMachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineDisplaySetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineDisplaySetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineMasterTable:
		var dbObjects []ToolingMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineHmiTable:
		var dbObjects []ToolingMachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHMI{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineHMI{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineHmiSettingTable:
		var dbObjects []ToolingMachineHmiSetting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHmiSetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineHmiSetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineDisplaySettingTable:
		var dbObjects []ToolingMachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineDisplaySetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingMachineDisplaySetting{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineLineTable:
		var dbObjects []AssemblyMachineLines
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineLines{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyMachineLines{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineConnectStatusTable:
		var dbObjects []MachineConnectStatus
		if len(objectCount) > 0 {
			err = database.Model(&MachineConnectStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineConnectStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyLineTypeTable:
		var dbObjects []AssemblyLineType
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyLineType{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyLineType{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyEquipmentNameTable:
		var dbObjects []AssemblyEquipmentName
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyEquipmentName{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyEquipmentName{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyEquipmentTypeMasterTable:
		var dbObjects []AssemblyEquipmentTypeMaster
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyEquipmentTypeMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyEquipmentTypeMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMachineBrandTable:
		var dbObjects []MouldMachineBrand
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineBrand{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineBrand{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineFilterTable:
		var dbObjects []MachineFilter
		if len(objectCount) > 0 {
			err = database.Model(&MachineFilter{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MachineFilter{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMachineSettingTable:
		var dbObjects []MouldMachineSetting
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineSetting{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMachineSetting{}).Where(condition).Find(&dbObjects).Error
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

func GetObjectsById(database *gorm.DB, table string, idList []int, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == MachineMasterTable:
		var dbObjects []MachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&MachineMaster{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineMaster{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineParameterTable:
		var dbObjects []MachineParameter
		if len(objectCount) > 0 {
			err = database.Model(&MachineParameter{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineParameter{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineCategoryTable:
		var dbObjects []MachineCategory
		if len(objectCount) > 0 {
			err = database.Model(&MachineCategory{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineCategory{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineSubCategoryTable:
		var dbObjects []MachineSubCategory
		if len(objectCount) > 0 {
			err = database.Model(&MachineSubCategory{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineSubCategory{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MachineStatusTable:
		var dbObjects []MachineStatus
		if len(objectCount) > 0 {
			err = database.Model(&MachineStatus{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineStatus{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineSchedulerTable:
		var dbObjects []MachineScheduler
		if len(objectCount) > 0 {
			err = database.Model(&MachineScheduler{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineScheduler{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMIStopReasonTable:
		var dbObjects []MachineHMIStopReason
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMIStopReason{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineHMIStopReason{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineDisplaySettingTable:
		var dbObjects []MachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineDisplaySetting{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineDisplaySetting{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMITable:
		var dbObjects []MachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMI{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineHMI{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineHMISettingSettingTable:
		var dbObjects []MachineHMISetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineHMISetting{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineHMISetting{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineModuleSettingTable:
		var dbObjects []MachineModuleSetting
		if len(objectCount) > 0 {
			err = database.Model(&MachineModuleSetting{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineModuleSetting{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineWidgetTable:
		var dbObjects []MachineWidget
		if len(objectCount) > 0 {
			err = database.Model(&MachineWidget{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineWidget{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineDashboardTable:
		var dbObjects []MachineDashboard
		if len(objectCount) > 0 {
			err = database.Model(&MachineDashboard{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineDashboard{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineEmailTemplateTable:
		var dbObjects []MachineEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&MachineEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineEmailTemplate{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineEmailMasterFieldsTable:
		var dbObjects []MachineEmailMasterFields
		if len(objectCount) > 0 {
			err = database.Model(&MachineEmailMasterFields{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineEmailMasterFields{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineMasterTable:
		var dbObjects []AssemblyMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineMaster{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&AssemblyMachineMaster{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineMasterTable:
		var dbObjects []ToolingMachineMaster
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineMaster{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&ToolingMachineMaster{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineHmiTable:
		var dbObjects []ToolingMachineHMI
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHMI{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&ToolingMachineHMI{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineHmiSettingTable:
		var dbObjects []ToolingMachineHmiSetting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineHmiSetting{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&ToolingMachineHmiSetting{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ToolingMachineDisplaySettingTable:
		var dbObjects []ToolingMachineDisplaySetting
		if len(objectCount) > 0 {
			err = database.Model(&ToolingMachineDisplaySetting{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&ToolingMachineDisplaySetting{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyMachineLineTable:
		var dbObjects []AssemblyMachineLines
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyMachineLines{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&AssemblyMachineLines{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineConnectStatusTable:
		var dbObjects []MachineConnectStatus
		if len(objectCount) > 0 {
			err = database.Model(&MachineConnectStatus{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineConnectStatus{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AssemblyLineTypeTable:
		var dbObjects []AssemblyLineType
		if len(objectCount) > 0 {
			err = database.Model(&AssemblyLineType{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&AssemblyLineType{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MachineFilterTable:
		var dbObjects []MachineFilter
		if len(objectCount) > 0 {
			err = database.Model(&MachineFilter{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MachineFilter{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMachineBrandTable:
		var dbObjects []MouldMachineBrand
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineBrand{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MouldMachineBrand{}).Find(&dbObjects, idList).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMachineSettingTable:
		var dbObjects []MouldMachineSetting
		if len(objectCount) > 0 {
			err = database.Model(&MouldMachineSetting{}).Limit(objectCount[0]).Find(&dbObjects, idList).Error
		} else {
			err = database.Model(&MouldMachineSetting{}).Find(&dbObjects, idList).Error
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

func CountByCondition(database *gorm.DB, table, condition string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == MachineHMITable:
		err = database.Model(&MachineHMI{}).Where(condition).Count(&numberOfRecords).Error
	case table == AssemblyMachineHmiTable:
		err = database.Model(&AssemblyMachineHmi{}).Where(condition).Count(&numberOfRecords).Error
	case table == ToolingMachineHmiTable:
		err = database.Model(&ToolingMachineHMI{}).Where(condition).Count(&numberOfRecords).Error
	default:
		return -1
	}

	if err != nil {
		return -1
	}
	return numberOfRecords
}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == MachineMasterTable:
		err = database.Model(&MachineMaster{}).Count(&numberOfRecords).Error
	case table == MachineStatusTable:
		err = database.Model(&MachineStatus{}).Count(&numberOfRecords).Error
	case table == MachineCategoryTable:
		err = database.Model(&MachineCategory{}).Count(&numberOfRecords).Error
	case table == MachineSubCategoryTable:
		err = database.Model(&MachineSubCategory{}).Count(&numberOfRecords).Error
	case table == MachineParameterTable:
		err = database.Model(&MachineParameter{}).Count(&numberOfRecords).Error
	case table == MachineSchedulerTable:
		err = database.Model(&MachineScheduler{}).Count(&numberOfRecords).Error
	case table == MachineHMIStopReasonTable:
		err = database.Model(&MachineHMIStopReason{}).Count(&numberOfRecords).Error
	case table == MachineDisplaySettingTable:
		err = database.Model(&MachineDisplaySetting{}).Count(&numberOfRecords).Error
	case table == MachineHMISettingSettingTable:
		err = database.Model(&MachineHMISetting{}).Count(&numberOfRecords).Error
	case table == MachineModuleSettingTable:
		err = database.Model(&MachineModuleSetting{}).Count(&numberOfRecords).Error
	case table == MachineEmailTemplateTable:
		err = database.Model(&MachineEmailTemplate{}).Count(&numberOfRecords).Error
	case table == MachineEmailMasterFieldsTable:
		err = database.Model(&MachineEmailMasterFields{}).Count(&numberOfRecords).Error
	case table == AssemblyMachineMasterTable:
		err = database.Model(&AssemblyMachineMaster{}).Count(&numberOfRecords).Error
	case table == ToolingMachineMasterTable:
		err = database.Model(&ToolingMachineMaster{}).Count(&numberOfRecords).Error
	case table == ToolingMachineHmiTable:
		err = database.Model(&ToolingMachineHMI{}).Count(&numberOfRecords).Error
	case table == ToolingMachineDisplaySettingTable:
		err = database.Model(&ToolingMachineDisplaySetting{}).Count(&numberOfRecords).Error
	case table == ToolingMachineHmiSettingTable:
		err = database.Model(&ToolingMachineHmiSetting{}).Count(&numberOfRecords).Error
	case table == MachineHMITable:
		err = database.Model(&MachineHMI{}).Count(&numberOfRecords).Error
	case table == AssemblyMachineLineTable:
		err = database.Model(&AssemblyMachineLines{}).Count(&numberOfRecords).Error
	case table == MachineConnectStatusTable:
		err = database.Model(&MachineConnectStatus{}).Count(&numberOfRecords).Error
	case table == AssemblyLineTypeTable:
		err = database.Model(&AssemblyLineType{}).Count(&numberOfRecords).Error
	case table == AssemblyEquipmentTypeMasterTable:
		err = database.Model(&AssemblyEquipmentTypeMaster{}).Count(&numberOfRecords).Error
	case table == AssemblyEquipmentNameTable:
		err = database.Model(&AssemblyEquipmentName{}).Count(&numberOfRecords).Error
	case table == MachineFilterTable:
		err = database.Model(&MachineFilter{}).Count(&numberOfRecords).Error
	case table == MouldMachineBrandTable:
		err = database.Model(&MouldMachineBrand{}).Count(&numberOfRecords).Error
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
	switch {
	case table == MachineMasterTable:
		err = database.Model(&MachineMaster{}).Delete(&MachineMaster{Id: objectInterface.Id}).Error
	case table == MachineStatusTable:
		err = database.Model(&MachineStatus{}).Delete(&MachineStatus{Id: objectInterface.Id}).Error
	case table == MachineCategoryTable:
		err = database.Model(&MachineCategory{}).Delete(&MachineCategory{Id: objectInterface.Id}).Error
	case table == MachineSubCategoryTable:
		err = database.Model(&MachineSubCategory{}).Delete(&MachineSubCategory{Id: objectInterface.Id}).Error
	case table == MachineParameterTable:
		err = database.Model(&MachineParameter{}).Delete(&MachineParameter{Id: objectInterface.Id}).Error
	case table == MachineSchedulerTable:
		err = database.Model(&MachineScheduler{}).Delete(&MachineScheduler{Id: objectInterface.Id}).Error
	case table == MachineHMIStopReasonTable:
		err = database.Model(&MachineHMIStopReason{}).Delete(&MachineHMIStopReason{Id: objectInterface.Id}).Error
	case table == MachineDisplaySettingTable:
		err = database.Model(&MachineDisplaySetting{}).Delete(&MachineDisplaySetting{Id: objectInterface.Id}).Error
	case table == MachineHMISettingSettingTable:
		err = database.Model(&MachineHMISetting{}).Delete(&MachineHMISetting{Id: objectInterface.Id}).Error
	case table == MachineModuleSettingTable:
		err = database.Model(&MachineModuleSetting{}).Delete(&MachineModuleSetting{Id: objectInterface.Id}).Error
	case table == MachineEmailTemplateTable:
		err = database.Model(&MachineEmailTemplate{}).Delete(&MachineEmailTemplate{Id: objectInterface.Id}).Error
	case table == MachineEmailMasterFieldsTable:
		err = database.Model(&MachineEmailMasterFields{}).Delete(&MachineEmailMasterFields{Id: objectInterface.Id}).Error
	case table == AssemblyMachineMasterTable:
		err = database.Model(&AssemblyMachineMaster{}).Delete(&AssemblyMachineMaster{Id: objectInterface.Id}).Error
	case table == AssemblyLineTypeTable:
		err = database.Model(&AssemblyLineType{}).Delete(&AssemblyLineType{Id: objectInterface.Id}).Error
	case table == AssemblyEquipmentTypeMasterTable:
		err = database.Model(&AssemblyEquipmentTypeMaster{}).Delete(&AssemblyEquipmentTypeMaster{Id: objectInterface.Id}).Error
	case table == AssemblyEquipmentNameTable:
		err = database.Model(&AssemblyEquipmentName{}).Delete(&AssemblyEquipmentName{Id: objectInterface.Id}).Error
	case table == MachineFilterTable:
		err = database.Model(&MachineFilter{}).Delete(&MachineFilter{Id: objectInterface.Id}).Error
	case table == MouldMachineBrandTable:
		err = database.Model(&MouldMachineBrand{}).Delete(&MouldMachineBrand{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == MachineMasterTable:
		err = database.Model(&MachineMaster{}).Take(&MachineMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineStatusTable:
		err = database.Model(&MachineStatus{}).Take(&MachineStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineCategoryTable:
		err = database.Model(&MachineCategory{}).Take(&MachineCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineSubCategoryTable:
		err = database.Model(&MachineSubCategory{}).Take(&MachineSubCategory{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineParameterTable:
		err = database.Model(&MachineParameter{}).Take(&MachineParameter{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineSchedulerTable:
		err = database.Model(&MachineScheduler{}).Take(&MachineScheduler{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineHMIStopReasonTable:
		err = database.Model(&MachineHMIStopReason{}).Take(&MachineHMIStopReason{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineDisplaySettingTable:
		err = database.Model(&MachineDisplaySetting{}).Take(&MachineDisplaySetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineHMISettingSettingTable:
		err = database.Model(&MachineHMISetting{}).Take(&MachineHMISetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineModuleSettingTable:
		err = database.Model(&MachineModuleSetting{}).Take(&MachineModuleSetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineWidgetTable:
		err = database.Model(&MachineWidget{}).Take(&MachineWidget{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineDashboardTable:
		err = database.Model(&MachineDashboard{}).Take(&MachineDashboard{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineEmailTemplateTable:
		err = database.Model(&MachineEmailTemplate{}).Take(&MachineEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineEmailMasterFieldsTable:
		err = database.Model(&MachineEmailMasterFields{}).Take(&MachineEmailMasterFields{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineComponentTable:
		err = database.Model(&MachineComponent{}).Take(&MachineComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineHMITable:
		err = database.Model(&MachineHMI{}).Take(&MachineHMI{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineMasterTable:
		err = database.Model(&AssemblyMachineMaster{}).Take(&AssemblyMachineMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineHmiTable:
		err = database.Model(&AssemblyMachineHmi{}).Take(&AssemblyMachineHmi{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineHmiSettingTable:
		err = database.Model(&AssemblyMachineHmiSetting{}).Take(&AssemblyMachineHmiSetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineDisplaySettingTable:
		err = database.Model(&AssemblyMachineDisplaySetting{}).Take(&AssemblyMachineDisplaySetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingMachineMasterTable:
		err = database.Model(&ToolingMachineMaster{}).Take(&ToolingMachineMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingMachineHmiTable:
		err = database.Model(&ToolingMachineHMI{}).Take(&ToolingMachineHMI{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingMachineHmiSettingTable:
		err = database.Model(&ToolingMachineHmiSetting{}).Take(&ToolingMachineHmiSetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingMachineDisplaySettingTable:
		err = database.Model(&ToolingMachineDisplaySetting{}).Take(&ToolingMachineDisplaySetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineLineTable:
		err = database.Model(&AssemblyMachineLines{}).Take(&AssemblyMachineLines{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineConnectStatusTable:
		err = database.Model(&MachineConnectStatus{}).Take(&MachineConnectStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyLineTypeTable:
		err = database.Model(&AssemblyLineType{}).Take(&AssemblyLineType{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyEquipmentNameTable:
		err = database.Model(&AssemblyEquipmentName{}).Take(&AssemblyEquipmentName{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyEquipmentTypeMasterTable:
		err = database.Model(&AssemblyEquipmentTypeMaster{}).Take(&AssemblyEquipmentTypeMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineFilterTable:
		err = database.Model(&MachineFilter{}).Take(&MachineFilter{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMachineBrandTable:
		err = database.Model(&MouldMachineBrand{}).Take(&MouldMachineBrand{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMachineSettingTable:
		err = database.Model(&MouldMachineSetting{}).Take(&MouldMachineSetting{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func NoOfReferenceObjects(database *gorm.DB, table string, referenceField string, id int) int {
	conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$." + referenceField + "\")) =" + strconv.Itoa(id)
	listOfObjects, err := GetConditionalObjects(database, table, conditionString)
	if err != nil {
		return -1
	} else {
		return len(*listOfObjects)
	}
}

func ArchiveReferenceObjects(database *gorm.DB, table string, referenceField string, id int) error {
	var err error

	conditionString := " object_info->>'$." + referenceField + "\" =" + strconv.Itoa(id)
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
		case table == MachineMasterTable:
			err = database.Model(&MachineMaster{}).Take(&MachineMaster{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineStatusTable:
			err = database.Model(&MachineStatus{}).Take(&MachineStatus{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineCategoryTable:
			err = database.Model(&MachineCategory{}).Take(&MachineCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineSubCategoryTable:
			err = database.Model(&MachineSubCategory{}).Take(&MachineSubCategory{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineParameterTable:
			err = database.Model(&MachineParameter{}).Take(&MachineParameter{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineSchedulerTable:
			err = database.Model(&MachineScheduler{}).Take(&MachineScheduler{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineHMIStopReasonTable:
			err = database.Model(&MachineHMIStopReason{}).Take(&MachineHMIStopReason{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineDisplaySettingTable:
			err = database.Model(&MachineDisplaySetting{}).Take(&MachineDisplaySetting{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineHMISettingSettingTable:
			err = database.Model(&MachineHMISetting{}).Take(&MachineHMISetting{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineModuleSettingTable:
			err = database.Model(&MachineModuleSetting{}).Take(&MachineModuleSetting{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineWidgetTable:
			err = database.Model(&MachineWidget{}).Take(&MachineWidget{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineDashboardTable:
			err = database.Model(&MachineDashboard{}).Take(&MachineDashboard{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineEmailTemplateTable:
			err = database.Model(&MachineEmailTemplate{}).Take(&MachineEmailTemplate{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MachineEmailMasterFieldsTable:
			err = database.Model(&MachineEmailMasterFields{}).Take(&MachineEmailMasterFields{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == AssemblyMachineMasterTable:
			err = database.Model(&AssemblyMachineMaster{}).Take(&AssemblyMachineMaster{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == MouldMachineBrandTable:
			err = database.Model(&MouldMachineBrand{}).Take(&MouldMachineBrand{Id: object.Id}).UpdateColumns(updateObject).Error
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
	case table == MachineMasterTable:
		err = database.Model(&MachineMaster{}).Take(&MachineMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineStatusTable:
		err = database.Model(&MachineStatus{}).Take(&MachineStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineCategoryTable:
		err = database.Model(&MachineCategory{}).Take(&MachineCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineSubCategoryTable:
		err = database.Model(&MachineSubCategory{}).Take(&MachineSubCategory{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineParameterTable:
		err = database.Model(&MachineParameter{}).Take(&MachineParameter{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineSchedulerTable:
		err = database.Model(&MachineScheduler{}).Take(&MachineScheduler{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineHMIStopReasonTable:
		err = database.Model(&MachineHMIStopReason{}).Take(&MachineHMIStopReason{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineDisplaySettingTable:
		err = database.Model(&MachineDisplaySetting{}).Take(&MachineDisplaySetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineHMISettingSettingTable:
		err = database.Model(&MachineHMISetting{}).Take(&MachineHMISetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineModuleSettingTable:
		err = database.Model(&MachineModuleSetting{}).Take(&MachineModuleSetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineWidgetTable:
		err = database.Model(&MachineWidget{}).Take(&MachineWidget{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineDashboardTable:
		err = database.Model(&MachineDashboard{}).Take(&MachineDashboard{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineEmailTemplateTable:
		err = database.Model(&MachineEmailTemplate{}).Take(&MachineEmailTemplate{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineEmailMasterFieldsTable:
		err = database.Model(&MachineEmailMasterFields{}).Take(&MachineEmailMasterFields{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineMasterTable:
		err = database.Model(&AssemblyMachineMaster{}).Take(&AssemblyMachineMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingMachineMasterTable:
		err = database.Model(&ToolingMachineMaster{}).Take(&ToolingMachineMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingMachineHmiTable:
		err = database.Model(&ToolingMachineHMI{}).Take(&ToolingMachineHMI{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingMachineHmiSettingTable:
		err = database.Model(&ToolingMachineHmiSetting{}).Take(&ToolingMachineHmiSetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingMachineDisplaySettingTable:
		err = database.Model(&ToolingMachineDisplaySetting{}).Take(&ToolingMachineDisplaySetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AssemblyMachineLineTable:
		err = database.Model(&AssemblyMachineLines{}).Take(&AssemblyMachineLines{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineConnectStatusTable:
		err = database.Model(&MachineConnectStatus{}).Take(&MachineConnectStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AssemblyLineTypeTable:
		err = database.Model(&AssemblyLineType{}).Take(&AssemblyLineType{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineFilterTable:
		err = database.Model(&MachineFilter{}).Take(&MachineFilter{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MouldMachineBrandTable:
		err = database.Model(&MouldMachineBrand{}).Take(&MouldMachineBrand{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err

}

func (v *MachineService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			if constraint.IsObjectField {
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
			} else {
				err, referenceCommonObject := Get(dbConnection, referenceTable, recordId)
				if err == nil {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + 1
					v.checkReference(dbConnection, referenceComponent, referenceCommonObject.Id, dependencyComponents, dependencyRecords)
				}

			}

		}
	}
}

func (v *MachineService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
