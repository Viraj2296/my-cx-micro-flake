package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/production_order/handler/utils"
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
	case table == ProductionOrderRecordTrailTable:
		var dbObjects []ProductionOrderRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ProductionOrderRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ProductionOrderRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ProductionOrderMasterTable:
		var dbObjects []ProductionOrderMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ProductionOrderMaster{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ProductionOrderMaster{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AssemblyProductionOrderTable:
		var dbObjects []AssemblyProductionOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AssemblyProductionOrder{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&AssemblyProductionOrder{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingOrderMasterTable:
		var dbObjects []ToolingOrderMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingOrderMaster{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&ToolingOrderMaster{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface ProductionOrderRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == ProductionOrderMasterTable:
		object := ProductionOrderMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ScheduledOrderEventTable:
		object := ScheduledOrderEvent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ProductionOrderStatusTable:
		object := ProductionOrderStatus{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == MaterialMasterTable:
		object := MaterialMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyProductionOrderTable:
		object := AssemblyProductionOrder{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AssemblyScheduledOrderEventTable:
		object := AssemblyScheduledOrderEvent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingBomMasterTable:
		object := ToolingBomMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingPartMasterTable:
		object := ToolingPartMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingOrderMasterTable:
		object := ToolingOrderMaster{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == ToolingScheduledOrderEventTable:
		object := ToolingScheduledOrderEvent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == utils.AssemblyManualOrderCompletedQuantityHistoryTable:
		object := AssemblyManualOrderCompletedQuantityHistory{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == ProductionOrderMasterTable:
		dbObject := ProductionOrderMaster{Id: recordId}
		err = database.Debug().Model(&ProductionOrderMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ProductionOrderComponentTable:
		dbObject := ProductionOrderComponent{Id: recordId}
		err = database.Debug().Model(&ProductionOrderComponent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ProductionOrderStatusTable:
		dbObject := ProductionOrderStatus{Id: recordId}
		err = database.Debug().Model(&ProductionOrderStatus{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ScheduledOrderEventTable:
		dbObject := ScheduledOrderEvent{Id: recordId}
		err = database.Debug().Model(&ScheduledOrderEvent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == MaterialMasterTable:
		dbObject := MaterialMaster{Id: recordId}
		err = database.Debug().Model(&MaterialMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyProductionOrderTable:
		dbObject := AssemblyProductionOrder{Id: recordId}
		err = database.Debug().Model(&AssemblyProductionOrder{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AssemblyScheduledOrderEventTable:
		dbObject := AssemblyScheduledOrderEvent{Id: recordId}
		err = database.Debug().Model(&AssemblyScheduledOrderEvent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingBomMasterTable:
		dbObject := ToolingBomMaster{Id: recordId}
		err = database.Debug().Model(&ToolingBomMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingPartMasterTable:
		dbObject := ToolingPartMaster{Id: recordId}
		err = database.Debug().Model(&ToolingPartMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingOrderMasterTable:
		dbObject := ToolingOrderMaster{Id: recordId}
		err = database.Debug().Model(&ToolingOrderMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == ToolingScheduledOrderEventTable:
		dbObject := ToolingScheduledOrderEvent{Id: recordId}
		err = database.Debug().Model(&ToolingScheduledOrderEvent{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == ProductionOrderMasterTable:

		var dbObjects []ProductionOrderMaster
		if len(objectCount) > 0 {
			err = database.Model(&ProductionOrderMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ProductionOrderMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ProductionOrderComponentTable:
		var dbObjects []ProductionOrderComponent
		if len(objectCount) > 0 {
			err = database.Model(&ProductionOrderComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ProductionOrderComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ProductionOrderStatusTable:
		var dbObjects []ProductionOrderStatus
		if len(objectCount) > 0 {
			err = database.Model(&ProductionOrderStatus{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ProductionOrderStatus{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == ScheduledOrderEventTable:
		var dbObjects []ScheduledOrderEvent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ScheduledOrderEvent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ScheduledOrderEvent{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaterialMasterTable:
		var dbObjects []MaterialMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaterialMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaterialMaster{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AssemblyProductionOrderTable:
		var dbObjects []AssemblyProductionOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AssemblyProductionOrder{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyProductionOrder{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AssemblyScheduledOrderEventTable:
		var dbObjects []AssemblyScheduledOrderEvent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AssemblyScheduledOrderEvent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyScheduledOrderEvent{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingBomMasterTable:
		var dbObjects []ToolingBomMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingBomMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingBomMaster{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingPartMasterTable:
		var dbObjects []ToolingPartMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingPartMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingPartMaster{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingOrderMasterTable:
		var dbObjects []ToolingOrderMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingOrderMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingOrderMaster{}).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingScheduledOrderEventTable:
		var dbObjects []ToolingScheduledOrderEvent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingScheduledOrderEvent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingScheduledOrderEvent{}).Find(&dbObjects).Error
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
	case table == ProductionOrderMasterTable:
		var dbObjects []ProductionOrderMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ProductionOrderMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ProductionOrderMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ProductionOrderComponentTable:
		var dbObjects []ProductionOrderComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ProductionOrderComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ProductionOrderComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ProductionOrderStatusTable:
		var dbObjects []ProductionOrderStatus
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ProductionOrderStatus{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ProductionOrderStatus{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ScheduledOrderEventTable:
		var dbObjects []ScheduledOrderEvent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ScheduledOrderEvent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ScheduledOrderEvent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == MaterialMasterTable:
		var dbObjects []MaterialMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&MaterialMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&MaterialMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AssemblyProductionOrderTable:
		var dbObjects []AssemblyProductionOrder
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AssemblyProductionOrder{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyProductionOrder{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AssemblyScheduledOrderEventTable:
		var dbObjects []AssemblyScheduledOrderEvent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AssemblyScheduledOrderEvent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyScheduledOrderEvent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingBomMasterTable:
		var dbObjects []ToolingBomMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingBomMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingBomMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingPartMasterTable:
		var dbObjects []ToolingPartMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingPartMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingPartMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingOrderMasterTable:
		var dbObjects []ToolingOrderMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingOrderMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingOrderMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == ToolingScheduledOrderEventTable:
		var dbObjects []ToolingScheduledOrderEvent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&ToolingScheduledOrderEvent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&ToolingScheduledOrderEvent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == utils.AssemblyManualOrderCompletedQuantityHistoryTable:
		var dbObjects []AssemblyManualOrderCompletedQuantityHistory
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AssemblyManualOrderCompletedQuantityHistory{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AssemblyManualOrderCompletedQuantityHistory{}).Where(condition).Find(&dbObjects).Error
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

func ArchiveObject(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	updateObject := make(map[string]interface{})
	var objectFields map[string]interface{}
	json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
	objectFields["objectStatus"] = "Archived"
	serializedObject, _ := json.Marshal(objectFields)
	updateObject["object_info"] = serializedObject

	switch {
	case table == ProductionOrderMasterTable:
		err = database.Debug().Model(&ProductionOrderMaster{}).Take(&ProductionOrderMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ProductionOrderStatusTable:
		err = database.Debug().Model(&ProductionOrderStatus{}).Take(&ProductionOrderStatus{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ScheduledOrderEventTable:
		err = database.Debug().Model(&ScheduledOrderEvent{}).Take(&ScheduledOrderEvent{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == MaterialMasterTable:
		err = database.Debug().Model(&MaterialMaster{}).Take(&MaterialMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AssemblyProductionOrderTable:
		err = database.Debug().Model(&AssemblyProductionOrder{}).Take(&AssemblyProductionOrder{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AssemblyScheduledOrderEventTable:
		err = database.Debug().Model(&AssemblyScheduledOrderEvent{}).Take(&AssemblyScheduledOrderEvent{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingBomMasterTable:
		err = database.Debug().Model(&ToolingBomMaster{}).Take(&ToolingBomMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingPartMasterTable:
		err = database.Debug().Model(&ToolingPartMaster{}).Take(&ToolingPartMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingOrderMasterTable:
		err = database.Debug().Model(&ToolingOrderMaster{}).Take(&ToolingOrderMaster{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == ToolingScheduledOrderEventTable:
		err = database.Debug().Model(&ToolingScheduledOrderEvent{}).Take(&ToolingScheduledOrderEvent{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}

func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	switch {
	case table == ProductionOrderMasterTable:
		err = database.Debug().Model(&ProductionOrderMaster{}).Delete(&ProductionOrderMaster{Id: objectInterface.Id}).Error
	case table == ProductionOrderStatusTable:
		err = database.Debug().Model(&ProductionOrderStatus{}).Delete(&ProductionOrderStatus{Id: objectInterface.Id}).Error
	case table == ScheduledOrderEventTable:
		err = database.Debug().Model(&ScheduledOrderEvent{}).Delete(&ScheduledOrderEvent{Id: objectInterface.Id}).Error
	case table == MaterialMasterTable:
		err = database.Debug().Model(&MaterialMaster{}).Delete(&MaterialMaster{Id: objectInterface.Id}).Error
	case table == AssemblyProductionOrderTable:
		err = database.Debug().Model(&AssemblyProductionOrder{}).Delete(&AssemblyProductionOrder{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == ProductionOrderMasterTable:
		err = database.Debug().Model(&ProductionOrderMaster{}).Count(&numberOfRecords).Error
	case table == ProductionOrderStatusTable:
		err = database.Debug().Model(&ProductionOrderStatus{}).Count(&numberOfRecords).Error
	case table == ScheduledOrderEventTable:
		err = database.Debug().Model(&ScheduledOrderEvent{}).Count(&numberOfRecords).Error
	case table == MaterialMasterTable:
		err = database.Debug().Model(&MaterialMaster{}).Count(&numberOfRecords).Error
	case table == AssemblyScheduledOrderEventTable:
		err = database.Debug().Model(&AssemblyScheduledOrderEvent{}).Count(&numberOfRecords).Error
	case table == ToolingOrderMasterTable:
		err = database.Debug().Model(&ToolingOrderMaster{}).Count(&numberOfRecords).Error
	case table == AssemblyProductionOrderTable:
		err = database.Debug().Model(&AssemblyProductionOrder{}).Count(&numberOfRecords).Error
	case table == ToolingScheduledOrderEventTable:
		err = database.Debug().Model(&ToolingScheduledOrderEvent{}).Count(&numberOfRecords).Error
	case table == ToolingPartMasterTable:
		err = database.Debug().Model(&ToolingPartMaster{}).Count(&numberOfRecords).Error
	case table == ToolingBomMasterTable:
		err = database.Debug().Model(&ToolingBomMaster{}).Count(&numberOfRecords).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	if err != nil {
		return -1
	}
	return numberOfRecords
}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == ProductionOrderMasterTable:
		err = database.Debug().Model(&ProductionOrderMaster{}).Take(&ProductionOrderMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ProductionOrderStatusTable:
		err = database.Debug().Model(&ProductionOrderStatus{}).Take(&ProductionOrderStatus{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ScheduledOrderEventTable:
		err = database.Debug().Model(&ScheduledOrderEvent{}).Take(&ScheduledOrderEvent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == MaterialMasterTable:
		err = database.Debug().Model(&MaterialMaster{}).Take(&MaterialMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ProductionOrderComponentTable:
		err = database.Debug().Model(&ProductionOrderComponent{}).Take(&ProductionOrderComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyProductionOrderTable:
		err = database.Debug().Model(&AssemblyProductionOrder{}).Take(&AssemblyProductionOrder{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AssemblyScheduledOrderEventTable:
		err = database.Debug().Model(&AssemblyScheduledOrderEvent{}).Take(&AssemblyScheduledOrderEvent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingBomMasterTable:
		err = database.Debug().Model(&ToolingBomMaster{}).Take(&ToolingBomMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingPartMasterTable:
		err = database.Debug().Model(&ToolingPartMaster{}).Take(&ToolingPartMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingOrderMasterTable:
		err = database.Debug().Model(&ToolingOrderMaster{}).Take(&ToolingOrderMaster{Id: recordId}).UpdateColumns(updateObject).Error
	case table == ToolingScheduledOrderEventTable:
		err = database.Debug().Model(&ToolingScheduledOrderEvent{}).Take(&ToolingScheduledOrderEvent{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *ProductionOrderService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *ProductionOrderService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
