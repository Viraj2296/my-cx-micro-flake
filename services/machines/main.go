package machines

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/machines/handler"
	"cx-micro-flake/services/machines/handler/jobs"
	"fmt"
	"os"

	"go.cerex.io/transcendflow/logging"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/machines/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "machine", err.Error())
	}
	baseService.Logger.Info("creating machine service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/machines/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", logging.String("error", err.Error()))
		os.Exit(1)
	}
	var calConfig handler.CalculationConfig
	if err := config.Get("calculation").Scan(&calConfig); err != nil {
		baseService.Logger.Error("unable to read the calculation config", logging.String("error", err.Error()))
		os.Exit(1)
	}

	var emailConfig handler.EmailConfig
	if err := config.Get("email").Scan(&emailConfig); err != nil {
		baseService.Logger.Error("unable to read the email config", logging.String("error", err.Error()))
		os.Exit(1)
	}

	var groupPermissionConfig int
	if err := config.Get("preventiveWorkOrderAllowedGroup").Scan(&groupPermissionConfig); err != nil {
		baseService.Logger.Error("unable to read the groupPermissionConfig config", logging.String("error", err.Error()))
		os.Exit(1)
	}

	var assemblyConfig []handler.AssemblyData
	if err := config.Get("assemblyData").Scan(&assemblyConfig); err != nil {
		baseService.Logger.Error("unable to read the assembly config", logging.String("error", err.Error()))
		os.Exit(1)
	}

	var jobsConfig jobs.JobsConfig
	if err := config.Get("jobs").Scan(&jobsConfig); err != nil {
		baseService.Logger.Error("unable to read the jobs config", logging.String("error", err.Error()))
		os.Exit(1)
	}

	machineService := handler.MachineService{
		BaseService:                  baseService,
		ComponentContentConfig:       contentConfig,
		CalculationConfig:            &calConfig,
		EmailConfig:                  &emailConfig,
		GroupPermissionConfig:        groupPermissionConfig,
		AssemblyMachineConfiguration: assemblyConfig,
		JobsConfig:                   jobsConfig,
	}

	machineService.OnTimer()
	machineService.InitRouter(router)
	var machineServiceInterface common.MachineInterface
	machineServiceInterface = &machineService
	service := component.ServiceConfig{
		ServiceType:      "machines",
		ServiceName:      "machines_module",
		ServiceInterface: machineServiceInterface,
	}

	common.RegisterService(&service)

	return nil
}
