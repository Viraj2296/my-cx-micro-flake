package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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
func (v *ProductionOrderService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var err error
	listOfProductionOrders, err := GetObjects(dbConnection, ProductionOrderMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfOrders = len(*listOfProductionOrders)
	var noOfPartiallyScheduled int
	var noOfScheduled int
	var noOfConfirmed int
	var noOfRunning int
	var noOfProductionStopped int
	var noOfCompleted int

	for _, orderInterface := range *listOfProductionOrders {
		productionOrder := ProductionOrderMaster{ObjectInfo: orderInterface.ObjectInfo}
		if productionOrder.getProductionOrderInfo().OrderStatus == 2 {
			noOfPartiallyScheduled = noOfPartiallyScheduled + 1
		}
		if productionOrder.getProductionOrderInfo().OrderStatus == 3 {
			noOfScheduled = noOfScheduled + 1
		}
		if productionOrder.getProductionOrderInfo().OrderStatus == 4 {
			noOfConfirmed = noOfConfirmed + 1
		}
		if productionOrder.getProductionOrderInfo().OrderStatus == 5 {
			noOfRunning = noOfRunning + 1
		}
		if productionOrder.getProductionOrderInfo().OrderStatus == 6 {
			noOfProductionStopped = noOfProductionStopped + 1
		}
		if productionOrder.getProductionOrderInfo().OrderStatus == 7 {
			noOfCompleted = noOfCompleted + 1
		}

	}

	overviewData := make([]component.OverviewData, 0)
	totalOrders := getKPIData(numberOfOrders, "Total Orders")
	partiallyScheduled := getKPIData(noOfPartiallyScheduled, "Partially Scheduled")
	scheduled := getKPIData(noOfScheduled, "Scheduled")
	confirmed := getKPIData(noOfConfirmed, "Confirmed")
	running := getKPIData(noOfRunning, "Running")
	productionStopped := getKPIData(noOfProductionStopped, "Production Stopped")
	completed := getKPIData(noOfCompleted, "Completed")
	overviewData = append(overviewData, totalOrders)
	partiallyScheduled.BackgroundColor = getBackgroundColour(dbConnection, 2, ProductionOrderStatusTable)
	overviewData = append(overviewData, partiallyScheduled)
	scheduled.BackgroundColor = getBackgroundColour(dbConnection, 3, ProductionOrderStatusTable)
	overviewData = append(overviewData, scheduled)
	confirmed.BackgroundColor = getBackgroundColour(dbConnection, 4, ProductionOrderStatusTable)
	overviewData = append(overviewData, confirmed)
	running.BackgroundColor = getBackgroundColour(dbConnection, 5, ProductionOrderStatusTable)
	overviewData = append(overviewData, running)
	productionStopped.BackgroundColor = getBackgroundColour(dbConnection, 6, ProductionOrderStatusTable)
	overviewData = append(overviewData, productionStopped)
	completed.BackgroundColor = getBackgroundColour(dbConnection, 7, ProductionOrderStatusTable)
	overviewData = append(overviewData, completed)
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "Orders",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
