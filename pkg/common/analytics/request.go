package analytics

type QueryResponseRequest struct {
	Format     string                 `json:"format"`
	Param      map[string]interface{} `json:"param"`
	Query      string                 `json:"query"`
	GroupByCol string                 `json:"groupByCol"`
}

type ChartFormatRequest struct {
	XColumn    string   `json:"xColumn"`
	YColumns   []string `json:"yColumns"`
	GroupByCol string   `json:"groupByCol"`
}

type TableFormatRequest struct {
	TableColumns []string `json:"tableColumns"`
}

type TimelineFormatRequest struct {
	TimelineColumn    string `json:"timelineColumn"`
	NameColumn        string `json:"nameColumn"`
	LabelColumn       string `json:"labelColumn"`
	DescriptionColumn string `json:"descriptionColumn"`
}
