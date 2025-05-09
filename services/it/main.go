package it

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/it/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/it/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize technology service [%s] due to error [%s]", "error", err.Error())
	}
	baseService.Logger.Info("creating technology service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/it/conf/config.json"),
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

	itService := handler.ITService{
		BaseService:             baseService,
		ComponentContentConfig:  contentConfig,
		EmailNotificationDomain: emailNotificationDomain,
	}
	itService.InitRouter(router)

	var itServiceInterface common.ITServiceInterface
	itServiceInterface = &itService

	service := component.ServiceConfig{
		ServiceType:      "it_service",
		ServiceName:      "it_service_module",
		ServiceInterface: itServiceInterface,
	}

	common.RegisterService(&service)

	return nil
}
