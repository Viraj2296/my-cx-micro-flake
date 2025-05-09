package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
func (v *MouldService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var err error
	userId := common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)
	var conditionString = " object_info->>'$.site' =" + strconv.Itoa(basicUserInfo.Site)
	listOfMoulds, err := database.GetConditionalObjects(dbConnection, const_util.MouldMasterTable, conditionString)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfMoulds = len(*listOfMoulds)
	var noOfPrototyping int
	var noOfQualification int
	var noOfCustomerApproval int
	var noOfActive int
	var noOfMaintenance int
	var noOfRepair int

	for _, mouldInterface := range *listOfMoulds {
		mouldMaster := database.MouldMaster{ObjectInfo: mouldInterface.ObjectInfo}
		if mouldMaster.GetMouldMasterInfo().MouldStatus == const_util.MouldStatusPrototyping {
			noOfPrototyping = noOfPrototyping + 1
		}
		if mouldMaster.GetMouldMasterInfo().MouldStatus == const_util.MouldStatusQualification {
			noOfQualification = noOfQualification + 1
		}
		if mouldMaster.GetMouldMasterInfo().MouldStatus == const_util.MouldStatusCustomerApproval {
			noOfCustomerApproval = noOfCustomerApproval + 1
		}
		if mouldMaster.GetMouldMasterInfo().MouldStatus == const_util.MouldStatusActive {
			noOfActive = noOfActive + 1
		}

		if mouldMaster.GetMouldMasterInfo().MouldStatus == const_util.MouldStatusMaintenance {
			noOfMaintenance = noOfMaintenance + 1
		}
		if mouldMaster.GetMouldMasterInfo().MouldStatus == const_util.MouldStatusRepair {
			noOfRepair = noOfRepair + 1
		}

	}

	overviewData := make([]component.OverviewData, 0)
	totalMoulds := getKPIData(numberOfMoulds, "Total Moulds")
	v1 := getKPIData(noOfPrototyping, "Prototyping")
	v2 := getKPIData(noOfQualification, "Qualification")
	v3 := getKPIData(noOfCustomerApproval, "Customer Approval")
	v4 := getKPIData(noOfActive, "Active")
	v5 := getKPIData(noOfMaintenance, "Maintenance")
	v6 := getKPIData(noOfRepair, "Repair")

	overviewData = append(overviewData, totalMoulds)
	v1.BackgroundColor = getBackgroundColour(dbConnection, const_util.MouldStatusPrototyping, const_util.MouldStatusTable)
	overviewData = append(overviewData, v1)
	v2.BackgroundColor = getBackgroundColour(dbConnection, const_util.MouldStatusQualification, const_util.MouldStatusTable)
	overviewData = append(overviewData, v2)
	v3.BackgroundColor = getBackgroundColour(dbConnection, const_util.MouldStatusCustomerApproval, const_util.MouldStatusTable)
	overviewData = append(overviewData, v3)
	v4.BackgroundColor = getBackgroundColour(dbConnection, const_util.MouldStatusActive, const_util.MouldStatusTable)
	overviewData = append(overviewData, v4)
	v5.BackgroundColor = getBackgroundColour(dbConnection, const_util.MouldStatusMaintenance, const_util.MouldStatusTable)
	overviewData = append(overviewData, v5)
	v6.BackgroundColor = getBackgroundColour(dbConnection, const_util.MouldStatusRepair, const_util.MouldStatusTable)
	overviewData = append(overviewData, v6)
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "Moulds Summary",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}

func (v *MouldService) getMouldsStatsSummary(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)
	var conditionString = " object_info->>'$.site' =" + strconv.Itoa(basicUserInfo.Site)
	allMouldData, _ := database.GetConditionalObjects(dbConnection, const_util.MouldMasterTable, conditionString)
	mouldStatus, _ := database.GetObjects(dbConnection, const_util.MouldStatusTable)

	statusCountMap := make(map[int]int)
	mouldStatusMap := make(map[int][]map[string]interface{})
	statusIdMap := make(map[int]database.MouldStatusInfo)

	productionOrderService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	for _, status := range *mouldStatus {
		mouldStatusInfo := database.MouldStatusInfo{}
		json.Unmarshal(status.ObjectInfo, &mouldStatusInfo)

		statusCountMap[status.Id] = 0
		statusIdMap[status.Id] = mouldStatusInfo
		mouldStatusMap[status.Id] = make([]map[string]interface{}, 0)
	}

	for _, mould := range *allMouldData {
		//find counts for each status
		mouldMaster := database.MouldMasterInfo{}
		json.Unmarshal(mould.ObjectInfo, &mouldMaster)
		statusCountMap[mouldMaster.MouldStatus] += 1

		mouldStatusInfo := statusIdMap[mouldMaster.MouldStatus]

		//clasify mould based on status
		mouldStatusObject := make(map[string]interface{})
		mouldStatusObject["toolNo"] = mouldMaster.ToolNo
		mouldStatusObject["status"] = mouldStatusInfo.Status
		mouldStatusObject["mouldImage"] = mouldMaster.MouldImage
		mouldStatusObject["colorCode"] = mouldStatusInfo.ColorCode

		_, lastUpdatedAt := productionOrderService.GetSchedulerUpdatedDate(projectId, mould.Id)
		mouldStatusObject["lastUpdatedAt"] = lastUpdatedAt

		mouldStatusMap[mouldMaster.MouldStatus] = append(mouldStatusMap[mouldMaster.MouldStatus], mouldStatusObject)
	}

	mouldStatusGroupList := make([]interface{}, 0)
	for key, object := range mouldStatusMap {
		groupByView := component.GroupByView{}
		if key == const_util.MouldStatusPrototyping {
			groupByView.GroupByField = "prototyping"
			groupByView.DisplayField = "PROTOTYPING"
		} else if key == const_util.MouldStatusQualification {
			groupByView.GroupByField = "qualification"
			groupByView.DisplayField = "QUALIFICATION"
		} else if key == const_util.MouldStatusCustomerApproval {
			groupByView.GroupByField = "customerApproval"
			groupByView.DisplayField = "CUSTOMER APPROVAL"
		} else if key == const_util.MouldStatusActive {
			groupByView.GroupByField = "active"
			groupByView.DisplayField = "ACTIVE"
		} else if key == const_util.MouldStatusMaintenance {
			groupByView.GroupByField = "maintenance"
			groupByView.DisplayField = "MAINTENANCE"
		} else if key == const_util.MouldStatusRepair {
			groupByView.GroupByField = "repair"
			groupByView.DisplayField = "REPAIR"
		}

		groupByView.Cards = object

		mouldStatusGroupList = append(mouldStatusGroupList, groupByView)
	}

	ctx.JSON(http.StatusOK, mouldStatusGroupList)
}
