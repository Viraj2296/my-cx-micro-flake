package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"github.com/gin-gonic/gin"
	"net/http"
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
func (v *BatchManagementService) summaryResponse(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	dbConnection := v.BaseService.ReferenceDatabase

	var err error
	err, listOfRawMaterials := database.GetObjects(dbConnection, const_util.BatchManagementRawMaterialTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	err, listOfMouldBatch := database.GetObjects(dbConnection, const_util.BatchManagementMouldTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfRawMaterialCreated int
	var numberOfRawMaterialPrinted int
	for _, rawMaterialInterface := range listOfRawMaterials {
		rawMaterialInfo := database.GeRawMaterialBatchInfo(rawMaterialInterface.ObjectInfo)
		if rawMaterialInfo.LabelStatus == const_util.CreatedStatus {
			numberOfRawMaterialCreated = numberOfRawMaterialCreated + 1
		}
		if rawMaterialInfo.LabelStatus == const_util.PrintedStatus {
			numberOfRawMaterialPrinted = numberOfRawMaterialPrinted + 1
		}

	}

	var noOfMouldBatchCreated int
	var noOfMouldBatchPrinted int
	for _, mouldBatchInterface := range listOfMouldBatch {
		mouldBatchInfo := database.GeRawMouldBatchInfo(mouldBatchInterface.ObjectInfo)
		if mouldBatchInfo.LabelStatus == const_util.CreatedStatus {
			noOfMouldBatchCreated = noOfMouldBatchCreated + 1
		}
		if mouldBatchInfo.LabelStatus == const_util.PrintedStatus {
			noOfMouldBatchPrinted = noOfMouldBatchPrinted + 1
		}

	}

	rawMaterialOverviewData := make([]component.OverviewData, 0)
	totalRawMaterialCreated := getKPIData(numberOfRawMaterialCreated, "Total Raw Material Created")
	totalRawMaterialPrinted := getKPIData(numberOfRawMaterialPrinted, "Total Raw Material Printed")
	rawMaterialOverviewData = append(rawMaterialOverviewData, totalRawMaterialCreated)
	rawMaterialOverviewData = append(rawMaterialOverviewData, totalRawMaterialPrinted)

	mouldBatchOverviewData := make([]component.OverviewData, 0)
	totalMouldBatchCreated := getKPIData(noOfMouldBatchCreated, "Total Mould Created")
	totalMouldBatchPrinted := getKPIData(noOfMouldBatchPrinted, "Total Mould Printed")
	mouldBatchOverviewData = append(mouldBatchOverviewData, totalMouldBatchCreated)
	mouldBatchOverviewData = append(mouldBatchOverviewData, totalMouldBatchPrinted)

	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  rawMaterialOverviewData,
		Label: "Raw Material Batch",
	})
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  mouldBatchOverviewData,
		Label: "Mould Batch",
	})

	ctx.JSON(http.StatusOK, overviewResponse)

}
