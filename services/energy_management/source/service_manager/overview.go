package service_manager

import (
	"cx-micro-flake/pkg/common/component"
	"github.com/gin-gonic/gin"
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
func (v *Service) getOverview(ctx *gin.Context) {
	var overviewResponse = make([]component.OverviewResponse, 0)

	ctx.JSON(http.StatusOK, overviewResponse)

}
