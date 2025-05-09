package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type WorkHourAssignmentService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
}

func (v *WorkHourAssignmentService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, WorkHourAssignmentComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, WorkHourAssignmentComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *WorkHourAssignmentService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, WorkHourAssignmentComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, WorkHourAssignmentComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *WorkHourAssignmentService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	WorkHourAssignmentGeneral := routerEngine.Group("/project/:projectId/work_hour_assignment")
	WorkHourAssignmentGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.loadFile)
	WorkHourAssignmentGeneral.GET("/work_hour_assignment_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getItOverview)
	WorkHourAssignment := routerEngine.Group("/project/:projectId/work_hour_assignment/component/:componentName")

	// table component  requests
	WorkHourAssignment.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getObjects)
	WorkHourAssignment.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardView)
	WorkHourAssignment.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNewRecord)
	WorkHourAssignment.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getRecordFormData)
	WorkHourAssignment.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateResource)
	WorkHourAssignment.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteResource)
	WorkHourAssignment.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.createNewResource)
	WorkHourAssignment.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.importObjects)
	WorkHourAssignment.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTableImportSchema)
	WorkHourAssignment.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.exportObjects)
	WorkHourAssignment.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSearchResults)
	WorkHourAssignment.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getExportSchema)
	WorkHourAssignment.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(ModuleName), v.removeInternalArrayReference)
	WorkHourAssignment.GET("/record/:recordId/action/delete_validation", v.deleteValidation)

	WorkHourAssignment.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.recordPOSTActionHandler)
	WorkHourAssignment.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getComponentRecordTrails)
	//get the records

}
