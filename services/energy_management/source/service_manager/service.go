package service_manager

import (
	"cx-micro-flake/services/energy_management/source/actions"
	"cx-micro-flake/services/energy_management/source/jobs"
	"cx-micro-flake/services/energy_management/source/module_config"
	"cx-micro-flake/services/energy_management/source/repository"

	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/service"
	time_series "go.cerex.io/transcendflow/time-series"
	"go.uber.org/zap"
)

type Service struct {
	BaseService       *service.BaseService
	Config            module_config.Configuration
	ComponentManager  *component_processor.ComponentManager
	Actions           *actions.Actions
	Jobs              *jobs.Jobs // Add JobService field here
	RealtimeDBManager *time_series.RealtimeDBManager
}

func (v *Service) InitService(routerEngine *gin.Engine) {
	// check whether we need to read from file
	repo := repository.NewRepository(v.BaseService.Logger, v.BaseService.ServiceDatabase)
	actionService := actions.NewActions(v.BaseService.Logger, v.BaseService.ServiceDatabase, repo, v.RealtimeDBManager, v.Config.EnergyStartDate)
	err := v.RealtimeDBManager.Connect()
	if err != nil {
		v.BaseService.Logger.Error("connecting to db manager failed", zap.Error(err))
	}
	jobService := jobs.NewJobs(v.BaseService.Logger, v.BaseService.ServiceDatabase, repo)
	jobService.Init()

	componentManager := component_processor.ComponentManager{RouterEngine: routerEngine, BaseService: v.BaseService}
	componentManager.InitComponents(v.Config.AppConfig)
	componentManager.SetActionHandler(actionService)
	componentManager.Start()
	v.ComponentManager = &componentManager
	v.Actions = actionService
	v.Jobs = jobService

	v.InitRouter()
}
