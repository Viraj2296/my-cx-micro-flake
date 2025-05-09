package dto

type LineChartResponse struct {
	Chart struct {
		Type string `json:"type"`
	} `json:"chart"`
	Title struct {
		Text string `json:"text"`
	} `json:"title"`
	XAxis struct {
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Categories []string `json:"categories"`
	} `json:"xAxis"`
	YAxis struct {
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Min int `json:"min"`
	} `json:"yAxis"`
	Credits struct {
		Enabled bool `json:"enabled"`
	} `json:"credits"`
	Series []Series `json:"series"`
}

type Series struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
	Type string    `json:"type"`
}
