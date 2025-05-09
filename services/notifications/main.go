package notifications

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/notifications/handler"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()

	err := baseService.Init("../services/notifications/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "notifications_service", err.Error())
	}
	baseService.Logger.Info("creating notifications service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/notifications/conf/config.json"),
	))

	var serviceConfig handler.ServiceConfig
	if err := config.Get().Scan(&serviceConfig); err != nil {
		baseService.Logger.Error("unable to read the notifications config", zap.String("error", err.Error()))
		os.Exit(1)
	}

	notificationManager := handler.NotificationManager{
		ServiceDatabase:        baseService.ServiceDatabases["906d0fd569404c59956503985b330132"],
		EmailConfig:            serviceConfig.Email,
		Logger:                 baseService.Logger,
		PushNotificationConfig: serviceConfig.PushNotificationConfig,
	}

	notificationManager.Init()
	if serviceConfig.Email.PoolingEnabled {
		go notificationManager.PoolNotification()
	}
	if serviceConfig.PushNotificationConfig.PoolingEnabled {
		go notificationManager.PoolPushNotification()
	}
	notificationService := handler.NotificationService{
		BaseService:   baseService,
		ServiceConfig: serviceConfig,
	}
	notificationService.InitRouter(router)

	var notificationInterface common.NotificationInterface
	notificationInterface = &notificationService
	service := component.ServiceConfig{
		ServiceType:      "notification",
		ServiceName:      "notification_module",
		ServiceInterface: notificationInterface,
	}

	// let's send some test email
	notificationManager.SendTestEmail()
	common.RegisterService(&service)

	return nil

}
