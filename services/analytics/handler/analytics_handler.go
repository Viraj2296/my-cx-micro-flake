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
	"github.com/google/uuid"
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
func (as *AnalyticsService) loadFile(ctx *gin.Context) {
	loadDataFileCommand := component.LoadDataFileCommand{}
	if err := ctx.ShouldBindBodyWith(&loadDataFileCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err, errorCode, loadFileResponse := as.ComponentManager.ProcessLoadFile(loadDataFileCommand)
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
func (as *AnalyticsService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	projectId := ctx.Param("projectId")

	componentName := ctx.Param("componentName")
	//targetTable := as.ComponentManager.GetTargetTable(componentName)
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	err, errorCode, _ := as.ComponentManager.ImportData(dbConnection, componentName, importDataCommand)
	if err != nil {
		as.BaseService.Logger.Error("unable to import data", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	//var failedRecords int
	//var recordId int
	//for _, object := range listOfObjects {
	//	err, recordId = Create(dbConnection, targetTable, object)
	//
	//	if err != nil {
	//		as.BaseService.Logger.Error("unable to create record", "error", err.Error())
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
func (as *AnalyticsService) exportObjects(ctx *gin.Context) {
	// we will get the uploaded url
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	exportCommand := component.ExportDataCommand{}

	if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var condition string
	err, errorCode, exportDataResponse := as.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
	if err != nil {
		as.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
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
func (as *AnalyticsService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := as.ComponentManager.GetTableImportSchema(componentName)
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
func (as *AnalyticsService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := as.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)
}

// getDataBaseSchema ShowAccount godoc
// @Summary Get the table schema
// @Description based on user permission, user will get the table related fields
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/database_schema [get]
func (as *AnalyticsService) getDataBaseSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	recordId := util.GetRecordIdString(ctx)
	recordIdInt, _ := strconv.Atoi(recordId)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	targetTable := as.ComponentManager.GetTargetTable(componentName)

	err, generalObject := Get(dbConnection, targetTable, recordIdInt)
	fmt.Println("targetTable: ", targetTable)
	fmt.Println("recordIdInt: ", recordIdInt)
	if err != nil {
		fmt.Println("err: ", err)
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	datasourceInfo := DatasourceInfo{}
	connectionParam := make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &datasourceInfo)

	byteData, _ := json.Marshal(datasourceInfo.ConnectionParam)

	json.Unmarshal(byteData, &connectionParam)

	databaseNodeSchema := as.ComponentManager.GetDatabaseTableSchema(dbConnection, util.InterfaceToString(connectionParam["schema"]), datasourceInfo.ConnectedDatabaseTables)
	ctx.JSON(http.StatusOK, databaseNodeSchema)

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
func (as *AnalyticsService) getObjects(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := as.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	condition := ctx.Query("condition")
	outFields := ctx.Query("out_fields")
	format := ctx.Query("format")
	searchFields := ctx.Query("search")
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
	//Have to next flag
	isNext := true
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
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
		searchQuery := as.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		searchWithBaseQuery := searchQuery + " AND " + baseCondition
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), limitVal)
		}

		currentRecordCount := len(*listOfObjects)
		limitVal, _ := strconv.Atoi(limitValue)
		if currentRecordCount < limitVal {
			isNext = false
		} else if currentRecordCount == limitVal {

			totalRecordObjects, _ := GetObjects(dbConnection, targetTable)

			lenTotalRecord := len(*totalRecordObjects)
			if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
				isNext = false
			}
		}

	}
	as.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := as.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		_, tableRecordsResponse := as.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)

		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
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
func (as *AnalyticsService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := as.ComponentManager.GetTargetTable(componentName)
	projectId := util.GetProjectId(ctx)
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	searchFields := ctx.Query("search")

	dbConnection := as.BaseService.ServiceDatabases[projectId]
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
		searchQuery := as.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchQuery, limitVal)
	} else {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
	}

	as.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := as.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

}

func (as *AnalyticsService) deleteValidation(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordIdString(ctx)
	recordIdInt, _ := strconv.Atoi(recordId)
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	err, _ := Get(dbConnection, targetTable, recordIdInt)

	if err != nil {
		ctx.JSON(http.StatusCreated, response.GeneralResponse{
			Code:    100,
			Message: "Error looking dependencies, check your resource id. This is mainly due to internal system error or some process deleted your resource while you looking at",
		})
		return
	}

	// check constraints and proceed
	listOfConstraints := as.ComponentManager.GetConstraints(componentName)
	for _, constraint := range listOfConstraints {
		referenceComponent := constraint.Reference
		referenceField := constraint.ReferenceProperty
		referenceTable := as.ComponentManager.GetTargetTable(referenceComponent)

		numberOfRecords := NoOfReferenceObjects(dbConnection, referenceTable, referenceField, recordIdInt)
		if numberOfRecords == -1 {
			ctx.JSON(http.StatusCreated, response.GeneralResponse{
				Code:    100,
				Message: "There are dependencies bound to the resource that you are trying to remove. But, system couldn't able to verify the number of resources component, proceed with your own risk",
			})

		} else {
			ctx.JSON(http.StatusCreated, response.GeneralResponse{
				Code:    100,
				Message: "There are dependencies bound to the resource that you are trying to remove. This resources is used under " + constraint.ReferenceComponentDisplayName + " in " + strconv.Itoa(numberOfRecords) + " places, Please understand the risk of deleting as all the dependencies would be arhvied immediately, and this process is not reversible",
			})
		}

	}
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		CanDelete: true,
		Code:      100,
		Message:   "There are no dependencies bound to the resource that you are trying to remove. You can proceed",
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
func (as *AnalyticsService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordIdString(ctx)
	recordIdInt, _ := strconv.Atoi(recordId)
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, recordIdInt)

	if err != nil {
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	// check constraints and proceed
	listOfConstraints := as.ComponentManager.GetConstraints(componentName)
	for _, constraint := range listOfConstraints {
		referenceComponent := constraint.Reference
		referenceField := constraint.ReferenceProperty
		referenceTable := as.ComponentManager.GetTargetTable(referenceComponent)

		ArchiveReferenceObjects(dbConnection, referenceTable, referenceField, recordIdInt)
	}
	as.ComponentManager.ProcessDeleteDependencyInjection(dbConnection, recordIdInt, componentName)
	updatedObjectInfo := common.UpdateMetaInfoFromSerializedObject(generalObject.ObjectInfo, ctx)
	err = ArchiveObject(dbConnection, targetTable, component.GeneralObject{Id: generalObject.Id, ObjectInfo: updatedObjectInfo})

	if err != nil {
		as.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
		return
	}

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
func (as *AnalyticsService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	userId := common.GetUserId(ctx)
	userIdStr := strconv.Itoa(userId)
	dbConnection := as.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	updatingData := make(map[string]interface{})

	err, objectInterface := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		as.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
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

	if componentName == AnalyticsDataSourceComponent {
		// Get the permissions for the dashboard
		condition := " object_info ->> '$.assignedUserId' = " + userIdStr + " AND object_info ->> '$.dashboardId'  = " + recordIdString
		listOfDashboardPermission, _ := GetConditionalObjects(dbConnection, AnalyticsDashboardPermissionTable, condition)

		if len(*listOfDashboardPermission) == 0 {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("user doesn't have permission to add widget into dashborad"), ErrorCreatingObjectInformation)
			return
		}

		for _, dashboardPermission := range *listOfDashboardPermission {
			dashboardPermissionInfo := make(map[string]interface{})
			json.Unmarshal(dashboardPermission.ObjectInfo, &dashboardPermissionInfo)
			permissionLeval := util.InterfaceToInt(dashboardPermissionInfo["permissionLevel"])

			if permissionLeval != WritePermission {
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("user doesn't have permission to add widget into dashborad"), ErrorCreatingObjectInformation)
				return
			}
		}

	}

	updatingData = make(map[string]interface{})

	serializedObject := as.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject

	err = as.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	err = Update(as.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}
	//rawInfo, _ := json.Marshal(updateRequest)
	//recordMessage := componentName + " is updated"
	//err = CreateUserRecordTrail(as, projectId, recordIdString, componentName, recordMessage, &component.GeneralObject{Id: intRecordId, ObjectInfo: objectInterface.ObjectInfo}, &component.GeneralObject{Id: intRecordId, ObjectInfo: rawInfo})
	//if err != nil {
	//	as.BaseService.Logger.Error("error in create record trail", zap.String("error", err.Error()))
	//}

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
func (as *AnalyticsService) createNewResource(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	fmt.Println("targetTable : ", targetTable)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// here we should do the validation
	err := as.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}
	fmt.Println("createRequest : ", createRequest)
	//Added Preprocessed request
	processedDefaultValues := as.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	serializedRequest, _ := json.Marshal(processedDefaultValues)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(serializedRequest, ctx)

	fmt.Println("preprocessedRequest : ", preprocessedRequest)
	if componentName == AnalyticsDataSourceComponent {
		var preProcessedRequestFields = make(map[string]interface{})
		json.Unmarshal(preprocessedRequest, &preProcessedRequestFields)
		preProcessedRequestFields["status"] = "Connected"

		// get the datasourceMaster
		datasourceMasterId := util.InterfaceToInt(preProcessedRequestFields["datasourceMaster"])
		fmt.Println("datasourceMasterId : ", datasourceMasterId)
		_, datasourceMasterObject := Get(dbConnection, AnalyticsDatasourcesMasterTable, datasourceMasterId)
		datasourceMaster := AnalyticsDatasourcesMaster{ObjectInfo: datasourceMasterObject.ObjectInfo}
		fmt.Println("datasourceMaster.getDatasourceMasterInfo():", datasourceMaster.getDatasourceMasterInfo())
		if datasourceMaster.getDatasourceMasterInfo().Type == "csv" {
			// create the table
			err = as.handleCSV2DatabaseRecordCreation(ctx, dbConnection, preProcessedRequestFields)
			if err != nil {
				return
			}
			preProcessedRequestFields["csvFileConnectionParameters"] = preProcessedRequestFields["connectionParam"]
			preprocessedRequest, _ = json.Marshal(preProcessedRequestFields)
		} else if datasourceMaster.getDatasourceMasterInfo().Type == "mysql" {
			// create the table
			preProcessedRequestFields["mysqlConnectionParam"] = preProcessedRequestFields["connectionParam"]
			preprocessedRequest, _ = json.Marshal(preProcessedRequestFields)
		}

	} else if componentName == AnalyticsWidgetComponent {
		var preProcessedRequestFields = make(map[string]interface{})
		json.Unmarshal(preprocessedRequest, &preProcessedRequestFields)
		var modifiedFilters []interface{}
		if listOfFilters, ok := preProcessedRequestFields["filters"]; ok {
			fmt.Println("listOfFilters", listOfFilters)
			for _, filter := range listOfFilters.([]interface{}) {
				var filterFields = filter.(map[string]interface{})
				filterFields["filterId"] = uuid.New().String()
				modifiedFilters = append(modifiedFilters, filterFields)
			}
		}

		preProcessedRequestFields["filters"] = modifiedFilters
		preprocessedRequest, _ = json.Marshal(preProcessedRequestFields)

	}
	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}

	err, _ = Create(dbConnection, targetTable, object)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Internal Error"), ErrorCreatingObjectInformation, err.Error())
		return
	}
	response.SendObjectCreationMessage(ctx)
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
func (as *AnalyticsService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := as.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

func (as *AnalyticsService) getHMIView(ctx *gin.Context) {

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
func (as *AnalyticsService) getRecordFormData(ctx *gin.Context) {

	// first get the record
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}
	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := as.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)
	var objectFields = make(map[string]interface{})
	json.Unmarshal(rawJSONObject, &objectFields)
	as.getAdditionalRecordFunctionResponse(componentName, objectFields, response)
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
func (as *AnalyticsService) getSearchResults(ctx *gin.Context) {

	var searchFieldCommand []component.SearchKeys
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	if len(searchFieldCommand) == 0 {
		// reset the search
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		as.getObjects(ctx)
		return
	}

	format := ctx.Query("format")
	searchQuery := as.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
	listOfObjects, err := GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := as.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), ErrorGettingObjectsInformation)
			return

		}
	}

	_, searchResponse := as.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)
}

func (as *AnalyticsService) queryResponse(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var queryResponseRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&queryResponseRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err, queryResponse := as.ComponentManager.GetQueryResponse(dbConnection, queryResponseRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Query Execution Failed"), QueryExecutionFailed, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, queryResponse)
}

func (as *AnalyticsService) removeInternalArrayReference(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	//projectId := ctx.Param("projectId")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := as.BaseService.ReferenceDatabase

	var removeInternalReferenceRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&removeInternalReferenceRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		as.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	fmt.Println("objectInterface:", objectInterface)
	serializedObject := as.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject
	err = Update(as.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}
	var updatingObjectFields map[string]interface{}
	json.Unmarshal(serializedObject, &updatingObjectFields)
	ctx.JSON(http.StatusOK, updatingObjectFields)

}

func matchingSelectedIds(selectedIds []string, schemaId float32) bool {

	for _, id := range selectedIds {
		findDataType := strings.Split(id, ".")

		if len(findDataType) == 1 {
			numericId, err := strconv.Atoi(findDataType[0])

			if err != nil {
				return false
			}

			if float32(numericId) == schemaId {
				return true
			}
		} else {
			floatId, _ := strconv.ParseFloat(id, 2)
			if float32(floatId) == schemaId {
				return true
			}
		}

	}
	return false
}

func (as *AnalyticsService) refreshExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := as.ComponentManager.GetTableExportSchema(componentName)
	var appendExportSchema []component.ExportSchema

	exportRefreshPayload := RefreshExportRequest{}
	if err := ctx.ShouldBindBodyWith(&exportRefreshPayload, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, schema := range exportSchema {
		if len(schema.Children) > 0 {
			var appendChildExportSchema []component.Childern
			for _, childSchema := range schema.Children {
				checkMatch := matchingSelectedIds(exportRefreshPayload.SelectedId, childSchema.Id)
				if checkMatch {
					continue
				}
				appendChildExportSchema = append(appendChildExportSchema, childSchema)
			}
			schema.Children = appendChildExportSchema
			appendExportSchema = append(appendExportSchema, schema)
		} else {
			checkMatch := matchingSelectedIds(exportRefreshPayload.SelectedId, schema.Id)
			if checkMatch {
				continue
			}
			appendExportSchema = append(appendExportSchema, schema)
		}

	}

	ctx.JSON(http.StatusOK, appendExportSchema)
}
