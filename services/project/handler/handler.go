package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProjectService struct {
	BaseService            *common.BaseService
	ComponentManager       *common.ComponentManager
	ComponentContentConfig component.UpstreamContentConfig
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *ProjectService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ProjectComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ProjectComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	// init the component config
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *ProjectService) InitRouter(routerEngine *gin.Engine) {

	v.InitComponents()
	projectGeneral := routerEngine.Group("/project")
	projectGeneral.POST("/loadFile", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.loadFile)
	projectGeneral.GET("/overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSummaryResponse)
	project := routerEngine.Group("/project/component/:componentName")

	// table component  requests
	project.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getObjects)
	project.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardView)
	project.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNewRecord)
	project.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getRecordFormData)
	project.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateResource)
	project.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteResource)
	project.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.createNewResource)
	project.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.importObjects)
	project.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTableImportSchema)
	project.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.exportObjects)
	project.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSearchResults)
	project.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getExportSchema)
	project.GET("/record/:recordId/action/delete_validation", v.deleteValidation)

	project.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getComponentRecordTrails)
	//get the records
}

func (v *ProjectService) GetProjectDatasourceInfo() []common.ProjectDatasourceConfig {

	listOfProjects, err := GetObjects(v.BaseService.ReferenceDatabase, ProjectTable)
	var listOfProjectDatasourceConfig []common.ProjectDatasourceConfig
	v.BaseService.Logger.Info("loading all the projects", zap.Any("size", *(listOfProjects)))
	if err == nil {
		for _, projectInterface := range *listOfProjects {
			project := Project{ObjectInfo: projectInterface.ObjectInfo}
			projectDatasource := project.getProjectInfo().ProjectDatasource
			projectDatasourceConfig := common.ProjectDatasourceConfig{ProjectId: project.getProjectInfo().ProjectReferenceId, DatasourceConfig: projectDatasource}

			listOfProjectDatasourceConfig = append(listOfProjectDatasourceConfig, projectDatasourceConfig)

			v.BaseService.Logger.Info("loading all the projects", zap.Any("size", *(listOfProjects)))
		}
	} else {
		v.BaseService.Logger.Error("error getting projects", zap.String("error", err.Error()))
	}
	return listOfProjectDatasourceConfig

}
