package handler

import (
	"cx-micro-flake/pkg/common/analytics"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type MysqlConnectionParameters struct {
	HostName string `json:"hostName"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Schema   string `json:"schema"`
	Username string `json:"username"`
}

type CSVFileConnectionParameters struct {
	URL string `json:"url"`
}

type AnalyticsRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type AnalyticsComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AnalyticsWidget struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AnalyticsDatasourcesMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type DatasourceMasterInfo struct {
	Type               string    `json:"type"`
	Category           string    `json:"category"`
	CreatedAt          time.Time `json:"createdAt"`
	IsVisible          bool      `json:"isVisible"`
	DisplayName        string    `json:"displayName"`
	ObjectStatus       string    `json:"objectStatus"`
	LastUpdatedAt      time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy      int       `json:"lastUpdatedBy"`
	ConnectionField    string    `json:"connectionField"`
	LongDescription    string    `json:"longDescription"`
	ShortDescription   string    `json:"shortDescription"`
	DatasourceImageUrl string    `json:"datasourceImageUrl"`
}

type AnalyticsReport struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (v *AnalyticsDatasourcesMaster) getDatasourceMasterInfo() *DatasourceMasterInfo {
	datasourceMasterInfo := DatasourceMasterInfo{}
	json.Unmarshal(v.ObjectInfo, &datasourceMasterInfo)
	return &datasourceMasterInfo
}

func (v *AnalyticsWidget) getWidgetInfo() *analytics.WidgetInfo {
	widgetInfo := analytics.WidgetInfo{}
	json.Unmarshal(v.ObjectInfo, &widgetInfo)
	return &widgetInfo
}

type AnalyticsDatasource struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type DatasourcePermission struct {
	User   int `json:"user"`
	Access []struct {
		Id     int      `json:"id"`
		Table  string   `json:"table"`
		Fields []string `json:"fields"`
	}
}
type DatasourceInfo struct {
	Name                    string                 `json:"name"`
	Image                   string                 `json:"image"`
	Status                  string                 `json:"status"`
	CreatedAt               time.Time              `json:"createdAt"`
	CreatedBy               int                    `json:"createdBy"`
	Description             string                 `json:"description"`
	Permissions             []DatasourcePermission `json:"permissions"`
	ObjectStatus            string                 `json:"objectStatus"`
	ConnectionParam         interface{}            `json:"connectionParam"`
	DatasourceMaster        int                    `json:"datasourceMaster"`
	ConnectedDatabaseTables []string               `json:"connectedDatabaseTables"`
}

func (v *AnalyticsDatasource) getDatasourceInfo() *DatasourceInfo {
	datasourceInfo := DatasourceInfo{}
	json.Unmarshal(v.ObjectInfo, &datasourceInfo)
	return &datasourceInfo
}

type AnalyticsDashboard struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type SpcStatDatasource struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type SpcResourceDatasource struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AnalyticsDashboardPermission struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type RefreshExportRequest struct {
	SelectedId []string `json:"selectEdId"`
}

type AnalyticsReportInfo struct {
	Title                           string   `json:"title"`
	Format                          string   `json:"format"`
	Message                         string   `json:"message"`
	Subject                         string   `json:"subject"`
	Interval                        string   `json:"interval"`
	WidgetId                        *int     `json:"widgetId"`
	CreatedAt                       string   `json:"createdAt"`
	CreatedBy                       int      `json:"createdBy"`
	Recipient                       []int    `json:"recipient"`
	DaysOfWeek                      []string `json:"daysOfWeek"`
	DashboardId                     *int     `json:"dashboardId"`
	ObjectStatus                    string   `json:"objectStatus"`
	LastUpdatedAt                   string   `json:"lastUpdatedAt"`
	LastUpdatedBy                   int      `json:"lastUpdatedBy"`
	IsDashboardSource               bool     `json:"isDashboardSource"`
	ScheduledDateTime               *string  `json:"scheduledDateTime"`
	ScheduledMonthlyDate            *string  `json:"scheduledMonthlyDate"`
	ScheduledMonthlyDateSpecificDay *int     `json:"scheduledMonthlyDateSpecificDay"`
}
