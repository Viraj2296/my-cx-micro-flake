package handler

import (
	"cx-micro-flake/pkg/common/component"

	"gorm.io/datatypes"
)

type SearchKeys struct {
	Field string `json:"key"`
	Value string `json:"value"`
}

type AccessToken struct {
	Token string `json:"token"`
}

type LoginCommand struct {
	Username             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"required"`
	Remember             bool   `json:"remember"`
	Fingerprint          string `json:"fingerprint"`
	Platform             string `json:"platform"`
	AppName              string `json:"appName"`
	DeviceToken          string `json:"deviceToken"`
	UiApplicationVersion string `json:"uiApplicationVersion"`
}

type ForgetPassword struct {
	Email string `json:"email" binding:"required"`
}

type CreateUserCommand struct {
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	ZipCode              string `json:"zipCode"`
	Address1             string `json:"address1"`
	Address2             string `json:"address2"`
	Email                string `json:"email"`
	UserRole             string `json:"userRole"`
	Unsubscribed         bool   `json:"unsubscribed"`
	SendWelcomeEmail     bool   `json:"sendWelcomeEmail"`
	LogoutAllSessions    bool   `json:"logoutAllSessions"`
	SimultaneousSessions bool   `json:"simultaneousSessions"`
	TimeZone             string `json:"timeZone"`
	MuteSounds           bool   `json:"muteSounds"`
	Password             string `json:"password"`
	About                string `json:"about"`
	Phone                string `json:"phone"`
	Gender               string `json:"gender"`
	FullName             string `json:"fullName"`
	Language             string `json:"language"`
	LastName             string `json:"lastName"`
	Location             string `json:"location"`
	MobileNo             string `json:"mobileNo"`
	Username             string `json:"username"`
	AvatarUrl            string `json:"avatarUrl"`
	BirthDate            string `json:"birthDate"`
	FirstName            string `json:"firstName"`
	MiddleName           string `json:"middleName"`
	RoleProfileName      string `json:"roleProfileName"`
	Youtube              string `json:"youtube"`
	Facebook             string `json:"facebook"`
	Twitter              string `json:"twitter"`
	InvitationStatus     string `json:"invitationStatus"`
}

type UpdateUserCommand struct {
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	ZipCode              string `json:"zipCode"`
	Address1             string `json:"address1"`
	Address2             string `json:"address2"`
	Email                string `json:"email"`
	UserRole             string `json:"userRole"`
	Unsubscribed         bool   `json:"unsubscribed"`
	SendWelcomeEmail     bool   `json:"sendWelcomeEmail"`
	LogoutAllSessions    bool   `json:"logoutAllSessions"`
	SimultaneousSessions bool   `json:"simultaneousSessions"`
	TimeZone             string `json:"timeZone"`
	MuteSounds           bool   `json:"muteSounds"`
	Password             string `json:"password"`
	NewPassword          string `json:"newPassword"`
	About                string `json:"about"`
	Phone                string `json:"phone"`
	Gender               string `json:"gender"`
	FullName             string `json:"fullName"`
	Language             string `json:"language"`
	LastName             string `json:"lastName"`
	Location             string `json:"location"`
	MobileNo             string `json:"mobileNo"`
	Username             string `json:"username"`
	AvatarUrl            string `json:"avatarUrl"`
	BirthDate            string `json:"birthDate"`
	FirstName            string `json:"firstName"`
	MiddleName           string `json:"middleName"`
	RoleProfileName      string `json:"roleProfileName"`
	Youtube              string `json:"youtube"`
	Facebook             string `json:"facebook"`
	Twitter              string `json:"twitter"`
}

type TableObjectResponse struct {
	TotalRowCount   int64                   `json:"totalRowCount"`
	CurrentRowCount int64                   `json:"currentRowCount"`
	Header          []component.TableSchema `json:"header"`
	Data            []datatypes.JSON        `json:"data"`
}

type ModuleReports struct {
	MenuId      string `json:"menuId"`
	DashboardId int    `json:"dashboardId"`
}
type Reports struct {
	ModuleId      int             `json:"moduleId"`
	ModuleReports []ModuleReports `json:"moduleReports"`
}
type SystemConfig struct {
	SchedulerViewEndDate                  string               `json:"schedulerViewEndDate"`
	SchedulerViewStartDate                string               `json:"schedulerViewStartDate"`
	TVDisplay                             interface{}          `json:"tvDisplay"`
	AllowedMenusIds                       []string             `json:"allowedMenusIds"`
	AllowedDepartments                    component.RecordInfo `json:"allowedDepartments"`
	AllowedModules                        []int                `json:"allowedModules"`
	MaintenanceMessage                    string               `json:"maintenanceMessage"`
	Reports                               []Reports            `json:"reports"`
	ListOfTVLabourManagementAssemblyLines []int                `json:"listOfTVLabourManagementLines"`
	EnergyManagementMachines              []int                `json:"energyManagementMachines"`
}
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`

	SystemConfig SystemConfig `json:"systemConfig"`
}

type CreateUserGroupCommand struct {
	GroupName        string   `json:"groupName"`
	Permissions      []string `json:"permissions"`
	GroupDescription string   `json:"groupPermission"`
}

type UpdateUserGroupCommand struct {
	GroupName        string   `json:"groupName"`
	Permissions      []string `json:"permissions"`
	GroupDescription string   `json:"groupPermission"`
}
