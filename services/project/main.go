package project

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/project/handler"
	"fmt"
	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()

	err := baseService.Init("../services/project/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "project", err.Error())
	}
	baseService.Logger.Info("creating project service")

	projectService := handler.ProjectService{
		BaseService: baseService,
	}
	projectService.InitRouter(router)

	var projectInterface common.ProjectInterface
	projectInterface = &projectService
	service := component.ServiceConfig{
		ServiceType:      "project",
		ServiceName:      "project",
		ServiceInterface: projectInterface,
	}

	common.RegisterService(&service)

	return nil
}
