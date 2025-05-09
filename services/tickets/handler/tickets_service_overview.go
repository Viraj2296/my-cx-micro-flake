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
func (ts *TicketsService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := ts.BaseService.ServiceDatabases[projectId]

	var err error
	listOfRequests, err := GetObjects(dbConnection, TicketsServiceRequestTable)
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
		serviceRequest := TicketsServiceRequest{ObjectInfo: requestInterface.ObjectInfo}
		if serviceRequest.getServiceRequestInfo().ServiceStatus == 1 {
			noOfRequestUnderUser = noOfRequestUnderUser + 1
		}
		if serviceRequest.getServiceRequestInfo().ServiceStatus == 2 {
			noOfRequestUnderHOD = noOfRequestUnderHOD + 1
		}
		if serviceRequest.getServiceRequestInfo().ServiceStatus == 3 {
			noOfRequestUnderReviewParty = noOfRequestUnderReviewParty + 1
		}
		if serviceRequest.getServiceRequestInfo().ServiceStatus == 4 {
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
	user.BackgroundColor = getBackgroundColour(dbConnection, 1, TicketsServiceRequestStatusTable)
	overviewData = append(overviewData, user)
	hod.BackgroundColor = getBackgroundColour(dbConnection, 2, TicketsServiceRequestStatusTable)
	overviewData = append(overviewData, hod)
	reviewParty.BackgroundColor = getBackgroundColour(dbConnection, 3, TicketsServiceRequestStatusTable)
	overviewData = append(overviewData, reviewParty)
	executionParty.BackgroundColor = getBackgroundColour(dbConnection, 4, TicketsServiceRequestStatusTable)
	overviewData = append(overviewData, executionParty)
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "Case Management Summary",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
