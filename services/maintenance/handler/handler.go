package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MaintenanceService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
	ToolingSupervisorGroup  int
}

func (v *MaintenanceService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MaintenanceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MaintenanceComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *MaintenanceService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MaintenanceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MaintenanceComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}
func (v *MaintenanceService) GetMaintenanceEventsForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info ->> '$.moduleName' = 'machines' and object_info->>'$.objectStatus' = 'Active' "
	listOfEvents, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderStatusCache = make(map[int]*MaintenanceWorkOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceStatus := MaintenanceWorkOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderStatusCache[statusGeneralObject.Id] = maintenanceStatus.getMaintenanceWorkOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMaintenanceOrder(scheduledEvent)

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false

		if value, ok := scheduledEvent["workOrderStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if maintenanceWorkOrderStatusInfo, ok := maintenanceWorkOrderStatusCache[eventStatusId]; ok {
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canRelease"] = true
				} else {
					composedScheduleEvent["canRelease"] = false
				}
				if eventStatusId == WorkOrderScheduled {
					composedScheduleEvent["canHold"] = true
				} else {
					composedScheduleEvent["canHold"] = false
				}
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canUpdate"] = true
				} else {
					composedScheduleEvent["canUpdate"] = false
				}

				composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderStatusInfo.Status
				composedScheduleEvent["eventColor"] = maintenanceWorkOrderStatusInfo.ColorCode
				composedScheduleEvent["eventSourceId"] = scheduledEventObject.Id
				composedScheduleEvent["machineId"] = util.InterfaceToInt(scheduledEvent["assetId"])
				serializedEventObject, _ := json.Marshal(composedScheduleEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject

}

func (v *MaintenanceService) GetMaintenanceEventsForAssemblyScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info ->> '$.moduleName' = 'assemblyMachines' "
	listOfEvents, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderStatusCache = make(map[int]*MaintenanceWorkOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceStatus := MaintenanceWorkOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderStatusCache[statusGeneralObject.Id] = maintenanceStatus.getMaintenanceWorkOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMaintenanceOrder(scheduledEvent)

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false

		if value, ok := scheduledEvent["workOrderStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if maintenanceWorkOrderStatusInfo, ok := maintenanceWorkOrderStatusCache[eventStatusId]; ok {
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canRelease"] = true
				} else {
					composedScheduleEvent["canRelease"] = false
				}
				if eventStatusId == WorkOrderScheduled {
					composedScheduleEvent["canHold"] = true
				} else {
					composedScheduleEvent["canHold"] = false
				}
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canUpdate"] = true
				} else {
					composedScheduleEvent["canUpdate"] = false
				}

				composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderStatusInfo.Status
				composedScheduleEvent["eventColor"] = maintenanceWorkOrderStatusInfo.ColorCode
				composedScheduleEvent["eventSourceId"] = scheduledEventObject.Id
				composedScheduleEvent["machineId"] = util.InterfaceToInt(scheduledEvent["assetId"])
				serializedEventObject, _ := json.Marshal(composedScheduleEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject

}

func (v *MaintenanceService) GetMaintenanceEventsForToolingScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info ->> '$.moduleName' = 'toolingMachines' "
	listOfEvents, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderStatusCache = make(map[int]*MaintenanceWorkOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceStatus := MaintenanceWorkOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderStatusCache[statusGeneralObject.Id] = maintenanceStatus.getMaintenanceWorkOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMaintenanceOrder(scheduledEvent)

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false

		if value, ok := scheduledEvent["workOrderStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if maintenanceWorkOrderStatusInfo, ok := maintenanceWorkOrderStatusCache[eventStatusId]; ok {
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canRelease"] = true
				} else {
					composedScheduleEvent["canRelease"] = false
				}
				if eventStatusId == WorkOrderScheduled {
					composedScheduleEvent["canHold"] = true
				} else {
					composedScheduleEvent["canHold"] = false
				}
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canUpdate"] = true
				} else {
					composedScheduleEvent["canUpdate"] = false
				}

				composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderStatusInfo.Status
				composedScheduleEvent["eventColor"] = maintenanceWorkOrderStatusInfo.ColorCode
				composedScheduleEvent["eventSourceId"] = scheduledEventObject.Id
				composedScheduleEvent["machineId"] = util.InterfaceToInt(scheduledEvent["assetId"])
				serializedEventObject, _ := json.Marshal(composedScheduleEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject

}

func composeEventByMaintenanceOrder(scheduledEventObject map[string]interface{}) map[string]interface{} {
	createEvent := make(map[string]interface{})

	createEvent["name"] = util.InterfaceToString(scheduledEventObject["description"])

	createEvent["iconCls"] = "fa fa-cogs"
	createEvent["eventType"] = "maintenance_work_order"

	createEvent["draggable"] = true
	createEvent["module"] = "maintenance"
	createEvent["componentName"] = "maintenance_work_order"
	createEvent["startDate"] = util.InterfaceToString(scheduledEventObject["workOrderScheduledStartDate"])
	createEvent["endDate"] = util.InterfaceToString(scheduledEventObject["workOrderScheduledEndDate"])
	createEvent["machineId"] = util.InterfaceToInt(scheduledEventObject["assetId"])
	createEvent["objectStatus"] = "Active"
	createEvent["canComplete"] = util.InterfaceToBool(scheduledEventObject["canComplete"])
	createEvent["workOrderStatus"] = util.InterfaceToInt(scheduledEventObject["workOrderStatus"])
	createEvent["canRelease"] = util.InterfaceToBool(scheduledEventObject["canRelease"])
	createEvent["canUnRelease"] = util.InterfaceToBool(scheduledEventObject["canUnRelease"])

	return createEvent
}

func composeEventByMouldMaintenanceOrder(scheduledEventObject map[string]interface{}) map[string]interface{} {
	createEvent := make(map[string]interface{})

	createEvent["name"] = util.InterfaceToString(scheduledEventObject["description"])

	createEvent["iconCls"] = "fa fa-cogs"
	createEvent["eventType"] = "maintenance_work_order"

	createEvent["draggable"] = true
	createEvent["module"] = "maintenance"
	createEvent["componentName"] = "mould_preventive_maintenance_work_order"
	createEvent["startDate"] = util.InterfaceToString(scheduledEventObject["workOrderScheduledStartDate"])
	createEvent["endDate"] = util.InterfaceToString(scheduledEventObject["workOrderScheduledEndDate"])
	createEvent["mouldId"] = util.InterfaceToInt(scheduledEventObject["mouldId"])
	createEvent["objectStatus"] = "Active"
	createEvent["canComplete"] = util.InterfaceToBool(scheduledEventObject["canComplete"])
	createEvent["workOrderStatus"] = util.InterfaceToInt(scheduledEventObject["workOrderStatus"])
	createEvent["canRelease"] = util.InterfaceToBool(scheduledEventObject["canRelease"])
	createEvent["canUnRelease"] = util.InterfaceToBool(scheduledEventObject["canUnRelease"])

	return createEvent
}

func (v *MaintenanceService) GetMaintenanceOrderForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info ->> '$.objectStatus' = 'Active'"
	listOfEvents, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderTaskStatusCache = make(map[int]*MaintenanceWorkOrderTaskStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceTaskStatus := MaintenanceWorkOrderTaskStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderTaskStatusCache[statusGeneralObject.Id] = maintenanceTaskStatus.getMaintenanceWorkOrderTaskStatusInfo()
	}

	var scheduledEvent = make(map[string]interface{})
	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMaintenanceOrder(scheduledEvent)

		eventStatusId := util.InterfaceToInt(composedScheduleEvent["workOrderStatus"])
		maintenanceWorkOrderTaskStatusInfo := maintenanceWorkOrderTaskStatusCache[eventStatusId]
		//statusPreference := util.InterfaceToInt(maintenanceWorkOrderTaskStatusInfo.Preference)
		//if statusPreference == ScheduleStatusPreferenceThree {
		//	composedScheduleEvent["canRelease"] = true
		//} else {
		//	composedScheduleEvent["canRelease"] = false
		//}
		//if statusPreference == ScheduleStatusPreferenceFour {
		//	composedScheduleEvent["canHold"] = true
		//} else {
		//	composedScheduleEvent["canHold"] = false
		//}

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false

		composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderTaskStatusInfo.Status
		composedScheduleEvent["eventColor"] = maintenanceWorkOrderTaskStatusInfo.ColorCode
		composedScheduleEvent["componentName"] = MaintenanceWorkOrderTable
		composedScheduleEvent["canUpdate"] = false

		serializedEventObject, _ := json.Marshal(composedScheduleEvent)
		arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
	}

	return err, &arrayOfStatusObject

}

func (v *MaintenanceService) GetMaintenanceCorrectiveOrderForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info ->> '$.moduleName' = 'machines' AND object_info ->> '$.objectStatus' = 'Active' "
	listOfEvents, err := GetConditionalObjects(dbConnection, MaintenanceCorrectiveWorkOrderComponent, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderStatusCache = make(map[int]*MaintenanceWorkOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceStatus := MaintenanceWorkOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderStatusCache[statusGeneralObject.Id] = maintenanceStatus.getMaintenanceWorkOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMouldMaintenanceOrder(scheduledEvent)

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false
		composedScheduleEvent["componentName"] = MaintenanceCorrectiveWorkOrderComponent

		if value, ok := scheduledEvent["workOrderStatus"]; ok {
			eventStatusId := util.InterfaceToInt(value)

			if maintenanceWorkOrderStatusInfo, ok := maintenanceWorkOrderStatusCache[eventStatusId]; ok {
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canRelease"] = true
				} else {
					composedScheduleEvent["canRelease"] = false
				}
				if eventStatusId == WorkOrderScheduled {
					composedScheduleEvent["canHold"] = true
				} else {
					composedScheduleEvent["canHold"] = false
				}
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canUpdate"] = true
				} else {
					composedScheduleEvent["canUpdate"] = false
				}

				composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderStatusInfo.Status
				composedScheduleEvent["eventColor"] = maintenanceWorkOrderStatusInfo.ColorCode
				composedScheduleEvent["eventSourceId"] = scheduledEventObject.Id
				composedScheduleEvent["machineId"] = util.InterfaceToInt(scheduledEvent["assetId"])
				serializedEventObject, _ := json.Marshal(composedScheduleEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject

}

func (v *MaintenanceService) GetMouldPreventiveMaintenanceOrderForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	conditionString := " object_info ->> '$.objectStatus' = 'Active'"
	listOfEvents, err := GetConditionalObjects(dbConnection, MouldMaintenancePreventiveWorkOrderTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderTaskStatusCache = make(map[int]*MaintenanceWorkOrderTaskStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceTaskStatus := MaintenanceWorkOrderTaskStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderTaskStatusCache[statusGeneralObject.Id] = maintenanceTaskStatus.getMaintenanceWorkOrderTaskStatusInfo()
	}

	var scheduledEvent = make(map[string]interface{})
	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMouldMaintenanceOrder(scheduledEvent)

		eventStatusId := util.InterfaceToInt(composedScheduleEvent["workOrderStatus"])
		maintenanceWorkOrderTaskStatusInfo := maintenanceWorkOrderTaskStatusCache[eventStatusId]

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MouldMaintenancePreventiveWorkOrderTaskTable, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false

		composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderTaskStatusInfo.Status
		composedScheduleEvent["eventColor"] = maintenanceWorkOrderTaskStatusInfo.ColorCode
		composedScheduleEvent["componentName"] = MouldMaintenancePreventiveWorkOrderTable
		composedScheduleEvent["canUpdate"] = false

		serializedEventObject, _ := json.Marshal(composedScheduleEvent)
		arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
	}

	return err, &arrayOfStatusObject

}

func (v *MaintenanceService) GetMouldCorrectiveMaintenanceOrderForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	// conditionString := " object_info ->> '$.moduleName' = 'machines' AND object_info ->> '$.objectStatus' = 'Active' "
	// conditionString := " object_info ->> '$.moduleName' = 'moulds' AND object_info ->> '$.objectStatus' = 'Active' "
	conditionString := "(object_info ->> '$.moduleName' = 'machines' OR object_info ->> '$.moduleName' = 'moulds') AND object_info ->> '$.objectStatus' = 'Active'"

	listOfEvents, err := GetConditionalObjects(dbConnection, MouldMaintenanceCorrectiveWorkOrderTable, conditionString)
	listOfStatusCode, err := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)

	var maintenanceWorkOrderStatusCache = make(map[int]*MaintenanceWorkOrderStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		maintenanceStatus := MaintenanceWorkOrderStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		maintenanceWorkOrderStatusCache[statusGeneralObject.Id] = maintenanceStatus.getMaintenanceWorkOrderStatusInfo()
	}

	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		var scheduledEvent = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMaintenanceOrder(scheduledEvent)

		// get all the tasks and see how much completed
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(scheduledEventObject.Id) + "\""
		listOfTasks, _ := GetConditionalObjects(dbConnection, MouldMaintenanceCorrectiveWorkOrderTaskTable, condition)
		var totalTasks = len(*listOfTasks)
		var completedTasks = 0
		var todoTask = 0
		var inProgressTask = 0
		for _, taskGeneralObjectInterface := range *listOfTasks {
			workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: taskGeneralObjectInterface.ObjectInfo}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
				completedTasks = completedTasks + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskToDo {
				todoTask = todoTask + 1
			}
			if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskInProgress {
				inProgressTask = inProgressTask + 1
			}
		}
		if totalTasks == 0 {
			composedScheduleEvent["percentDone"] = 0
		} else {
			composedScheduleEvent["percentDone"] = (completedTasks / totalTasks) * 100
		}

		composedScheduleEvent["noOfToDoTasks"] = todoTask
		composedScheduleEvent["noOfTotalTasks"] = totalTasks
		composedScheduleEvent["noOfCompletedTasks"] = completedTasks
		composedScheduleEvent["noOfInProgressTasks"] = inProgressTask
		composedScheduleEvent["isAbortEnabled"] = false
		composedScheduleEvent["componentName"] = MouldMaintenanceCorrectiveWorkOrderTable

		if value, ok := scheduledEvent["workOrderStatus"]; ok {

			eventStatusId := util.InterfaceToInt(value)

			if maintenanceWorkOrderStatusInfo, ok := maintenanceWorkOrderStatusCache[eventStatusId]; ok {
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canRelease"] = true
				} else {
					composedScheduleEvent["canRelease"] = false
				}
				if eventStatusId == WorkOrderScheduled {
					composedScheduleEvent["canHold"] = true
				} else {
					composedScheduleEvent["canHold"] = false
				}
				if eventStatusId == WorkOrderCreated {
					composedScheduleEvent["canUpdate"] = true
				} else {
					composedScheduleEvent["canUpdate"] = false
				}

				composedScheduleEvent["eventStatusName"] = maintenanceWorkOrderStatusInfo.Status
				composedScheduleEvent["eventColor"] = maintenanceWorkOrderStatusInfo.ColorCode
				composedScheduleEvent["eventSourceId"] = scheduledEventObject.Id
				composedScheduleEvent["machineId"] = util.InterfaceToInt(scheduledEvent["assetId"])
				composedScheduleEvent["mouldId"] = util.InterfaceToInt(scheduledEvent["mouldId"])
				serializedEventObject, _ := json.Marshal(composedScheduleEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}
		}

	}
	return err, &arrayOfStatusObject

}

func (v *MaintenanceService) GetMaintenanceWorkOrderByMachineId(machineId int, table string) (error, string) { //err, start date, end date, description
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	condition := "object_info ->> '$.assetId' = " + strconv.Itoa(machineId) + " AND object_info ->> '$.moduleName' = 'machines' AND CURDATE() BETWEEN DATE(JSON_UNQUOTE(object_info->> '$.workOrderScheduledStartDate')) AND DATE(JSON_UNQUOTE(object_info->> '$.workOrderScheduledEndDate')) ORDER BY id DESC"
	generalObject, err := GetConditionalObjects(dbConnection, table, condition, 1)

	if err != nil {
		return err, ""
	}
	if generalObject == nil || len(*generalObject) == 0 {
		return fmt.Errorf("no data found"), ""
	}
	var scheduledEvent = make(map[string]interface{})
	json.Unmarshal((*generalObject)[0].ObjectInfo, &scheduledEvent)
	correctiveWorkOrder := composeEventByMouldMaintenanceOrder(scheduledEvent)

	name := correctiveWorkOrder["name"].(string)
	startOrderDate := util.InterfaceToString(correctiveWorkOrder["startDate"])
	endOrderDate := util.InterfaceToString(correctiveWorkOrder["endDate"])

	startDate := util.ConvertStringToDateTime(startOrderDate)
	endDate := util.ConvertStringToDateTime(endOrderDate)

	timeNow := time.Now().Unix()

	if (startDate.DateTimeEpoch <= timeNow) && (endDate.DateTimeEpoch > timeNow) {

		return nil, name
	}
	return err, ""

}

func getAssetIdFromMaintenanceList(maintenanceOrderId int, maintenanceOrderList *[]component.GeneralObject) (int, bool) {
	machineId := 0
	canRelease := false
	var orderObject = make(map[string]interface{})
	for _, maintenanceOrderObject := range *maintenanceOrderList {
		if maintenanceOrderObject.Id == maintenanceOrderId {
			json.Unmarshal(maintenanceOrderObject.ObjectInfo, &orderObject)
			return util.InterfaceToInt(orderObject["assetId"]), util.InterfaceToBool(orderObject["canRelease"])
		}
	}
	return machineId, canRelease
}

func composeEventByMaintenanceOrderTask(scheduledEventObject map[string]interface{}) map[string]interface{} {
	createEvent := make(map[string]interface{})

	createEvent["name"] = util.InterfaceToString(scheduledEventObject["taskName"])
	createEvent["iconCls"] = "fa fa-cogs"
	createEvent["eventType"] = "maintenance_work_order_task"

	createEvent["draggable"] = true
	createEvent["module"] = "maintenance"
	createEvent["componentName"] = "maintenance_work_order_task"
	createEvent["startDate"] = util.InterfaceToString(scheduledEventObject["taskDate"])
	createEvent["endDate"] = util.InterfaceToString(scheduledEventObject["estimatedTaskEndDate"])
	createEvent["assignedUserId"] = util.InterfaceToInt(scheduledEventObject["assignedUserId"])

	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0

	return createEvent
}

func (v *MaintenanceService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	maintenanceGeneral := routerEngine.Group("/project/:projectId/maintenance")
	maintenanceGeneral.POST("/loadFile", v.loadFile)
	maintenanceGeneral.GET("/maintenance_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getMaintenanceOverview)
	maintenanceGeneral.GET("/scheduler_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSchedulerView)
	maintenanceGeneral.GET("/work_order_task_kanban_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getMaintenanceWorkOrderTaskKanbanView)
	maintenance := routerEngine.Group("/project/:projectId/maintenance/component/:componentName")

	// table component  requests
	maintenance.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getObjects)
	maintenance.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardView)
	maintenance.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNewRecord)
	maintenance.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getRecordFormData)
	maintenance.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateResource)
	maintenance.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteResource)
	maintenance.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteValidation)
	maintenance.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.createNewResource)
	maintenance.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.importObjects)
	maintenance.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTableImportSchema)
	maintenance.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.exportObjects)
	maintenance.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSearchResults)
	maintenance.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getExportSchema)
	maintenance.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(ModuleName), v.removeInternalArrayReference)
	//get the records

	maintenance.POST("/record/:recordId/action/kanban_move_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.kanbanMoveTask)

	maintenance.POST("/record/:recordId/action/update_scheduler_event", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateSchedulerEvent)
	maintenance.POST("/record/:recordId/action/release", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.releaseWorkOrder)
	// maintenance.POST("/record/:recordId/action/complete_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ms.completeTask)
	maintenance.POST("/record/:recordId/action/checkout_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.completeTask)
	maintenance.POST("/record/:recordId/action/checkin_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.checkInTask)

	maintenance.POST("/records/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getGroupBy)
	maintenance.POST("/card_view/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardViewGroupBy)

	maintenance.POST("/record/:recordId/action/approve_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.approveTask)
	maintenance.POST("/record/:recordId/action/reject", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.rejectTask)

	maintenance.POST("/record/:recordId/action/hold", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.holdOrder)
	maintenance.POST("/record/:recordId/action/complete", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.completeOrder)
	maintenance.POST("/record/:recordId/action/force_complete", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.forceCompleteOrder)

	maintenance.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), v.getComponentRecordTrails)

}
