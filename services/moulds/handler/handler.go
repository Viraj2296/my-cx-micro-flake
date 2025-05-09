package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"cx-micro-flake/services/moulds/handler/jobs"
	"cx-micro-flake/services/moulds/handler/notification"
	"cx-micro-flake/services/moulds/handler/workflow_actions"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type MouldService struct {
	BaseService                     *common.BaseService
	ComponentContentConfig          component.UpstreamContentConfig
	ComponentManager                *common.ComponentManager
	EmailNotificationDomain         string
	WorkflowActionHandler           *workflow_actions.ActionService
	LifeNotificationPoolingInterval int
}

func (v *MouldService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.MouldMasterTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.MouldMasterTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

func composeEventByMould(scheduledEventObject map[string]interface{}) map[string]interface{} {
	createEvent := make(map[string]interface{})
	createEvent["iconCls"] = "fa fa-cogs"
	createEvent["eventType"] = "mould_test_request"

	createEvent["draggable"] = true
	createEvent["module"] = "moulds"
	createEvent["componentName"] = "mould_test_request"
	createEvent["startDate"] = util.InterfaceToString(scheduledEventObject["requestTestStartDate"])
	createEvent["endDate"] = util.InterfaceToString(scheduledEventObject["requestTestEndDate"])
	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0

	return createEvent
}

func getMouldNameById(mouldId int, mouldList *[]component.GeneralObject) string {
	mouldName := ""

	for _, mould := range *mouldList {
		if mould.Id == mouldId {
			mouldMaster := database.MouldMasterInfo{}
			json.Unmarshal(mould.ObjectInfo, &mouldMaster)
			return mouldMaster.ToolNo
		}

	}

	return mouldName
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *MouldService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.MouldComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := database.GetObjects(dbConnection, const_util.MouldComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
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
	var jobService = jobs.JobService{
		Database:        dbConnection,
		Logger:          v.BaseService.Logger,
		PoolingInterval: v.LifeNotificationPoolingInterval,
		EmailHandler:    &notificationHandler,
	}
	jobService.Init()

	v.WorkflowActionHandler = &actionService
	// go jobService.SendMouldLifeNotification()
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *MouldService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	mouldGeneral := routerEngine.Group("/project/:projectId/moulds")
	mouldGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), v.loadFile)
	mouldGeneral.GET("/mould_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getOverview)
	mouldGeneral.GET("/mould_stats_summary", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getMouldsStatsSummary)

	generalComponents := routerEngine.Group("/project/:projectId/moulds/component/:componentName")

	// table component  requests
	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getObjects)

	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getNewRecord)

	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.updateResource)
	generalComponents.PUT("/record/:recordId/remove_internal_array_reference", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.removeInternalArrayReference)

	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.createNewResource)
	generalComponents.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.importObjects)
	generalComponents.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getTableImportSchema)
	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSearchResults)
	generalComponents.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getExportSchema)
	generalComponents.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.exportObjects)
	generalComponents.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getComponentRecordTrails)
	generalComponents.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteValidation)

	// action handler
	generalComponents.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordPOSTActionHandler)
	generalComponents.PUT("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordPUTActionHandler)

}
