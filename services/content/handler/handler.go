package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ContentConfig struct {
	StorageDirectory            string `json:"storageDirectory"`
	ApplicationStorageDirectory string `json:"applicationStorageDirectory"`
	DefaultPreviewUrl           string `json:"defaultPreviewUrl"`
	DomainUrl                   string `json:"domainUrl"`
}

// ContentService ModuleConfig this should be generic, we can have any number of module based configuration
type ContentService struct {
	BaseService      *common.BaseService
	ContentConfig    ContentConfig
	ComponentManager *common.ComponentManager
	TestLogger       *zap.SugaredLogger
}

func (v *ContentService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ContentComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ContentComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *ContentService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ContentComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, ContentComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *ContentService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	contentGeneral := routerEngine.Group("/project/:projectId/content")
	contentGeneral.GET("/content_overview", v.getOverview)

	generalComponents := routerEngine.Group("/project/:projectId/content/component/:componentName")
	generalComponents.GET("/record/:recordId", v.getRecordFormData)
	generalComponents.PUT("/record/:recordId", v.updateResource)
	generalComponents.POST("/records", v.createNewResource)
	generalComponents.POST("/records/multiple", v.createMultipleNewResource)
	generalComponents.GET("/record/:recordId/action/:actionName", v.handleGetIndividualRecordAction)
	generalComponents.GET("/action/:actionName", middlewares.TokenAuthMiddleware(), v.handleGetAction)
	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), v.getCardView)
	generalComponents.GET("/:recordId/card_view", middlewares.TokenAuthMiddleware(), v.getChildCardView)
	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), v.getObjects)
	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), v.getSearchResults)
	generalComponents.GET("/:recordId/records", middlewares.TokenAuthMiddleware(), v.getChildObjects)
	generalComponents.POST("/action/:actionName", middlewares.TokenAuthMiddleware(), v.handlePOSTAction)
	generalComponents.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.handleIndividualRecordPOSTAction)
	generalComponents.DELETE("/record/:recordId", v.deleteResource)
}
