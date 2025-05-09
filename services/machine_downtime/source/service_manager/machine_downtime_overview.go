package service_manager

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/models"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/orm"
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

type EventStatusCount struct {
	EventID       int `json:"eventId"`
	EventStatusID int `json:"EventStatusId"`
	TotalCount    int `json:"totalCount"`
}

// summaryResponse this should contains shift information, and plus total planned lines, and active lines
func (v *MachineDowntimeService) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)

	var err error
	err, listOfDowntimeMaster := orm.GetObjects(v.BaseService.ServiceDatabase, consts.MachineDownTimeMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var numberOfDowns int
	var numberOfPendingTasks int
	var numberOfCompletedTasks int
	for index, shiftInterface := range listOfDowntimeMaster {
		numberOfDowns = index + 1
		machineDowntimeInfo := models.GetMachineDowntimeInfo(shiftInterface.ObjectInfo)
		if machineDowntimeInfo.CheckOutDate == "" {
			numberOfPendingTasks = numberOfPendingTasks + 1
		}
		if machineDowntimeInfo.CheckOutDate != "" {
			numberOfCompletedTasks = numberOfCompletedTasks + 1
		}

	}

	labourManagementOverviewData := make([]component.OverviewData, 0)
	totalFailuresReported := getKPIData(numberOfDowns, "Total Failures Reported")
	totalCompletedTasks := getKPIDataWithColorCode(numberOfCompletedTasks, "Completed Tasks", "#9fcf99")
	totalPendingTasks := getKPIDataWithColorCode(numberOfPendingTasks, "Pending Tasks", "#c8d497")
	labourManagementOverviewData = append(labourManagementOverviewData, totalFailuresReported)
	labourManagementOverviewData = append(labourManagementOverviewData, totalCompletedTasks)
	labourManagementOverviewData = append(labourManagementOverviewData, totalPendingTasks)

	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  labourManagementOverviewData,
		Label: "Downtime History",
	})

	ctx.JSON(http.StatusOK, overviewResponse)

}
