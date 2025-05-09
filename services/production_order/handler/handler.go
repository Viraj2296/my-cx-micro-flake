package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/services/production_order/handler/views"
	"github.com/gin-gonic/gin"
)

type ProductionOrderService struct {
	BaseService            *common.BaseService
	ComponentContentConfig component.UpstreamContentConfig
	ComponentManager       *common.ComponentManager
	ViewManager            *views.ViewManager
}

func (v *ProductionOrderService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ProductionOrderComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ProductionOrderComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *ProductionOrderService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ProductionOrderComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ProductionOrderComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
	viewManager := views.ViewManager{
		DbConn: v.BaseService.ServiceDatabases[ProjectID],
		Logger: v.BaseService.Logger,
	}
	v.ViewManager = &viewManager
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *ProductionOrderService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	//po.generatePermissions()
	//os.Exit(0)
	productionOrders := routerEngine.Group("/project/:projectId/production_order/component/:componentName")

	productionOrderGeneral := routerEngine.Group("/project/:projectId/production_order")
	productionOrderGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), v.loadFile)

	//get the dashboard snapshots
	productionOrderGeneral.GET("/production_order_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getOverview)

	// table component  requests

	productionOrders.GET("/records", middlewares.TokenAuthMiddleware(), v.getObjects)
	productionOrders.GET("/card_view", middlewares.TokenAuthMiddleware(), v.getCardView)
	productionOrders.GET("/new_record", middlewares.TokenAuthMiddleware(), v.getNewRecord)
	productionOrders.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), v.getRecordFormData)

	productionOrders.POST("/record/:recordId/action/split_schedule", middlewares.TokenAuthMiddleware(), v.splitScheduleOrders)
	productionOrders.POST("/record/:recordId/action/reset", middlewares.TokenAuthMiddleware(), v.resetSchedule)

	productionOrders.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), v.updateResource)
	productionOrders.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), v.deleteResource)
	productionOrders.POST("/records", middlewares.TokenAuthMiddleware(), v.createNewResource)
	productionOrders.POST("/import", middlewares.TokenAuthMiddleware(), v.importObjects)
	productionOrders.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), v.getTableImportSchema)
	productionOrders.POST("/export", middlewares.TokenAuthMiddleware(), v.exportObjects)
	productionOrders.POST("/search", middlewares.TokenAuthMiddleware(), v.getSearchResults)
	productionOrders.GET("/export_schema", middlewares.TokenAuthMiddleware(), v.getExportSchema)

	productionOrders.POST("/record/:recordId/action/update_scheduler_event", middlewares.TokenAuthMiddleware(), v.updateSchedulerEvent)
	productionOrders.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteValidation)
	productionOrders.POST("/record/:recordId/action/release", middlewares.TokenAuthMiddleware(), v.releaseOrder)
	productionOrders.POST("/record/:recordId/action/hold", middlewares.TokenAuthMiddleware(), v.holdOrder)
	productionOrders.POST("/record/:recordId/action/complete", middlewares.TokenAuthMiddleware(), v.completeOrder)
	productionOrders.POST("/record/:recordId/action/abort", middlewares.TokenAuthMiddleware(), v.abortOrder)
	productionOrders.POST("/record/:recordId/action/force_stop", middlewares.TokenAuthMiddleware(), v.forceStopOrder)
	//get the records
	productionOrders.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), v.getComponentRecordTrails)

}
