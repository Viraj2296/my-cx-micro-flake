package service_manager

import (
	"cx-micro-flake/services/machine_downtime/source/actions"
	"cx-micro-flake/services/machine_downtime/source/jobs"
	"cx-micro-flake/services/machine_downtime/source/module_config"
	"cx-micro-flake/services/machine_downtime/source/repository"

	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/service"
)

type MachineDowntimeService struct {
	BaseService      *service.BaseService
	Config           module_config.MachineDowntimeConfig
	ComponentManager *component_processor.ComponentManager
	Actions          *actions.Actions
	Jobs             *jobs.Jobs // Add JobService field here
}

func (v *MachineDowntimeService) InitService(routerEngine *gin.Engine) {
	// check whether we need to read from file
	repo := repository.NewRepository(v.BaseService.Logger, v.BaseService.ServiceDatabase)
	actionService := actions.NewActions(v.BaseService.Logger, v.BaseService.ServiceDatabase, repo)

	jobService := jobs.NewJobs(v.BaseService.Logger, v.BaseService.ServiceDatabase, v.Config.DowntimeConfig, repo)
	jobService.Init()

	componentManager := component_processor.ComponentManager{RouterEngine: routerEngine, BaseService: v.BaseService}
	componentManager.InitComponents(v.Config.AppConfig)
	componentManager.SetActionHandler(actionService)
	componentManager.Start()
	v.ComponentManager = &componentManager
	v.Actions = actionService
	v.Jobs = jobService

	// Start polling function in a separate goroutine
	go jobService.GenerateNotificationMessage()
	if v.Config.DowntimeConfig.JobServiceConfig.EnableEascalation {
		// Start escalation job in a separate goroutine
		v.BaseService.Logger.Info("escalation job is enabled")
		go jobService.SendEmailEscalationJob()
	} else {
		v.BaseService.Logger.Info("escalation job is disabled")
	}
	// now init the rest of the routing
	v.InitRouter()
}
