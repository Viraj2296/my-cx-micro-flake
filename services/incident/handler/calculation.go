package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

// This calculation index are same to all
func (v *IncidentService) calculateSafetyIndex() {
	//
}

func (v *IncidentService) calculateQualityIndex(year, month int) {
	// first load all the quality_category from table
	// then each reason, loop, and get all the incident_quality record with year, and month condition
	// calculate individual percentage totalnumberOfOccurance/ total -(no work day)* 100
	// then take the average of each category
	// Not meeet the exact  requirment  : 12 - 89
	// Not delivered to cutomer : 3 - 35
	// 89 +35/2 = index value

}

func (v *IncidentService) calculateDeliveryIndex() {

}

func (v *IncidentService) calculateInventoryIndex() {

}

func (v *IncidentService) calculateProductivityIndex() {

}

func (v *IncidentService) calculateSafetyTrend() {
	// take that particular day, sum up and divided by total
}

func (v *IncidentService) calculateMatrixDayView(dateQuery string, incidentObject *[]component.GeneralObject, incidentCategoryObject *[]component.GeneralObject) map[string]interface{} {
	result := make(map[string]interface{})

	if incidentObject == nil {
		return result
	}

	// dateQueryList := strings.Split(dateQuery, "/")
	// currentMonth, _ := strconv.Atoi(dateQueryList[0])
	// currentYear, _ := strconv.Atoi(dateQueryList[1])
	// currentLocation := time.Now().UTC().Location()
	// firstDayOfMonth := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, currentLocation)
	// lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	// lastDayOfMonth = lastDayOfMonth.AddDate(0, 0, 1)

	incidentDataMap := make(map[string][][]string)
	for _, object := range *incidentCategoryObject {
		//Intiate incidentDataMap with [key(description)] := [][]string
		dayDefectArray := make([][]string, 0)
		var categoryInfo = make(map[string]interface{})
		json.Unmarshal(object.ObjectInfo, &categoryInfo)
		categoryName := util.InterfaceToString(categoryInfo["name"])
		incidentDataMap[categoryName] = dayDefectArray
	}

	for _, incident := range *incidentObject {
		var incidentInfo = make(map[string]interface{})
		json.Unmarshal(incident.ObjectInfo, &incidentInfo)

		day := util.InterfaceToInt(incidentInfo["day"])

		if incidentInfo["incidents"] != nil {
			incidentList := incidentInfo["incidents"].([]interface{})
			isNoWorkDay := util.InterfaceToBool(incidentInfo["isNoWorkDay"])
			for _, object := range incidentList {
				incidentObject := object.(map[string]interface{})
				categoryId := util.InterfaceToInt(incidentObject["id"])
				categoryInfo := findCategoryObject(incidentCategoryObject, categoryId)
				if categoryNameInterface, ok := categoryInfo["name"]; ok {
					var categoryName = util.InterfaceToString(categoryNameInterface)
					if isNoWorkDay {
						// if it is no work day, then send the no work day
						dayDefectPairArray := []string{strconv.Itoa(day), "NO WORK"}
						incidentDataMap[categoryName] = append(incidentDataMap[categoryName], [][]string{dayDefectPairArray}...)
					} else {
						if !incidentObject["indicator"].(bool) {
							dayDefectPairArray := []string{strconv.Itoa(day), "GOOD"}
							incidentDataMap[categoryName] = append(incidentDataMap[categoryName], [][]string{dayDefectPairArray}...)
						} else {
							dayDefectPairArray := []string{strconv.Itoa(day), "ISSUE"}
							incidentDataMap[categoryName] = append(incidentDataMap[categoryName], [][]string{dayDefectPairArray}...)
						}
					}
				}

			}
		}
	}

	descriptionHeader := map[string]interface{}{"name": "Description",
		"type":         "text",
		"display":      true,
		"property":     "description",
		"routeEnabled": false,
		"columSize":    0}

	headerArray := []interface{}{descriptionHeader}

	dataArray := make([]interface{}, 0)

	headerLoopFlag := 1
	for key, dayDefectArray := range incidentDataMap {
		//dataResult formate
		// {"description": "Falling Down",
		// "1": "GOOD",
		// "2": "ISSUE"}
		dataResult := make(map[string]string)
		dataResult["description"] = key

		for _, dayDefectValue := range dayDefectArray {
			dataResult[dayDefectValue[0]] = dayDefectValue[1]

			if headerLoopFlag == 1 {
				//This one creating colums in the header
				colorCodeHeader := map[string]interface{}{
					"GOOD":    "#38A976",
					"ISSUE":   "#FF6265",
					"NO WORK": "#CAC824"}
				objectHeader := map[string]interface{}{"name": dayDefectValue[0],
					"type":         "background_color",
					"display":      true,
					"property":     dayDefectValue[0],
					"routeEnabled": false,
					"columSize":    0,
					"objectList":   colorCodeHeader}

				headerArray = append(headerArray, objectHeader)
			}

		}
		headerLoopFlag = 0
		dataArray = append(dataArray, dataResult)

	}

	result["totalRowCount"] = len(incidentDataMap)
	result["currentRowCount"] = 1

	result["header"] = headerArray

	result["data"] = dataArray

	return result
}

func (v *IncidentService) getIncidentOverviewResult(generalList *[]component.GeneralObject, generalCategoryList *[]component.GeneralObject, incidentType string) map[string]interface{} {
	response := make(map[string]interface{})
	localZone, _ := time.LoadLocation("Asia/Singapore")
	_, _, currentDay := time.Now().In(localZone).Date()

	// generalList = filterByDepartmentId(*generalList, generalCategoryList)

	categoryNameIdMap := make(map[string]int)
	for _, object := range *generalCategoryList {
		categoryInfo := make(map[string]interface{})
		json.Unmarshal(object.ObjectInfo, &categoryInfo)
		name := util.InterfaceToString(categoryInfo["name"])
		categoryNameIdMap[name] = object.Id
	}

	incidentMap := make(map[int]interface{}, 0)
	dayCount := 1
	for dayCount <= 31 {
		incidentInfo := findIncidentObjectGivenDay(generalList, dayCount)

		fmt.Println("incidentInfo: , daycount ", dayCount, "incidentInfo: ", incidentInfo)
		if incidentInfo == nil {
			dayCount += 1
			continue
		}
		incidentResponse := make(map[string]interface{})

		isNoWorkDay := incidentInfo["isNoWorkDay"].(bool)

		day := util.InterfaceToInt(incidentInfo["day"])
		fmt.Println("daycount : ", dayCount, " day :", day)
		if dayCount <= day {
			if isNoWorkDay {
				//yellow
				incidentResponse["color"] = "#CAC824"
				incidentResponse["day"] = day
				incidentResponse["visible"] = true
			} else {
				if incidentInfo["incidents"] != nil {
					incidents := incidentInfo["incidents"].([]interface{})
					if len(incidents) == 0 {
						if currentDay < day {
							//dark grey
							incidentResponse["color"] = "#696e6c"
							incidentResponse["day"] = day
							incidentResponse["visible"] = true
						} else {
							//green
							incidentResponse["color"] = "#38A976"
							incidentResponse["day"] = day
							incidentResponse["visible"] = true
						}

					} else {
						//red
						var isIssueHappened bool
						isIssueHappened = false
						var emptyIncident int
						emptyIncident = 0
						for _, incident := range incidents {
							incidentFields := incident.(map[string]interface{})
							indicator := util.InterfaceToBool(incidentFields["indicator"])
							categoryId := util.InterfaceToInt(incidentFields["id"])

							// check the category id is active or not
							categoryInfo := findCategoryObject(generalCategoryList, categoryId)

							if categoryInfo != nil {
								if indicator {
									isIssueHappened = true
									// something is already happened
									break
								}
							} else {
								emptyIncident = emptyIncident + 1
							}

						}
						if emptyIncident == len(incidents) {
							//grey
							incidentResponse["color"] = "#696e6c"
							incidentResponse["day"] = day
							incidentResponse["visible"] = true
						} else {
							if isIssueHappened {
								incidentResponse["color"] = "#FF6265"
								incidentResponse["day"] = day
								incidentResponse["visible"] = true
							} else {
								incidentResponse["color"] = "#38A976"
								incidentResponse["day"] = day
								incidentResponse["visible"] = true
							}

						}

					}
				} else {
					incidentResponse["color"] = "#c5c9c8"
					incidentResponse["day"] = day
					incidentResponse["visible"] = true

				}

			}
		} else {
			//grey
			incidentResponse["color"] = "#c5c9c8"
			incidentResponse["day"] = day
			incidentResponse["visible"] = true
		}

		incidentMap[day] = incidentResponse
		dayCount += 1
	}

	switch {
	case incidentType == "safety":
		response["incident"] = getSafetyChart(incidentMap)
	case incidentType == "quality":
		response["incident"] = getQualityChart(incidentMap)
	case incidentType == "delivery":
		response["incident"] = getDeliveryChart(incidentMap)
	case incidentType == "inventory":
		response["incident"] = getInventoryChart(incidentMap)
	case incidentType == "productivity":
		response["incident"] = getProductivityChart(incidentMap)
	}

	// response["incident"] = getSafetyChart(incidentMap)
	incidentMatrix := make(map[string]interface{})
	headerObjectList := make(map[string]interface{})
	matrixData := make([]map[string]interface{}, 0)
	for _, object := range *generalList {
		incidentInfo := make(map[string]interface{})
		json.Unmarshal(object.ObjectInfo, &incidentInfo)

		isNoWorkDay := incidentInfo["isNoWorkDay"].(bool)
		day := int(incidentInfo["day"].(float64))

		if incidentInfo["incidents"] != nil {
			incidents := incidentInfo["incidents"].([]interface{})
			if !isNoWorkDay && len(incidents) > 0 {
				for _, incident := range incidents {
					incidentObj := incident.(map[string]interface{})
					incidentMatrixResponse := make(map[string]interface{})
					if incidentObj["indicator"].(bool) {
						incidentId := int(incidentObj["id"].(float64))
						categoryInfo := findCategoryObject(generalCategoryList, incidentId)

						if categoryNameInterface, ok := categoryInfo["name"]; ok {
							categoryName := util.InterfaceToString(categoryNameInterface)
							categoryColor := categoryInfo["colorCode"].(string)
							var dayString string

							switch {
							case day == 1:
								dayString = strconv.Itoa(day) + "st"
							case day == 2:
								dayString = strconv.Itoa(day) + "nd"
							case day == 3:
								dayString = strconv.Itoa(day) + "rd"
							default:
								dayString = strconv.Itoa(day) + "th"
							}

							idStr := strconv.Itoa(categoryNameIdMap[categoryName])
							incidentMatrixResponse["id"] = idStr
							incidentMatrixResponse["day"] = dayString
							incidentMatrixResponse["description"] = categoryName

							headerObjectList[idStr] = categoryColor
							matrixData = append(matrixData, incidentMatrixResponse)

						}

					}
				}
			}
		}

	}
	incidentMatrix["currentRowCount"] = len(matrixData)
	incidentMatrix["data"] = matrixData

	header := make([]interface{}, 0)
	descriptionHeaderObject := map[string]interface{}{
		"name":         "Description",
		"type":         "text",
		"display":      true,
		"property":     "description",
		"routeEnabled": false,
		"columSize":    0}

	dayHeaderObject := map[string]interface{}{"name": "Day",
		"type":         "index_background_color",
		"display":      true,
		"property":     "day",
		"routeEnabled": false,
		"objectList":   headerObjectList,
		"columSize":    0}

	header = append(header, descriptionHeaderObject, dayHeaderObject)
	incidentMatrix["header"] = header
	incidentMatrix["totalRowCount"] = len(matrixData)

	response["matrix"] = incidentMatrix
	response["dailyTrend"] = v.getIncidentTrendResult(generalList)

	return response
}

func (v *IncidentService) getIncidentTrendResult(generalList *[]component.GeneralObject) map[string]interface{} {
	safetyTrendData := make([][]float64, 0)
	safetyTargetData := make([][]int, 0)
	for _, incidentObject := range *generalList {
		var incidentInfo = make(map[string]interface{})
		json.Unmarshal(incidentObject.ObjectInfo, &incidentInfo)

		totalCountStatForEachDay := make([]float64, 0)
		targetForEachDay := make([]int, 0)

		day := int(incidentInfo["day"].(float64))
		target := util.InterfaceToInt(incidentInfo["target"])
		totalCountStatForEachDay = append(totalCountStatForEachDay, float64(day))
		targetForEachDay = append(targetForEachDay, day)

		if !incidentInfo["isNoWorkDay"].(bool) {
			if incidentInfo["incidents"] != nil {
				incidentList := incidentInfo["incidents"].([]interface{})
				noIssueCount := 0
				totalNoOfDescription := len(incidentList)
				//Find number of ok counts in a day(noIssueCount)
				for _, incident := range incidentList {
					incidentInfo := incident.(map[string]interface{})
					incidentIndicator := incidentInfo["indicator"].(bool)

					if !incidentIndicator {
						noIssueCount += 1
					}
				}
				var totalCountPer float64
				if noIssueCount != 0 {
					totalCountPer = (float64(noIssueCount) / float64(totalNoOfDescription)) * 100
					totalCountPer = math.Round(totalCountPer*100) / 100

				}
				totalCountStatForEachDay = append(totalCountStatForEachDay, totalCountPer)
			} else {
				totalCountStatForEachDay = append(totalCountStatForEachDay, 0)
			}

		} else {
			totalCountStatForEachDay = append(totalCountStatForEachDay, 0)
		}
		targetForEachDay = append(targetForEachDay, target)

		safetyTrendData = append(safetyTrendData, [][]float64{totalCountStatForEachDay}...)
		safetyTargetData = append(safetyTargetData, [][]int{targetForEachDay}...)
	}

	seriesData := make(map[string]interface{})
	seriesData["name"] = "Daily Trend"
	seriesData["data"] = safetyTrendData

	seriesTargetData := make(map[string]interface{})
	seriesTargetData["name"] = "Daily Trend Target"
	seriesTargetData["data"] = safetyTargetData

	series := []interface{}{seriesData, seriesTargetData}
	safetyTrend := map[string]interface{}{"series": series}

	return safetyTrend
}

func getSafetyChart(safetyList map[int]interface{}) []interface{} {
	result := make([]interface{}, 0)

	matrixDayMap := map[int]int{2: 7, 3: 6, 4: 5, 5: 4, 7: 9, 8: 8, 11: 3, 12: 2, 13: 10, 18: 1, 19: 11, 25: 12, 31: 13, 32: 14, 33: 15, 34: 16, 39: 17, 40: 18, 41: 19, 42: 20, 48: 21, 54: 22, 60: 23, 61: 31, 62: 30, 65: 25, 66: 24, 68: 29, 69: 28, 70: 27, 71: 26}

	matrixIndex := 1
	dayIndex := 1

	for i := 1; i <= 12; i++ {
		rawData := make(map[string]interface{})
		for j := 1; j <= 6; j++ {
			var colData map[string]interface{}

			if (i == 1 && j == 1) || (i == 1 && j == 6) || (i == 2 && j == 3) || (i == 2 && j == 4) || (i == 3 && j == 2) || (i == 3 && j == 3) || (i == 3 && j == 4) || (i == 3 && j == 5) || (i == 4 && j == 2) || (i == 4 && j == 3) || (i == 4 && j == 4) || (i == 4 && j == 5) || (i == 4 && j == 6) || (i == 5 && j == 2) || (i == 5 && j == 3) || (i == 5 && j == 4) || (i == 5 && j == 5) || (i == 5 && j == 6) || (i == 6 && j == 5) || (i == 6 && j == 6) || (i == 7 && j == 1) || (i == 7 && j == 2) || (i == 8 && j == 1) || (i == 8 && j == 2) || (i == 8 && j == 3) || (i == 8 && j == 4) || (i == 8 && j == 5) || (i == 9 && j == 1) || (i == 9 && j == 2) || (i == 9 && j == 3) || (i == 9 && j == 4) || (i == 9 && j == 5) || (i == 10 && j == 1) || (i == 10 && j == 2) || (i == 10 && j == 3) || (i == 10 && j == 4) || (i == 10 && j == 5) || (i == 11 && j == 3) || (i == 11 && j == 4) || (i == 12 && j == 1) || (i == 12 && j == 6) {
				colData = map[string]interface{}{
					"visible": false,
					"day":     0,
					"color":   "transparent",
				}
			} else {
				dayIndex = matrixDayMap[matrixIndex]
				if safetyObject, ok := safetyList[dayIndex]; ok {
					colData = safetyObject.(map[string]interface{})
				} else {
					colData = map[string]interface{}{
						"visible": true,
						"day":     dayIndex,
						"color":   "#818589",
					}
				}
			}
			matrixIndex += 1

			colArrayData := []interface{}{colData}
			rawData["col"+strconv.Itoa(j)] = colArrayData
		}

		result = append(result, rawData)
	}

	return result
}

func getQualityChart(safetyList map[int]interface{}) []interface{} {
	result := make([]interface{}, 0)
	matrixDayMap := map[int]int{2: 1, 3: 2, 4: 3, 5: 4, 7: 30, 8: 31, 11: 5, 12: 6, 13: 29, 18: 7, 19: 28, 24: 8, 25: 27, 30: 9, 31: 26, 36: 10, 37: 25, 42: 11, 43: 24, 48: 12, 49: 23, 54: 13, 55: 22, 60: 14, 61: 21, 62: 20, 65: 16, 68: 19, 69: 18, 70: 17, 72: 15}

	matrixIndex := 1
	dayIndex := 1

	for i := 1; i <= 12; i++ {
		rawData := make(map[string]interface{})
		for j := 1; j <= 6; j++ {
			var colData map[string]interface{}

			if (i == 1 && j == 1) || (i == 1 && j == 6) || (i == 2 && j == 3) || (i == 2 && j == 4) || (i == 3 && j == 2) || (i == 3 && j == 3) || (i == 3 && j == 4) || (i == 3 && j == 5) || (i == 4 && j == 2) || (i == 4 && j == 3) || (i == 4 && j == 4) || (i == 4 && j == 5) || (i == 5 && j == 2) || (i == 5 && j == 3) || (i == 5 && j == 4) || (i == 5 && j == 5) || (i == 6 && j == 2) || (i == 6 && j == 3) || (i == 6 && j == 4) || (i == 6 && j == 5) || (i == 7 && j == 2) || (i == 7 && j == 3) || (i == 7 && j == 4) || (i == 7 && j == 5) || (i == 8 && j == 2) || (i == 8 && j == 3) || (i == 8 && j == 4) || (i == 8 && j == 5) || (i == 9 && j == 2) || (i == 9 && j == 3) || (i == 9 && j == 4) || (i == 9 && j == 5) || (i == 10 && j == 2) || (i == 10 && j == 3) || (i == 10 && j == 4) || (i == 10 && j == 5) || (i == 11 && j == 3) || (i == 11 && j == 4) || (i == 11 && j == 6) || (i == 12 && j == 1) || (i == 12 && j == 5) {
				colData = map[string]interface{}{
					"visible": false,
					"day":     0,
					"color":   "transparent",
				}
			} else {
				dayIndex = matrixDayMap[matrixIndex]
				if safetyObject, ok := safetyList[dayIndex]; ok {
					colData = safetyObject.(map[string]interface{})
				} else {
					colData = map[string]interface{}{
						"visible": true,
						"day":     dayIndex,
						"color":   "#818589",
					}
				}
			}
			matrixIndex += 1

			colArrayData := []interface{}{colData}
			rawData["col"+strconv.Itoa(j)] = colArrayData
		}

		result = append(result, rawData)
	}

	return result
}

func getDeliveryChart(safetyList map[int]interface{}) []interface{} {
	result := make([]interface{}, 0)

	matrixDayMap := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 7: 31, 11: 6, 12: 7, 13: 30, 18: 8, 19: 29, 24: 9, 25: 28, 30: 10, 31: 27, 36: 11, 37: 26, 42: 12, 43: 25, 48: 13, 49: 24, 54: 14, 55: 23, 60: 15, 61: 22, 66: 16, 67: 21, 68: 20, 69: 19, 70: 18, 71: 17}

	matrixIndex := 1
	dayIndex := 1

	for i := 1; i <= 12; i++ {
		rawData := make(map[string]interface{})
		for j := 1; j <= 6; j++ {
			var colData map[string]interface{}

			if (i == 1 && j == 6) || (i == 2 && j == 2) || (i == 2 && j == 3) || (i == 2 && j == 4) || (i == 3 && j == 2) || (i == 3 && j == 3) || (i == 3 && j == 4) || (i == 3 && j == 5) || (i == 4 && j == 2) || (i == 4 && j == 3) || (i == 4 && j == 4) || (i == 4 && j == 5) || (i == 5 && j == 2) || (i == 5 && j == 3) || (i == 5 && j == 4) || (i == 5 && j == 5) || (i == 6 && j == 2) || (i == 6 && j == 3) || (i == 6 && j == 4) || (i == 6 && j == 5) || (i == 7 && j == 2) || (i == 7 && j == 3) || (i == 7 && j == 4) || (i == 7 && j == 5) || (i == 8 && j == 2) || (i == 8 && j == 3) || (i == 8 && j == 4) || (i == 8 && j == 5) || (i == 9 && j == 2) || (i == 9 && j == 3) || (i == 9 && j == 4) || (i == 9 && j == 5) || (i == 10 && j == 2) || (i == 10 && j == 3) || (i == 10 && j == 4) || (i == 10 && j == 5) || (i == 11 && j == 2) || (i == 11 && j == 3) || (i == 11 && j == 4) || (i == 11 && j == 5) || (i == 12 && j == 6) {
				colData = map[string]interface{}{
					"visible": false,
					"day":     0,
					"color":   "transparent",
				}
			} else {
				dayIndex = matrixDayMap[matrixIndex]
				if safetyObject, ok := safetyList[dayIndex]; ok {
					colData = safetyObject.(map[string]interface{})
				} else {
					colData = map[string]interface{}{
						"visible": true,
						"day":     dayIndex,
						"color":   "#818589",
					}
				}
			}
			matrixIndex += 1

			colArrayData := []interface{}{colData}
			rawData["col"+strconv.Itoa(j)] = colArrayData
		}

		result = append(result, rawData)
	}

	return result
}

func getInventoryChart(safetyList map[int]interface{}) []interface{} {
	result := make([]interface{}, 0)
	matrixDayMap := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 9: 8, 10: 7, 15: 10, 16: 9, 21: 12, 22: 11, 27: 14, 28: 13, 33: 16, 34: 15, 39: 18, 40: 17, 45: 20, 46: 19, 51: 22, 52: 21, 57: 24, 58: 23, 63: 26, 64: 25, 67: 27, 68: 28, 69: 29, 70: 30, 71: 31}

	dayIndex := 1
	matrixIndex := 1

	for i := 1; i <= 12; i++ {
		rawData := make(map[string]interface{})
		for j := 1; j <= 6; j++ {
			var colData map[string]interface{}

			if (i == 2 && j == 1) || (i == 2 && j == 2) || (i == 2 && j == 5) || (i == 2 && j == 6) || (i == 3 && j == 1) || (i == 3 && j == 2) || (i == 3 && j == 5) || (i == 3 && j == 6) || (i == 4 && j == 1) || (i == 4 && j == 2) || (i == 4 && j == 5) || (i == 4 && j == 6) || (i == 5 && j == 1) || (i == 5 && j == 2) || (i == 5 && j == 5) || (i == 5 && j == 6) || (i == 6 && j == 1) || (i == 6 && j == 2) || (i == 6 && j == 5) || (i == 6 && j == 6) || (i == 7 && j == 1) || (i == 7 && j == 2) || (i == 7 && j == 5) || (i == 7 && j == 6) || (i == 8 && j == 1) || (i == 8 && j == 2) || (i == 8 && j == 5) || (i == 8 && j == 6) || (i == 9 && j == 1) || (i == 9 && j == 2) || (i == 9 && j == 5) || (i == 9 && j == 6) || (i == 10 && j == 1) || (i == 10 && j == 2) || (i == 10 && j == 5) || (i == 10 && j == 6) || (i == 11 && j == 1) || (i == 11 && j == 2) || (i == 11 && j == 5) || (i == 11 && j == 6) {
				colData = map[string]interface{}{
					"visible": false,
					"day":     0,
					"color":   "transparent",
				}
			} else if i == 12 && j == 6 {
				colData = map[string]interface{}{
					"visible": true,
					"day":     0,
					"color":   "#818589",
				}
			} else {
				dayIndex = matrixDayMap[matrixIndex]
				if safetyObject, ok := safetyList[dayIndex]; ok {
					colData = safetyObject.(map[string]interface{})
				} else {
					colData = map[string]interface{}{
						"visible": true,
						"day":     dayIndex,
						"color":   "#818589",
					}
				}
			}
			matrixIndex += 1

			colArrayData := []interface{}{colData}
			rawData["col"+strconv.Itoa(j)] = colArrayData
		}

		result = append(result, rawData)
	}

	return result
}

func getProductivityChart(safetyList map[int]interface{}) []interface{} {
	result := make([]interface{}, 0)
	matrixDayMap := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 7: 31, 11: 6, 12: 7, 13: 30, 18: 8, 19: 29, 24: 9, 25: 28, 30: 10, 31: 27, 36: 11, 37: 26, 38: 15, 39: 14, 40: 13, 41: 12, 43: 25, 44: 16, 49: 24, 50: 17, 55: 23, 56: 18, 61: 22, 62: 19, 67: 21, 68: 20}

	dayIndex := 1
	matrixIndex := 1

	for i := 1; i <= 12; i++ {
		rawData := make(map[string]interface{})
		for j := 1; j <= 6; j++ {
			var colData map[string]interface{}

			if (i == 1 && j == 6) || (i == 2 && j == 2) || (i == 2 && j == 3) || (i == 2 && j == 4) || (i == 3 && j == 2) || (i == 3 && j == 3) || (i == 3 && j == 4) || (i == 3 && j == 5) || (i == 4 && j == 2) || (i == 4 && j == 3) || (i == 4 && j == 4) || (i == 4 && j == 5) || (i == 5 && j == 2) || (i == 5 && j == 3) || (i == 5 && j == 4) || (i == 5 && j == 5) || (i == 6 && j == 2) || (i == 6 && j == 3) || (i == 6 && j == 4) || (i == 6 && j == 5) || (i == 7 && j == 6) || (i == 8 && j == 3) || (i == 8 && j == 4) || (i == 8 && j == 5) || (i == 8 && j == 6) || (i == 9 && j == 3) || (i == 9 && j == 4) || (i == 9 && j == 5) || (i == 9 && j == 6) || (i == 10 && j == 3) || (i == 10 && j == 4) || (i == 10 && j == 5) || (i == 10 && j == 6) || (i == 11 && j == 3) || (i == 11 && j == 4) || (i == 11 && j == 5) || (i == 11 && j == 6) || (i == 12 && j == 3) || (i == 12 && j == 4) || (i == 12 && j == 5) || (i == 12 && j == 6) {
				colData = map[string]interface{}{
					"visible": false,
					"day":     0,
					"color":   "transparent",
				}
			} else {
				dayIndex = matrixDayMap[matrixIndex]
				if safetyObject, ok := safetyList[dayIndex]; ok {
					colData = safetyObject.(map[string]interface{})
				} else {
					colData = map[string]interface{}{
						"visible": true,
						"day":     dayIndex,
						"color":   "#818589",
					}
				}
			}
			matrixIndex += 1

			colArrayData := []interface{}{colData}
			rawData["col"+strconv.Itoa(j)] = colArrayData
		}

		result = append(result, rawData)
	}

	return result
}

func findIncidentObjectGivenDay(generalList *[]component.GeneralObject, searchValue int) map[string]interface{} {
	for _, object := range *generalList {
		incidentInfo := make(map[string]interface{})
		json.Unmarshal(object.ObjectInfo, &incidentInfo)
		if util.InterfaceToInt(incidentInfo["day"]) == searchValue {
			return incidentInfo
		}

	}
	return nil
}
