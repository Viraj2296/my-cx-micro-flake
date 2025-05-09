package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/analytics"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"strings"
)

type ExecuteDashboardWidgetsVisualisationRequest struct {
	WidgetId int                    `json:"widgetId"`
	Param    map[string]interface{} `json:"param"`
}

func getDatasourceType(dbConnection *gorm.DB, id int) string {
	err, datasourceMasterGeneralObject := Get(dbConnection, AnalyticsDatasourcesMasterTable, id)
	if err == nil {
		analyticsDatasourceMaster := AnalyticsDatasourcesMaster{ObjectInfo: datasourceMasterGeneralObject.ObjectInfo}
		return analyticsDatasourceMaster.getDatasourceMasterInfo().Type
	}
	return ""
}
func getExistingDatabaseTables(service *AnalyticsService, objectFields map[string]interface{}, functionArguments []component.FunctionArguments) map[string][]string {
	var tableFields = make(map[string][]string)
	// take the value from
	var recordId int
	if functionArguments[0].Value == nil {
		field := util.InterfaceToString(functionArguments[0].Field)
		recordId = util.InterfaceToInt(objectFields[field])
	}

	dbConnection := service.BaseService.ServiceDatabases[ProjectID]
	err, datasourceGeneralObject := Get(dbConnection, AnalyticsDataSourceTable, recordId)
	fmt.Println("datasourceGeneralObject: ", datasourceGeneralObject, "err:", err)
	if err != nil {
		return tableFields
	}
	analyticsDatasource := AnalyticsDatasource{ObjectInfo: datasourceGeneralObject.ObjectInfo}
	fmt.Println("analyticsDatasource: ", analyticsDatasource.getDatasourceInfo())
	if getDatasourceType(dbConnection, analyticsDatasource.getDatasourceInfo().DatasourceMaster) == "mysql" {
		// get the connection param by casting the interface
		//datasourceConnection := service.DatasourceConnectionCache[recordId].(*gorm.DB)
		serialisedConnectionParam, _ := json.Marshal(analyticsDatasource.getDatasourceInfo().ConnectionParam)
		mysqlConnectionParam := MysqlConnectionParameters{}
		json.Unmarshal(serialisedConnectionParam, &mysqlConnectionParam)
		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConnectionParam.Username, mysqlConnectionParam.Password, mysqlConnectionParam.HostName, mysqlConnectionParam.Port, mysqlConnectionParam.Schema)
		datasourceConnection, _ := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
		tableInformationSchema := "SELECT * FROM information_schema.tables WHERE table_schema =\"" + mysqlConnectionParam.Schema + "\""
		var queryResults []map[string]interface{}
		datasourceConnection.Raw(tableInformationSchema).Scan(&queryResults)
		var listOfTables []string
		for _, result := range queryResults {
			tableName := result["TABLE_NAME"]
			listOfTables = append(listOfTables, util.InterfaceToString(tableName))
		}
		fmt.Println("listOfTables: ", listOfTables)
		for _, tableName := range listOfTables {
			internalDataQuery := "select column_name  from information_schema.columns where table_schema =\"" + mysqlConnectionParam.Schema + "\" and table_name = \"" + tableName + "\""
			var columnNames []string
			err := datasourceConnection.Raw(internalDataQuery).Scan(&columnNames).Error
			fmt.Println("err: :err: ", err)
			fmt.Println("internalDataQuery: ", internalDataQuery)

			if err == nil {
				tableFields[tableName] = columnNames
			}
		}

	}
	return tableFields

}

func (as *AnalyticsService) getDashboardWidgetVisualisation(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	var executeDashboardWidgetsVisualisationRequest []ExecuteDashboardWidgetsVisualisationRequest
	if err := ctx.ShouldBindBodyWith(&executeDashboardWidgetsVisualisationRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var dashboardWidgetVisualisationResponse = make([]interface{}, 0)
	for _, widgetRequests := range executeDashboardWidgetsVisualisationRequest {
		widgetId := widgetRequests.WidgetId
		err, generalWidgetObject := Get(dbConnection, AnalyticsWidgetTable, widgetId)
		if err != nil {
			// lets send empty visualisation objects, don't send empty
		}
		analyticsWidget := AnalyticsWidget{ObjectInfo: generalWidgetObject.ObjectInfo}
		query := analyticsWidget.getWidgetInfo().Query

		if len(widgetRequests.Param) > 0 {
			// no param defined
			for key, value := range widgetRequests.Param {
				query = strings.Replace(query, "{{"+key+"}}", util.InterfaceToString(value), -1)
			}

		}
		var queryResults []map[string]interface{}
		err = dbConnection.Raw(query).Scan(&queryResults).Error
		if err != nil {
			as.BaseService.Logger.Error("error happened during widget query execution", zap.String("query", query))
			continue
		}
		seriesMapping := analyticsWidget.getWidgetInfo().SeriesMapping
		baseBuilder := analytics.BaseVisualisationBuilder{}
		baseBuilder.QueryResults = queryResults
		var seriesMappingFields map[string]interface{}
		json.Unmarshal(seriesMapping, &seriesMappingFields)
		baseBuilder.SeriesMapping = seriesMappingFields
		baseBuilder.Visualisation = analyticsWidget.getWidgetInfo().Visualization
		baseBuilder.Init()

		var generatedSeries interface{}
		if analyticsWidget.getWidgetInfo().VisualizationType == "gauge" {
			responseBuilder := analytics.GaugeVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
			_, generatedSeries = responseBuilder.BuildResponse()
		} else if analyticsWidget.getWidgetInfo().VisualizationType == "table" {
			responseBuilder := analytics.TableVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
			_, generatedSeries = responseBuilder.BuildResponse()
		} else if analyticsWidget.getWidgetInfo().VisualizationType == "bullet_graph" {
			responseBuilder := analytics.BulletGraphBuilder{BaseVisualisationBuilder: &baseBuilder}
			_, generatedSeries = responseBuilder.BuildResponse()
		} else {
			responseBuilder := analytics.ChartVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
			_, generatedSeries = responseBuilder.BuildResponse()
		}

		visualObjects := analyticsWidget.getWidgetInfo().Visualization
		visualObjects["series"] = generatedSeries
		var internalResponse = make(map[string]interface{})
		internalResponse["id"] = widgetId
		internalResponse["visualisation"] = visualObjects
		dashboardWidgetVisualisationResponse = append(dashboardWidgetVisualisationResponse, internalResponse)

	}

	ctx.JSON(http.StatusOK, dashboardWidgetVisualisationResponse)
}

func (as *AnalyticsService) recordPOSTActionHandler(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == "refresh_data" {
		as.RefreshData(ctx)
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
