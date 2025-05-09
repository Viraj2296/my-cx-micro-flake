package maintenance

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/maintenance/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/maintenance/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "maintenance", err.Error())
	}
	baseService.Logger.Info("creating machine service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/maintenance/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	var emailNotificationDomain string
	if err := config.Get("emailNotificationDomain").Scan(&emailNotificationDomain); err != nil {
		baseService.Logger.Error("unable to read the email notification domain", zap.String("error", err.Error()))
		os.Exit(1)
	}

	var toolingSupervisorGroup int
	if err := config.Get("toolingSupervisorGroup").Scan(&toolingSupervisorGroup); err != nil {
		baseService.Logger.Error("unable to read the toolingSupervisorGroup", zap.String("error", err.Error()))
		os.Exit(1)
	}

	maintenanceService := handler.MaintenanceService{
		BaseService:             baseService,
		ComponentContentConfig:  contentConfig,
		EmailNotificationDomain: emailNotificationDomain,
		ToolingSupervisorGroup:  toolingSupervisorGroup,
	}
	maintenanceService.OnTimer()
	maintenanceService.InitRouter(router)

	var maintenanceServiceInterface common.MaintenanceInterface
	maintenanceServiceInterface = &maintenanceService

	service := component.ServiceConfig{
		ServiceType:      "maintenance",
		ServiceName:      "maintenance_module",
		ServiceInterface: maintenanceServiceInterface,
	}

	common.RegisterService(&service)

	return nil
}
