package handler

import (
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (as *AnalyticsService) getDatasourceMasterCardView(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)

	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var err error

	listOfDatasourceMasterObjects, err := GetObjects(dbConnection, AnalyticsDatasourcesMasterTable)

	if err != nil {
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}

	var groupByCardResponse []GroupByCardView
	for _, datasourceMasterObject := range *listOfDatasourceMasterObjects {
		datasourceMaster := AnalyticsDatasourcesMaster{ObjectInfo: datasourceMasterObject.ObjectInfo}

		category := util.InterfaceToString(datasourceMaster.getDatasourceMasterInfo().Category)
		var datasourceMasterFields map[string]interface{}
		json.Unmarshal(datasourceMasterObject.ObjectInfo, &datasourceMasterFields)
		var isElementFound bool
		isElementFound = false
		for index, mm := range groupByCardResponse {
			if mm.GroupByField == category {

				groupByCardResponse[index].Cards = append(groupByCardResponse[index].Cards, datasourceMasterFields)
				isElementFound = true
			}
		}
		if !isElementFound {
			xl := GroupByCardView{}
			xl.GroupByField = category
			xl.Cards = append(xl.Cards, datasourceMasterFields)
			groupByCardResponse = append(groupByCardResponse, xl)
		}
	}

	ctx.JSON(http.StatusOK, groupByCardResponse)
}
