package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

// loadFile ShowAccount godoc
// @Summary load the file and get the schema information with data(currently only csv format)
// @Description based on user permission, user will allow importing csv file url to populate machine register
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param User body  component.LoadDataFileCommand true "Send the following fields"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/loadFile [post]
func (v *TraceabilityService) loadFile(ctx *gin.Context) {
	loadDataFileCommand := component.LoadDataFileCommand{}
	if err := ctx.ShouldBindBodyWith(&loadDataFileCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err, errorCode, loadFileResponse := v.ComponentManager.ProcessLoadFile(loadDataFileCommand)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	ctx.JSON(http.StatusOK, loadFileResponse)
	return
}

// importObjects ShowAccount godoc
// @Summary import machine register information (currently only csv format)
// @Description based on user permission, user will allow importing csv file url to populate machine register
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param User body  component.ImportDataCommand true "Send the following fields"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/import [get]
func (v *TraceabilityService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	projectId := ctx.Param("projectId")

	componentName := ctx.Param("componentName")
	//targetTable := ms.ComponentManager.GetTargetTable(componentName)
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, errorCode, _ := v.ComponentManager.ImportData(dbConnection, componentName, importDataCommand)
	if err != nil {
		v.BaseService.Logger.Error("unable to import data", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	//var failedRecords int
	//var recordId int
	//for _, object := range listOfObjects {
	//	err, recordId = Create(dbConnection, targetTable, object)
	//
	//	if err != nil {
	//		ms.BaseService.Logger.Error("unable to create record", zap.String("error", err.Error()))
	//		failedRecords = failedRecords + 1
	//	}
	//	recordIdInString := strconv.Itoa(recordId)
	//	CreateBotRecordTrail(projectId, recordIdInString, componentName, "machine master is created")
	//}
	//importDataResponse := component.ImportDataResponse{
	//	TotalRecords:  totalRecords,
	//	FailedRecords: failedRecords,
	//	Message:       "data is successfully imported",
	//}

	//ctx.JSON(http.StatusOK, importDataResponse)
}

// exportObjects ShowAccount godoc
// @Summary export machine related information (currently only csv format)
// @Description based on user permission, user will allow importing csv file url to populate machine register
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param User body  component.ExportDataCommand true "Send the following fields"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/import [get]
func (v *TraceabilityService) exportObjects(ctx *gin.Context) {
	// we will get the uploaded url
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	exportCommand := component.ExportDataCommand{}

	if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var condition string
	err, errorCode, exportDataResponse := v.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
	if err != nil {
		v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	ctx.JSON(http.StatusOK, exportDataResponse)
}

// getTableSchema ShowAccount godoc
// @Summary Get the table schema
// @Description based on user permission, user will get the table related fields
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/table_import_schema [get]
func (v *TraceabilityService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := v.ComponentManager.GetTableImportSchema(componentName)
	ctx.JSON(http.StatusOK, tableImportSchema)
}

// getExportSchema ShowAccount godoc
// @Summary Get the table schema
// @Description based on user permission, user will get the table related fields
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/getExportSchema [get]
func (v *TraceabilityService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)
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
func (v *TraceabilityService) getObjects(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
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
		listOfObjects, err = GetObjects(dbConnection, targetTable)

		if listOfObjects != nil {
			totalRecords = int64(len(*listOfObjects))
		}

	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), limitVal)
		}

		var currentRecordCount int
		if listOfObjects != nil {
			currentRecordCount = len(*listOfObjects)
		}

		conditionString := component.TableCondition(offsetValue, fields, values, condition)
		limitVal, _ := strconv.Atoi(limitValue)
		if currentRecordCount < limitVal {
			isNext = false
		} else if currentRecordCount == limitVal {
			andClauses := strings.Split(conditionString, "AND")
			var totalRecordObjects *[]component.GeneralObject
			if len(andClauses) > 1 {
				totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, conditionString)

			} else {
				totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
			}
			lenTotalRecord := len(*totalRecordObjects)
			if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
				isNext = false
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
func (v *TraceabilityService) getGroupBy(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")

	groupByAction := component.GroupByAction{}
	if err := ctx.ShouldBindBodyWith(&groupByAction, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// if the group by is empty, then return the normal default 30 records per page results.
	if len(groupByAction.GroupBy) == 0 {
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getObjects(ctx)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetObjects(dbConnection, targetTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	fmt.Println("groupByAction : ", groupByAction)
	var totalRecords = Count(dbConnection, targetTable)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
	tableResponse := component.TableObjectResponse{}
	json.Unmarshal(tableRecordsResponse, &tableResponse)
	groupByColumns := groupByAction.GroupBy
	finalResponse := component.TableObjectResponse{}
	if len(groupByColumns) == 1 {
		results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range results {
			groupByChildren := component.GroupByChildren{}

			groupByChildren.Data = level1Value
			groupByChildren.Type = "json"

			tableGroupResponse := component.TableGroupByResponse{}
			tableGroupResponse.Label = level1Key
			tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Data = append(finalResponse.Data, rawData)
		}
	} else if len(groupByColumns) == 2 {
		level1Results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range level1Results {
			tableGroupResponse := component.TableGroupByResponse{}
			if len(level1Value) > 1 {
				var internalGroupResponse []interface{}
				// here we need to group again

				level2Results := component.GetGroupByResultsFromInterface(groupByColumns[1], level1Value)
				for level2Key, level2Value := range level2Results {
					level2Children := component.GroupByChildren{}
					level2Children.Data = level2Value
					level2Children.Type = "json"

					internalTableGroupResponse := component.TableGroupByResponse{}
					internalTableGroupResponse.Label = level2Key
					internalTableGroupResponse.Children = append(internalTableGroupResponse.Children, level2Children)
					internalGroupResponse = append(internalGroupResponse, internalTableGroupResponse)
				}
				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = internalGroupResponse
			} else {
				groupByChildren := component.GroupByChildren{}
				groupByChildren.Data = level1Value
				groupByChildren.Type = "json"

				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			}

			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Data = append(finalResponse.Data, rawData)
		}
	}

	finalResponse.Header = tableResponse.Header
	finalResponse.TotalRowCount = tableResponse.TotalRowCount
	finalResponse.CurrentRowCount = tableResponse.CurrentRowCount

	ctx.JSON(http.StatusOK, finalResponse)
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
func (v *TraceabilityService) getCardViewGroupBy(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")

	groupByAction := component.GroupByAction{}
	if err := ctx.ShouldBindBodyWith(&groupByAction, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// if the group by is empty, then return the normal default 30 records per page results.
	if len(groupByAction.GroupBy) == 0 {
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getCardView(ctx)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetObjects(dbConnection, targetTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	fmt.Println("groupByAction : ", groupByAction)
	var totalRecords = Count(dbConnection, targetTable)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
	tableResponse := component.TableObjectResponse{}
	json.Unmarshal(tableRecordsResponse, &tableResponse)
	groupByColumns := groupByAction.GroupBy
	finalResponse := component.CardViewGroupResponse{}
	if len(groupByColumns) == 1 {
		results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range results {
			groupByChildren := component.GroupByChildren{}

			groupByChildren.Data = v.ComponentManager.GetCardViewFromListOfInterface(level1Value, componentName)
			groupByChildren.Type = "json"

			tableGroupResponse := component.TableGroupByResponse{}
			tableGroupResponse.Label = level1Key
			tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Cards = append(finalResponse.Cards, rawData)
		}
	} else if len(groupByColumns) == 2 {
		level1Results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range level1Results {
			tableGroupResponse := component.TableGroupByResponse{}
			if len(level1Value) > 1 {
				var internalGroupResponse []interface{}
				// here we need to group again

				level2Results := component.GetGroupByResultsFromInterface(groupByColumns[1], level1Value)
				for level2Key, level2Value := range level2Results {
					level2Children := component.GroupByChildren{}
					level2Children.Data = v.ComponentManager.GetCardViewFromListOfInterface(level2Value, componentName)
					level2Children.Type = "json"

					internalTableGroupResponse := component.TableGroupByResponse{}
					internalTableGroupResponse.Label = level2Key
					internalTableGroupResponse.Children = append(internalTableGroupResponse.Children, level2Children)
					internalGroupResponse = append(internalGroupResponse, internalTableGroupResponse)
				}
				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = internalGroupResponse
			} else {
				groupByChildren := component.GroupByChildren{}
				groupByChildren.Data = level1Value
				groupByChildren.Type = "json"

				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			}

			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Cards = append(finalResponse.Cards, rawData)
		}
	}

	ctx.JSON(http.StatusOK, finalResponse)
}

type GroupByCardView struct {
	GroupByField string                   `json:"groupByField"`
	Cards        []map[string]interface{} `json:"cards"`
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
func (v *TraceabilityService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := util.GetProjectId(ctx)
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	searchFields := ctx.Query("search")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	if searchFields != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchQuery, limitVal)
	} else {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
	}

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}

	cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

}

func (v *TraceabilityService) deleteValidation(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, recordId)

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Resource",
				Description: "The resource that you are trying to delete doesn't exist, Please check refresh page and try again",
			})
		return
	}
	if component.IsArchived(generalObject.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Resource Archived",
				Description: "The resource that you are trying to delete is already archived. This operation is not allowed",
			})
		return
	}
	var dependencyComponents []string
	var dependencyRecords int
	v.checkReference(dbConnection, componentName, recordId, &dependencyComponents, &dependencyRecords)
	if dependencyRecords > 0 {
		var dependencyString string
		dependencyComponents = util.RemoveDuplicateString(dependencyComponents)
		dependencyString = " ["
		for index, dependencyComponent := range dependencyComponents {
			if index == len(dependencyComponents)-1 {
				dependencyString += dependencyComponent
			} else {
				dependencyString += dependencyComponent + " ->"
			}
		}
		dependencyString += " ]"
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			CanDelete: false,
			Code:      100,
			Message:   "There are dependencies bound to the resource that you are trying to remove. Removing this resource would create the chain removal on following resources " + dependencyString + " in " + strconv.Itoa(dependencyRecords) + " places, Please understand the risk of deleting this resource as all the dependencies would be achieved immediately, and this process is not reversible",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		CanDelete: true,
		Code:      100,
		Message:   "There are no dependencies bound to the resource that you are trying to remove. So, removing this resource won't affect others resource now, you can proceed !!",
	})

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
func (v *TraceabilityService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, recordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		sendResourceNotFound(ctx)
		return
	}
	// first archive all the dependency
	userId := common.GetUserId(ctx)
	v.archiveReferences(userId, dbConnection, componentName, recordId)

	err = ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		sendArchiveFailed(ctx)
		return
	}

	v.CreateUserRecordMessage(ProjectID, componentName, "Resource is archived, no further modification allowed", recordId, userId, nil, nil)
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
func (v *TraceabilityService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	err, objectInterface := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	if !common.ValidateObjectStatus(objectInterface.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      getError(common.InvalidObjectStatusError).Error(),
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})

	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Updating Resource Failed",
				Description: "Error updating resource information due to internal system error. Please report this error to system administrator",
			})
		return
	}
	userId := common.GetUserId(ctx)
	v.CreateUserRecordMessage(ProjectID, componentName, "Resource got updated", intRecordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

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
func (v *TraceabilityService) createNewResource(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// here we should do the validation
	err := v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	var createdRecordId int

	processedDefaultValues := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	serializedRequest, _ := json.Marshal(processedDefaultValues)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(serializedRequest, ctx)
	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}

	err, createdRecordId = Create(dbConnection, targetTable, object)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
		return
	}

	userId := common.GetUserId(ctx)

	v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	component.SendObjectCreationResponse(ctx, projectId, componentName, createdRecordId)
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
func (v *TraceabilityService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

// getRecordFormData ShowAccount godoc
// @Summary Get the record form data to facilitate the update
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId} [get]
func (v *TraceabilityService) getRecordFormData(ctx *gin.Context) {

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

// getSearchResults ShowAccount godoc
// @Summary Get the search results based on given input
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param SearchField body SearchKeys true "Pass the array of key and values"
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/search [post]
func (v *TraceabilityService) getSearchResults(ctx *gin.Context) {

	var searchFieldCommand []component.SearchKeys
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	format := ctx.Query("format")
	if format != "" {
		if format == "card_view" {
			if len(searchFieldCommand) == 0 {
				// reset the search
				ctx.Set("offset", 1)
				ctx.Set("limit", 30)
				v.getObjects(ctx)
				return
			}

			searchList := v.ComponentManager.GetSearchQueryV2(dbConnection, componentName, searchFieldCommand)
			searchQuery := " id in " + searchList
			if searchList == "()" {
				cardViewResponse := component.CardViewResponse{}
				ctx.JSON(http.StatusOK, cardViewResponse)
				return
			}
			listOfObjects, err := GetConditionalObjects(dbConnection, targetTable, searchQuery)
			// fmt.Println("listOfObjects:", len(*listOfObjects))
			if err != nil {
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
				return
			}
			cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), ErrorGettingObjectsInformation)
			return

		}
	}
	if len(searchFieldCommand) == 0 {
		// reset the search
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getObjects(ctx)
		return
	}

	searchList := v.ComponentManager.GetSearchQueryV2(dbConnection, componentName, searchFieldCommand)
	// if searchList == "()" {
	// 	tableObject := component.TableObjectResponse{}
	// 	tableObjectResponse, _ := json.Marshal(tableObject)
	// 	ctx.JSON(http.StatusOK, tableObjectResponse)
	// 	return
	// }
	listOfObject := make([]component.GeneralObject, 0)
	listOfObjects := &listOfObject
	var err error
	if searchList != "()" {
		searchQuery := " id in " + searchList
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchQuery)
		// fmt.Println("listOfObjects:", len(*listOfObjects))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
			return
		}
	}

	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)
}
