package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
func (v *FactoryService) getFactoryOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var err error
	listOfSites, err := GetObjects(dbConnection, FactorySiteTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	for _, siteInterface := range *listOfSites {
		factorySite := FactorySite{ObjectInfo: siteInterface.ObjectInfo}

		plantCondition := "object_info->>'$.site'=" + strconv.Itoa(siteInterface.Id)
		listOfPlants, err := GetConditionalObjects(dbConnection, FactoryPlantTable, plantCondition)
		var numberOfPlants int
		if err != nil {
			numberOfPlants = 0
		} else {
			numberOfPlants = len(*listOfPlants)
		}
		overviewData := make([]component.OverviewData, 0)
		totalPlants := getKPIData(numberOfPlants, "Total Plants")
		overviewData = append(overviewData, totalPlants)
		overviewResponse = append(overviewResponse, component.OverviewResponse{
			Data:  overviewData,
			Label: factorySite.getFactorySiteInfo().Name + " - " + factorySite.getFactorySiteInfo().Country,
		})
	}

	ctx.JSON(http.StatusOK, overviewResponse)

}
