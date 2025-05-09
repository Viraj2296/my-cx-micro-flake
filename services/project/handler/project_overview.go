package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (v *ProjectService) getSummaryResponse(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	userId := common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	isHead, departmentUsers := authService.GetDepartmentUsers(userId)
	fmt.Println("department users :", departmentUsers)

	var listOfObjects *[]component.GeneralObject
	var err error
	if isHead {
		// this user is department head
		// this is not department head
		listOfObjects, err = GetObjects(dbConnection, ProjectTable)

		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
	} else {
		// this is not department head
		userBasedQuery := " object_info ->>'$.createdBy' = " + strconv.Itoa(userId) + " "
		listOfObjects, err = GetConditionalObjects(dbConnection, ProjectTable, userBasedQuery)

		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

	}
	overviewData := make([]component.OverviewData, 0)
	var assignmentDataArray []map[string]interface{}
	var assignmentData = make(map[string]interface{}, 0)
	assignmentData["v1"] = len(*listOfObjects)
	assignmentDataArray = append(assignmentDataArray, assignmentData)
	overviewData = append(overviewData, component.OverviewData{
		Value:           assignmentDataArray,
		IsVisible:       true,
		Label:           "Total Assignments",
		Icon:            "bx:task",
		BackgroundColor: "#49C4ED",
	})
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "Assignments",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
