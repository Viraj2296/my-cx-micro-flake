package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/scheduler/handler/const_util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// getMouldingOrderSchedulerView ShowAccount godoc
// @Summary Get the record message
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Param   RecordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/scheduler_view [get]
func (v *SchedulerService) getMouldingOrderSchedulerView(ctx *gin.Context) {

	//TODO
	// events  - ScheduledOrderEvent(under production order) == machine_time_event
	// resources -> resources should come from several places, machine (list of machines)
	// assignment -> which event is assigned to which resource ?
	// events (ScheduledOrderEvent)  ->machineId , we can create mapping of machineId and resourcs
	// if eventsAllocated is true, only send the events and resources which have events > 0

	// hard code the generation
	projectId := ctx.Param("projectId")
	requestedMonth := ctx.Query("requestedMonth")
	mode := ctx.Query("mode")

	eventsAllocated := ctx.Query("eventsAllocated")
	var filterEventAllocated bool
	if eventsAllocated == "true" {
		filterEventAllocated = true
	} else {
		filterEventAllocated = false
	}

	singaporeLocation, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(singaporeLocation)

	requestedMonthInt := 1
	if requestedMonth != "" {
		requestedMonthInt, _ = strconv.Atoi(requestedMonth)
	}

	startOfMonth := now.AddDate(0, -requestedMonthInt, 0)

	productionOrderService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)
	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

	var listOfScheduledOrderEvents *[]component.GeneralObject
	_, listOfScheduledOrderEvents = productionOrderService.GetAllSchedulerEventForScheduler(projectId)

	if mode == "read_only" {
		listOfScheduledOrderEvents = getReadOnlySchedulerEvents(listOfScheduledOrderEvents)
	}

	_, listOfWorkOrderEvents := maintenanceService.GetMaintenanceEventsForScheduler(projectId)
	//_, listOfWorkOrderTaskEvents := maintenanceService.GetMaintenanceOrderTaskEventsForScheduler(projectId)
	_, listOfMouldTestEvents := mouldService.GetMouldTestEventsForScheduler(projectId)

	//*listOfEvents = append(*listOfEvents, *listOfWorkOrderEvents...)
	*listOfScheduledOrderEvents = append(*listOfScheduledOrderEvents, *listOfWorkOrderEvents...)
	*listOfScheduledOrderEvents = append(*listOfScheduledOrderEvents, *listOfMouldTestEvents...)

	_, listOfResources := machineService.GetListOfMachines(projectId)

	_, listOfSubCategory := machineService.GetListOfMachineSubCategory(projectId)

	var arrayOfEvents []interface{}
	var arrayOfResource []interface{}
	var arrayOfAssignment []interface{}
	var eventResourceMap = make(map[int]int)
	assigmentId := 1
	for _, event := range *listOfScheduledOrderEvents {
		var objects = make(map[string]interface{})
		var assignmentObject = make(map[string]interface{})
		json.Unmarshal(event.ObjectInfo, &objects)
		if objects["objectStatus"].(string) != "Archived" {
			objects["id"] = assigmentId
			objects["eventId"] = event.Id
			startDate := objects["startDate"].(string)
			endDate := objects["endDate"].(string)
			objectStartDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
			objectEndDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)
			objects["startDate"] = objectStartDate
			objects["endDate"] = objectEndDate

			startDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["startDate"]))
			endDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["endDate"]))

			if startDateTime.Unix() > startOfMonth.Unix() || endDateTime.Unix() > startOfMonth.Unix() {
				if objects["eventType"] == "production_schedule" {
					machineId := util.InterfaceToInt(objects["machineId"])
					productionOrderId := util.InterfaceToInt(objects["eventSourceId"])

					_, productionGeneralObject := productionOrderService.GetMachineProductionOrderInfo(projectId, productionOrderId, machineId)

					var productionObjectInfo = make(map[string]interface{})
					json.Unmarshal(productionGeneralObject.ObjectInfo, &productionObjectInfo)
					partId := util.InterfaceToInt(productionObjectInfo["partNumber"])

					_, partMasterObject := mouldService.GetPartInfo(projectId, partId)
					var partMasterInfo = make(map[string]interface{})
					json.Unmarshal(partMasterObject.ObjectInfo, &partMasterInfo)

					partName := util.InterfaceToString(partMasterInfo["partNumber"])
					objects["partNo"] = partName
				}
				objects["permissions"] = common.GetPermissions(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]))

				keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease"}
				objects = common.RemoveKeys(objects, keysToRemove)

				arrayOfEvents = append(arrayOfEvents, objects)
				machineId := util.InterfaceToInt(objects["machineId"])

				//Create assignment object
				assignmentObject["id"] = assigmentId
				assignmentObject["event"] = assigmentId
				assignmentObject["resource"] = machineId
				arrayOfAssignment = append(arrayOfAssignment, assignmentObject)
				assigmentId += 1
				if _, ok := eventResourceMap[machineId]; !ok {
					//do something here
					eventResourceMap[machineId] = event.Id
				}
			}

			//if mouldUpTime, ok := objects["mouldUp"]; ok {
			//	if mouldUpTime != "" {
			//
			//		mouldUpObject := make(map[string]interface{})
			//
			//		mouldUpStr := util.InterfaceToString(mouldUpTime)
			//		mouldUpObject["startDate"] = objectStartDate
			//
			//		startDateTime, _ := time.Parse("2006-01-02T15:04:05", objectStartDate)
			//
			//		timeList := strings.Split(mouldUpStr, ":")
			//
			//		hour, _ := strconv.Atoi(timeList[0])
			//		min, _ := strconv.Atoi(timeList[1])
			//
			//		endTime := startDateTime.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min))
			//
			//		mouldUpObject["endDate"] = endTime.Format("2006-01-02T15:04:05")
			//		mouldUpObject["eventType"] = "Mould Up"
			//		mouldUpObject["name"] = objects["name"]
			//		mouldUpObject["draggable"] = false
			//		mouldUpObject["canComplete"] = false
			//		mouldUpObject["isUpdate"] = false
			//		arrayOfEvents = append(arrayOfEvents, mouldUpObject)
			//
			//		objects["startDate"] = endTime.Format("2006-01-02T15:04:05")
			//	}
			//
			//}

			//if mouldDown, ok := objects["mouldDown"]; ok {
			//	if mouldDown != "" {
			//		mouldDownObject := make(map[string]interface{})
			//
			//		mouldUpStr := util.InterfaceToString(mouldDown)
			//		mouldDownObject["endDate"] = objectEndDate
			//		endDateTime, _ := time.Parse("2006-01-02T15:04:05", objectEndDate)
			//
			//		timeList := strings.Split(mouldUpStr, ":")
			//		hour, _ := strconv.Atoi(timeList[0])
			//		min, _ := strconv.Atoi(timeList[1])
			//		startTime := endDateTime.Add(-time.Hour*time.Duration(hour) - time.Minute*time.Duration(min))
			//		mouldDownObject["startDate"] = startTime.Format("2006-01-02T15:04:05")
			//		mouldDownObject["eventType"] = "Mould Down"
			//		mouldDownObject["name"] = objects["name"]
			//		mouldDownObject["draggable"] = false
			//		mouldDownObject["canComplete"] = false
			//		mouldDownObject["isUpdate"] = false
			//		arrayOfEvents = append(arrayOfEvents, mouldDownObject)
			//
			//		objects["endDate"] = startTime.Format("2006-01-02T15:04:05")
			//	}
			//
			//}

		}

	}
	for _, resource := range listOfResources {
		var objects = make(map[string]interface{})
		var responseObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &objects)
		responseObject["id"] = resource.Id
		responseObject["name"] = objects["newMachineId"]
		responseObject["type"] = GetSubCategoryNameById(listOfSubCategory, util.InterfaceToInt(objects["subCategory"]))
		if filterEventAllocated {
			if _, ok := eventResourceMap[resource.Id]; ok {
				//do something here
				arrayOfResource = append(arrayOfResource, responseObject)
			}
		} else {
			arrayOfResource = append(arrayOfResource, responseObject)
		}

	}

	// for index, assignment := range *listOfAssignment {
	// 	var objects = make(map[string]interface{})
	// 	json.Unmarshal(assignment.ObjectInfo, &objects)
	// 	objects["id"] = assignment.Id
	// 	arrayOfAssignment = append(arrayOfAssignment, objects)
	// }

	schedulerResponse := component.SchedulerResponse{}
	schedulerResponse.Events.Rows = arrayOfEvents
	schedulerResponse.Resources.Rows = arrayOfResource
	schedulerResponse.Assignments.Rows = arrayOfAssignment

	ctx.JSON(http.StatusOK, schedulerResponse)
}

// getMouldingOrderSchedulerView ShowAccount godoc
// @Summary Get the record message
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Param   RecordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/scheduler_view [get]
func (v *SchedulerService) getSchedulerOverView(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	var response = make(map[string]interface{}, 100)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	// conditionString := " JSON_EXTRACT(object_info, \"$.eventSourceId\") = '" + strconv.Itoa(productionOrderId) + "'"
	listOfTotalObjects, _ := GetObjects(dbConnection, const_util.ScheduledOrderTable)

	if listOfTotalObjects != nil {
		response["totalNoScheduledOrder"] = map[string]int{"data": len(*listOfTotalObjects)}
	} else {
		response["totalNoScheduledOrder"] = map[string]int{"data": 0}
	}

	productionService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	recordIdForRunning := productionService.GetOrderStatusIdFromPreferenceLevel(projectId, const_util.ScheduleStatusPreferenceFive)
	recordIdForCompleted := productionService.GetOrderStatusIdFromPreferenceLevel(projectId, const_util.ScheduleStatusPreferenceSeven)

	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, mouldTestEvent := mouldService.GetMouldTestEventsForScheduler(projectId)

	conditionRunningStatusString := " JSON_EXTRACT(object_info, \"$.eventStatus\") = " + strconv.Itoa(recordIdForRunning)
	listOfRunningObjects, _ := GetConditionalObjects(dbConnection, const_util.ScheduledOrderTable, conditionRunningStatusString)
	if listOfRunningObjects != nil {
		response["totalNoRunningScheduledOrder"] = map[string]int{"data": len(*listOfRunningObjects)}
	} else {
		response["totalNoRunningScheduledOrder"] = map[string]int{"data": 0}
	}

	conditionCompletedStatusString := " JSON_EXTRACT(object_info, \"$.eventStatus\") = " + strconv.Itoa(recordIdForCompleted)
	listOfCompletedObjects, _ := GetConditionalObjects(dbConnection, const_util.ScheduledOrderTable, conditionCompletedStatusString)

	if listOfCompletedObjects != nil {
		response["totalNoCompletedScheduledOrder"] = map[string]int{"data": len(*listOfCompletedObjects)}
	} else {
		response["totalNoCompletedScheduledOrder"] = map[string]int{"data": 0}
	}

	if mouldTestEvent != nil {
		response["totalNoMouldTestRequest"] = map[string]int{"data": len(*mouldTestEvent)}
	} else {
		response["totalNoMouldTestRequest"] = map[string]int{"data": 0}
	}

	ctx.JSON(http.StatusOK, response)
}

func (v *SchedulerService) getMaintenanceView(ctx *gin.Context) {

	//TODO
	// events  - ScheduledOrderEvent(under production order) == machine_time_event
	// resources -> resources should come from several places, machine (list of machines)
	// assignment -> which event is assigned to which resource ?
	// events (ScheduledOrderEvent)  ->machineId , we can create mapping of machineId and resourcs
	// if eventsAllocated is true, only send the events and resources which have events > 0

	// hard code the generation
	projectId := ctx.Param("projectId")
	requestedMonth := ctx.Query("requestedMonth")

	eventsAllocated := ctx.Query("eventsAllocated")
	var filterEventAllocated bool
	if eventsAllocated == "true" {
		filterEventAllocated = true
	} else {
		filterEventAllocated = false
	}

	//currentTimeString := util.GetZoneCurrentTime("Asia/Singapore")
	//currentTime := util.ConvertStringToDateTime(currentTimeString)
	//monthQuery := int(currentTime.DateTime.Month())
	//
	//requestedMonthInt := 1
	//if requestedMonth != "" {
	//	requestedMonthInt, _ = strconv.Atoi(requestedMonth)
	//}
	singaporeLocation, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(singaporeLocation)

	requestedMonthInt := 1
	if requestedMonth != "" {
		requestedMonthInt, _ = strconv.Atoi(requestedMonth)
	}

	startOfMonth := now.AddDate(0, -requestedMonthInt, 0)

	maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)

	_, listOfEvents := maintenanceService.GetMaintenanceOrderForScheduler(projectId)
	_, listOfCorrectiveEvents := maintenanceService.GetMaintenanceCorrectiveOrderForScheduler(projectId)

	*listOfEvents = append(*listOfEvents, *listOfCorrectiveEvents...)
	_, listOfResources := machineService.GetListOfMachines(projectId)

	_, listOfSubCategory := machineService.GetListOfMachineSubCategory(projectId)

	var arrayOfEvents []interface{}
	var arrayOfResource []interface{}
	var arrayOfAssignment []interface{}
	var eventResourceMap = make(map[int]int)
	assigmentId := 1
	for _, event := range *listOfEvents {
		var objects = make(map[string]interface{})
		var assignmentObject = make(map[string]interface{})
		json.Unmarshal(event.ObjectInfo, &objects)
		if objects["objectStatus"].(string) != "Archived" {
			objects["id"] = assigmentId
			objects["eventId"] = event.Id
			startDate := objects["startDate"].(string)
			endDate := objects["endDate"].(string)
			objects["startDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
			objects["endDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)

			startDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["startDate"]))
			endDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["endDate"]))

			if startDateTime.Unix() > startOfMonth.Unix() || endDateTime.Unix() > startOfMonth.Unix() {

				objects["permissions"] = common.GetPermissionsForMaintenance(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]), util.InterfaceToBool(objects["canUnRelease"]))

				keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease", "canUnRelease"}
				objects = common.RemoveKeys(objects, keysToRemove)

				arrayOfEvents = append(arrayOfEvents, objects)
				machineId := util.InterfaceToInt(objects["machineId"])

				//Create assignment object
				assignmentObject["id"] = assigmentId
				assignmentObject["event"] = assigmentId
				assignmentObject["resource"] = machineId
				arrayOfAssignment = append(arrayOfAssignment, assignmentObject)
				assigmentId += 1
				if _, ok := eventResourceMap[machineId]; !ok {
					//do something here
					eventResourceMap[machineId] = event.Id
				}
			}
		}

	}
	for _, resource := range listOfResources {
		var objects = make(map[string]interface{})
		var responseObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &objects)
		responseObject["id"] = resource.Id
		responseObject["name"] = objects["newMachineId"]
		responseObject["type"] = GetSubCategoryNameById(listOfSubCategory, util.InterfaceToInt(objects["subCategory"]))
		if filterEventAllocated {
			if _, ok := eventResourceMap[resource.Id]; ok {
				//do something here
				arrayOfResource = append(arrayOfResource, responseObject)
			}
		} else {
			arrayOfResource = append(arrayOfResource, responseObject)
		}

	}

	schedulerResponse := component.SchedulerResponse{}
	schedulerResponse.Events.Rows = arrayOfEvents
	schedulerResponse.Resources.Rows = arrayOfResource
	schedulerResponse.Assignments.Rows = arrayOfAssignment

	ctx.JSON(http.StatusOK, schedulerResponse)
}

func (v *SchedulerService) getMouldMaintenanceView(ctx *gin.Context) {

	// hard code the generation
	projectId := ctx.Param("projectId")
	requestedMonth := ctx.Query("requestedMonth")

	eventsAllocated := ctx.Query("eventsAllocated")
	var filterEventAllocated bool
	if eventsAllocated == "true" {
		filterEventAllocated = true
	} else {
		filterEventAllocated = false
	}

	singaporeLocation, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(singaporeLocation)

	requestedMonthInt := 1
	if requestedMonth != "" {
		requestedMonthInt, _ = strconv.Atoi(requestedMonth)
	}

	startOfMonth := now.AddDate(0, -requestedMonthInt, 0)

	maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)
	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

	_, listOfEvents := maintenanceService.GetMouldPreventiveMaintenanceOrderForScheduler(projectId)
	_, listOfCorrectiveEvents := maintenanceService.GetMouldCorrectiveMaintenanceOrderForScheduler(projectId)

	*listOfEvents = append(*listOfEvents, *listOfCorrectiveEvents...)
	_, listOfResources := mouldService.GetListOfMoulds(projectId)

	_, listOfSubCategory := mouldService.GetListOfMouldCategory(projectId)

	var arrayOfEvents []interface{}
	var arrayOfResource []interface{}
	var arrayOfAssignment []interface{}
	var eventResourceMap = make(map[int]int)
	assigmentId := 1

	for _, event := range *listOfEvents {

		var objects = make(map[string]interface{})
		var assignmentObject = make(map[string]interface{})
		json.Unmarshal(event.ObjectInfo, &objects)
		if objects["objectStatus"].(string) != "Archived" {

			objects["id"] = assigmentId
			objects["eventId"] = event.Id
			startDate := objects["startDate"].(string)
			endDate := objects["endDate"].(string)
			objects["startDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
			objects["endDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)

			startDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["startDate"]))
			endDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["endDate"]))

			if startDateTime.Unix() > startOfMonth.Unix() || endDateTime.Unix() > startOfMonth.Unix() {

				objects["permissions"] = common.GetPermissionsForMaintenance(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]), util.InterfaceToBool(objects["canUnRelease"]))

				keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease", "canUnRelease"}
				objects = common.RemoveKeys(objects, keysToRemove)

				arrayOfEvents = append(arrayOfEvents, objects)
				mouldId := util.InterfaceToInt(objects["mouldId"])

				//Create assignment object
				assignmentObject["id"] = assigmentId
				assignmentObject["event"] = assigmentId
				assignmentObject["resource"] = mouldId
				arrayOfAssignment = append(arrayOfAssignment, assignmentObject)
				assigmentId += 1
				if _, ok := eventResourceMap[mouldId]; !ok {
					//do something here
					eventResourceMap[mouldId] = event.Id
				}
			}
		}

	}
	for _, resource := range listOfResources {
		var objects = make(map[string]interface{})
		var responseObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &objects)
		responseObject["id"] = resource.Id
		responseObject["name"] = objects["toolNo"]
		responseObject["type"] = GetSubCategoryNameById(listOfSubCategory, util.InterfaceToInt(objects["category"]))
		if filterEventAllocated {
			if _, ok := eventResourceMap[resource.Id]; ok {
				//do something here
				arrayOfResource = append(arrayOfResource, responseObject)
			}
		} else {
			arrayOfResource = append(arrayOfResource, responseObject)
		}

	}

	schedulerResponse := component.SchedulerResponse{}
	schedulerResponse.Events.Rows = arrayOfEvents
	schedulerResponse.Resources.Rows = arrayOfResource
	schedulerResponse.Assignments.Rows = arrayOfAssignment

	ctx.JSON(http.StatusOK, schedulerResponse)
}
func (v *SchedulerService) getAssemblyView(ctx *gin.Context) {

	//TODO
	// events  - ScheduledOrderEvent(under production order) == machine_time_event
	// resources -> resources should come from several places, machine (list of machines)
	// assignment -> which event is assigned to which resource ?
	// events (ScheduledOrderEvent)  ->machineId , we can create mapping of machineId and resourcs
	// if eventsAllocated is true, only send the events and resources which have events > 0

	// hard code the generation
	projectId := ctx.Param("projectId")
	requestedMonth := ctx.Query("requestedMonth")
	mode := ctx.Query("mode")
	eventsAllocated := ctx.Query("eventsAllocated")
	var filterEventAllocated bool
	if eventsAllocated == "true" {
		filterEventAllocated = true
	} else {
		filterEventAllocated = false
	}

	//currentTimeString := util.GetZoneCurrentTime("Asia/Singapore")
	//currentTime := util.ConvertStringToDateTime(currentTimeString)
	//monthQuery := int(currentTime.DateTime.Month())
	singaporeLocation, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(singaporeLocation)

	requestedMonthInt := 1
	if requestedMonth != "" {
		requestedMonthInt, _ = strconv.Atoi(requestedMonth)
	}

	startOfMonth := now.AddDate(0, -requestedMonthInt, 0)

	productionOrderService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	manufacturingService := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
	maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)

	_, listOfEvents := productionOrderService.GetAllAssemblyEventForScheduler(projectId)
	if mode == "read_only" {
		listOfEvents = getReadOnlySchedulerEvents(listOfEvents)
	}
	_, listOfWorkOrderEvents := maintenanceService.GetMaintenanceEventsForAssemblyScheduler(projectId)
	_, listOfResources := machineService.GetListOfAssemblyMachines(projectId)

	*listOfEvents = append(*listOfEvents, *listOfWorkOrderEvents...)

	var arrayOfEvents []interface{}
	var arrayOfResource []interface{}
	var arrayOfAssignment []interface{}
	var eventResourceMap = make(map[int]int)
	assigmentId := 1
	for _, event := range *listOfEvents {
		var objects = make(map[string]interface{})
		var assignmentObject = make(map[string]interface{})
		json.Unmarshal(event.ObjectInfo, &objects)
		if objects["objectStatus"].(string) != "Archived" {
			objects["id"] = assigmentId
			objects["eventId"] = event.Id
			startDate := objects["startDate"].(string)
			endDate := objects["endDate"].(string)
			objects["startDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
			objects["endDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)

			startDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["startDate"]))
			endDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["endDate"]))

			if startDateTime.Unix() > startOfMonth.Unix() || endDateTime.Unix() > startOfMonth.Unix() {

				machineId := util.InterfaceToInt(objects["machineId"])
				productionOrderId := util.InterfaceToInt(objects["eventSourceId"])

				_, productionGeneralObject := productionOrderService.GetAssemblyProductionOrderInfo(projectId, productionOrderId, machineId)

				var productionObjectInfo = make(map[string]interface{})
				json.Unmarshal(productionGeneralObject.ObjectInfo, &productionObjectInfo)
				partId := util.InterfaceToInt(productionObjectInfo["partNumber"])

				_, partMasterObject := manufacturingService.GetAssemblyPartInfo(projectId, partId)
				var partMasterInfo = make(map[string]interface{})
				json.Unmarshal(partMasterObject.ObjectInfo, &partMasterInfo)

				partName := util.InterfaceToString(partMasterInfo["partNumber"])
				objects["partNo"] = partName

				objects["name"] = partName + "_" + util.InterfaceToString(objects["name"])
				objects["permissions"] = common.GetPermissionsForAssembly(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]))

				keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease"}
				objects = common.RemoveKeys(objects, keysToRemove)

				fmt.Println("Passed")
				arrayOfEvents = append(arrayOfEvents, objects)
				resourceId := util.InterfaceToInt(objects["machineId"])

				//Create assignment object
				assignmentObject["id"] = assigmentId
				assignmentObject["event"] = assigmentId
				assignmentObject["resource"] = resourceId
				arrayOfAssignment = append(arrayOfAssignment, assignmentObject)
				assigmentId += 1
				if _, ok := eventResourceMap[resourceId]; !ok {
					//do something here
					eventResourceMap[resourceId] = event.Id
				}

			}

		}

	}
	for _, resource := range listOfResources {
		var objects = make(map[string]interface{})
		var responseObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &objects)
		responseObject["id"] = resource.Id
		responseObject["name"] = objects["newMachineId"]
		responseObject["type"] = "Default"
		if filterEventAllocated {
			if _, ok := eventResourceMap[resource.Id]; ok {
				//do something here
				arrayOfResource = append(arrayOfResource, responseObject)
			}
		} else {
			arrayOfResource = append(arrayOfResource, responseObject)
		}

	}

	schedulerResponse := component.SchedulerResponse{}
	schedulerResponse.Events.Rows = arrayOfEvents
	schedulerResponse.Resources.Rows = arrayOfResource
	schedulerResponse.Assignments.Rows = arrayOfAssignment

	ctx.JSON(http.StatusOK, schedulerResponse)
}

// getToolingView ShowAccount godoc
// @Summary Get the record message
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Param   RecordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/scheduler_view [get]
func (v *SchedulerService) getToolingView(ctx *gin.Context) {

	//TODO
	// events  - ScheduledOrderEvent(under production order) == machine_time_event
	// resources -> resources should come from several places, machine (list of machines)
	// assignment -> which event is assigned to which resource ?
	// events (ScheduledOrderEvent)  ->machineId , we can create mapping of machineId and resourcs
	// if eventsAllocated is true, only send the events and resources which have events > 0

	// hard code the generation
	projectId := ctx.Param("projectId")
	requestedMonth := ctx.Query("requestedMonth")
	mode := ctx.Query("mode")
	eventsAllocated := ctx.Query("eventsAllocated")
	var filterEventAllocated bool
	if eventsAllocated == "true" {
		filterEventAllocated = true
	} else {
		filterEventAllocated = false
	}

	//currentTimeString := util.GetZoneCurrentTime("Asia/Singapore")
	//currentTime := util.ConvertStringToDateTime(currentTimeString)
	//monthQuery := int(currentTime.DateTime.Month())
	singaporeLocation, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(singaporeLocation)

	requestedMonthInt := 1
	if requestedMonth != "" {
		requestedMonthInt, _ = strconv.Atoi(requestedMonth)
	}

	startOfMonth := now.AddDate(0, -requestedMonthInt, 0)

	productionOrderService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)

	_, listOfEvents := productionOrderService.GetAllToolingEventForScheduler(projectId)
	if mode == "read_only" {
		listOfEvents = getReadOnlySchedulerEvents(listOfEvents)
	}
	_, listOfWorkOrderEvents := maintenanceService.GetMaintenanceEventsForToolingScheduler(projectId)
	_, listOfResources := machineService.GetListOfToolingMachines(projectId)

	*listOfEvents = append(*listOfEvents, *listOfWorkOrderEvents...)

	var arrayOfEvents []interface{}
	var arrayOfResource []interface{}
	var arrayOfAssignment []interface{}
	var eventResourceMap = make(map[int]int)
	assigmentId := 1
	for _, event := range *listOfEvents {
		var objects = make(map[string]interface{})
		var assignmentObject = make(map[string]interface{})
		json.Unmarshal(event.ObjectInfo, &objects)
		if objects["objectStatus"].(string) != "Archived" {
			objects["id"] = assigmentId
			objects["eventId"] = event.Id
			startDate := objects["startDate"].(string)
			endDate := objects["endDate"].(string)
			objects["startDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
			objects["endDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)

			startDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["startDate"]))
			endDateTime, _ := time.Parse(const_util.Layout, util.InterfaceToString(objects["endDate"]))

			if startDateTime.Unix() > startOfMonth.Unix() || endDateTime.Unix() > startOfMonth.Unix() {
				fmt.Println("Passed")

				partId := util.InterfaceToInt(objects["partId"])
				eventSourceId := util.InterfaceToInt(objects["eventSourceId"])

				_, partGeneralObject := productionOrderService.GetToolingPartById(projectId, partId)

				var partObjectInfo = make(map[string]interface{})
				json.Unmarshal(partGeneralObject.ObjectInfo, &partObjectInfo)
				partName := util.InterfaceToString(partObjectInfo["name"])
				objects["partNo"] = partName

				_, bomGeneralObject := productionOrderService.GetBomInfo(projectId, eventSourceId)
				var bomObjectInfo = make(map[string]interface{})
				json.Unmarshal(bomGeneralObject.ObjectInfo, &bomObjectInfo)
				bomName := util.InterfaceToString(bomObjectInfo["name"])
				objects["bomName"] = bomName

				objects["permissions"] = common.GetPermissions(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]))

				keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease"}
				objects = common.RemoveKeys(objects, keysToRemove)

				arrayOfEvents = append(arrayOfEvents, objects)
				resourceId := util.InterfaceToInt(objects["machineId"])

				//Create assignment object
				assignmentObject["id"] = assigmentId
				assignmentObject["event"] = assigmentId
				assignmentObject["resource"] = resourceId
				arrayOfAssignment = append(arrayOfAssignment, assignmentObject)
				assigmentId += 1
				if _, ok := eventResourceMap[resourceId]; !ok {
					//do something here
					eventResourceMap[resourceId] = event.Id
				}

			}

		}

	}
	for _, resource := range listOfResources {
		var objects = make(map[string]interface{})
		var responseObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &objects)
		responseObject["id"] = resource.Id
		responseObject["name"] = objects["newMachineId"]
		responseObject["type"] = "Default"
		if filterEventAllocated {
			if _, ok := eventResourceMap[resource.Id]; ok {
				//do something here
				arrayOfResource = append(arrayOfResource, responseObject)
			}
		} else {
			arrayOfResource = append(arrayOfResource, responseObject)
		}

	}

	schedulerResponse := component.SchedulerResponse{}
	schedulerResponse.Events.Rows = arrayOfEvents
	schedulerResponse.Resources.Rows = arrayOfResource
	schedulerResponse.Assignments.Rows = arrayOfAssignment

	ctx.JSON(http.StatusOK, schedulerResponse)
}

func GetSubCategoryNameById(listOfSubCategory []component.GeneralObject, id int) string {
	subCategory := ""
	for _, subCategory := range listOfSubCategory {
		if subCategory.Id == id {
			var objects = make(map[string]interface{})
			json.Unmarshal(subCategory.ObjectInfo, &objects)
			return util.InterfaceToString(objects["name"])
		}
	}

	return subCategory
}

func getReadOnlySchedulerEvents(listOfScheduledOrderEvents *[]component.GeneralObject) *[]component.GeneralObject {
	var modifiedScheduledOrderEvent []component.GeneralObject
	for _, scheduledOrderEvent := range *listOfScheduledOrderEvents {
		var orderEventMap = make(map[string]interface{})
		json.Unmarshal(scheduledOrderEvent.ObjectInfo, &orderEventMap)
		if _, ok := orderEventMap["canRelease"]; ok {
			orderEventMap["canRelease"] = false
		}
		if _, ok := orderEventMap["canHold"]; ok {
			orderEventMap["canHold"] = false
		}
		if _, ok := orderEventMap["canForceStop"]; ok {
			orderEventMap["canForceStop"] = false
		}
		if _, ok := orderEventMap["draggable"]; ok {
			orderEventMap["draggable"] = false
		}
		if _, ok := orderEventMap["isAbortEnabled"]; ok {
			orderEventMap["isAbortEnabled"] = false
		}
		if _, ok := orderEventMap["canComplete"]; ok {
			orderEventMap["canComplete"] = false
		}
		if _, ok := orderEventMap["isUpdate"]; ok {
			orderEventMap["isUpdate"] = false
		}
		serialisedData, _ := json.Marshal(orderEventMap)
		modifiedScheduledOrderEvent = append(modifiedScheduledOrderEvent, component.GeneralObject{Id: scheduledOrderEvent.Id, ObjectInfo: serialisedData})

	}
	return &modifiedScheduledOrderEvent
}
