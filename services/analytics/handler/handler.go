package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type AnalyticsService struct {
	BaseService               *common.BaseService
	ComponentContentConfig    component.UpstreamContentConfig
	ComponentManager          *common.ComponentManager
	DatasourceConnectionCache map[int]interface{}
}

func (as *AnalyticsService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, AnalyticsComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	as.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, AnalyticsComponentTable)
		if err == nil {
			as.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (as *AnalyticsService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, AnalyticsComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	as.ComponentManager = &common.ComponentManager{}
	as.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, AnalyticsComponentTable)
		if err == nil {
			as.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	as.DatasourceConnectionCache = make(map[int]interface{}, 0)
	// init the component config
	as.ComponentManager.ComponentContentConfig = as.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (as *AnalyticsService) InitRouter(routerEngine *gin.Engine) {
	as.InitComponents()
	analyticsGeneral := routerEngine.Group("/project/:projectId/analytics")

	analyticsGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), as.loadFile)
	analyticsGeneral.POST("/query_response", middlewares.TokenAuthMiddleware(), as.queryResponse)
	analyticsGeneral.GET("/dashboard_snapshots", middlewares.TokenAuthMiddleware(), as.getDashboardSnapshots)
	analyticsGeneral.GET("/dashboard_snapshot/record/:recordId", middlewares.TokenAuthMiddleware(), as.getDashboardSnapshotById)
	// analyticsGeneral.GET("/database_schema", middlewares.TokenAuthMiddleware(), as.getDataBaseSchema)
	analyticsGeneral.POST("/action/dashboard_widget_visualisation", middlewares.TokenAuthMiddleware(), as.getDashboardWidgetVisualisation)

	analyticsGeneral.POST("/:datasourceName/test_connection", as.testDatasourceConnection)
	analyticsGeneral.GET("/datasources_master_group_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getDatasourceMasterCardView)
	generalComponents := routerEngine.Group("/project/:projectId/analytics/component/:componentName")

	// table component  requests
	generalComponents.GET("/records", as.getObjects)
	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), as.getCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), as.getNewRecord)
	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), as.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), as.updateResource)
	generalComponents.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), as.deleteValidation)
	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), as.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), as.createNewResource)
	generalComponents.POST("/import", middlewares.TokenAuthMiddleware(), as.importObjects)
	generalComponents.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), as.getTableImportSchema)
	generalComponents.POST("/export", middlewares.TokenAuthMiddleware(), as.exportObjects)
	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), as.getSearchResults)
	generalComponents.GET("/export_schema", middlewares.TokenAuthMiddleware(), as.getExportSchema)
	generalComponents.POST("/refresh_export_schema", middlewares.TokenAuthMiddleware(), as.refreshExportSchema)
	generalComponents.PUT("/record/:recordId/remove_internal_array_reference", as.removeInternalArrayReference)
	generalComponents.POST("/record/:recordId/action/dashboard_widget_visualisation", middlewares.TokenAuthMiddleware(), as.getDashboardWidgetVisualisation)

	generalComponents.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), as.recordPOSTActionHandler)
	//get the records
	generalComponents.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), as.getComponentRecordTrails)

	generalComponents.GET("/record/:recordId/database_schema", middlewares.TokenAuthMiddleware(), as.getDataBaseSchema)

}
