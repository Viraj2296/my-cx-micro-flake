package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getSchedulerView ShowAccount godoc
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
// @Router /project/{projectId}/maintenance/component/{componentName}/scheduler_view [get]
func (v *MaintenanceService) getSchedulerView(ctx *gin.Context) {

	// hard code the generation
	projectId := ctx.Param("projectId")
	//eventsAllocated := ctx.Query("eventsAllocated")
	//var filterEventAllocated bool
	//if eventsAllocated == "true" {
	//	filterEventAllocated = true
	//} else {
	//	filterEventAllocated = false
	//}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfWorkOrders, _ := GetObjects(dbConnection, MaintenanceWorkOrderTable)

	machineInterface := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	_, listOfMachinesObjects := machineInterface.GetListOfMachines(projectId)
	mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, listOfMoulds := mouldInterface.GetListOfMoulds(projectId)
	var arrayOfEvents []interface{}
	var arrayOfResource []interface{}
	var arrayOfAssignment []interface{}
	//var eventResourceMap = make(map[int]int)

	var assignmentId int
	for _, workOrderObject := range *listOfWorkOrders {
		var workOrderEvents = make(map[string]interface{})
		json.Unmarshal(workOrderObject.ObjectInfo, &workOrderEvents)
		if workOrderEvents["objectStatus"].(string) != "Archived" {
			var eventObject = make(map[string]interface{})
			eventObject["id"] = workOrderObject.Id
			startDate := workOrderEvents["workOrderScheduledStartDate"].(string)
			endDate := workOrderEvents["workOrderScheduledEndDate"].(string)
			eventObject["startDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
			eventObject["endDate"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)
			eventObject["eventType"] = "work_order"
			eventObject["iconCls"] = "fa fa-cogs"
			eventObject["eventColor"] = "#CE8196"
			arrayOfEvents = append(arrayOfEvents, eventObject)
			machineId := util.InterfaceToInt(workOrderEvents["assetId"])
			fmt.Println("machineId: ", machineId)
			//if _, ok := eventResourceMap[machineId]; !ok {
			//	//do something here
			//	eventResourceMap[machineId] = workOrderObject.Id
			//}
			var assignmentObject = make(map[string]interface{})
			assignmentObject["id"] = assignmentId
			assignmentObject["event"] = workOrderObject.Id
			assignmentObject["resource"] = machineId
			arrayOfAssignment = append(arrayOfAssignment, assignmentObject)
			assignmentId += 1
		}

	}

	for _, resource := range listOfMachinesObjects {
		var machineObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &machineObject)
		var resourceObject = make(map[string]interface{})
		resourceObject["id"] = resource.Id
		resourceObject["name"] = machineObject["newMachineId"]
		resourceObject["type"] = "machines"
		arrayOfResource = append(arrayOfResource, resourceObject)
		//if filterEventAllocated {
		//	if _, ok := eventResourceMap[resource.Id]; ok {
		//		//do something here
		//		arrayOfResource = append(arrayOfResource, objects)
		//	}
		//} else {
		//	arrayOfResource = append(arrayOfResource, objects)
		//}

	}
	for _, resource := range listOfMoulds {
		var machineObject = make(map[string]interface{})
		json.Unmarshal(resource.ObjectInfo, &machineObject)
		var resourceObject = make(map[string]interface{})
		resourceObject["id"] = resource.Id
		resourceObject["name"] = machineObject["toolNo"]
		resourceObject["type"] = "moulds"
		arrayOfResource = append(arrayOfResource, resourceObject)
		//if filterEventAllocated {
		//	if _, ok := eventResourceMap[resource.Id]; ok {
		//		//do something here
		//		arrayOfResource = append(arrayOfResource, objects)
		//	}
		//} else {
		//	arrayOfResource = append(arrayOfResource, objects)
		//}

	}

	schedulerResponse := component.SchedulerResponse{}
	fmt.Println("arrayOfEvents: ", arrayOfEvents)
	schedulerResponse.Events.Rows = arrayOfEvents
	schedulerResponse.Resources.Rows = arrayOfResource
	fmt.Println("arrayOfResource: ", arrayOfResource)
	schedulerResponse.Assignments.Rows = arrayOfAssignment
	fmt.Println("arrayOfAssignment: ", arrayOfAssignment)

	ctx.JSON(http.StatusOK, schedulerResponse)
}
