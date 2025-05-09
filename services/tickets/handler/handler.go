package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type TicketsService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
}

func (ts *TicketsService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range ts.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, TicketsServiceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	ts.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range ts.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, TicketsServiceComponentTable)
		if err == nil {
			ts.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (ts *TicketsService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range ts.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, TicketsServiceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	ts.ComponentManager = &common.ComponentManager{}
	ts.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range ts.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, TicketsServiceComponentTable)
		if err == nil {
			ts.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	// init the component config
	ts.ComponentManager.ComponentContentConfig = ts.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (ts *TicketsService) InitRouter(routerEngine *gin.Engine) {
	ts.InitComponents()
	TicketsServiceGeneral := routerEngine.Group("/project/:projectId/tickets")
	TicketsServiceGeneral.POST("/loadFile", ts.loadFile)
	TicketsServiceGeneral.GET("/tickets_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getOverview)
	ticketsService := routerEngine.Group("/project/:projectId/tickets/component/:componentName")

	// table component  requests
	ticketsService.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getObjects)
	ticketsService.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getCardView)
	ticketsService.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getNewRecord)
	ticketsService.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getRecordFormData)
	ticketsService.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.updateResource)
	ticketsService.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.deleteResource)
	ticketsService.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.createNewResource)
	ticketsService.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.importObjects)
	ticketsService.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getTableImportSchema)
	ticketsService.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.exportObjects)
	ticketsService.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getSearchResults)
	ticketsService.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.getExportSchema)
	ticketsService.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(ModuleName), ts.removeInternalArrayReference)

	ticketsService.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ts.recordPOSTActionHandler)
	ticketsService.GET("/record_messages/:recordId", ts.getComponentRecordTrails)
	//get the records

}
