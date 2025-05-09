package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/machines/handler/jobs"
	"cx-micro-flake/services/machines/handler/views"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type MachineService struct {
	BaseService                  *common.BaseService
	ComponentContentConfig       component.UpstreamContentConfig
	ComponentManager             *common.ComponentManager
	CalculationConfig            *CalculationConfig
	EmailConfig                  *EmailConfig
	GroupPermissionConfig        int
	AssemblyMachineConfiguration []AssemblyData
	ViewManager                  *views.ViewManager
	JobService                   *jobs.JobService
	MouldShotCountCache          map[string]int
	JobsConfig                   jobs.JobsConfig
}

type EmailConfig struct {
	FilePath string `json:"filePath"`
}

type AssemblyData struct {
	MessageFlag string   `json:"messageFlag"`
	Area        []string `json:"area"`
}

type AssemblyConfig struct {
	AssemblyData []AssemblyData `json:"assemblyData"`
}

type MachineConfig struct {
	EnableCalculation   bool   `json:"enableCalculation"`
	CalculationInterval string `json:"calculationInterval"`
	MachineType         string `json:"machineType"`
}

type CalculationConfig struct {
	MachineConfiguration []MachineConfig `json:"machineCalculation"`
}

func (v *MachineService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MachineComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MachineComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *MachineService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MachineComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, MachineComponentTable)
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
	v.MouldShotCountCache = make(map[string]int)
	jobService := jobs.JobService{Logger: v.BaseService.Logger, Database: v.BaseService.ServiceDatabases[ProjectID], JobConfig: v.JobsConfig}
	err := jobService.Init()
	if err != nil {
		v.BaseService.Logger.Error("error init the job server, system won't generate help signal", zap.Error(err))
	} else {
		v.BaseService.Logger.Info("starting generate help signal history")
		go jobService.GenerateHelpSignalHistory()
	}
}

func (v *MachineService) createSystemNotification(projectId, header, description string, recordId int) error {
	systemNotification := common.SystemNotification{}
	systemNotification.Name = header
	systemNotification.ColorCode = "#14F44E"
	systemNotification.IconCls = "icon-park-outline:transaction-order"
	systemNotification.RecordId = recordId
	systemNotification.RouteLinkComponent = "timeline"
	systemNotification.Component = "Production Order"
	systemNotification.Description = description
	systemNotification.GeneratedTime = util.GetCurrentTime(time.RFC822)
	rawSystemNotification, _ := json.Marshal(systemNotification)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	err := notificationService.CreateSystemNotification(projectId, rawSystemNotification)
	return err
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *MachineService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	//ms.generatePermissions()
	machineGeneral := routerEngine.Group("/project/:projectId/machines")
	machineGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.loadFile)
	machineGeneral.GET("/machine_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getOverview)
	machineGeneral.GET("/machine_stats_summary", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getMachineStatsSummary)
	machineGeneral.GET("/machine_summary", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getMachineSummary)

	//machineGeneral.GET("/:machineId/dashboard", middlewares.TokenAuthMiddleware(), ms.getMachineDashboard)
	machineGeneral.GET("/dashboard/list_of_tv_machines", middlewares.TokenAuthMiddleware(), v.getIntialMachineDashboardInfo)
	machineGeneral.GET("/record/:recordId/moulding_machine_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getMouldingMachineOverview)
	//Hmi info
	//machineGeneral.GET("/:machineId/hmi_info", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ms.getHmiInfo)
	//machineGeneral.GET("/:machineId/hmi_info/action/previous/:eventId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ms.getPreviousHmi)
	//machineGeneral.GET("/:machineId/hmi_info/action/next/:eventId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), ms.getNextHmi)
	machineGeneral.GET("/hmi_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getHMIView)

	// Filter assigned machines according to user id
	machineGeneral.GET("/get_assigned_machines", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getAssignedMachines)

	// Assembly hmi info
	machineGeneral.GET("/assembly_hmi_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getAssemblyHMIView)

	machineGeneral.GET("/tooling_hmi_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getToolingHMIView)

	generalComponents := routerEngine.Group("/project/:projectId/:moduleName/component/:componentName")

	// table component  requests
	generalComponents.GET("/record/:recordId/dashboard", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getMachineDashboard)
	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getObjects)
	generalComponents.POST("/records/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getGroupBy)
	generalComponents.POST("/card_view/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardViewGroupBy)
	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNewRecord)
	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateResource)
	generalComponents.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteValidation)
	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.createNewResource)
	generalComponents.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.importObjects)
	generalComponents.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTableImportSchema)
	generalComponents.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.exportObjects)

	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSearchResults)
	generalComponents.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getExportSchema)
	generalComponents.POST("/refresh_export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.refreshExportSchema)
	generalComponents.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(ModuleName), v.removeInternalArrayReference)

	generalComponents.GET("/action/production_dashboard_display_setting", middlewares.TokenAuthMiddleware(), v.getProductionDashboardDisplaySetting)
	//get the records
	generalComponents.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), v.getComponentRecordTrails)

	generalComponents.POST("/reset_machine", middlewares.TokenAuthMiddleware(), v.addingManualAssemblyResetMessage)

	// Get hmi for moulding and assembly
	generalComponents.GET("/record/:recordId/hmi_info", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getHmiInfo)
	generalComponents.GET("/record/:recordId/hmi_info/action/previous/:eventId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getPreviousHmi)
	generalComponents.GET("/record/:recordId/hmi_info/action/next/:eventId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNextHmi)

	generalComponents.POST("/records/filter", middlewares.TokenAuthMiddleware(), v.getFilterObjects)
	generalComponents.GET("/filters", middlewares.TokenAuthMiddleware(), v.getCacheFilters)
	routerEngine.Routes()

}
