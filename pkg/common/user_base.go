package common

type UserBasicInfo struct {
	FullName                string `json:"fullName"`
	Username                string `json:"userName"`
	Email                   string `json:"email"`
	AvatarUrl               string `json:"avatarUrl"`
	UserId                  int    `json:"userId"`
	NotificationLimit       int    `json:"notificationLimit"`
	Department              []int  `json:"department"`
	Section                 int    `json:"section"`
	SecondaryDepartmentList []int  `json:"secondaryDepartmentList"`
	Site                    int    `json:"site"`
	JobRole                 int    `json:"jobRole"`
	EmployeeNumber          string `json:"employeeNumber"`
	EmployeeTypeId          int    `json:"employeeTypeId"`
	SecondarySiteList       []int  `json:"secondarySiteList"`
}

type UserBasicInfoServiceResponse struct {
	FullName  string `json:"fullName"`
	Username  string `json:"userName"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatarUrl"`
	Id        int    `json:"id"`
}

type SystemModuleInfo struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	ModuleName  string `json:"moduleName"`
	Id          int    `json:"id"`
}
type NotificationMetaInfo struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"createdAt"`
}

// ExtractNotificationIds extracts the Ids from a slice of NotificationMetaInfo.
func ExtractNotificationIds(metaInfos []NotificationMetaInfo) []int {
	ids := make([]int, len(metaInfos))
	for i, info := range metaInfos {
		ids[i] = info.Id
	}
	return ids
}
