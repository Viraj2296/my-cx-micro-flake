package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/pkg/util"

	"github.com/gin-gonic/gin"
)

type QAService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
}

func (v *QAService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, QAComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, QAComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

func composeEventByMould(scheduledEventObject map[string]interface{}) map[string]interface{} {
	createEvent := make(map[string]interface{})
	createEvent["iconCls"] = "fa fa-cogs"
	createEvent["eventType"] = "mould_test_request"

	createEvent["draggable"] = true
	createEvent["module"] = "moulds"
	createEvent["componentName"] = "mould_test_request"
	createEvent["startDate"] = util.InterfaceToString(scheduledEventObject["requestTestStartDate"])
	createEvent["endDate"] = util.InterfaceToString(scheduledEventObject["requestTestEndDate"])
	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0

	return createEvent
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *QAService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, QAComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, QAComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
}

func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *QAService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	mouldGeneral := routerEngine.Group("/project/:projectId/qa")
	mouldGeneral.GET("/overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.summaryResponse)

	generalComponents := routerEngine.Group("/project/:projectId/qa/component/:componentName")

	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getObjects)
	generalComponents.POST("/records/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getGroupBy)
	generalComponents.POST("/card_view/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardViewGroupBy)
	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getNewRecord)

	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.updateResource)
	generalComponents.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteValidation)
	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.createNewResource)
	generalComponents.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.importObjects)
	generalComponents.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getTableImportSchema)
	generalComponents.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.exportObjects)
	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getSearchResults)
	// action handler
	generalComponents.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordPOSTActionHandler)

}
