package http

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/spare_part_inventory/source/consts"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func getKPIDataWithColorCode(value interface{}, label string, colorCode string) component.OverviewData {
	var arrayResponse []map[string]interface{}
	var numberOfUsersData = make(map[string]interface{}, 0)
	numberOfUsersData["v1"] = value
	arrayResponse = append(arrayResponse, numberOfUsersData)

	return component.OverviewData{
		Value:           arrayResponse,
		IsVisible:       true,
		Label:           label,
		Icon:            "bx:task",
		BackgroundColor: colorCode,
	}
}

func (v *Service) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	var sparePartInventoryMasterTable = v.ComponentManager.GetTargetTable(consts.SparePartInventoryMasterStatusComponent)
	var err error
	err, listOfMaster := v.Repository.GetResources(sparePartInventoryMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}
	var sparePartInventoryTransactionTable = v.ComponentManager.GetTargetTable(consts.SparePartInventorTransactionComponent)
	err, listOfTransaction := v.Repository.GetResources(sparePartInventoryTransactionTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfTransactionIn int
	var numberOfTransactionOut int
	for _, transaction := range listOfTransaction {
		transactionInfo := models.GetSparePartInventoryTransactionInfo(transaction.ObjectInfo)
		if consts.TransactionStatusIn == transactionInfo.Transaction {
			numberOfTransactionIn += 1
		}
		if consts.TransactionStatusOut == transactionInfo.Transaction {
			numberOfTransactionOut += 1
		}

	}

	sparePartInventoryOverviewData := make([]component.OverviewData, 0)

	totalShiftsCreated := getKPIData(len(listOfMaster), "Total Spare Parts")
	numberOfTransactionReceived := getKPIData(numberOfTransactionIn, "Total Transaction Received")
	numberOfTransactionDelivered := getKPIData(numberOfTransactionOut, "Total Transaction Out")
	sparePartInventoryOverviewData = append(sparePartInventoryOverviewData, totalShiftsCreated)
	sparePartInventoryOverviewData = append(sparePartInventoryOverviewData, numberOfTransactionReceived)
	sparePartInventoryOverviewData = append(sparePartInventoryOverviewData, numberOfTransactionDelivered)

	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  sparePartInventoryOverviewData,
		Label: "Shift History",
	})

	ctx.JSON(http.StatusOK, overviewResponse)

}
