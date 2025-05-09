package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"strings"
)

type PermissionResource struct {
	Action            string `json:"action"`
	Method            string `json:"method"`
	Pattern           string `json:"pattern"`
	ModuleId          int    `json:"moduleId"`
	Resource          string `json:"resource"`
	CreatedAt         string `json:"createdAt"`
	CreatedBy         int    `json:"createdBy"`
	ProjectId         string `json:"projectId"`
	ResourceId        string `json:"resourceId"`
	LastUpdatedAt     string `json:"lastUpdatedAt"`
	LastUpdatedBy     int    `json:"lastUpdatedBy"`
	IsRouteEnabled    bool   `json:"isRouteEnabled"`
	ComponentAction   string `json:"componentAction"`
	ResourceDisplay   string `json:"resourceDisplay"`
	RoutingComponent  string `json:"routingComponent"`
	ActionDescription string `json:"actionDescription"`
}

func (v *PermissionResource) updateToDB() {
	serialisedPermission, _ := json.Marshal(v)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	authService.UpdateComponentResource(serialisedPermission)
}

func getNewRecordPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "GET"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/new_record"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = ""
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Get " + permissionResource.ResourceDisplay + " New Record"
	permissionResource.ActionDescription = "Get " + permissionResource.ResourceDisplay + " New Record"

	return permissionResource
}

func getHMICardViewPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "GET"
	permissionResource.Pattern = "/project/:projectId/machines/hmi_card_view"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = ""
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	permissionResource.ResourceDisplay = "Machines"
	permissionResource.ComponentAction = "Get HMI Card View"
	permissionResource.ActionDescription = "Get HMI Card View"

	return permissionResource
}

func getRecordPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "GET"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/record/:recordId"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = "*"
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Get " + permissionResource.ResourceDisplay + " Record"
	permissionResource.ActionDescription = "Get " + permissionResource.ResourceDisplay + " Record"

	return permissionResource
}

func updateRecordPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "PUT"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/record/:recordId"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = "*"
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Update " + permissionResource.ResourceDisplay + " Record"
	permissionResource.ActionDescription = "Update " + permissionResource.ResourceDisplay + " Record"

	return permissionResource
}

func deleteRecordPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "DELETE"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/record/:recordId"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = "*"
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Delete " + permissionResource.ResourceDisplay + " Record"
	permissionResource.ActionDescription = "Delete " + permissionResource.ResourceDisplay + " Record"

	return permissionResource
}

func createRecordPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "POST"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/records"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = ""
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Create " + permissionResource.ResourceDisplay + " Record"
	permissionResource.ActionDescription = "Create " + permissionResource.ResourceDisplay + " Record"

	return permissionResource
}

func getReadAllPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "GET"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/records"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = ""
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Read All"
	permissionResource.ActionDescription = "Read All"

	return permissionResource
}

func getRecordMessageTrailPermission(componentName string) PermissionResource {
	permissionResource := PermissionResource{}
	permissionResource.Resource = componentName

	// GET New Record
	permissionResource.Action = ""
	permissionResource.Method = "GET"
	permissionResource.Pattern = "/project/:projectId/machines/component/:componentName/record_messages/:recordId"
	permissionResource.ModuleId = 1
	permissionResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	permissionResource.CreatedBy = 1
	permissionResource.LastUpdatedBy = 1
	permissionResource.ResourceId = "*"
	permissionResource.ProjectId = "906d0fd569404c59956503985b330132"
	permissionResource.IsRouteEnabled = false
	permissionResource.RoutingComponent = "*"

	modifiedComponentName := strings.Replace(componentName, "_", " ", -1)
	permissionResource.ResourceDisplay = strings.Title(modifiedComponentName)
	permissionResource.ComponentAction = "Get " + permissionResource.ResourceDisplay + " Record Trail"
	permissionResource.ActionDescription = "Delete " + permissionResource.ResourceDisplay + " Record Trail"

	return permissionResource
}
func (v *MachineService) generatePermissions() {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfComponents, _ := GetObjects(dbConnection, MachineComponentTable)
	permissionResource := getHMICardViewPermission("")
	permissionResource.updateToDB()
	for _, componentObject := range *listOfComponents {
		var componentFields = make(map[string]interface{})
		json.Unmarshal(componentObject.ObjectInfo, &componentFields)
		componentName := componentFields["componentName"].(string)
		permissionResource = getNewRecordPermission(componentName)
		permissionResource.updateToDB()
		permissionResource = getRecordPermission(componentName)
		permissionResource.updateToDB()
		permissionResource = updateRecordPermission(componentName)
		permissionResource.updateToDB()
		permissionResource = getReadAllPermission(componentName)
		permissionResource.updateToDB()
		permissionResource = createRecordPermission(componentName)
		permissionResource.updateToDB()
		permissionResource = deleteRecordPermission(componentName)
		permissionResource.updateToDB()
		permissionResource = getRecordMessageTrailPermission(componentName)
		permissionResource.updateToDB()

	}

}
