package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
)

func (v *MaintenanceService) getSummaryResponse(projectId string) map[string]interface{} {
	var response = make(map[string]interface{})
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	allMaintenanceOrder, _ := GetObjects(dbConnection, MaintenanceWorkOrderTable)
	allCorrectiveMaintenanceOrder, _ := GetObjects(dbConnection, MaintenanceCorrectiveWorkOrderTable)
	allMaintenanceStatus, _ := GetObjects(dbConnection, MaintenanceWorkOrderStatusTable)
	allMaintenanceOrderTask, _ := GetObjects(dbConnection, MaintenanceWorkOrderTaskTable)
	allCorrectiveMaintenanceOrderTask, _ := GetObjects(dbConnection, MaintenanceWorkOrderCorrectiveTaskTable)
	allMaintenanceOrderTaskStatus, _ := GetObjects(dbConnection, MaintenanceWorkOrderTaskStatusTable)

	statusIdMap := make(map[int]string)
	correctiveWorkOrderStatusCount := make(map[int]int)
	for _, status := range *allMaintenanceStatus {
		maintenanceWorkOrderStatusInfo := MaintenanceWorkOrderStatusInfo{}
		json.Unmarshal(status.ObjectInfo, &maintenanceWorkOrderStatusInfo)

		statusIdMap[status.Id] = maintenanceWorkOrderStatusInfo.Status
		correctiveWorkOrderStatusCount[status.Id] = 0
	}

	preventiveWorkOrderStatusCount := make(map[int]int)
	for _, status := range *allMaintenanceStatus {
		maintenanceWorkOrderStatusInfo := MaintenanceWorkOrderStatusInfo{}
		json.Unmarshal(status.ObjectInfo, &maintenanceWorkOrderStatusInfo)

		statusIdMap[status.Id] = maintenanceWorkOrderStatusInfo.Status
		preventiveWorkOrderStatusCount[status.Id] = 0
	}

	correctiveWorkOrderTaskStatusCount := make(map[int]int)
	for _, taskStatus := range *allMaintenanceOrderTaskStatus {
		maintenanceWorkOrderTaskStatusInfo := MaintenanceWorkOrderTaskStatusInfo{}
		json.Unmarshal(taskStatus.ObjectInfo, &maintenanceWorkOrderTaskStatusInfo)

		statusIdMap[taskStatus.Id] = maintenanceWorkOrderTaskStatusInfo.Status
		correctiveWorkOrderTaskStatusCount[taskStatus.Id] = 0
	}

	preventiveWorkOrderTaskStatusCount := make(map[int]int)
	for _, taskStatus := range *allMaintenanceOrderTaskStatus {
		maintenanceWorkOrderTaskStatusInfo := MaintenanceWorkOrderTaskStatusInfo{}
		json.Unmarshal(taskStatus.ObjectInfo, &maintenanceWorkOrderTaskStatusInfo)

		statusIdMap[taskStatus.Id] = maintenanceWorkOrderTaskStatusInfo.Status
		preventiveWorkOrderTaskStatusCount[taskStatus.Id] = 0
	}

	noPreventiveTypeOrders := 0
	noCorrectiveTypeOrders := 0
	noPreventiveTypeOrderTasks := 0
	noCorrectiveTypeOrdersTasks := 0

	filterOverview := make([]interface{}, 0)

	//Filter for machine and mould
	machineMaintenanceStatusOverview := MaintenanceStatusOverview{}
	machineMaintenanceStatusOverview.GroupByField = "MACHINE ASSETS"
	machineGroup := make([]map[string]interface{}, 0)

	mouldMaintenanceStatusOverview := MaintenanceStatusOverview{}
	mouldMaintenanceStatusOverview.GroupByField = "MOULD ASSETS"
	mouldGroup := make([]map[string]interface{}, 0)

	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

	for _, order := range *allMaintenanceOrder {
		filteredAssetInfo := make(map[string]interface{})
		workOrderInfo := WorkOrderInfo{}
		json.Unmarshal(order.ObjectInfo, &workOrderInfo)

		noPreventiveTypeOrders += 1

		if _, ok := preventiveWorkOrderStatusCount[workOrderInfo.WorkOrderStatus]; ok {
			preventiveWorkOrderStatusCount[workOrderInfo.WorkOrderStatus] += 1
		}

		filteredAssetInfo["status"] = statusIdMap[workOrderInfo.WorkOrderStatus]
		filteredAssetInfo["startDate"] = workOrderInfo.CreatedAt

		assetInfo := make(map[string]interface{})
		if workOrderInfo.ModuleName == "machines" {
			_, machineData := machineService.GetMachineInfoById(projectId, workOrderInfo.AssetId)
			json.Unmarshal(machineData.ObjectInfo, &assetInfo)

			filteredAssetInfo["toolNo"] = util.InterfaceToString(assetInfo["newMachineId"])
			filteredAssetInfo["imageUrl"] = util.InterfaceToString(assetInfo["machineImage"])

			machineGroup = append(machineGroup, filteredAssetInfo)

		} else {
			_, mouldData := mouldService.GetMouldInfoById(projectId, workOrderInfo.AssetId)
			json.Unmarshal(mouldData.ObjectInfo, &assetInfo)

			filteredAssetInfo["toolNo"] = util.InterfaceToString(assetInfo["toolNo"])
			filteredAssetInfo["imageUrl"] = util.InterfaceToString(assetInfo["mouldImage"])

			mouldGroup = append(mouldGroup, filteredAssetInfo)
		}

	}

	// This is for corrective work order
	for _, order := range *allCorrectiveMaintenanceOrder {
		filteredAssetInfo := make(map[string]interface{})
		correctiveWorkOrderInfo := make(map[string]interface{})
		json.Unmarshal(order.ObjectInfo, &correctiveWorkOrderInfo)

		noCorrectiveTypeOrders += 1

		if _, ok := correctiveWorkOrderStatusCount[util.InterfaceToInt(correctiveWorkOrderInfo["workOrderStatus"])]; ok {
			correctiveWorkOrderStatusCount[util.InterfaceToInt(correctiveWorkOrderInfo["workOrderStatus"])] += 1
		}

		filteredAssetInfo["status"] = statusIdMap[util.InterfaceToInt(correctiveWorkOrderInfo["workOrderStatus"])]
		filteredAssetInfo["startDate"] = util.InterfaceToString(correctiveWorkOrderInfo["createdAt"])

		assetInfo := make(map[string]interface{})

		_, machineData := machineService.GetMachineInfoById(projectId, util.InterfaceToInt(correctiveWorkOrderInfo["assetId"]))
		json.Unmarshal(machineData.ObjectInfo, &assetInfo)

		filteredAssetInfo["toolNo"] = util.InterfaceToString(assetInfo["newMachineId"])
		filteredAssetInfo["imageUrl"] = util.InterfaceToString(assetInfo["machineImage"])

		machineGroup = append(machineGroup, filteredAssetInfo)

	}

	machineMaintenanceStatusOverview.Cards = machineGroup
	mouldMaintenanceStatusOverview.Cards = mouldGroup

	filterOverview = append(filterOverview, machineMaintenanceStatusOverview)
	filterOverview = append(filterOverview, mouldMaintenanceStatusOverview)

	for _, task := range *allMaintenanceOrderTask {
		workOrderTaskInfo := WorkOrderTaskInfo{}
		json.Unmarshal(task.ObjectInfo, &workOrderTaskInfo)
		noPreventiveTypeOrderTasks += 1
		if _, ok := preventiveWorkOrderTaskStatusCount[workOrderTaskInfo.TaskStatus]; ok {
			preventiveWorkOrderTaskStatusCount[workOrderTaskInfo.TaskStatus] += 1
		}

	}

	// For corrective work order task
	for _, task := range *allCorrectiveMaintenanceOrderTask {
		workOrderTaskInfo := make(map[string]interface{})
		json.Unmarshal(task.ObjectInfo, &workOrderTaskInfo)
		noCorrectiveTypeOrdersTasks += 1

		if _, ok := correctiveWorkOrderTaskStatusCount[util.InterfaceToInt(workOrderTaskInfo["taskStatus"])]; ok {

			correctiveWorkOrderTaskStatusCount[util.InterfaceToInt(workOrderTaskInfo["taskStatus"])] += 1
		}

	}

	preventiveOrder := make(map[string]interface{})
	preventiveOrder["data"] = noPreventiveTypeOrders
	response["totalPreventiveTypeOrders"] = preventiveOrder

	correctiveOrder := make(map[string]interface{})
	correctiveOrder["data"] = noCorrectiveTypeOrders
	response["totalCorrectiveTypeOrders"] = correctiveOrder

	preventiveOrderTask := make(map[string]interface{})
	preventiveOrderTask["data"] = noPreventiveTypeOrderTasks
	response["totalPreventiveTypeOrderTasks"] = preventiveOrderTask

	correctiveOrderTask := make(map[string]interface{})
	correctiveOrderTask["data"] = noCorrectiveTypeOrdersTasks
	response["totalCorrectiveTypeOrderTasks"] = correctiveOrderTask

	ordersCorrectiveCreated := make(map[string]interface{})
	ordersCorrectiveCreated["data"] = correctiveWorkOrderStatusCount[WorkOrderCreated]
	response["correctiveCreatedWorkOrders"] = ordersCorrectiveCreated

	ordersCorrectiveScheduled := make(map[string]interface{})
	ordersCorrectiveScheduled["data"] = correctiveWorkOrderStatusCount[WorkOrderScheduled]
	response["correctiveScheduledWorkOrders"] = ordersCorrectiveScheduled

	ordersCorrectiveInProgress := make(map[string]interface{})
	ordersCorrectiveInProgress["data"] = correctiveWorkOrderStatusCount[WorkOrderInProgress]
	response["correctiveInProgressWorkOrders"] = ordersCorrectiveInProgress

	ordersCorrectiveDone := make(map[string]interface{})
	ordersCorrectiveDone["data"] = correctiveWorkOrderStatusCount[WorkOrderDone]
	response["correctiveCompletedWorkOrders"] = ordersCorrectiveDone

	ordersPreventiveCreated := make(map[string]interface{})
	ordersPreventiveCreated["data"] = preventiveWorkOrderStatusCount[WorkOrderCreated]
	response["preventiveCreatedWorkOrders"] = ordersPreventiveCreated

	ordersPreventiveScheduled := make(map[string]interface{})
	ordersPreventiveScheduled["data"] = preventiveWorkOrderStatusCount[WorkOrderScheduled]
	response["preventiveScheduledWorkOrders"] = ordersPreventiveScheduled

	ordersPreventiveInProgress := make(map[string]interface{})
	ordersPreventiveInProgress["data"] = preventiveWorkOrderStatusCount[WorkOrderInProgress]
	response["preventiveInProgressWorkOrders"] = ordersPreventiveInProgress

	ordersPreventiveDone := make(map[string]interface{})
	ordersPreventiveDone["data"] = preventiveWorkOrderStatusCount[WorkOrderDone]
	response["preventiveCompletedWorkOrders"] = ordersPreventiveDone

	tasksCorrectiveToDo := make(map[string]interface{})
	tasksCorrectiveToDo["data"] = correctiveWorkOrderTaskStatusCount[WorkOrderTaskToDo]
	response["correctiveToDoTasks"] = tasksCorrectiveToDo

	tasksCorrectiveInProgress := make(map[string]interface{})
	tasksCorrectiveInProgress["data"] = correctiveWorkOrderTaskStatusCount[WorkOrderTaskInProgress]
	response["correctiveInProgressTasks"] = tasksCorrectiveInProgress

	tasksCorrectiveDone := make(map[string]interface{})
	tasksCorrectiveDone["data"] = correctiveWorkOrderTaskStatusCount[WorkOrderTaskDone]
	response["correctiveCompletedTasks"] = tasksCorrectiveDone

	tasksPreventiveToDo := make(map[string]interface{})
	tasksPreventiveToDo["data"] = preventiveWorkOrderTaskStatusCount[WorkOrderTaskToDo]
	response["preventiveToDoTasks"] = tasksPreventiveToDo

	tasksPreventiveInProgress := make(map[string]interface{})
	tasksPreventiveInProgress["data"] = preventiveWorkOrderTaskStatusCount[WorkOrderTaskInProgress]
	response["preventiveInProgressTasks"] = tasksPreventiveInProgress

	tasksPreventiveDone := make(map[string]interface{})
	tasksPreventiveDone["data"] = preventiveWorkOrderTaskStatusCount[WorkOrderTaskDone]
	response["preventiveCompletedTasks"] = tasksPreventiveDone

	response["filterOverview"] = filterOverview

	return response
}
