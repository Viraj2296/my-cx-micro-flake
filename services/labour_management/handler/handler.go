package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/services/labour_management/handler/cache"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"cx-micro-flake/services/labour_management/handler/jobs"
	"cx-micro-flake/services/labour_management/handler/notification"
	"cx-micro-flake/services/labour_management/handler/workflow_actions"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LabourManagementService struct {
	BaseService              *common.BaseService
	ComponentContentConfig   component.UpstreamContentConfig
	ComponentManager         *common.ComponentManager
	EmailNotificationDomain  string
	WorkflowActionHandler    *workflow_actions.ActionService
	MachineStatsPoolingTimer int
	EscalationEmailTemplate  string
}

func (v *LabourManagementService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.LabourManagementComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.LabourManagementComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchemaV1(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *LabourManagementService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.LabourManagementComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.LabourManagementComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchemaV1(listOfComponents)
		}
	}
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	notificationHandler := notification.EmailHandler{
		Logger:                  v.BaseService.Logger,
		EmailNotificationDomain: v.EmailNotificationDomain,
		ComponentManager:        v.ComponentManager,
	}
	// only 1 setting
	err, settingInterface := database.Get(dbConnection, const_util.LabourManagementSettingTable, 1)
	if err != nil {
		v.BaseService.Logger.Error("error getting labour management setting", zap.Error(err))
		os.Exit(0)
	}

	var actualShiftValue = make(map[int][]database.ActualShiftValueCache)
	var machineStatsCache = cache.NewMachineStatsCache()
	actionService := workflow_actions.ActionService{
		Logger:                      v.BaseService.Logger,
		Database:                    dbConnection,
		EmailHandler:                &notificationHandler,
		LabourManagementSettingInfo: database.GetLabourManagementSettingInfo(settingInterface.ObjectInfo),
		ShiftActualValueCache:       actualShiftValue,
		MachineStatsCache:           machineStatsCache,
	}
	var jobService = jobs.JobService{
		Database:                dbConnection,
		Logger:                  v.BaseService.Logger,
		PoolingInterval:         v.MachineStatsPoolingTimer,
		EscalationEmailTemplate: v.EscalationEmailTemplate,
		MachineStatsCache:       machineStatsCache,
	}
	jobService.Init()
	go jobService.LoadMachineStats()
	go jobService.StopShiftAuto()
	v.WorkflowActionHandler = &actionService
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *LabourManagementService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	general := routerEngine.Group("/project/:projectId/labour_management")
	general.POST("/loadFile", middlewares.TokenAuthMiddleware(), v.loadFile)
	general.GET("/overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.summaryResponse)

	general.POST("/action/:actionName", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.handleModulePOSTAction)
	generalComponents := routerEngine.Group("/project/:projectId/labour_management/component/:componentName")

	// table component  requests
	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getObjects)

	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getNewRecord)

	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.updateResource)
	generalComponents.PUT("/record/:recordId/remove_internal_array_reference", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.removeInternalArrayReference)
	generalComponents.POST("/records/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getGroupBy)
	generalComponents.POST("/card_view/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardViewGroupBy)
	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.createNewResource)
	generalComponents.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.importObjects)
	generalComponents.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getTableImportSchema)
	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSearchResults)
	generalComponents.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getComponentRecordTrails)
	generalComponents.GET("/record/:recordId/action/delete_validation", v.deleteValidation)
	generalComponents.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.exportObjectsWithQueryResults)
	generalComponents.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getExportSchema)
	generalComponents.POST("/refresh_export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.refreshExportSchema)
	// action handler
	generalComponents.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordPOSTActionHandler)

	generalComponents.POST("/action/:actionName", middlewares.TokenAuthMiddleware(), v.handleComponentAction)

	generalComponents.GET("/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordGetActionHandler)

}
