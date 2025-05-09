package content

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/content/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()

	err := baseService.Init("../services/content/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "contenet_service", err.Error())
	}
	baseService.Logger.Info("creating content service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/content/conf/config.json"),
	))

	var contentConfig handler.ContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	contentService := handler.ContentService{
		BaseService:   baseService,
		ContentConfig: contentConfig,
	}

	contentService.InitRouter(router)

	var contentInterface common.ContentInterface
	contentInterface = &contentService
	service := component.ServiceConfig{
		ServiceType:      "content",
		ServiceName:      "content_module",
		ServiceInterface: contentInterface,
	}

	common.RegisterService(&service)

	return nil

}
