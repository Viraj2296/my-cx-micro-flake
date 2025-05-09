package source

import (
	"cx-micro-flake/services/spare_part_inventory/source/module_config"
	"cx-micro-flake/services/spare_part_inventory/source/service_manager"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/config"
	"go.cerex.io/transcendflow/service"
	"go.cerex.io/transcendflow/util/color"
	"go.uber.org/zap"
	"os"
)

func Bootstrap(router *gin.Engine, configPath string) error {
	baseService := service.New()
	err := baseService.Init(configPath)
	if err != nil {
		color.Red("failed to initialize service spare part inventory due to error [%s]", err.Error())
		os.Exit(0)
	}

	applicationConfig := config.GetApplicationConfig(configPath)
	conf := module_config.SparePartInventoryConfig{}
	if err := applicationConfig.Unmarshal(&conf); err != nil {
		baseService.Logger.Error("invalid service configuration", zap.String("error", err.Error()))
		os.Exit(0)
	}

	moduleService := service_manager.ModuleService{
		BaseService: baseService,
		Config:      conf,
	}
	moduleService.InitService(router)
	baseService.Logger.Info("creating spare part inventory service done")
	return nil
}
