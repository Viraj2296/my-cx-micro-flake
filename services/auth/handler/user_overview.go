package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func hasPermission(permissions []string, target string) bool {
	for _, permission := range permissions {
		if permission == target {
			return true
		}
	}
	return false
}

func getKPIData(value interface{}, label string, isClickable bool, module, comp, apiQuery, menuId string) component.OverviewData {
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
		IsClickable:     isClickable,
		Module:          module,
		Component:       comp,
		ApiQuery:        apiQuery,
		MenuId:          menuId,
	}

}
func (as *AuthService) summaryResponse(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	dbConnection := as.BaseService.ReferenceDatabase
	userId := common.GetUserId(ctx)
	var isClickable bool

	var err error
	err, listOfUsers := GetObjects(dbConnection, UserTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}
	_, userCurrentRecord := Get(dbConnection, UserTable, userId)
	userCurrentInfo := GetUserInfo(userCurrentRecord.ObjectInfo)

	var numberOfUsers int
	var numberOfActiveUsers int
	var numberOfInActiveUsers int
	var numberOfGroups int
	for _, userInterface := range listOfUsers {
		userInfo := GetUserInfo(userInterface.ObjectInfo)
		if userInfo.Type == "user" {
			numberOfUsers += 1
			if userInfo.Status == "enabled" {
				numberOfActiveUsers += 1
			} else {
				numberOfInActiveUsers += 1
			}
		}

	}

	err, listOfGroups := GetObjects(dbConnection, UserGroupTable)
	if err != nil {
		numberOfGroups = 0
	} else {
		numberOfGroups = len(listOfGroups)
	}

	if userCurrentInfo.Type == "system_admin" {
		isClickable = true
	} else {
		userPermissions := as.getListOfAllowedMenus(userId)
		isClickable = hasPermission(userPermissions, "user_overview")
	}

	userOverviewData := make([]component.OverviewData, 0)
	totalUsers := getKPIData(numberOfUsers, "Total Users", isClickable, ModuleName, "user", "", "user_register")
	totalActiveUsers := getKPIData(numberOfActiveUsers, "Total Active Users", isClickable, ModuleName, "user", "filter=status=enabled", "user_register")
	totalInActiveUsers := getKPIData(numberOfInActiveUsers, "Total In-Active Users", isClickable, ModuleName, "user", "filter=status!=enabled", "user_register")
	userOverviewData = append(userOverviewData, totalUsers)
	userOverviewData = append(userOverviewData, totalActiveUsers)
	userOverviewData = append(userOverviewData, totalInActiveUsers)

	groupOverviewData := make([]component.OverviewData, 0)
	totalGroups := getKPIData(numberOfGroups, "Total Groups", isClickable, ModuleName, "user_group", "", "group_register")
	groupOverviewData = append(groupOverviewData, totalGroups)
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  userOverviewData,
		Label: "Users",
	})
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  groupOverviewData,
		Label: "Groups",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
