package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type ServiceConfig struct {
	Email                  EmailConfig                     `json:"email"`
	Content                component.UpstreamContentConfig `json:"content"`
	PushNotificationConfig PushNotificationConfig          `json:"pushNotification"`
}
type MailProcessing struct {
	BaseUrlField string `json:"baseUrlField"`
	BaseUrl      string `json:"baseUrl"`
}

type ContentConfig struct {
	Directory string `json:"directory"`
	DomainUrl string `json:"domainUrl"`
}

type NotificationService struct {
	BaseService      *common.BaseService
	ServiceConfig    ServiceConfig
	ComponentManager *common.ComponentManager
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *NotificationService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, NotificationComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, NotificationComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	v.ComponentManager.ComponentContentConfig = v.ServiceConfig.Content
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *NotificationService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()

	generalComponents := routerEngine.Group("/project/:projectId/notification/component/:componentName")

	// table component  requests
	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), v.getObjects)
	generalComponents.GET("/list", middlewares.TokenAuthMiddleware(), v.getNotificationList)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), v.updateResource)
	generalComponents.POST("/action/:actionName", middlewares.TokenAuthMiddleware(), v.handleComponentAction)
}
