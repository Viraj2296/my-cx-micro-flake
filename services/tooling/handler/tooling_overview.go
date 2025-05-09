package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
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
func reverseSlice(generalObjects *[]component.GeneralObject) *[]component.GeneralObject {
	if generalObjects != nil {
		for i, j := 0, len(*generalObjects)-1; i < j; i, j = i+1, j-1 {
			(*generalObjects)[i], (*generalObjects)[j] = (*generalObjects)[j], (*generalObjects)[i]
		}
	}

	return generalObjects
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
func (v *ToolingService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var err error
	listOfProjects, err := GetObjects(dbConnection, ToolingProjectTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfProjects = len(*listOfProjects)

	overviewData := make([]component.OverviewData, 0)
	totalProjects := getKPIData(numberOfProjects, "Projects")

	overviewData = append(overviewData, totalProjects)

	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "Projects Summary",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
