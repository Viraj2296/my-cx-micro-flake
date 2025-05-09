package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"os"
	"regexp"
	"strings"
)

type SystemMenuInfo struct {
	Name          string `json:"name"`
	MenuId        string `json:"menuId"`
	ModuleId      int    `json:"moduleId"`
	CreatedAt     string `json:"createdAt"`
	CreatedBy     int    `json:"createdBy"`
	Description   string `json:"description"`
	DisplayName   string `json:"displayName"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

func GetSystemMenuInfo(serialisedData datatypes.JSON) *SystemMenuInfo {
	systemMenuInfo := SystemMenuInfo{}
	json.Unmarshal(serialisedData, &systemMenuInfo)
	return &systemMenuInfo
}
func (as *AuthService) getUserModuleAccess(userId int) []int {
	var userModuleMap = make(map[int][]int, 0)

	err, listOfGroups := GetObjects(as.BaseService.ReferenceDatabase, UserGroupTable)
	var listOfModules = make([]int, 0)
	if err == nil {
		for _, object := range listOfGroups {
			userGroup := UserGroup{
				Id:         object.Id,
				ObjectInfo: object.ObjectInfo,
			}
			userList := userGroup.GetGroupInfo().Users
			approvedModules := userGroup.GetGroupInfo().ApprovedModules
			for _, userId := range userList {
				for _, moduleId := range approvedModules {
					userModuleMap[userId] = append(userModuleMap[userId], moduleId)
				}
			}
		}
		for key, value := range userModuleMap {
			userModuleMap[key] = util.RemoveDuplicateInt(value)
		}

		if value, ok := userModuleMap[userId]; ok {
			return value
		} else {
			return listOfModules
		}
	} else {
		as.BaseService.Logger.Error("get user module access failed", zap.String("err", err.Error()))
	}
	return listOfModules
}

func (as *AuthService) getListOfAllowedMenus(userId int) []string {
	var userMap = make(map[int][]int, 0)

	err, listOfGroups := GetObjects(as.BaseService.ReferenceDatabase, UserGroupTable)
	var listOfAllowedMenus = make([]string, 0)
	if err == nil {
		for _, object := range listOfGroups {
			userGroup := UserGroup{
				Id:         object.Id,
				ObjectInfo: object.ObjectInfo,
			}
			userList := userGroup.GetGroupInfo().Users
			rolesList := userGroup.GetGroupInfo().Roles
			for _, userId := range userList {
				for _, roleId := range rolesList {
					userMap[userId] = append(userMap[userId], roleId)
				}
			}
		}
		for key, value := range userMap {
			userMap[key] = util.RemoveDuplicateInt(value)
		}

		listOfRoles := userMap[userId]

		for _, roleId := range listOfRoles {
			_, roleInterface := Get(as.BaseService.ReferenceDatabase, RoleTable, roleId)
			userRole := Role{ObjectInfo: roleInterface.ObjectInfo}
			allowedMenuIds := userRole.getRoleInfo().ListOfAllowedMenus

			for _, systemMenuId := range allowedMenuIds {
				err, systemMenuObject := Get(as.BaseService.ServiceDatabases[ProjectID], "system_menu", systemMenuId)
				if err == nil {
					systemMenuInfo := GetSystemMenuInfo(systemMenuObject.ObjectInfo)
					if systemMenuInfo.MenuId != "" {
						listOfAllowedMenus = append(listOfAllowedMenus, systemMenuInfo.MenuId)
					}
				}

			}
		}
	} else {
		as.BaseService.Logger.Error("error getting user group data from database", zap.String("error", err.Error()))
	}

	return listOfAllowedMenus
}
func (as *AuthService) loadUserAccess() {
	err, listOfGroups := GetObjects(as.BaseService.ReferenceDatabase, UserGroupTable)

	if err != nil {
		as.BaseService.Logger.Error("unable to load the group information", zap.String("error", err.Error()))
		os.Exit(0)
	}

	var userMap = make(map[int][]int, 0)

	for _, object := range listOfGroups {
		userGroup := UserGroup{
			Id:         object.Id,
			ObjectInfo: object.ObjectInfo,
		}
		userList := userGroup.GetGroupInfo().Users
		rolesList := userGroup.GetGroupInfo().Roles
		for _, userId := range userList {
			for _, roleId := range rolesList {
				userMap[userId] = append(userMap[userId], roleId)
			}
		}
	}
	for key, value := range userMap {
		as.BaseService.Logger.Info("user map keys", zap.Any("key", value))
		userMap[key] = util.RemoveDuplicateInt(value)

		as.BaseService.Logger.Info("user map keys", zap.Any("roles", userMap[key]))
	}
	fmt.Println("userMap: ,userMap:", userMap)
	var userResourceMap = make(map[int][]int, 0)
	for userId, listOfRoles := range userMap {
		for _, roleId := range listOfRoles {
			// get the role info
			_, generalObject := Get(as.BaseService.ReferenceDatabase, RoleTable, roleId)
			role := Role{ObjectInfo: generalObject.ObjectInfo}
			fmt.Println("role.getRoleInfo(): ", role.getRoleInfo())
			listOfPermission := role.getRoleInfo().PermissionResources
			fmt.Println("listOfPermission: ", listOfPermission)
			for _, permissionId := range listOfPermission {
				_, permissionObject := Get(as.BaseService.ReferenceDatabase, PermissionTable, permissionId)
				permission := Permission{ObjectInfo: permissionObject.ObjectInfo}
				listOfPermissionResource := permission.getPermissionInfo().PermissionResources
				for _, resourceId := range listOfPermissionResource {
					userResourceMap[userId] = append(userResourceMap[userId], resourceId)
				}
			}
		}
	}

	for key, value := range userResourceMap {
		as.BaseService.Logger.Info("userResourceMap", zap.Any("key", value))
		userResourceMap[key] = util.RemoveDuplicateInt(value)
		as.BaseService.Logger.Info("userResourceMap, RemoveDuplicateInt", zap.Any("roles", userResourceMap[key]))
	}

	for userId, resourceIdList := range userResourceMap {
		for _, resourceId := range resourceIdList {
			_, generalObject := Get(as.BaseService.ReferenceDatabase, ComponentResourceTable, resourceId)
			componentResource := ComponentResource{ObjectInfo: generalObject.ObjectInfo}
			as.PermissionCache[userId] = append(as.PermissionCache[userId], componentResource.GetComponentResourceInfo())
			fmt.Println(" componentResource.GetComponentResourceInfo(): ", componentResource.GetComponentResourceInfo())
			as.BaseService.Logger.Info("user permission", zap.Any("user_id", userId), zap.Any("component_info", componentResource.GetComponentResourceInfo()))
		}
	}

	var listOfModules []common.SystemModuleInfo
	err, listOfModulesObjects := GetObjects(as.BaseService.ServiceDatabases[ProjectID], "system_module")
	if err == nil {
		for _, moduleObject := range listOfModulesObjects {
			systemModuleInfo := GetSystemModuleInfo(moduleObject.ObjectInfo)
			listOfModules = append(listOfModules, common.SystemModuleInfo{Id: moduleObject.Id, DisplayName: systemModuleInfo.DisplayName, Description: systemModuleInfo.Description, ModuleName: systemModuleInfo.Name})
		}
	}

	for _, moduleInfo := range listOfModules {
		as.ModuleCache[moduleInfo.ModuleName] = moduleInfo.Id
	}

}

func (as *AuthService) IsAllowed(userId int, projectId, moduleName, componentName string, action string, resourceId string, method string, path string) (bool, bool, string) {
	as.BaseService.Logger.Info("validation info", zap.Any("user_id", userId), zap.Any("component_name", componentName), zap.Any("resource_access_point", action), zap.Any("method", method))
	listOfResources := as.PermissionCache[userId]
	var isRouteEnabled bool
	var routingComponent string
	isRouteEnabled = false
	routingComponent = "*"
	for _, resourceObject := range listOfResources {
		as.BaseService.Logger.Info("validation info", zap.Any("resource_object", resourceObject.Resource))
		moduleId := as.ModuleCache[moduleName]
		if moduleId == resourceObject.ModuleId {
			// check the action
			if resourceObject.Method == method {
				if resourceObject.Resource == componentName {
					// now check the action
					if resourceObject.Action == action {
						if resourceObject.ResourceId == "" {
							// here not resource configured, so match the parth and return
							urlPath := resourceObject.Pattern
							urlPath = strings.Replace(urlPath, ":projectId", projectId, -1)
							if resourceObject.Resource != "" {
								urlPath = strings.Replace(urlPath, ":componentName", resourceObject.Resource, -1)
							}
							if resourceObject.Action != "" {
								urlPath = strings.Replace(urlPath, ":actionName", resourceObject.Action, -1)
							}

							// check new record is requested

							if urlPath == path {
								return true, resourceObject.IsRouteEnabled, resourceObject.RoutingComponent
							}
						} else if resourceObject.ResourceId == "*" {
							// all the resources
							urlPath := resourceObject.Pattern
							urlPath = strings.Replace(urlPath, ":projectId", projectId, -1)
							if resourceObject.Resource != "" {
								urlPath = strings.Replace(urlPath, ":componentName", resourceObject.Resource, -1)
							}
							if resourceObject.Action != "" {
								urlPath = strings.Replace(urlPath, ":actionName", resourceObject.Action, -1)
							}
							urlPath = strings.Replace(urlPath, ":recordId", "([0-9]+)", -1)
							match, _ := regexp.MatchString(urlPath, path)
							if match {
								return true, resourceObject.IsRouteEnabled, resourceObject.RoutingComponent
							}
						} else {
							urlPath := resourceObject.Pattern
							urlPath = strings.Replace(urlPath, ":projectId", projectId, -1)
							if resourceObject.Resource != "" {
								urlPath = strings.Replace(urlPath, ":componentName", resourceObject.Resource, -1)
							}
							if resourceObject.Action != "" {
								urlPath = strings.Replace(urlPath, ":actionName", resourceObject.Action, -1)
							}

							urlPath = strings.Replace(urlPath, ":recordId", resourceObject.ResourceId, -1)
							if urlPath == path {
								return true, resourceObject.IsRouteEnabled, resourceObject.RoutingComponent
							}
						}

					}
				}
			}

		}

	}

	_, listOfComponentResourcesObject := GetObjects(as.BaseService.ReferenceDatabase, ComponentResourceTable)
	for _, componentResourceObject := range listOfComponentResourcesObject {
		componentResource := ComponentResource{ObjectInfo: componentResourceObject.ObjectInfo}
		componentResourceInfo := componentResource.GetComponentResourceInfo()
		urlPath := componentResourceInfo.Pattern
		urlPath = strings.Replace(urlPath, ":projectId", componentResourceInfo.ProjectId, -1)

		//fmt.Println("default routing : ", urlPath, " path :", path)
		if urlPath == path {
			return false, componentResourceInfo.IsRouteEnabled, componentResourceInfo.RoutingComponent
		}
	}

	return false, isRouteEnabled, routingComponent
}
