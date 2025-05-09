package labour_management

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/labour_management/handler"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/labour_management/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "machine", err.Error())
	}
	baseService.Logger.Info("creating labour_management service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/labour_management/conf/config.json"),
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

	var machineStatsPoolingTimer int
	if err := config.Get("machineStatsPoolingTimer").Scan(&machineStatsPoolingTimer); err != nil {
		baseService.Logger.Error("unable to read the machineStatsPoolingTimer", zap.String("error", err.Error()))
		os.Exit(1)
	}
	var escalationEmailTemplate string
	if err := config.Get("escalationEmailTemplate").Scan(&escalationEmailTemplate); err != nil {
		baseService.Logger.Error("unable to read the escalationEmailTemplate", zap.String("error", err.Error()))
	}
	labourManagementService := handler.LabourManagementService{
		BaseService:              baseService,
		ComponentContentConfig:   contentConfig,
		EmailNotificationDomain:  emailNotificationDomain,
		MachineStatsPoolingTimer: machineStatsPoolingTimer,
		EscalationEmailTemplate:  escalationEmailTemplate,
	}
	labourManagementService.InitRouter(router)

	var labourManagementInterface common.LabourManagementInterface
	labourManagementInterface = &labourManagementService

	service := component.ServiceConfig{
		ServiceType:      "labour_management",
		ServiceName:      "labour_management_module",
		ServiceInterface: labourManagementInterface,
	}

	common.RegisterService(&service)

	return nil
}
