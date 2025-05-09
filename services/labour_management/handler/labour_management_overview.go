package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func (v *LabourManagementService) summaryResponse(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)
	dbConnection := v.BaseService.ReferenceDatabase

	var err error
	err, listOfShifts := database.GetObjects(dbConnection, const_util.LabourManagementShiftMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	err, shiftStatusMaster := database.GetObjects(dbConnection, const_util.LabourManagementShiftStatusTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	var shiftStatusMap = make(map[int]string)
	for _, shiftStatusInterface := range shiftStatusMaster {
		shiftStatusInfo := database.GetShiftStatusInfo(shiftStatusInterface.ObjectInfo)
		shiftStatusMap[shiftStatusInterface.Id] = shiftStatusInfo.ColorCode
	}

	var numberOfShifts int

	var completedShifts int
	var pendingShifts int
	for index, shiftInterface := range listOfShifts {
		numberOfShifts = index + 1
		shiftInfo := database.GetShiftMasterInfo(shiftInterface.ObjectInfo)
		if shiftInfo.ShiftStatus == const_util.ShiftStatusPending {
			pendingShifts = pendingShifts + 1
		}
		if shiftInfo.ShiftStatus == const_util.ShiftStatusCompleted {
			completedShifts = completedShifts + 1
		}

	}

	labourManagementOverviewData := make([]component.OverviewData, 0)
	totalShiftsCreated := getKPIData(numberOfShifts, "Total Shifts")
	totalCompletedShifts := getKPIDataWithColorCode(completedShifts, "Completed Shifts", shiftStatusMap[const_util.ShiftStatusCompleted])
	totalPendingShifts := getKPIDataWithColorCode(pendingShifts, "Pending Shifts", shiftStatusMap[const_util.ShiftStatusPending])
	labourManagementOverviewData = append(labourManagementOverviewData, totalShiftsCreated)
	labourManagementOverviewData = append(labourManagementOverviewData, totalCompletedShifts)
	labourManagementOverviewData = append(labourManagementOverviewData, totalPendingShifts)

	shiftEventOverview := make([]component.OverviewData, 0)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	arrayOfEventSummary := productionOrderInterface.GetAssemblyEventHistorySummary()
	var eventStatusSummary []EventStatusCount
	json.Unmarshal(arrayOfEventSummary, &eventStatusSummary)

	v.BaseService.Logger.Info("shift event summary ", zap.Any("shift_event_summary", arrayOfEventSummary))
	if len(eventStatusSummary) == 0 {
		scheduled := getKPIData(0, "Total Scheduled")
		running := getKPIData(0, "Total Running")
		completed := getKPIData(0, "Total Completed")
		shiftEventOverview = append(shiftEventOverview, scheduled)
		shiftEventOverview = append(shiftEventOverview, running)
		shiftEventOverview = append(shiftEventOverview, completed)
	} else {
		scheduled := getKPIData(GetTotalEventsByStatus(eventStatusSummary, 4), "Scheduled")
		running := getKPIData(GetTotalEventsByStatus(eventStatusSummary, 5), "Running")
		completed := getKPIData(GetTotalEventsByStatus(eventStatusSummary, 7), "Completed")
		shiftEventOverview = append(shiftEventOverview, scheduled)
		shiftEventOverview = append(shiftEventOverview, running)
		shiftEventOverview = append(shiftEventOverview, completed)
	}

	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  labourManagementOverviewData,
		Label: "Shift History",
	})
	overviewResponse = append(overviewResponse, component.OverviewResponse{
		Data:  shiftEventOverview,
		Label: "Shift Scheduled Events",
	})

	ctx.JSON(http.StatusOK, overviewResponse)

}

// GetTotalEventsByStatus calculates the total count of events for a given event status ID.
func GetTotalEventsByStatus(eventStatusCounts []EventStatusCount, statusID int) int {
	total := 0
	for _, eventStatus := range eventStatusCounts {
		if eventStatus.EventStatusID == statusID {
			total += eventStatus.TotalCount
		}
	}
	return total
}
