package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	AnalyticsRecordTrailTable = "analytics_record_trail"

	AnalyticsComponentTable           = "analytics_component"
	AnalyticsWidgetTable              = "analytics_widget"
	AnalyticsDataSourceTable          = "analytics_datasource"
	AnalyticsDashboardTable           = "analytics_dashboard"
	AnalyticsDatasourcesMasterTable   = "analytics_datasources_master"
	AnalyticsReportTable              = "analytics_report"
	AnalyticsDashboardPermissionTable = "analytics_dashboard_permission"

	SPCStatDataSourceTable     = "spc_stat_datasource"
	SPCResourceDataSourceTable = "spc_resource_datasource"

	AnalyticsDataSourceComponent = "datasource"
	AnalyticsWidgetComponent     = "widget"

	InvalidSourceError = "Invalid Source"

	InvalidSchedulePosition = "Invalid Schedule Position"
	InvalidComponent        = 6010

	InvalidScheduleStatus    = 6054
	ErrorGettingActionFields = 6056
	QueryExecutionFailed     = 6057
	InvalidEventId           = 6058
	InvalidMachineId         = 6059

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008
	InvalidDatasourceType                   = 5009

	ConnectingDatasourceFailed = 5009
	DecodingFailed             = 6070

	ProjectID  = "906d0fd569404c59956503985b330132"
	TimeLayout = "2006-01-02T15:04:05.000Z"

	FunctionGetExistingDatabaseTables = "getExistingDatabaseTables"

	ModuleName = "analytics"

	ReadPermission    = 1
	WritePermission   = 2
	ExecutePermission = 3
)

type InvalidRequest struct {
	Message string `json:"message"`
}

func getError(errorString string) error {
	return errors.New(errorString)
}

func sendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}
func sendArchiveFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Archived Failed",
			Description: "Internal system error during archive process. This is normally happen when the system is not configured properly. Please report to system administrator",
		})
	return
}
