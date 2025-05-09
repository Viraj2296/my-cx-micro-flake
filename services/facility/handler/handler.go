package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"cx-micro-flake/services/facility/handler/notification"
	"cx-micro-flake/services/facility/handler/workflow_actions"
	"github.com/gin-gonic/gin"
)

type FacilityService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
	WorkfowActionHandler    *workflow_actions.ActionService
	NotificationHandler     *notification.EmailHandler
}

func (v *FacilityService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.FacilityServiceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.FacilityServiceComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *FacilityService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.FacilityServiceComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.FacilityServiceComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}

	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig

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
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *FacilityService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	itServiceGeneral := routerEngine.Group("/project/:projectId/facility")
	itServiceGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.loadFile)
	itServiceGeneral.GET("/facility_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getOverview)
	facilityService := routerEngine.Group("/project/:projectId/facility/component/:componentName")

	// table component  requests
	facilityService.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getObjects)
	facilityService.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardView)
	facilityService.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getNewRecord)            // this one support dynamic records also, if any dynamic records configured, it will load that.
	facilityService.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getRecordFormData) // this also load the dynamic records
	facilityService.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.updateResource)
	facilityService.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteResource)
	facilityService.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.createNewResource)
	facilityService.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.importObjects)
	facilityService.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getTableImportSchema)
	facilityService.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.exportObjects)
	facilityService.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSearchResults)
	facilityService.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getExportSchema)
	facilityService.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(const_util.ModuleName), v.removeInternalArrayReference)
	facilityService.GET("/group_by_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.DefaultGroupByCardView)

	// To default ExportSchema and Export Object
	facilityService.POST("/export_v1", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.exportObjectsV1)
	facilityService.GET("/export_schema_v1", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getExportSchemaV1)

	//this API will be executing as an action for internal table
	facilityService.POST("/record/:recordId/internal_table_record_ordering", middlewares.PermissionMiddleware(const_util.ModuleName), v.internalTableRecordOrdering)

	facilityService.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteValidation)

	facilityService.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.recordPOSTActionHandler)
	facilityService.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getComponentRecordTrails)
	facilityService.GET("/record_versions/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSourceObjectVersions)

}
