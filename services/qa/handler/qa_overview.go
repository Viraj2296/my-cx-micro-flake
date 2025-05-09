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
func (v *QAService) summaryResponse(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	dbConnection := v.BaseService.ReferenceDatabase

	var err error
	listOfQABatchRecords, err := GetObjects(dbConnection, QABatchTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var pendingApproval int
	var approved int
	var rejected int
	for _, qaBatchInterface := range *listOfQABatchRecords {
		var qaBatchFields = make(map[string]interface{})
		json.Unmarshal(qaBatchInterface.ObjectInfo, &qaBatchFields)
		qualityStatus := util.InterfaceToInt(qaBatchFields["qualityStatus"])

		if qualityStatus == 1 {
			pendingApproval = pendingApproval + 1
		}
		if qualityStatus == 2 {
			approved = approved + 1
		}
		if qualityStatus == 3 {
			rejected = rejected + 1
		}

	}

	overviewData := make([]component.OverviewData, 0)
	pendingApprovalKPI := getKPIData(pendingApproval, "Pending Approval")
	approvedKPI := getKPIData(approved, "Approved")
	rejectedKPI := getKPIData(rejected, "Rejected")
	overviewData = append(overviewData, pendingApprovalKPI)
	overviewData = append(overviewData, approvedKPI)
	overviewData = append(overviewData, rejectedKPI)

	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  overviewData,
		Label: "QA Status",
	})
	ctx.JSON(http.StatusOK, overviewResponse)

}
