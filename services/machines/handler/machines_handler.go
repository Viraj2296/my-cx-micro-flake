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
	"gorm.io/datatypes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
func (v *MachineService) loadFile(ctx *gin.Context) {
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
func (v *MachineService) importObjects(ctx *gin.Context) {
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
func (v *MachineService) exportObjects(ctx *gin.Context) {
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

	if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
		condition = " where object_info ->> '$.hmiStatus' = 'started' OR object_info ->> '$.hmiStatus' = 'stopped'"
	} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
		condition = " where object_info ->> '$.rejectedQuantity' > 0 "
	}
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
func (v *MachineService) getTableImportSchema(ctx *gin.Context) {
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
func (v *MachineService) getExportSchema(ctx *gin.Context) {
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
func (v *MachineService) getObjects(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	orderValue := ctx.Query("order")
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
		totalRecords = int64(len(*listOfObjects))
	} else {
		if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
			statusCondition := " object_info ->> '$.hmiStatus' = 'started' OR object_info ->> '$.hmiStatus' = 'stopped'"
			totalRecords = CountByCondition(dbConnection, targetTable, statusCondition)
		} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
			statusCondition := " object_info ->> '$.rejectedQuantity' > 0 "
			totalRecords = CountByCondition(dbConnection, targetTable, statusCondition)
		} else {
			totalRecords = Count(dbConnection, targetTable)
		}

		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			queryCondition := component.TableCondition(offsetValue, fields, values, condition)

			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				offsetValue = strconv.Itoa(int(totalRecords) - limitVal + 1)

				limitVal = limitVal - offsetVal
				queryCondition = component.TableCondition(offsetValue, fields, values, condition)
				if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
					queryCondition = queryCondition + " AND (object_info ->> '$.hmiStatus' = 'started' OR object_info ->> '$.hmiStatus' = 'stopped')"

				} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
					queryCondition = queryCondition + " AND (object_info ->> '$.rejectedQuantity' > 0)"

				}
				orderBy := "object_info ->> '$.createdAt' desc"
				listOfObjects, err = GetConditionalObjectsOrderBy(dbConnection, targetTable, queryCondition, orderBy, limitVal)
				currentRecordCount := len(*listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(queryCondition, "AND")
					var totalRecordObjects *[]component.GeneralObject
					if len(andClauses) > 1 {
						totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition)

					} else {
						totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
					}

					if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[0].Id {
						isNext = false
					}
				}

			} else {
				if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
					queryCondition = queryCondition + " AND (object_info ->> '$.rejectedQuantity' > 0)"

				}
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, queryCondition, limitVal)

				currentRecordCount := len(*listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(queryCondition, "AND")
					var totalRecordObjects *[]component.GeneralObject
					if len(andClauses) > 1 {
						totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition)

					} else {
						totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
					}
					lenTotalRecord := len(*totalRecordObjects)
					if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
						isNext = false
					}
				}
			}
			//limitVal, _ := strconv.Atoi(limitValue)
			//listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), limitVal)
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
func (v *MachineService) getGroupBy(ctx *gin.Context) {

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
func (v *MachineService) getCardViewGroupBy(ctx *gin.Context) {

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
// @Router /project/{projectId}/machines/component/{componentName}/group_by_card_view [get]
func (v *MachineService) getGroupByCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := util.GetProjectId(ctx)
	//groupByFields := ctx.Query("groupByFields")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	searchFields := ctx.Query("search")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	listOfGroupByObjects, _ := GetObjects(dbConnection, MachineSubCategoryTable)
	var cacheData = make(map[int]string, 0)
	for _, groupByObject := range *listOfGroupByObjects {
		machineCategory := MachineSubCategory{ObjectInfo: groupByObject.ObjectInfo}
		cacheData[groupByObject.Id] = machineCategory.getSubCategoryInfo().Name
	}

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
	//d

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponseMap := v.ComponentManager.GetCardViewArrayOfMapInterface(listOfObjects, componentName)
	var dd []GroupByCardView
	for _, responseMap := range cardViewResponseMap {
		field6 := util.InterfaceToInt(responseMap["field6"])
		field6val := util.InterfaceToString(cacheData[field6])

		var isElementFound bool
		isElementFound = false
		for index, mm := range dd {
			if mm.GroupByField == field6val {
				dd[index].Cards = append(dd[index].Cards, responseMap)
				isElementFound = true
			}
		}
		if !isElementFound {
			xl := GroupByCardView{}
			xl.GroupByField = field6val
			xl.Cards = append(xl.Cards, responseMap)
			dd = append(dd, xl)
		}
	}

	ctx.JSON(http.StatusOK, dd)

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
func (v *MachineService) getCardView(ctx *gin.Context) {

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

func (v *MachineService) deleteValidation(ctx *gin.Context) {

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
func (v *MachineService) deleteResource(ctx *gin.Context) {

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
func (v *MachineService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	if targetTable == AssemblyMachineHmiTable {
		v.updateHmiInfo(ctx)
		return
	}

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
func (v *MachineService) createNewResource(ctx *gin.Context) {

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

	if componentName == MachineHMIComponent {
		v.handleNewHMIResource(ctx)
		return
	}

	if componentName == AssemblyMachineHmiComponent {
		v.handleNewAssemblyHMIResource(ctx)
		return
	}

	if componentName == ToolingMachineHmiComponent {
		v.handleNewToolingHMIResource(ctx)
		return
	}

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
	if componentName == MachineMasterComponent {
		// if it is the machine master , then create machine hmi setting, and machine display setting
		v.BaseService.Logger.Info("creating hmi setting for new machine", zap.Any("machine_id", createdRecordId))
		var hmiOperators = make([]int, 0)
		var hmiStopReasons = make([]int, 0)

		hmiSettingInfo := HmiSettingInfo{
			CreatedAt:                      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:                      userId,
			HmiOperators:                   hmiOperators,
			ObjectStatus:                   common.ObjectStatusActive,
			LastUpdatedAt:                  util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedBy:                  userId,
			HmiStopReasons:                 hmiStopReasons,
			HmiAutoStopPeriod:              "120m",
			MachineLiveDetectionInterval:   "5m",
			WarningMessageGenerationPeriod: "15m",
		}
		serialisedObject, _ := json.Marshal(hmiSettingInfo)
		commonObject := component.GeneralObject{Id: createdRecordId, ObjectInfo: serialisedObject}
		CreateWithId(dbConnection, MachineHMISettingSettingTable, commonObject)
		displaySettingInfo := DisplaySettingInfo{
			DisplayEnabled:  true,
			DisplayInterval: 5,
			ObjectStatus:    common.ObjectStatusActive,
			CreatedAt:       util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:       userId,
			LastUpdatedAt:   util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedBy:   userId,
		}
		serialisedObject, _ = json.Marshal(displaySettingInfo)
		commonObject = component.GeneralObject{Id: createdRecordId, ObjectInfo: serialisedObject}
		CreateWithId(dbConnection, MachineDisplaySettingTable, commonObject)
		v.BaseService.Logger.Info("creating machine display  setting for new machine", zap.Any("machine_id", createdRecordId))
	} else if componentName == AssemblyMachineMasterComponent {
		// if it is the machine master , then create machine hmi setting, and machine display setting
		v.BaseService.Logger.Info("creating hmi setting for new assembly machine", zap.Any("machine_id", createdRecordId))
		var hmiOperators = make([]int, 0)
		var hmiStopReasons = make([]int, 0)

		hmiSettingInfo := HmiSettingInfo{
			CreatedAt:                      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:                      userId,
			HmiOperators:                   hmiOperators,
			ObjectStatus:                   common.ObjectStatusActive,
			LastUpdatedAt:                  util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedBy:                  userId,
			HmiStopReasons:                 hmiStopReasons,
			HmiAutoStopPeriod:              "120m",
			MachineLiveDetectionInterval:   "5m",
			WarningMessageGenerationPeriod: "15m",
		}
		serialisedObject, _ := json.Marshal(hmiSettingInfo)
		commonObject := component.GeneralObject{Id: createdRecordId, ObjectInfo: serialisedObject}
		CreateWithId(dbConnection, AssemblyMachineHmiSettingTable, commonObject)
		displaySettingInfo := DisplaySettingInfo{
			DisplayEnabled:  true,
			DisplayInterval: 5,
			ObjectStatus:    common.ObjectStatusActive,
			CreatedAt:       util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:       userId,
			LastUpdatedAt:   util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedBy:   userId,
		}
		serialisedObject, _ = json.Marshal(displaySettingInfo)
		commonObject = component.GeneralObject{Id: createdRecordId, ObjectInfo: serialisedObject}
		CreateWithId(dbConnection, AssemblyMachineHmiSettingTable, commonObject)

		// create the view now
		err = v.ViewManager.CreateNewView(createdRecordId)
		if err != nil {
			v.BaseService.Logger.Error("error creating new assembly machine view", zap.Error(err))
		}
		v.BaseService.Logger.Info("creating assembly machine display  setting for new machine", zap.Any("machine_id", createdRecordId))
	} else if componentName == ToolingMachineMasterComponent {
		// if it is the machine master , then create machine hmi setting, and machine display setting
		v.BaseService.Logger.Info("creating hmi setting for new tooling machine", zap.Any("machine_id", createdRecordId))
		var hmiOperators = make([]int, 0)
		var hmiStopReasons = make([]int, 0)

		hmiSettingInfo := HmiSettingInfo{
			CreatedAt:                      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:                      userId,
			HmiOperators:                   hmiOperators,
			ObjectStatus:                   common.ObjectStatusActive,
			LastUpdatedAt:                  util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedBy:                  userId,
			HmiStopReasons:                 hmiStopReasons,
			HmiAutoStopPeriod:              "120m",
			MachineLiveDetectionInterval:   "5m",
			WarningMessageGenerationPeriod: "15m",
		}
		serialisedObject, _ := json.Marshal(hmiSettingInfo)
		commonObject := component.GeneralObject{Id: createdRecordId, ObjectInfo: serialisedObject}
		CreateWithId(dbConnection, ToolingMachineHmiTable, commonObject)
		displaySettingInfo := DisplaySettingInfo{
			DisplayEnabled:  true,
			DisplayInterval: 5,
			ObjectStatus:    common.ObjectStatusActive,
			CreatedAt:       util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:       userId,
			LastUpdatedAt:   util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedBy:   userId,
		}
		serialisedObject, _ = json.Marshal(displaySettingInfo)
		commonObject = component.GeneralObject{Id: createdRecordId, ObjectInfo: serialisedObject}
		CreateWithId(dbConnection, ToolingMachineHmiSettingTable, commonObject)
		v.BaseService.Logger.Info("creating tooling machine display  setting for new machine", zap.Any("machine_id", createdRecordId))
	}
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
func (v *MachineService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

func (v *MachineService) getHMIView(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var err error

	listOfSubcategoryObjects, _ := GetObjects(dbConnection, MachineSubCategoryTable)
	var cacheData = make(map[int]string, 0)
	for _, groupByObject := range *listOfSubcategoryObjects {
		machineCategory := MachineSubCategory{ObjectInfo: groupByObject.ObjectInfo}
		cacheData[groupByObject.Id] = machineCategory.getSubCategoryInfo().Name
	}

	var statusData = make(map[int]string, 0)
	var statusNameData = make(map[int]string, 0)
	listOfConnectStatusObjects, _ := GetObjects(dbConnection, MachineConnectStatusTable)
	for _, statusObject := range *listOfConnectStatusObjects {
		statusObjectInfo := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, &statusObjectInfo)
		statusName := statusObject.Id
		colorCode := util.InterfaceToString(statusObjectInfo["colorCode"])
		status := util.InterfaceToString(statusObjectInfo["status"])
		statusNameData[statusObject.Id] = status
		statusData[statusName] = colorCode
	}
	//statusData["Maintenance"] = "#edf10e"
	//statusData["Live"] = "#32602e"
	//statusData["Waiting For Feed"] = "#E60E18"

	listOfMachineMasterObjects, err := GetObjects(dbConnection, MachineMasterTable)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

	var dd []GroupByCardView
	for _, machineMasterObject := range *listOfMachineMasterObjects {
		machineMaster := MachineMaster{ObjectInfo: machineMasterObject.ObjectInfo}
		var cardResponse = make(map[string]interface{})
		if machineMaster.getMachineMasterInfo().ObjectStatus != "Archived" {
			var isBatchManaged = false
			machineImage := machineMaster.getMachineMasterInfo().MachineImage
			machineName := machineMaster.getMachineMasterInfo().NewMachineId
			machineConnectStatus := machineMaster.getMachineMasterInfo().MachineConnectStatus
			subCategoryId := machineMaster.getMachineMasterInfo().SubCategory
			err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineMasterObject.Id)
			if err != nil {
				cardResponse["productionOrder"] = "-"
				cardResponse["orderStatus"] = "-"
			} else {
				scheduledOrderInfo := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
				cardResponse["productionOrder"] = scheduledOrderInfo.Name
				orderStatusString := productionOrderInterface.OrderStatusId2String(projectId, scheduledOrderInfo.EventStatus)
				cardResponse["orderStatus"] = orderStatusString
				productionOrderId := scheduledOrderInfo.EventSourceId
				err, productionGeneralObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, productionOrderId, machineMasterObject.Id)

				if err == nil {
					var productionOrderInfo = make(map[string]interface{})
					json.Unmarshal(productionGeneralObject.ObjectInfo, &productionOrderInfo)

					partNumber := util.InterfaceToInt(productionOrderInfo["partNumber"])

					if partNumber != 0 {
						_, partGeneralObject := mouldInterface.GetPartInfo(projectId, partNumber)
						partInfo := make(map[string]interface{})
						json.Unmarshal(partGeneralObject.ObjectInfo, &partInfo)

						isBatchManaged = util.InterfaceToBool(partInfo["isBatchManaged"])
					}

				}

			}

			schedulerOrderEventInfo := getMouldingMachineOrderDetails(machineMasterObject.Id)
			if schedulerOrderEventInfo != nil {
				cardResponse["scheduledQty"] = schedulerOrderEventInfo.ScheduledQty
				cardResponse["completedQty"] = schedulerOrderEventInfo.CompletedQty
				cardResponse["progressPercentage"] = schedulerOrderEventInfo.PercentDone
				cardResponse["cycleCount"] = schedulerOrderEventInfo.CompletedQty
			} else {
				cardResponse["scheduledQty"] = 0
				cardResponse["completedQty"] = 0
				cardResponse["progressPercentage"] = 0
				cardResponse["cycleCount"] = 0
			}

			cardResponse["machineImage"] = machineImage
			cardResponse["id"] = machineMasterObject.Id
			cardResponse["colorCode"] = statusData[machineConnectStatus]
			cardResponse["machineName"] = machineName
			cardResponse["isBatchManaged"] = isBatchManaged
			cardResponse["lastUpdatedAt"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", machineMaster.getMachineMasterInfo().LastUpdatedAt)
			cardResponse["machineConnectStatus"] = statusNameData[machineConnectStatus]

			subCategoryValue := util.InterfaceToString(cacheData[subCategoryId])

			var isElementFound bool
			isElementFound = false
			for index, mm := range dd {
				if mm.GroupByField == subCategoryValue {
					dd[index].Cards = append(dd[index].Cards, cardResponse)
					isElementFound = true
				}
			}
			if !isElementFound {
				xl := GroupByCardView{}
				xl.GroupByField = subCategoryValue
				xl.Cards = append(xl.Cards, cardResponse)
				dd = append(dd, xl)
			}
		}

	}

	ctx.JSON(http.StatusOK, dd)
}

// To get assigned machines based on the user Id
func (v *MachineService) getAssignedMachines(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	userId := common.GetUserId(ctx)

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var err error

	listOfSubcategoryObjects, _ := GetObjects(dbConnection, MachineSubCategoryTable)
	var cacheData = make(map[int]string)
	for _, groupByObject := range *listOfSubcategoryObjects {
		machineCategory := MachineSubCategory{ObjectInfo: groupByObject.ObjectInfo}
		cacheData[groupByObject.Id] = machineCategory.getSubCategoryInfo().Name
	}

	var statusData = make(map[int]string)
	var statusNameData = make(map[int]string)
	listOfConnectStatusObjects, _ := GetObjects(dbConnection, MachineConnectStatusTable)
	for _, statusObject := range *listOfConnectStatusObjects {
		statusObjectInfo := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, &statusObjectInfo)
		statusName := statusObject.Id
		colorCode := util.InterfaceToString(statusObjectInfo["colorCode"])
		status := util.InterfaceToString(statusObjectInfo["status"])
		statusNameData[statusObject.Id] = status
		statusData[statusName] = colorCode
	}

	listOfMachineMasterObjects, err := GetObjects(dbConnection, MachineMasterTable)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

	var dd []GroupByCardView
	for _, machineMasterObject := range *listOfMachineMasterObjects {
		machineMaster := MachineMaster{ObjectInfo: machineMasterObject.ObjectInfo}
		var cardResponse = make(map[string]interface{})
		if machineMaster.getMachineMasterInfo().ObjectStatus != "Archived" {
			var isBatchManaged = false
			machineImage := machineMaster.getMachineMasterInfo().MachineImage
			machineName := machineMaster.getMachineMasterInfo().NewMachineId
			machineConnectStatus := machineMaster.getMachineMasterInfo().MachineConnectStatus
			subCategoryId := machineMaster.getMachineMasterInfo().SubCategory

			// Fetch HMI settings and check if user is assigned
			err, hmiSettingObject := Get(dbConnection, MachineHMISettingSettingTable, machineMasterObject.Id)
			if err != nil {
				continue // Skip this machine if there is an error fetching HMI settings
			}
			machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
			listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
			var operatorConfigured bool
			operatorConfigured = false
			for _, operatorId := range listOfOperators {
				if userId == operatorId {
					operatorConfigured = true
					break
				}
			}
			if !operatorConfigured {
				continue // Skip this machine if the user is not an HMI operator
			}

			err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineMasterObject.Id)
			if err != nil {
				cardResponse["productionOrder"] = "-"
				cardResponse["orderStatus"] = "-"
			} else {
				scheduledOrderInfo := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
				cardResponse["productionOrder"] = scheduledOrderInfo.Name
				orderStatusString := productionOrderInterface.OrderStatusId2String(projectId, scheduledOrderInfo.EventStatus)
				cardResponse["orderStatus"] = orderStatusString

				productionOrderId := scheduledOrderInfo.EventSourceId
				err, productionGeneralObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, productionOrderId, machineMasterObject.Id)

				if err == nil {
					var productionOrderInfo = make(map[string]interface{})
					json.Unmarshal(productionGeneralObject.ObjectInfo, &productionOrderInfo)

					partNumber := util.InterfaceToInt(productionOrderInfo["partNumber"])

					if partNumber != 0 {
						_, partGeneralObject := mouldInterface.GetPartInfo(projectId, partNumber)
						partInfo := make(map[string]interface{})
						json.Unmarshal(partGeneralObject.ObjectInfo, &partInfo)

						isBatchManaged = util.InterfaceToBool(partInfo["isBatchManaged"])
					}

				}

			}
			schedulerOrderEventInfo := getMachineOrderDetails(machineMasterObject.Id)
			if schedulerOrderEventInfo != nil {
				cardResponse["scheduledQty"] = schedulerOrderEventInfo.ScheduledQty
				cardResponse["completedQty"] = schedulerOrderEventInfo.CompletedQty
				cardResponse["progressPercentage"] = schedulerOrderEventInfo.PercentDone
				cardResponse["cycleCount"] = schedulerOrderEventInfo.CompletedQty
			} else {
				cardResponse["scheduledQty"] = 0
				cardResponse["completedQty"] = 0
				cardResponse["progressPercentage"] = 0
				cardResponse["cycleCount"] = 0
			}
			cardResponse["machineImage"] = machineImage
			cardResponse["id"] = machineMasterObject.Id
			cardResponse["colorCode"] = statusData[machineConnectStatus]
			cardResponse["machineName"] = machineName
			cardResponse["isBatchManaged"] = isBatchManaged
			cardResponse["lastUpdatedAt"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", machineMaster.getMachineMasterInfo().LastUpdatedAt)
			cardResponse["machineConnectStatus"] = statusNameData[machineConnectStatus]

			subCategoryValue := util.InterfaceToString(cacheData[subCategoryId])

			var isElementFound bool
			isElementFound = false
			for index, mm := range dd {
				if mm.GroupByField == subCategoryValue {
					dd[index].Cards = append(dd[index].Cards, cardResponse)
					isElementFound = true
				}
			}
			if !isElementFound {
				xl := GroupByCardView{}
				xl.GroupByField = subCategoryValue
				xl.Cards = append(xl.Cards, cardResponse)
				dd = append(dd, xl)
			}
		}

	}

	ctx.JSON(http.StatusOK, dd)
}

func (v *MachineService) getAssemblyHMIView(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var err error

	listOfSubcategoryObjects, _ := GetObjects(dbConnection, AssemblyMachineLineTable)
	var cacheData = make(map[int]string, 0)
	for _, groupByObject := range *listOfSubcategoryObjects {
		machineCategory := make(map[string]interface{})
		json.Unmarshal(groupByObject.ObjectInfo, &machineCategory)
		cacheData[groupByObject.Id] = util.InterfaceToString(machineCategory["name"])
	}

	var statusData = make(map[int]string, 0)
	var statusNameData = make(map[int]string, 0)
	listOfConnectStatusObjects, _ := GetObjects(dbConnection, MachineConnectStatusTable)
	for _, statusObject := range *listOfConnectStatusObjects {
		statusObjectInfo := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, &statusObjectInfo)
		statusName := statusObject.Id
		colorCode := util.InterfaceToString(statusObjectInfo["colorCode"])
		status := util.InterfaceToString(statusObjectInfo["status"])
		statusNameData[statusName] = status
		statusData[statusName] = colorCode
	}

	listOfMachineMasterObjects, err := GetObjects(dbConnection, AssemblyMachineMasterTable)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	var dd []GroupByCardView
	for _, machineMasterObject := range *listOfMachineMasterObjects {
		machineMaster := AssemblyMachineMaster{ObjectInfo: machineMasterObject.ObjectInfo}
		var cardResponse = make(map[string]interface{})
		if machineMaster.getAssemblyMachineMasterInfo().ObjectStatus != "Archived" && machineMaster.getAssemblyMachineMasterInfo().IsEnabled {
			machineImage := machineMaster.getAssemblyMachineMasterInfo().MachineImage
			machineName := machineMaster.getAssemblyMachineMasterInfo().NewMachineId
			machineConnectStatus := machineMaster.getAssemblyMachineMasterInfo().MachineConnectStatus
			subCategoryId := machineMaster.getAssemblyMachineMasterInfo().AssemblyLineOption
			err, scheduledEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineMasterObject.Id)
			if err != nil {
				cardResponse["productionOrder"] = "-"
				cardResponse["orderStatus"] = "-"
			} else {
				scheduledOrderInfo := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
				cardResponse["productionOrder"] = scheduledOrderInfo.Name
				orderStatusString := productionOrderInterface.OrderStatusId2String(projectId, scheduledOrderInfo.EventStatus)
				cardResponse["orderStatus"] = orderStatusString
			}

			cardResponse["machineImage"] = machineImage
			cardResponse["id"] = machineMasterObject.Id
			cardResponse["colorCode"] = statusData[machineConnectStatus]
			cardResponse["machineName"] = machineName
			schedulerOrderEventInfo := getMachineOrderDetails(machineMasterObject.Id)
			if schedulerOrderEventInfo != nil {
				cardResponse["scheduledQty"] = schedulerOrderEventInfo.ScheduledQty
				cardResponse["completedQty"] = schedulerOrderEventInfo.CompletedQty
				cardResponse["progressPercentage"] = schedulerOrderEventInfo.PercentDone
				cardResponse["cycleCount"] = schedulerOrderEventInfo.CompletedQty
			} else {
				cardResponse["scheduledQty"] = 0
				cardResponse["completedQty"] = 0
				cardResponse["progressPercentage"] = 0
				cardResponse["cycleCount"] = 0
			}

			cardResponse["lastUpdatedAt"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", machineMaster.getAssemblyMachineMasterInfo().LastUpdatedAt)
			cardResponse["machineConnectStatus"] = statusNameData[machineConnectStatus]

			subCategoryValue := util.InterfaceToString(cacheData[subCategoryId])

			var isElementFound bool
			isElementFound = false
			for index, mm := range dd {
				if mm.GroupByField == subCategoryValue {
					dd[index].Cards = append(dd[index].Cards, cardResponse)
					isElementFound = true
				}
			}
			if !isElementFound {
				xl := GroupByCardView{}
				xl.GroupByField = subCategoryValue
				xl.Cards = append(xl.Cards, cardResponse)
				dd = append(dd, xl)
			}
		}

	}

	ctx.JSON(http.StatusOK, dd)
}

func (v *MachineService) getToolingHMIView(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var err error

	listOfMachineMasterObjects, err := GetObjects(dbConnection, ToolingMachineMasterTable)

	//var statusData = make(map[string]string, 0)
	//statusData["Maintenance"] = "#edf10e"
	//statusData["Live"] = "#32602e"
	//statusData["Waiting For Feed"] = "#E60E18"
	var statusData = make(map[int]string, 0)
	var statusNameData = make(map[int]string, 0)
	listOfConnectStatusObjects, _ := GetObjects(dbConnection, MachineConnectStatusTable)
	for _, statusObject := range *listOfConnectStatusObjects {
		statusObjectInfo := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, &statusObjectInfo)
		statusName := statusObject.Id
		colorCode := util.InterfaceToString(statusObjectInfo["colorCode"])
		status := util.InterfaceToString(statusObjectInfo["status"])
		statusNameData[statusName] = status
		statusData[statusName] = colorCode
	}

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	var dd []interface{}
	for _, machineMasterObject := range *listOfMachineMasterObjects {
		machineMaster := AssemblyMachineMaster{ObjectInfo: machineMasterObject.ObjectInfo}
		var cardResponse = make(map[string]interface{})
		if machineMaster.getAssemblyMachineMasterInfo().ObjectStatus != "Archived" {
			machineImage := machineMaster.getAssemblyMachineMasterInfo().MachineImage
			machineName := machineMaster.getAssemblyMachineMasterInfo().NewMachineId
			machineStatus := machineMaster.getAssemblyMachineMasterInfo().MachineConnectStatus
			//machineConnectStatus := machineMaster.getAssemblyMachineMasterInfo().MachineConnectStatus
			err, scheduledEventObject := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineMasterObject.Id)
			if err != nil {
				cardResponse["productionOrder"] = "-"
				cardResponse["orderStatus"] = "-"
			} else {
				scheduledOrderInfo := make(map[string]interface{})
				json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledOrderInfo)
				fmt.Println(scheduledOrderInfo["name"])
				cardResponse["productionOrder"] = scheduledOrderInfo["name"]
				orderStatusString := productionOrderInterface.OrderStatusId2String(projectId, util.InterfaceToInt(scheduledOrderInfo["eventStatus"]))
				cardResponse["orderStatus"] = orderStatusString
			}
			schedulerOrderEventInfo := getMachineOrderDetails(machineMasterObject.Id)
			if schedulerOrderEventInfo != nil {
				cardResponse["scheduledQty"] = schedulerOrderEventInfo.ScheduledQty
				cardResponse["completedQty"] = schedulerOrderEventInfo.CompletedQty
				cardResponse["progressPercentage"] = schedulerOrderEventInfo.PercentDone
				cardResponse["cycleCount"] = schedulerOrderEventInfo.CompletedQty
			} else {
				cardResponse["scheduledQty"] = 0
				cardResponse["completedQty"] = 0
				cardResponse["progressPercentage"] = 0
				cardResponse["cycleCount"] = 0
			}
			cardResponse["machineImage"] = machineImage
			cardResponse["id"] = machineMasterObject.Id
			cardResponse["machineName"] = machineName
			cardResponse["colorCode"] = statusData[machineStatus]
			cardResponse["lastUpdatedAt"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", machineMaster.getAssemblyMachineMasterInfo().LastUpdatedAt)
			cardResponse["machineConnectStatus"] = statusNameData[machineMaster.getAssemblyMachineMasterInfo().MachineConnectStatus]

			dd = append(dd, cardResponse)

		}

	}

	ctx.JSON(http.StatusOK, dd)
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
func (v *MachineService) getRecordFormData(ctx *gin.Context) {

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

	if targetTable == MachineMasterTable {
		//groupId, _ := strconv.Atoi(ms.GroupPermissionConfig)
		authInterface := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		listOfUserBasicInfo := authInterface.GetUserInfoFromGroupId(v.GroupPermissionConfig)

		canCreateWorkOrder := response["canCreateWorkOrder"].(component.RecordInfo)
		canCreateWorkOrder.Value = false
		response["canCreateWorkOrder"] = canCreateWorkOrder

		for _, userInfo := range listOfUserBasicInfo {

			if userInfo.UserId == userId {
				canCreateWorkOrder = response["canCreateWorkOrder"].(component.RecordInfo)
				canCreateWorkOrder.Value = true
				response["canCreateWorkOrder"] = canCreateWorkOrder
			}
		}
	}

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
func (v *MachineService) getSearchResults(ctx *gin.Context) {

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

// getDashboard ShowAccount godoc
// @Summary Get the record form data to display dashboard
// @Description based on user permission, user will get machine dashboard data
// @Tags Dashboards
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   machineId     path    string     true        "Machine Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/{machineId}/dashboard [get]
func (v *MachineService) getMachineDashboard(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	machineId, _ := strconv.Atoi(recordId)

	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	var summaryResponse datatypes.JSON
	if targetTable == AssemblyMachineStatisticsTable {
		summaryResponse = v.getAssemblyMachineDashboardResponse(projectId, machineId)
	} else if targetTable == ToolingMachineStatisticsTable {
		summaryResponse = v.getToolingMachineDashboardResponse(projectId, machineId)
	} else {
		summaryResponse = v.getMachineDashboardResponse(projectId, machineId)
	}

	ctx.JSON(http.StatusOK, summaryResponse)
}

// getDashboard ShowAccount godoc
// @Summary Get the record form data to display dashboard
// @Description based on user permission, user will get machine dashboard data
// @Tags Dashboards
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   machineId     path    string     true        "Machine Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/{machineId}/dashboard [get]
func (v *MachineService) getIntialMachineDashboardInfo(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	userId := common.GetUserId(ctx)
	summaryResponse := v.getIntialMachineDashboardResponse(projectId, userId)
	ctx.JSON(http.StatusOK, summaryResponse)
}

// getHMIInfo ShowAccount godoc
// @Summary Get the hmi info to display in hmi interface
// @Description based on user permission, user will get machine dashboard data
// @Tags HMI
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   machineId     path    string     true        "Machine Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/{machineId}/hmi_info [get]
func (v *MachineService) getHmiInfo(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	machineId, _ := strconv.Atoi(recordId)

	userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	emptyHMIResponse, _ := json.Marshal(HMIInfoResponse{})
	var listOfResetOperators = make([]int, 0)
	var machineMasterInfo map[string]interface{}
	if targetTable == MachineHMITable {
		err, machineMasterGeneralObject := Get(dbConnection, MachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}
		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, MachineHMISettingSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	} else if targetTable == ToolingMachineHmiTable {
		err, machineMasterGeneralObject := Get(dbConnection, ToolingMachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}
		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, ToolingMachineHmiSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	} else {
		err, machineMasterGeneralObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}

		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, AssemblyMachineHmiSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		listOfResetOperators = machineHMISetting.getHMISettingInfo().ResetOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	}

	hmiResponse := v.getHMIInfoResponse(projectId, machineId, userId, machineMasterInfo, dbConnection, targetTable)
	if targetTable == AssemblyMachineHmiComponent {
		var jsonHmiResponse = make(map[string]interface{})
		json.Unmarshal(hmiResponse, &jsonHmiResponse)
		var operatorConfigured = false
		operatorConfigured = false
		for _, operatorId := range listOfResetOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		jsonHmiResponse["canReset"] = component.GetBoolRecordInfo(operatorConfigured, "bool")
		hmiResponse, _ = json.Marshal(jsonHmiResponse)
	}
	ctx.JSON(http.StatusOK, hmiResponse)
}

// getHMIInfo ShowAccount godoc
// @Summary Get the hmi info to display in hmi interface
// @Description based on user permission, user will get machine dashboard data
// @Tags HMI
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   machineId     path    string     true        "Machine Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/{machineId}/hmi_info [get]
func (v *MachineService) getNextHmi(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	machineId, _ := strconv.Atoi(recordId)
	eventId := util.GetEventId(ctx)
	userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	emptyHMIResponse, _ := json.Marshal(HMIInfoResponse{})

	var machineMasterInfo map[string]interface{}
	if targetTable == MachineHMITable {
		err, machineMasterGeneralObject := Get(dbConnection, MachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}
		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, MachineHMISettingSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	} else if targetTable == ToolingMachineHmiTable {
		err, machineMasterGeneralObject := Get(dbConnection, ToolingMachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}
		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, ToolingMachineHmiSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	} else {
		err, machineMasterGeneralObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}

		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, AssemblyMachineHmiSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	}
	hmiResponse := v.getNextHMIResponse(projectId, eventId, machineId, userId, machineMasterInfo, dbConnection, targetTable)
	if targetTable == AssemblyMachineHmiComponent {
		var jsonHmiResponse = make(map[string]interface{})
		json.Unmarshal(hmiResponse, &jsonHmiResponse)

		jsonHmiResponse["canReset"] = component.GetBoolRecordInfo(false, "bool")
		hmiResponse, _ = json.Marshal(jsonHmiResponse)
	}
	ctx.JSON(http.StatusOK, hmiResponse)
}

// getHMIInfo ShowAccount godoc
// @Summary Get the hmi info to display in hmi interface
// @Description based on user permission, user will get machine dashboard data
// @Tags HMI
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   machineId     path    string     true        "Machine Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/{machineId}/hmi_info [get]
func (v *MachineService) getPreviousHmi(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	machineId, _ := strconv.Atoi(recordId)
	eventId := util.GetEventId(ctx)
	userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	emptyHMIResponse, _ := json.Marshal(HMIInfoResponse{})

	var machineMasterInfo map[string]interface{}
	if targetTable == MachineHMITable {
		err, machineMasterGeneralObject := Get(dbConnection, MachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}
		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, MachineHMISettingSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	} else if targetTable == ToolingMachineHmiTable {
		err, machineMasterGeneralObject := Get(dbConnection, ToolingMachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}
		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, ToolingMachineHmiSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	} else {
		err, machineMasterGeneralObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)

		if err != nil {
			// this is something wrong, not found machine means, something manual operation happened
			ctx.JSON(http.StatusOK, emptyHMIResponse)
			return
		}

		json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineMasterInfo)
		err, hmiSettingObject := Get(dbConnection, AssemblyMachineHmiSettingTable, machineId)
		machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingObject.ObjectInfo}
		listOfOperators := machineHMISetting.getHMISettingInfo().HmiOperators
		var operatorConfigured bool
		operatorConfigured = false
		for _, operatorId := range listOfOperators {
			if userId == operatorId {
				operatorConfigured = true
			}
		}
		if !operatorConfigured {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
			return
		}
	}
	hmiResponse := v.getPreviousHMIResponse(projectId, eventId, machineId, userId, machineMasterInfo, dbConnection, targetTable)
	if targetTable == AssemblyMachineHmiComponent {
		var jsonHmiResponse = make(map[string]interface{})
		json.Unmarshal(hmiResponse, &jsonHmiResponse)

		jsonHmiResponse["canReset"] = component.GetBoolRecordInfo(false, "bool")
		hmiResponse, _ = json.Marshal(jsonHmiResponse)
	}
	ctx.JSON(http.StatusOK, hmiResponse)
}

func (v *MachineService) removeInternalArrayReference(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	//projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ReferenceDatabase

	var removeInternalReferenceRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&removeInternalReferenceRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	fmt.Println("objectInterface:", objectInterface)
	serializedObject := v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject
	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
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

func (v *MachineService) refreshExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
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

//func (v *MachineService) AdvanceSearchResults(ctx *gin.Context) {
//	var searchFieldCommand common.FilterCriteria
//	userId := common.GetUserId(ctx)
//	zone := getUserTimezone(userId)
//	dbConnection := v.BaseService.ReferenceDatabase
//	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
//		ctx.AbortWithError(http.StatusBadRequest, err)
//		return
//	}
//	var listOfObjects *[]component.GeneralObject
//	var err error
//
//	if len(searchFieldCommand.Filters) == 0 {
//		// reset the search
//		queryParam := ctx.Request.URL.Query()
//		queryParam.Set("offset", "1")
//		queryParam.Set("limit", "30")
//		ctx.Request.URL.RawQuery = queryParam.Encode()
//
//		v.defaultAdvanceSearch(ctx)
//		return
//	}
//
//	searchList := v.ComponentManager.ExecuteFilter(dbConnection, searchFieldCommand.ComponentName, searchFieldCommand)
//	fmt.Println("searchList: ", searchList)
//	if searchList != "()" {
//		searchQuery := " id in " + searchList
//		listOfObjects, err = GetConditionalObjects(dbConnection, searchFieldCommand.ComponentName, searchQuery)
//		fmt.Println("searchQuery: ", searchQuery)
//		fmt.Println("err: ", err)
//		if err != nil {
//			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
//			return
//		}
//	}
//
//	var lenObject = 0
//	if listOfObjects != nil {
//		lenObject = len(*listOfObjects)
//	}
//	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(lenObject), searchFieldCommand.ComponentName, "", zone)
//	ctx.JSON(http.StatusOK, searchResponse)
//}

func (v *MachineService) defaultAdvanceSearch(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	orderValue := ctx.Query("order")
	condition := ctx.Query("condition")
	outFields := ctx.Query("out_fields")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	var err error

	err, cacheFilter := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
		return
	}
	var cacheFilterInfo = common.FilterInfo{}
	json.Unmarshal(cacheFilter.ObjectInfo, &cacheFilterInfo)
	targetTable = cacheFilterInfo.ComponentName

	//Have to next flag
	isNext := true
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64

	if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
		statusCondition := " object_info ->> '$.hmiStatus' = 'started' OR object_info ->> '$.hmiStatus' = 'stopped'"
		totalRecords = CountByCondition(dbConnection, targetTable, statusCondition)
	} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
		statusCondition := " object_info ->> '$.rejectedQuantity' > 0 "
		totalRecords = CountByCondition(dbConnection, targetTable, statusCondition)
	} else {
		totalRecords = Count(dbConnection, targetTable)
	}

	if limitValue == "" {
		listOfObjects, _ = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

	} else {
		limitVal, _ := strconv.Atoi(limitValue)
		queryCondition := component.TableCondition(offsetValue, fields, values, condition)

		if orderValue == "desc" {
			offsetVal, _ := strconv.Atoi(offsetValue)
			offsetValue = strconv.Itoa(int(totalRecords) - limitVal + 1)

			limitVal = limitVal - offsetVal
			queryCondition = component.TableCondition(offsetValue, fields, values, condition)
			if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
				queryCondition = queryCondition + " AND (object_info ->> '$.hmiStatus' = 'started' OR object_info ->> '$.hmiStatus' = 'stopped')"

			} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
				queryCondition = queryCondition + " AND (object_info ->> '$.rejectedQuantity' > 0)"

			}
			orderBy := "object_info ->> '$.createdAt' desc"
			listOfObjects, _ = GetConditionalObjectsOrderBy(dbConnection, targetTable, queryCondition, orderBy, limitVal)
			currentRecordCount := len(*listOfObjects)

			if currentRecordCount < limitVal {
				isNext = false
			} else if currentRecordCount == limitVal {
				andClauses := strings.Split(queryCondition, "AND")
				var totalRecordObjects *[]component.GeneralObject
				if len(andClauses) > 1 {
					totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition)

				} else {
					totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
				}

				if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[0].Id {
					isNext = false
				}
			}

		} else {
			if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
				queryCondition = queryCondition + " AND (object_info ->> '$.rejectedQuantity' > 0)"

			}
			listOfObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition, limitVal)

			currentRecordCount := len(*listOfObjects)

			if currentRecordCount < limitVal {
				isNext = false
			} else if currentRecordCount == limitVal {
				andClauses := strings.Split(queryCondition, "AND")
				var totalRecordObjects *[]component.GeneralObject
				if len(andClauses) > 1 {
					totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition)

				} else {
					totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
				}
				lenTotalRecord := len(*totalRecordObjects)
				if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
					isNext = false
				}
			}
		}
		//limitVal, _ := strconv.Atoi(limitValue)
		//listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), limitVal)
	}

	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)

	tableObjectResponse := component.TableObjectResponse{}
	json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

	tableObjectResponse.IsNext = isNext
	tableRecordsResponse, _ = json.Marshal(tableObjectResponse)
	ctx.JSON(http.StatusOK, tableRecordsResponse)

}

func (v *MachineService) getCacheFilters(ctx *gin.Context) {

	// first get the record
	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	//targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	var conditionalString = "object_info ->> '$.createdBy' = " + strconv.Itoa(userId) + " and object_info ->> '$.componentName' = '" + componentName + "' and object_info ->> '$.objectStatus' = 'Active'"
	fmt.Println("conditionalString: ", conditionalString)
	generalObject, err := GetConditionalObjects(dbConnection, MachineFilterTable, conditionalString)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}
	tableObjectResponse := component.TableObjectResponse{}
	var dataResponse = make([]datatypes.JSON, 0)
	for _, cacheResponse := range *generalObject {
		var responseObject = make(map[string]interface{})
		json.Unmarshal(cacheResponse.ObjectInfo, &responseObject)
		responseObject["id"] = cacheResponse.Id
		jsonResponse, _ := json.Marshal(responseObject)
		dataResponse = append(dataResponse, jsonResponse)
	}
	tableObjectResponse.Data = dataResponse
	fmt.Println("tableObjectResponse: ", tableObjectResponse)
	//responseObjectData, _ := json.Marshal(tableObjectResponse)

	ctx.JSON(http.StatusOK, tableObjectResponse)

}
