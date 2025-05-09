package common

import (
	"cx-micro-flake/pkg/orm"
	"encoding/json"
	"gorm.io/datatypes"
)

type ModuleConfiguration struct {
	ModuleId            string         `json:"moduleId"`
	ModuleConfiguration datatypes.JSON `json:"moduleConfiguration"`
}

func (pm *ModuleConfiguration) Key() string {
	return "module_id='" + pm.ModuleId + "'"
}

func (pm *ModuleConfiguration) GetModuleConfiguration() *ModuleConfiguration {
	moduleConfig := ModuleConfiguration{}
	json.Unmarshal(pm.ModuleConfiguration, &moduleConfig)
	return &moduleConfig
}

type Config struct {
	Info struct {
		Icon        string `json:"icon"`
		Name        string `json:"name"`
		Author      string `json:"author"`
		Version     string `json:"version"`
		Description string `json:"description"`
		DisplayName string `json:"displayName"`
		RouteLink   string `json:"routeLink"`
		InstalledBy string `json:"installedBy"`
		InstalledOn string `json:"installedOn"`
	} `json:"info"`
	Menu []struct {
		Label      string `json:"label"`
		RouterLink string `json:"routerLink,omitempty"`
		Component  struct {
			DataSource struct {
				Get string `json:"get"`
			} `json:"dataSource"`
			ComponentType string `json:"componentType"`
			ComponentId   string `json:"componentId"`
			Config        struct {
				DashboardId    string `json:"dashboardId"`
				RefreshEnabled bool   `json:"refreshEnabled"`
				RefreshTime    string `json:"refreshTime"`
			} `json:"config"`
		} `json:"component,omitempty"`
		Items []struct {
			Label      string `json:"label"`
			Icon       string `json:"icon"`
			RouterLink string `json:"routerLink"`
			Component  struct {
				DataSource struct {
					Update string `json:"update"`
					Get    string `json:"get"`
					Post   string `json:"post"`
					Delete string `json:"delete"`
				} `json:"dataSource"`

				ComponentType string `json:"componentType"`
				ComponentId   string `json:"componentId"`
				Config        struct {
					Export bool `json:"export"`
					Import bool `json:"import"`
					Search bool `json:"search"`
				} `json:"config"`
			} `json:"component"`
		} `json:"items,omitempty"`
	} `json:"menu"`
}

type LinkedValues struct {
	FieldMapping map[string]map[int]string
}
type ProjectDatasourceConfig struct {
	ProjectId        string
	DatasourceConfig orm.DatabaseConfig
}
