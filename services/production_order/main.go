package production_order

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/production_order/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/production_order/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "production_order", err.Error())
	}
	baseService.Logger.Info("creating production_order service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/production_order/conf/config.json"),
	))

	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	productionOrderService := handler.ProductionOrderService{
		BaseService:            baseService,
		ComponentContentConfig: contentConfig,
	}
	productionOrderService.InitRouter(router)
	var productionOrderInterface common.ProductionOrderInterface
	productionOrderInterface = &productionOrderService
	service := component.ServiceConfig{
		ServiceType:      "production_order",
		ServiceName:      "production_order_module",
		ServiceInterface: productionOrderInterface,
	}

	common.RegisterService(&service)

	return nil
}
