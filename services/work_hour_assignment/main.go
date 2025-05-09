package work_hour_assignment

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/services/work_hour_assignment/handler"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// InitStandaloneService init standalone service needed to have gin common router to support API end points.
func InitStandaloneService(router *gin.Engine, projectConfig []common.ProjectDatasourceConfig) error {
	baseService := common.New()
	err := baseService.Init("../services/work_hour_assignment/conf/config.json", projectConfig)
	if err != nil {
		fmt.Printf("failed to initialize technology service [%s] due to error [%s]", "error", err.Error())
	}
	baseService.Logger.Info("creating technology service")

	//we don't get any error here as we already loaded the config without any errors
	config.Load(file.NewSource(
		file.WithPath("../services/work_hour_assignment/conf/config.json"),
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

	whaService := handler.WorkHourAssignmentService{
		BaseService:             baseService,
		ComponentContentConfig:  contentConfig,
		EmailNotificationDomain: emailNotificationDomain,
	}
	whaService.InitRouter(router)

	var workHourAssignmentInterface common.WorkHourAssignmentInterface
	workHourAssignmentInterface = &whaService

	//dbObject := component.GeneralObject{Id: 1}
	//dbConnection := baseService.ServiceDatabases[handler.ProjectID]
	//type WorkHourAssignmentRequest struct {
	//	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	//	ObjectInfo datatypes.JSON `json:"objectInfo"`
	//}
	//err = dbConnection.Debug().Model(&WorkHourAssignmentRequest{}).Limit(100).Find(&dbObject).Error
	//if err != nil {
	//	fmt.Println("error :", err.Error())
	//	os.Exit(0)
	//}
	//fmt.Println("db object: ", string(dbObject.ObjectInfo))
	//os.Exit(0)
	service := component.ServiceConfig{
		ServiceType:      "wha_service",
		ServiceName:      "wha_service_module",
		ServiceInterface: workHourAssignmentInterface,
	}

	common.RegisterService(&service)

	return nil
}
