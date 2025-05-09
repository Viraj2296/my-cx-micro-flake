package auth

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/services/auth/handler"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/auth/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize service [%s] due to error [%s]", "auth", err.Error())
	}
	baseService.Logger.Info("creating auth service")

	var emailNotificationDomain string
	if err := config.Get("emailNotificationDomain").Scan(&emailNotificationDomain); err != nil {
		baseService.Logger.Error("unable to read the email notification domain", zap.String("error", err.Error()))
		os.Exit(1)
	}

	var profileGroupId int
	if err := config.Get("profileGroupId").Scan(&profileGroupId); err != nil {
		baseService.Logger.Error("unable to read the profileGroupId", zap.String("error", err.Error()))
		os.Exit(1)
	}
	var supervisorJobRoleId int
	if err := config.Get("supervisorJobRoleId").Scan(&supervisorJobRoleId); err != nil {
		baseService.Logger.Error("unable to read the supervisorJobRoleId", zap.String("error", err.Error()))
		os.Exit(1)
	}
	var defaultRefreshTokenExpiry float64
	if err := config.Get("defaultRefreshTokenExpiry").Scan(&defaultRefreshTokenExpiry); err != nil {
		baseService.Logger.Error("unable to read the default Refresh Token Expiry", zap.String("error", err.Error()))
		os.Exit(1)
	}
	var defaultTokenExpiry float64
	if err := config.Get("defaultTokenExpiry").Scan(&defaultTokenExpiry); err != nil {
		baseService.Logger.Error("unable to read the default Token Expiry", zap.String("error", err.Error()))
		os.Exit(1)
	}
	var contentConfig component.UpstreamContentConfig
	if err := config.Get("content").Scan(&contentConfig); err != nil {
		baseService.Logger.Error("unable to read the content config", zap.String("error", err.Error()))

		os.Exit(1)
	}
	var reports []handler.Report
	if err := config.Get("report").Scan(&reports); err != nil {
		baseService.Logger.Error("unable to read the report config", zap.String("error", err.Error()))
		os.Exit(1)
	}
	baseService.Logger.Info("service initialization", zap.Any("ProfileGroupId", profileGroupId))
	authService := handler.AuthService{
		BaseService:               baseService,
		EmailNotificationDomain:   emailNotificationDomain,
		Report:                    reports,
		ProfileGroupId:            profileGroupId,
		SupervisorJobRoleId:       supervisorJobRoleId,
		DefaultRefreshTokenExpiry: defaultRefreshTokenExpiry,
		DefaultTokenExpiry:        defaultTokenExpiry,
		ComponentContentConfig:    contentConfig,
	}

	authService.InitRouter(router)

	var authServiceInterface common.AuthInterface
	authServiceInterface = &authService
	service := component.ServiceConfig{
		ServiceType:      "auth",
		ServiceName:      "general_auth",
		ServiceInterface: authServiceInterface,
	}

	common.RegisterService(&service)

	return nil

}
