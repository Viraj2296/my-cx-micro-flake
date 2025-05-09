package facility

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/facility/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/facility/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to facility service [%s] due to error [%s]", "error", err.Error())
	}
	baseService.Logger.Info("creating facility service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/facility/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))

		os.Exit(1)
	}
	var emailNotificationDomain string
	if err := config.Get("emailNotificationDomain").Scan(&emailNotificationDomain); err != nil {
		baseService.Logger.Error("unable to read the email notification domain", zap.String("error", err.Error()))

		os.Exit(1)
	}

	facilityService := handler.FacilityService{
		BaseService:             baseService,
		ComponentContentConfig:  contentConfig,
		EmailNotificationDomain: emailNotificationDomain,
	}
	facilityService.InitRouter(router)

	var facilityServiceInterface common.FacilityServiceInterface
	facilityServiceInterface = &facilityService

	service := component.ServiceConfig{
		ServiceType:      "facility_service",
		ServiceName:      "facility_service_module",
		ServiceInterface: facilityServiceInterface,
	}

	common.RegisterService(&service)

	return nil
}
