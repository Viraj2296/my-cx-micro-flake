package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (as *AnalyticsService) RefreshData(ctx *gin.Context) {

	var updateRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	_, datasourceInterface := Get(dbConnection, AnalyticsDataSourceTable, recordId)

	dataSourceInfo := DatasourceInfo{}
	json.Unmarshal(datasourceInterface.ObjectInfo, &dataSourceInfo)
	_, datasourceMasterInterface := Get(dbConnection, AnalyticsDataSourceTable, dataSourceInfo.DatasourceMaster)

	datasourceMasterInfo := DatasourceMasterInfo{}
	json.Unmarshal(datasourceMasterInterface.ObjectInfo, &datasourceMasterInfo)

	if datasourceMasterInfo.Type == "csv" {
		csvDatasource := CSVDataSource{}
		json.Unmarshal(datasourceInterface.ObjectInfo, &csvDatasource)
		url := util.InterfaceToString(updateRequest["url"])

		csvConnection := CSVFileConnection{}
		_, err := csvConnection.ImportRefreshData(url, csvDatasource, dbConnection)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Failed to process data"), ErrorCreatingObjectInformation, err.Error())
			return
		}
	} else {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Failed to process data"), ErrorCreatingObjectInformation, "Unknown Data Source")
		return
	}

	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully refresh data, please check the data tables for new data",
		Error:   0,
	})

}
