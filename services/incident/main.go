package incident

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/incident/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/incident/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to incident management service [%s] due to error [%s]", "error", err.Error())
	}
	baseService.Logger.Info("creating technology service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/incident/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	incidentService := handler.IncidentService{
		BaseService:            baseService,
		ComponentContentConfig: contentConfig,
	}
	incidentService.InitRouter(router)

	var incidentServiceInterface common.IncidentServiceInterface
	incidentServiceInterface = &incidentService

	service := component.ServiceConfig{
		ServiceType:      "incident_service",
		ServiceName:      "incident_service_module",
		ServiceInterface: incidentServiceInterface,
	}

	common.RegisterService(&service)

	return nil
}
