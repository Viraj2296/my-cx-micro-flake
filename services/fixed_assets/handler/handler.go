package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type FixedAssetService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	QRCodeDomainUrl         string
	EmailNotificationDomain string
}

func (v *FixedAssetService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FixedAssetsComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FixedAssetsComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *FixedAssetService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FixedAssetsComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FixedAssetsComponentTable)
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

func (v *FixedAssetService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	fixedAssetsGeneral := routerEngine.Group("/project/:projectId/fixed_assets")
	fixedAssetsGeneral.POST("/loadFile", v.loadFile)
	fixedAssetsGeneral.GET("/fixed_asset_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getOverview)
	fixedAsstes := routerEngine.Group("/project/:projectId/fixed_assets/component/:componentName")

	// table component  requests
	fixedAsstes.GET("/records", middlewares.TokenAuthMiddleware(), v.getObjects)
	fixedAsstes.GET("/card_view", middlewares.TokenAuthMiddleware(), v.getCardView)
	fixedAsstes.GET("/new_record", middlewares.TokenAuthMiddleware(), v.getNewRecord)

	fixedAsstes.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), v.getRecordFormData)
	fixedAsstes.GET("/record/:recordId/dynamic_fields", middlewares.TokenAuthMiddleware(), v.getDynamicFields)

	fixedAsstes.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), v.updateResource)
	fixedAsstes.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), v.deleteResource)
	fixedAsstes.POST("/records", middlewares.TokenAuthMiddleware(), v.createNewResource)
	fixedAsstes.POST("/import", middlewares.TokenAuthMiddleware(), v.importObjects)
	fixedAsstes.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), v.getTableImportSchema)
	fixedAsstes.POST("/export", middlewares.TokenAuthMiddleware(), v.exportObjects)
	fixedAsstes.POST("/search", middlewares.TokenAuthMiddleware(), v.getSearchResults)
	fixedAsstes.GET("/export_schema", middlewares.TokenAuthMiddleware(), v.getExportSchema)
	fixedAsstes.PUT("/record/:recordId/remove_internal_array_reference", v.removeInternalArrayReference)
	fixedAsstes.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordPOSTActionHandler)
	//get the records
	fixedAsstes.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), v.getComponentRecordTrails)

}
