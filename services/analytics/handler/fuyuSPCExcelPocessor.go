package handler

//type FuYuSPCConnectionRequest struct {
//	ConnectionParam struct {
//		Url string `json:"url"`
//	} `json:"connectionParam"`
//	Name                string `json:"name"`
//	Description         string `json:"description"`
//	DatasourceAttribute struct {
//		TargetTable        string   `json:"targetTable"`
//		UniqueIndexColumns []string `json:"uniqueIndexColumns"`
//	} `json:"datasourceAttribute"`
//}
//
//type FuYuSPCDatasource struct {
//	Name            string        `json:"name"`
//	Status          string        `json:"status"`
//	CreatedAt       time.Time     `json:"createdAt"`
//	CreatedBy       int           `json:"createdBy"`
//	Description     string        `json:"description"`
//	Permissions     []interface{} `json:"permissions"`
//	ObjectStatus    string        `json:"objectStatus"`
//	LastUpdatedAt   time.Time     `json:"lastUpdatedAt"`
//	LastUpdatedBy   int           `json:"lastUpdatedBy"`
//	ConnectionParam struct {
//		Url string `json:"url"`
//	} `json:"connectionParam"`
//	DatasourceMaster    int `json:"datasourceMaster"`
//	DatasourceAttribute struct {
//		TargetTable        string   `json:"targetTable"`
//		UniqueIndexColumns []string `json:"uniqueIndexColumns"`
//	} `json:"datasourceAttribute"`
//	ConnectedDatabaseTables []string `json:"connectedDatabaseTables"`
//}
//
//func (as *AnalyticsService) handleFuYuSPCDatabaseRecordCreation(ctx *gin.Context, dbConnection *gorm.DB, request map[string]interface{}) error {
//
//	serialisedRequest, _ := json.Marshal(request)
//	fmt.Println("serialised request : ", string(serialisedRequest))
//	fuyuSPCConnectionRequest := FuYuSPCConnectionRequest{}
//	json.Unmarshal(serialisedRequest, &fuyuSPCConnectionRequest)
//	// create the table
//	fmt.Println("ExcelFileConnectionRequest request : ", fuyuSPCConnectionRequest)
//
//	// now download the file in /tmp folder
//	savedFileName := "/tmp/" + uuid.New().String()
//
//	client := &getter.Client{
//		Ctx: context.Background(),
//		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
//		Dst: savedFileName,
//		Dir: false,
//		//the repository with a subdirectory I would like to clone only
//		Src:  fuyuSPCConnectionRequest.ConnectionParam.Url,
//		Mode: getter.ClientModeFile,
//		////define the type of detectors go getter should use, in this case only github is needed
//
//		////provide the getter needed to download the files
//		Getters:  GetGetter("https"),
//		Insecure: false,
//	}
//	//download the files
//	if err := client.Get(); err != nil {
//		return errors.New("failed to process given file, check the file again")
//	}
//
//	dataSourceName := strings.Join(strings.Fields(fuyuSPCConnectionRequest.Name), "_")
//
//	targetTable1 := dataSourceName + "_table1"
//	targetTable2 := dataSourceName + "_table2"
//
//	// err, recordId, columnList := readStatsData(as, savedFileName, targetTable1, dbConnection)
//	err := readStatsData(as, savedFileName, targetTable1, dbConnection)
//
//	if err != nil {
//		//TODO remove table is created in the previous step
//		return err
//	}
//
//	err = readSourceData(as, savedFileName, dataSourceName, dbConnection)
//
//	if err != nil {
//		//TODO remove table is created in the previous step
//		// deleteObj := component.GeneralObject{Id: recordId}
//		// Delete(dbConnection, SPCStatDataSourceTable, deleteObj)
//		dbConnection.Exec("DROP TABLE " + targetTable1)
//		return err
//	}
//
//	err = addDashboardData(fuyuSPCConnectionRequest.Name, dbConnection, targetTable1, targetTable2)
//
//	// var connectedTables []string
//	// connectedTables = append(connectedTables, fuyuSPCConnectionRequest.DatasourceAttribute.TargetTable)
//	// request["connectedDatabaseTables"] = connectedTables
//
//	// import the data
//	// fileUrl := fuyuSPCConnectionRequest.ConnectionParam.Url
//	// csvConnection := CSVFileConnection{}
//	// _, err = csvConnection.ImportData(fileUrl, fuyuSPCConnectionRequest.DatasourceAttribute.TargetTable, fuyuSPCConnectionRequest.DatasourceAttribute.DestinationSchema, dbConnection)
//	// if err != nil {
//	// 	dbConnection.Exec("DROP TABLE " + fuyuSPCConnectionRequest.DatasourceAttribute.TargetTable)
//	// 	response.SendDetailedError(ctx, http.StatusBadRequest, getError("Failed to process data"), ErrorCreatingObjectInformation, err.Error())
//	// 	return err
//	// }
//	return err
//}
//
//func readStatsData(as *AnalyticsService, fileName string, targetTable string, dbConnection *gorm.DB) error {
//	f, err := excelize.OpenFile(fileName)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	defer func() {
//		// Close the spreadsheet.
//		if err := f.Close(); err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	columnList := []string{"stat_key"}
//	columnCount := 1
//	statKeys := make([]string, 0)
//	allData := make([][]string, 0)
//	for rowIndex := 2; rowIndex <= 19; rowIndex++ {
//		rowData := make([]string, 0)
//		for colIndex := 9; colIndex <= 31; colIndex++ {
//			if rowIndex == 19 && colIndex == 9 {
//				continue
//			}
//
//			cellIndex, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
//			cellValue, err := f.GetCellValue("sheet1", cellIndex)
//			if err != nil {
//				fmt.Println(err)
//				cellValue = ""
//			}
//
//			if rowIndex == 19 {
//				renameColumn := strings.Join(strings.Fields(cellValue), "_") + "_" + strconv.Itoa(columnCount)
//				columnList = append(columnList, renameColumn)
//				columnCount += 1
//			} else if colIndex == 9 {
//				if cellValue != "" {
//					statKeys = append(statKeys, cellValue)
//				}
//			} else {
//				rowData = append(rowData, cellValue)
//			}
//
//		}
//		if len(rowData) != 0 {
//			allData = append(allData, rowData)
//		}
//	}
//
//	// targetTable := fuyuSPCConnectionRequest.DatasourceAttribute.TargetTable + "_table1"
//	// check table is there created already
//	//CREATE TABLE `fuyu_mes`.`new_table` (`id` INT NOT NULL AUTO_INCREMENT,`student_marks` VARCHAR(45) NULL,PRIMARY KEY (`id`));
//	var tableExecutionString = "CREATE TABLE `" + targetTable + "` (`record_id` INT NOT NULL AUTO_INCREMENT,"
//
//	fmt.Println("columnList", columnList)
//	// Create column for table
//	insertColumns := " ("
//	for _, colName := range columnList {
//		tableExecutionString += "`" + colName + "` VARCHAR(512) NULL,"
//		insertColumns += "`" + colName + "`,"
//	}
//	insertColumns = util.TrimSuffix(insertColumns, ",") + ") "
//
//	tableExecutionString = tableExecutionString[:len(tableExecutionString)-1]
//	tableExecutionString += ",PRIMARY KEY (`record_id`));"
//	err = dbConnection.Exec(tableExecutionString).Error
//
//	as.BaseService.Logger.Infow("created table creation string", "execution_string", tableExecutionString)
//	if err != nil {
//		// response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Schema"), ErrorCreatingObjectInformation, err.Error())
//		return err
//	}
//
//	for index, allStatsData := range allData {
//		// after successfull table creation, we needed to compose the insert statment and write.
//		insertStatement := "INSERT INTO " + targetTable + insertColumns + " VALUES("
//		values := "'" + statKeys[index] + "',"
//		for _, rowStatsData := range allStatsData {
//			values += "'" + rowStatsData + "',"
//		}
//
//		insertStatement = insertStatement + util.TrimSuffix(values, ",") + ");"
//
//		fmt.Println(insertStatement)
//		err = dbConnection.Exec(insertStatement).Error
//
//		if err != nil {
//			dbConnection.Exec("DROP TABLE " + targetTable)
//			fmt.Println(err)
//			return err
//		}
//
//	}
//
//	return nil
//
//}
//
//// func readStatsData(as *AnalyticsService, fileName string, dataSource string, dbConnection *gorm.DB) (error, int, []string) {
//// 	f, err := excelize.OpenFile(fileName)
//// 	columnList := []string{"stat_key"}
//// 	if err != nil {
//// 		return err, 0, columnList
//// 	}
//// 	defer func() {
//// 		// Close the spreadsheet.
//// 		if err := f.Close(); err != nil {
//// 			fmt.Println(err)
//// 		}
//// 	}()
//
//// 	columnCount := 1
//// 	statKeys := make([]string, 0)
//// 	allData := make([][]string, 0)
//// 	for rowIndex := 2; rowIndex <= 19; rowIndex++ {
//// 		rowData := make([]string, 0)
//// 		for colIndex := 9; colIndex <= 31; colIndex++ {
//// 			if rowIndex == 19 && colIndex == 9 {
//// 				continue
//// 			}
//
//// 			cellIndex, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
//// 			cellValue, err := f.GetCellValue("sheet1", cellIndex)
//// 			if err != nil {
//// 				fmt.Println(err)
//// 				cellValue = ""
//// 			}
//
//// 			if rowIndex == 19 {
//// 				renameColumn := strings.Join(strings.Fields(cellValue), "_") + "_" + strconv.Itoa(columnCount)
//// 				columnList = append(columnList, renameColumn)
//// 				columnCount += 1
//// 			} else if colIndex == 9 {
//// 				if cellValue != "" {
//// 					statKeys = append(statKeys, cellValue)
//// 				}
//// 			} else {
//// 				rowData = append(rowData, cellValue)
//// 			}
//
//// 		}
//// 		if len(rowData) != 0 {
//// 			allData = append(allData, rowData)
//// 		}
//// 	}
//
//// 	//var createTableCommand = "CREATE TABLE `fuyu_mes`.`[TABLE_NAME]` (`id` INT NOT NULL AUTO_INCREMENT,`object_info` json DEFAULT NULL,PRIMARY KEY (`id`))"
//// 	//createTableCommand = strings.Replace(createTableCommand, "[TABLE_NAME]", targetTable, -1)
//// 	//err = dbConnection.Exec(createTableCommand).Error
//// 	//fmt.Println("creating a table :", createTableCommand)
//// 	//if err != nil {
//// 	//	fmt.Println(err)
//// 	//	return err, 0
//// 	//}
//// 	partNumberIndex, _ := excelize.CoordinatesToCellName(6, 14)
//// 	partNumber, _ := f.GetCellValue("sheet1", partNumberIndex)
//// 	insertData := make(map[string]interface{})
//// 	// insertStatement := "INSERT INTO " + SPCStatDataSourceTable + " (object_info)" + " VALUES('{\"name\":\"" + dataSource + "\",\"partNumber\":\"" + partNumber + "\",\"resourceData\":["
//// 	// insertStatement := "'{\"name\":\"" + dataSource + "\",\"partNumber\":\"" + partNumber + "\",\"resourceData\":["
//// 	insertData["name"] = dataSource
//// 	insertData["partNumber"] = partNumber
//// 	resourceDataList := make([]interface{}, 0)
//// 	for index, allStatsData := range allData {
//// 		// after successfull table creation, we needed to compose the insert statment and write.
//// 		resourceData := make(map[string]interface{})
//// 		resourceData["stat_key"] = statKeys[index]
//// 		// values := "{\"stat_key\":" + "\"" + statKeys[index] + "\","
//// 		for stat_index, rowStatsData := range allStatsData {
//// 			resourceData[columnList[stat_index+1]] = rowStatsData
//// 			// values += "\"" + columnList[stat_index+1] + "\":" + "\"" + rowStatsData + "\","
//// 		}
//// 		resourceDataList = append(resourceDataList, resourceData)
//
//// 	}
//// 	insertData["resourceData"] = resourceDataList
//
//// 	objectInfo, _ := json.Marshal(insertData)
//
//// 	insertObj := component.GeneralObject{ObjectInfo: objectInfo}
//
//// 	// fmt.Println(insertStatement)
//// 	err, recordId := Create(dbConnection, SPCStatDataSourceTable, insertObj)
//
//// 	if err != nil {
//// 		fmt.Println(err)
//// 		return err, 0, columnList
//// 	}
//
//// 	return nil, recordId, columnList
//
//// }
//
//func readSourceData(as *AnalyticsService, fileName string, targetTable string, dbConnection *gorm.DB) error {
//	f, err := excelize.OpenFile(fileName)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	defer func() {
//		// Close the spreadsheet.
//		if err := f.Close(); err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	columnList := make([]string, 0)
//	allData := make([][]string, 0)
//	columnCount := 1
//	for rowIndex := 21; rowIndex <= 51; rowIndex++ {
//		rowData := make([]string, 0)
//		for colIndex := 4; colIndex <= 34; colIndex++ {
//			cellIndex, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
//			cellValue, err := f.GetCellValue("sheet1", cellIndex)
//			if err != nil {
//				fmt.Println(err)
//				cellValue = ""
//			}
//
//			if rowIndex == 21 {
//				if cellValue == "Inspected by" {
//					columnList = append(columnList, "inspected_by")
//				} else if cellValue == "Remark" {
//					columnList = append(columnList, "remark")
//				} else if cellValue == "Working Shift" {
//					columnList = append(columnList, "working_shift")
//				} else {
//					renameColumn := strings.Join(strings.Fields(cellValue), "_") + "_" + strconv.Itoa(columnCount)
//					columnList = append(columnList, renameColumn)
//				}
//
//				columnCount += 1
//			} else {
//				rowData = append(rowData, cellValue)
//			}
//
//		}
//		if len(rowData) != 0 {
//			allData = append(allData, rowData)
//		}
//
//	}
//
//	// targetTable := fuyuSPCConnectionRequest.DatasourceAttribute.TargetTable + "_table2"
//	// check table is there created already
//	//CREATE TABLE `fuyu_mes`.`new_table` (`id` INT NOT NULL AUTO_INCREMENT,`student_marks` VARCHAR(45) NULL,PRIMARY KEY (`id`));
//	var tableExecutionString = "CREATE TABLE `" + targetTable + "` (`record_id` INT NOT NULL AUTO_INCREMENT,"
//
//	// Create column for table
//	insertColumns := " ("
//	for _, colName := range columnList {
//		tableExecutionString += "`" + colName + "` VARCHAR(512) NULL,"
//		insertColumns += "`" + colName + "`,"
//	}
//	insertColumns = util.TrimSuffix(insertColumns, ",") + ") "
//
//	tableExecutionString = tableExecutionString[:len(tableExecutionString)-1]
//	tableExecutionString += ",PRIMARY KEY (`record_id`));"
//	err = dbConnection.Exec(tableExecutionString).Error
//
//	fmt.Println("===================================================")
//	fmt.Println("tableExecutionString: ", tableExecutionString)
//
//	as.BaseService.Logger.Infow("created table creation string", "execution_string", tableExecutionString)
//	if err != nil {
//		// response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Schema"), ErrorCreatingObjectInformation, err.Error())
//		return err
//	}
//	fmt.Println("allData", len(allData))
//	for _, allStatsData := range allData {
//		// after successfull table creation, we needed to compose the insert statment and write.
//		insertStatement := "INSERT INTO " + targetTable + insertColumns + " VALUES("
//		values := ""
//		for _, rowStatsData := range allStatsData {
//			values += "'" + rowStatsData + "',"
//		}
//
//		insertStatement = insertStatement + util.TrimSuffix(values, ",") + ");"
//
//		fmt.Println(insertStatement)
//		err = dbConnection.Exec(insertStatement).Error
//
//		if err != nil {
//			dbConnection.Exec("DROP TABLE " + targetTable)
//			fmt.Println(err)
//			return err
//		}
//
//	}
//
//	return nil
//}
//
//// func readSourceData(as *AnalyticsService, fileName string, dataSource string, dbConnection *gorm.DB, statsColumnList []string) error {
//// 	f, err := excelize.OpenFile(fileName)
//// 	if err != nil {
//// 		fmt.Println(err)
//// 		return err
//// 	}
//// 	defer func() {
//// 		// Close the spreadsheet.
//// 		if err := f.Close(); err != nil {
//// 			fmt.Println(err)
//// 		}
//// 	}()
//
//// 	columnList := make([]string, 0)
//// 	originalColumnList := make(map[string]string)
//// 	allData := make([][]string, 0)
//// 	columnCount := 1
//// 	statsColumnCount := 1
//// 	for rowIndex := 21; rowIndex <= 51; rowIndex++ {
//// 		rowData := make([]string, 0)
//// 		for colIndex := 4; colIndex <= 34; colIndex++ {
//// 			cellIndex, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
//// 			cellValue, err := f.GetCellValue("sheet1", cellIndex)
//// 			if err != nil {
//// 				fmt.Println(err)
//// 				cellValue = ""
//// 			}
//
//// 			if rowIndex == 21 {
//// 				if cellValue == "Inspected by" {
//// 					columnList = append(columnList, "inspected_by")
//// 				} else if cellValue == "Remark" {
//// 					columnList = append(columnList, "remark")
//// 				} else if cellValue == "Mold #" {
//// 					columnList = append(columnList, "Mold")
//// 				} else if cellValue == "Machine #" {
//// 					columnList = append(columnList, "Machine")
//// 				} else if cellValue == "Sample#" {
//// 					columnList = append(columnList, "Sample")
//// 				} else if cellValue == "Working Shift" {
//// 					columnList = append(columnList, "working_shift")
//// 				} else {
//// 					renameColumn := strings.Join(strings.Fields(cellValue), "_") + "_" + strconv.Itoa(columnCount)
//
//// 					if (31 >= colIndex) && (colIndex >= 10) {
//// 						columnList = append(columnList, statsColumnList[statsColumnCount])
//// 						originalColumnList[statsColumnList[statsColumnCount]] = renameColumn
//// 						statsColumnCount += 1
//// 					} else {
//// 						columnList = append(columnList, renameColumn)
//// 					}
//
//// 				}
//
//// 				columnCount += 1
//// 			} else {
//// 				rowData = append(rowData, cellValue)
//// 			}
//
//// 		}
//// 		if len(rowData) != 0 {
//// 			allData = append(allData, rowData)
//// 		}
//
//// 	}
//// 	// originalColumnStr := fmt.Sprint(originalColumnList)
//
//// 	partNumberIndex, _ := excelize.CoordinatesToCellName(6, 14)
//// 	partNumber, _ := f.GetCellValue("sheet1", partNumberIndex)
//// 	// insertStatement := "INSERT INTO " + SPCResourceDataSourceTable + " (object_info)" + " VALUES('{\"name\":\"" + dataSource + "\",\"originalColumn\":\"" + originalColumnStr + "\",\"partNumber\":\"" + partNumber + "\",\"resourceData\":["
//// 	insertData := make(map[string]interface{})
//// 	insertData["name"] = dataSource
//// 	insertData["partNumber"] = partNumber
//// 	insertData["originalColumn"] = originalColumnList
//// 	resourceDataList := make([]interface{}, 0)
//// 	for _, allStatsData := range allData {
//// 		resourceData := make(map[string]interface{})
//// 		// after successfull table creation, we needed to compose the insert statment and write.
//// 		// values := "{"
//// 		for index, rowStatsData := range allStatsData {
//// 			// values += "\"" + columnList[index] + "\":" + "\"" + rowStatsData + "\","
//// 			resourceData[columnList[index]] = rowStatsData
//// 		}
//
//// 		resourceDataList = append(resourceDataList, resourceData)
//// 	}
//
//// 	// insertStatement = util.TrimSuffix(insertStatement, ",") + "]}');"
//
//// 	// err = dbConnection.Exec(insertStatement).Error
//// 	insertData["resourceData"] = resourceDataList
//// 	objectInfo, _ := json.Marshal(insertData)
//// 	insertObj := component.GeneralObject{ObjectInfo: objectInfo}
//
//// 	// fmt.Println(insertStatement)
//// 	err, _ = Create(dbConnection, SPCResourceDataSourceTable, insertObj)
//
//// 	if err != nil {
//// 		// dbConnection.Exec("DROP TABLE " + targetTable)
//// 		fmt.Println(err)
//// 		return err
//// 	}
//
//// 	return nil
//// }
//
//type WidgetTemplate struct {
//	Name          string        `json:"name"`
//	Tags          []interface{} `json:"tags"`
//	Query         string        `json:"query"`
//	Status        string        `json:"status"`
//	Filters       []interface{} `json:"filters"`
//	CreatedAt     time.Time     `json:"createdAt"`
//	CreatedBy     int           `json:"createdBy"`
//	Datasource    interface{}   `json:"datasource"`
//	Description   interface{}   `json:"description"`
//	ObjectStatus  string        `json:"objectStatus"`
//	LastUpdatedAt time.Time     `json:"lastUpdatedAt"`
//	LastUpdatedBy int           `json:"lastUpdatedBy"`
//	SeriesMapping struct {
//		Column string `json:"column"`
//	} `json:"seriesMapping"`
//	Visualization struct {
//		Type   string `json:"type"`
//		Title  string `json:"title"`
//		Series []struct {
//			Data      string `json:"data"`
//			Column    string `json:"column"`
//			FontSize  string `json:"fontSize"`
//			Position  string `json:"position"`
//			FontColor string `json:"fontColor"`
//		} `json:"series"`
//		CardBackground string `json:"cardBackground"`
//	} `json:"visualization"`
//	VisualizationType string `json:"visualizationType"`
//}
//
//var cavity1Cpk = `{
//    "name": "[$TABLE_NAME_1]Cavity 1 - CpK",
//    "tags": [],
//    "query": "SELECT  cast(cav_1_1 as decimal(5,2)) as val FROM fuyu_mes.$TABLE_NAME_1 where stat_key = 'Cpk' ",
//    "status": "Un-Published",
//    "filters": [],
//    "createdAt": "2023-04-02T07:05:16.190Z",
//    "createdBy": 2,
//    "datasource": null,
//    "description": null,
//    "objectStatus": "Active",
//    "lastUpdatedAt": "2023-04-02T07:05:16.190Z",
//    "lastUpdatedBy": 2,
//    "seriesMapping": {
//        "column": "val"
//    },
//    "visualization": {
//        "type": "number_card",
//        "title": "Cavity 1 - Cpk ",
//        "series": [
//            {
//                "data": "3.09",
//                "column": "val",
//                "fontSize": "60px",
//                "position": "center",
//                "fontColor": "#8DFEEA"
//            }
//        ],
//        "cardBackground": "#004239"
//    },
//    "visualizationType": "number_card"
//}`
//
//var inspectionCount = `{
//    "name": "[$TABLE_NAME_2] Number Of Inspection",
//    "tags": [],
//    "query": "SELECT inspected_by, count(*) noOfTest FROM fuyu_mes.$TABLE_NAME_2 where inspected_by <> \"\" group by inspected_by",
//    "status": "Un-Published",
//    "filters": [],
//    "createdAt": "2023-04-02T20:31:58.799Z",
//    "createdBy": 2,
//    "datasource": null,
//    "description": null,
//    "objectStatus": "Active",
//    "lastUpdatedAt": "2023-04-02T20:31:58.799Z",
//    "lastUpdatedBy": 2,
//    "seriesMapping": {
//        "xColumn": "inspected_by",
//        "yColumns": [
//            "noOfTest"
//        ]
//    },
//    "visualization": {
//        "type": "chart",
//        "chart": {},
//        "title": {
//            "text": "",
//            "zoomType": "xy"
//        },
//        "xAxis": {
//            "type": "auto",
//            "title": {
//                "enabled": true
//            }
//        },
//        "yAxis": [
//            {
//                "type": "auto",
//                "index": 0,
//                "title": {
//                    "text": "noOfTest"
//                },
//                "labels": {},
//                "opposite": false,
//                "gridLineWidth": 1
//            }
//        ],
//        "legend": {
//            "enabled": true
//        },
//        "series": [
//            {
//                "id": "e8cc6fe3c2eb3a6987b5d02c1e01b664",
//                "data": [
//                    [
//                        "Sara",
//                        2,
//                        "Sara",
//                        2
//                    ],
//                    [
//                        "Jin",
//                        2,
//                        "Jin",
//                        2
//                    ]
//                ],
//                "name": "noOfTest",
//                "type": "pie",
//                "color": "rgba(75, 239, 210, 1)"
//            }
//        ],
//        "credits": {
//            "enabled": false
//        },
//        "tooltip": {
//            "shared": true
//        },
//        "plotOptions": {
//            "series": {
//                "turboThreshold": 0
//            }
//        }
//    },
//    "visualizationType": "chart"
//}`
//
//type InspectionPieChartTemplate struct {
//	Name          string        `json:"name"`
//	Tags          []interface{} `json:"tags"`
//	Query         string        `json:"query"`
//	Status        string        `json:"status"`
//	Filters       []interface{} `json:"filters"`
//	CreatedAt     time.Time     `json:"createdAt"`
//	CreatedBy     int           `json:"createdBy"`
//	Datasource    interface{}   `json:"datasource"`
//	Description   interface{}   `json:"description"`
//	ObjectStatus  string        `json:"objectStatus"`
//	LastUpdatedAt time.Time     `json:"lastUpdatedAt"`
//	LastUpdatedBy int           `json:"lastUpdatedBy"`
//	SeriesMapping struct {
//		XColumn  string   `json:"xColumn"`
//		YColumns []string `json:"yColumns"`
//	} `json:"seriesMapping"`
//	Visualization struct {
//		Type  string `json:"type"`
//		Chart struct {
//		} `json:"chart"`
//		Title struct {
//			Text     string `json:"text"`
//			ZoomType string `json:"zoomType"`
//		} `json:"title"`
//		XAxis struct {
//			Type  string `json:"type"`
//			Title struct {
//				Enabled bool `json:"enabled"`
//			} `json:"title"`
//		} `json:"xAxis"`
//		YAxis []struct {
//			Type  string `json:"type"`
//			Index int    `json:"index"`
//			Title struct {
//				Text string `json:"text"`
//			} `json:"title"`
//			Labels struct {
//			} `json:"labels"`
//			Opposite      bool `json:"opposite"`
//			GridLineWidth int  `json:"gridLineWidth"`
//		} `json:"yAxis"`
//		Legend struct {
//			Enabled bool `json:"enabled"`
//		} `json:"legend"`
//		Series []struct {
//			Id    string          `json:"id"`
//			Data  [][]interface{} `json:"data"`
//			Name  string          `json:"name"`
//			Type  string          `json:"type"`
//			Color string          `json:"color"`
//		} `json:"series"`
//		Credits struct {
//			Enabled bool `json:"enabled"`
//		} `json:"credits"`
//		Tooltip struct {
//			Shared bool `json:"shared"`
//		} `json:"tooltip"`
//		PlotOptions struct {
//			Series struct {
//				TurboThreshold int `json:"turboThreshold"`
//			} `json:"series"`
//		} `json:"plotOptions"`
//	} `json:"visualization"`
//	VisualizationType string `json:"visualizationType"`
//}
//
//func addDashboardData(dashboardName string, dbConnection *gorm.DB, table1Name string, table2Name string) error {
//
//	// if adding widget into dashboard fails, then remove all the widget created, otherwise it will in-consistency
//	cavity1CpkWidgetTemplate := WidgetTemplate{}
//	json.Unmarshal([]byte(cavity1Cpk), &cavity1CpkWidgetTemplate)
//	cavity1CpkWidgetTemplate.Query = strings.Replace(cavity1CpkWidgetTemplate.Query, "$TABLE_NAME_1", table1Name, -1)
//	cavity1CpkWidgetTemplate.Name = strings.Replace(cavity1CpkWidgetTemplate.Name, "$TABLE_NAME_1", table1Name, -1)
//	serialisedData, _ := json.Marshal(cavity1CpkWidgetTemplate)
//	generalObject := component.GeneralObject{ObjectInfo: serialisedData}
//	err, cavity1CpkWidgetId := Create(dbConnection, AnalyticsWidgetTable, generalObject)
//	if err != nil {
//		return errors.New("error creating a widget template for given datasource")
//	}
//
//	inspectionPieChartTemplate := InspectionPieChartTemplate{}
//	json.Unmarshal([]byte(inspectionCount), &inspectionPieChartTemplate)
//
//	inspectionPieChartTemplate.Query = strings.Replace(inspectionPieChartTemplate.Query, "$TABLE_NAME_2", table2Name, -1)
//	inspectionPieChartTemplate.Name = strings.Replace(inspectionPieChartTemplate.Name, "$TABLE_NAME_2", table2Name, -1)
//
//	serialisedInspectionData, _ := json.Marshal(inspectionPieChartTemplate)
//	generalSerialiseInspectionObject := component.GeneralObject{ObjectInfo: serialisedInspectionData}
//	err, inspectionWidgetId := Create(dbConnection, AnalyticsWidgetTable, generalSerialiseInspectionObject)
//	if err != nil {
//		return errors.New("error creating a widget template for given datasource")
//	}
//
//	objectInfo := map[string]interface{}{"name": dashboardName,
//		"status":          "Pending",
//		"createdAt":       "2023-03-22T10:39:40.342Z",
//		"createdBy":       2,
//		"description":     dashboardName,
//		"objectStatus":    "Active",
//		"lastUpdatedAt":   "2023-03-22T10:39:40.342Z",
//		"refreshInterval": "5m",
//		"snapshotEnabled": true,
//		"dashboardWidgets": []interface{}{
//			map[string]interface{}{
//				"x":                         0,
//				"y":                         0,
//				"id":                        2,
//				"cols":                      2,
//				"rows":                      2,
//				"width":                     625.6666666666666,
//				"height":                    330,
//				"widgetId":                  inspectionWidgetId,
//				"visualizationType":         "chart",
//				"dashboardWidgetFilterInfo": make([]interface{}, 0),
//			},
//			map[string]interface{}{
//
//				"x":                         2,
//				"y":                         0,
//				"id":                        6,
//				"cols":                      2,
//				"rows":                      2,
//				"width":                     625.6666666666666,
//				"height":                    330,
//				"widgetId":                  cavity1CpkWidgetId,
//				"visualizationType":         "number_card",
//				"dashboardWidgetFilterInfo": make([]interface{}, 0),
//			},
//		},
//		"allowPublicAccess": nil}
//	preprocessedRequest, _ := json.Marshal(objectInfo)
//	object := component.GeneralObject{
//		ObjectInfo: preprocessedRequest,
//	}
//
//	err, _ = Create(dbConnection, AnalyticsDashboardTable, object)
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//	/*
//		ap[string]interface{}{
//						"x":                         2,
//						"y":                         0,
//						"id":                        5,
//						"cols":                      2,
//						"rows":                      2,
//						"width":                     402.3333333333333,
//						"height":                    330,
//						"widgetId":                  7,
//						"visualizationType":         "number_card",
//						"dashboardWidgetFilterInfo": make([]interface{}, 0),
//					},
//					map[string]interface{}{
//						"x":                         4,
//						"y":                         0,
//						"id":                        4,
//						"cols":                      2,
//						"rows":                      2,
//						"width":                     402.3333333333333,
//						"height":                    330,
//						"widgetId":                  6,
//						"visualizationType":         "number_card",
//						"dashboardWidgetFilterInfo": make([]interface{}, 0),
//					},
//					map[string]interface{}{
//						"x":                         0,
//						"y":                         2,
//						"id":                        3,
//						"cols":                      5,
//						"rows":                      3,
//						"width":                     1020.8333333333332,
//						"height":                    500,
//						"widgetId":                  5,
//						"visualizationType":         "chart",
//						"dashboardWidgetFilterInfo": make([]interface{}, 0),
//					},
//					map[string]interface{}{
//						"x":                         5,
//						"y":                         2,
//						"id":                        2,
//						"cols":                      1,
//						"rows":                      3,
//						"width":                     196.16666666666664,
//						"height":                    500,
//						"widgetId":                  4,
//						"visualizationType":         "chart",
//						"dashboardWidgetFilterInfo": make([]interface{}, 0),
//					},
//					map[string]interface{}{
//						"x":                         0,
//						"y":                         5,
//						"id":                        1,
//						"cols":                      6,
//						"rows":                      4,
//						"width":                     1227,
//						"height":                    670,
//						"widgetId":                  3,
//						"visualizationType":         "chart",
//						"dashboardWidgetFilterInfo": make([]interface{}, 0),
//					},
//	*/
//
//}
