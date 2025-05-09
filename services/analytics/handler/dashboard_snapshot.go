package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/analytics"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (as *AnalyticsService) generateChartData(visualisationType string, dbConnection *gorm.DB, widgetId int, listOfFilterInfo []analytics.DashboardWidgetFilterInfo) (error, interface{}) {
	err, generalObject := Get(dbConnection, AnalyticsWidgetTable, widgetId)
	if err != nil {
		return err, nil
	}
	widgetInfo := analytics.WidgetInfo{}
	json.Unmarshal(generalObject.ObjectInfo, &widgetInfo)
	seriesMapping := widgetInfo.SeriesMapping
	filterArray := widgetInfo.Filters
	var param = make(map[string]string)
	for _, dashboardFilter := range listOfFilterInfo {
		defaultValue := dashboardFilter.DefaultValues
		for _, individualFilter := range filterArray {
			if individualFilter.FilterId == dashboardFilter.FilterId {
				parameter := individualFilter.Parameter
				param[parameter] = defaultValue
			}
		}

	}

	if len(param) > 0 {
		// no param defined
		for key, value := range param {
			widgetInfo.Query = strings.Replace(widgetInfo.Query, "{{"+key+"}}", util.InterfaceToString(value), -1)
		}

	}
	var queryResults []map[string]interface{}
	transaction := dbConnection.Raw(widgetInfo.Query).Scan(&queryResults)
	err = transaction.Error
	if err != nil {
		return err, nil
	}

	baseBuilder := analytics.BaseVisualisationBuilder{}
	baseBuilder.QueryResults = queryResults
	var seriesMappingFields map[string]interface{}
	json.Unmarshal(seriesMapping, &seriesMappingFields)
	baseBuilder.SeriesMapping = seriesMappingFields
	baseBuilder.Visualisation = widgetInfo.Visualization
	baseBuilder.Init()
	var generatedSeries interface{}
	if visualisationType == "chart" {
		responseBuilder := analytics.ChartVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
		_, generatedSeries = responseBuilder.BuildResponse()
	} else if visualisationType == "gauge" {
		responseBuilder := analytics.GaugeVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
		_, generatedSeries = responseBuilder.BuildResponse()
	} else if visualisationType == "number_card" {
		responseBuilder := analytics.NumberCardVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
		_, generatedSeries = responseBuilder.BuildResponse()
	} else if visualisationType == "timeline" {
		responseBuilder := analytics.TimelineVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
		_, generatedSeries = responseBuilder.BuildResponse()
	} else if visualisationType == "table" {
		responseBuilder := analytics.TableVisualisationBuilder{BaseVisualisationBuilder: &baseBuilder}
		_, generatedSeries = responseBuilder.BuildResponse()
	} else if visualisationType == "bullet_graph" {
		responseBuilder := analytics.BulletGraphBuilder{BaseVisualisationBuilder: &baseBuilder}
		_, generatedSeries = responseBuilder.BuildResponse()
	}

	visualObjects := widgetInfo.Visualization
	visualObjects["series"] = generatedSeries
	return nil, visualObjects

}
func (as *AnalyticsService) getDashboardSnapshots(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	userId := common.GetUserId(ctx)
	userIdStr := strconv.Itoa(userId)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	listOfDashboards, _ := GetObjects(dbConnection, AnalyticsDashboardTable)

	var listOfDashboardResponse = make([]analytics.DashboardWidgetSnapshotResponse, 0)

	// Get the permissions for the dashboard
	condition := " object_info ->> '$.assignedUserId' = " + userIdStr
	listOfDashboardPermission, _ := GetConditionalObjects(dbConnection, AnalyticsDashboardPermissionTable, condition)

	as.BaseService.Logger.Info("number of dashboard permission", zap.Any("len", len(*listOfDashboardPermission)))

	if len(*listOfDashboardPermission) == 0 {
		ctx.JSON(http.StatusOK, listOfDashboardResponse)
		return
	}

	for _, dashboardInterface := range *listOfDashboards {
		permissionFlag := false
		for _, permission := range *listOfDashboardPermission {
			dashboardPermissionInfo := make(map[string]interface{})
			json.Unmarshal(permission.ObjectInfo, &dashboardPermissionInfo)
			dashboardId := util.InterfaceToInt(dashboardPermissionInfo["dashboardId"])
			permissionLeval := util.InterfaceToInt(dashboardPermissionInfo["permissionLevel"])

			if dashboardId == dashboardInterface.Id {
				if permissionLeval == ReadPermission {
					permissionFlag = true
					break
				}
			}
		}

		if !permissionFlag {
			continue
		}
		// list of dashboard
		dashboardInfo := analytics.DashboardInfo{}
		json.Unmarshal(dashboardInterface.ObjectInfo, &dashboardInfo)
		if !dashboardInfo.SnapshotEnabled {
			continue
		}

		var tmpDashboardFilters []*analytics.Filter
		var individualWidgetFilterArray []*analytics.Filter
		var dashboardWidgetSnapshots []analytics.DashboardWidgetSnapshot
		for _, dashboardWidget := range dashboardInfo.DashboardWidgets {
			// list of widgets

			err, generalWidgetObject := Get(dbConnection, AnalyticsWidgetTable, dashboardWidget.WidgetId)
			if err != nil {
				as.BaseService.Logger.Error("failed load widget,", zap.String("error", err.Error()))
				return
			}

			widgetInfo := analytics.WidgetInfo{}
			json.Unmarshal(generalWidgetObject.ObjectInfo, &widgetInfo)

			//each dashboard widgets are having filter based on whether dashbaord widget or widget filter
			for _, dashboardWidgetFilterInfo := range dashboardWidget.DashboardWidgetFilterInfo {
				valueSource := dashboardWidgetFilterInfo.ValueSource
				filterId := dashboardWidgetFilterInfo.FilterId
				// the one configured during widget creation
				widgetFilter := widgetInfo.GetFilter(filterId)
				fmt.Println(widgetInfo.Name)

				widgetFilter.WidgetId = dashboardWidget.WidgetId
				if valueSource == "dashboard" {
					// only dashboard values are going for dropped down
					dashboardFilter := widgetInfo.GetDashboardFilter(widgetFilter)
					if dashboardWidgetFilterInfo.DefaultValues != "" {
						dashboardFilter.DefaultValues = dashboardWidgetFilterInfo.DefaultValues
					}
					if dashboardWidgetFilterInfo.Title != "" {
						dashboardFilter.FilterLabel = dashboardWidgetFilterInfo.Title
					}
					tmpDashboardFilters = append(tmpDashboardFilters, dashboardFilter)
				} else if valueSource == "widget" {
					individualWidgetFilters := widgetInfo.GetDashboardFilter(widgetFilter)
					if dashboardWidgetFilterInfo.DefaultValues != "" {
						individualWidgetFilters.DefaultValues = dashboardWidgetFilterInfo.DefaultValues
					}
					if dashboardWidgetFilterInfo.Title != "" {
						individualWidgetFilters.FilterLabel = dashboardWidgetFilterInfo.Title
					}
					individualWidgetFilterArray = append(individualWidgetFilterArray, individualWidgetFilters)

				}

			}

			// Get dbConnection using datasource id in widget
			_, generalDataSourceObject := Get(dbConnection, AnalyticsDataSourceTable, widgetInfo.Datasource)
			dataSourceInfo := make(map[string]interface{})
			json.Unmarshal(generalDataSourceObject.ObjectInfo, &dataSourceInfo)
			// connectionFields := dataSourceInfo["connectionParam"]

			// serialisedData, _ := json.Marshal(connectionFields)
			// datasourceConn := createDbConnectionByDatasource(serialisedData)

			var visualisationObject interface{}
			fmt.Println("dashboardWidget: ", dashboardWidget)
			fmt.Println("generating chart data , chart type :", widgetInfo.GetVisualisationType())
			err, visualisationObject = as.generateChartData(widgetInfo.GetVisualisationType(), dbConnection, generalWidgetObject.Id, dashboardWidget.DashboardWidgetFilterInfo)
			if err != nil {
				fmt.Println("generate series failed", err.Error())
				continue
			}
			fmt.Println("visualisationObject: ", visualisationObject)
			dashboardWidgetSnapshot := analytics.DashboardWidgetSnapshot{}
			dashboardWidgetSnapshot.WidgetId = dashboardWidget.WidgetId
			dashboardWidgetSnapshot.Width = dashboardWidget.Width
			dashboardWidgetSnapshot.Cols = dashboardWidget.Cols
			dashboardWidgetSnapshot.Name = widgetInfo.Name
			dashboardWidgetSnapshot.Id = dashboardWidget.Id
			dashboardWidgetSnapshot.X = dashboardWidget.X
			dashboardWidgetSnapshot.Y = dashboardWidget.Y
			dashboardWidgetSnapshot.Rows = dashboardWidget.Rows
			dashboardWidgetSnapshot.Height = dashboardWidget.Height
			dashboardWidgetSnapshot.VisualizationType = dashboardWidget.VisualizationType
			dashboardWidgetSnapshot.Visualization = visualisationObject
			dashboardWidgetSnapshot.WidgetFilterInfo = individualWidgetFilterArray
			dashboardWidgetSnapshots = append(dashboardWidgetSnapshots, dashboardWidgetSnapshot)
			// now we need to fill the data for widget

		}

		// now go through all the dashboard filters and find the same name one.
		var filterLabels []string
		var filterWidgetMapping = make(map[string][]int)
		for _, tmpDashboardFilter := range tmpDashboardFilters {

			if _, ok := filterWidgetMapping[tmpDashboardFilter.FilterLabel]; ok {
				filterWidgetMapping[tmpDashboardFilter.FilterLabel] = append(filterWidgetMapping[tmpDashboardFilter.FilterLabel], tmpDashboardFilter.WidgetId)
			} else {
				var listOfWidgets []int
				listOfWidgets = append(listOfWidgets, tmpDashboardFilter.WidgetId)
				filterWidgetMapping[tmpDashboardFilter.FilterLabel] = listOfWidgets
			}
		}
		var newDashboardFilters []analytics.DashboardFilter
		for _, tmpDashboardFilter := range tmpDashboardFilters {
			if !util.StringArrayContains(filterLabels, tmpDashboardFilter.FilterLabel) {
				as.BaseService.Logger.Info("not in the list, so adding as a dashboard filter label", zap.String("filter_label", tmpDashboardFilter.FilterLabel))
				dashboardFilter := analytics.DashboardFilter{}
				dashboardFilter.FilterType = tmpDashboardFilter.FilterType
				dashboardFilter.FilterValues = tmpDashboardFilter.FilterValues
				dashboardFilter.Param = tmpDashboardFilter.Parameter
				dashboardFilter.MultiSelect = tmpDashboardFilter.MultiSelect
				dashboardFilter.DefaultValues = tmpDashboardFilter.DefaultValues
				dashboardFilter.FilterLabel = tmpDashboardFilter.FilterLabel
				dashboardFilter.QuotationFormat = tmpDashboardFilter.QuotationFormat
				dashboardFilter.LinkedWidgets = filterWidgetMapping[tmpDashboardFilter.FilterLabel]
				newDashboardFilters = append(newDashboardFilters, dashboardFilter)
			}
			filterLabels = append(filterLabels, tmpDashboardFilter.FilterLabel)
		}
		dashboardWidgetSnapshotInfo := analytics.DashboardWidgetSnapshotInfo{}
		dashboardWidgetSnapshotInfo.RefreshInterval = dashboardInfo.RefreshInterval
		dashboardWidgetSnapshotInfo.Name = dashboardInfo.Name
		dashboardWidgetSnapshotInfo.DashboardWidgets = dashboardWidgetSnapshots
		dashboardWidgetSnapshotInfo.DashboardFilters = newDashboardFilters

		dashboardSnapshotResponse := analytics.DashboardWidgetSnapshotResponse{DashboardId: dashboardInterface.Id, DashboardInfo: dashboardWidgetSnapshotInfo}
		listOfDashboardResponse = append(listOfDashboardResponse, dashboardSnapshotResponse)
	}

	ctx.JSON(http.StatusOK, listOfDashboardResponse)

}

func (as *AnalyticsService) getDashboardSnapshotById(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordId)
	userId := common.GetUserId(ctx)
	userIdStr := strconv.Itoa(userId)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	_, dashboardInterface := Get(dbConnection, AnalyticsDashboardTable, intRecordId)

	var listOfDashboardResponse = make([]analytics.DashboardWidgetSnapshotResponse, 0)

	// Get the permissions for the dashboard
	condition := " object_info ->> '$.assignedUserId' = " + userIdStr
	listOfDashboardPermission, _ := GetConditionalObjects(dbConnection, AnalyticsDashboardPermissionTable, condition)

	fmt.Println("=============================================== Loading Dashboard permission start ====================")
	fmt.Println(len(*listOfDashboardPermission))
	fmt.Println("=============================================== Loading Dashboard permission done ====================")
	if len(*listOfDashboardPermission) == 0 {
		ctx.JSON(http.StatusOK, listOfDashboardResponse)
		return
	}

	permissionFlag := false
	for _, permission := range *listOfDashboardPermission {
		dashboardPermissionInfo := make(map[string]interface{})
		json.Unmarshal(permission.ObjectInfo, &dashboardPermissionInfo)
		dashboardId := util.InterfaceToInt(dashboardPermissionInfo["dashboardId"])
		permissionLeval := util.InterfaceToInt(dashboardPermissionInfo["permissionLevel"])

		if dashboardId == dashboardInterface.Id {
			if permissionLeval == ReadPermission {
				permissionFlag = true
				break
			}
		}
	}

	if !permissionFlag {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("permission denied"), ErrorGettingIndividualObjectInformation)
		return

	}
	// list of dashboard
	dashboardInfo := analytics.DashboardInfo{}
	json.Unmarshal(dashboardInterface.ObjectInfo, &dashboardInfo)
	if !dashboardInfo.SnapshotEnabled {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("snapshot not enabled"), ErrorGettingIndividualObjectInformation)
		return
	}

	var tmpDashboardFilters []*analytics.Filter
	var individualWidgetFilterArray []*analytics.Filter
	var dashboardWidgetSnapshots []analytics.DashboardWidgetSnapshot
	for _, dashboardWidget := range dashboardInfo.DashboardWidgets {
		// list of widgets

		err, generalWidgetObject := Get(dbConnection, AnalyticsWidgetTable, dashboardWidget.WidgetId)
		if err != nil {
			as.BaseService.Logger.Error("failed load widget,", zap.String("error", err.Error()))
			return
		}
		widgetInfo := analytics.WidgetInfo{}
		json.Unmarshal(generalWidgetObject.ObjectInfo, &widgetInfo)

		//each dashboard widgets are having filter based on whether dashbaord widget or widget filter
		for _, dashboardWidgetFilterInfo := range dashboardWidget.DashboardWidgetFilterInfo {
			valueSource := dashboardWidgetFilterInfo.ValueSource
			filterId := dashboardWidgetFilterInfo.FilterId
			// the one configured during widget creation
			widgetFilter := widgetInfo.GetFilter(filterId)
			fmt.Println(widgetInfo.Name)

			widgetFilter.WidgetId = dashboardWidget.WidgetId
			if valueSource == "dashboard" {
				// only dashboard values are going for dropped down
				dashboardFilter := widgetInfo.GetDashboardFilter(widgetFilter)
				if dashboardWidgetFilterInfo.DefaultValues != "" {
					dashboardFilter.DefaultValues = dashboardWidgetFilterInfo.DefaultValues
				}
				if dashboardWidgetFilterInfo.Title != "" {
					dashboardFilter.FilterLabel = dashboardWidgetFilterInfo.Title
				}
				tmpDashboardFilters = append(tmpDashboardFilters, dashboardFilter)
			} else if valueSource == "widget" {
				individualWidgetFilters := widgetInfo.GetDashboardFilter(widgetFilter)
				if dashboardWidgetFilterInfo.DefaultValues != "" {
					individualWidgetFilters.DefaultValues = dashboardWidgetFilterInfo.DefaultValues
				}
				if dashboardWidgetFilterInfo.Title != "" {
					individualWidgetFilters.FilterLabel = dashboardWidgetFilterInfo.Title
				}
				individualWidgetFilterArray = append(individualWidgetFilterArray, individualWidgetFilters)

			}

		}

		// Get dbConnection using datasource id in widget
		_, generalDataSourceObject := Get(dbConnection, AnalyticsDataSourceTable, widgetInfo.Datasource)
		dataSourceInfo := make(map[string]interface{})
		json.Unmarshal(generalDataSourceObject.ObjectInfo, &dataSourceInfo)
		// connectionFields := dataSourceInfo["connectionParam"]

		// serialisedData, _ := json.Marshal(connectionFields)
		// datasourceConn := createDbConnectionByDatasource(serialisedData)

		var visualisationObject interface{}
		fmt.Println("dashboardWidget: ", dashboardWidget)
		fmt.Println("generating chart data , chart type :", widgetInfo.GetVisualisationType())
		err, visualisationObject = as.generateChartData(widgetInfo.GetVisualisationType(), dbConnection, generalWidgetObject.Id, dashboardWidget.DashboardWidgetFilterInfo)
		if err != nil {
			fmt.Println("generate series failed", err.Error())
			continue
		}
		fmt.Println("visualisationObject: ", visualisationObject)
		dashboardWidgetSnapshot := analytics.DashboardWidgetSnapshot{}
		dashboardWidgetSnapshot.WidgetId = dashboardWidget.WidgetId
		dashboardWidgetSnapshot.Width = dashboardWidget.Width
		dashboardWidgetSnapshot.Cols = dashboardWidget.Cols
		dashboardWidgetSnapshot.Name = widgetInfo.Name
		dashboardWidgetSnapshot.Id = dashboardWidget.Id
		dashboardWidgetSnapshot.X = dashboardWidget.X
		dashboardWidgetSnapshot.Y = dashboardWidget.Y
		dashboardWidgetSnapshot.Rows = dashboardWidget.Rows
		dashboardWidgetSnapshot.Height = dashboardWidget.Height
		dashboardWidgetSnapshot.VisualizationType = dashboardWidget.VisualizationType
		dashboardWidgetSnapshot.Visualization = visualisationObject
		dashboardWidgetSnapshot.WidgetFilterInfo = individualWidgetFilterArray
		dashboardWidgetSnapshots = append(dashboardWidgetSnapshots, dashboardWidgetSnapshot)
		// now we need to fill the data for widget

	}

	// now go through all the dashboard filters and find the same name one.
	var filterLabels []string
	var filterWidgetMapping = make(map[string][]int)
	for _, tmpDashboardFilter := range tmpDashboardFilters {

		if _, ok := filterWidgetMapping[tmpDashboardFilter.FilterLabel]; ok {
			filterWidgetMapping[tmpDashboardFilter.FilterLabel] = append(filterWidgetMapping[tmpDashboardFilter.FilterLabel], tmpDashboardFilter.WidgetId)
		} else {
			var listOfWidgets []int
			listOfWidgets = append(listOfWidgets, tmpDashboardFilter.WidgetId)
			filterWidgetMapping[tmpDashboardFilter.FilterLabel] = listOfWidgets
		}
	}
	var newDashboardFilters []analytics.DashboardFilter
	for _, tmpDashboardFilter := range tmpDashboardFilters {
		if !util.StringArrayContains(filterLabels, tmpDashboardFilter.FilterLabel) {
			as.BaseService.Logger.Info("not in the list, so adding as a dashboard filter label", zap.String("filter_label", tmpDashboardFilter.FilterLabel))
			dashboardFilter := analytics.DashboardFilter{}
			dashboardFilter.FilterType = tmpDashboardFilter.FilterType
			dashboardFilter.FilterValues = tmpDashboardFilter.FilterValues
			dashboardFilter.Param = tmpDashboardFilter.Parameter
			dashboardFilter.MultiSelect = tmpDashboardFilter.MultiSelect
			dashboardFilter.DefaultValues = tmpDashboardFilter.DefaultValues
			dashboardFilter.FilterLabel = tmpDashboardFilter.FilterLabel
			dashboardFilter.QuotationFormat = tmpDashboardFilter.QuotationFormat
			dashboardFilter.LinkedWidgets = filterWidgetMapping[tmpDashboardFilter.FilterLabel]
			newDashboardFilters = append(newDashboardFilters, dashboardFilter)
		}
		filterLabels = append(filterLabels, tmpDashboardFilter.FilterLabel)
	}
	dashboardWidgetSnapshotInfo := analytics.DashboardWidgetSnapshotInfo{}
	dashboardWidgetSnapshotInfo.RefreshInterval = dashboardInfo.RefreshInterval
	dashboardWidgetSnapshotInfo.Name = dashboardInfo.Name
	dashboardWidgetSnapshotInfo.DashboardWidgets = dashboardWidgetSnapshots
	dashboardWidgetSnapshotInfo.DashboardFilters = newDashboardFilters

	dashboardSnapshotResponse := analytics.DashboardWidgetSnapshotResponse{DashboardId: dashboardInterface.Id, DashboardInfo: dashboardWidgetSnapshotInfo}
	listOfDashboardResponse = append(listOfDashboardResponse, dashboardSnapshotResponse)

	ctx.JSON(http.StatusOK, listOfDashboardResponse)

}

// func createDbConnectionByDatasource(connectionParameter datatypes.JSON) *gorm.DB {
// 	mysqlConnectionParameter := MysqlConnectionParameters{}

// 	json.Unmarshal(connectionParameter, &mysqlConnectionParameter)

// 	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConnectionParameter.Username, mysqlConnectionParameter.Password, mysqlConnectionParameter.HostName, mysqlConnectionParameter.Port, mysqlConnectionParameter.Schema)
// 	dbConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
// 	if err != nil {
// 		return nil
// 	}

// 	return dbConnection
// }
