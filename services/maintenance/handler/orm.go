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
	case table == MaintenanceRecordTrailTable:
		var dbObjects []MaintenanceRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaintenanceWorkOrderTable:
		var dbObjects []MaintenanceWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrder{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceWorkOrder{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaintenanceWorkOrderTaskTable:
		var dbObjects []MaintenanceWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaintenanceCorrectiveWorkOrderTable:
		var dbObjects []MaintenanceCorrectiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaintenanceWorkOrderCorrectiveTaskTable:
		var dbObjects []MaintenanceWorkOrderCorrectiveTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrderCorrectiveTask{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceWorkOrderCorrectiveTask{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MouldMaintenancePreventiveWorkOrderTable:
		var dbObjects []MouldMaintenancePreventiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		var dbObjects []MouldMaintenancePreventiveWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
	case table == MouldMaintenanceCorrectiveWorkOrderTable:
		var dbObjects []MouldMaintenanceCorrectiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MouldMaintenanceCorrectiveWorkOrderTaskTable:
		var dbObjects []MouldMaintenanceCorrectiveWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		var dbObjects []MouldCorrectiveMaintenanceJrOption
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldCorrectiveMaintenanceJrOption{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldCorrectiveMaintenanceJrOption{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
	default:
		return nil, errors.New(GetUnknownObjectType)
	}
	return nil, err
}
func CreateRecordTrail(database *gorm.DB, objectInterface MaintenanceRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == MaintenanceWorkOrderTable:
		object := MaintenanceWorkOrder{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceWorkOrderTaskTable:
		object := MaintenanceWorkOrderTask{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceWorkOrderTaskStatusTable:
		object := MaintenanceWorkOrderTaskStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceWorkOrderStatusTable:
		object := MaintenanceWorkOrderStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceEmailTemplateTable:
		object := MaintenanceEmailTemplate{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceCorrectiveWorkOrderTable:
		object := MaintenanceCorrectiveWorkOrder{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceWorkOrderCorrectiveTaskTable:
		object := MaintenanceWorkOrderCorrectiveTask{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenanceFaultCodeTable:
		object := MaintenanceFaultCode{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaintenancePreventiveWorkOrderStatusTable:
		object := MaintenancePreventiveWorkOrderStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMaintenancePreventiveWorkOrderTable:
		object := MouldMaintenancePreventiveWorkOrder{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		object := MouldMaintenancePreventiveWorkOrderTask{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMaintenanceCorrectiveWorkOrderTable:
		object := MouldMaintenanceCorrectiveWorkOrder{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMaintenanceCorrectiveWorkOrderTaskTable:
		object := MouldMaintenanceCorrectiveWorkOrderTask{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		object := MouldCorrectiveMaintenanceJrOption{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MachineMaintenanceSettingTable:
		object := MachineMaintenanceSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MouldMaintenanceSettingTable:
		object := MouldMaintenanceSetting{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == MaintenanceWorkOrderTable:
		dbObject := MaintenanceWorkOrder{Id: recordId}
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceWorkOrderTaskTable:
		dbObject := MaintenanceWorkOrderTask{Id: recordId}
		err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceWorkOrderTaskStatusTable:
		dbObject := MaintenanceWorkOrderTaskStatus{Id: recordId}
		err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceWorkOrderStatusTable:
		dbObject := MaintenanceWorkOrderStatus{Id: recordId}
		err = database.Debug().Model(&MaintenanceWorkOrderStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceComponentTable:
		dbObject := MaintenanceComponent{Id: recordId}
		err = database.Debug().Model(&MaintenanceComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceEmailTemplateTable:
		dbObject := MaintenanceEmailTemplate{Id: recordId}
		err = database.Debug().Model(&MaintenanceEmailTemplate{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceEmailTemplateFieldTable:
		dbObject := MaintenanceEmailTemplateField{Id: recordId}
		err = database.Debug().Model(&MaintenanceEmailTemplateField{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceCorrectiveWorkOrderComponent:
		dbObject := MaintenanceCorrectiveWorkOrder{Id: recordId}
		err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceWorkOrderCorrectiveTaskComponent:
		dbObject := MaintenanceWorkOrderCorrectiveTask{Id: recordId}
		err = database.Debug().Model(&MaintenanceWorkOrderCorrectiveTask{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenanceFaultCodeTable:
		dbObject := MaintenanceFaultCode{Id: recordId}
		err = database.Debug().Model(&MaintenanceFaultCode{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaintenancePreventiveWorkOrderStatusTable:
		dbObject := MaintenancePreventiveWorkOrderStatus{Id: recordId}
		err = database.Debug().Model(&MaintenancePreventiveWorkOrderStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMaintenancePreventiveWorkOrderTable:
		dbObject := MouldMaintenancePreventiveWorkOrder{Id: recordId}
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMaintenanceCorrectiveWorkOrderComponent:
		dbObject := MouldMaintenanceCorrectiveWorkOrder{Id: recordId}
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		dbObject := MouldMaintenancePreventiveWorkOrderTask{Id: recordId}
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMaintenanceCorrectiveWorkOrderTaskTable:
		dbObject := MouldMaintenanceCorrectiveWorkOrderTask{Id: recordId}
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MachineMaintenanceSettingTable:
		dbObject := MachineMaintenanceSetting{Id: recordId}
		err = database.Debug().Model(&MachineMaintenanceSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MouldMaintenanceSettingTable:
		dbObject := MouldMaintenanceSetting{Id: recordId}
		err = database.Debug().Model(&MouldMaintenanceSetting{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == MaintenanceWorkOrderTable:

		var dbObjects []MaintenanceWorkOrder
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceWorkOrder{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrder{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err

	case table == MaintenanceWorkOrderTaskTable:
		var dbObjects []MaintenanceWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceWorkOrderTask{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderTask{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceEmailTemplateTable:
		var dbObjects []MaintenanceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceEmailTemplate{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceEmailTemplate{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceWorkOrderStatusTable:
		var dbObjects []MaintenanceWorkOrderStatus
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceWorkOrderStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceWorkOrderTaskStatusTable:
		var dbObjects []MaintenanceWorkOrderTaskStatus
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceWorkOrderTaskStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderTaskStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceComponentTable:
		var dbObjects []MaintenanceComponent
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceCorrectiveWorkOrderComponent:
		var dbObjects []MaintenanceCorrectiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceCorrectiveWorkOrder{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceCorrectiveWorkOrder{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceWorkOrderCorrectiveTaskComponent:
		var dbObjects []MaintenanceWorkOrderCorrectiveTask
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceWorkOrderCorrectiveTask{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderCorrectiveTask{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceFaultCodeTable:
		var dbObjects []MaintenanceFaultCode
		if len(objectCount) > 0 {
			err = database.Model(&MaintenanceFaultCode{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceFaultCode{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenancePreventiveWorkOrderStatusTable:
		var dbObjects []MaintenancePreventiveWorkOrderStatus
		if len(objectCount) > 0 {
			err = database.Model(&MaintenancePreventiveWorkOrderStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenancePreventiveWorkOrderStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenancePreventiveWorkOrderTable:

		var dbObjects []MouldMaintenancePreventiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Model(&MouldMaintenancePreventiveWorkOrder{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenancePreventiveWorkOrder{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenanceCorrectiveWorkOrderComponent:
		var dbObjects []MouldMaintenanceCorrectiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Model(&MouldMaintenanceCorrectiveWorkOrder{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenanceCorrectiveWorkOrder{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		var dbObjects []MouldMaintenancePreventiveWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Model(&MouldMaintenancePreventiveWorkOrderTask{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenancePreventiveWorkOrderTask{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		var dbObjects []MouldCorrectiveMaintenanceJrOption
		if len(objectCount) > 0 {
			err = database.Model(&MouldCorrectiveMaintenanceJrOption{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldCorrectiveMaintenanceJrOption{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenanceCorrectiveWorkOrderTaskTable:
		var dbObjects []MouldMaintenanceCorrectiveWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Find(&dbObjects).Error
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
	case table == MaintenanceWorkOrderTable:
		var dbObjects []MaintenanceWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrder{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceWorkOrder{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaintenanceWorkOrderTaskTable:
		var dbObjects []MaintenanceWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderTask{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceWorkOrderTaskStatusTable:
		var dbObjects []MaintenanceWorkOrderTaskStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceWorkOrderStatusTable:
		var dbObjects []MaintenanceWorkOrderStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrderStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceEmailTemplateTable:
		var dbObjects []MaintenanceEmailTemplate
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceEmailTemplate{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceEmailTemplate{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceComponentTable:
		var dbObjects []MaintenanceComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceCorrectiveWorkOrderComponent:
		var dbObjects []MaintenanceCorrectiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceCorrectiveWorkOrder{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceWorkOrderCorrectiveTaskComponent:
		var dbObjects []MaintenanceWorkOrderCorrectiveTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceWorkOrderCorrectiveTask{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceWorkOrderCorrectiveTask{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenanceFaultCodeTable:
		var dbObjects []MaintenanceFaultCode
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenanceFaultCode{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenanceFaultCode{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MaintenancePreventiveWorkOrderStatusTable:
		var dbObjects []MaintenancePreventiveWorkOrderStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaintenancePreventiveWorkOrderStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaintenancePreventiveWorkOrderStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenancePreventiveWorkOrderTable:
		var dbObjects []MouldMaintenancePreventiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		var dbObjects []MouldMaintenancePreventiveWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenancePreventiveWorkOrderTask{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenanceCorrectiveWorkOrderComponent:
		var dbObjects []MouldMaintenanceCorrectiveWorkOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenanceCorrectiveWorkOrder{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldMaintenanceCorrectiveWorkOrderTaskTable:
		var dbObjects []MouldMaintenanceCorrectiveWorkOrderTask
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		var dbObjects []MouldCorrectiveMaintenanceJrOption
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MouldCorrectiveMaintenanceJrOption{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MouldCorrectiveMaintenanceJrOption{}).Where(condition).Find(&dbObjects).Error
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
	case table == MaintenanceWorkOrderTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Delete(&MaintenanceWorkOrder{Id: objectInterface.Id}).Error
	case table == MaintenanceWorkOrderTaskTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Delete(&MaintenanceWorkOrderTask{Id: objectInterface.Id}).Error
	case table == MaintenanceWorkOrderTaskStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderStatus{}).Delete(&MaintenanceWorkOrderStatus{Id: objectInterface.Id}).Error
	case table == MaintenanceWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Delete(&MaintenanceWorkOrderTaskStatus{Id: objectInterface.Id}).Error
	case table == MouldMaintenancePreventiveWorkOrderTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Delete(&MouldMaintenancePreventiveWorkOrder{Id: objectInterface.Id}).Error
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Delete(&MouldMaintenancePreventiveWorkOrderTask{Id: objectInterface.Id}).Error
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Delete(&MouldCorrectiveMaintenanceJrOption{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == MaintenanceWorkOrderTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Count(&numberOfRecords).Error
	case table == MaintenanceWorkOrderTaskTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Count(&numberOfRecords).Error
	case table == MaintenanceWorkOrderTaskStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Count(&numberOfRecords).Error
	case table == MaintenanceWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderStatus{}).Count(&numberOfRecords).Error
	case table == MaintenanceEmailTemplateTable:
		err = database.Debug().Model(&MaintenanceEmailTemplate{}).Count(&numberOfRecords).Error
	case table == MaintenanceCorrectiveWorkOrderTable:
		err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Count(&numberOfRecords).Error
	case table == MaintenanceWorkOrderCorrectiveTaskComponent:
		err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Count(&numberOfRecords).Error
	case table == MaintenanceFaultCodeTable:
		err = database.Debug().Model(&MaintenanceFaultCode{}).Count(&numberOfRecords).Error
	case table == MaintenancePreventiveWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenancePreventiveWorkOrderStatus{}).Count(&numberOfRecords).Error
	case table == MouldMaintenancePreventiveWorkOrderTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Count(&numberOfRecords).Error
	case table == MouldMaintenanceCorrectiveWorkOrderTable:
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Count(&numberOfRecords).Error
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Count(&numberOfRecords).Error
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		err = database.Debug().Model(&MouldCorrectiveMaintenanceJrOption{}).Count(&numberOfRecords).Error
	default:
		return -1
	}
	if err != nil {
		return -1
	}
	return numberOfRecords
}

func CountByCondition(database *gorm.DB, table, condition string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == MaintenanceWorkOrderTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Where(condition).Count(&numberOfRecords).Error
	case table == MaintenanceCorrectiveWorkOrderTable:
		err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Where(condition).Count(&numberOfRecords).Error
	case table == MouldMaintenanceCorrectiveWorkOrderTable:
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Where(condition).Count(&numberOfRecords).Error
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
	case table == MaintenanceWorkOrderTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Take(&MaintenanceWorkOrder{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderTaskTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Take(&MaintenanceWorkOrderTask{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderTaskStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Take(&MaintenanceWorkOrderTaskStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderStatus{}).Take(&MaintenanceWorkOrderStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceEmailTemplateTable:
		err = database.Debug().Model(&MaintenanceEmailTemplate{}).Take(&MaintenanceEmailTemplate{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceComponentTable:
		err = database.Debug().Model(&MaintenanceComponent{}).Take(&MaintenanceComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceCorrectiveWorkOrderComponent:
		err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Take(&MaintenanceCorrectiveWorkOrder{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderCorrectiveTaskComponent:
		err = database.Debug().Model(&MaintenanceWorkOrderCorrectiveTask{}).Take(&MaintenanceWorkOrderCorrectiveTask{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenanceFaultCodeTable:
		err = database.Debug().Model(&MaintenanceFaultCode{}).Take(&MaintenanceFaultCode{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaintenancePreventiveWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenancePreventiveWorkOrderStatus{}).Take(&MaintenancePreventiveWorkOrderStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMaintenancePreventiveWorkOrderTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Take(&MouldMaintenancePreventiveWorkOrder{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMaintenanceCorrectiveWorkOrderComponent:
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Take(&MouldMaintenanceCorrectiveWorkOrder{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Take(&MouldMaintenancePreventiveWorkOrderTask{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMaintenanceCorrectiveWorkOrderTaskTable:
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrderTask{}).Take(&MouldMaintenanceCorrectiveWorkOrderTask{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		err = database.Debug().Model(&MouldCorrectiveMaintenanceJrOption{}).Take(&MouldCorrectiveMaintenanceJrOption{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MachineMaintenanceSettingTable:
		err = database.Debug().Model(&MachineMaintenanceSetting{}).Take(&MachineMaintenanceSetting{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MouldMaintenanceSettingTable:
		err = database.Debug().Model(&MouldMaintenanceSetting{}).Take(&MouldMaintenanceSetting{Id: recordId}).UpdateColumns(updateObject).Error
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
	case table == MaintenanceWorkOrderTable:
		err = database.Debug().Model(&MaintenanceWorkOrder{}).Take(&MaintenanceWorkOrder{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderTaskTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTask{}).Take(&MaintenanceWorkOrderTask{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderTaskStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderTaskStatus{}).Take(&MaintenanceWorkOrderTaskStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenanceWorkOrderStatus{}).Take(&MaintenanceWorkOrderStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaintenanceCorrectiveWorkOrderComponent:
		err = database.Debug().Model(&MaintenanceCorrectiveWorkOrder{}).Take(&MaintenanceCorrectiveWorkOrder{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaintenanceWorkOrderCorrectiveTaskComponent:
		err = database.Debug().Model(&MaintenanceWorkOrderCorrectiveTask{}).Take(&MaintenanceWorkOrderCorrectiveTask{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaintenancePreventiveWorkOrderStatusTable:
		err = database.Debug().Model(&MaintenancePreventiveWorkOrderStatus{}).Take(&MaintenancePreventiveWorkOrderStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MouldMaintenancePreventiveWorkOrderTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrder{}).Take(&MouldMaintenancePreventiveWorkOrder{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MouldMaintenanceCorrectiveWorkOrderComponent:
		err = database.Debug().Model(&MouldMaintenanceCorrectiveWorkOrder{}).Take(&MouldMaintenanceCorrectiveWorkOrder{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MouldMaintenancePreventiveWorkOrderTaskTable:
		err = database.Debug().Model(&MouldMaintenancePreventiveWorkOrderTask{}).Take(&MouldMaintenancePreventiveWorkOrderTask{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MouldCorrectiveMaintenanceJrOptionTable:
		err = database.Debug().Model(&MouldCorrectiveMaintenanceJrOption{}).Take(&MouldCorrectiveMaintenanceJrOption{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MachineMaintenanceSettingTable:
		err = database.Debug().Model(&MachineMaintenanceSetting{}).Take(&MachineMaintenanceSetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MouldMaintenanceSettingTable:
		err = database.Debug().Model(&MouldMaintenanceSetting{}).Take(&MouldMaintenanceSetting{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *MaintenanceService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId) + " AND object_info ->>'$.objectStatus' = 'Active'"
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

func (v *MaintenanceService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
