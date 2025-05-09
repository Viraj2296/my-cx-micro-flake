package service_manager

import (
	"cx-micro-flake/services/spare_part_inventory/source/actions"
	"cx-micro-flake/services/spare_part_inventory/source/jobs"
	"cx-micro-flake/services/spare_part_inventory/source/module_config"
	"cx-micro-flake/services/spare_part_inventory/source/repository"
	"cx-micro-flake/services/spare_part_inventory/source/service_manager/http"
	"cx-micro-flake/services/spare_part_inventory/source/service_manager/rpc"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/service"
)

type ModuleService struct {
	BaseService      *service.BaseService
	Config           module_config.SparePartInventoryConfig
	ComponentManager *component_processor.ComponentManager
}

func (v *ModuleService) InitService(routerEngine *gin.Engine) {

	repo := repository.NewSparePartInventoryRepository(v.BaseService.ServiceDatabase, v.BaseService.Logger)

	componentManager := component_processor.ComponentManager{RouterEngine: routerEngine, BaseService: v.BaseService}
	actionService := actions.NewActions(v.BaseService.Logger, v.BaseService.ServiceDatabase, repo, &componentManager)
	componentManager.InitComponents(v.Config.AppConfig)
	componentManager.SetActionHandler(actionService)
	componentManager.Start()
	v.ComponentManager = &componentManager

	jobService := jobs.NewJobs(v.BaseService.Logger, v.BaseService.ServiceDatabase, v.Config.SparePartConfig.JobServiceConfig, repo, v.ComponentManager)
	jobService.Init()

	httpService := http.NewHTTPService(v.BaseService.Logger, repo, &componentManager)
	httpService.Init()
	httpService.InitRouter()
	v.BaseService.Logger.Info("started the HTTP service successfully")

	rpcService := rpc.NewRPCService(v.BaseService.Logger, repo, &componentManager)
	rpcService.Init()
	v.BaseService.Logger.Info("started the RPC service successfully")

	// go v.JobService.SendEmailEscalationJob()
	//v.ComponentManager.ClientService.RegisterFunction("xx", rpcService.xx())
	go jobService.GenerateNotificationMessage()
}
