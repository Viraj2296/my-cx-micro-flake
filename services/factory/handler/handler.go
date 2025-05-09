package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FactoryService struct {
	BaseService            *common.BaseService
	ComponentContentConfig component.UpstreamContentConfig
	ComponentManager       *common.ComponentManager
}

func (v *FactoryService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FactoryComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}
	}
	v.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FactoryComponentTable)
		if err == nil {
			v.ComponentManager.LoadTableSchema(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (v *FactoryService) InitComponents() {

	var totalComponents int
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FactoryComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(*listOfComponents)
		}

	}
	v.ComponentManager = &common.ComponentManager{}
	v.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	for _, dbConnection := range v.BaseService.ServiceDatabases {
		listOfComponents, err := GetObjects(dbConnection, FactoryComponentTable)
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

func (v *FactoryService) GetDepartmentName(departmentId int) string {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	err, generalObject := Get(dbConnection, FactoryDepartmentTable, departmentId)
	if err == nil {
		factoryDepartment := FactoryDepartment{ObjectInfo: generalObject.ObjectInfo}
		return factoryDepartment.getDepartmentInfo().Name
	}

	return ""
}

func (v *FactoryService) GetSiteName(siteId int) string {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	err, generalObject := Get(dbConnection, FactorySiteTable, siteId)
	if err == nil {
		factorySite := FactorySite{ObjectInfo: generalObject.ObjectInfo}
		return factorySite.getFactorySiteInfo().Name
	}

	return ""
}
func (v *FactoryService) GetBuildingInfo() *[]component.GeneralObject {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject, err := GetObjects(dbConnection, FactoryBuildingTable)
	if err != nil {
		return nil
	}

	return generalObject
}

func (v *FactoryService) IsFactoryBuildingExist(recordId int) bool {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	conditionalString := "id = " + strconv.Itoa(recordId)
	fmt.Println("conditionalString", conditionalString)
	generalObject, err := GetConditionalObjects(dbConnection, FactoryBuildingTable, conditionalString)
	if err != nil {
		fmt.Println("E1")
		return false
	}

	if len(*generalObject) > 0 {
		fmt.Println("E2")
		return true
	}
	fmt.Println("E3")
	return false
}

func (v *FactoryService) InitRouter(routerEngine *gin.Engine) {
	v.InitComponents()
	factoryServiceGeneral := routerEngine.Group("/project/:projectId/factory")
	factoryServiceGeneral.POST("/loadFile", v.loadFile)
	factoryServiceGeneral.GET("/factory_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), v.getFactoryOverview)
	factoryService := routerEngine.Group("/project/:projectId/factory/component/:componentName")

	// table component  requests
	factoryService.GET("/records", middlewares.TokenAuthMiddleware(), v.getObjects)
	factoryService.GET("/card_view", middlewares.TokenAuthMiddleware(), v.getCardView)
	factoryService.GET("/new_record", middlewares.TokenAuthMiddleware(), v.getNewRecord)
	factoryService.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), v.getRecordFormData)
	factoryService.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), v.updateResource)
	factoryService.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), v.deleteResource)
	factoryService.POST("/records", middlewares.TokenAuthMiddleware(), v.createNewResource)
	factoryService.POST("/import", middlewares.TokenAuthMiddleware(), v.importObjects)
	factoryService.GET("/table_import_schema", middlewares.TokenAuthMiddleware(), v.getTableImportSchema)
	factoryService.POST("/export", middlewares.TokenAuthMiddleware(), v.exportObjects)
	factoryService.POST("/search", middlewares.TokenAuthMiddleware(), v.getSearchResults)
	factoryService.GET("/export_schema", middlewares.TokenAuthMiddleware(), v.getExportSchema)
	factoryService.PUT("/record/:recordId/remove_internal_array_reference", v.removeInternalArrayReference)

	//get the records
	factoryService.GET("/record_messages/:recordId", v.getComponentRecordTrails)

}
