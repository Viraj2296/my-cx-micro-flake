package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/datatypes"

	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomExportResult struct {
	Id               int            // id
	ObjectInfo       datatypes.JSON // object_info
	ListOfAttendance string         //list_of_attendance
}

func (v *LabourManagementService) loadFile(ctx *gin.Context) {
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

func (v *LabourManagementService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")
	//
	//componentName := ctx.Param("componentName")

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")

}

func (v *LabourManagementService) exportObjectsWithQueryResults(ctx *gin.Context) {
	// Get the uploaded URL parameters
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	exportCommand := component.ExportDataCommand{}

	// Bind the JSON body to the exportCommand struct
	if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get the database connection for the project
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var condition string

	if componentName == const_util.LabourManagementAttendanceTable {
		err, errorCode, exportDataResponse := v.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
		if err != nil {
			v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
			return
		}

		ctx.JSON(http.StatusOK, exportDataResponse)
		return
	}

	// Initialize conditions slice
	var conditions []string

	// Handle startDate condition
	startDateAttr, startDateExists := exportCommand.Attributes["startDate"]
	if startDateExists && startDateAttr.Type == "date_time" {
		if value, exists := startDateAttr.Value.(string); exists && len(value) > 0 {
			// Add condition for startDate on createdAt field inside object_info
			conditions = append(conditions, fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(lsp.object_info, '$.createdAt')) >= '%s'", value))
		}
	}

	// Handle endDate condition
	endDateAttr, endDateExists := exportCommand.Attributes["endDate"]
	if endDateExists && endDateAttr.Type == "date_time" {
		if value, exists := endDateAttr.Value.(string); exists && len(value) > 0 {
			// Add condition for endDate on createdAt field inside object_info
			conditions = append(conditions, fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(lsp.object_info, '$.createdAt')) <= '%s'", value))
		}
	}

	// If both startDate and endDate are not provided, retrieve all records
	startDateHasValue := startDateExists && exportCommand.Attributes["startDate"].Value != nil
	endDateHasValue := endDateExists && exportCommand.Attributes["endDate"].Value != nil

	if !startDateHasValue && !endDateHasValue {
		// No values for startDate and endDate, retrieve all records
		condition = "" // No conditions mean no WHERE clause
	} else {
		// Join the conditions with AND to create the full query condition
		if len(conditions) > 0 {
			condition = strings.Join(conditions, " AND ")
		}
	}

	// Log the condition for debugging
	v.BaseService.Logger.Info("passing the condition to get the objects", zap.String("condition", condition))

	baseWhere := `
WHERE 
    (lma.object_info->>'$.manufacturingLines' IS NULL 
    OR JSON_CONTAINS(lma.object_info->>'$.manufacturingLines', CAST(amm.object_info->>'$.assemblyLineOption' AS JSON)))
`

	if len(condition) > 0 {
		baseWhere += " AND " + condition
	}
	selectQuery := fmt.Sprintf(`
	SELECT 
		lsp.id AS shift_production_id,
		lsp.object_info AS object_info,
		CONCAT('"', GROUP_CONCAT(DISTINCT u.object_info->>'$.fullName' ORDER BY u.object_info->>'$.fullName' SEPARATOR ', '), '"') AS list_of_attendance
	FROM 
		labour_management_shift_production lsp
	LEFT JOIN 
		labour_management_shift_production_attendance lmsa ON lmsa.shift_production_id = JSON_UNQUOTE(lsp.object_info->>'$.shiftId')
	LEFT JOIN 
		labour_management_attendance lma ON lma.id = lmsa.shift_attendance_id
	LEFT JOIN 
		cx_micro_flake.user u ON u.id = JSON_UNQUOTE(lma.object_info->>'$.userResourceId')
	LEFT JOIN 
		assembly_machine_master amm ON JSON_UNQUOTE(lsp.object_info->>'$.machineId') = amm.id
	%s
	GROUP BY 
		lsp.id,
		amm.object_info->>'$.assemblyLineOption',
		lma.object_info->>'$.manufacturingLines'
	ORDER BY 
		lsp.id
`, baseWhere)
	// SQL query with date filtering and GROUP_CONCAT for listOfAttendance
	// 	selectQuery := `
	// 	SELECT
	//     lsp.id AS shift_production_id,  -- Grouping by shift_production_id
	//     lsp.object_info AS object_info,  -- Include object_info column
	//     CONCAT('"', GROUP_CONCAT(DISTINCT u.object_info->>'$.fullName' ORDER BY u.object_info->>'$.fullName' SEPARATOR ', '), '"') AS list_of_attendance  -- Concatenate unique full names without quotes around each name
	// FROM
	//     labour_management_shift_production lsp
	// LEFT JOIN
	//     labour_management_shift_production_attendance lmsa ON lmsa.shift_production_id = JSON_UNQUOTE(lsp.object_info->>'$.shiftId')
	// LEFT JOIN
	//     labour_management_attendance lma ON lma.id = lmsa.shift_attendance_id
	// LEFT JOIN
	//     cx_micro_flake.user u ON u.id = JSON_UNQUOTE(lma.object_info->>'$.userResourceId')
	// LEFT JOIN
	//     assembly_machine_master amm ON JSON_UNQUOTE(lsp.object_info->>'$.machineId') = amm.id
	// WHERE
	//     lma.object_info->>'$.manufacturingLines' IS NULL
	//     OR JSON_CONTAINS(lma.object_info->>'$.manufacturingLines', CAST(amm.object_info->>'$.assemblyLineOption' AS JSON))  -- Filter records by matching manufacturingLines and assemblyLineOption
	// GROUP BY
	//     lsp.id,  -- Group by shift_production_id
	//     amm.object_info->>'$.assemblyLineOption',  -- Group by assemblyLineOption
	//     lma.object_info->>'$.manufacturingLines'  -- Group by manufacturingLines
	// ORDER BY
	//     lsp.id`

	// Apply condition for filtering by date
	// if len(condition) > 0 {
	// 	selectQuery += " WHERE " + condition
	// }

	var queryResults []CustomExportResult
	dbConnection.Raw(selectQuery).Scan(&queryResults)

	var modifiedQueryResults []component.GeneralObject

	// Process the query results
	for _, queryResult := range queryResults {
		// Add the "listOfAttendance" field to the object_info JSON
		var modifiedObjectInfo = common.AddFieldJSONObject(queryResult.ObjectInfo, "listOfAttendance", queryResult.ListOfAttendance)
		modifiedQueryResults = append(modifiedQueryResults, component.GeneralObject{
			Id:         queryResult.Id,
			ObjectInfo: modifiedObjectInfo,
		})
	}

	// Call the ExportData function to get the filtered data
	err, errorCode, exportDataResponse := v.ComponentManager.ExportDataFromQueryResults(dbConnection, componentName, exportCommand, modifiedQueryResults)
	if err != nil {
		v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}

	// Return the filtered data as a JSON response
	ctx.JSON(http.StatusOK, exportDataResponse)
}

func (v *LabourManagementService) exportObjects(ctx *gin.Context) {
	// Get the uploaded URL parameters
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	exportCommand := component.ExportDataCommand{}

	// Bind the JSON body to the exportCommand struct
	if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get the database connection for the project
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var condition string

	// Initialize conditions slice
	var conditions []string

	// Handle startDate condition
	startDateAttr, startDateExists := exportCommand.Attributes["startDate"]
	if startDateExists && startDateAttr.Type == "date_time" {
		if value, exists := startDateAttr.Value.(string); exists && len(value) > 0 {
			// Add condition for startDate on createdAt field inside object_info
			conditions = append(conditions, fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(object_info, '$.createdAt')) >= '%s'", value))
		}
	}

	// Handle endDate condition
	endDateAttr, endDateExists := exportCommand.Attributes["endDate"]
	if endDateExists && endDateAttr.Type == "date_time" {
		if value, exists := endDateAttr.Value.(string); exists && len(value) > 0 {
			// Add condition for endDate on createdAt field inside object_info
			conditions = append(conditions, fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(object_info, '$.createdAt')) <= '%s'", value))
		}
	}

	// If both startDate and endDate are present without values, retrieve all records
	startDateHasValue := startDateExists && exportCommand.Attributes["startDate"].Value != nil
	endDateHasValue := endDateExists && exportCommand.Attributes["endDate"].Value != nil

	if startDateExists && !startDateHasValue && endDateExists && !endDateHasValue {
		// No values for startDate and endDate, retrieve all records
		condition = "" // No conditions means no WHERE clause
	} else {
		// Join the conditions with AND to create the full query condition
		if len(conditions) > 0 {
			condition = strings.Join(conditions, " AND ")
		}
	}

	// Log the condition for debugging
	v.BaseService.Logger.Info("passing the condition to get the objects", zap.String("condition", condition))

	// Call the ExportData function to get the filtered data
	err, errorCode, exportDataResponse := v.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
	if err != nil {
		v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}

	// Return the filtered data as a JSON response
	ctx.JSON(http.StatusOK, exportDataResponse)
}

func (v *LabourManagementService) getTableImportSchema(ctx *gin.Context) {
	//componentName := ctx.Param("componentName")

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")

}

func (v *LabourManagementService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)
}

func (v *LabourManagementService) getObjects(ctx *gin.Context) {
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
	orderValue := ctx.Query("order")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//Have to next flag
	isNext := true
	var listOfObjects []component.GeneralObject
	var totalRecords int64
	var err error
	userId := common.GetUserId(ctx)
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
		err, listOfObjects = database.GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		err, listOfObjects = database.GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(listOfObjects))
	} else {

		totalRecords = database.Count(dbConnection, targetTable)
		if limitValue == "" {
			err, listOfObjects = database.GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			var conditionString string
			limitVal, _ := strconv.Atoi(limitValue)

			totalRecords = database.CountByCondition(dbConnection, targetTable, conditionString)
			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				var tableCondition string
				if conditionString != "" {
					if offsetVal == -1 {
						tableCondition = component.TableConditionV1(offsetValue, fields, values, condition)
					} else {
						tableCondition = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
					}
					if tableCondition != "" {
						conditionString = tableCondition + " AND " + conditionString
					}
				} else {
					if offsetVal == -1 {
						conditionString = component.TableConditionV1(offsetValue, fields, values, condition)
					} else {
						conditionString = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
					}

				}

				orderBy := "object_info ->> '$.createdAt' desc"

				err, listOfObjects = database.GetConditionalObjectsOrderBy(dbConnection, targetTable, conditionString, orderBy, limitVal)
				currentRecordCount := len(listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(conditionString, "AND")
					var totalRecordObjects []component.GeneralObject
					if len(andClauses) > 1 {
						err, totalRecordObjects = database.GetConditionalObjects(dbConnection, targetTable, conditionString)

					} else {
						err, totalRecordObjects = database.GetObjects(dbConnection, targetTable)
					}

					if (listOfObjects)[currentRecordCount-1].Id == (totalRecordObjects)[0].Id {
						isNext = false
					}
				}
				//listOfObjects = reverseSlice(listOfObjects)
			} else {
				err, listOfObjects = database.GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)
				currentRecordCount := len(listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(conditionString, "AND")
					var totalRecordObjects []component.GeneralObject
					if len(andClauses) > 1 {
						err, totalRecordObjects = database.GetConditionalObjects(dbConnection, targetTable, conditionString)

					} else {
						err, totalRecordObjects = database.GetObjects(dbConnection, targetTable)
					}
					lenTotalRecord := len(totalRecordObjects)
					if (listOfObjects)[currentRecordCount-1].Id == (totalRecordObjects)[lenTotalRecord-1].Id {
						isNext = false
					}
				}
			}
		}

	}
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArrayV1(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		userId = common.GetUserId(ctx)
		zone := getUserTimezone(userId)
		_, tableRecordsResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (v *LabourManagementService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects []component.GeneralObject
	var err error
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if offsetValue != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		err, listOfObjects = database.GetObjects(dbConnection, targetTable, limitVal)
	} else {
		err, listOfObjects = database.GetObjects(dbConnection, targetTable)
	}

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponseV1(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)
}

func (v *LabourManagementService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, recordId)
	userId := common.GetUserId(ctx)
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}

	err = database.ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), const_util.ErrorRemovingObjectInformation)
		return
	}
	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource is archived, no further modification allowed", recordId, userId, nil, nil)
	ctx.Status(http.StatusNoContent)

}

func (v *LabourManagementService) removeInternalArrayReference(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var userId = common.GetUserId(ctx)
	var removeInternalReferenceRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&removeInternalReferenceRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}
	initializedObject := v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = initializedObject
	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource get updated", intRecordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	ctx.JSON(http.StatusOK, updatingData)

}

func (v *LabourManagementService) updateResource(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}

	if !common.ValidateObjectStatus(objectInterface.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      const_util.GetError(common.InvalidObjectStatusError).Error(),
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}

	var updateRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})

	//Adding update process request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("error updating resource", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating resource information"), const_util.ErrorUpdatingObjectInformation)
		return
	}

	if componentName == const_util.LabourManagementSettingComponent {
		v.BaseService.Logger.Info("updating labour management setting...", zap.Any("setting", string(serializedObject)))
		v.WorkflowActionHandler.LabourManagementSettingInfo = database.GetLabourManagementSettingInfo(serializedObject)
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the resource",
	})

}

func (v *LabourManagementService) createNewResource(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, createRequest := component.GetRequestFields(ctx)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// here we should do the validation
	err = v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Validation Failed"), const_util.FieldValidationFailed, err.Error())
		return
	}
	processedDefaultValues := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	processedDefaultObjectMappingFields := v.ComponentManager.ProcessDefaultObjectMapping(dbConnection, processedDefaultValues, componentName)

	processedDefaultObjectMappingFields["objectStatus"] = common.ObjectStatusActive

	if componentName == const_util.LabourManagementShiftMasterTable {
		processedDefaultObjectMappingFields["actionRemarks"] = make([]interface{}, 0)
	}
	var userId = common.GetUserId(ctx)

	rawCreateRequest, _ := json.Marshal(processedDefaultObjectMappingFields)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)

	var createdRecordId int

	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}

	err, createdRecordId = database.CreateFromGeneralObject(dbConnection, targetTable, object)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error in Resource Creation "), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}

	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:     0,
		RecordId: createdRecordId,
		Message:  "New resource is successfully created",
	})
}

func (v *LabourManagementService) getNewRecord(ctx *gin.Context) {
	componentName := util.GetComponentName(ctx)
	projectId := util.GetProjectId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

func (v *LabourManagementService) getRecordFormData(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting table information", zap.String("error", err.Error()), zap.String("table", targetTable))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}

	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)

	ctx.JSON(http.StatusOK, response)

}

func (v *LabourManagementService) getSearchResults(ctx *gin.Context) {
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

	if len(searchFieldCommand) == 0 {
		// reset the search
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getObjects(ctx)
		return
	}

	format := ctx.Query("format")
	searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
	err, listOfObjects := database.GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), const_util.ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := v.ComponentManager.GetCardViewResponseV1(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), const_util.ErrorGettingObjectsInformation)
			return

		}
	}

	_, searchResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, int64(len(listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)

}

func (v *LabourManagementService) deleteValidation(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, recordId)

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

func (v *LabourManagementService) getGroupBy(ctx *gin.Context) {

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
	err, listOfObjects := database.GetObjects(dbConnection, targetTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	var totalRecords = database.Count(dbConnection, targetTable)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
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

func (v *LabourManagementService) getCardViewGroupBy(ctx *gin.Context) {

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
	err, listOfObjects := database.GetObjects(dbConnection, targetTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	var totalRecords = database.Count(dbConnection, targetTable)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
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

func (v *LabourManagementService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			err, listOfObjects := database.GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(listOfObjects)
					for _, referenceObject := range listOfObjects {
						v.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (v *LabourManagementService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			err, listOfObjects := database.GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(listOfObjects) > 0 {
					for _, referenceObject := range listOfObjects {
						err := database.ArchiveObject(dbConnection, referenceTable, referenceObject)
						if err == nil {
							err := v.CreateUserRecordMessage(const_util.ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
							if err == nil {
								v.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
							}
						}

					}
				}
			}

		}
	}
}
func (v *LabourManagementService) refreshExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
	var appendExportSchema []component.ExportSchema

	exportRefreshPayload := database.RefreshExportRequest{}
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
