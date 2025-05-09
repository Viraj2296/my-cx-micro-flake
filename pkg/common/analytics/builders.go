package analytics

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/pkg/util/encryption"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"strconv"
	"time"
)

type QueryResponseBuilder interface {
	BuildResponse() (error, interface{})
}

type BaseBuilder struct {
	ListOfColumns        []string
	QueryResults         []map[string]interface{}
	QueryResponseRequest datatypes.JSON
	SchemaList           []component.TableSchema
}

type BaseVisualisationBuilder struct {
	ListOfColumns []string
	QueryResults  []map[string]interface{}
	SchemaList    []component.TableSchema
	Visualisation map[string]interface{}
	SeriesMapping map[string]interface{}
}

func (bb *BaseVisualisationBuilder) Init() {
	for _, result := range bb.QueryResults {
		for key, _ := range result {
			tableSchema := component.TableSchema{Name: key, Property: key, Display: true, Type: "text"}
			bb.SchemaList = append(bb.SchemaList, tableSchema)
			bb.ListOfColumns = append(bb.ListOfColumns, key)
		}
		break
	}
	fmt.Println("bb.ListOfColumns: ", bb.ListOfColumns)
}
func (bvb BaseVisualisationBuilder) getGroupByColumnUniqueValues(groupByCol string) interface{} {
	var listOfValues []interface{}
	for _, result := range bvb.QueryResults {
		interfaceValue := result[groupByCol]
		listOfValues = append(listOfValues, interfaceValue)
	}
	return util.Unique(listOfValues)
}

func (bb *BaseBuilder) Init() {
	for _, result := range bb.QueryResults {
		for key, _ := range result {
			tableSchema := component.TableSchema{Name: key, Property: key, Display: true, Type: "text"}
			bb.SchemaList = append(bb.SchemaList, tableSchema)
			bb.ListOfColumns = append(bb.ListOfColumns, key)
		}
		break
	}
	fmt.Println("bb.ListOfColumns: ", bb.ListOfColumns)
}

type ChartVisualisationBuilder struct {
	BaseVisualisationBuilder *BaseVisualisationBuilder
}

type TableVisualisationBuilder struct {
	BaseVisualisationBuilder *BaseVisualisationBuilder
}

type GaugeVisualisationBuilder struct {
	BaseVisualisationBuilder *BaseVisualisationBuilder
}

type NumberCardVisualisationBuilder struct {
	BaseVisualisationBuilder *BaseVisualisationBuilder
}

type BulletGraphBuilder struct {
	BaseVisualisationBuilder *BaseVisualisationBuilder
}

type TimelineVisualisationBuilder struct {
	BaseVisualisationBuilder *BaseVisualisationBuilder
}

func (tvb TableVisualisationBuilder) BuildResponse() (error, interface{}) {
	fmt.Println("cvb.TableVisualisationBuilder.QueryResults : ", tvb.BaseVisualisationBuilder.QueryResults)

	//tableColumns := util.InterfaceToStringArray(tvb.BaseVisualisationBuilder.SeriesMapping["tableColumns"])
	/*
			"seriesMapping": {
			        "tableColumns": [
			            {
			                "name": "Completed Quantity",
			                "type": "number"
			            },
			            {
			                "name": "End Time",
			                "type": "date"
			            },
			            {
			                "name": "Machine Name",
			                "type": "text"
			            }
			        ]
			    },
		it will accept string array or array of object with name, and data type
	*/
	tableColumnsData := tvb.BaseVisualisationBuilder.SeriesMapping["tableColumns"]
	var tableColumnNames []string
	var tableDataMap = make(map[string]string)
	switch v := tableColumnsData.(type) {
	case []interface{}:
		if areObjects(v) {
			tableColumns := tableColumnsData.([]interface{})
			for _, tableColData := range tableColumns {
				tableColMapData := tableColData.(map[string]interface{})
				if value, ok := tableColMapData["name"]; ok {
					tableColumnNames = append(tableColumnNames, util.InterfaceToString(value))
					if typeValue, ok := tableColMapData["type"]; ok {
						tableDataMap[util.InterfaceToString(value)] = util.InterfaceToString(typeValue)
					} else {
						tableDataMap[util.InterfaceToString(value)] = "text"
					}
				}
			}
		} else {
			tableColumnNames = util.InterfaceToStringArray(tableColumnsData)
		}

	default:
		fmt.Println("Invalid type: It's neither []string nor an array of objects. {}", v)
	}

	tableDataResponse := component.TableDataResponse{}
	for _, result := range tvb.BaseVisualisationBuilder.QueryResults {
		var objectKeyValues = make(map[string]interface{}, len(tvb.BaseVisualisationBuilder.QueryResults))

		for _, tableCol := range tableColumnNames {
			for key, value := range result {
				if tableCol == key {
					if typeValue, ok := tableDataMap[tableCol]; ok {
						if typeValue == "date" {
							var timeValue = value.(time.Time)
							parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeValue.String())
							if err != nil {
								fmt.Println("Error parsing time:", err)
								objectKeyValues[key] = value
							}

							// Format the time as a string without fractional seconds and timezone offset
							formattedTime := parsedTime.Format("2006-01-02 15:04:05")
							objectKeyValues[key] = formattedTime
						} else {
							objectKeyValues[key] = value
						}
					} else {
						objectKeyValues[key] = value
					}

				}
			}
		}

		raw, _ := json.Marshal(objectKeyValues)
		tableDataResponse.Data = append(tableDataResponse.Data, raw)
	}

	var modifiedTableSchema []component.TableSchema
	for _, tableCol := range tableColumnNames {
		for _, schemaElement := range tvb.BaseVisualisationBuilder.SchemaList {
			if tableCol == schemaElement.Property {
				if dataType, ok := tableDataMap[tableCol]; ok {
					schemaElement.Type = dataType
				}
				modifiedTableSchema = append(modifiedTableSchema, schemaElement)
			}
		}
	}
	tableDataResponse.Header = modifiedTableSchema

	tableDataResponse.TotalRowCount = len(tvb.BaseVisualisationBuilder.QueryResults)

	//TODO we need to make sure, we don't supply non-integer values to avoid graph issues
	visualizationSeries := tvb.BaseVisualisationBuilder.Visualisation["series"].([]interface{})
	// we should have only one series
	existingVisualSeries := visualizationSeries[0].(map[string]interface{})
	existingVisualSeries["data"] = tableDataResponse.Data
	existingVisualSeries["header"] = tableDataResponse.Header
	existingVisualSeries["totalRowCount"] = tableDataResponse.TotalRowCount
	fmt.Println("visualizationSeries: ", existingVisualSeries)
	var builtNewVisualArray []interface{}
	builtNewVisualArray = append(builtNewVisualArray, existingVisualSeries)
	return nil, builtNewVisualArray

}

func areObjects(arr []interface{}) bool {
	for _, item := range arr {
		if _, ok := item.(interface{}); !ok {
			return false
		}
	}
	return true
}

func (cvb TimelineVisualisationBuilder) BuildResponse() (error, interface{}) {
	fmt.Println("cvb.GaugeVisualisationBuilder.QueryResults : ", cvb.BaseVisualisationBuilder.QueryResults)

	timelineColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "timelineColumn")
	labelColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "labelColumn")
	nameColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "nameColumn")
	descriptionColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "descriptionColumn")
	var dataArray []interface{}
	for _, result := range cvb.BaseVisualisationBuilder.QueryResults {
		var internalData = make(map[string]interface{})
		if timelineColumn != "" {
			value := result[timelineColumn]
			internalData["x"] = util.InterfaceToInt(value)
		}
		if labelColumn != "" {
			value := result[labelColumn]
			internalData["label"] = value
		}
		if nameColumn != "" {
			value := result[nameColumn]
			internalData["name"] = value
		}
		if descriptionColumn != "" {
			value := result[descriptionColumn]
			internalData["description"] = value
		}
		dataArray = append(dataArray, internalData)

	}

	//TODO we need to make sure, we don't supply non-integer values to avoid graph issues
	visualizationSeries := cvb.BaseVisualisationBuilder.Visualisation["series"].([]interface{})
	// we should have only one series
	existingVisualSeries := visualizationSeries[0].(map[string]interface{})
	existingVisualSeries["data"] = dataArray
	fmt.Println("visualizationSeries: ", existingVisualSeries)
	var builtNewVisualArray []interface{}
	builtNewVisualArray = append(builtNewVisualArray, existingVisualSeries)
	return nil, builtNewVisualArray

}
func (cvb NumberCardVisualisationBuilder) BuildResponse() (error, interface{}) {
	numberCardColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "column")
	var numberCardValue interface{}
	fmt.Println("cvb.GaugeVisualisationBuilder.QueryResults : ", cvb.BaseVisualisationBuilder.QueryResults)
	if len(cvb.BaseVisualisationBuilder.QueryResults) > 0 {
		// got more than one records, it is not possible , anyway, inorder to safe gurd us, lets take the first one
		numberCardValue = cvb.BaseVisualisationBuilder.QueryResults[0][numberCardColumn]
	} else if len(cvb.BaseVisualisationBuilder.QueryResults) == 0 {
		numberCardValue = ""
	}
	var dataArray []interface{}
	dataArray = append(dataArray, numberCardValue)
	//TODO we need to make sure, we don't supply non-integer values to avoid graph issues
	visualizationSeries := cvb.BaseVisualisationBuilder.Visualisation["series"].([]interface{})
	// we should have only one series
	existingVisualSeries := visualizationSeries[0].(map[string]interface{})
	existingVisualSeries["data"] = dataArray
	fmt.Println("visualizationSeries: ", existingVisualSeries)
	var builtNewVisualArray []interface{}
	builtNewVisualArray = append(builtNewVisualArray, existingVisualSeries)
	return nil, builtNewVisualArray

}

func (cvb BulletGraphBuilder) BuildResponse() (error, interface{}) {
	actualColum := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "actual")
	targetColum := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "target")
	//expectedColum := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "target")
	var actualValue interface{}
	var targetValue interface{}
	fmt.Println("cvb.GaugeVisualisationBuilder.QueryResults : ", cvb.BaseVisualisationBuilder.QueryResults)
	if len(cvb.BaseVisualisationBuilder.QueryResults) > 0 {
		// got more than one records, it is not possible , anyway, inorder to safe gurd us, lets take the first one
		actualValue = util.InterfaceToInt(cvb.BaseVisualisationBuilder.QueryResults[0][actualColum])
		targetValue = util.InterfaceToInt(cvb.BaseVisualisationBuilder.QueryResults[0][targetColum])
	} else if len(cvb.BaseVisualisationBuilder.QueryResults) == 0 {
		actualValue = 0
		targetValue = 0
	}
	type bulletData struct {
		Y      interface{} `json:"y"`
		Target interface{} `json:"target"`
	}
	var dataArray []interface{}
	dataArray = append(dataArray, bulletData{Y: actualValue, Target: targetValue})
	//TODO we need to make sure, we don't supply non-integer values to avoid graph issues
	visualizationSeries := cvb.BaseVisualisationBuilder.Visualisation["series"].([]interface{})
	// we should have only one series
	existingVisualSeries := visualizationSeries[0].(map[string]interface{})
	existingVisualSeries["data"] = dataArray
	fmt.Println("visualizationSeries: ", existingVisualSeries)
	var builtNewVisualArray []interface{}
	builtNewVisualArray = append(builtNewVisualArray, existingVisualSeries)
	return nil, builtNewVisualArray
}

func (cvb GaugeVisualisationBuilder) BuildResponse() (error, interface{}) {
	gaugeColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "column")
	var gaugeValue interface{}
	fmt.Println("cvb.GaugeVisualisationBuilder.QueryResults : ", cvb.BaseVisualisationBuilder.QueryResults)
	if len(cvb.BaseVisualisationBuilder.QueryResults) > 0 {
		// got more than one records, it is not possible , anyway, inorder to safe gurd us, lets take the first one
		gaugeValue = cvb.BaseVisualisationBuilder.QueryResults[0][gaugeColumn]
	} else if len(cvb.BaseVisualisationBuilder.QueryResults) == 0 {
		gaugeValue = 0
	}

	//TODO we need to make sure, we don't supply non-integer values to avoid graph issues
	var dataArray []interface{}
	dataArray = append(dataArray, gaugeValue)
	visualizationSeries := cvb.BaseVisualisationBuilder.Visualisation["series"].([]interface{})
	// we should have only one series
	existingVisualSeries := visualizationSeries[0].(map[string]interface{})
	existingVisualSeries["data"] = dataArray
	fmt.Println("visualizationSeries: ", existingVisualSeries)
	var builtNewVisualArray []interface{}
	builtNewVisualArray = append(builtNewVisualArray, existingVisualSeries)
	return nil, builtNewVisualArray

}

func (cvb ChartVisualisationBuilder) BuildResponse() (error, interface{}) {
	xColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "xColumn")
	yColumns := util.InterfaceToStringArray(cvb.BaseVisualisationBuilder.SeriesMapping["yColumns"])
	groupByColumn := util.MapInterfaceToString(cvb.BaseVisualisationBuilder.SeriesMapping, "groupByColumn")

	var dynamicSeriesName = make(map[string]interface{}, 0)
	var seriesMap = make(map[string]interface{}, 0)
	if xColumn != "" && len(yColumns) > 0 {
		if util.StringArrayContains(cvb.BaseVisualisationBuilder.ListOfColumns, xColumn) {
			if groupByColumn == "" {
				for _, yCol := range yColumns {
					var seriesData []interface{}
					for _, result := range cvb.BaseVisualisationBuilder.QueryResults {
						xColValue := result[xColumn]
						yColValue := result[yCol]
						var data []interface{}
						switch v := yColValue.(type) {

						case float64:
							floatValue := util.InterfaceToFloat(v)
							data = append(data, xColValue)
							data = append(data, floatValue)
						case int64:
							intValue := util.InterfaceToInt(v)
							data = append(data, xColValue)
							data = append(data, intValue)
						default:
							if value, err := strconv.Atoi(yColValue.(string)); err == nil {
								// we got the number
								data = append(data, xColValue)
								data = append(data, value)
							} else if value, err := strconv.ParseFloat(yColValue.(string), 10); err == nil {
								data = append(data, xColValue)
								data = append(data, value)
							} else {
								data = append(data, xColValue)
								data = append(data, util.InterfaceToString(yColValue))
							}

						}

						seriesData = append(seriesData, data)
					}
					key := xColumn + ":" + yCol
					seriesIdHash := encryption.GetMD5Hash(key)
					dynamicSeriesName[seriesIdHash] = yCol
					seriesMap[seriesIdHash] = seriesData
				}
			} else {
				// we have group by column
				groupByColUniqValues := cvb.BaseVisualisationBuilder.getGroupByColumnUniqueValues(groupByColumn)
				interfaceValueArray := groupByColUniqValues.([]interface{})
				for _, groupByValue := range interfaceValueArray {
					var filteredFrames []map[string]interface{}
					for _, result := range cvb.BaseVisualisationBuilder.QueryResults {
						interfaceValue := result[groupByColumn]
						if interfaceValue == groupByValue {
							filteredFrames = append(filteredFrames, result)
						}
					}
					for _, yCol := range yColumns {
						var seriesData []interface{}
						for _, filteredResult := range filteredFrames {
							xColValue := filteredResult[xColumn]
							yColValue := filteredResult[yCol]
							var data []interface{}
							data = append(data, xColValue)
							switch v := yColValue.(type) {

							case float64:
								floatValue := util.InterfaceToFloat(v)

								data = append(data, floatValue)
							case int64:
								intValue := util.InterfaceToInt(v)

								data = append(data, intValue)
							default:
								if value, err := strconv.Atoi(yColValue.(string)); err == nil {
									// we got the number

									data = append(data, value)
								} else if value, err := strconv.ParseFloat(yColValue.(string), 10); err == nil {

									data = append(data, value)
								} else {

									data = append(data, util.InterfaceToString(yColValue))
								}

							}
							seriesData = append(seriesData, data)
						}
						key := xColumn + ":" + yCol + ":" + groupByValue.(string)
						seriesIdHash := encryption.GetMD5Hash(key)
						dynamicSeriesName[seriesIdHash] = groupByValue
						seriesMap[seriesIdHash] = seriesData

					}

				}
			}
		}
	}
	fmt.Println("seriesMap: ", seriesMap)
	fmt.Println("cvb.BaseVisualisationBuilder.Visualisation_series", cvb.BaseVisualisationBuilder.Visualisation["series"])
	visualizationSeries := cvb.BaseVisualisationBuilder.Visualisation["series"].([]interface{})
	var builtNewVisualArray []interface{}
	for _, seriesVisual := range visualizationSeries {
		seriesFields := seriesVisual.(map[string]interface{})
		if len(seriesFields) > 0 {
			if seriesId, ok := seriesFields["id"]; ok {
				id := seriesId.(string)
				if generatedSeriesData, ok := seriesMap[id]; ok {
					fmt.Println("series visual :", generatedSeriesData)
					seriesFields["data"] = generatedSeriesData
					builtNewVisualArray = append(builtNewVisualArray, seriesVisual)
					seriesMap[id] = nil
				}
			}

		}

	}

	var seriesVisual = make(map[string]interface{})

	for seriesId, _ := range seriesMap {
		if seriesMap[seriesId] != nil {
			seriesVisual["yAxis"] = 0
			seriesVisual["type"] = "line"
			seriesVisual["name"] = dynamicSeriesName[seriesId]
			generatedSeriesData := seriesMap[seriesId]
			seriesVisual["data"] = generatedSeriesData

			builtNewVisualArray = append(builtNewVisualArray, seriesVisual)
		}

	}
	fmt.Println("visualizationSeries: ", builtNewVisualArray)
	return nil, builtNewVisualArray

}

type ChartResponseBuilder struct {
	BaseBuilder *BaseBuilder
}

type TimelineResponseBuilder struct {
	BaseBuilder *BaseBuilder
}

type TableResponseBuilder struct {
	BaseBuilder *BaseBuilder
}

func (trb TimelineResponseBuilder) BuildResponse() (error, interface{}) {
	timelineResponseRequest := TimelineFormatRequest{}
	json.Unmarshal(trb.BaseBuilder.QueryResponseRequest, &timelineResponseRequest)
	var dataArray []interface{}

	var seriesResponse = make(map[string]interface{}, 1)
	for _, result := range trb.BaseBuilder.QueryResults {
		var internalData = make(map[string]interface{})
		if timelineResponseRequest.TimelineColumn != "" {
			value := result[timelineResponseRequest.TimelineColumn]
			internalData["x"] = util.InterfaceToInt(value)
		}
		if timelineResponseRequest.LabelColumn != "" {
			value := result[timelineResponseRequest.LabelColumn]
			internalData["label"] = value
		}
		if timelineResponseRequest.NameColumn != "" {
			value := result[timelineResponseRequest.NameColumn]
			internalData["name"] = value
		}
		if timelineResponseRequest.DescriptionColumn != "" {
			value := result[timelineResponseRequest.DescriptionColumn]
			internalData["description"] = value
		}
		dataArray = append(dataArray, internalData)

	}
	var individualSeriesData = make(map[string]interface{})
	var dataLabelMap = make(map[string]interface{})
	dataLabelMap["connectorColor"] = "#2AD31F"
	dataLabelMap["connectorWidth"] = "5"
	individualSeriesData["dataLabels"] = dataLabelMap
	individualSeriesData["data"] = dataArray
	seriesResponse["series"] = individualSeriesData

	return nil, seriesResponse
}
func (trb TableResponseBuilder) BuildResponse() (error, interface{}) {
	tableResponseRequest := TableFormatRequest{}
	json.Unmarshal(trb.BaseBuilder.QueryResponseRequest, &tableResponseRequest)

	fmt.Println("tableResponseRequest: ", tableResponseRequest)

	tableResponse := component.WidgetTableObjectResponse{}
	totalRecords := int64(len(trb.BaseBuilder.QueryResults))
	tableResponse.TotalRowCount = totalRecords

	// now we built the header

	for _, result := range trb.BaseBuilder.QueryResults {
		var objectKeyValues = make(map[string]interface{}, len(trb.BaseBuilder.QueryResults))
		if len(tableResponseRequest.TableColumns) > 0 {
			for _, tableCol := range tableResponseRequest.TableColumns {
				for key, value := range result {
					if tableCol == key {
						objectKeyValues[key] = value
					}
				}
			}
		} else {
			for key, value := range result {

				var typedValue interface{}
				switch v := value.(type) {
				case int:
					typedValue = util.InterfaceToInt(v)
					break
				case int32:
					typedValue = util.InterfaceToInt(v)
					break
				case float64:
					typedValue = util.InterfaceToFloat(v)
					break
				case float32:
					typedValue = util.InterfaceToFloat(v)
					break
				case int64:
					typedValue = util.InterfaceToInt(v)
					break
				default:
					if stringValue, ok := v.(string); ok {
						if valueInt, err := strconv.Atoi(stringValue); err == nil {
							typedValue = util.InterfaceToInt(valueInt)
						} else if valueFloat, err := strconv.ParseFloat(stringValue, 64); err == nil {
							typedValue = util.InterfaceToFloat(valueFloat)
						} else {
							// Handle the case when stringValue is not a valid number
							typedValue = util.InterfaceToString(v)
						}
					} else if timeValue, ok := v.(time.Time); ok {
						// Handle the case when v is of type time.Time
						// For example, you might want to convert it to a string or handle it differently
						typedValue = util.InterfaceToTime(timeValue)
					} else {
						// Handle other cases
						typedValue = util.InterfaceToString(v)
					}

					break
				}

				objectKeyValues[key] = typedValue
			}
		}
		raw, _ := json.Marshal(objectKeyValues)
		tableResponse.Data = append(tableResponse.Data, raw)
	}

	if len(tableResponseRequest.TableColumns) > 0 {
		var modifiedTableSchema []component.TableSchema
		for _, tableCol := range tableResponseRequest.TableColumns {
			for _, schemaElement := range trb.BaseBuilder.SchemaList {
				if tableCol == schemaElement.Property {
					modifiedTableSchema = append(modifiedTableSchema, schemaElement)
				}
			}
		}
		tableResponse.Header = modifiedTableSchema
	} else {
		tableResponse.Header = trb.BaseBuilder.SchemaList
	}

	for _, schemaElement := range trb.BaseBuilder.SchemaList {
		tableResponse.XColumnList = append(tableResponse.XColumnList, schemaElement.Property)
		tableResponse.YColumnList = append(tableResponse.YColumnList, schemaElement.Property)
	}
	tableResponse.TotalRowCount = totalRecords
	return nil, tableResponse
}

func (crb ChartResponseBuilder) getGroupByColumnUniqueValues(groupByCol string) interface{} {
	var listOfValues []interface{}
	for _, result := range crb.BaseBuilder.QueryResults {
		interfaceValue := result[groupByCol]
		listOfValues = append(listOfValues, interfaceValue)
	}
	return util.Unique(listOfValues)
}
func (crb ChartResponseBuilder) BuildResponse() (error, interface{}) {
	var response []interface{}
	charResponseRequest := ChartFormatRequest{}
	json.Unmarshal(crb.BaseBuilder.QueryResponseRequest, &charResponseRequest)
	fmt.Println("============================================")
	fmt.Println("charResponseRequest", charResponseRequest)
	fmt.Println("============================================")
	if charResponseRequest.XColumn != "" && len(charResponseRequest.YColumns) > 0 {
		if util.StringArrayContains(crb.BaseBuilder.ListOfColumns, charResponseRequest.XColumn) {
			if charResponseRequest.GroupByCol == "" {

				for _, yCol := range charResponseRequest.YColumns {
					var seriesList = make(map[string]interface{}, len(charResponseRequest.YColumns))
					var seriesData []interface{}
					for _, result := range crb.BaseBuilder.QueryResults {

						xColValue := result[charResponseRequest.XColumn]
						yColValue := result[yCol]
						var data []interface{}
						switch v := yColValue.(type) {

						case float64:
							floatValue := util.InterfaceToFloat(v)
							data = append(data, xColValue)
							data = append(data, floatValue)

						case float32:
							floatValue := util.InterfaceToFloat(v)
							data = append(data, xColValue)
							data = append(data, floatValue)

						case int64:
							intValue := util.InterfaceToInt(v)
							data = append(data, xColValue)
							data = append(data, intValue)

						default:
							if value, err := strconv.Atoi(yColValue.(string)); err == nil {
								// we got the number

								data = append(data, xColValue)
								data = append(data, value)

							} else if value, err := strconv.ParseFloat(yColValue.(string), 10); err == nil {

								data = append(data, xColValue)
								data = append(data, value)

							} else {
								data = append(data, xColValue)
								fmt.Println("string:", xColValue)
								data = append(data, util.InterfaceToString(yColValue))
							}

						}

						//data = append(data, xColValue)
						fmt.Println("data1", data)
						//data = append(data, yColValue)

						seriesData = append(seriesData, data)
					}
					seriesList["data"] = seriesData
					key := charResponseRequest.XColumn + ":" + yCol
					seriesList["id"] = encryption.GetMD5Hash(key)
					seriesList["name"] = yCol
					response = append(response, seriesList)

				}
			} else {
				groupByColUniqValues := crb.getGroupByColumnUniqueValues(charResponseRequest.GroupByCol)
				interfaceValueArray := groupByColUniqValues.([]interface{})
				for _, groupByValue := range interfaceValueArray {
					var filteredFrames []map[string]interface{}
					for _, result := range crb.BaseBuilder.QueryResults {
						interfaceValue := result[charResponseRequest.GroupByCol]
						if interfaceValue == groupByValue {
							filteredFrames = append(filteredFrames, result)
						}
					}
					for _, yCol := range charResponseRequest.YColumns {
						var seriesList = make(map[string]interface{}, len(charResponseRequest.YColumns))
						var seriesData []interface{}
						for _, filteredResult := range filteredFrames {
							xColValue := filteredResult[charResponseRequest.XColumn]
							yColValue := filteredResult[yCol]
							var data []interface{}
							data = append(data, xColValue)
							switch v := yColValue.(type) {

							case float64:
								floatValue := util.InterfaceToFloat(v)

								data = append(data, floatValue)

							case float32:
								floatValue := util.InterfaceToFloat(v)

								data = append(data, floatValue)

							case int64:
								intValue := util.InterfaceToInt(v)

								data = append(data, intValue)

							default:
								if value, err := strconv.Atoi(yColValue.(string)); err == nil {
									// we got the number
									data = append(data, value)

								} else if value, err := strconv.ParseFloat(yColValue.(string), 10); err == nil {

									data = append(data, value)

								} else {

									fmt.Println("string:", xColValue)
									data = append(data, util.InterfaceToString(yColValue))
								}

							}

							seriesData = append(seriesData, data)
						}
						seriesList["data"] = seriesData
						key := charResponseRequest.XColumn + ":" + yCol + ":" + groupByValue.(string)
						seriesList["id"] = encryption.GetMD5Hash(key)
						seriesList["name"] = groupByValue
						response = append(response, seriesList)
					}

				}
			}
		}
	}
	fmt.Println()
	var seriesResponse = make(map[string]interface{}, 1)
	seriesResponse["series"] = response

	return nil, seriesResponse
}
