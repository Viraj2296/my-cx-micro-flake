package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/scheduler/handler/const_util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getBackgroundColour(dbConnection *gorm.DB, id int, table string) string {
	err, objectInterface := Get(dbConnection, table, id)
	if err != nil {
		return "#49C4ED"
	} else {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
		return util.InterfaceToString(objectFields["colorCode"])
	}

}

func getKPIData(value interface{}, label string) component.OverviewData {
	var arrayResponse []map[string]interface{}
	var numberOfUsersData = make(map[string]interface{}, 0)
	numberOfUsersData["v1"] = value
	arrayResponse = append(arrayResponse, numberOfUsersData)

	return component.OverviewData{
		Value:           arrayResponse,
		IsVisible:       true,
		Label:           label,
		Icon:            "bx:task",
		BackgroundColor: "#49C4ED",
	}

}

func (v *SchedulerService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var err error
	listOfObjects, err := GetObjects(dbConnection, const_util.ScheduledOrderTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfScheduledOrders = len(*listOfObjects)

	productionService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	recordIdForRunning := productionService.GetOrderStatusIdFromPreferenceLevel(projectId, const_util.ScheduleStatusPreferenceFive)
	recordIdForCompleted := productionService.GetOrderStatusIdFromPreferenceLevel(projectId, const_util.ScheduleStatusPreferenceSeven)
	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, mouldTestEvent := mouldService.GetMouldTestEventsForScheduler(projectId)
	conditionRunningStatusString := " object_info->>'$.eventStatus' = " + strconv.Itoa(recordIdForRunning)
	listOfRunningObjects, _ := GetConditionalObjects(dbConnection, const_util.ScheduledOrderTable, conditionRunningStatusString)

	var numberOfRunningOrders = len(*listOfRunningObjects)

	conditionCompletedStatusString := " object_info->>'$.eventStatus' = " + strconv.Itoa(recordIdForCompleted)
	listOfCompletedObjects, _ := GetConditionalObjects(dbConnection, const_util.ScheduledOrderTable, conditionCompletedStatusString)

	var numberOfCompletedOrders = len(*listOfCompletedObjects)
	var numberOfMouldTestEvents = len(*mouldTestEvent)

	overviewData := make([]component.OverviewData, 0)
	totalScheduledOrders := getKPIData(numberOfScheduledOrders, "Total Scheduled Orders")
	runningOrders := getKPIData(numberOfRunningOrders, "Total Running Orders")
	completedOrders := getKPIData(numberOfCompletedOrders, "Total Completed Orders")
	mouldTestSchedules := getKPIData(numberOfMouldTestEvents, "Mould Test Schedules")

	overviewData = append(overviewData, totalScheduledOrders)
	overviewData = append(overviewData, runningOrders)
	overviewData = append(overviewData, completedOrders)
	overviewData = append(overviewData, mouldTestSchedules)
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "Scheduling Summary",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
