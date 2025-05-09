package handler

import (
	"cx-micro-flake/pkg/common/analytics"
	"cx-micro-flake/pkg/util"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (as *AnalyticsService) OnTimer() {
	//c := cron.New()
	//c.AddFunc("@every 20s", func() { as.loadDatasourceConnection() })
	//c.Start()
}

// This is for purpose to check the initial and progressive
func (as *AnalyticsService) loadDatasourceConnection() {
	dbConnection := as.BaseService.ServiceDatabases[ProjectID]
	listOfDatasource, err := GetObjects(dbConnection, AnalyticsDataSourceTable)
	fmt.Println("loaded_listOfDatasources", listOfDatasource)
	if err == nil {
		for _, datasourceInterface := range *listOfDatasource {
			analyticsDatasource := AnalyticsDatasource{ObjectInfo: datasourceInterface.ObjectInfo}
			if getDatasourceType(dbConnection, analyticsDatasource.getDatasourceInfo().DatasourceMaster) == "mysql" {
				serialisedConnectionParam, _ := json.Marshal(analyticsDatasource.getDatasourceInfo().ConnectionParam)
				mysqlConnectionParam := MysqlConnectionParameters{}
				json.Unmarshal(serialisedConnectionParam, &mysqlConnectionParam)
				fmt.Println("mysqlConnectionParam: ", mysqlConnectionParam)
				if _, ok := as.DatasourceConnectionCache[datasourceInterface.Id]; !ok {
					// not created the connection, we will establish again
					fmt.Println("DatasourceConnectionCache:", as.DatasourceConnectionCache)
					dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConnectionParam.Username, mysqlConnectionParam.Password, mysqlConnectionParam.HostName, mysqlConnectionParam.Port, mysqlConnectionParam.Schema)
					datasourceConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
					if err == nil {
						//no error, so update the cache
						as.BaseService.Logger.Info("datasource connection established, type = MySQL")
						as.DatasourceConnectionCache[datasourceInterface.Id] = datasourceConnection
					}
				}
			}
		}
	}
}

// TODO, we need to have another cron to check the connection status, and update the status

// Scheduler for send analytics reports
func sendAnalyticsReport(as *AnalyticsService, dbConnection *gorm.DB) {
	analyticsReportObjects, err := GetObjects(dbConnection, AnalyticsReportTable)

	if err != nil {
		return
	}

	for _, analyticReport := range *analyticsReportObjects {
		currentTimeNow := time.Now().UTC()
		currentYear, currentMonth, _ := currentTimeNow.Date()
		currentLocation := currentTimeNow.Location()

		analyticsReportInfo := AnalyticsReportInfo{}
		json.Unmarshal(analyticReport.ObjectInfo, &analyticsReportInfo)
		lastExecutionTime := util.ConvertStringToDateTime(analyticsReportInfo.LastUpdatedAt)

		switch {
		case analyticsReportInfo.Interval == "daily":
			if currentTimeNow.Unix() > lastExecutionTime.DateTimeEpoch {
				//send file
				if analyticsReportInfo.IsDashboardSource {
					// err , dashboardObject := Get(dbConnection, AnalyticsDashboardTable, *analyticsReportInfo.DashboardId)
					// if err != nil{
					// 	continue
					// }

				} else {
					_, widgetObject := Get(dbConnection, AnalyticsWidgetTable, *analyticsReportInfo.WidgetId)
					widgetInfo := analytics.WidgetInfo{}
					json.Unmarshal(widgetObject.ObjectInfo, &widgetInfo)
					_, widgetData := as.generateChartData(widgetInfo.GetVisualisationType(), dbConnection, widgetObject.Id, nil)
					dataBuffer := createTableDataBuffer(widgetData)

					switch {
					case analyticsReportInfo.Format == "csv":
						createCSVFile(as, dataBuffer)

					case analyticsReportInfo.Format == "excel":
						createExcelFile(as, dataBuffer)
					}

				}
			}

		case analyticsReportInfo.Interval == "weekly":
			if currentTimeNow.Unix() > lastExecutionTime.DateTimeEpoch {
				//send file
			}

		case analyticsReportInfo.Interval == "monthly":
			firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
			lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
			fmt.Println(lastOfMonth)
		}
	}
}

func createCSVFile(as *AnalyticsService, dataBuffer [][]string) {
	csvFile, err := os.Create("timeline.csv")

	if err != nil {
		as.BaseService.Logger.Error("Can't create csv file", zap.String("error", err.Error()))
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, empRow := range dataBuffer {
		_ = csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	csvFile.Close()

}

func createExcelFile(as *AnalyticsService, dataBuffer [][]string) {
	if len(dataBuffer) == 0 {
		return
	}
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	row := len(dataBuffer)
	col := len(dataBuffer[0])
	// Set value of a cell.

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			cellIndex, _ := excelize.CoordinatesToCellName(i, j)
			f.SetCellValue("Sheet1", cellIndex, dataBuffer[i][j])
		}

	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Timeline.xlsx"); err != nil {
		fmt.Println(err)
	}

}

func createTableDataBuffer(visualizationObject interface{}) [][]string {
	// map[chart:map[] credits:map[enabled:false] legend:map[enabled:true] plotOptions:map[series:map[turboThreshold:0]]
	// series:[map[data:[[2 2.2] [3 2.6] [4 1.19] [5 1.5] [6 1.2] [7 1.3] [8 1.9]]
	// id:fd0e41c07c06afeb48f2f0038d56663f name:energy type:spline]]
	// title:map[text: zoomType:xy] tooltip:map[shared:true] type:chart xAxis:map[title:map[enabled:true] type:auto]
	// yAxis:[map[gridLineWidth:1 index:0 labels:map[] opposite:false title:map[text:energy] type:auto]]]

	csvDataPayload := make([][]string, 0)
	visualization := visualizationObject.(map[string]interface{})
	series := visualization["series"].(map[string]interface{})
	seriesData := series["data"].([]map[string]interface{})

	switch {
	case util.InterfaceToString(series["type"]) == "timeline":
		//Get keys of series data
		keys := make([]string, 0)
		if len(seriesData) != 0 {
			for key, _ := range seriesData[0] {
				keys = append(keys, key)
			}
			csvDataPayload = append(csvDataPayload, keys)
		}

		//Inserting series data
		for _, timeLineData := range seriesData {
			timeLineRowData := make([]interface{}, 0)
			for _, timeLineKey := range keys {
				timeLineRowData = append(timeLineRowData, timeLineData[timeLineKey].(string))
			}
		}
	}
	return csvDataPayload
}
