package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/services/scheduler/handler/const_util"
	"cx-micro-flake/services/scheduler/handler/workflow_actions"

	"github.com/gin-gonic/gin"
)

type SchedulerService struct {
	BaseService            *common.BaseService
	ComponentContentConfig component.UpstreamContentConfig
	ComponentManager       *common.ComponentManager
	WorkflowActionHandler  *workflow_actions.ActionService
}

func (v *SchedulerService) LoadInitComponents() {

	//var totalComponents int
	//for _, dbConnection := range ss.BaseService.ServiceDatabases {
	//	listOfComponents, err := GetObjects(dbConnection, ProductionOrderComponentTable)
	//	if err == nil {
	//		totalComponents = totalComponents + len(*listOfComponents)
	//	}
	//}
	//for _, dbConnection := range ss.BaseService.ServiceDatabases {
	//	listOfComponents, err := GetObjects(dbConnection, ProductionOrderComponentTable)
	//	if err == nil {
	//		ss.ComponentManager.LoadTableSchema(listOfComponents)
	//	}
	//}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *SchedulerService) InitComponents() {

	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(0)
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]

	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
	actionService := workflow_actions.ActionService{
		Logger:   v.BaseService.Logger,
		Database: dbConnection,
	}
	v.WorkflowActionHandler = &actionService
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}
func (v *SchedulerService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	schedulerGeneral := routerEngine.Group("/project/:projectId/scheduler")
	schedulerGeneral.GET("/scheduler_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getOverview)
	schedulerGeneral.GET("/order_scheduler_view", middlewares.TokenAuthMiddleware(), v.getMouldingOrderSchedulerView)
	// TODO , we need to send the similar scheduler response , but the difference, instead of machine id as a source, we need to send userId,"name":Rishoban, "avtarUrl":
	// schedulerGeneral.GET("/machine_maintenance", middlewares.TokenAuthMiddleware(), ss.getMaintenanceView)
	schedulerGeneral.GET("/maintenance_scheduler_view", middlewares.TokenAuthMiddleware(), v.getMaintenanceView)
	schedulerGeneral.GET("/mould_maintenance_scheduler_view", middlewares.TokenAuthMiddleware(), v.getMouldMaintenanceView)
	schedulerGeneral.GET("/assembly_scheduler_view", middlewares.TokenAuthMiddleware(), v.getAssemblyView)
	schedulerGeneral.GET("/tooling_scheduler_view", middlewares.TokenAuthMiddleware(), v.getToolingView)

	// TODO , we need to send the similar scheduler response , but the difference, instead of machine id as a source, we need to send userId,"name":Rishoban, "avtarUrl":
	//schedulerGeneral.GET("/mould_maintenance", middlewares.TokenAuthMiddleware(), ss.getMouldingOrderSchedulerView)
	generalComponents := routerEngine.Group("/project/:projectId/scheduler/component/:componentName")
	generalComponents.POST("/action/:actionName", middlewares.TokenAuthMiddleware(), v.handleComponentAction)
}
