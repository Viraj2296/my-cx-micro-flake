package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"cx-micro-flake/services/it/handler/notification"
	"cx-micro-flake/services/it/handler/workflow_actions"
	"github.com/gin-gonic/gin"
)

type ITService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
	WorkfowActionHandler    *workflow_actions.ActionService
	NotificationHandler     *notification.EmailHandler
}

func (v *ITService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.ITServiceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.ITServiceComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *ITService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.ITServiceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.ITServiceComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	notificationHandler := notification.EmailHandler{
		Logger:                  v.BaseService.Logger,
		EmailNotificationDomain: v.EmailNotificationDomain,
		ComponentManager:        v.ComponentManager,
	}
	actionService := workflow_actions.ActionService{
		Logger:       v.BaseService.Logger,
		Database:     dbConnection,
		EmailHandler: &notificationHandler,
	}
	v.WorkfowActionHandler = &actionService
	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *ITService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	itServiceGeneral := routerEngine.Group("/project/:projectId/it")
	itServiceGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.loadFile)
	itServiceGeneral.GET("/it_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getOverview)
	itService := routerEngine.Group("/project/:projectId/it/component/:componentName")

	// table component  requests
	itService.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getObjects)
	itService.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardView)
	itService.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getNewRecord)            // this one support dynamic records also, if any dynamic records configured, it will load that.
	itService.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getRecordFormData) // this also load the dynamic records
	itService.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.updateResource)
	itService.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteResource)
	itService.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.createNewResource)
	itService.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.importObjects)
	itService.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getTableImportSchema)
	itService.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.exportObjects)
	itService.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSearchResults)
	itService.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getExportSchema)
	itService.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(const_util.ModuleName), v.removeInternalArrayReference)

	//this API will be executing as an action for internal table
	itService.POST("/record/:recordId/internal_table_record_ordering", middlewares.PermissionMiddleware(const_util.ModuleName), v.internalTableRecordOrdering)

	itService.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteValidation)

	itService.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.recordPOSTActionHandler)
	itService.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getComponentRecordTrails)
	itService.GET("/record_versions/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSourceObjectVersions)
	itService.GET("/group_by_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.DefaultGroupByCardView)
	// in-line email action handling.
	// action token would contain record id and action name
	//itService.GET("/token_parameters/:actionToken", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTokenParameters)
	//get the records

}
