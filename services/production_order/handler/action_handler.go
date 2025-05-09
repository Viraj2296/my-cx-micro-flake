package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type Assignment struct {
	Event    int `json:"int"`
	Id       int `json:"id"`
	Resource int `json:"resource"`
}

type OrderScheduledEventUpdateRequest struct {
	EndDate         string      `json:"endDate"`
	StartDate       string      `json:"startDate"`
	Assignment      *Assignment `json:"assignment"`
	PlannedManpower int         `json:"plannedManpower"`
	PriorityLevel   int         `json:"priorityLevel"`
}

func (v *ProductionOrderService) updateSchedulerEvent(ctx *gin.Context) {

	updateScheduleEventRequest := OrderScheduledEventUpdateRequest{}
	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, scheduledEventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	scheduledOrderInfo := make(map[string]interface{})
	//scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	//scheduledOrderInfo := scheduledOrderEvent.getScheduledOrderEventInfo()
	json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledOrderInfo)

	//Can't modify archived machine time line
	if util.InterfaceToString(scheduledOrderInfo["objectStatus"]) == "Archived" {
		response.SendAlreadyArchivedError(ctx)
		return
	}

	preferenceOrderStatusId3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
	preferenceOrderStatusId8 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceEight)
	preferenceOrderStatusId7 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSeven)

	if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) > preferenceOrderStatusId3 {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSchedulePosition), InvalidScheduleStatus, "Schedule order passed the update stage")
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateScheduleEventRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	orderInfoStartDate := scheduledOrderInfo["startDate"]
	orderInfoEndDate := scheduledOrderInfo["endDate"]

	// we should be receiving singapore time zone time, so we should change and update that in to UTC time
	if updateScheduleEventRequest.StartDate != "" {
		updateScheduleEventRequest.StartDate = util.ConvertSingaporeTimeToUTC(updateScheduleEventRequest.StartDate)
		scheduledOrderInfo["startDate"] = updateScheduleEventRequest.StartDate
	}
	if updateScheduleEventRequest.EndDate != "" {
		updateScheduleEventRequest.EndDate = util.ConvertSingaporeTimeToUTC(updateScheduleEventRequest.EndDate)
		scheduledOrderInfo["endDate"] = updateScheduleEventRequest.EndDate
	}

	currentTime := time.Now().Unix()

	startTime, _ := time.Parse(TimeLayout, updateScheduleEventRequest.StartDate)
	endTime, _ := time.Parse(TimeLayout, updateScheduleEventRequest.EndDate)
	v.BaseService.Logger.Info("date validation", zap.Any("start_date", updateScheduleEventRequest.StartDate), zap.Any("end_date", updateScheduleEventRequest.EndDate),
		zap.Any("start_time", startTime), zap.Any("end_time", endTime), zap.Any("current_time", currentTime))

	reqStart, _ := time.Parse(time.RFC3339, updateScheduleEventRequest.StartDate)
	reqEnd, _ := time.Parse(time.RFC3339, updateScheduleEventRequest.EndDate)
	dbStart, _ := time.Parse(time.RFC3339, util.InterfaceToString(orderInfoStartDate))
	dbEnd, _ := time.Parse(time.RFC3339, util.InterfaceToString(orderInfoEndDate))

	if !reqStart.Equal(dbStart) {

		if currentTime > dbStart.Unix() {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSourceError), InvalidScheduleStatus, "You are trying to modify the event which is already passed, This modification is rejected")
			return
		}
	} else if !reqEnd.Equal(dbEnd) {

		if currentTime > dbEnd.Unix() {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSourceError), InvalidScheduleStatus, "You are trying to modify the event which is already passed, This modification is rejected")
			return
		}
	}

	// update is very critical, we need to verify whether there is overlapping scheduled, get the
	var conditionStartDate, conditionEndDate string
	if updateScheduleEventRequest.StartDate != "" {
		conditionStartDate = updateScheduleEventRequest.StartDate
	} else {
		conditionStartDate = util.InterfaceToString(scheduledOrderInfo["startDate"])
	}
	if updateScheduleEventRequest.EndDate != "" {
		conditionEndDate = updateScheduleEventRequest.EndDate
	} else {
		conditionEndDate = util.InterfaceToString(scheduledOrderInfo["endDate"])
	}
	machineId := util.InterfaceToInt(scheduledOrderInfo["machineId"])

	dateConditionString := "id != " + strconv.Itoa(recordId) + " AND object_info->>'$.machineId' = " + strconv.Itoa(machineId) + " AND (( CONVERT_TZ(\"" + conditionStartDate + "\",'+00:00','+00:00') > CONVERT_TZ(object_info->>'$.startDate','+00:00','+00:00')   AND CONVERT_TZ(\"" + conditionStartDate + "\",'+00:00','+00:00') <  CONVERT_TZ(object_info->>'$.endDate','+00:00','+00:00'))  OR ( CONVERT_TZ(\"" + conditionEndDate + "\",'+00:00','+00:00') > CONVERT_TZ(object_info->>'$.startDate','+00:00','+00:00')   AND CONVERT_TZ(\"" + conditionEndDate + "\",'+00:00','+00:00') <  CONVERT_TZ(object_info->>'$.endDate','+00:00','+00:00')))"
	listOfScheduledEvents, _ := GetConditionalObjects(dbConnection, targetTable, dateConditionString)

	if len(*listOfScheduledEvents) != 0 {
		var numberOfUnconfirmedEvents int
		// check if all are confirmed, then error, otherwise it is okay.
		for _, scheduledEvent := range *listOfScheduledEvents {
			tmpScheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEvent.ObjectInfo}

			if tmpScheduledOrderEvent.getScheduledOrderEventInfo().EventStatus == preferenceOrderStatusId3 || tmpScheduledOrderEvent.getScheduledOrderEventInfo().EventStatus == preferenceOrderStatusId8 || tmpScheduledOrderEvent.getScheduledOrderEventInfo().EventStatus == preferenceOrderStatusId7 {

				numberOfUnconfirmedEvents = numberOfUnconfirmedEvents + 1
			}
		}
		if numberOfUnconfirmedEvents != len(*listOfScheduledEvents) {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSchedulePosition), InvalidScheduleStatus, "You are trying to confirm or changing the schedule with overlapping date range, make sure you are finding exact empty slot to assign the event, single resource can not have or confirmed with overlapping schedule period")
			return
		}

	}

	if updateScheduleEventRequest.Assignment != nil {
		if updateScheduleEventRequest.Assignment.Resource != util.InterfaceToInt(scheduledOrderInfo["machineId"]) {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSourceError), InvalidScheduleStatus, "Changing resource which already assigned is not permitted, consider moving along with the resources")
			return
		}
	}
	userId := common.GetUserId(ctx)

	updatingData := make(map[string]interface{})
	scheduledOrderInfo["lastUpdatedBy"] = userId
	scheduledOrderInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

	//update plannedManpower, and schedule priority level
	if plannedManpower, ok := scheduledOrderInfo["plannedManpower"]; ok {
		var finalPlannedManpower int

		if util.InterfaceToInt(plannedManpower) == 0 {
			machineInterface := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			finalPlannedManpower = machineInterface.GetAssemblyMachineDefaultManpower(machineId)
		} else if updateScheduleEventRequest.PlannedManpower != 0 {
			finalPlannedManpower = updateScheduleEventRequest.PlannedManpower
		} else {
			finalPlannedManpower = util.InterfaceToInt(plannedManpower)
		}

		scheduledOrderInfo["plannedManpower"] = finalPlannedManpower
	}

	if _, ok := scheduledOrderInfo["priorityLevel"]; ok {
		scheduledOrderInfo["priorityLevel"] = updateScheduleEventRequest.PriorityLevel
	}
	updatingData["object_info"], _ = json.Marshal(scheduledOrderInfo)

	err = Update(dbConnection, targetTable, recordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
		return
	}
	v.createSystemNotification(projectId, util.InterfaceToString(scheduledOrderInfo["name"]), util.InterfaceToString(scheduledOrderInfo["name"])+" is updated", recordId)
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: util.InterfaceToString(scheduledOrderInfo["name"]) + " is successfully updated",
		Code:    0,
	})
}
func (v *ProductionOrderService) releaseOrder(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, scheduledEventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	scheduledOrderInfo := make(map[string]interface{})
	json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledOrderInfo)
	//Can't modify archived machine time line
	if util.InterfaceToString(scheduledOrderInfo["objectStatus"]) == "Archived" {
		response.SendAlreadyArchivedError(ctx)
		return
	}

	preferenceOrderStatusId3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
	preferenceOrderStatusId4 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFour)
	_, preferenceLevel4ColorCode := v.getColorCode(dbConnection, preferenceOrderStatusId4)
	v.BaseService.Logger.Info("order status info", zap.Any("event_status", scheduledOrderInfo["eventStatus"]), zap.Any("scheduler_order_info", scheduledOrderInfo),
		zap.Any("preference_level", preferenceOrderStatusId3))

	if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) != preferenceOrderStatusId3 {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Status"), InvalidScheduleStatus, "You are releasing order where status is not properly aligned")
		return
	}
	if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == preferenceOrderStatusId4 {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Status"), InvalidScheduleStatus, "You have already made it released, doing modification is not permitted")
		return
	}

	scheduledOrderInfo["eventStatus"] = preferenceOrderStatusId4
	scheduledOrderInfo["eventColor"] = preferenceLevel4ColorCode

	if targetTable == AssemblyScheduledOrderEventTable {
		v.ViewManager.CreateLabourManagementShiftLinesHistory(recordId, preferenceOrderStatusId4)
	}
	userId := common.GetUserId(ctx)

	updatingData := make(map[string]interface{})
	scheduledOrderInfo["lastUpdatedBy"] = userId
	scheduledOrderInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"], _ = json.Marshal(scheduledOrderInfo)

	err = Update(dbConnection, targetTable, recordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
		return
	}
}

func (v *ProductionOrderService) holdOrder(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	err, scheduledEventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	scheduledOrderInfo := make(map[string]interface{})
	json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledOrderInfo)
	//Can't modify archived machine time line
	if util.InterfaceToString(scheduledOrderInfo["objectStatus"]) == "Archived" {
		response.SendAlreadyArchivedError(ctx)
		return
	}

	preferenceOrderStatusId3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
	preferenceOrderStatusId4 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFour)
	_, preferenceLevel3ColorCode := v.getColorCode(dbConnection, preferenceOrderStatusId3)
	scheduledOrderInfo["eventStatus"] = preferenceOrderStatusId4
	if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) != preferenceOrderStatusId4 {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Status"), InvalidScheduleStatus, "You are releasing order where status is not properly aligned")
		return
	}
	if scheduledOrderInfo["eventStatus"] == preferenceOrderStatusId3 {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Status"), InvalidScheduleStatus, "You have already made it released, doing modification is not permitted")
		return
	}

	scheduledOrderInfo["eventStatus"] = preferenceOrderStatusId3
	scheduledOrderInfo["eventColor"] = preferenceLevel3ColorCode
	if targetTable == AssemblyScheduledOrderEventTable {
		v.ViewManager.CreateLabourManagementShiftLinesHistory(recordId, preferenceOrderStatusId3)
	}
	userId := common.GetUserId(ctx)

	updatingData := make(map[string]interface{})
	scheduledOrderInfo["lastUpdatedBy"] = userId
	scheduledOrderInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"], _ = json.Marshal(scheduledOrderInfo)

	err = Update(dbConnection, targetTable, recordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
		return
	}
}

// completeOrder ShowAccount godoc
// @Summary schedule orders
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId}/action/schedule [get]
func (v *ProductionOrderService) completeOrder(ctx *gin.Context) {
	userId := common.GetUserId(ctx)
	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	err, generalScheduledEventObject := Get(dbConnection, targetTable, recordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}
	scheduledOrderEventInfo := make(map[string]interface{})
	json.Unmarshal(generalScheduledEventObject.ObjectInfo, &scheduledOrderEventInfo)

	preferenceSevenOrderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSeven)
	preferenceFiveOrderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFive)
	if util.InterfaceToInt(scheduledOrderEventInfo["eventStatus"]) == preferenceSevenOrderStatusId {
		//send error saying, it is already scheduled
		response.DispatchDetailedError(ctx, AlreadyScheduled,
			&response.DetailedError{
				Header:      "Invalid Operation",
				Description: "The current order status is already in a same status that your trying to change. Redundant actions are not performed",
			})
		return
	}
	if util.InterfaceToInt(scheduledOrderEventInfo["eventStatus"]) == preferenceFiveOrderStatusId {
		response.DispatchDetailedError(ctx, AlreadyScheduled,
			&response.DetailedError{
				Header:      "Invalid Operation",
				Description: "Currently, the order is running, please stop this order to complete",
			})
		return
	}
	scheduledOrderEventInfo["eventStatus"] = preferenceSevenOrderStatusId
	scheduledOrderEventInfo["isAbortEnabled"] = false
	scheduledOrderEventInfo["canComplete"] = false
	scheduledOrderEventInfo["isUpdate"] = false
	scheduledOrderEventInfo["draggable"] = false
	if targetTable == AssemblyScheduledOrderEventTable {
		v.ViewManager.CreateLabourManagementShiftLinesHistory(recordId, preferenceSevenOrderStatusId)
	}
	productionOrderTable := ProductionOrderMasterTable

	if targetTable == AssemblyScheduledOrderEventTable {
		productionOrderTable = AssemblyProductionOrderTable
	} else if targetTable == ToolingScheduledOrderEventTable {
		productionOrderTable = ToolingOrderMasterTable
	}

	if targetTable == AssemblyScheduledOrderEventTable || targetTable == ProductionOrderMasterTable {
		addCarryForwardRemainingQty(dbConnection, scheduledOrderEventInfo, userId, productionOrderTable)
	}

	updatingData := make(map[string]interface{})
	scheduledOrderEventInfo["lastUpdatedBy"] = userId
	scheduledOrderEventInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"], _ = json.Marshal(scheduledOrderEventInfo)

	err = Update(dbConnection, targetTable, recordId, updatingData)

	// now generate the label for completed or, this mean created the label to print. Not sure we have needed this flag.

	// if this event has the mould batch resource ID created, then use that
	//if mouldBatchResourceId, ok := scheduledOrderEventInfo["mouldBatchResourceId"]; ok {
	//	v.BaseService.Logger.Info("updating the mould batch, and generating the label", zap.Any("resource_id", mouldBatchResourceId))
	//	// if it is available, take that update the end time
	//	batchManagementInterface := common.GetService("batch_management_module").ServiceInterface.(common.BatchManagementInterface)
	//	var resourceId = util.InterfaceToInt(mouldBatchResourceId)
	//	batchManagementInterface.GenerateMouldBatchLabel(projectId, resourceId)
	//}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Production order has successfully completed",
	})

}

func addCarryForwardRemainingQty(dbConnection *gorm.DB, scheduledOrderEventInfo map[string]interface{}, userId int, targetTable string) {
	err, productionOrder := Get(dbConnection, targetTable, util.InterfaceToInt(scheduledOrderEventInfo["eventSourceId"]))
	if err != nil {
		return
	}

	var productionOrderInfo ProductionOrderInfo
	json.Unmarshal(productionOrder.ObjectInfo, &productionOrderInfo)

	productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty + (util.InterfaceToInt(scheduledOrderEventInfo["scheduledQty"]) - util.InterfaceToInt(scheduledOrderEventInfo["completedQty"]))

	_ = Update(dbConnection, targetTable, productionOrder.Id, productionOrderInfo.DatabaseSerialize(userId))
}

// completeOrder ShowAccount godoc
// @Summary schedule orders
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId}/action/schedule [get]
func (v *ProductionOrderService) abortOrder(ctx *gin.Context) {
	userId := common.GetUserId(ctx)
	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	err, generalScheduledEventObject := Get(dbConnection, targetTable, recordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}

	scheduledOrderEventInfo := make(map[string]interface{})
	json.Unmarshal(generalScheduledEventObject.ObjectInfo, &scheduledOrderEventInfo)
	preferenceSevenOrderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSeven)

	if util.InterfaceToInt(scheduledOrderEventInfo["eventStatus"]) == preferenceSevenOrderStatusId {
		response.DispatchDetailedError(ctx, AlreadyScheduled,
			&response.DetailedError{
				Header:      "Invalid Operation",
				Description: "The current order status is already completed.",
			})
		return
	}
	scheduledOrderEventInfo["eventStatus"] = ScheduleStatusPreferenceEight
	scheduledOrderEventInfo["isAbortEnabled"] = false

	productionOrderTable := ProductionOrderMasterTable

	if targetTable == AssemblyScheduledOrderEventTable {
		v.ViewManager.CreateLabourManagementShiftLinesHistory(recordId, ScheduleStatusPreferenceEight)
		productionOrderTable = AssemblyProductionOrderTable
	} else if targetTable == ToolingScheduledOrderEventTable {
		productionOrderTable = ToolingOrderMasterTable
	}

	if targetTable == AssemblyScheduledOrderEventTable || targetTable == ScheduledOrderEventTable {
		addCarryForwardRemainingQty(dbConnection, scheduledOrderEventInfo, userId, productionOrderTable)
	}

	updatingData := make(map[string]interface{})
	scheduledOrderEventInfo["lastUpdatedBy"] = userId
	scheduledOrderEventInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"], _ = json.Marshal(scheduledOrderEventInfo)

	err = Update(dbConnection, targetTable, recordId, updatingData)

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Scheduled order has successfully aborted",
	})
}

func (v *ProductionOrderService) forceStopOrder(ctx *gin.Context) {
	userId := common.GetUserId(ctx)
	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	err, generalScheduledEventObject := Get(dbConnection, targetTable, recordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}

	scheduledOrderEventInfo := make(map[string]interface{})
	json.Unmarshal(generalScheduledEventObject.ObjectInfo, &scheduledOrderEventInfo)
	preferenceSevenOrderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSeven)

	if util.InterfaceToBool(scheduledOrderEventInfo["canForceStop"]) == false {
		response.DispatchDetailedError(ctx, AlreadyScheduled,
			&response.DetailedError{
				Header:      "Invalid Operation",
				Description: "Can't stop the order.",
			})
		return
	}
	scheduledOrderEventInfo["eventStatus"] = preferenceSevenOrderStatusId
	scheduledOrderEventInfo["isAbortEnabled"] = false
	scheduledOrderEventInfo["canForceStop"] = false

	machineId := util.InterfaceToInt(scheduledOrderEventInfo["machineId"])
	eventId := generalScheduledEventObject.Id

	machineInterface := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	if componentName == ScheduledOrderEventComponent {
		machineInterface.AddMouldingForceStop(projectId, machineId, eventId)
	} else if componentName == AssemblyScheduledOrderEventTable {
		machineInterface.AddAssemblyForceStop(projectId, machineId, eventId)
		v.ViewManager.CreateLabourManagementShiftLinesHistory(recordId, preferenceSevenOrderStatusId)
	} else if componentName == ToolingScheduledOrderEventTable {
		machineInterface.AddToolingForceStop(projectId, machineId, eventId)
	} else {
		response.DispatchDetailedError(ctx, AlreadyScheduled,
			&response.DetailedError{
				Header:      "Invalid Request",
				Description: "Can't stop the order.",
			})
		return
	}

	updatingData := make(map[string]interface{})
	scheduledOrderEventInfo["lastUpdatedBy"] = userId
	scheduledOrderEventInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"], _ = json.Marshal(scheduledOrderEventInfo)

	err = Update(dbConnection, targetTable, recordId, updatingData)

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Scheduled order has successfully aborted",
	})
}
