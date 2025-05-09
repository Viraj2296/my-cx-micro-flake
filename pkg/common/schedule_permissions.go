package common

type FieldPermission struct {
	CanUpdate bool `json:"canUpdate"`
}

type ActionPermission struct {
	CanComplete  bool `json:"canComplete"`
	CanForceStop bool `json:"canForceStop"`
	CanHold      bool `json:"canHold"`
	CanRelease   bool `json:"canRelease"`
	CanUnRelease bool `json:"canUnRelease"`
}

type Permissions struct {
	Fields  map[string]FieldPermission `json:"fields"`
	Actions ActionPermission           `json:"actions"`
}

func GetPermissions(status int, canComplete bool, canForceStop bool, canHold bool, canRelease bool) Permissions {
	permissions := Permissions{
		Fields: map[string]FieldPermission{
			"plannedManpower": {CanUpdate: true},
			"startDate":       {CanUpdate: true},
			"endDate":         {CanUpdate: true},
			"name":            {CanUpdate: true},
			"partNo":          {CanUpdate: true},
			"productionOrder": {CanUpdate: true},
			"completedQty":    {CanUpdate: true},
			"scheduledQty":    {CanUpdate: true},
		},
		Actions: ActionPermission{
			CanComplete:  canComplete,
			CanForceStop: canForceStop,
			CanHold:      canHold,
			CanRelease:   canRelease,
		},
	}

	if status >= 4 {
		for key := range permissions.Fields {
			permissions.Fields[key] = FieldPermission{CanUpdate: false}
		}
	}

	return permissions
}
func RemoveKeys(data map[string]interface{}, keysToRemove []string) map[string]interface{} {
	for _, key := range keysToRemove {
		delete(data, key)
	}
	return data
}

func GetPermissionsForMaintenance(status int, canComplete bool, canForceStop bool, canHold bool, canRelease bool, canUnRelease bool) Permissions {
	permissions := Permissions{
		Fields: map[string]FieldPermission{
			"startDate": {CanUpdate: true},
			"endDate":   {CanUpdate: true},
			"name":      {CanUpdate: true},
		},
		Actions: ActionPermission{
			CanComplete:  canComplete,
			CanForceStop: canForceStop,
			CanHold:      canHold,
			CanRelease:   canRelease,
			CanUnRelease: canUnRelease,
		},
	}

	if status >= 4 {
		for key := range permissions.Fields {
			permissions.Fields[key] = FieldPermission{CanUpdate: false}
		}
	}

	return permissions
}

func GetPermissionsForAssembly(status int, canComplete bool, canForceStop bool, canHold bool, canRelease bool) Permissions {
	permissions := Permissions{
		Fields: map[string]FieldPermission{
			"plannedManpower": {CanUpdate: true},
			"startDate":       {CanUpdate: true},
			"endDate":         {CanUpdate: true},
			"name":            {CanUpdate: true},
			"partNo":          {CanUpdate: true},
			"productionOrder": {CanUpdate: true},
			"module":          {CanUpdate: true},
			"scheduledQty":    {CanUpdate: true},
			"completedQty":    {CanUpdate: true},
			"priorityLevel":   {CanUpdate: true},
		},
		Actions: ActionPermission{
			CanComplete:  canComplete,
			CanForceStop: canForceStop,
			CanHold:      canHold,
			CanRelease:   canRelease,
		},
	}

	if status >= 4 {
		for key := range permissions.Fields {
			permissions.Fields[key] = FieldPermission{CanUpdate: false}
		}
	}

	return permissions
}
