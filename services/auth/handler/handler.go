package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/service_zmq"
)

const (
	loginError = "An authentication error occurred. Try updating your credentials and logging in again."
)

type AuthService struct {
	BaseService               *common.BaseService
	PermissionCache           map[int][]*ComponentResourceInfo
	ModuleCache               map[string]int
	ComponentContentConfig    component.UpstreamContentConfig
	ComponentManager          *common.ComponentManager
	EmailNotificationDomain   string
	Report                    []Report
	ProfileGroupId            int
	SupervisorJobRoleId       int
	DefaultRefreshTokenExpiry float64
	DefaultTokenExpiry        float64
	ClientService             *service_zmq.ClientService
}

type Report struct {
	MenuId      string `json:"menuId"`
	DashboardId int    `json:"dashboardId"`
}

func (as *AuthService) LoadInitComponents() {

	var totalComponents int
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		err, listOfComponents := GetObjects(dbConnection, IAMComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(listOfComponents)
		}
	}
	as.ComponentManager.InitComponentManager(totalComponents)
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		err, listOfComponents := GetObjects(dbConnection, IAMComponentTable)
		if err == nil {
			as.ComponentManager.LoadTableSchemaV1(listOfComponents)
		}
	}
}

// InitComponents if we are writing this as component, we should init the component manager
func (as *AuthService) InitComponents() {
	var totalComponents int
	for _, dbConnection := range as.BaseService.ServiceDatabases {
		err, listOfComponents := GetObjects(dbConnection, IAMComponentTable)
		if err == nil {
			totalComponents = totalComponents + len(listOfComponents)
		}
	}
	as.ComponentManager = &common.ComponentManager{}
	as.ComponentManager.InitComponentManager(totalComponents)
	// next init table schema and linked fields, should we handled the errors
	// TODO, better we handle the errors in the component level
	dbConnection := as.BaseService.ReferenceDatabase
	err, listOfComponents := GetObjects(dbConnection, IAMComponentTable)
	if err == nil {
		as.ComponentManager.LoadTableSchemaV1(listOfComponents)
	}
	as.ModuleCache = make(map[string]int, 0)
	as.PermissionCache = make(map[int][]*ComponentResourceInfo, 100)
	as.ComponentManager.ComponentContentConfig = as.ComponentContentConfig

}

func (as *AuthService) InitRouter(routerEngine *gin.Engine) {
	as.InitComponents()
	as.loadUserAccess()
	// common routes without project ids
	routerEngine.POST("/login", as.login)
	routerEngine.GET("/health", as.healthHandler)
	routerEngine.GET("/renew_token", middlewares.TokenAuthMiddleware(), as.renewToken)
	routerEngine.POST("/forget_password", as.forgetPassword)
	routerEngine.POST("/validate_reset_token", as.validateResetToken)
	routerEngine.POST("/new_password", middlewares.TokenAuthMiddleware(), as.createNewPassword)

	// user component actions

	routerEngine.POST("/logout", middlewares.TokenAuthMiddleware(), as.userLogout)
	routerEngine.GET("/profile", middlewares.TokenAuthMiddleware(), as.getUserProfile)
	routerEngine.PUT("/update_profile", middlewares.TokenAuthMiddleware(), as.updateProfile)
	routerEngine.POST("/update_password", middlewares.TokenAuthMiddleware(), as.updatePassword)

	// projects based routes
	generalComponents := routerEngine.Group("/project/:projectId/iam/component/:componentName")

	generalComponents.GET("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getObjects)
	generalComponents.GET("/card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getCardView)
	generalComponents.GET("/group_by_card_view", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getGroupByCardView)
	generalComponents.GET("/new_record", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getNewRecord)
	generalComponents.GET("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getRecordFormData)
	generalComponents.PUT("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.updateResource)
	generalComponents.PUT("/record/:recordId/remove_internal_array_reference", middlewares.PermissionMiddleware(ModuleName), middlewares.TokenAuthMiddleware(), as.removeInternalArrayReference)

	generalComponents.DELETE("/record/:recordId", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.deleteResource)
	generalComponents.POST("/records", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.createNewResource)

	generalComponents.POST("/record/:recordId/action/reset_password", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.resetIndividualUserPassword)

	generalComponents.GET("/record/:recordId/action/send_email", middlewares.PermissionMiddleware(ModuleName), as.sendEmailToUser)
	generalComponents.POST("/record/:recordId/action/send_email_invitation", middlewares.PermissionMiddleware(ModuleName), as.sendEmailInvitation)

	generalComponents.POST("/search", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getSearchResults)

	generalComponents.POST("/export", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.exportObjects)
	generalComponents.GET("/export_schema", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.getExportSchema)

	iamGeneral := routerEngine.Group("/project/:projectId/iam")
	iamGeneral.GET("/user_overview", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.summaryResponse)

	generalComponents.GET("/record_messages/:recordId", as.getComponentRecordTrails)

	generalComponents.GET("/record/:recordId/action/delete_validation", middlewares.TokenAuthMiddleware(), middlewares.PermissionMiddleware(ModuleName), as.deleteValidation)
}
