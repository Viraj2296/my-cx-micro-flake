package analytics

import (
	"gorm.io/datatypes"
)

type Filter struct {
	FilterId        string   `json:"filterId"`
	FilterLabel     string   `json:"filterLabel"`
	MultiSelect     bool     `json:"multiSelect"`
	QuotationFormat string   `json:"quotationFormat"`
	WidgetId        int      `json:"widgetId"`
	FilterValues    []string `json:"filterValues"`
	Parameter       string   `json:"parameter"`
	FilterType      string   `json:"filterType"`
	DefaultValues   string   `json:"defaultValues"`
}

type WidgetInfo struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Status            string                 `json:"status"`
	Query             string                 `json:"query"`
	Filters           []Filter               `json:"filters"`
	Visualization     map[string]interface{} `json:"visualization"`
	SeriesMapping     datatypes.JSON         `json:"seriesMapping"`
	Tags              []string               `json:"tag"`
	Created           int                    `json:"created"`
	CreatedAt         string                 `json:"createdAt"`
	LastUpdatedAt     string                 `json:"lastUpdatedAt"`
	LastUpdatedBy     int                    `json:"lastUpdatedBy"`
	Datasource        int                    `json:"datasource"`
	VisualizationType string                 `json:"visualizationType"`
}

func (wi *WidgetInfo) GetFilter(filterId string) *Filter {
	for _, filter := range wi.Filters {
		if filter.FilterId == filterId {
			return &filter
		}
	}
	return nil
}

func (wi *WidgetInfo) GetVisualisationType() string {
	if visualisationType, ok := wi.Visualization["type"]; ok {
		return visualisationType.(string)
	}
	return ""
}

func (wi *WidgetInfo) GetDashboardFilter(filter *Filter) *Filter {
	dashboardFilter := Filter{}
	if filter.FilterType != "" {
		dashboardFilter.FilterType = filter.FilterType
	}
	if len(filter.FilterValues) > 0 {
		dashboardFilter.FilterValues = filter.FilterValues
	}
	if filter.QuotationFormat != "" {
		dashboardFilter.QuotationFormat = filter.QuotationFormat
	}
	dashboardFilter.MultiSelect = filter.MultiSelect
	if filter.Parameter != "" {
		dashboardFilter.Parameter = filter.Parameter
	}

	dashboardFilter.WidgetId = filter.WidgetId
	return &dashboardFilter
}
