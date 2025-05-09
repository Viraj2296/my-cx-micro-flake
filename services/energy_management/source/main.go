package source

import (
	"cx-micro-flake/services/energy_management/source/module_config"
	"cx-micro-flake/services/energy_management/source/service_manager"
	"os"

	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/config"
	"go.cerex.io/transcendflow/service"
	time_series "go.cerex.io/transcendflow/time-series"
	"go.cerex.io/transcendflow/util/color"
	"go.uber.org/zap"
)

func Bootstrap(router *gin.Engine, configPath string) error {
	baseService := service.New()
	err := baseService.Init(configPath)
	if err != nil {
		color.Red("Failed to initialize service machine-downtime service due to error [%s]", err.Error())
		os.Exit(0)
	}

	applicationConfig := config.GetApplicationConfig(configPath)
	conf := module_config.Configuration{}
	if err := applicationConfig.Unmarshal(&conf); err != nil {
		baseService.Logger.Error("invalid service configuration", zap.String("error", err.Error()))
		os.Exit(0)
	}

	DBManager, err := time_series.NewRealtimeDBManager(conf.InfluxConfig, baseService.Logger)
	if err != nil {
		baseService.Logger.Error("error creating influx db manager", zap.String("error", err.Error()))
		os.Exit(1)
	}

	srv := service_manager.Service{
		BaseService:       baseService,
		Config:            conf,
		RealtimeDBManager: DBManager,
	}
	srv.InitService(router)
	baseService.Logger.Info("creating downtime service done")
	return nil
}
