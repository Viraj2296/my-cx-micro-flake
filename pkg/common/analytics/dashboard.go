package analytics

type DashboardWidgetFilterInfo struct {
	FilterId      string `json:"filterId"`
	ValueSource   string `json:"valueSource"`
	StaticValue   string `json:"staticValue"`
	Title         string `json:"title"`
	DefaultValues string `json:"defaultValues"`
}
type TempDashboardFilter struct {
	FilterId        string   `json:"filterId"`
	WidgetId        int      `json:"widgetId"`
	FilterType      string   `json:"filterType"`  // it says, dropped_down, number, text, date,date-range
	FilterLabel     string   `json:"filterLabel"` // Surface Types
	MultiSelect     bool     `json:"multiSelect"` // whether this multi select or not
	FilterValues    []string `json:"filterValues"`
	Param           string   `json:"param"`
	QuotationFormat string   `json:"quotationFormat"`
	DefaultValues   string   `json:"defaultValues"`
}
type DashboardFilter struct {
	FilterId        string   `json:"filterId"`
	FilterType      string   `json:"filterType"`  // it says, dropped_down, number, text, date,date-range
	FilterLabel     string   `json:"filterLabel"` // Surface Types
	MultiSelect     bool     `json:"multiSelect"` // whether this multi select or not
	FilterValues    []string `json:"filterValues"`
	Param           string   `json:"param"`
	DefaultValues   string   `json:"defaultValues"`
	QuotationFormat string   `json:"quotationFormat"`
	LinkedWidgets   []int    `json:"linkedWidgets"`
}

// this will be used to store the dashboard widget, we don;t need to store the visualisation object in the dashboard itself, it should be generated to added
type DashboardWidget struct {
	Id                        int                         `json:"id"`
	X                         int                         `json:"x"`
	Y                         int                         `json:"y"`
	Cols                      int                         `json:"cols"`
	Rows                      int                         `json:"rows"`
	Width                     float64                     `json:"width"`
	Height                    float64                     `json:"height"`
	WidgetId                  int                         `json:"widgetId"` // this will be used to send during update
	Name                      string                      `json:"name"`
	VisualizationType         string                      `json:"visualizationType"`
	DashboardWidgetFilterInfo []DashboardWidgetFilterInfo `json:"dashboardWidgetFilterInfo"` // when adding widget, we should specify what are the filters we are using for dashboard
}

type DashboardWidgetSnapshot struct {
	Id                int         `json:"id"`
	X                 int         `json:"x"`
	Y                 int         `json:"y"`
	Cols              int         `json:"cols"`
	Rows              int         `json:"rows"`
	Width             float64     `json:"width"`
	Height            float64     `json:"height"`
	WidgetId          int         `json:"widgetId"` // this will be used to send during update
	Name              string      `json:"name"`
	Visualization     interface{} `json:"visualization"`
	WidgetFilterInfo  []*Filter   `json:"widgetFilterInfo"`
	VisualizationType string      `json:"visualizationType"`
}

type DashboardInfo struct {
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Status            string            `json:"status"`
	RefreshInterval   string            `json:"refreshInterval"`
	AllowPublicAccess bool              `json:"allowPublicAccess"`
	PublicUrl         string            `json:"publicUrl"`
	Created           int               `json:"created"`
	CreatedAt         string            `json:"createdAt"`
	LastUpdatedAt     string            `json:"lastUpdatedAt"`
	LastUpdatedBy     int               `json:"lastUpdatedBy"`
	DashboardWidgets  []DashboardWidget `json:"dashboardWidgets"`
	SnapshotEnabled   bool              `json:"snapshotEnabled"`
}

type DashboardWidgetSnapshotInfo struct {
	Name             string                    `json:"name"`
	RefreshInterval  string                    `json:"refreshInterval"`
	DashboardFilters []DashboardFilter         `json:"dashboardFilters"`
	DashboardWidgets []DashboardWidgetSnapshot `json:"dashboardWidgets"`
}

type DashboardWidgetSnapshotResponse struct {
	DashboardId   int                         `json:"dashboardId"`
	DashboardInfo DashboardWidgetSnapshotInfo `json:"dashboardInfo"`
}
