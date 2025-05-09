package scheduler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/scheduler/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/scheduler/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "scheduler", err.Error())
	}
	baseService.Logger.Info("creating scheduler service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/scheduler/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	schedulerService := handler.SchedulerService{
		BaseService:            baseService,
		ComponentContentConfig: contentConfig,
	}
	schedulerService.InitRouter(router)

	var schedulingInterface common.SchedulingInterface
	schedulingInterface = &schedulerService

	service := component.ServiceConfig{
		ServiceType:      "scheduler",
		ServiceName:      "scheduler_module",
		ServiceInterface: schedulingInterface,
	}

	common.RegisterService(&service)

	return nil
}
