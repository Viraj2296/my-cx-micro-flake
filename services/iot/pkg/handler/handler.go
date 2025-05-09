package handler

import (
	"cx-micro-flake/services/iot/resources"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/auth_util"
	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/config"
	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/service"
)

type ModuleConfig struct {
	AppConfig config.AppConfig `mapstructure:"app"`
}

type IoTService struct {
	BaseService      *service.BaseService
	ModuleConfig     ModuleConfig
	ComponentManager *component_processor.ComponentManager
}

func (v *IoTService) InitService(routerEngine *gin.Engine) {
	var componentDataResource []byte
	var err error
	if v.ModuleConfig.AppConfig.LoadComponentSchema {
		componentDataResource, err = resources.GetComponents()
		if err != nil {
			v.BaseService.Logger.Error("Failed to read resource", logging.String("error", err.Error()))
		}
	}
	componentManager := component_processor.ComponentManager{RouterEngine: routerEngine, BaseService: v.BaseService, ComponentResource: componentDataResource}
	componentManager.InitComponents(v.ModuleConfig.AppConfig)
	componentManager.Start()
	v.ComponentManager = &componentManager
	// now init the rest of the routing
	v.InitRouter(routerEngine)
}

func (v *IoTService) InitRouter(routerEngine *gin.Engine) {
	iotRouter := routerEngine.Group(v.ModuleConfig.AppConfig.ModuleName)
	iotRouter.GET("/overview", auth_util.TokenAuthMiddleware(), v.overview)

}
