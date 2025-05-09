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

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// getDashboardSnapshot ShowAccount godoc
// @Summary Get the record form data to facilitate the update
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags Dashboards
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/machine_overview [get]
func (v *IncidentService) getIncidentDashboard(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	dateQuery := ctx.Query("dateQuery")
	department := ctx.Query("department")
	summaryResponse := v.getSummaryResponse(projectId, dateQuery, department)
	ctx.JSON(http.StatusOK, summaryResponse)
}

// getNewRecord ShowAccount godoc
// @Summary Get the new record based on record schema
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/new_record [get]
func (v *IncidentService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

// getMachineRegister ShowAccount godoc
// @Summary Get all the machine related information
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/records [get]
func (v *IncidentService) getObjects(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	condition := ctx.Query("condition")
	outFields := ctx.Query("out_fields")
	format := ctx.Query("format")
	searchFields := ctx.Query("search")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	//Have to next flag
	isNext := true

	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
	if searchFields != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)
		baseCondition := component.TableCondition(offsetValue, fields, values, condition)
		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		searchWithBaseQuery := searchQuery + " AND " + baseCondition
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		fmt.Println("targetTable:", targetTable, "dbconnection:", dbConnection)
		listOfObjects, err = GetObjects(dbConnection, targetTable)
		if err != nil {
			fmt.Println("error :", err.Error())
			totalRecords = 0
		} else {
			totalRecords = int64(len(*listOfObjects))
		}

	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			condition = component.TableCondition(offsetValue, fields, values, condition)
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, condition, limitVal)

			currentRecordCount := len(*listOfObjects)

			if currentRecordCount < limitVal {
				isNext = false
			} else if currentRecordCount == limitVal {
				andClauses := strings.Split(condition, "AND")
				var totalRecordObjects *[]component.GeneralObject
				if len(andClauses) > 1 {
					totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, condition)

				} else {
					totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
				}
				lenTotalRecord := len(*totalRecordObjects)
				if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
					isNext = false
				}
			}
		}

	}
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		userId := common.GetUserId(ctx)
		zone := getUserTimezone(userId)
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)

		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

// deleteResource ShowAccount godoc
// @Summary Delete the any given resource using resource id
// @Description based on user permission, user can perform delete operations
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId} [delete]
func (v *IncidentService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordIdString(ctx)
	recordIdInt, _ := strconv.Atoi(recordId)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, recordIdInt)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	// check constraints and proceed
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	for _, constraint := range listOfConstraints {

		referenceComponent := constraint.Reference
		referenceField := constraint.ReferenceProperty
		referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)

		ArchiveReferenceObjects(dbConnection, referenceTable, referenceField, recordIdInt)
	}
	v.ComponentManager.ProcessDeleteDependencyInjection(dbConnection, recordIdInt, componentName)
	updatedObjectInfo := common.UpdateMetaInfoFromSerializedObject(generalObject.ObjectInfo, ctx)
	err = ArchiveObject(dbConnection, targetTable, component.GeneralObject{Id: generalObject.Id, ObjectInfo: updatedObjectInfo})

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
		return
	}

	//recordMessage := componentName + " is deleted"
	//err = CreateUserRecordTrail(ms, projectId, recordId, componentName, recordMessage, nil, nil)
	//if err != nil {
	//	ms.BaseService.Logger.Error("error in create record trail", zap.String("error", err.Error()))
	//}
	ctx.Status(http.StatusNoContent)

}

// updateResource ShowAccount godoc
// @Summary update given resource based on resource id
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   resourceId     path    string     true        "Resource Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId} [put]
func (v *IncidentService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//Adding update preprocess request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
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

// getCardView ShowAccount godoc
// @Summary Get all the machine information in a card view
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/card_view [get]
func (v *IncidentService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	if offsetValue != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		listOfObjects, err = GetObjects(dbConnection, targetTable, limitVal)
	} else {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
	}

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

}

func (v *IncidentService) getRecordFormData(ctx *gin.Context) {

	// first get the record
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}

	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)

	ctx.JSON(http.StatusOK, response)

}
