package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
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
func (v *FacilityService) loadFile(ctx *gin.Context) {
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
func (v *FacilityService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")

	//componentName := ctx.Param("componentName")
	//targetTable := ms.ComponentManager.GetTargetTable(componentName)
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

}

func (v *FacilityService) exportObjects(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	categoryIdStr := ctx.Query("categoryId")
	componentName := util.GetComponentName(ctx)
	exportCommand := component.ExportDataCommand{}

	// Check if categoryId is provided
	if categoryIdStr == "" {
		int64ComponentId := v.ComponentManager.ComponentNameIdMapping[componentName]
		componentSchema := v.ComponentManager.ComponentSchema[int64ComponentId]
		targetTable := componentSchema.TargetTable
		dbConnection := v.BaseService.ServiceDatabases[projectId]

		if componentName == const_util.FacilityServiceMyExecutionRequestComponent {
			userId := common.GetUserId(ctx)
			selectedRequestIds := v.getExecutionRequests(dbConnection, userId)

			// Retrieve the selected columns from the payload
			var payload struct {
				Format string                   `json:"format"`
				Data   []component.ExportSchema `json:"data"`
			}

			if err := ctx.BindJSON(&payload); err != nil {
				v.BaseService.Logger.Error("invalid payload", zap.Error(err))
				response.SendSimpleError(ctx, http.StatusBadRequest, err, const_util.InvalidPayload)
				return
			}

			if len(payload.Data) == 0 {
				ctx.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
				return
			}

			if len(selectedRequestIds) > 0 {
				// Formulate condition string correctly
				var condition string
				condition = "id IN ("
				for i, requestId := range selectedRequestIds {
					if i > 0 {
						condition += ", "
					}
					condition += strconv.Itoa(requestId)
				}
				condition += ")"

				if condition == "" {
					condition = " (object_info ->> '$.assignedExecutionParty' = " + strconv.Itoa(userId) + " or object_info ->> '$.assignedExecutionParty' = 0)"
				} else {
					condition = condition + " AND (object_info ->> '$.assignedExecutionParty' = " + strconv.Itoa(userId) + " or object_info ->> '$.assignedExecutionParty' = 0)"
				}

				queryResults, err := database.GetConditionalObjects(dbConnection, targetTable, condition)
				if err != nil {
					v.BaseService.Logger.Error("Failed to get conditional objects", zap.Error(err))
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
					return
				}

				// Prepare ExportDataCommand
				exportDataCommand := component.ExportDataCommand{
					Format: payload.Format, // Use the format from the payload
					Data:   payload.Data,
				}
				err, errorCode, exportDataResponse := v.ComponentManager.ExportDataV3(dbConnection, componentName, exportDataCommand, *queryResults)
				if err != nil {
					v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to handle export", "errorCode": errorCode})
					return
				}
				// Return the export data response
				ctx.JSON(http.StatusOK, exportDataResponse)
			} else {
				v.BaseService.Logger.Warn("No execution requests found for the user")
				ctx.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
				return
			}
		} else if componentName == const_util.FacilityServiceRequestTable {

			if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}

			var condition string
			err, errorCode, exportDataResponse := v.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
			if err != nil {
				v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
				response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
				return
			}
			ctx.JSON(http.StatusOK, exportDataResponse)
			return
		} else {
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

	} else {
		// Convert categoryId to int
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			v.BaseService.Logger.Error("invalid category id", zap.Int("category_id", categoryId), zap.Error(err))
			response.SendSimpleError(ctx, http.StatusBadRequest, err, const_util.ErrorConvertingField)
			return
		}

		// Retrieve database connection
		dbConnection := v.BaseService.ServiceDatabases[projectId]

		// Retrieve the selected columns from the payload
		var payload struct {
			Format string                   `json:"format"`
			Data   []component.ExportSchema `json:"data"`
		}
		if err := ctx.BindJSON(&payload); err != nil {
			v.BaseService.Logger.Error("invalid payload", zap.Error(err))
			response.SendSimpleError(ctx, http.StatusBadRequest, err, const_util.InvalidPayload)
			return
		}

		if len(payload.Data) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
			return
		}

		tableName := v.ComponentManager.GetTargetTable(componentName)
		conditionString := "object_info->>'$.categoryId'=" + categoryIdStr

		// Get all the records where requested category ID
		listOfObjects, err := database.GetConditionalObjects(dbConnection, tableName, conditionString)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error Getting Resources"), const_util.ErrorCreatingObjectInformation, err.Error())
			return
		}

		// Cache templates
		templateCache := make(map[int][]interface{})
		listOfTemplates, _ := database.GetObjects(dbConnection, const_util.FacilityServiceCategoryTemplateTable)
		for _, templateInterface := range *listOfTemplates {
			var objectFields map[string]interface{}
			json.Unmarshal(templateInterface.ObjectInfo, &objectFields)
			templateRecords := objectFields["templateFields"].([]interface{})
			templateCache[templateInterface.Id] = templateRecords
		}

		// Modify object fields based on templates
		var modifiedObjectFields []component.GeneralObject
		for _, objectInterface := range *listOfObjects {
			var objectFields map[string]interface{}
			json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
			templateFieldId := util.InterfaceToInt(objectFields["templateFields"])
			templateRecords := templateCache[templateFieldId]

			for key, value := range objectFields {
				for _, templateField := range templateRecords {
					var templateRecord common.TemplateRecords
					serializedData, _ := json.Marshal(templateField)
					json.Unmarshal(serializedData, &templateRecord)

					if key == templateRecord.Property {
						for index, val := range templateRecord.DynamicDroppedDownAttributes.ManualFieldsSource {
							if index == util.InterfaceToInt(value) {
								objectFields[key] = val
							}
						}
					}
				}
			}

			raw, _ := json.Marshal(objectFields)
			componentObject := component.GeneralObject{Id: objectInterface.Id, ObjectInfo: raw}
			modifiedObjectFields = append(modifiedObjectFields, componentObject)
		}

		// Prepare ExportDataCommand
		exportDataCommand := component.ExportDataCommand{
			Format: payload.Format, // Use the format from the payload
			Data:   payload.Data,
		}

		err, errorCode, exportDataResponse := v.ComponentManager.ExportDataV2(dbConnection, componentName, exportDataCommand, modifiedObjectFields)
		if err != nil {
			v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to handle export", "errorCode": errorCode})
			return
		}

		// Return the export data response
		ctx.JSON(http.StatusOK, exportDataResponse)
	}
}

func (v *FacilityService) exportObjectsV1(ctx *gin.Context) {
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
func (v *FacilityService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := v.ComponentManager.GetTableImportSchema(componentName)
	ctx.JSON(http.StatusOK, tableImportSchema)
}

// func (v *FacilityService) getTableImportSchema(ctx *gin.Context) {
//     componentName := ctx.Param("componentName")
//     categoryID := ctx.Query("categoryID") // Assuming the category ID is provided as a query parameter

//     // Modify this line to get table import schema based on categoryID
//     tableImportSchema := v.ComponentManager.GetTableImportSchemaByCategory(componentName, categoryID)

//     ctx.JSON(http.StatusOK, tableImportSchema)
// }

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

// type CategoryTemplate struct {
// 	ID       int    `gorm:"column:id"`
// 	Property string `gorm:"column:property"`
// 	Label    string `gorm:"column:label"`
// }

func (v *FacilityService) getExportSchema(ctx *gin.Context) {
	// Get the categoryId from the query parameter
	categoryIdStr := ctx.Query("categoryId")
	if categoryIdStr == "" {
		componentName := ctx.Param("componentName")
		exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
		ctx.JSON(http.StatusOK, exportSchema)
		return
	}

	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid categoryId"})
		return
	}

	// Get the projectId from the context
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	// Fetch data from the it_service_category_template table using the provided dbConnection
	var objectInfo string
	query := dbConnection.Table(const_util.FacilityServiceCategoryTemplateTable).Select("object_info").Where("id = ?", categoryId).Row().Scan(&objectInfo)
	if query != nil {
		log.Println("Error fetching data:", query)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	// Decode the JSON stored in object_info
	var templateInfo map[string]interface{}
	err = json.Unmarshal([]byte(objectInfo), &templateInfo)
	if err != nil {
		log.Println("Error decoding object_info:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode object_info"})
		return
	}

	// Extract templateFields from the decoded JSON
	templateFields, ok := templateInfo["templateFields"].([]interface{})
	if !ok {
		log.Println("Error extracting templateFields from object_info")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract templateFields"})
		return
	}

	// Get the componentName from the context parameters
	componentName := ctx.Param("componentName")

	// Retrieve the table export schema
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)

	// Add index to export schema starting from 1
	exportSchema = addIndexExportSchema(exportSchema, 1)

	// Convert exportSchema to map[string]interface{}
	var exportSchemaMap []map[string]interface{}
	for _, schema := range exportSchema {
		schemaMap := map[string]interface{}{
			"label":         schema.Label,
			"expandedIcon":  schema.ExpandedIcon,
			"collapsedIcon": schema.CollapsedIcon,
			"droppable":     schema.Droppable,
			"data":          schema.Data,
			"id":            schema.Id,
			"targetTable":   schema.TargetTable,
			"dataType":      schema.DataType,
			"linkedMapFlag": schema.LinkedMapFlag,
		}
		if len(schema.Children) > 0 {
			var children []map[string]interface{}
			for _, child := range schema.Children {
				childMap := map[string]interface{}{
					"label": child.Label,
					"icon":  "pi pi-minus-circle", // Use icon for child
					"data":  child.Data,
					"id":    child.Id,
				}
				children = append(children, childMap)
			}
			schemaMap["children"] = children
		}
		exportSchemaMap = append(exportSchemaMap, schemaMap)
	}

	// Construct the specificResponses list with indexes starting after the last exportSchema index
	startId := float32(len(exportSchemaMap) + 1)
	var specificResponses []map[string]interface{}
	for i, field := range templateFields {
		fieldMap, ok := field.(map[string]interface{})
		if !ok {
			log.Println("Error parsing template field")
			continue
		}
		label, _ := fieldMap["label"].(string)
		property, _ := fieldMap["property"].(string)

		specificResponse := map[string]interface{}{
			"id":            startId + float32(i),
			"data":          property,
			"label":         label,
			"droppable":     false,
			"expandedIcon":  "pi pi-folder-open",
			"collapsedIcon": "pi pi-tag",
		}
		specificResponses = append(specificResponses, specificResponse)
	}

	// Merge exportSchemaMap and specificResponses
	finalResponse := append(exportSchemaMap, specificResponses...)

	// Respond with the merged list
	ctx.JSON(http.StatusOK, finalResponse)
}
func addIndexExportSchema(exportSchema []component.ExportSchema, startId float32) []component.ExportSchema {
	mainId := startId
	for index, schema := range exportSchema {
		schema.Id = mainId
		if len(schema.Children) > 0 {
			subId := 1
			for childIndex, schemaChildren := range schema.Children {
				stringIndex := fmt.Sprintf("%.0f", mainId) + "." + strconv.Itoa(subId)
				childrenIndex, _ := strconv.ParseFloat(stringIndex, 32)
				schemaChildren.Id = float32(childrenIndex)
				schema.Children[childIndex] = schemaChildren

				subId += 1
			}
		}
		exportSchema[index] = schema
		mainId += 1
	}
	return exportSchema
}

func (v *FacilityService) getExportSchemaV1(ctx *gin.Context) {
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
func (v *FacilityService) getObjects(ctx *gin.Context) {

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
	userId := common.GetUserId(ctx)

	//Have to next flag
	isNext := true

	userBasedQuery := " object_info ->>'$.createdBy' = " + strconv.Itoa(userId) + " "
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

		if componentName == const_util.FacilityServiceMyRequestComponent {
			searchWithBaseQuery = searchWithBaseQuery + " AND " + userBasedQuery
		} else if componentName == const_util.FacilityServiceMyDepartmentRequestComponent {
			authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
			userInfo := authService.GetUserInfoById(userId)

			userBasedInQuery := " object_info ->>'$.hodEmail' = '" + userInfo.Email + "'"
			searchWithBaseQuery = searchWithBaseQuery + " AND " + userBasedInQuery
		}
		listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		if componentName == const_util.FacilityServiceMyRequestComponent {
			listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, userBasedQuery)
			totalRecords = int64(len(*listOfObjects))
		} else {
			listOfObjects, err = database.GetObjects(dbConnection, targetTable)
			totalRecords = int64(len(*listOfObjects))
		}

	} else {
		totalRecords = database.Count(dbConnection, targetTable)
		if limitValue == "" {
			if componentName == const_util.FacilityServiceMyRequestComponent {
				listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition)+" AND "+userBasedQuery)
			} else {
				listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))
			}

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			var conditionString string
			if componentName == const_util.FacilityServiceMyRequestComponent {
				conditionString = userBasedQuery
			} else if componentName == const_util.FacilityServiceMyDepartmentRequestComponent {
				authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
				userInfo := authService.GetUserInfoById(userId)

				userBasedInQuery := " object_info ->>'$.hodEmail' = '" + userInfo.Email + "'"
				conditionString = userBasedInQuery
			} else if componentName == const_util.FacilityServiceMyTechnicianRequestComponent {

				selectedRequestIds := v.getExecutionRequests(dbConnection, userId)

				var customQuery string
				var idInQuery = ""

				for index, jobRequestId := range selectedRequestIds {
					if index == len(selectedRequestIds)-1 {
						idInQuery += strconv.Itoa(jobRequestId)
					} else {
						idInQuery += strconv.Itoa(jobRequestId) + ","
					}

				}

				if len(selectedRequestIds) != 0 {
					customQuery = " id IN (" + idInQuery + ")"
					if conditionString == "" {
						conditionString = customQuery
					} else {
						conditionString = conditionString + " AND " + customQuery
					}

				}

				if conditionString == "" {
					conditionString = " (object_info ->> '$.assignedExecutionParty' = " + strconv.Itoa(userId) + " or object_info ->> '$.assignedExecutionParty' = 0)"
				} else {
					conditionString = conditionString + " AND (object_info ->> '$.assignedExecutionParty' = " + strconv.Itoa(userId) + " or object_info ->> '$.assignedExecutionParty' = 0)"
				}

			} else if componentName == const_util.FacilityServiceMyEHSManagerRequestComponent {
				//conditionString = " object_info->>'$.serviceStatus' = 4"
				selectedRequestIds := v.getEHSManagerRequests(dbConnection, userId)

				var customQuery string
				var idInQuery = ""

				for index, jobRequestId := range selectedRequestIds {
					if index == len(selectedRequestIds)-1 {
						idInQuery += strconv.Itoa(jobRequestId)
					} else {
						idInQuery += strconv.Itoa(jobRequestId) + ","
					}

				}

				if len(selectedRequestIds) != 0 {
					customQuery = " id IN (" + idInQuery + ")"
					conditionString = customQuery
				}

				if conditionString == "" {
					conditionString = " JSON_LENGTH(object_info->>'$.isSafetyRelated') > 0"
				} else {
					conditionString = conditionString + " AND JSON_LENGTH(object_info->>'$.isSafetyRelated') > 0"
				}

				// facility all safety service
			} else if componentName == const_util.FacilityServiceAllSafetyRequestComponent {
				conditionString = " JSON_LENGTH(object_info->>'$.isSafetyRelated') > 0"
			} else if componentName == const_util.FacilityServiceMyFacilityManagerRequestComponent {
				//conditionString = " object_info->>'$.serviceStatus' = 5"
				selectedRequestIds := v.getFacilityManagerRequests(dbConnection, userId)

				var customQuery string
				var idInQuery = ""

				for index, jobRequestId := range selectedRequestIds {
					if index == len(selectedRequestIds)-1 {
						idInQuery += strconv.Itoa(jobRequestId)
					} else {
						idInQuery += strconv.Itoa(jobRequestId) + ","
					}

				}

				if len(selectedRequestIds) != 0 {
					customQuery = " id IN (" + idInQuery + ")"
					conditionString = customQuery
				}
			}

			v.BaseService.Logger.Info("executing condition string", zap.String("condition_string", conditionString))
			totalRecords = database.CountByCondition(dbConnection, targetTable, conditionString)
			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				//if offsetVal == -1 {
				//	offsetValue = strconv.Itoa(int(totalRecords) - limitVal + 1)
				//}
				//limitVal = limitVal - offsetVal
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

				listOfObjects, err = database.GetConditionalObjectsOrderBy(dbConnection, targetTable, conditionString, orderBy, limitVal)

				currentRecordCount := len(*listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(conditionString, "AND")
					var totalRecordObjects *[]component.GeneralObject
					if len(andClauses) > 1 {
						totalRecordObjects, _ = database.GetConditionalObjects(dbConnection, targetTable, conditionString)

					} else {
						totalRecordObjects, _ = database.GetObjects(dbConnection, targetTable)
					}

					if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[0].Id {
						isNext = false
					}
				}

			} else {
				listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)
				if err == nil {
					currentRecordCount := len(*listOfObjects)
					if currentRecordCount < limitVal {
						isNext = false
					} else if currentRecordCount == limitVal {
						andClauses := strings.Split(conditionString, "AND")
						var totalRecordObjects *[]component.GeneralObject
						if len(andClauses) > 1 {
							totalRecordObjects, _ = database.GetConditionalObjects(dbConnection, targetTable, conditionString)

						} else {
							totalRecordObjects, _ = database.GetObjects(dbConnection, targetTable)
						}
						lenTotalRecord := len(*totalRecordObjects)
						if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
							isNext = false
						}
					}
				} else {
					v.BaseService.Logger.Error("error getting resource records", zap.String("error", err.Error()))
				}

			}

		}

	}

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		//userId := common.GetUserId(ctx)
		zone := getUserTimezone(userId)
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)

		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}

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
func (v *FacilityService) getCardView(ctx *gin.Context) {

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
		listOfObjects, err = database.GetObjects(dbConnection, targetTable, limitVal)
	} else {
		listOfObjects, err = database.GetObjects(dbConnection, targetTable)
	}
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

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
func (v *FacilityService) deleteResource(ctx *gin.Context) {

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
func (v *FacilityService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	updatingData := make(map[string]interface{})
	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
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
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Validation Failed"), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	userId := common.GetUserId(ctx)
	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource got updated", intRecordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully updated",
		Error:   0,
	})

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
func (v *FacilityService) createNewResource(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var createdRecordId int
	var err error

	updatedRequest := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	// here we should do the validation
	err = v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Validation Failed"), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}

	objectField, _ := json.Marshal(createRequest)
	err = v.ComponentManager.TemplateMandatoryFieldValidation(dbConnection, objectField, componentName)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Validation Failed"), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}
	// attach the levelCounter, this mean, no actions done yet, when the record is created, it should be created with initial record counter
	updatedRequest["levelCounter"] = 1

	rawCreateRequest, _ := json.Marshal(updatedRequest)

	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)

	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}

	err, createdRecordId = database.Create(dbConnection, targetTable, object)
	if err != nil {
		v.BaseService.Logger.Error("error creating resource", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), const_util.ErrorCreatingObjectInformation)
		return
	}

	switch componentName {
	case const_util.FacilityServiceMyRequestComponent:

		_, generalWorkOrder := database.Get(dbConnection, const_util.FacilityServiceRequestTable, createdRecordId)
		facilityJobRequestInfo := make(map[string]interface{})
		json.Unmarshal(generalWorkOrder.ObjectInfo, &facilityJobRequestInfo)

		if createdRecordId < 10 {
			facilityJobRequestInfo["serviceRequestId"] = "JR0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			facilityJobRequestInfo["serviceRequestId"] = "JR000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			facilityJobRequestInfo["serviceRequestId"] = "JR00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			facilityJobRequestInfo["serviceRequestId"] = "JR0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			facilityJobRequestInfo["serviceRequestId"] = "JR" + strconv.Itoa(createdRecordId)
		}
		facilityJobRequestInfo["templateFields"] = createRequest["templateFields"]
		facilityJobRequestInfo["actionStatus"] = const_util.ActionPendingSubmission
		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(facilityJobRequestInfo)
		updatingData["object_info"] = rawWorkOrderInfo

		database.Update(dbConnection, const_util.FacilityServiceRequestTable, createdRecordId, updatingData)
	}
	userId := common.GetUserId(ctx)
	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	ctx.JSON(http.StatusOK, component.GeneralResourceCreateResponse{
		Message: "The resource is successfully created",
		Error:   0,
		Id:      createdRecordId,
	})
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
func (v *FacilityService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	templateId := header_parser.GetQueryField(ctx, "templateId")
	var iTemplateId = -1
	if templateId != "" {
		iTemplateId, _ = strconv.Atoi(templateId)
	}

	newRecordResponse := v.ComponentManager.GetNewRecordResponse_v1(zone, dbConnection, componentName, iTemplateId)
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
func (v *FacilityService) getRecordFormData(ctx *gin.Context) {

	// first get the record
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}
	rawObjectInfo := generalObject.ObjectInfo

	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)

	if targetTable == const_util.FacilityServiceRequestTable {
		categoryInfo := make(map[string]interface{})
		json.Unmarshal(generalObject.ObjectInfo, &categoryInfo)
		templateId := util.InterfaceToInt(categoryInfo["categoryId"])
		_, generalTemplateObject := database.Get(dbConnection, const_util.FacilityServiceRequestCategoryTable, templateId)
		templateInfoObj := make(map[string]interface{})
		json.Unmarshal(generalTemplateObject.ObjectInfo, &templateInfoObj)

		response["categoryImage"] = component.RecordInfo{Value: templateInfoObj["categoryImage"]}
		response["name"] = component.RecordInfo{Value: templateInfoObj["name"]}
		response["expectedDeliveryDays"] = component.RecordInfo{Value: templateInfoObj["expectedDeliveryDays"]}
		response["description"] = component.RecordInfo{Value: templateInfoObj["description"]}
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
func (v *FacilityService) getSearchResults(ctx *gin.Context) {

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
	listOfObjects, err := database.GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), const_util.ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), const_util.ErrorGettingObjectsInformation)
			return

		}
	}

	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)

}

func (v *FacilityService) internalTableRecordOrdering(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	//projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ReferenceDatabase
	type ReorderRequest struct {
		SrcId int `json:"srcId"`
		DstId int `json:"dstId"`
	}
	var reorderRequest = ReorderRequest{}
	if err := ctx.ShouldBindBodyWith(&reorderRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// first get the main record
	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}
	// then get the internal array records using reference filed, assume now is templateFields
	var objectFields = make(map[string]interface{})
	err = json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
	if err != nil {
		v.BaseService.Logger.Error("error unmarshal resource record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting resource field information"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}
	var sourceObject interface{}
	var modifiedArray []interface{}
	if templateFields, ok := objectFields["templateFields"]; ok {
		//templateFields is the object array
		var templateObjectFields = templateFields.([]interface{})
		for _, templateField := range templateObjectFields {
			var templateIndividualFields = templateField.(map[string]interface{})
			// get the id
			if idValue, isId := templateIndividualFields["id"]; isId {
				// we took out the source from array
				if util.InterfaceToInt(idValue) != reorderRequest.SrcId {
					modifiedArray = append(modifiedArray, templateField)
				} else {
					sourceObject = templateField
				}
			}
		}

		var finalModifiedArray []interface{}
		// now run through the modified one, and add the source after destination
		for _, newTemplateField := range modifiedArray {
			var templateIndividualFields = newTemplateField.(map[string]interface{})
			if idValue, isId := templateIndividualFields["id"]; isId {
				if util.InterfaceToInt(idValue) != reorderRequest.DstId {
					finalModifiedArray = append(finalModifiedArray, newTemplateField)
				} else {
					// yes found the destination, now add the source after destination
					finalModifiedArray = append(finalModifiedArray, newTemplateField)
					finalModifiedArray = append(finalModifiedArray, sourceObject)
				}
			}
		}
		updatingData := make(map[string]interface{})
		objectFields["templateFields"] = finalModifiedArray
		serializedObject, _ := json.Marshal(objectFields)
		fmt.Println("updating modified array from it service handgle", string(serializedObject))
		updatingData["object_info"] = serializedObject
		err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating resource information"), const_util.ErrorUpdatingObjectInformation)
			return
		} else {
			ctx.JSON(http.StatusOK, component.GeneralResponse{
				Message: "Successfully updated the resource",
				Error:   0,
			})
			return
		}
	}

	response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid internal records"), const_util.ErrorGettingIndividualObjectInformation)
	return
}

func (v *FacilityService) removeInternalArrayReference(ctx *gin.Context) {

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
	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}
	serializedObject := v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject
	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	var updatingObjectFields map[string]interface{}
	json.Unmarshal(serializedObject, &updatingObjectFields)
	ctx.JSON(http.StatusOK, updatingObjectFields)

}

func (v *FacilityService) deleteValidation(ctx *gin.Context) {

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

func (v *FacilityService) DefaultGroupByCardView(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	groupByFields := ctx.Query("groupByFields")
	searchFields := ctx.Query("search")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	if searchFields != "" {

		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := v.ComponentManager.GetAbsoluteSearchQuery(componentName, searchFieldCommand)
		listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, searchQuery)
	} else {
		listOfObjects, err = database.GetObjects(dbConnection, targetTable)
	}

	if err != nil {
		response.DispatchDetailedError(ctx, 1900,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	cardViewResponseMap := v.ComponentManager.GetCardViewArrayOfMapInterface(listOfObjects, componentName)
	var groupByCardViewResponse = make([]component.GroupByCardView, 0)
	for _, responseMap := range cardViewResponseMap {
		groupByFieldValue := util.InterfaceToString(responseMap[groupByFields])

		var isElementFound bool
		isElementFound = false
		for index, mm := range groupByCardViewResponse {
			if mm.GroupByField == groupByFieldValue {
				groupByCardViewResponse[index].Cards = append(groupByCardViewResponse[index].Cards, responseMap)
				isElementFound = true
			}
		}
		if !isElementFound {
			groupByCardView := component.GroupByCardView{}
			groupByCardView.GroupByField = groupByFieldValue
			groupByCardView.Cards = append(groupByCardView.Cards, responseMap)
			groupByCardViewResponse = append(groupByCardViewResponse, groupByCardView)
		}
	}

	ctx.JSON(http.StatusOK, groupByCardViewResponse)
}

func (v *FacilityService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := database.GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(*listOfObjects)
					for _, referenceObject := range *listOfObjects {
						v.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (v *FacilityService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := database.GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					for _, referenceObject := range *listOfObjects {
						database.ArchiveObject(dbConnection, referenceTable, referenceObject)
						v.CreateUserRecordMessage(const_util.ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
						v.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
					}
				}
			}

		}
	}
}
