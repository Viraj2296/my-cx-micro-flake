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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type IncidentMatrixResponse struct {
	IsNoWorkDay component.RecordInfo   `json:"IsNoWorkDay"`
	Incidents   []component.RecordInfo `json:"incidents"`
	Target      component.RecordInfo   `json:"target"`
}

type Incident struct {
	Id        int    `json:"id"`
	Indicator bool   `json:"indicator"`
	Name      string `json:"name"`
}

func (v *IncidentService) getUpdateMatrixInfo(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	dateQuery := ctx.Query("dateQuery")
	department := ctx.Query("department")

	var dateList []string
	if dateQuery == "" {
		currentDate := util.GetCurrentDate()
		dateList = strings.Split(currentDate, "-")
	} else {
		dateList = strings.Split(dateQuery, "/")
	}

	localZone, _ := time.LoadLocation("Asia/Singapore")
	currentYear, currentMonth, currentDay := time.Now().In(localZone).Date()

	if dateList[0] > strconv.Itoa(currentDay) && (dateList[1] > strconv.Itoa(int(currentMonth))) && (dateList[2] > strconv.Itoa(currentYear)) {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("requested date in future"), ErrorGettingIndividualObjectInformation)
		return
	}

	// if not sent , current date
	condition := " object_info ->> '$.department' = " + department + " AND object_info ->> '$.year' = " + dateList[2] + " AND object_info ->> '$.month' = " + dateList[1] + " AND object_info ->> '$.day' = " + dateList[0]
	listOfObjectInterface, err := GetConditionalObjects(dbConnection, targetTable, condition)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)

	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)

	var incidentCategoryObjects *[]component.GeneralObject
	var incidentType int
	categoryCondition := " object_info ->> '$.department' = " + department

	switch {
	case componentName == IncidentSafetyComponent:
		incidentType = Safety
		incidentCategoryObjects, _ = GetConditionalObjects(dbConnection, IncidentSafetyCategoryTable, categoryCondition)
	case componentName == IncidentQualityComponent:
		incidentType = Quality
		incidentCategoryObjects, _ = GetConditionalObjects(dbConnection, IncidentQualityCategoryTable, categoryCondition)
	case componentName == IncidentDeliveryComponent:
		incidentType = Delivery
		incidentCategoryObjects, _ = GetConditionalObjects(dbConnection, IncidentDeliveryCategoryTable, categoryCondition)
	case componentName == IncidentInventoryComponent:
		incidentType = Inventory
		incidentCategoryObjects, _ = GetConditionalObjects(dbConnection, IncidentInventoryCategoryTable, categoryCondition)
	case componentName == IncidentProductivityComponent:
		incidentType = Productivity
		incidentCategoryObjects, _ = GetConditionalObjects(dbConnection, IncidentProductivityCategoryTable, categoryCondition)
	default:
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("requested incident category not available"), ErrorGettingIndividualObjectInformation)
		return
	}

	//If the data isn't available the we shoud send schema to insert new data
	if len(*listOfObjectInterface) == 0 {
		// incidentInfo := IncidentInfo{}
		defaultTarget := 0
		targetCondition := " object_info ->> '$.type' = " + strconv.Itoa(incidentType) + " AND object_info ->> '$.department' = " + department + " "
		incidentTargetObjects, _ := GetConditionalObjects(dbConnection, IncidentTargetTable, targetCondition)

		if len(*incidentTargetObjects) != 0 {
			targetInfo := make(map[string]interface{})
			json.Unmarshal((*incidentTargetObjects)[0].ObjectInfo, &targetInfo)

			defaultTarget = util.InterfaceToInt(targetInfo["target"])
		}

		incidentArray := make([]interface{}, 0)

		for _, category := range *incidentCategoryObjects {
			categoryInfo := make(map[string]interface{})
			json.Unmarshal(category.ObjectInfo, &categoryInfo)
			if util.InterfaceToString(categoryInfo["objectStatus"]) == common.ObjectStatusActive {
				incidentArray = append(incidentArray, map[string]interface{}{"id": category.Id, "indicator": false, "name": categoryInfo["name"]})
			}

		}

		incidentObject := map[string]interface{}{"data": incidentArray, "isEdit": true, "isExternal": false, "type": "object_array"}
		targetObject := map[string]interface{}{"isEdit": true, "value": defaultTarget, "isExternal": false, "type": "number"}

		newRecordResponse["incidents"] = incidentObject
		newRecordResponse["target"] = targetObject

		ctx.JSON(http.StatusOK, newRecordResponse)
		return
	}

	generalObject := (*listOfObjectInterface)[0]
	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", generalObject.Id)

	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, generalObject.Id, componentName, rawJSONObject)
	recordInfo := response["incidents"].(component.RecordInfo)
	incidentsArray := recordInfo.Data.([]interface{})
	fmt.Println("incidentsArray: ", incidentsArray)
	arrayOfIncidents := make([]interface{}, 0)
	if len(incidentsArray) != 0 {
		// should be new one, then send all the category available

		for _, incidentCategory := range *incidentCategoryObjects {
			incidentCategoryInfo := make(map[string]interface{})
			json.Unmarshal(incidentCategory.ObjectInfo, &incidentCategoryInfo)
			// incidentObject := incidentObjectInterface.(map[string]interface{})

			incidentObject := findIncidentObject(incidentsArray, incidentCategory.Id)

			// incidentCategory := findCategoryObject(incidentSafetyCategoryObjects, util.InterfaceToInt(incidentObject["id"]))

			if incidentObject != nil {
				if util.InterfaceToString(incidentCategoryInfo["objectStatus"]) == common.ObjectStatusActive {

					arrayOfIncidents = append(arrayOfIncidents, Incident{
						Id:        util.InterfaceToInt(incidentObject["id"]),
						Indicator: incidentObject["indicator"].(bool),
						Name:      util.InterfaceToString(incidentCategoryInfo["name"]),
					})
				}
			} else {
				if util.InterfaceToString(incidentCategoryInfo["objectStatus"]) == common.ObjectStatusActive {
					arrayOfIncidents = append(arrayOfIncidents, Incident{
						Id:        incidentCategory.Id,
						Indicator: false,
						Name:      util.InterfaceToString(incidentCategoryInfo["name"]),
					})
				}

			}
		}
	}

	recordInfo.Data = arrayOfIncidents
	response["incidents"] = recordInfo

	ctx.JSON(http.StatusOK, response)

}

func (v *IncidentService) updateMatrixInfo(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	year := util.InterfaceToInt(updateRequest["year"])
	month := util.InterfaceToInt(updateRequest["month"])
	day := util.InterfaceToInt(updateRequest["day"])
	department := util.InterfaceToInt(updateRequest["department"])

	condition := "object_info ->> '$.year' = " + strconv.Itoa(year) + " AND object_info ->> '$.month' = " + strconv.Itoa(month) + " AND object_info ->> '$.day' = " + strconv.Itoa(day) + " AND object_info ->> '$.department' = " + strconv.Itoa(department)
	orderCondition := "object_info -> '$.day'"
	listOfObjectInterface, err := GetConditionalObjectsOrderBy(dbConnection, targetTable, condition, orderCondition)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	//If the incident not available we have to insert data
	if len(*listOfObjectInterface) == 0 {
		// updatedRequest := is.ComponentManager.PreprocessCreateRequestFields(updateRequest, componentName)

		rawCreateRequest, _ := json.Marshal(updateRequest)
		preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
		object := component.GeneralObject{
			ObjectInfo: preprocessedRequest,
		}

		err, _ = Create(dbConnection, targetTable, object)

		if err != nil {
			fmt.Println("error creating ...", err.Error())
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error creating resource information"), ErrorUpdatingObjectInformation)
			return
		}

		ctx.JSON(http.StatusOK, component.GeneralResponse{
			Message: "Successfully created",
			Error:   0,
		})
		return

	}

	generalObject := (*listOfObjectInterface)[0]

	//Adding update preprocess request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, generalObject.ObjectInfo, componentName)

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, generalObject.Id, updatingData)
	if err != nil {
		fmt.Println("error updating ...", err.Error())
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating resource information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully updated",
		Error:   0,
	})
}

// func departmentWiseFilter(incidentObjects *[]component.GeneralObject, departmentId string, dbConnection *gorm.DB, componentName string) *[]component.GeneralObject {
// 	filteredObjects := make([]component.GeneralObject, 0)
// 	var targetTableName string

// 	switch{
// 	case componentName == IncidentSafetyComponent:
// 		targetTableName =
// 	}

// 	condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, '$.department')) = " + departmentId
// 	listOfObjectInterface, err := GetConditionalObjects(dbConnection, targetTable, condition)

// 	for _, incident := range *incidentObjects {
// 		incidentInfo := make(map[string]interface{})
// 		json.Unmarshal(incident.ObjectInfo, &incidentInfo)

// 		for _,
// 	}

// 	return &filteredObjects
// }
