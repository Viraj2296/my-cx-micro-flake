package analytics

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/analytics/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

var ProjectID = "906d0fd569404c59956503985b330132"

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/analytics/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "analytics", err.Error())
	}
	baseService.Logger.Info("creating analytics service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/analytics/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	analyticsService := handler.AnalyticsService{
		BaseService:            baseService,
		ComponentContentConfig: contentConfig,
	}

	analyticsService.OnTimer()
	analyticsService.InitRouter(router)

	var analyticsServiceInterface common.AnalyticsServiceInterface
	analyticsServiceInterface = &analyticsService
	service := component.ServiceConfig{
		ServiceType:      "analytics",
		ServiceName:      "analytics_module",
		ServiceInterface: analyticsServiceInterface,
	}

	common.RegisterService(&service)

	return nil
}
