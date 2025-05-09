package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (v *IncidentService) getSummaryResponse(projectId string, dateQuery string, department string) map[string]interface{} {
	var decodedResponse = make(map[string]interface{})

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	dateQueryList := strings.Split(dateQuery, "/")

	dateCondition := " object_info->>'$.department' = " + department + " AND object_info->>'$.month' = " + dateQueryList[0] + " AND object_info->>'$.year' =" + dateQueryList[1] + " "
	orderCondition := "object_info -> '$.day'"

	listOfSafetyObjects, err := GetConditionalObjectsOrderBy(dbConnection, IncidentSafetyTable, dateCondition, orderCondition)
	listOfQualityObjects, err := GetConditionalObjectsOrderBy(dbConnection, IncidentQualityTable, dateCondition, orderCondition)
	listOfDeliveryObjects, err := GetConditionalObjectsOrderBy(dbConnection, IncidentDeliveryTable, dateCondition, orderCondition)
	listOfInventoryObjects, err := GetConditionalObjectsOrderBy(dbConnection, IncidentInventoryTable, dateCondition, orderCondition)
	listOfProductivityObjects, err := GetConditionalObjectsOrderBy(dbConnection, IncidentProductivityTable, dateCondition, orderCondition)

	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		return decodedResponse
	}

	//Get incident category
	condition := " object_info->>'$.department' = " + department + " AND object_info->>'$.objectStatus' = 'Active'"
	incidentSafetyCategoryObjects, err := GetConditionalObjects(dbConnection, IncidentSafetyCategoryTable, condition)
	incidentQualityCategoryObjects, err := GetConditionalObjects(dbConnection, IncidentQualityCategoryTable, condition)
	incidentDeliveryCategoryObjects, err := GetConditionalObjects(dbConnection, IncidentDeliveryCategoryTable, condition)
	incidentInventoryCategoryObjects, err := GetConditionalObjects(dbConnection, IncidentInventoryCategoryTable, condition)
	incidentProductivityCategoryObjects, err := GetConditionalObjects(dbConnection, IncidentProductivityCategoryTable, condition)

	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		return decodedResponse
	}

	safetyOverview := v.getIncidentOverviewResult(listOfSafetyObjects, incidentSafetyCategoryObjects, "safety")
	qualityOverview := v.getIncidentOverviewResult(listOfQualityObjects, incidentQualityCategoryObjects, "quality")
	deliveryOverview := v.getIncidentOverviewResult(listOfDeliveryObjects, incidentDeliveryCategoryObjects, "delivery")
	inventoryOverview := v.getIncidentOverviewResult(listOfInventoryObjects, incidentInventoryCategoryObjects, "inventory")
	productivityOverview := v.getIncidentOverviewResult(listOfProductivityObjects, incidentProductivityCategoryObjects, "productivity")

	decodedResponse["delivery"] = deliveryOverview["incident"]
	decodedResponse["deliveryDailyTrend"] = deliveryOverview["dailyTrend"]
	decodedResponse["deliveryMatrix"] = deliveryOverview["matrix"]

	decodedResponse["inventory"] = inventoryOverview["incident"]
	decodedResponse["inventoryDailyTrend"] = inventoryOverview["dailyTrend"]
	decodedResponse["inventoryMatrix"] = inventoryOverview["matrix"]

	decodedResponse["productivity"] = productivityOverview["incident"]
	decodedResponse["productivityDailyTrend"] = productivityOverview["dailyTrend"]
	decodedResponse["productivityMatrix"] = productivityOverview["matrix"]

	decodedResponse["quality"] = qualityOverview["incident"]
	decodedResponse["qualityDailyTrend"] = qualityOverview["dailyTrend"]
	decodedResponse["qualityMatrix"] = qualityOverview["matrix"]

	decodedResponse["safety"] = safetyOverview["incident"]
	decodedResponse["safetyDailyTrend"] = safetyOverview["dailyTrend"]
	decodedResponse["safetyMatrix"] = safetyOverview["matrix"]

	// json.Unmarshal([]byte(predefinedResponse), &decodedResponse)
	return decodedResponse
}

// need to have query param which accepts  monthYear , 09/2022 or sep - 2022
func (v *IncidentService) getSQDIPOverview(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	dateQuery := ctx.Query("dateQuery")
	department := ctx.Query("department")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var safetyResponse = make(map[string]interface{})
	var dateQueryList []string
	if dateQuery != "" {
		dateQueryList = strings.Split(dateQuery, "/")
	} else {
		v.BaseService.Logger.Error("request doesn't have date query")
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("Please send date query"), ErrorGettingIndividualObjectInformation)
		return
	}
	var incidentTable string
	var incidentCategoryTable string

	switch {
	case componentName == IncidentSafetyComponent:
		incidentTable = IncidentSafetyTable
		incidentCategoryTable = IncidentSafetyCategoryTable
	case componentName == IncidentQualityComponent:
		incidentTable = IncidentQualityTable
		incidentCategoryTable = IncidentQualityCategoryTable
	case componentName == IncidentDeliveryComponent:
		incidentTable = IncidentDeliveryTable
		incidentCategoryTable = IncidentDeliveryCategoryTable
	case componentName == IncidentInventoryComponent:
		incidentTable = IncidentInventoryTable
		incidentCategoryTable = IncidentInventoryCategoryTable
	case componentName == IncidentProductivityComponent:
		incidentTable = IncidentProductivityTable
		incidentCategoryTable = IncidentProductivityCategoryTable
	}

	//Get list of incident objects
	dateCondition := " object_info ->> '$.department' = " + department + " AND object_info->>'$.month' = " + dateQueryList[0] + " AND object_info->>'$.year' =" + dateQueryList[1] + " "
	orderCondition := "object_info -> '$.day'"
	listOfObjects, err := GetConditionalObjectsOrderBy(dbConnection, incidentTable, dateCondition, orderCondition)

	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	fmt.Println("Step 1")
	if len(*listOfObjects) == 0 {
		pieDefaultResponse := map[string]interface{}{"data": make([]int, 0), "name": ""}
		safetyResponse["pieDistribution"] = []interface{}{pieDefaultResponse}
		safetyResponse["safetyPercentage"] = 0
		safetyResponse["safetyTrend"] = map[string]interface{}{"series": make([]interface{}, 0)}
		safetyResponse["safetyParetoAnalysis"] = map[string]interface{}{"series": make([]interface{}, 0)}
		safetyResponse["safetyMatrixDayView"] = map[string]interface{}{"data": make([]interface{}, 0), "header": make([]interface{}, 0), "totalRowCount": 0}
		ctx.JSON(http.StatusOK, safetyResponse)
		return
	}
	fmt.Println("Step 2")
	//Get incident category
	categoryCondition := " object_info ->> '$.department' = " + department + " AND object_info->>'$.objectStatus' = 'Active'"
	incidentCategoryObjects, err := GetConditionalObjects(dbConnection, incidentCategoryTable, categoryCondition)

	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	fmt.Println("Step 3")
	if len(*incidentCategoryObjects) == 0 {
		pieDefaultResponse := map[string]interface{}{"data": make([]int, 0), "name": ""}
		safetyResponse["pieDistribution"] = []interface{}{pieDefaultResponse}
		safetyResponse["safetyPercentage"] = 0
		safetyResponse["safetyTrend"] = map[string]interface{}{"series": make([]interface{}, 0)}
		safetyResponse["safetyParetoAnalysis"] = map[string]interface{}{"series": make([]interface{}, 0)}
		safetyResponse["safetyMatrixDayView"] = map[string]interface{}{"data": make([]interface{}, 0), "header": make([]interface{}, 0), "totalRowCount": 0}
		ctx.JSON(http.StatusOK, safetyResponse)
		// response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("requested record is empty"), ErrorGettingIndividualObjectInformation)
		return
	}

	//Get list of categories
	categoryMap := make(map[string]int, 0)
	for _, category := range *incidentCategoryObjects {
		var safetyCategoryInfo = make(map[string]interface{})
		json.Unmarshal(category.ObjectInfo, &safetyCategoryInfo)
		categoryName := util.InterfaceToString(safetyCategoryInfo["name"])
		categoryMap[categoryName] = 0
	}
	fmt.Println("Step 4")
	//creating object for pieDistribution############################################
	//Creating master data stats
	//This will return table stats occurStatForEachDescription
	//{"WIP": { "TotalOccurence": 1, "TotalNormal": 12}}
	occurStatForEachDescription := make(map[string]interface{}, 0)
	//Total count stat
	// safetyTrendData := make([][]int, 0)
	totalDefects := 0

	for _, safetyObject := range *listOfObjects {
		var incidentSafetyInfo = make(map[string]interface{})
		json.Unmarshal(safetyObject.ObjectInfo, &incidentSafetyInfo)

		// totalCountStatForEachDay := make([]int, 0)
		// day := util.InterfaceToInt(incidentSafetyInfo["day"])
		// totalCountStatForEachDay = append(totalCountStatForEachDay, day)

		if !incidentSafetyInfo["isNoWorkDay"].(bool) {
			incidentList := incidentSafetyInfo["incidents"].([]interface{})
			noIssueCount := 0
			//Loop through the incident array
			for _, object := range incidentList {
				incidentObject := object.(map[string]interface{})

				categoryId := util.InterfaceToInt(incidentObject["id"])
				categoryInfo := findCategoryObject(incidentCategoryObjects, categoryId)
				if categoryNameInterface, ok := categoryInfo["name"]; ok {
					var categoryName = util.InterfaceToString(categoryNameInterface)
					categoryMap[categoryName] = categoryMap[categoryName] + 1
					incidentIndicator := incidentObject["indicator"].(bool)

					//Creating occurrence stats map####################################################################
					if val, ok := occurStatForEachDescription[categoryName]; ok {
						tableMap := val.(map[string]interface{})
						if !incidentIndicator {
							tableMap["TotalOccurrence"] = util.InterfaceToInt(tableMap["TotalOccurrence"]) + 1

						} else {
							tableMap["TotalNormal"] = util.InterfaceToInt(tableMap["TotalNormal"]) + 1
							totalDefects += 1
						}
					} else {
						tableMap := make(map[string]interface{}, 0)

						if !incidentIndicator {
							tableMap["TotalOccurrence"] = 1
							tableMap["TotalNormal"] = 0

						} else {
							tableMap["TotalNormal"] = 1
							tableMap["TotalOccurrence"] = 0
							totalDefects += 1
						}
						occurStatForEachDescription[categoryName] = tableMap

					}
					//###################################################################################

					//Create total stats########################################################################

					if incidentIndicator {
						noIssueCount += 1

					}
				}
			}
			// totalCountStatForEachDay = append(totalCountStatForEachDay, noIssueCount)

		}
		// else {
		// 	totalCountStatForEachDay = append(totalCountStatForEachDay, 0)
		// }
		// safetyTrendData = append(safetyTrendData, [][]int{totalCountStatForEachDay}...)
	}
	fmt.Println("Step 5")
	pieDistributionArray := make([]map[string]interface{}, 0)
	pieDistributionObject := make(map[string]interface{}, 0)
	pieDistributionObject["name"] = "Category Distribution"
	pieDistributionObject["colorByPoint"] = true

	pieDistributionData := make([]map[string]interface{}, 0)

	for key, element := range categoryMap {
		pieDataObject := make(map[string]interface{}, 0)
		pieDataObject["name"] = key
		pieDataObject["y"] = element

		pieDistributionData = append(pieDistributionData, pieDataObject)
	}
	fmt.Println("Step 6")
	pieDistributionObject["data"] = pieDistributionData
	pieDistributionArray = append(pieDistributionArray, pieDistributionObject)
	safetyResponse["pieDistribution"] = pieDistributionArray
	//##############################################################################

	//safetyPercentage###########################################################
	occurrence := 0
	normal := 0
	for _, element := range occurStatForEachDescription {
		stats := element.(map[string]interface{})
		occurrence = occurrence + util.InterfaceToInt(stats["TotalOccurrence"])
		normal = normal + util.InterfaceToInt(stats["TotalNormal"])
	}
	fmt.Println("Step 7")
	var occurrenceRatio float64

	if (occurrence + normal) > 0 {
		occurrenceRatio = float64(normal) / float64(occurrence+normal)
	}
	// safetyPercentage := (1 - occurrenceRatio) * 100
	safetyPercentage := occurrenceRatio * 100

	safetyResponse["safetyPercentage"] = math.Round(100 - math.Round(safetyPercentage*100)/100)

	//##############################################################################

	//safetyTrend######################################################################

	safetyResponse["safetyTrend"] = v.getIncidentTrendResult(listOfObjects)
	//###################################################################################

	//safetyPareto######################################################################
	paretoXAxis := make([]string, 0)
	for key, _ := range categoryMap {
		paretoXAxis = append(paretoXAxis, key)
	}
	fmt.Println("Step 8")
	sortedOccurStatForEachDescription := sortDescriptionBasedOnOccurence(occurStatForEachDescription, paretoXAxis)
	defectList := make([]int, 0)

	//category list must be same order as data
	paretoXAxis = make([]string, 0)
	fmt.Println("Step 9")
	cumulativeDefects := make([]float64, 0)
	previousDefectValue := 0

	for i := len(sortedOccurStatForEachDescription) - 1; i >= 0; i-- {
		description := sortedOccurStatForEachDescription[i]
		tableMap := description.(map[string]interface{})
		paretoXAxis = append(paretoXAxis, tableMap["key"].(string))

		totalOccurrenceForDescription := tableMap["TotalNormal"].(int)
		defectList = append(defectList, totalOccurrenceForDescription)

		cumulativeDefectCount := previousDefectValue + totalOccurrenceForDescription

		var cumulativeDefectPercentApprox float64

		if totalDefects != 0 {
			cumulativeDefectPercent := (float64(cumulativeDefectCount) / float64(totalDefects)) * 100
			cumulativeDefectPercentApprox = math.Round(cumulativeDefectPercent*100) / 100
		}

		cumulativeDefects = append(cumulativeDefects, cumulativeDefectPercentApprox)

		previousDefectValue = cumulativeDefectCount
	}
	fmt.Println("Step 10")
	// var intialCumulativeDefect float64
	// if totalDefects != 0 {
	// 	intialCumulativeDefect = math.Round((float64(defectList[0])/float64(totalDefects))*100) / 100
	// }

	// cumulativeDefects = append(cumulativeDefects, intialCumulativeDefect*100)

	// cumulativeData := defectList[0]
	// for i := 1; i < len(paretoXAxis); i++ {
	// 	cumulativeData += defectList[i]
	// 	calc := math.Round((float64(cumulativeData)/float64(totalDefects))*100) / 100
	// 	cumulativeDefects = append(cumulativeDefects, calc*100)
	// }

	defectMap := make(map[string]interface{}, 0)
	defectMap["name"] = "No of Defects"
	defectMap["data"] = defectList

	cumulativeDefectMap := make(map[string]interface{}, 0)
	cumulativeDefectMap["name"] = "Cumulative Defect in %"
	cumulativeDefectMap["data"] = cumulativeDefects

	paretoSeries := make([]interface{}, 0)
	paretoSeries = append(paretoSeries, defectMap)
	paretoSeries = append(paretoSeries, cumulativeDefectMap)

	safetyParetoAnalysis := make(map[string]interface{})
	safetyParetoAnalysis["xAxis"] = paretoXAxis
	safetyParetoAnalysis["series"] = paretoSeries
	safetyResponse["safetyParetoAnalysis"] = safetyParetoAnalysis
	//###################################################################################
	safetyResponse["safetyMatrixDayView"] = v.calculateMatrixDayView(dateQuery, listOfObjects, incidentCategoryObjects)
	fmt.Println("Step 11")
	// json.Unmarshal([]byte(safetyHardCodedResponse), &safetyResponse)
	ctx.JSON(http.StatusOK, safetyResponse)
}

func findCategoryObject(generalList *[]component.GeneralObject, searchValue int) map[string]interface{} {
	for _, object := range *generalList {
		if object.Id == searchValue {
			var safetyCategoryInfo = make(map[string]interface{})
			json.Unmarshal(object.ObjectInfo, &safetyCategoryInfo)
			return safetyCategoryInfo
		}
	}
	return nil
}

func findIncidentObject(incidentArray []interface{}, searchValue int) map[string]interface{} {
	for _, object := range incidentArray {
		incidentObject := object.(map[string]interface{})
		if util.InterfaceToInt(incidentObject["id"]) == searchValue {
			return incidentObject
		}
	}
	return nil
}

func sortDescriptionBasedOnOccurence(safetyTable map[string]interface{}, category []string) []interface{} {
	sliceMap := convertMapToSlice(safetyTable)

	sort.SliceStable(sliceMap, func(i, j int) bool {
		firstValue := sliceMap[i].(map[string]interface{})
		secondValue := sliceMap[j].(map[string]interface{})

		return util.InterfaceToInt(firstValue["TotalNormal"]) < util.InterfaceToInt(secondValue["TotalNormal"])
	})

	return sliceMap
}

// createNewResource ShowAccount godoc
// @Summary create new resource
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   recordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/records [post]
func (v *IncidentService) createNewResource(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var err error

	updatedRequest := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	// here we should do the validation
	err = v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}
	rawCreateRequest, _ := json.Marshal(updatedRequest)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}

	err, _ = Create(dbConnection, targetTable, object)
	// TODO we need to create the record trail
	if err != nil {
		v.BaseService.Logger.Error("error creating resource", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
		return
	}

	switch componentName {
	case IncidentSafetyCategoryComponent:

	}
	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully created",
		Error:   0,
	})
}

func convertMapToSlice(occurStatForEachDescription map[string]interface{}) []interface{} {
	result := make([]interface{}, 0)
	for key, description := range occurStatForEachDescription {
		object := description.(map[string]interface{})
		object["key"] = key
		result = append(result, object)
	}

	return result
}
