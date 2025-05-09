package handler

import (
	"cx-micro-flake/pkg/auth"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"

	"gorm.io/datatypes"
)

func (as *AuthService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := as.BaseService.ReferenceDatabase
	err, listOfObjects := GetObjects(dbConnection, IAMComponentTable)
	if err == nil {
		for _, objectInterface := range listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (as *AuthService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := as.BaseService.ReferenceDatabase
	err, listOfObjects := GetConditionalObjects(dbConnection, IAMComponentTable, conditionQuery)
	if err == nil {
		if len(listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, IAMComponentTable, (listOfObjects)[0].Id, updatingData)
		as.LoadInitComponents()
		return err
	}
	return err

}

func (as *AuthService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := as.BaseService.ReferenceDatabase
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := CreateFromGeneralObject(dbConnection, IAMComponentTable, generalObject)
	if err == nil {
		as.LoadInitComponents()
	}
	return recordId, err
}

func (as *AuthService) GetBasicInfo2TableFromUserList(userList reflect.Value) component.TableDataResponse {
	tableDataResponse := component.TableDataResponse{}
	var headerSchema []component.TableSchema
	userListSlice := userList.Interface().([]int)
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Email", "email"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Full Name", "fullName"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Avatar Url", "avatarUrl"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Id", "id"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Username", "username"))

	if len(userListSlice) == 0 {
		var emptyArray = make([]datatypes.JSON, 0)
		tableDataResponse.Data = emptyArray
		tableDataResponse.Header = headerSchema
		return tableDataResponse
	}
	replacementValue := util.InterfaceArrayToCommaSeperatedString(userListSlice)
	condition := " id IN (" + replacementValue + ")"
	err, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)

	if err != nil {
		var emptyArray = make([]datatypes.JSON, 0)
		tableDataResponse.Data = emptyArray
		tableDataResponse.Header = headerSchema
	} else {
		for _, userObject := range listOfUsers {
			userInfo := GetUserInfo(userObject.ObjectInfo)
			userBasicInfo := common.UserBasicInfoServiceResponse{
				Username:  userInfo.Username,
				FullName:  userInfo.FullName,
				AvatarUrl: userInfo.AvatarUrl,
				Id:        userObject.Id,
				Email:     userInfo.Email,
			}
			rawUserInfo, _ := json.Marshal(userBasicInfo)
			tableDataResponse.Header = headerSchema
			tableDataResponse.Data = append(tableDataResponse.Data, rawUserInfo)
		}

	}
	return tableDataResponse
}

func (as *AuthService) GetUserType(userId int) string {

	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return ""
	}

	return GetUserInfo(generalObject.ObjectInfo).Type

}
func (as *AuthService) GetAllUserBasicInfo2Table() component.TableDataResponse {
	err, listOfUsers := GetObjects(as.BaseService.ReferenceDatabase, UserTable)

	tableDataResponse := component.TableDataResponse{}
	var headerSchema []component.TableSchema

	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Email", "email"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Full Name", "fullName"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Avatar Url", "avatarUrl"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Id", "id"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Username", "username"))

	if err != nil {
		var emptyArray = make([]datatypes.JSON, 0)
		tableDataResponse.Data = emptyArray
		tableDataResponse.Header = headerSchema
	} else {
		for _, userObject := range listOfUsers {
			userInfo := GetUserInfo(userObject.ObjectInfo)
			if userInfo.Type == "user" {
				userBasicInfo := common.UserBasicInfoServiceResponse{
					Username:  userInfo.Username,
					FullName:  userInfo.FullName,
					AvatarUrl: userInfo.AvatarUrl,
					Id:        userObject.Id,
					Email:     userInfo.Email,
				}
				rawUserInfo, _ := json.Marshal(userBasicInfo)
				tableDataResponse.Header = headerSchema
				tableDataResponse.Data = append(tableDataResponse.Data, rawUserInfo)
			}

		}

	}
	return tableDataResponse
}

func (as *AuthService) GetAllUserBasicInfo2QueryResults() []datatypes.JSON {
	err, listOfUsers := GetObjects(as.BaseService.ReferenceDatabase, UserTable)

	var userInfoResults []datatypes.JSON

	if err == nil {

		for _, userObject := range listOfUsers {
			userInfo := GetUserInfo(userObject.ObjectInfo)
			if userInfo.Type == "user" {
				userBasicInfo := common.UserBasicInfoServiceResponse{
					Username:  userInfo.Username,
					FullName:  userInfo.FullName,
					AvatarUrl: userInfo.AvatarUrl,
					Id:        userObject.Id,
					Email:     userInfo.Email,
				}
				rawUserInfo, _ := json.Marshal(userBasicInfo)
				userInfoResults = append(userInfoResults, rawUserInfo)
			}

		}

	}
	return userInfoResults
}

func (as *AuthService) GetBasicInfo2TableFromUserId(userId int) component.TableDataResponse {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)

	tableDataResponse := component.TableDataResponse{}
	var headerSchema []component.TableSchema

	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Email", "email"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Full Name", "fullName"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Avatar Url", "avatarUrl"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Id", "id"))
	headerSchema = append(headerSchema, component.GetTableHeaderSchema("Username", "username"))

	if err != nil {
		var emptyArray = make([]datatypes.JSON, 0)
		tableDataResponse.Data = emptyArray
		tableDataResponse.Header = headerSchema
	} else {
		userInfo := GetUserInfo(generalObject.ObjectInfo)
		userBasicInfo := common.UserBasicInfoServiceResponse{
			Username:  userInfo.Username,
			FullName:  userInfo.FullName,
			AvatarUrl: userInfo.AvatarUrl,
			Id:        generalObject.Id,
			Email:     userInfo.Email,
		}
		rawUserInfo, _ := json.Marshal(userBasicInfo)
		tableDataResponse.Data = append(tableDataResponse.Data, rawUserInfo)
	}

	return tableDataResponse
}
func (as *AuthService) IsAnyDepartmentHeadExist(departmentId int) bool {
	condition := " " + strconv.Itoa(departmentId) + " MEMBER OF(object_info->>'$.hodDepartment') " + " AND object_info ->>'$.isDepartmentHead' =  'true'"
	err, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	if err == nil {
		if len(listOfUsers) == 0 {
			return false
		} else {
			return true
		}
	}
	return false

}
func (as *AuthService) GetHeadOfDepartment(departmentId int) common.UserBasicInfo {
	condition := " " + strconv.Itoa(departmentId) + " MEMBER OF(object_info->>'$.hodDepartment') " + " AND object_info ->>'$.isDepartmentHead' =  'true'"
	err, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	if err == nil {
		userInfo := GetUserInfo((listOfUsers)[0].ObjectInfo)
		userBasicInfo := common.UserBasicInfo{
			Username:          userInfo.Username,
			FullName:          userInfo.FullName,
			AvatarUrl:         userInfo.AvatarUrl,
			UserId:            (listOfUsers)[0].Id,
			NotificationLimit: userInfo.NotificationLimit,
			Email:             userInfo.Email,
		}
		return userBasicInfo
	} else {

		userBasicInfo := common.UserBasicInfo{
			Username:                "",
			AvatarUrl:               "",
			FullName:                "",
			UserId:                  0,
			NotificationLimit:       0,
			Email:                   "",
			Department:              make([]int, 0),
			Section:                 0,
			SecondaryDepartmentList: make([]int, 0),
			Site:                    0,
			JobRole:                 0,
		}
		return userBasicInfo
	}
}

func (as *AuthService) GetDepartmentUsers(userId int) (bool, []common.UserBasicInfo) {

	var listOfDepartmentUsers []common.UserBasicInfo
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	var isDepartmentHead bool
	isDepartmentHead = false
	if err == nil {
		userInfo := GetUserInfo(generalObject.ObjectInfo)
		departmentIdList := userInfo.Department
		isDepartmentHead = userInfo.IsDepartmentHead
		for _, departmentId := range departmentIdList {
			condition := " " + strconv.Itoa(departmentId) + " MEMBER OF(object_info->>'$.department') "
			_, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
			for _, userInterfaceObject := range listOfUsers {
				userInternalInfo := GetUserInfo(userInterfaceObject.ObjectInfo)
				listOfDepartmentUsers = append(listOfDepartmentUsers, common.UserBasicInfo{
					Username:          userInternalInfo.Username,
					FullName:          userInternalInfo.FullName,
					AvatarUrl:         userInternalInfo.AvatarUrl,
					UserId:            userInterfaceObject.Id,
					NotificationLimit: userInternalInfo.NotificationLimit,
					Email:             userInternalInfo.Email,
				})

			}
		}

	}

	return isDepartmentHead, listOfDepartmentUsers
}

func (as *AuthService) GetUsersByDepartment(departmentId int) []common.UserBasicInfo {

	var listOfDepartmentUsers []common.UserBasicInfo

	condition := " " + strconv.Itoa(departmentId) + " MEMBER OF(object_info->>'$.department') "
	_, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	for _, userInterfaceObject := range listOfUsers {
		userInfo := GetUserInfo(userInterfaceObject.ObjectInfo)
		listOfDepartmentUsers = append(listOfDepartmentUsers, common.UserBasicInfo{
			Username:          userInfo.Username,
			FullName:          userInfo.FullName,
			AvatarUrl:         userInfo.AvatarUrl,
			UserId:            userInterfaceObject.Id,
			NotificationLimit: userInfo.NotificationLimit,
			Email:             userInfo.Email,
		})

	}

	return listOfDepartmentUsers
}

func (as *AuthService) GetHeadOfDepartments(userId int) []common.UserBasicInfo {
	var listOfHeads []common.UserBasicInfo
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err == nil {
		userInfo := GetUserInfo(generalObject.ObjectInfo)
		departmentIdList := userInfo.Department
		for _, departmentId := range departmentIdList {
			condition := " " + strconv.Itoa(departmentId) + " MEMBER OF(object_info->>'$.hodDepartment') " + " AND object_info ->>'$.isDepartmentHead' =  'true'"
			_, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
			for _, userInterfaceObject := range listOfUsers {
				interalUserInfo := GetUserInfo(userInterfaceObject.ObjectInfo)
				listOfHeads = append(listOfHeads, common.UserBasicInfo{
					Username:          interalUserInfo.Username,
					FullName:          interalUserInfo.FullName,
					AvatarUrl:         interalUserInfo.AvatarUrl,
					UserId:            userInterfaceObject.Id,
					NotificationLimit: interalUserInfo.NotificationLimit,
					Email:             interalUserInfo.Email,
				})
			}
		}
	}

	return listOfHeads
}

func (as *AuthService) GetHeadOfSections(userId int) []common.UserBasicInfo {
	var listOfHeads []common.UserBasicInfo
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err == nil {
		userInfo := GetUserInfo(generalObject.ObjectInfo)
		sectionId := userInfo.Section
		condition := " object_info ->>'$.section' = " + strconv.Itoa(sectionId) + " AND object_info ->>'$.isSectionHead' =  'true'"

		_, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)

		for _, userInterfaceObject := range listOfUsers {
			userInfo := GetUserInfo(userInterfaceObject.ObjectInfo)
			listOfHeads = append(listOfHeads, common.UserBasicInfo{
				Username:          userInfo.Username,
				FullName:          userInfo.FullName,
				AvatarUrl:         userInfo.AvatarUrl,
				UserId:            userInterfaceObject.Id,
				NotificationLimit: userInfo.NotificationLimit,
				Email:             userInfo.Email,
			})
		}
	}

	return listOfHeads
}

func (as *AuthService) UpdateComponentResource(rawComponentResourceInfo datatypes.JSON) error {

	componentResourceInfo := ComponentResourceInfo{}
	json.Unmarshal(rawComponentResourceInfo, &componentResourceInfo)
	condition := "object_info ->>'$.moduleId' =" + strconv.Itoa(componentResourceInfo.ModuleId) + " AND object_info ->> '$.method' = '" + componentResourceInfo.Method + "' AND object_info ->> '$.pattern' ='" + componentResourceInfo.Pattern + "' AND object_info ->>'$.resource' = '" + componentResourceInfo.Resource + "'"
	_, listOfObjects := GetConditionalObjects(as.BaseService.ReferenceDatabase, ComponentResourceTable, condition)
	var err error
	if len(listOfObjects) == 0 {
		// nothing inserted, just insert
		commonObject := component.GeneralObject{ObjectInfo: rawComponentResourceInfo}
		err, _ = CreateFromGeneralObject(as.BaseService.ReferenceDatabase, ComponentResourceTable, commonObject)
	} else {
		recordId := (listOfObjects)[0].Id
		err = Update(as.BaseService.ReferenceDatabase, ComponentResourceTable, recordId, componentResourceInfo.DatabaseSerialize())
	}

	return err

}
func (as *AuthService) GetUserTimezone(userId int) string {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {

		return "Asia/Singapore"
	}

	return GetUserInfo(generalObject.ObjectInfo).TimeZone
}

func (as *AuthService) ConvertToUserTimezoneToISO(userId int, datetime string) string {

	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {

		return "Asia/Singapore"
	}

	userTimezone := GetUserInfo(generalObject.ObjectInfo).TimeZone
	return util.ConvertUserTimezoneToUTC(userTimezone, datetime)
}
func (as *AuthService) GetJobRoleName(jobRoleId int) string {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, JobRoleTable, jobRoleId)
	if err != nil {
		return ""
	} else {
		return GetJobRoleInfo(generalObject.ObjectInfo).JobTitleName
	}
}
func (as *AuthService) GetJobRoleHierarchy(jobRoleId int) int {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, JobRoleTable, jobRoleId)
	if err != nil {
		return -1
	} else {
		return GetJobRoleInfo(generalObject.ObjectInfo).HierarchyLevel
	}
}
func (as *AuthService) GetUserInfoById(userRecordId int) common.UserBasicInfo {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userRecordId)
	if err != nil {
		userBasicInfo := common.UserBasicInfo{
			Username:                "",
			AvatarUrl:               "",
			FullName:                "",
			UserId:                  0,
			NotificationLimit:       0,
			Email:                   "",
			Department:              make([]int, 0),
			Section:                 0,
			SecondaryDepartmentList: make([]int, 0),
			Site:                    0,
			JobRole:                 0,
			SecondarySiteList:       make([]int, 0),
		}
		return userBasicInfo
	}
	userInfo := GetUserInfo(generalObject.ObjectInfo)
	userBasicInfo := common.UserBasicInfo{
		Username:                userInfo.Username,
		FullName:                userInfo.FullName,
		AvatarUrl:               userInfo.AvatarUrl,
		UserId:                  userRecordId,
		NotificationLimit:       userInfo.NotificationLimit,
		Email:                   userInfo.Email,
		Department:              userInfo.Department,
		Section:                 userInfo.Section,
		SecondaryDepartmentList: userInfo.SecondaryDepartmentList,
		Site:                    userInfo.Site,
		JobRole:                 userInfo.JobRoleId,
		SecondarySiteList:       userInfo.SecondarySiteList,
	}
	return userBasicInfo

}

func (as *AuthService) GetUserInfoByEmployeeId(employeeId string) common.UserBasicInfo {

	var conditionString = " object_info->>'$.employeeNumber'='" + employeeId + "'"
	err, listOfRecords := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, conditionString)
	if err != nil {
		userBasicInfo := common.UserBasicInfo{
			Username:                "",
			AvatarUrl:               "",
			FullName:                "",
			UserId:                  0,
			NotificationLimit:       0,
			Email:                   "",
			Department:              make([]int, 0),
			Section:                 0,
			SecondaryDepartmentList: make([]int, 0),
			Site:                    0,
			JobRole:                 0,
		}
		return userBasicInfo
	}
	if len(listOfRecords) == 1 {
		userInfo := GetUserInfo((listOfRecords)[0].ObjectInfo)
		userBasicInfo := common.UserBasicInfo{
			Username:                userInfo.Username,
			FullName:                userInfo.FullName,
			AvatarUrl:               userInfo.AvatarUrl,
			UserId:                  (listOfRecords)[0].Id,
			NotificationLimit:       userInfo.NotificationLimit,
			Email:                   userInfo.Email,
			Department:              userInfo.Department,
			Section:                 userInfo.Section,
			SecondaryDepartmentList: userInfo.SecondaryDepartmentList,
			Site:                    userInfo.Site,
			JobRole:                 userInfo.JobRoleId,
		}
		return userBasicInfo
	} else {
		userBasicInfo := common.UserBasicInfo{
			Username:                "",
			AvatarUrl:               "",
			FullName:                "",
			UserId:                  0,
			NotificationLimit:       0,
			Email:                   "",
			Department:              make([]int, 0),
			Section:                 0,
			SecondaryDepartmentList: make([]int, 0),
			Site:                    0,
		}
		return userBasicInfo
	}

}

func (as *AuthService) GetUserList() []common.UserBasicInfo {
	var arrayOfUserInfo []common.UserBasicInfo
	err, listOfUsers := GetObjects(as.BaseService.ReferenceDatabase, UserTable)
	if err != nil {
		userBasicInfo := common.UserBasicInfo{
			FullName:  "",
			Username:  "",
			AvatarUrl: "",
			Email:     "",
			UserId:    0,
		}
		arrayOfUserInfo = append(arrayOfUserInfo, userBasicInfo)
		return arrayOfUserInfo
	}
	for _, userObject := range listOfUsers {
		userInfo := GetUserInfo(userObject.ObjectInfo)
		userBasicInfo := common.UserBasicInfo{
			FullName:  userInfo.FullName,
			Username:  userInfo.Username,
			AvatarUrl: userInfo.AvatarUrl,
			Email:     userInfo.Email,
			UserId:    userObject.Id,
		}
		arrayOfUserInfo = append(arrayOfUserInfo, userBasicInfo)
	}

	return arrayOfUserInfo

}

func (as *AuthService) IsEmailExist(email string) bool {
	condition := "object_info ->>'$.email' = \"" + email + "\""
	err, listOfObjects := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	if err != nil {
		return false
	}
	if len(listOfObjects) == 1 {
		return true
	}
	return false
}

func (as *AuthService) EmailToUserId(email string) int {
	condition := "object_info ->>'$.email' = \"" + email + "\""
	err, listOfObjects := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	if err != nil {
		return -1
	}
	if len(listOfObjects) == 1 {
		return (listOfObjects)[0].Id
	}
	return -1
}

func (as *AuthService) GetUserInfoFromGroupId(groupId int) []common.UserBasicInfo {
	var arrayOfUserInfo []common.UserBasicInfo
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserGroupTable, groupId)
	if err != nil {
		userBasicInfo := common.UserBasicInfo{
			Username:  "",
			AvatarUrl: "",
			Email:     "",
			UserId:    0,
		}
		arrayOfUserInfo = append(arrayOfUserInfo, userBasicInfo)
		return arrayOfUserInfo
	}
	userGroup := UserGroup{ObjectInfo: generalObject.ObjectInfo}
	for _, userRecordId := range userGroup.GetGroupInfo().Users {
		_, generalObject = Get(as.BaseService.ReferenceDatabase, UserTable, userRecordId)
		userInfo := GetUserInfo(generalObject.ObjectInfo)
		userBasicInfo := common.UserBasicInfo{
			Username:  userInfo.Username,
			AvatarUrl: userInfo.AvatarUrl,
			Email:     userInfo.Email,
			UserId:    userRecordId,
		}
		arrayOfUserInfo = append(arrayOfUserInfo, userBasicInfo)
	}

	return arrayOfUserInfo

}

func (as *AuthService) GetBotInfo() common.UserBasicInfo {
	user := User{}
	var conditionString = " object_info->>'$.type' = 'bot'"
	err, objectsInterface := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, conditionString)
	if err != nil {
		return common.UserBasicInfo{}
	}
	if len(objectsInterface) > 0 {
		userInfo := GetUserInfo(objectsInterface[0].ObjectInfo)
		userBasicInfo := common.UserBasicInfo{
			Username:  userInfo.Username,
			AvatarUrl: userInfo.AvatarUrl,
			UserId:    user.Id,
		}
		return userBasicInfo
	} else {
		return common.UserBasicInfo{}
	}

}

func (as *AuthService) Authenticate(username, password string) (string, error) {

	err, listOfUsers := GetObjects(as.BaseService.ReferenceDatabase, UserTable)

	if err != nil {
		as.BaseService.Logger.Error("Loading users had failed due to :", zap.String("error", err.Error()))
		return "", err
	}
	for _, userInterface := range listOfUsers {
		userInfo := GetUserInfo(userInterface.ObjectInfo)
		if userInfo.Username == username {
			as.BaseService.Logger.Info("checking user password ", zap.String("user", userInfo.Username))
			err = VerifyPassword(userInfo.Password, password)
			if err != nil {
				as.BaseService.Logger.Info("verify password had failed ", zap.String("error", err.Error()))
				return "", err
			}
			zAuthToken, _ := auth.CreateToken(userInterface.Id, as.DefaultTokenExpiry)
			// zRefreshAuthToken, _ := auth.CreateRefreshToken(userInterface.Id, as.DefaultRefreshTokenExpiry)

			// user has successfully logged in, so update the last login and last active
			userInfo.LastLogin = time.Now().String()
			userInfo.LastActive = time.Now().String()
			Update(as.BaseService.ReferenceDatabase, UserTable, userInterface.Id, userInfo.DatabaseSerialize())

			// save the user session
			//userSession := UserSession{UserId: userObject.UserId, SessionId: zAuthToken}
			//sessionInfo := SessionInfo{SessionStatus: "logged-in", RequestHeader: ctx.Request.Header}
			//var clientIpArray []string
			//clientIp, _ := util.GetClientIPHelper(ctx.Request)
			//clientIpArray = append(clientIpArray, clientIp)
			//sessionInfo.RequestHeader["ClientIp"] = clientIpArray
			//userSession.SessionTime = time.Now().Unix()
			//userSession.SessionInfo = sessionInfo.Serialize()
			//as.BaseService.ReferenceDatabase.Model(&UserSession{}).Create(userSession)

			// return map[string]string{
			// 	"access_token":  zAuthToken,
			// 	"refresh_token": zRefreshAuthToken,
			// }, nil
			return zAuthToken, nil

		}
	}
	return "", errors.New("system error")
}

func (as *AuthService) GetSectionUsers(userId int) ([]int, bool) {
	sectionalUserList := make([]int, 0)
	isSectionalHead := false

	headOfSections := as.GetHeadOfSections(userId)
	for _, userInfoOfHead := range headOfSections {
		if userInfoOfHead.UserId == userId {
			isSectionalHead = true
		}
	}

	basicUserInfo := as.GetUserInfoById(userId)
	condition := " object_info ->>'$.section' = " + strconv.Itoa(basicUserInfo.Section)
	_, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)

	for _, userInfoOfSectionMate := range listOfUsers {
		sectionalUserList = append(sectionalUserList, userInfoOfSectionMate.Id)
	}

	return sectionalUserList, isSectionalHead
}

func (as *AuthService) GetAllDeviceTokenFromJobRole(jobRole int) []string {
	var deviceTokens = make([]string, 0)

	var condition = "object_info->>'$.jobRoleId' = '" + strconv.Itoa(jobRole) + "'"

	err, listOfUsers := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	if err != nil {
		return deviceTokens
	}

	for _, user := range listOfUsers {
		userId := user.Id
		if userId == 0 {
			continue
		}
		var deviceTokenRecords []UserDeviceToken
		if err := as.BaseService.ReferenceDatabase.Where("user_id = ?", userId).Find(&deviceTokenRecords).Error; err != nil {
			continue
		}
		for _, deviceTokenRecord := range deviceTokenRecords {
			deviceTokens = append(deviceTokens, deviceTokenRecord.DeviceToken)
		}
	}
	return deviceTokens
}

func (as *AuthService) GetUsersDeviceTokens(listOfUsers []int) []string {
	var listOfUserDeviceTokens = make([]string, 0)
	var deviceTokenRecords []UserDeviceToken
	err := as.BaseService.ReferenceDatabase.Where("id IN ?", listOfUsers).Find(&deviceTokenRecords).Error
	if err != nil {
		return listOfUserDeviceTokens
	}
	for _, deviceTokenRecord := range deviceTokenRecords {
		listOfUserDeviceTokens = append(listOfUserDeviceTokens, deviceTokenRecord.DeviceToken)
	}

	return listOfUserDeviceTokens
}

func (as *AuthService) GetUserOneSignalSubscriptionIds(userId int) []string {
	var deviceTokenRecord UserDeviceToken
	err := as.BaseService.ReferenceDatabase.Where("user_id = ?", userId).First(&deviceTokenRecord).Error
	if err != nil {
		as.BaseService.Logger.Error("Failed to fetch user device token", logging.Int("userId", userId), logging.Error(err))
		return []string{}
	}
	as.BaseService.Logger.Info("Fetched user device token", logging.Int("userId", userId), logging.String("deviceToken", deviceTokenRecord.DeviceToken))
	var subscriptionIds []string
	if deviceTokenRecord.SubscriptionId != nil {
		err = json.Unmarshal([]byte(deviceTokenRecord.SubscriptionId), &subscriptionIds)
		if err != nil {
			as.BaseService.Logger.Error("Failed to unmarshal subscription IDs", logging.Int("userId", userId), logging.Error(err))
			return []string{}
		}
	}
	if len(subscriptionIds) == 0 {
		as.BaseService.Logger.Warn("No subscription IDs found for user", logging.Int("userId", userId))
	}
	return subscriptionIds
}

func (as *AuthService) GetUsersFromJobRoleId(jobRole int) []int {
	var listOfUsers = make([]int, 0)

	var condition = "object_info->>'$.jobRoleId'= '" + strconv.Itoa(jobRole) + "'"

	err, userInterfaceObjects := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserTable, condition)
	if err != nil {
		return listOfUsers
	}

	as.BaseService.Logger.Info("listOfUsers : ", zap.Any("listOfUsers", listOfUsers))

	for _, user := range userInterfaceObjects {
		listOfUsers = append(listOfUsers, user.Id)
	}
	return listOfUsers
}

func (as *AuthService) GetNotificationList(userId int) []common.NotificationMetaInfo {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return make([]common.NotificationMetaInfo, 0)
	}
	return GetUserInfo(generalObject.ObjectInfo).NotificationIds
}

func (as *AuthService) GetViewNotificationList(userId int) []common.NotificationMetaInfo {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return make([]common.NotificationMetaInfo, 0)
	}
	return GetUserInfo(generalObject.ObjectInfo).ViewNotificationIds
}

func (as *AuthService) GetPushNotificationList(userId int) []common.NotificationMetaInfo {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return make([]common.NotificationMetaInfo, 0)
	}
	return GetUserInfo(generalObject.ObjectInfo).PushNotificationIds
}

func (as *AuthService) GetViewPushNotificationList(userId int) []common.NotificationMetaInfo {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return make([]common.NotificationMetaInfo, 0)
	}
	return GetUserInfo(generalObject.ObjectInfo).ViewPushNotificationIds
}

func (as *AuthService) AddNotificationIds(userId, notificationId int) error {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return err
	}
	userInfo := GetUserInfo(generalObject.ObjectInfo)
	notificationMeta := common.NotificationMetaInfo{Id: notificationId, CreatedAt: util.GetCurrentTime("2006-01-02T15:04:05.000Z")}
	userInfo.NotificationIds = append(userInfo.NotificationIds, notificationMeta)
	userInfo.NotificationIds = RemoveDuplicates(userInfo.NotificationIds)
	updatingData := make(map[string]interface{})
	rawUserInfo, _ := json.Marshal(userInfo)
	updatingData["object_info"] = rawUserInfo
	err = Update(as.BaseService.ReferenceDatabase, UserTable, userId, updatingData)
	return err
}

func (as *AuthService) AddPushNotificationIds(userId, notificationId int) error {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return err
	}
	userInfo := GetUserInfo(generalObject.ObjectInfo)
	notificationMeta := common.NotificationMetaInfo{Id: notificationId, CreatedAt: util.GetCurrentTime("2006-01-02T15:04:05.000Z")}
	userInfo.PushNotificationIds = append(userInfo.PushNotificationIds, notificationMeta)
	userInfo.PushNotificationIds = RemoveDuplicates(userInfo.PushNotificationIds)
	updatingData := make(map[string]interface{})
	rawUserInfo, err := json.Marshal(userInfo)
	if err == nil {
		updatingData["object_info"] = rawUserInfo
		err = Update(as.BaseService.ReferenceDatabase, UserTable, userId, updatingData)
	}
	return err
}

func (as *AuthService) AddViewPushNotificationIds(userId, notificationId int) error {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return err
	}
	userInfo := GetUserInfo(generalObject.ObjectInfo)
	notificationMeta := common.NotificationMetaInfo{Id: notificationId, CreatedAt: util.GetCurrentTime("2006-01-02T15:04:05.000Z")}
	userInfo.ViewPushNotificationIds = append(userInfo.ViewPushNotificationIds, notificationMeta)
	userInfo.ViewPushNotificationIds = RemoveDuplicates(userInfo.ViewPushNotificationIds)
	updatingData := make(map[string]interface{})
	rawUserInfo, err := json.Marshal(userInfo)
	if err == nil {
		updatingData["object_info"] = rawUserInfo
		err = Update(as.BaseService.ReferenceDatabase, UserTable, userId, updatingData)
	}
	return err
}
func (as *AuthService) AddViewNotificationIds(userId, notificationId int) error {
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)
	if err != nil {
		return err
	}
	userInfo := GetUserInfo(generalObject.ObjectInfo)
	notificationMeta := common.NotificationMetaInfo{Id: notificationId, CreatedAt: util.GetCurrentTime("2006-01-02T15:04:05.000Z")}
	userInfo.ViewNotificationIds = append(userInfo.ViewNotificationIds, notificationMeta)
	userInfo.ViewNotificationIds = RemoveDuplicates(userInfo.ViewNotificationIds)
	updatingData := make(map[string]interface{})
	rawUserInfo, _ := json.Marshal(userInfo)
	updatingData["object_info"] = rawUserInfo
	err = Update(as.BaseService.ReferenceDatabase, UserTable, userId, updatingData)
	return err

}
func (as *AuthService) getLatestVersion() (string, error) {
	condition := "object_info ->>'$.objectStatus' = 'Active' ORDER BY id DESC"

	err, generalObject := GetConditionalObjects(as.BaseService.ServiceDatabases[ProjectID], SystemReleaseTable, condition, 1)
	if err != nil {
		return "", err
	}
	systemReleaseInfo := GetSystemReleaseInfo(generalObject[0].ObjectInfo)

	return systemReleaseInfo.Version, nil
}
