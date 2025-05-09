package source

import (
	"cx-micro-flake/services/machine_downtime/source/module_config"
	"cx-micro-flake/services/machine_downtime/source/service_manager"
	"go.cerex.io/transcendflow/config"
	"go.cerex.io/transcendflow/service"
	"go.cerex.io/transcendflow/util/color"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

func Bootstrap(router *gin.Engine, configPath string) error {
	baseService := service.New()
	err := baseService.Init(configPath)
	if err != nil {
		color.Red("Failed to initialize service machine-downtime service due to error [%s]", err.Error())
		os.Exit(0)
	}

	applicationConfig := config.GetApplicationConfig(configPath)
	conf := module_config.MachineDowntimeConfig{}
	if err := applicationConfig.Unmarshal(&conf); err != nil {
		baseService.Logger.Error("invalid service configuration", zap.String("error", err.Error()))
		os.Exit(0)
	}

	downtimeService := service_manager.MachineDowntimeService{
		BaseService: baseService,
		Config:      conf,
	}
	downtimeService.InitService(router)
	baseService.Logger.Info("creating downtime service done")
	return nil
}
