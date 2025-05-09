package handler

import (
	"cx-micro-flake/pkg/common"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type AuthRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type UserDeviceTokens struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type UserDeviceToken struct {
	Id             int            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	UserId         int            `json:"userId" gorm:"not null;index:user_device_idx,unique"`
	DeviceToken    string         `json:"deviceToken" gorm:"not null;index:user_device_idx,unique"`
	CreatedAt      time.Time      `json:"createdAt" gorm:"not null"`
	LastUsedAt     time.Time      `json:"lastUsedAt" gorm:"not null"`
	CreatedBy      int            `json:"createdBy" gorm:"not null"`
	LastUsedBy     int            `json:"lastUsedBy" gorm:"not null"`
	SubscriptionId datatypes.JSON `json:"subscriptionId" gorm:"not null"`
}
type User struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func GetUserInfo(serialisedData datatypes.JSON) *UserInfo {
	userInfo := UserInfo{}
	err := json.Unmarshal(serialisedData, &userInfo)
	if err != nil {
		fmt.Print("error serializing userInfo", err.Error())
		return &UserInfo{}
	}
	return &userInfo
}
func GetJobRoleInfo(serialisedData datatypes.JSON) *JobRoleInfo {
	jobRoleInfo := JobRoleInfo{}
	err := json.Unmarshal(serialisedData, &jobRoleInfo)
	if err != nil {
		return &JobRoleInfo{}
	}
	return &jobRoleInfo
}

type GroupInfo struct {
	Name            string    `json:"name"`
	Roles           []int     `json:"roles"`
	Users           []int     `json:"users"`
	CreatedAt       time.Time `json:"createdAt"`
	CreatedBy       int       `json:"createdBy"`
	GroupImage      string    `json:"groupImage"`
	Description     string    `json:"description"`
	ObjectStatus    string    `json:"objectStatus"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy   int       `json:"lastUpdatedBy"`
	ApprovedModules []int     `json:"approvedModules"`
}

func (u *UserGroup) GetGroupInfo() *GroupInfo {
	groupInfo := GroupInfo{}
	json.Unmarshal(u.ObjectInfo, &groupInfo)
	return &groupInfo
}

type JobRoleInfo struct {
	CreatedAt      time.Time `json:"createdAt"`
	CreatedBy      int       `json:"createdBy"`
	Description    string    `json:"description"`
	JobTitleName   string    `json:"jobTitleName"`
	HierarchyLevel int       `json:"hierarchyLevel"`
	ObjectStatus   string    `json:"objectStatus"`
	LastUpdatedAt  string    `json:"lastUpdatedAt"`
	LastUpdatedBy  int       `json:"lastUpdatedBy"`
}
type UserInfo struct {
	City                    string                        `json:"city"`
	Country                 string                        `json:"country"`
	ZipCode                 string                        `json:"zipCode"`
	Address1                string                        `json:"address1"`
	Address2                string                        `json:"address2"`
	LastLogin               string                        `json:"lastLogin"`
	LastActive              string                        `json:"lastActive"`
	Email                   string                        `json:"email"`
	Status                  string                        `json:"status"`
	UserRoles               []int                         `json:"userRoles"`
	Unsubscribed            bool                          `json:"unsubscribed"`
	SendWelcomeEmail        bool                          `json:"sendWelcomeEmail"`
	LogoutAllSessions       bool                          `json:"logoutAllSessions"`
	SimultaneousSessions    bool                          `json:"simultaneousSessions"`
	TimeZone                string                        `json:"timeZone"`
	MuteSounds              bool                          `json:"muteSounds"`
	Password                string                        `json:"password"`
	NewPassword             string                        `json:"newPassword"`
	PlainPassword           string                        `json:"plainPassword"`
	LastPasswordResetDate   string                        `json:"lastPasswordResetDate"`
	ResetPasswordKey        string                        `json:"resetPasswordKey"`
	About                   string                        `json:"about"`
	Phone                   string                        `json:"phone"`
	Gender                  string                        `json:"gender"`
	FullName                string                        `json:"fullName"`
	Language                string                        `json:"language"`
	LastName                string                        `json:"lastName"`
	Location                string                        `json:"location"`
	MobileNo                string                        `json:"mobileNo"`
	Username                string                        `json:"username"`
	AvatarUrl               string                        `json:"avatarUrl"`
	BirthDate               string                        `json:"birthDate"`
	FirstName               string                        `json:"firstName"`
	MiddleName              string                        `json:"middleName"`
	Youtube                 string                        `json:"youtube"`
	Facebook                string                        `json:"facebook"`
	Twitter                 string                        `json:"twitter"`
	InvitationStatus        string                        `json:"invitationStatus"`
	InvitationToken         string                        `json:"invitationToken"`
	NotificationIds         []common.NotificationMetaInfo `json:"notificationIds"`
	ViewNotificationIds     []common.NotificationMetaInfo `json:"viewNotificationIds"`
	CreatedAt               string                        `json:"createdAt"`
	CreatedBy               int                           `json:"createdBy"`
	LastUpdatedAt           string                        `json:"lastUpdatedAt"`
	LastUpdatedBy           int                           `json:"lastUpdatedBy"`
	NotificationLimit       int                           `json:"notificationLimit"`
	Type                    string                        `json:"type"`
	Department              []int                         `json:"department"`
	HodDepartment           []int                         `json:"hodDepartment"`
	EmployeeNumber          string                        `json:"employeeNumber"`
	Position                int                           `json:"position"`
	IsDepartmentHead        bool                          `json:"isDepartmentHead"`
	SecondaryDepartmentList []int                         `json:"secondaryDepartmentList"`
	ObjectStatus            string                        `json:"objectStatus"`
	Section                 int                           `json:"section"`
	IsSectionHead           bool                          `json:"isSectionHead"`
	IsSessionTimeOut        bool                          `json:"isSessionTimeOut"`
	Site                    int                           `json:"site"`
	JobRoleId               int                           `json:"jobRoleId"`
	SecondarySiteList       []int                         `json:"secondarySiteList"`
	PushNotificationIds     []common.NotificationMetaInfo `json:"pushNotificationIds"`
	ViewPushNotificationIds []common.NotificationMetaInfo `json:"viewPushNotificationIds"`
	UiApplicationVersion    string                        `json:"uiApplicationVersion"`
	EnforceVersionUpdate    bool                          `json:"enforceVersionUpdate"`
}

func (userInfo *UserInfo) Serialize() []byte {
	rawData, _ := json.Marshal(userInfo)
	return rawData
}

func (userInfo *UserInfo) DatabaseSerialize() map[string]interface{} {
	updatingData := make(map[string]interface{})
	updatingData["object_info"] = userInfo.Serialize()
	return updatingData
}

func (userInfo *GroupInfo) Serialize() []byte {
	rawData, _ := json.Marshal(userInfo)
	return rawData
}

func (userInfo *GroupInfo) DatabaseSerialize() map[string]interface{} {
	updatingData := make(map[string]interface{})
	updatingData["object_info"] = userInfo.Serialize()
	return updatingData
}
func VerifyPassword(srcPassword, givenPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(srcPassword), []byte(givenPassword))
}

type IAMComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type UserGroup struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type UserEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type UserEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type UserModuleSetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ComponentResource struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (cr *ComponentResource) GetComponentResourceInfo() *ComponentResourceInfo {
	componentResourceInfo := ComponentResourceInfo{}
	json.Unmarshal(cr.ObjectInfo, &componentResourceInfo)
	return &componentResourceInfo
}

func (v *ComponentResourceInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *ComponentResourceInfo) DatabaseSerialize() map[string]interface{} {
	updatingData := make(map[string]interface{})
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

type Role struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type JobRole struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type EmployeeType struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type APIAccess struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type RoleInfo struct {
	Name                string    `json:"name"`
	CreatedAt           time.Time `json:"createdAt"`
	CreatedBy           int       `json:"createdBy"`
	Description         string    `json:"description"`
	PermissionResources []int     `json:"permissionResources"`
	LastUpdatedAt       time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy       int       `json:"lastUpdatedBy"`
	PermissionLevel     int       `json:"permissionLevel"`
	ListOfAllowedMenus  []int     `json:"listOfAllowedMenus"`
}

func (r *Role) getRoleInfo() *RoleInfo {
	roleInfo := RoleInfo{}
	json.Unmarshal(r.ObjectInfo, &roleInfo)
	return &roleInfo
}

type ComponentResourceInfo struct {
	ModuleId          int    `json:"moduleId"`
	Resource          string `json:"resource"`
	ResourceId        string `json:"resourceId"`
	RoutingComponent  string `json:"routingComponent"`
	IsRouteEnabled    bool   `json:"isRouteEnabled"`
	Method            string `json:"method"`
	Action            string `json:"action"`
	ProjectId         string `json:"projectId"`
	Pattern           string `json:"pattern"`
	CreatedAt         string `json:"createdAt"`
	CreatedBy         int    `json:"createdBy"`
	LastUpdatedAt     string `json:"lastUpdatedAt"`
	LastUpdatedBy     int    `json:"lastUpdatedBy"`
	ComponentAction   string `json:"componentAction"`
	ResourceDisplay   string `json:"resourceDisplay"`
	ActionDescription string `json:"actionDescription"`
}

type Permission struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type PermissionInfo struct {
	ComponentName       string `json:"componentName"`
	PermissionResources []int  `json:"permissionResources"`
}

func (pm *Permission) getPermissionInfo() *PermissionInfo {
	permissionInfo := PermissionInfo{}
	json.Unmarshal(pm.ObjectInfo, &permissionInfo)
	return &permissionInfo
}

// APIAccessInfo this hold the api access info
type APIAccessInfo struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	EndDate       string   `json:"endDate"`
	IsEnabled     bool     `json:"isEnabled"`
	LastUpdatedBy int      `json:"lastUpdatedBy"`
	LastUpdatedAt string   `json:"lastUpdatedAt"`
	CreatedAt     string   `json:"createdAt"`
	CreatedBy     int      `json:"createdBy"`
	ObjectStatus  string   `json:"objectStatus"`
	EndPoints     []string `json:"endPoints"`
	APIKey        string   `json:"apiKey"`
	UserId        int      `json:"userId"`
}

func GetAPIAccessInfo(serialisedData datatypes.JSON) *APIAccessInfo {
	apiAccessInfo := APIAccessInfo{}
	err := json.Unmarshal(serialisedData, &apiAccessInfo)
	if err != nil {
		return &APIAccessInfo{}
	}
	return &apiAccessInfo
}

func (v *APIAccessInfo) serialised() datatypes.JSON {
	serialisedData, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON{}
	}
	return serialisedData

}

type SystemModuleInfo struct {
	Name                string   `json:"name"`
	Labels              []string `json:"labels"`
	Version             string   `json:"version"`
	CreatedAt           string   `json:"createdAt"`
	CreatedBy           int      `json:"createdBy"`
	Attachment          string   `json:"attachment"`
	Description         string   `json:"description"`
	DisplayName         string   `json:"displayName"`
	ReleaseDate         string   `json:"releaseDate"`
	HasAttachment       bool     `json:"hasAttachment"`
	LastUpdatedAt       string   `json:"lastUpdatedAt"`
	DetailedDescription string   `json:"detailedDescription"`
	ServiceRegisterName string   `json:"serviceRegisterName"`
}

func GetSystemModuleInfo(serialisedData datatypes.JSON) *SystemModuleInfo {
	systemModuleInfo := SystemModuleInfo{}
	err := json.Unmarshal(serialisedData, &systemModuleInfo)
	if err != nil {
		return &SystemModuleInfo{}
	}
	return &systemModuleInfo
}

func RemoveDuplicates(metaInfos []common.NotificationMetaInfo) []common.NotificationMetaInfo {
	seen := make(map[int]bool)
	var result []common.NotificationMetaInfo

	for _, info := range metaInfos {
		if !seen[info.Id] {
			seen[info.Id] = true
			result = append(result, info)
		}
	}

	return result
}

type SystemReleaseInfo struct {
	Labels              []string `json:"labels"`
	Version             string   `json:"version"`
	CreatedAt           string   `json:"createdAt"`
	CreatedBy           int      `json:"createdBy"`
	Attachment          string   `json:"attachment"`
	Description         string   `json:"description"`
	ReleaseDate         string   `json:"releaseDate"`
	HasAttachment       bool     `json:"hasAttachment"`
	LastUpdatedAt       string   `json:"lastUpdatedAt"`
	DetailedDescription string   `json:"detailedDescription"`
	ObjectStatus        string   `json:"objectStatus"`
	LastUpdatedBy       string   `json:"lastUpdatedBy"`
	ReleaseModules      []int    `json:"releaseModules"`
}

func GetSystemReleaseInfo(serialisedData datatypes.JSON) *SystemReleaseInfo {
	systemReleaseInfo := SystemReleaseInfo{}
	err := json.Unmarshal(serialisedData, &systemReleaseInfo)
	if err != nil {
		return &SystemReleaseInfo{}
	}
	return &systemReleaseInfo
}
