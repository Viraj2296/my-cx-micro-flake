package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type ToolingService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
}

func (v *ToolingService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ToolingComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ToolingComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *ToolingService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ToolingComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ToolingComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *ToolingService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	toolingGeneral := routerEngine.Group("/project/:projectId/tooling")
	toolingGeneral.GET("/tooling_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getOverview)
	toolingGeneral.GET("/project_task_kanban_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTasksWorkOrderTaskKanbanView)

	toolingRouting := routerEngine.Group("/project/:projectId/tooling/component/:componentName")

	toolingRouting.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardView)
	toolingRouting.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNewRecord)
	toolingRouting.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getRecordFormData)
	toolingRouting.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getObjects)

	toolingRouting.POST("/record/:recordId/action/kanban_move_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.kanbanMoveTask)

	toolingRouting.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateResource)
	toolingRouting.POST("/records/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getGroupBy)
	toolingRouting.POST("/card_view/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardViewGroupBy)
	toolingRouting.PUT("/record/:recordId/remove_internal_array_reference", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.removeInternalArrayReference)
	toolingRouting.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteResource)
	toolingRouting.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.createNewResource)
	toolingRouting.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSearchResults)
	toolingRouting.POST("/record/:recordId/action/complete_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.completeTask)
	toolingRouting.POST("/record/:recordId/action/checkin_task", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.checkInTask)
	toolingRouting.POST("/record/:recordId/action/approve", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.approveTask)
	toolingRouting.POST("/record/:recordId/action/reject", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.rejectTask)

	toolingRouting.GET("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.recordGetActionHandler)

	toolingRouting.POST("/record/:recordId/action/activate_sprint", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.activateSprint)
	toolingRouting.POST("/record/:recordId/action/complete_sprint", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.completeSprint)

	toolingRouting.POST("/record/:recordId/action/kick_off", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.kickOffProject)

	toolingRouting.GET("/record_messages/:recordId", v.getComponentRecordTrails)
}
