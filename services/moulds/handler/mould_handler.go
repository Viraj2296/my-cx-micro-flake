package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"cx-micro-flake/services/moulds/handler/workflow_actions"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

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
func (v *MouldService) loadFile(ctx *gin.Context) {
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
func (v *MouldService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	projectId := ctx.Param("projectId")

	componentName := ctx.Param("componentName")
	switch componentName {
	case const_util.MouldMasterComponent,
		const_util.MouldCategoryComponent,
		const_util.MouldSubCategoryComponent,
		const_util.MouldStatusComponent,
		const_util.MouldRequestTestComponent:
		importDataCommand := component.ImportDataCommand{}
		if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		dbConnection := v.BaseService.ServiceDatabases[projectId]
		err, errorCode, importDataResonse := v.ComponentManager.ImportData(dbConnection, componentName, importDataCommand)
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

		ctx.JSON(http.StatusOK, importDataResonse)
	default:
		response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")
	}
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
func (v *MouldService) exportObjects(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")

	// Get the current user ID and fetch user info
	userId := common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)
	userSiteId := basicUserInfo.Site // Assuming you get the site ID here
	switch componentName {
	case const_util.MouldMasterComponent,
		const_util.MouldCategoryComponent,
		const_util.MouldSubCategoryComponent,
		const_util.MouldStatusComponent,
		const_util.MouldRequestTestComponent:
		exportCommand := component.ExportDataCommand{}

		if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		dbConnection := v.BaseService.ServiceDatabases[projectId]

		// Build the condition to filter by site ID
		var condition string
		if componentName == const_util.MouldMasterComponent {
			condition = "object_info ->> '$.site' = " + strconv.Itoa(userSiteId)
		}

		err, errorCode, exportDataResponse := v.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
		if err != nil {
			v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
			return
		}

		ctx.JSON(http.StatusOK, exportDataResponse)

	default:
		response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")
	}
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
func (v *MouldService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	switch componentName {
	case const_util.MouldMasterComponent,
		const_util.MouldCategoryComponent,
		const_util.MouldSubCategoryComponent,
		const_util.MouldStatusComponent,
		const_util.MouldRequestTestComponent:
		tableImportSchema := v.ComponentManager.GetTableImportSchema(componentName)
		ctx.JSON(http.StatusOK, tableImportSchema)

	default:
		response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")
	}
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
func (v *MouldService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	switch componentName {
	case const_util.MouldMasterComponent,
		const_util.MouldCategoryComponent,
		const_util.MouldSubCategoryComponent,
		const_util.MouldStatusComponent,
		const_util.MouldRequestTestComponent:
		exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
		ctx.JSON(http.StatusOK, exportSchema)
		return
	default:
		response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")
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

func (v *MouldService) getObjects(ctx *gin.Context) {
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
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
	v.BaseService.Logger.Info("getting objects", zap.Any("limits", limitValue), zap.Any("offset", offsetValue))
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
		listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		listOfObjects, err = database.GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
	} else {

		totalRecords = database.Count(dbConnection, targetTable)
		if limitValue == "" {
			listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))
		} else {
			var conditionString string
			limitVal, _ := strconv.Atoi(limitValue)
			if componentName == const_util.MyMouldTestRequestComponent {

				// send to all in the group
				err, c := database.Get(dbConnection, const_util.MouldSettingTable, 1)
				var testUserFound bool
				testUserFound = false
				if err == nil {
					authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
					mouldSettingInfo := database.GetMouldSettingInfo(c.ObjectInfo)
					var listOfUsers = authService.GetUserInfoFromGroupId(mouldSettingInfo.MouldTestingGroup)

					for _, v := range listOfUsers {
						if userId == v.UserId {
							testUserFound = true
						}
					}
				}
				if testUserFound {
					conditionString = " id > 0 "
				} else {
					conditionString = " id < 0"
				}

			} else if componentName == const_util.MouldCategoryComponent || componentName == const_util.MouldMasterComponent || componentName == const_util.PartMasterComponent {
				// get the user site
				authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

				basicUserInfo := authService.GetUserInfoById(userId)
				v.BaseService.Logger.Info("getting user info for id", zap.Int("user_id", userId), zap.Any("user_info", basicUserInfo))
				if len(basicUserInfo.SecondarySiteList) > 0 {
					// this user has more sites
					var whereInQuery = " IN ("
					var combinedSiteList = basicUserInfo.SecondarySiteList
					combinedSiteList = append(combinedSiteList, basicUserInfo.Site)
					for _, v := range combinedSiteList {
						whereInQuery += strconv.Itoa(v) + ","
					}
					// Remove the trailing comma and close the parenthesis
					whereInQuery = strings.TrimSuffix(whereInQuery, ",") + ") "
					userBasedQuery := " object_info ->>'$.site' " + whereInQuery

					conditionString = userBasedQuery
				} else {
					userBasedQuery := " object_info ->>'$.site' =" + strconv.Itoa(basicUserInfo.Site) + " "
					conditionString = userBasedQuery
				}

				v.BaseService.Logger.Info("loaded the user site information (included secondary site list)", zap.String("conditionString", conditionString))
			}

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

				listOfObjects, err = database.GetConditionalObjectsOrderBy(dbConnection, targetTable, conditionString, orderBy, limitVal)

				var currentRecordCount = 0
				if listOfObjects != nil {
					currentRecordCount = len(*listOfObjects)
				}

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

				//listOfObjects = reverseSlice(listOfObjects)
			} else {
				listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)
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
			}
		}

	}
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.Any("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		userId = common.GetUserId(ctx)
		zone := getUserTimezone(userId)
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (v *MouldService) getSearchReset(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	outFields := ctx.Query("out_fields")
	//Have to next flag
	isNext := true
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	userId := common.GetUserId(ctx)
	limitValue, _ := ctx.Get("limit")
	totalRecords = database.Count(dbConnection, targetTable)
	var conditionString string
	limitVal := util.InterfaceToInt(limitValue)
	if componentName == const_util.MyMouldTestRequestComponent {
		userBasedQuery := " object_info ->>'$.testedBy' =" + strconv.Itoa(userId) + " "
		conditionString = userBasedQuery
	} else if componentName == const_util.MouldCategoryComponent || componentName == const_util.MouldMasterComponent {
		// get the user site
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		basicUserInfo := authService.GetUserInfoById(userId)

		if len(basicUserInfo.SecondarySiteList) > 0 {
			// this user has more sites
			var whereInQuery = " IN ("
			var combinedSiteList = basicUserInfo.SecondarySiteList
			combinedSiteList = append(combinedSiteList, basicUserInfo.Site)
			for _, v := range combinedSiteList {
				whereInQuery += strconv.Itoa(v) + ","
			}
			// Remove the trailing comma and close the parenthesis
			whereInQuery = strings.TrimSuffix(whereInQuery, ",") + ") "
			userBasedQuery := " object_info ->>'$.site' " + whereInQuery

			conditionString = userBasedQuery
		} else {
			userBasedQuery := " object_info ->>'$.site' =" + strconv.Itoa(basicUserInfo.Site) + " "
			conditionString = userBasedQuery
		}

		v.BaseService.Logger.Info("loaded the user site information for search query", zap.String("conditionString", conditionString))
	}

	listOfObjects, _ = database.GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)
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

	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
	tableObjectResponse := component.TableObjectResponse{}
	json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

	tableObjectResponse.IsNext = isNext
	tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

	ctx.JSON(http.StatusOK, tableRecordsResponse)
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
func (v *MouldService) getCardView(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")

	// Get the current user ID and fetch user info
	userId := common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)
	var customCondition string
	if len(basicUserInfo.SecondarySiteList) > 0 {
		// this user has more sites
		var whereInQuery = " IN ("
		var combinedSiteList = basicUserInfo.SecondarySiteList
		combinedSiteList = append(combinedSiteList, basicUserInfo.Site)
		for _, v := range combinedSiteList {
			whereInQuery += strconv.Itoa(v) + ","
		}
		// Remove the trailing comma and close the parenthesis
		whereInQuery = strings.TrimSuffix(whereInQuery, ",") + ") "
		userBasedQuery := " object_info ->>'$.site' " + whereInQuery

		customCondition = userBasedQuery
	} else {
		userBasedQuery := " object_info ->>'$.site' =" + strconv.Itoa(basicUserInfo.Site) + " "
		customCondition = userBasedQuery
	}

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	// Convert limitValue to integer
	var limitVal int
	if limitValue != "" {
		limitVal, _ = strconv.Atoi(limitValue)
	}

	// Prepare the query based on offset and limit
	if offsetValue != "" && limitValue != "" {
		// Construct the user-based query to filter by site ID
		listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, customCondition, limitVal)
	} else {
		// Construct the user-based query without limit
		listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, customCondition)
	}
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))

		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
		return
	}

	// Assuming GetCardViewResponse needs fewer arguments
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
func (v *MouldService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, recordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		error_util.SendResourceNotFound(ctx)
		return
	}
	// first archive all the dependency
	userId := common.GetUserId(ctx)
	v.archiveReferences(userId, dbConnection, componentName, recordId)

	err = database.ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		error_util.SendArchiveFailed(ctx)
		return
	}

	if targetTable == const_util.MouldMasterTable {
		mouldMaster := database.MouldMaster{ObjectInfo: generalObject.ObjectInfo}
		mouldMasterInfo := mouldMaster.GetMouldMasterInfo()

		mouldMasterInfo.CanModify = false
		var serialisedData = mouldMasterInfo.DatabaseSerialize(userId)
		err = database.Update(dbConnection, const_util.MouldMasterTable, recordId, serialisedData)
	}

	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource is archived, no further modification allowed", recordId, userId, nil, nil)
	ctx.Status(http.StatusNoContent)

}

// updateResource ShowAccount godoc
// @Summary update given resource based on resource id
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MouldManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   resourceId     path    string     true        "Resource Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId}/remove_internal_array_reference [put]
func (v *MouldService) removeInternalArrayReference(ctx *gin.Context) {

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
func (v *MouldService) updateResource(ctx *gin.Context) {
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

	var initializedObjectMap map[string]interface{}
	err = json.Unmarshal(initializedObject, &initializedObjectMap)

	if componentName == const_util.MouldManualShotCountComponent {
		shotCount, ok := updateRequest["shotCount"].(float64)
		if shotCount > 0 {

			mouldTestRequest := database.MouldManualShotCount{ObjectInfo: objectInterface.ObjectInfo}
			mouldTestRequestInfo := mouldTestRequest.GetMouldManualShotCountInfo()
			mouldId, _ := strconv.Atoi(mouldTestRequestInfo.MouldId)

			err, objectMasterInterface := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)
			if err != nil {
				v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
				return
			}
			var objectFields = make(map[string]interface{})
			json.Unmarshal(objectMasterInterface.ObjectInfo, &objectFields)

			err, objectMouldShotView := database.Get(dbConnection, const_util.MouldShoutCountViewTable, mouldId)
			if err != nil {
				v.BaseService.Logger.Error("error getting given record in Mould shot count", zap.String("error", err.Error()))
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record in Mould shot count"), const_util.ErrorGettingIndividualObjectInformation)
				return
			}
			var objectMouldShotViewFields = make(map[string]interface{})
			json.Unmarshal(objectMouldShotView.ObjectInfo, &objectMouldShotViewFields)

			if !ok {
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("shotCount is not a valid float64"), const_util.ErrorGettingIndividualObjectInformation)
				return
			}
			newNoOfShotCount := 0.0
			objectNoOfCav := objectFields["noOfCav"]
			objectShotCount := objectMouldShotViewFields["currentShotCount"]

			//shotCount: 12 old shot: 20
			if int(shotCount) >= mouldTestRequestInfo.ShotCount { //
				newShotCount := shotCount - float64(mouldTestRequestInfo.ShotCount) //20-12=8

				noOfCav, ok := objectNoOfCav.(float64)
				if !ok || noOfCav == 0 {
					v.BaseService.Logger.Error("invalid cavity count: cannot update quantity due to invalid or zero cavity count")
					response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid cavity count: cannot update quantity due to invalid or zero cavity count"), const_util.ErrorGettingIndividualObjectInformation)
					return
				}
				newNoOfShotCount = newShotCount / noOfCav // 8/4=2
				if currentShotCount, ok := objectShotCount.(float64); ok {
					objectMouldShotViewFields["currentShotCount"] = currentShotCount + float64(newNoOfShotCount) // 100 + 2 = 102
				} else {
					v.BaseService.Logger.Error("shotCount is not a valid float64")
					return
				}
			} else {
				newShotCount := float64(mouldTestRequestInfo.ShotCount) - shotCount // 20-12=8               // 102-2=100
				noOfCav, ok := objectNoOfCav.(float64)

				if !ok || noOfCav == 0 {
					v.BaseService.Logger.Error("invalid cavity count: cannot update quantity due to invalid or zero cavity count")
					response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid cavity count: cannot update quantity due to invalid or zero cavity count"), const_util.ErrorGettingIndividualObjectInformation)
					return
				}
				newNoOfShotCount = newShotCount / noOfCav // 8/4=2
				if currentShotCount, ok := objectShotCount.(float64); ok {
					objectMouldShotViewFields["currentShotCount"] = currentShotCount - float64(newNoOfShotCount) // 100 + 2 = 102
				} else {
					v.BaseService.Logger.Error("shotCount is not a valid float64")
					return
				}
			}

			objectFieldsJSON, err := json.Marshal(objectMouldShotViewFields)
			if err != nil {
				v.BaseService.Logger.Error("Error marshaling objectFields:", zap.String("error", err.Error()))
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error processing objectFields"), const_util.ErrorUpdatingObjectInformation)
				return
			}

			updatingMasterData := make(map[string]interface{})
			updatingMasterData["object_info"] = string(objectFieldsJSON)
			if err != nil {
				v.BaseService.Logger.Error("Error marshaling objectInfo:", zap.String("error", err.Error()))
				return
			}
			err = database.Update(v.BaseService.ReferenceDatabase, const_util.MouldShoutCountViewTable, mouldId, updatingMasterData)
			if err != nil {
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating mould register information"), const_util.ErrorUpdatingObjectInformation)
				return
			}
		} else {

			initializedObjectMap["objectStatus"] = "Archived"

		}

	}

	updatedObjectInfo, err := json.Marshal(initializedObjectMap)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error processing updated object information"), const_util.ErrorUpdatingObjectInformation)
		return
	}

	// updatingData["object_info"] = initializedObject
	updatingData["object_info"] = updatedObjectInfo

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("error updating mould master information")
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating mould master information"), const_util.ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the resource",
	})

}

// createNewResource ShowAccount godoc
// @Summary create new resource
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MouldManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   recordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/moulds/component/{componentName}/records [post]
func (v *MouldService) createNewResource(ctx *gin.Context) {
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
	rawCreateRequest, _ := json.Marshal(processedDefaultObjectMappingFields)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)

	var createdRecordId int
	var userId = common.GetUserId(ctx)
	if componentName == const_util.MouldRequestTestComponent {

		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

		mouldTestRequest := database.MouldTestRequestInfo{}
		json.Unmarshal(preprocessedRequest, &mouldTestRequest)
		err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEventByMachineId(projectId, mouldTestRequest.MachineId, mouldTestRequest.RequestTestStartDate, mouldTestRequest.RequestTestEndDate)

		err, mouldMasterObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldTestRequest.MouldId)

		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      const_util.GetError(common.ObjectNotFoundError).Error(),
					Description: "The requested resource is not exist, please check your mould identifier",
				})
			return
		}
		userId = common.GetUserId(ctx)

		var testRequestReferenceId string
		numberOfRecords := int(database.GetCount(dbConnection, const_util.MouldTestRequestTable))
		numberOfRecords = numberOfRecords + 1
		if numberOfRecords < 10 {
			testRequestReferenceId = "MTR0000" + strconv.Itoa(numberOfRecords)
		} else if numberOfRecords < 100 {
			testRequestReferenceId = "MTR000" + strconv.Itoa(numberOfRecords)
		} else if numberOfRecords < 1000 {
			testRequestReferenceId = "MTR00" + strconv.Itoa(numberOfRecords)
		} else if numberOfRecords < 10000 {
			testRequestReferenceId = "MTR0" + strconv.Itoa(numberOfRecords)
		} else if numberOfRecords < 100000 {
			testRequestReferenceId = "MTR" + strconv.Itoa(numberOfRecords)
		}

		numberOfOffToolPerformed := database.GetCount(dbConnection, const_util.MouldTestRequestTable, "JSON_EXTRACT(object_info, \"$.mouldId\") =  "+strconv.Itoa(mouldTestRequest.MouldId))
		mouldTestRequest.TestRequestReferenceId = testRequestReferenceId
		mouldTestRequest.OffTool = int(numberOfOffToolPerformed)
		mouldTestRequest.ActionStatus = const_util.MouldTestRequestActionCreated
		mouldTestRequest.MouldTestStatus = const_util.MouldTestWorkFlowPlanner
		mouldTestRequest.IsUpdate = true
		mouldTestRequest.Draggable = true
		mouldTestRequest.CanView = false
		mouldTestRequest.RequestTestStartDate = util.ConvertSingaporeTimeToUTC(mouldTestRequest.RequestTestStartDate)
		mouldTestRequest.RequestTestEndDate = util.ConvertSingaporeTimeToUTC(mouldTestRequest.RequestTestEndDate)

		existingActionRemarks := mouldTestRequest.ActionRemarks

		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "REQUEST CREATED",
			UserId:        userId,
			Remarks:       "New test request is successfully created",
			ProcessedTime: workflow_actions.GetTimeDifference(mouldTestRequest.CreatedAt),
		})
		mouldTestRequest.ActionRemarks = existingActionRemarks

		serializedObject := component.GetGeneralObjectFromInterface(mouldTestRequest)

		err, createdRecordId = database.Create(dbConnection, const_util.MouldTestRequestTable, serializedObject)

		// now update the status of the mould into test request submit to false
		mouldMaster := database.MouldMaster{ObjectInfo: mouldMasterObject.ObjectInfo}
		var mouldMasterFields = make(map[string]interface{})
		json.Unmarshal(mouldMaster.ObjectInfo, &mouldMasterFields)

		mouldMasterFields[const_util.MouldMasterFieldCanSubmitTestRequest] = false
		mouldMasterFields[const_util.MouldMasterFieldCanCreateWorkOrder] = false
		mouldMasterFields[const_util.MouldMasterFieldMouldStatus] = const_util.MouldStatusQualification
		mouldMasterFields[const_util.CommonFieldLastUpdatedBy] = userId
		mouldMasterFields[const_util.CommonFieldLastUpdatedAt] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		serialisedMouldMasterData, _ := json.Marshal(mouldMasterFields)
		updatingMouldMasterData := make(map[string]interface{})
		updatingMouldMasterData["object_info"] = serialisedMouldMasterData
		err = database.Update(dbConnection, const_util.MouldMasterTable, mouldTestRequest.MouldId, updatingMouldMasterData)

		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      const_util.GetError(common.ErrorCreatingResource).Error(),
					Description: "Internal system error occurred during new mould test request creation. Please report the error code to system admin",
				})
			return
		}

		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userInfo := authService.GetUserInfoById(mouldTestRequest.CreatedBy)
		systemNotification := common.SystemNotification{}
		eventName := mouldMaster.GetMouldMasterInfo().ToolNo + " has been used to create a mould test request"
		systemNotification.Name = eventName
		systemNotification.ColorCode = "#1257E0"
		systemNotification.IconCls = "fluent:shape-exclude-16-regular"
		systemNotification.RecordId = createdRecordId
		systemNotification.Component = "Mould Test Request"
		systemNotification.RouteLinkComponent = "mould_test_request"
		systemNotification.GeneratedTime = util.GetCurrentTime(time.RFC822)
		systemNotification.Description = "New mould test is request is created by " + userInfo.FullName
		rawSystemNotification, _ := json.Marshal(systemNotification)
		notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
		notificationService.CreateSystemNotification(projectId, rawSystemNotification)

		if scheduledEventObject != nil {
			if len(*scheduledEventObject) > 0 {
				v.CreateUserRecordMessage(const_util.ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
				ctx.JSON(http.StatusCreated, response.GeneralResponse{
					Code:     const_util.DuplicateValueFound,
					RecordId: createdRecordId,
					Message:  "The Mould Test Request has been successfully created for the requested machine, even though a previously scheduled production order exists for the same machine.",
				})
				return
			}
		}

	} else if componentName == const_util.MouldManualShotCountTable {
		var mouldId = util.InterfaceToInt(createRequest["mouldId"])
		err, mouldMasterObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)

		if err == nil {
			var mouldMasterInfo = make(map[string]interface{})
			err = json.Unmarshal(mouldMasterObject.ObjectInfo, &mouldMasterInfo)

			if err == nil {
				err, objectMouldShotView := database.Get(dbConnection, const_util.MouldShoutCountViewTable, mouldId)
				if err != nil {
					v.BaseService.Logger.Error("error getting given record in Mould shot count", zap.String("error", err.Error()))
					response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record in Mould shot count"), const_util.ErrorGettingIndividualObjectInformation)
					return
				}
				if objectMouldShotView.ObjectInfo == nil {
					newRecord := database.MouldShotCountViewInfo{
						CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
						UpdatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
						CreatedBy:          userId,
						UpdatedBy:          userId,
						CurrentShotCount:   util.InterfaceToInt(createRequest["shotCount"]),
						IsNotificationSent: false,
					}
					shotCountObject := component.GeneralObject{Id: mouldId, ObjectInfo: newRecord.Serialised()}

					err, _ = database.Create(dbConnection, const_util.MouldShoutCountViewTable, shotCountObject)
					if err != nil {
						v.BaseService.Logger.Error("error creating new mould shout count record", zap.String("error", err.Error()))
					}
				} else {
					var objectMouldShotViewFields = make(map[string]interface{})
					json.Unmarshal(objectMouldShotView.ObjectInfo, &objectMouldShotViewFields)

					var shotCount = util.InterfaceToFloat(objectMouldShotViewFields["currentShotCount"])
					var noOfCav = util.InterfaceToFloat(mouldMasterInfo["noOfCav"])
					if noOfCav == 0 {
						noOfCav = 1
					}
					if shotCount > 0 {
						shotCount = shotCount + (util.InterfaceToFloat(createRequest["shotCount"]) / noOfCav)

						objectMouldShotViewFields["currentShotCount"] = shotCount

						serialisedMouldShotViewInfo, _ := json.Marshal(objectMouldShotViewFields)
						updatingMouldShotViewData := make(map[string]interface{})
						updatingMouldShotViewData["object_info"] = string(serialisedMouldShotViewInfo)
						err = database.Update(dbConnection, const_util.MouldShoutCountViewTable, mouldId, updatingMouldShotViewData)

						if err != nil {
							v.BaseService.Logger.Error("error updating mould master after shot count", zap.String("error", err.Error()))
						}
					} else {
						v.BaseService.Logger.Info("shot count isn't greater than 0", zap.Any("mould Id", mouldId))
					}
				}

			} else {
				v.BaseService.Logger.Error("error updating mould master", zap.String("error", err.Error()))
			}
		} else {
			v.BaseService.Logger.Error("error getting mould master", zap.String("error", err.Error()))
		}

		object := component.GeneralObject{
			ObjectInfo: preprocessedRequest,
		}

		err, createdRecordId = database.Create(dbConnection, targetTable, object)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error in Resource Creation "), const_util.ErrorCreatingObjectInformation, err.Error())
			return
		}
	} else {
		object := component.GeneralObject{
			ObjectInfo: preprocessedRequest,
		}

		err, createdRecordId = database.Create(dbConnection, targetTable, object)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error in Resource Creation "), const_util.ErrorCreatingObjectInformation, err.Error())
			return
		}
	}

	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:     0,
		RecordId: createdRecordId,
		Message:  "New resource is successfully created",
	})

	/*
		else if componentName == const_util.MouldManualShotCountTable {
				var mouldId = util.InterfaceToInt(createRequest["mouldId"])
				err, mouldMasterObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)

				if err == nil {
					var mouldMasterInfo = make(map[string]interface{})
					err = json.Unmarshal(mouldMasterObject.ObjectInfo, &mouldMasterInfo)

					if err == nil {
						err, objectMouldShotView := database.Get(dbConnection, const_util.MouldShoutCountViewTable, mouldId)
						if err != nil {
							v.BaseService.Logger.Error("error getting given record in Mould shot count", zap.String("error", err.Error()))
							response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record in Mould shot count"), const_util.ErrorGettingIndividualObjectInformation)
							return
						}
						if objectMouldShotView.ObjectInfo == nil {
							newRecord := database.MouldShotCountViewInfo{
								CreatedAt:          time.Now(),
								UpdatedAt:          time.Now(),
								CreatedBy:          userId,
								UpdatedBy:          userId,
								CurrentShotCount:   util.InterfaceToInt(createRequest["shotCount"]),
								IsNotificationSent: false,
							}
							shotCountObject := component.GeneralObject{Id: mouldId, ObjectInfo: newRecord.Serialised()}

							err, _ = database.Create(dbConnection, const_util.MouldShoutCountViewTable, shotCountObject)
							if err != nil {
								v.BaseService.Logger.Error("error creating new mould shout count record", zap.String("error", err.Error()))
							}
						} else {
							var objectMouldShotViewFields = make(map[string]interface{})
							json.Unmarshal(objectMouldShotView.ObjectInfo, &objectMouldShotViewFields)

							var shotCount = util.InterfaceToFloat(objectMouldShotViewFields["currentShotCount"])
							var noOfCav = util.InterfaceToFloat(mouldMasterInfo["noOfCav"])
							if noOfCav == 0 {
								noOfCav = 1
							}
							if shotCount > 0 {
								shotCount = shotCount + (util.InterfaceToFloat(createRequest["shotCount"]) / noOfCav)

								objectMouldShotViewFields["currentShotCount"] = shotCount

								serialisedMouldShotViewInfo, _ := json.Marshal(objectMouldShotViewFields)
								updatingMouldShotViewData := make(map[string]interface{})
								updatingMouldShotViewData["object_info"] = string(serialisedMouldShotViewInfo)
								err = database.Update(dbConnection, const_util.MouldShoutCountViewTable, mouldId, updatingMouldShotViewData)

								if err != nil {
									v.BaseService.Logger.Error("error updating mould master after shot count", zap.String("error", err.Error()))
								}
							} else {
								v.BaseService.Logger.Info("shot count isn't greater than 0", zap.Any("mould Id", mouldId))
							}
						}

					} else {
						v.BaseService.Logger.Error("error updating mould master", zap.String("error", err.Error()))
					}
				} else {
					v.BaseService.Logger.Error("error getting mould master", zap.String("error", err.Error()))
				}

				object := component.GeneralObject{
					ObjectInfo: preprocessedRequest,
				}

				err, createdRecordId = database.Create(dbConnection, targetTable, object)
				if err != nil {
					response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error in Resource Creation "), const_util.ErrorCreatingObjectInformation, err.Error())
					return
				}
			}
	*/
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
func (v *MouldService) getNewRecord(ctx *gin.Context) {
	componentName := util.GetComponentName(ctx)
	projectId := util.GetProjectId(ctx)
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
// @Router /project/{projectId}/moulds/component/{componentName}/record/{recordId} [get]
func (v *MouldService) getRecordFormData(ctx *gin.Context) {
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
	type resourcePermission struct {
		CanDelete bool `json:"canDelete"`
		CanUpdate bool `json:"canUpdate"`
	}
	var resourcePermissions resourcePermission
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)
	if componentName == const_util.MouldMasterComponent || componentName == const_util.MouldCategoryComponent {
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		basicUserInfo := authService.GetUserInfoById(userId)
		var objectFields = make(map[string]interface{})
		json.Unmarshal(generalObject.ObjectInfo, &objectFields)
		if siteInterface, ok := objectFields["site"]; ok {
			var siteId = util.InterfaceToInt(siteInterface)
			if util.HasInt(siteId, basicUserInfo.SecondarySiteList) {
				v.BaseService.Logger.Info("this site is belong to secondary one, we will have to send the record permission denied")

				resourcePermissions = resourcePermission{
					CanDelete: false,
					CanUpdate: false,
				}
			} else {
				if siteId == basicUserInfo.Site {
					v.BaseService.Logger.Info("this site is belongs to user primary site")
					resourcePermissions = resourcePermission{
						CanDelete: true,
						CanUpdate: true,
					}
				} else {
					v.BaseService.Logger.Warn("this site is not belongs to user primary site", zap.Int("siteId", siteId), zap.Int("user_site_id", basicUserInfo.Site))
					resourcePermissions = resourcePermission{
						CanDelete: false,
						CanUpdate: false,
					}
				}
			}
		} else {
			resourcePermissions = resourcePermission{
				CanDelete: true,
				CanUpdate: true,
			}
		}
	} else {
		resourcePermissions = resourcePermission{
			CanDelete: true,
			CanUpdate: true,
		}
	}
	response["resourcePermission"] = resourcePermissions
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
// @Router /project/{projectId}/moulds/component/{componentId}/search [post]
func (v *MouldService) getSearchResults(ctx *gin.Context) {
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
	v.BaseService.Logger.Info("search fields command size", zap.Any("search", len(searchFieldCommand)))
	format := ctx.Query("format")
	if len(searchFieldCommand) == 0 {
		// reset the search
		ctx.Set("offset", "1")
		ctx.Set("limit", "30")
		if format == "card_view" {
			if componentName == const_util.MouldMasterComponent {
				authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
				basicUserInfo := authService.GetUserInfoById(userId)
				var siteCondition string
				if len(basicUserInfo.SecondarySiteList) > 0 {
					// this user has more sites
					var whereInQuery = " IN ("
					var combinedSiteList = basicUserInfo.SecondarySiteList
					combinedSiteList = append(combinedSiteList, basicUserInfo.Site)
					for _, v := range combinedSiteList {
						whereInQuery += strconv.Itoa(v) + ","
					}
					// Remove the trailing comma and close the parenthesis
					whereInQuery = strings.TrimSuffix(whereInQuery, ",") + ") "
					userBasedQuery := " object_info ->>'$.site' " + whereInQuery

					siteCondition = userBasedQuery
				} else {
					userBasedQuery := " object_info ->>'$.site' =" + strconv.Itoa(basicUserInfo.Site) + " "
					siteCondition = userBasedQuery
				}

				listOfObjects, err := database.GetConditionalObjects(dbConnection, targetTable, siteCondition)
				if err != nil {
					response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), const_util.ErrorGettingObjectsInformation)
					return
				}
				cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
				ctx.JSON(http.StatusOK, cardViewResponse)
				return
			} else {
				listOfObjects, err := database.GetObjects(dbConnection, targetTable)
				if err != nil {
					response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), const_util.ErrorGettingObjectsInformation)
					return
				}
				cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
				ctx.JSON(http.StatusOK, cardViewResponse)
				return
			}

		} else {
			v.getSearchReset(ctx)
		}

		return
	}

	searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)

	if componentName == const_util.MouldMasterComponent || componentName == const_util.PartMasterComponent {
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		basicUserInfo := authService.GetUserInfoById(userId)
		var conditionQuery = ""
		if len(basicUserInfo.SecondarySiteList) > 0 {
			// this user has more sites
			var whereInQuery = " IN ("
			var combinedSiteList = basicUserInfo.SecondarySiteList
			combinedSiteList = append(combinedSiteList, basicUserInfo.Site)
			for _, v := range combinedSiteList {
				whereInQuery += strconv.Itoa(v) + ","
			}
			// Remove the trailing comma and close the parenthesis
			whereInQuery = strings.TrimSuffix(whereInQuery, ",") + ") "
			userBasedQuery := " object_info ->>'$.site' " + whereInQuery

			conditionQuery = userBasedQuery
		} else {
			userBasedQuery := " object_info ->>'$.site' =" + strconv.Itoa(basicUserInfo.Site) + " "
			conditionQuery = userBasedQuery
		}

		searchQuery = searchQuery + " AND " + conditionQuery
	}

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

func (v *MouldService) deleteValidation(ctx *gin.Context) {

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

func (v *MouldService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
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

func (v *MouldService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
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
						fmt.Println("referenceTable : ", referenceTable, " id :", referenceObject)
						database.ArchiveObject(dbConnection, referenceTable, referenceObject)
						v.CreateUserRecordMessage(const_util.ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
						v.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
					}
				}
			}

		}
	}
}
