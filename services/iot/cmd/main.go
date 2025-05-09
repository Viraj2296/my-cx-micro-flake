package cmd

import (
	"cx-micro-flake/services/iot/pkg/handler"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/config"
	"go.cerex.io/transcendflow/service"
	"go.cerex.io/transcendflow/util/color"
	"go.uber.org/zap"
	"os"
)

func Bootstrap(router *gin.Engine) error {
	baseService := service.New()
	err := baseService.Init("../services/iot/conf")
	if err != nil {
		color.Red("Failed to initialize service iot service due to error [%s]", err.Error())
		os.Exit(0)
	}

	applicationConfig := config.GetApplicationConfig("../services/iot/conf")
	conf := handler.ModuleConfig{}
	if err := applicationConfig.Unmarshal(&conf); err != nil {
		baseService.Logger.Error("invalid service configuration", zap.String("error", err.Error()))
		os.Exit(0)
	}

	factoryService := handler.IoTService{
		BaseService:  baseService,
		ModuleConfig: conf,
	}
	factoryService.InitService(router)
	baseService.Logger.Info("creating iot service")
	return nil
}
