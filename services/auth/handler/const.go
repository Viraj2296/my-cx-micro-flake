package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	AuthRecordTrailTable = "auth_record_trail"

	UserComponent          = "user"
	UserTable              = "user"
	UserGroupTable         = "user_group"
	RoleTable              = "role"
	JobRoleTable           = "job_role"
	EmployeeTypeTable      = "employee_type"
	PermissionTable        = "permission"
	ComponentResourceTable = "component_resource"
	IAMComponentTable      = "iam_component"
	APIAccessTable         = "api_access"
	APIAccessComponent     = "api_access"

	UserEmailTemplateTable       = "user_email_template"
	UserModuleSettingTable       = "user_module_setting"
	UserEmailTemplateFieldsTable = "user_email_template_field"
	ComponentName                = "iam"
	UserGroupComponent           = "user_group"
	SystemReleaseTable           = "system_release"

	UserStatusEnabled        = "enabled"
	ErrorGettingActionFields = 6056

	InvalidEmailSetting = 6057
	InvalidCredentials  = 6058

	LoadingUsersFailed                = 1000
	PasswordVerificationFailed        = 1001
	AuthenticationFailed              = 1002
	ErrorInUserInsertToDatabase       = 1003
	ErrorGettingObjectsInformation    = 1004
	ErrorDeletingUserFromDatabase     = 1005
	UserNotFoundInDatabase            = 1006
	ErrorUpdatingUserDataInToDatabase = 1007

	ErrorGettingPermissionsFromDatabase     = 1008
	InvalidUserRoleToGetPermissions         = 1009
	ErrorGettingIndividualObjectInformation = 1010
	ErrorRemovingObjectInformation          = 1011
	ErrorUpdatingObjectInformation          = 1012
	ErrorCreatingObjectInformation          = 1013
	InvalidAccess                           = 1014
	ExpectedFieldNotFound                   = 1015
	UnsupportedSchema                       = 1016
	InvalidToken                            = 1017
	InvalidUserStatus                       = 1018

	SystemError = 1019

	VersionMismatch = 1020

	AccessDenied            = "Access Denied"
	AccessDeniedDescription = "You are not authorised to access the requested resources, permission is not given to explore the resources, please check the permission level, and access principal"

	ForgetPasswordEmailTemplateType = 1
	WelcomeMESEmailTemplateType     = 2

	UserTypeBot         = "bot"
	UserTypeSystemAdmin = "system_admin"
	UserTypeUser        = "user"

	TimeLayout = "2006-01-02T15:04:05.000Z"
	ProjectID  = "906d0fd569404c59956503985b330132"

	ModuleName = "iam"
)

func getError(errorString string) error {
	return errors.New(errorString)
}

func sendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}
func sendArchiveFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Archived Failed",
			Description: "Internal system error during archive process. This is normally happen when the system is not configured properly. Please report to system administrator",
		})
	return
}
