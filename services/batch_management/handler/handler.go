package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/batch_management/handler/api_access"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"cx-micro-flake/services/batch_management/pkg/workflow_actions"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type BatchManagementService struct {
	BaseService             *common.BaseService
	ComponentContentConfig  component.UpstreamContentConfig
	ComponentManager        *common.ComponentManager
	EmailNotificationDomain string
	APIService              *api_access.APIService
	ActionService           *workflow_actions.ActionService
}

func (v *BatchManagementService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.BatchManagementComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.BatchManagementComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchemaV1(listOfComponents)
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
func (v *BatchManagementService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.BatchManagementComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(listOfComponents)
		}
	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		err, listOfComponents := database.GetObjects(dbConnection, const_util.BatchManagementComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchemaV1(listOfComponents)
		}
	}
	v.ComponentManager.ComponentContentConfig = v.ComponentContentConfig
	v.ActionService = &workflow_actions.ActionService{
		Logger:           v.BaseService.Logger,
		Database:         v.BaseService.ServiceDatabases[const_util.ProjectID],
		ComponentManager: v.ComponentManager,
	}
}

func (v *BatchManagementService) GenerateMouldBatchLabel(projectId string, resourceId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.BatchManagementMouldTable, resourceId)
	if err == nil {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(generalObject.ObjectInfo, &objectFields)
		objectFields["stopTime"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		var mouldBatchNumber = util.InterfaceToString(objectFields["mouldBatchNumber"])
		err, generatedImage := generateQRCode(mouldBatchNumber)
		if err == nil {
			objectFields["label"] = mouldBatchNumber
			objectFields["labelImage"] = generatedImage
			objectFields["canPrint"] = true
		}

		var updatingFields = make(map[string]interface{})
		serialisedData, _ := json.Marshal(objectFields)
		updatingFields["object_info"] = serialisedData
		database.Update(dbConnection, const_util.BatchManagementMouldTable, resourceId, updatingFields)

	}

	return err

}
func getUserTimezone(userId int) string {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	return authService.GetUserTimezone(userId)
}

func (v *BatchManagementService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	batchManagement := routerEngine.Group("/batch_management")
	batchManagement.GET("/v1/printers", v.APIService.GetPrinterDetails)
	batchManagement.POST("/v1/update_printer_status", v.APIService.UpdatePrinterStatus)
	batchManagement.POST("/v1/update_job_feedback", v.APIService.UpdateJobFeedback)

	batchManagementBase := routerEngine.Group("/project/:projectId/batch_management")
	batchManagementBase.POST("/loadFile", middlewares.TokenAuthMiddleware(), v.loadFile)
	batchManagementBase.GET("/batch_management_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.summaryResponse)

	generalComponents := routerEngine.Group("/project/:projectId/batch_management/component/:componentName")

	// table component  requests
	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getObjects)

	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getNewRecord)

	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.updateResource)
	generalComponents.PUT("/record/:recordId/remove_internal_array_reference", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.removeInternalArrayReference)
	generalComponents.POST("/records/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getGroupBy)
	generalComponents.POST("/card_view/group_by", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getCardViewGroupBy)
	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.createNewResource)
	generalComponents.POST("/import", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.importObjects)
	generalComponents.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getTableImportSchema)
	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getSearchResults)
	generalComponents.GET("/record_messages/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(const_util.ModuleName), v.getComponentRecordTrails)
	generalComponents.GET("/record/:recordId/action/delete_validation", v.deleteValidation)
	// action handler
	generalComponents.POST("/record/:recordId/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordPOSTActionHandler)

	generalComponents.POST("/action/:actionName", middlewares.TokenAuthMiddleware(), v.handleComponentAction)

	generalComponents.GET("/action/:actionName", middlewares.TokenAuthMiddleware(), v.recordGetActionHandler)

}
