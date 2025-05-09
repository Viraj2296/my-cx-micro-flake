package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func getBackgroundColour(dbConnection *gorm.DB, id int, table string) string {
	err, objectInterface := database.Get(dbConnection, table, id)
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
func (v *ITService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var err error
	listOfRequests, err := database.GetObjects(dbConnection, const_util.ITServiceRequestTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfRequests = len(*listOfRequests)
	var noOfRequestUnderUser int
	var noOfRequestUnderHOD int
	var noOfRequestUnderReviewParty int
	var noOfRequestUnderExecutionParty int

	for _, requestInterface := range *listOfRequests {
		serviceRequest := make(map[string]interface{})
		json.Unmarshal(requestInterface.ObjectInfo, &serviceRequest)
		if util.InterfaceToInt(serviceRequest["serviceStatus"]) == 1 {
			noOfRequestUnderUser = noOfRequestUnderUser + 1
		}
		if util.InterfaceToInt(serviceRequest["serviceStatus"]) == 2 {
			noOfRequestUnderHOD = noOfRequestUnderHOD + 1
		}
		if util.InterfaceToInt(serviceRequest["serviceStatus"]) == 3 {
			noOfRequestUnderReviewParty = noOfRequestUnderReviewParty + 1
		}
		if util.InterfaceToInt(serviceRequest["serviceStatus"]) == 4 {
			noOfRequestUnderExecutionParty = noOfRequestUnderExecutionParty + 1
		}

	}

	overviewData := make([]component.OverviewData, 0)
	totalRequests := getKPIData(numberOfRequests, "Total Requests")
	user := getKPIData(noOfRequestUnderUser, "Under User")
	hod := getKPIData(noOfRequestUnderHOD, "Under Head Of Departments")
	reviewParty := getKPIData(noOfRequestUnderReviewParty, "Under Review Party")
	executionParty := getKPIData(noOfRequestUnderExecutionParty, "Under Execution Party")

	overviewData = append(overviewData, totalRequests)
	user.BackgroundColor = getBackgroundColour(dbConnection, 1, const_util.ITServiceRequestStatusTable)
	overviewData = append(overviewData, user)
	hod.BackgroundColor = getBackgroundColour(dbConnection, 2, const_util.ITServiceRequestStatusTable)
	overviewData = append(overviewData, hod)
	reviewParty.BackgroundColor = getBackgroundColour(dbConnection, 3, const_util.ITServiceRequestStatusTable)
	overviewData = append(overviewData, reviewParty)
	executionParty.BackgroundColor = getBackgroundColour(dbConnection, 4, const_util.ITServiceRequestStatusTable)
	overviewData = append(overviewData, executionParty)
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "IT Service Summary",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
