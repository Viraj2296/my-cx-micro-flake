package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/production_order/handler/model"
	"cx-micro-flake/services/production_order/handler/rpc"
	"cx-micro-flake/services/production_order/handler/utils"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func (v *ProductionOrderService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, ProductionOrderComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *ProductionOrderService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, ProductionOrderComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, ProductionOrderComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *ProductionOrderService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, ProductionOrderComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *ProductionOrderService) GetChildEventsOfProductionOrder(projectId string, productionOrderId int) []int {
	conditionString := " object_info->>'$.eventSourceId' =  " + strconv.Itoa(productionOrderId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, _ := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)
	var listOfChildEvents []int
	for _, eventObject := range *listOfEvents {
		listOfChildEvents = append(listOfChildEvents, eventObject.Id)
	}
	return listOfChildEvents
}

func (v *ProductionOrderService) GetChildEventsOfToolingProductionOrder(projectId string, productionOrderId int) []int {
	conditionString := " object_info->>'$.eventSourceId' =  " + strconv.Itoa(productionOrderId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, _ := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionString)
	var listOfChildEvents []int
	for _, eventObject := range *listOfEvents {
		listOfChildEvents = append(listOfChildEvents, eventObject.Id)
	}
	return listOfChildEvents
}

func (v *ProductionOrderService) GetChildEventsOfAssemblyProductionOrder(projectId string, productionOrderId int) []int {
	conditionString := " object_info->>'$.eventSourceId' =  " + strconv.Itoa(productionOrderId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, _ := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionString)
	var listOfChildEvents []int
	for _, eventObject := range *listOfEvents {
		listOfChildEvents = append(listOfChildEvents, eventObject.Id)
	}
	return listOfChildEvents
}

func (v *ProductionOrderService) GetScheduledEvents(projectId string) (error, *[]component.GeneralObject) {
	preferenceFiveOrderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceFive)
	preferenceSixOrderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)
	eventQueryCondition := " object_info->>'$.eventStatus' =  " + strconv.Itoa(preferenceFiveOrderStatusId) + " OR   object_info->>'$.eventStatus' =  " + strconv.Itoa(preferenceSixOrderStatusId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, eventQueryCondition)

	return err, listOfEvents
}

func (v *ProductionOrderService) GetScheduledEventByProductionId(projectId string, productionOrderId int) (error, *[]component.GeneralObject) {
	eventQueryCondition := " object_info->>'$.eventSourceId' =  " + strconv.Itoa(productionOrderId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, eventQueryCondition)

	return err, listOfEvents
}

func (v *ProductionOrderService) GetAssemblyScheduledEvents(projectId string) (error, *[]component.GeneralObject) {
	preferenceFiveOrderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceFive)
	preferenceSixOrderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)
	eventQueryCondition := " object_info->>'$.eventStatus' =  " + strconv.Itoa(preferenceFiveOrderStatusId) + " OR   object_info->>'$.eventStatus' =  " + strconv.Itoa(preferenceSixOrderStatusId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, eventQueryCondition)

	return err, listOfEvents
}

func (v *ProductionOrderService) GetAssemblyScheduledEventByProductionId(projectId string, productionOrderId int) (error, *[]component.GeneralObject) {
	eventQueryCondition := " object_info->>'$.eventSourceId' =  " + strconv.Itoa(productionOrderId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, eventQueryCondition)

	return err, listOfEvents
}

func (v *ProductionOrderService) GetToolingScheduledEvents(projectId string) (error, *[]component.GeneralObject) {
	preferenceFiveOrderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceFive)
	preferenceSixOrderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)
	eventQueryCondition := " object_info->>'$.eventStatus' =  " + strconv.Itoa(preferenceFiveOrderStatusId) + " OR   object_info->>'$.eventStatus' =  " + strconv.Itoa(preferenceSixOrderStatusId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, eventQueryCondition)

	return err, listOfEvents
}
func (v *ProductionOrderService) GetToolingScheduledEventsByProductionId(projectId string, productionOrderId int) (error, *[]component.GeneralObject) {

	eventQueryCondition := " object_info->>'$.eventSourceId' =  " + strconv.Itoa(productionOrderId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfEvents, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, eventQueryCondition)

	return err, listOfEvents
}

func (v *ProductionOrderService) GetToolingPartById(projectId string, partId int) (error, component.GeneralObject) {

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, partMaster := Get(dbConnection, ToolingPartMasterTable, partId)

	return err, partMaster
}

// GetAllSchedulerEventForScheduler rename the function as this, because we are modifying the exisiting objects by adding extra logics.
func (v *ProductionOrderService) GetAllSchedulerEventForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.objectStatus' = 'Active' "
	listOfEvents, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, ProductionOrderStatusTable)
	var productionOrderStatusCache = make(map[int]*ProductionOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		productionStatus := ProductionOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		productionOrderStatusCache[statusGeneralObject.Id] = productionStatus.getProductionOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)

		if value, ok := scheduledEvent["eventStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if productionOrderInfo, ok := productionOrderStatusCache[eventStatusId]; ok {
				statusPreference := util.InterfaceToInt(productionOrderInfo.Preference)
				if statusPreference == ScheduleStatusPreferenceThree {
					scheduledEvent["canRelease"] = true
				} else {
					scheduledEvent["canRelease"] = false
				}
				if statusPreference == ScheduleStatusPreferenceFour {
					scheduledEvent["canHold"] = true
				} else {
					scheduledEvent["canHold"] = false
				}
				scheduledEvent["eventStatusName"] = productionOrderInfo.Status
				scheduledEvent["eventColor"] = productionOrderInfo.ColorCode
				scheduledEvent["module"] = "production_order"
				scheduledEvent["componentName"] = "scheduled_order_event"

				var recoveryScheduleName string
				if isRecoverySchedule, isFlag := scheduledEvent["isRecoverySchedule"]; isFlag {
					if util.InterfaceToBool(isRecoverySchedule) {
						recoveryScheduleName = getSchedulerName(listOfEvents, util.InterfaceToInt(scheduledEvent["recoveryScheduleId"]))
					}
				} else {
					scheduledEvent["isRecoverySchedule"] = false
				}

				scheduledEvent["recoveryScheduleId"] = recoveryScheduleName
				serializedEventObject, _ := json.Marshal(scheduledEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject
}

func getSchedulerName(listOfEvents *[]component.GeneralObject, id int) string {
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)

		if scheduledEventObject.Id == id {
			return util.InterfaceToString(scheduledEvent["name"])
		}
	}

	return ""
}

// GetAllAssemblyEventForScheduler rename the function as this, because we are modifying the exisiting objects by adding extra logics.
func (v *ProductionOrderService) GetAllAssemblyEventForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.objectStatus' = 'Active' "
	listOfEvents, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, ProductionOrderStatusTable)
	var productionOrderStatusCache = make(map[int]*ProductionOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		productionStatus := ProductionOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		productionOrderStatusCache[statusGeneralObject.Id] = productionStatus.getProductionOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)

		if value, ok := scheduledEvent["eventStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if productionOrderInfo, ok := productionOrderStatusCache[eventStatusId]; ok {
				statusPreference := util.InterfaceToInt(productionOrderInfo.Preference)
				if statusPreference == ScheduleStatusPreferenceThree {
					scheduledEvent["canRelease"] = true
				} else {
					scheduledEvent["canRelease"] = false
				}
				if statusPreference == ScheduleStatusPreferenceFour {
					scheduledEvent["canHold"] = true
				} else {
					scheduledEvent["canHold"] = false
				}
				scheduledEvent["eventStatusName"] = productionOrderInfo.Status
				scheduledEvent["eventColor"] = productionOrderInfo.ColorCode
				scheduledEvent["module"] = "production_order"
				scheduledEvent["componentName"] = "assembly_scheduled_order_event"
				if isRecoverySchedule, isFlag := scheduledEvent["isRecoverySchedule"]; isFlag {
					if util.InterfaceToBool(isRecoverySchedule) {
						recoveryScheduleName := getSchedulerName(listOfEvents, util.InterfaceToInt(scheduledEvent["recoveryScheduleId"]))
						scheduledEvent["recoveryScheduleId"] = recoveryScheduleName

					}
				} else {
					scheduledEvent["recoveryScheduleId"] = ""
					scheduledEvent["isRecoverySchedule"] = false
				}
				serializedEventObject, _ := json.Marshal(scheduledEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject
}

func (v *ProductionOrderService) GetAllToolingEventForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.objectStatus' = 'Active' "
	listOfEvents, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, ProductionOrderStatusTable)
	var productionOrderStatusCache = make(map[int]*ProductionOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		productionStatus := ProductionOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		productionOrderStatusCache[statusGeneralObject.Id] = productionStatus.getProductionOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)

		if value, ok := scheduledEvent["eventStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if productionOrderInfo, ok := productionOrderStatusCache[eventStatusId]; ok {
				statusPreference := util.InterfaceToInt(productionOrderInfo.Preference)
				if statusPreference == ScheduleStatusPreferenceThree {
					scheduledEvent["canRelease"] = true
				} else {
					scheduledEvent["canRelease"] = false
				}
				if statusPreference == ScheduleStatusPreferenceFour {
					scheduledEvent["canHold"] = true
				} else {
					scheduledEvent["canHold"] = false
				}
				scheduledEvent["eventStatusName"] = productionOrderInfo.Status
				scheduledEvent["eventColor"] = productionOrderInfo.ColorCode
				scheduledEvent["module"] = "production_order"
				scheduledEvent["componentName"] = "tooling_scheduled_order_event"
				serializedEventObject, _ := json.Marshal(scheduledEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject
}

func (v *ProductionOrderService) GetRejectQtyByProductionOrder(projectId string, productionOrderId int) int {
	allChildOrderCondition := " object_info->>'$.productionOrder' = " + strconv.Itoa(productionOrderId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfChildObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, allChildOrderCondition)
	if err != nil {
		return 0
	}

	var totalRejectQty int

	for _, scheduledOrderEventObject := range *listOfChildObjects {
		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledOrderEventObject.ObjectInfo}

		totalRejectQty += scheduledOrderEvent.getScheduledOrderEventInfo().RejectedQty
	}

	return totalRejectQty
}
func (v *ProductionOrderService) GetScheduledOrderInfo(projectId string, eventId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalScheduledOrderObject := Get(dbConnection, ScheduledOrderEventTable, eventId)

	return err, generalScheduledOrderObject
}

func (v *ProductionOrderService) GetAssemblyScheduledOrderInfo(projectId string, eventId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalScheduledOrderObject := Get(dbConnection, AssemblyScheduledOrderEventTable, eventId)

	return err, generalScheduledOrderObject
}

func (v *ProductionOrderService) GetToolingScheduledOrderInfo(projectId string, eventId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalScheduledOrderObject := Get(dbConnection, ToolingScheduledOrderEventTable, eventId)

	return err, generalScheduledOrderObject
}

func (v *ProductionOrderService) GetNextMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionQuery string
	if eventId == 0 {
		conditionQuery = " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	} else {
		conditionQuery = "  id > " + strconv.Itoa(eventId) + " and object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	}
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetNextAssemblyMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionQuery string
	if eventId == 0 {
		conditionQuery = " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	} else {
		conditionQuery = "  id > " + strconv.Itoa(eventId) + " and object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	}
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetFirstScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"

	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetFirstAssemblyScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"

	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetFirstToolingScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"

	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetLastScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id desc limit 1"

	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetLastAssemblyScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id desc limit 1"

	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetReleasedAssemblyOrderEventsBetween(projectId string, startDateTime string, endDateTime string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " STR_TO_DATE(object_info->>'$.startDate', '%Y-%m-%dT%H:%i:%s') >= '" + startDateTime + "' AND  STR_TO_DATE(object_info->>'$.endDate', '%Y-%m-%dT%H:%i:%s') <= '" + endDateTime + "' AND ( object_info->>'$.eventStatus' = " + strconv.Itoa(ScheduleStatusPreferenceFour) + " OR " + " object_info->>'$.eventStatus' = " + strconv.Itoa(ScheduleStatusPreferenceFive) + ")"
	v.BaseService.Logger.Info("getting assembly orders condition ", zap.String("condition", conditionQuery))
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetLastToolingScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionQuery := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id desc limit 1"

	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetPreviousMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionQuery string
	if eventId == 0 {
		conditionQuery = " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	} else {
		conditionQuery = "  id  < " + strconv.Itoa(eventId) + " and object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id desc limit 1"
	}
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetPreviousAssemblyMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionQuery string
	if eventId == 0 {
		conditionQuery = " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	} else {
		conditionQuery = "  id  < " + strconv.Itoa(eventId) + " and object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id desc limit 1"
	}
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetPreviousToolingMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionQuery string
	if eventId == 0 {
		conditionQuery = " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	} else {
		conditionQuery = "  id  < " + strconv.Itoa(eventId) + " and object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id desc limit 1"
	}
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) GetNextToolingMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionQuery string
	if eventId == 0 {
		conditionQuery = " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	} else {
		conditionQuery = "  id > " + strconv.Itoa(eventId) + " and object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " order by id asc limit 1"
	}
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionQuery)
	return err, listOfScheduledObjects
}

func (v *ProductionOrderService) UpdateOrderPreferenceLevel(projectId string, userId int, eventId int, preferenceLevel int) error {
	orderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, preferenceLevel)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, scheduledEventObject := Get(dbConnection, ScheduledOrderEventTable, eventId)
	schedulerEvent := ScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	scheduledOrderInfo := schedulerEvent.getScheduledOrderEventInfo()
	scheduledOrderInfo.EventStatus = orderStatusId
	return Update(dbConnection, ScheduledOrderEventTable, eventId, scheduledOrderInfo.DatabaseSerialize(userId))
}

func (v *ProductionOrderService) UpdateAssemblyOrderPreferenceLevel(projectId string, userId int, eventId int, preferenceLevel int) error {
	orderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, preferenceLevel)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, scheduledEventObject := Get(dbConnection, AssemblyScheduledOrderEventTable, eventId)
	schedulerEvent := ScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	scheduledOrderInfo := schedulerEvent.getScheduledOrderEventInfo()
	scheduledOrderInfo.EventStatus = orderStatusId
	v.BaseService.Logger.Info("updating UpdateAssemblyOrderPreferenceLevel", zap.Any("scheduler_order_event", scheduledOrderInfo), zap.Int("event_id", eventId))
	return Update(dbConnection, AssemblyScheduledOrderEventTable, eventId, scheduledOrderInfo.DatabaseSerialize(userId))
}

func (v *ProductionOrderService) UpdateToolingOrderPreferenceLevel(projectId string, userId int, eventId int, preferenceLevel int) error {
	orderStatusId := v.GetOrderStatusIdFromPreferenceLevel(projectId, preferenceLevel)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, scheduledEventObject := Get(dbConnection, ToolingScheduledOrderEventTable, eventId)
	schedulerEvent := ToolingScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	scheduledOrderInfo := schedulerEvent.getToolingScheduledOrderEventInfo()
	scheduledOrderInfo.EventStatus = orderStatusId
	return Update(dbConnection, ToolingScheduledOrderEventTable, eventId, scheduledOrderInfo.DatabaseSerialize(userId))
}

func (v *ProductionOrderService) UpdateScheduledOrderFields(projectId string, eventId int, updatingData datatypes.JSON) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, scheduledEventObject := Get(dbConnection, ScheduledOrderEventTable, eventId)
	var existingFields = make(map[string]interface{})
	json.Unmarshal(scheduledEventObject.ObjectInfo, &existingFields)
	var updateFields = make(map[string]interface{})
	json.Unmarshal(updatingData, &updateFields)
	for key, value := range updateFields {
		existingFields[key] = value
	}
	updatingObjectData := make(map[string]interface{})
	serializedObject, _ := json.Marshal(existingFields)
	updatingObjectData["object_info"] = serializedObject
	Update(dbConnection, ScheduledOrderEventTable, eventId, updatingObjectData)
	return nil
}

func (v *ProductionOrderService) UpdateAssemblyScheduledOrderFields(projectId string, eventId int, updatingData datatypes.JSON) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, scheduledEventObject := Get(dbConnection, AssemblyScheduledOrderEventTable, eventId)
	var existingFields = make(map[string]interface{})
	json.Unmarshal(scheduledEventObject.ObjectInfo, &existingFields)
	var updateFields = make(map[string]interface{})
	json.Unmarshal(updatingData, &updateFields)

	if eventStatus, ok := updateFields["eventStatus"]; ok {
		v.ViewManager.CreateLabourManagementShiftLinesHistory(eventId, util.InterfaceToInt(eventStatus))
	}

	for key, value := range updateFields {
		existingFields[key] = value
	}
	updatingObjectData := make(map[string]interface{})
	serializedObject, _ := json.Marshal(existingFields)
	updatingObjectData["object_info"] = serializedObject
	Update(dbConnection, AssemblyScheduledOrderEventTable, eventId, updatingObjectData)
	return nil
}
func (v *ProductionOrderService) GetAssemblyEventHistorySummary() datatypes.JSON {
	arrayEvent, err := v.ViewManager.GetEventStatusCounts()
	if err != nil {
		v.BaseService.Logger.Error("error getting assembly scheduler event count", zap.Error(err))
		return make([]byte, 0)
	}
	jsonData, err := json.Marshal(arrayEvent)
	if err != nil {
		v.BaseService.Logger.Error("error marshaling assembly scheduler event count", zap.Error(err))
		return make([]byte, 0)
	}
	return jsonData
}

func (v *ProductionOrderService) UpdateToolingScheduledOrderFields(projectId string, eventId int, updatingData datatypes.JSON) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	_, scheduledEventObject := Get(dbConnection, ToolingScheduledOrderEventTable, eventId)
	var existingFields = make(map[string]interface{})
	json.Unmarshal(scheduledEventObject.ObjectInfo, &existingFields)
	var updateFields = make(map[string]interface{})
	json.Unmarshal(updatingData, &updateFields)
	for key, value := range updateFields {
		existingFields[key] = value
	}
	updatingObjectData := make(map[string]interface{})
	serializedObject, _ := json.Marshal(existingFields)
	updatingObjectData["object_info"] = serializedObject
	Update(dbConnection, ToolingScheduledOrderEventTable, eventId, updatingObjectData)
	return nil
}

func (v *ProductionOrderService) GetCurrentToolingScheduledEvent(projectId string, machineId int) (error, component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionString)

	if err != nil {
		return err, component.GeneralObject{}
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), component.GeneralObject{}
	}
	for _, scheduledEvent := range *listOfObjects {
		scheduledOrderEvent := ToolingScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledEventInfo := scheduledOrderEvent.getToolingScheduledOrderEventInfo()

		if scheduledEventInfo.EventStatus == ScheduleStatusPreferenceEight || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceSeven {
			continue
		}

		startDate := util.ConvertStringToDateTime(scheduledEventInfo.StartDate)
		endDate := util.ConvertStringToDateTime(scheduledEventInfo.EndDate)

		if startDate.Error != nil || endDate.Error != nil {
			continue
		}

		timeNow := time.Now().Unix()

		if (startDate.DateTimeEpoch <= timeNow) && (endDate.DateTimeEpoch > timeNow) {
			return nil, scheduledEvent
		}
	}
	return errors.New("no scheduled events found"), component.GeneralObject{}
}

func (v *ProductionOrderService) GetCurrentScheduledEvent(projectId string, machineId int) (error, component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)

	if err != nil {
		return err, component.GeneralObject{}
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), component.GeneralObject{}
	}
	for _, scheduledEvent := range *listOfObjects {
		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

		if scheduledEventInfo.EventStatus == ScheduleStatusPreferenceEight || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceSeven {
			continue
		}

		startDate := util.ConvertStringToDateTime(scheduledEventInfo.StartDate)
		endDate := util.ConvertStringToDateTime(scheduledEventInfo.EndDate)

		if startDate.Error != nil || endDate.Error != nil {
			continue
		}

		timeNow := time.Now().Unix()
		if (startDate.DateTimeEpoch <= timeNow) && (endDate.DateTimeEpoch > timeNow) {
			return nil, scheduledEvent
		}
	}
	return errors.New("no scheduled events found"), component.GeneralObject{}
}

func (v *ProductionOrderService) GetGivenTimeIntervalScheduledEvent(projectId string, machineId int, energyStartDate string) (error, []component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)
	var listOfScheduleEvent = make([]component.GeneralObject, 0)

	if err != nil {
		return err, listOfScheduleEvent
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), listOfScheduleEvent
	}

	// Get the current time
	now := time.Now()

	// Get the start of the day (00:00:00)
	startOfDay := util.ConvertStringToDateTime(energyStartDate)

	// Convert the start of the day to a Unix timestamp
	startOfDayTimestamp := startOfDay.DateTimeEpoch

	// Get the end of the day (23:59:59)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Convert the end of the day to a Unix timestamp
	endOfDayTimestamp := endOfDay.Unix()
	for _, scheduledEvent := range *listOfObjects {
		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

		startDate := util.ConvertStringToDateTime(scheduledEventInfo.StartDate)
		endDate := util.ConvertStringToDateTime(scheduledEventInfo.EndDate)

		if startDate.Error != nil || endDate.Error != nil {
			continue
		}

		if (startDate.DateTimeEpoch >= startOfDayTimestamp) && (startDate.DateTimeEpoch < endOfDayTimestamp) {
			listOfScheduleEvent = append(listOfScheduleEvent, scheduledEvent)
		}
	}
	return errors.New("no scheduled events found"), listOfScheduleEvent
}

func (v *ProductionOrderService) GetCurrentAssemblyScheduledEvent(projectId string, assemblyMachineId int) (error, component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(assemblyMachineId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionString)
	if err != nil {
		return err, component.GeneralObject{}
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), component.GeneralObject{}
	}
	for _, scheduledEvent := range *listOfObjects {
		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

		if scheduledEventInfo.EventStatus == ScheduleStatusPreferenceEight || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceSeven {
			continue
		}

		startDate := util.ConvertStringToDateTime(scheduledEventInfo.StartDate)
		endDate := util.ConvertStringToDateTime(scheduledEventInfo.EndDate)

		if startDate.Error != nil || endDate.Error != nil {
			continue
		}

		timeNow := time.Now().Unix()
		if (startDate.DateTimeEpoch <= timeNow) && (endDate.DateTimeEpoch > timeNow) {
			return nil, scheduledEvent
		}
	}
	return errors.New("no scheduled events found"), component.GeneralObject{}
}

func (v *ProductionOrderService) GetNextToCurrentScheduledEvent(projectId string, machineId int, scheduledOrderId int) (error, component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' AND id > " + strconv.Itoa(scheduledOrderId) + " "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString, 1)
	if err != nil {
		return err, component.GeneralObject{}
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), component.GeneralObject{}
	}

	return nil, (*listOfObjects)[0]
}

func (v *ProductionOrderService) GetNextToCurrentAssemblyScheduledEvent(projectId string, machineId int, scheduledOrderId int) (error, component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' AND id > " + strconv.Itoa(scheduledOrderId) + " "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, AssemblyScheduledOrderEventTable, conditionString, 1)
	if err != nil {
		return err, component.GeneralObject{}
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), component.GeneralObject{}
	}

	return nil, (*listOfObjects)[0]
}

func (v *ProductionOrderService) GetNextToCurrentToolingScheduledEvent(projectId string, machineId int, scheduledOrderId int) (error, component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' AND id > " + strconv.Itoa(scheduledOrderId) + " "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionString, 1)
	if err != nil {
		return err, component.GeneralObject{}
	}
	if len(*listOfObjects) == 0 {
		return errors.New("no scheduled events found"), component.GeneralObject{}
	}

	return nil, (*listOfObjects)[0]
}

func (v *ProductionOrderService) GetProductionOrderStatus(projectId string, statusId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, statusGeneralObject := Get(dbConnection, ProductionOrderStatusTable, statusId)
	return err, statusGeneralObject
}

func (v *ProductionOrderService) GetSchedulerUpdatedDate(projectId string, mouldId int) (error, string) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.mouldId' =" + strconv.Itoa(mouldId)
	listOfProductionObjects, err := GetConditionalObjects(dbConnection, ProductionOrderMasterTable, conditionString)

	if err != nil || len(*listOfProductionObjects) == 0 {
		return getError("Record not found"), ""
	}

	productionObj := (*listOfProductionObjects)[0]

	schedulerConditionString := " object_info->>'$.eventSourceId' =" + strconv.Itoa(productionObj.Id)
	listOfSchedulerObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventComponent, schedulerConditionString)

	if err != nil || len(*listOfSchedulerObjects) == 0 {
		return getError("Record not found"), ""
	}

	for _, schedulerObj := range *listOfSchedulerObjects {
		scheduledOrderEventInfo := ScheduledOrderEventInfo{}
		json.Unmarshal(schedulerObj.ObjectInfo, &scheduledOrderEventInfo)

		if scheduledOrderEventInfo.LastUpdatedAt != "" {
			return nil, scheduledOrderEventInfo.LastUpdatedAt
		}
	}

	return nil, ""
}

func (v *ProductionOrderService) GetMachineProductionOrderInfo(projectId string, productionOrderId int, machineId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND id = " + strconv.Itoa(productionOrderId)
	listOfObjects, err := GetConditionalObjects(dbConnection, ProductionOrderMasterTable, conditionString)
	if len(*listOfObjects) == 1 {
		return nil, (*listOfObjects)[0]
	}
	return err, component.GeneralObject{}
}

func (v *ProductionOrderService) GetCompletedQuantity(projectId string, productionOrderId int) int {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.eventSourceId' =" + strconv.Itoa(productionOrderId)
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)

	if err != nil || len(*listOfObjects) == 0 {
		return 0
	}

	var totalCompletedQty = 0
	for _, schedulerObj := range *listOfObjects {
		var scheduledEvent = make(map[string]interface{})
		err = json.Unmarshal(schedulerObj.ObjectInfo, &scheduledEvent)
		if err == nil {
			if completedValue, ok := scheduledEvent["completedQty"]; ok {
				totalCompletedQty = totalCompletedQty + util.InterfaceToInt(completedValue)
			}
		}
	}
	return totalCompletedQty
}

func (v *ProductionOrderService) GetAssemblyProductionOrderInfo(projectId string, productionOrderId int, machineId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND id = " + strconv.Itoa(productionOrderId)
	listOfObjects, err := GetConditionalObjects(dbConnection, AssemblyProductionOrderTable, conditionString)
	if len(*listOfObjects) == 1 {
		return nil, (*listOfObjects)[0]
	}
	return err, component.GeneralObject{}
}

func (v *ProductionOrderService) GetAssemblyProductionOrderById(productionOrderId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	conditionString := " id = " + strconv.Itoa(productionOrderId)
	listOfObjects, err := GetConditionalObjects(dbConnection, AssemblyProductionOrderTable, conditionString)
	if len(*listOfObjects) == 1 {
		return nil, (*listOfObjects)[0]
	}
	return err, component.GeneralObject{}
}
func (v *ProductionOrderService) GetToolingProductionOrderInfo(projectId string, productionOrderId int, machineId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND id = " + strconv.Itoa(productionOrderId)
	listOfObjects, err := GetConditionalObjects(dbConnection, ToolingOrderMasterTable, conditionString)
	if len(*listOfObjects) == 1 {
		return nil, (*listOfObjects)[0]
	}
	return err, component.GeneralObject{}
}

func (v *ProductionOrderService) GetOrderStatusIdFromPreferenceLevel(projectId string, preferenceLevel int) int {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	return v.getOrderStatusId(dbConnection, preferenceLevel)
}

func (v *ProductionOrderService) OrderStatusId2String(projectId string, orderStatusId int) string {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, ProductionOrderStatusTable, orderStatusId)

	if err != nil || generalObject.Id == 0 {
		return ""
	}
	var productionStatusInfo ProductionOderStatusInfo
	err = json.Unmarshal(generalObject.ObjectInfo, &productionStatusInfo)

	if err != nil {
		return ""
	}
	return productionStatusInfo.Status
}

func (v *ProductionOrderService) GetBomInfo(projectId string, orderId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalOrderObject := Get(dbConnection, ToolingOrderMasterTable, orderId)

	var productionInfo map[string]interface{}
	err = json.Unmarshal(generalOrderObject.ObjectInfo, &productionInfo)

	bomId := util.InterfaceToInt(productionInfo["bomId"])

	err, generalBombject := Get(dbConnection, ToolingBomMasterTable, bomId)

	return err, generalBombject
}

func (v *ProductionOrderService) UpdateMouldScheduleOrderEventMouldBatchId(projectId string, scheduledOrderEventId int, mouldBatchResourceId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalOrderObject := Get(dbConnection, ScheduledOrderEventTable, scheduledOrderEventId)
	if err != nil {
		return err
	}

	var objectFields map[string]interface{}
	err = json.Unmarshal(generalOrderObject.ObjectInfo, &objectFields)

	objectFields["mouldBatchResourceId"] = mouldBatchResourceId
	serialisedData, _ := json.Marshal(objectFields)
	var updatingFields = make(map[string]interface{})
	updatingFields["object_info"] = serialisedData
	err = Update(dbConnection, ScheduledOrderEventTable, scheduledOrderEventId, updatingFields)
	if err != nil {
		return err
	}

	return err

}

func (v *ProductionOrderService) UpdateAssemblyManualOrderCompletedQuantity(serialisedData datatypes.JSON) (interface{}, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	err, updateAssemblyManualOrderQuantityRequest := rpc.GetUpdateAssemblyManualOrderQuantityRequest(serialisedData)
	if err != nil {
		return nil, err
	}
	var assemblyManualOrderHistoryInfo = model.AssemblyManualOrderCompletedQuantityHistoryInfo{
		LastUpdatedAt:     util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
		LastUpdatedBy:     updateAssemblyManualOrderQuantityRequest.RequestBy,
		CreatedAt:         util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
		CreatedBy:         updateAssemblyManualOrderQuantityRequest.RequestBy,
		ObjectStatus:      "Active",
		EventId:           updateAssemblyManualOrderQuantityRequest.EventId,
		CompletedQuantity: updateAssemblyManualOrderQuantityRequest.CompletedQuantity,
	}
	serialised, err := assemblyManualOrderHistoryInfo.Serialised()
	if err != nil {
		return nil, err
	}
	var generalObject = component.GeneralObject{
		ObjectInfo: serialised,
	}
	err, createdRecordId := Create(dbConnection, utils.AssemblyManualOrderCompletedQuantityHistoryTable, generalObject)
	if err != nil {
		return nil, err
	} else {
		v.BaseService.Logger.Debug("new history record is created successfully", zap.Int("record_id", createdRecordId))
	}

	return nil, err

}

func (v *ProductionOrderService) GetAssemblyManualOrderHistoryFromEventId(serialisedData datatypes.JSON) (interface{}, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	err, assemblyManualOrderHistoryByEventRequest := rpc.GetAssemblyManualOrderHistoryByEventRequest(serialisedData)
	if err != nil {
		return nil, err
	}
	var condition = " object_info->>'$.eventId' = " + strconv.Itoa(assemblyManualOrderHistoryByEventRequest.EventId)
	objects, err := GetConditionalObjects(dbConnection, utils.AssemblyManualOrderCompletedQuantityHistoryTable, condition)
	if err != nil {
		return nil, err
	}
	var historyList []rpc.AssemblyManualOrderHistoryByEventResponse
	for _, objectInterface := range *objects {
		err, a := model.GetAssemblyManualOrderCompletedQuantityHistoryInfo(objectInterface.ObjectInfo)
		if err == nil {
			historyList = append(historyList, rpc.AssemblyManualOrderHistoryByEventResponse{CreatedAt: a.CreatedAt, CompletedQuantity: a.CompletedQuantity})
		}
	}
	marshal, err := json.Marshal(historyList)
	if err != nil {
		return nil, err
	}

	return marshal, err

}

func (v *ProductionOrderService) GetNumberOfOdersByStatus(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	now := time.Now()
	utcFormat := "2006-01-02T15:04:05.000Z"
	sgt, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		v.BaseService.Logger.Error("error loading SGT time zone:", zap.Error(err))
		return err, nil
	}

	referenceDateTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 15, 0, 0, time.Local)
	var startDate, nextDate time.Time

	if now.Before(referenceDateTime) {
		startDate = referenceDateTime.AddDate(0, 0, -1)
		nextDate = referenceDateTime
	} else {
		startDate = referenceDateTime
		nextDate = referenceDateTime.AddDate(0, 0, 1)
	}
	startDateStringType := startDate.Format(utcFormat)
	nextDateStringType := nextDate.Format(utcFormat)
	// conditionQuery := fmt.Sprintf("STR_TO_DATE(object_info->>'$.startDate', '%%Y-%%m-%%dT%%H:%%i:%%s.000Z') BETWEEN '%s' AND '%s'",
	// 	startDateStringType, nextDateStringType,
	// )

	// conditionQuery := " DATE(STR_TO_DATE(object_info->>'$. ', '%Y-%m-%dT%H:%i:%s')) = '" + currentDate + "'"

	conditionQuery := " object_info->> '$.objectStatus' = 'Active'"
	v.BaseService.Logger.Info("getting schedule order events condition ", zap.String("condition", conditionQuery))
	listOfScheduledObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionQuery)
	if err != nil {
		v.BaseService.Logger.Error("error getting schedule list:", zap.Error(err))
		return err, nil
	}

	startDateTime, _ := time.Parse(utcFormat, startDateStringType)
	nextDateTime, _ := time.Parse(utcFormat, nextDateStringType)

	var filteredObjects []component.GeneralObject

	for _, scheduledEvent := range *listOfScheduledObjects {

		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledOrderInfo := scheduledOrderEvent.getScheduledOrderEventInfo()
		utcTime, err := time.Parse(utcFormat, scheduledOrderInfo.StartDate)
		if err != nil {
			v.BaseService.Logger.Error("error parsing UTC time:", zap.Error(err))
			continue
		}

		sgtTime := utcTime.In(sgt)
		sgtFormatted := sgtTime.Format(utcFormat)

		sgtTimeParsed, _ := time.Parse(utcFormat, sgtFormatted)
		if (sgtTimeParsed.After(startDateTime) || sgtTimeParsed.Equal(startDateTime)) && sgtTimeParsed.Before(nextDateTime) {
			filteredObjects = append(filteredObjects, scheduledEvent)
		}

	}

	return err, &filteredObjects
}

func (v *ProductionOrderService) GetTotalCompletedQuantity(projectId string, mouldId int, sincetime string) int64 {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var completedQty int64
	var query = "SELECT COALESCE(sum(object_info->>'$.completedQty'),0)  as completedQty as completedQty FROM scheduled_order_event  WHERE object_info->>'$.mouldId' =" + strconv.Itoa(mouldId) + " and  object_info->>'$.lastUpdatedAt' >" + sincetime
	result := dbConnection.Raw(query).Scan(&completedQty)

	if result.Error != nil {
		v.BaseService.Logger.Error("error running query to get the total quantity", zap.Error(result.Error))
		return 0
	} else {
		v.BaseService.Logger.Info("got the total completed quantity since", zap.String("since_time", sincetime), zap.Any("value", completedQty))
	}

	return completedQty
}

func (v *ProductionOrderService) GetCurrentScheduledEventByMachineId(projectId string, machineId int, testStartDate string, testEndDate string) (error, *[]component.GeneralObject) {
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventComponent, conditionString)
	var filteredObjects []component.GeneralObject
	// timeNow := time.Now().Unix()
	requestStartDate := util.ConvertStringToDateTime(testStartDate)
	requestEndDate := util.ConvertStringToDateTime(testEndDate)

	if err != nil {
		v.BaseService.Logger.Error("error running query to get the total quantity", zap.Error(err))
		return err, nil
	}
	if len(*listOfObjects) == 0 {
		v.BaseService.Logger.Error("no scheduled events found")
		return errors.New("no scheduled events found"), nil
	}

	for _, scheduledEvent := range *listOfObjects {
		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

		if scheduledEventInfo.EventStatus == ScheduleStatusPreferenceThree || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceFour || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceFive {
			startDate := util.ConvertStringToDateTime(scheduledEventInfo.StartDate)
			endDate := util.ConvertStringToDateTime(scheduledEventInfo.EndDate)

			if startDate.Error != nil || endDate.Error != nil {
				continue
			}

			if (startDate.DateTimeEpoch <= requestEndDate.DateTimeEpoch) && (endDate.DateTimeEpoch >= requestStartDate.DateTimeEpoch) {
				filteredObjects = append(filteredObjects, scheduledEvent)
			}
		}

	}
	return nil, &filteredObjects
}

func (v *ProductionOrderService) GetCurrentScheduledEventByMouldId(projectId string, mouldId int) (error, *[]component.GeneralObject) {
	conditionString := " object_info->>'$.mouldId' =" + strconv.Itoa(mouldId) + " AND object_info->>'$.objectStatus' = 'Active' "
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventComponent, conditionString)
	var filteredObjects []component.GeneralObject
	timeNow := time.Now().Unix()

	if err != nil {
		v.BaseService.Logger.Error("error running query to get the total quantity", zap.Error(err))
		return err, nil
	}
	if len(*listOfObjects) == 0 {
		v.BaseService.Logger.Error("no scheduled events found")
		return errors.New("no scheduled events found"), nil
	}

	for _, scheduledEvent := range *listOfObjects {
		scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}
		scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

		if scheduledEventInfo.EventStatus == ScheduleStatusPreferenceThree || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceFour || scheduledEventInfo.EventStatus == ScheduleStatusPreferenceFive {
			startDate := util.ConvertStringToDateTime(scheduledEventInfo.StartDate)
			endDate := util.ConvertStringToDateTime(scheduledEventInfo.EndDate)

			if startDate.Error != nil || endDate.Error != nil {
				continue
			}

			if (startDate.DateTimeEpoch <= timeNow) && (endDate.DateTimeEpoch > timeNow) {
				filteredObjects = append(filteredObjects, scheduledEvent)
			}
		}

	}
	return nil, &filteredObjects
}
