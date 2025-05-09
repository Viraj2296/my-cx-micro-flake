package handler

import (
	"cx-micro-flake/pkg/auth"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (as *AuthService) loadFile(ctx *gin.Context) {
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

func (as *AuthService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")

	componentName := ctx.Param("componentName")
	//targetTable := as.ComponentManager.GetTargetTable(componentName)
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ReferenceDatabase
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
	//		as.BaseService.Logger.Error("unable to create record", zap.String("error",err.Error()))
	//		failedRecords = failedRecords + 1
	//	}
	//	recordIdInString := strconv.Itoa(recordId)
	//	CreateBotRecordTrail(projectId, recordIdInString, componentName, "machine master is created")
	//}
	//importDataResponse := common.ImportDataResponse{
	//	TotalRecords:  totalRecords,
	//	FailedRecords: failedRecords,
	//	Message:       "data is successfully imported",
	//}
	//
	//ctx.JSON(http.StatusOK, importDataResponse)
}

func (as *AuthService) exportObjects(ctx *gin.Context) {
	// we will get the uploaded url
	componentName := ctx.Param("componentName")
	exportCommand := component.ExportDataCommand{}

	if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ReferenceDatabase
	var condition string
	err, errorCode, exportDataResponse := as.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
	if err != nil {
		as.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	ctx.JSON(http.StatusOK, exportDataResponse)
}

func (as *AuthService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := as.ComponentManager.GetTableImportSchema(componentName)
	ctx.JSON(http.StatusOK, tableImportSchema)
}

func (as *AuthService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := as.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)
}

func (as *AuthService) getObjects(ctx *gin.Context) {
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
	objectStatus := ctx.Query("objectStatus")
	filter := ctx.Query("filter")
	dbConnection := as.BaseService.ReferenceDatabase
	var listOfObjects []component.GeneralObject
	var totalRecords int64
	var err error
	isNext := true
	userId := common.GetUserId(ctx)
	zone := as.GetUserTimezone(userId)
	var componentCondition = ""
	if componentName == "user" {
		componentCondition = " object_info ->>'$.type' = 'user'"
	}

	statusQuery := " object_info ->>'$.objectStatus' = '" + objectStatus + "'"

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
		var searchWithBaseQuery = ""
		if componentCondition == "" {
			searchWithBaseQuery = searchQuery + " AND " + baseCondition
		} else {
			searchWithBaseQuery = searchQuery + " AND " + baseCondition + " AND " + componentCondition
		}

		err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		if objectStatus == "" {
			if componentCondition == "" {
				err, listOfObjects = GetObjects(dbConnection, targetTable)
				totalRecords = int64(len(listOfObjects))
			} else {
				err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, componentCondition)
				totalRecords = int64(len(listOfObjects))
			}

		} else {
			if componentCondition == "" {
				err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, statusQuery)
				totalRecords = int64(len(listOfObjects))
			} else {
				err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, statusQuery+" AND "+componentCondition)
				totalRecords = int64(len(listOfObjects))
			}

		}

	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			if objectStatus == "" {
				err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))
			} else {
				baseCondition := component.TableCondition(offsetValue, fields, values, condition)
				if componentCondition == "" {
					err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, baseCondition+" AND "+statusQuery)
				} else {
					err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, baseCondition+" AND "+statusQuery+" AND "+componentCondition)
				}

			}

		} else {
			var baseCondition = ""

			if componentCondition == "" {
				baseCondition = component.TableCondition(offsetValue, fields, values, condition)
			} else {
				baseCondition = component.TableCondition(offsetValue, fields, values, condition) + " AND " + componentCondition
			}
			if filter != "" {
				baseCondition = baseCondition + " AND " + component.TableConditionForFilter(filter)
			}

			if objectStatus != "" {
				baseCondition += " AND " + statusQuery
			}
			limitVal, _ := strconv.Atoi(limitValue)
			err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, baseCondition, limitVal)

			currentRecordCount := len(listOfObjects)

			if currentRecordCount < limitVal {
				isNext = false
			} else if currentRecordCount == limitVal {
				andClauses := strings.Split(baseCondition, "AND")
				var totalRecordObjects []component.GeneralObject
				if len(andClauses) > 1 {
					_, totalRecordObjects = GetConditionalObjects(dbConnection, targetTable, baseCondition)

				} else {
					_, totalRecordObjects = GetObjects(dbConnection, targetTable)
				}
				lenTotalRecord := len(totalRecordObjects)
				if (listOfObjects)[currentRecordCount-1].Id == (totalRecordObjects)[lenTotalRecord-1].Id {
					isNext = false
				}
			}
		}

	}
	as.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := as.ComponentManager.TableRecordsToArrayV1(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		_, tableRecordsResponse := as.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)

		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (as *AuthService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := as.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")

	dbConnection := as.BaseService.ReferenceDatabase
	err, listOfObjects := GetObjects(dbConnection, targetTable)
	as.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := as.ComponentManager.GetCardViewResponseV1(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

}

func (as *AuthService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	recordId := ctx.Param("recordId")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := as.BaseService.ReferenceDatabase
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	err = ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		as.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
		return
	}
	ctx.Status(http.StatusNoContent)

}

func (as *AuthService) removeInternalArrayReference(ctx *gin.Context) {

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

	updatingData["object_info"] = as.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	err = Update(as.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}
	as.loadUserAccess()
	ctx.JSON(http.StatusOK, updatingData)

}

func (as *AuthService) updateResource(ctx *gin.Context) {

	componentName := util.GetComponentName(ctx)
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	recordId := util.GetRecordId(ctx)
	dbConnection := as.BaseService.ReferenceDatabase

	var updateRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if componentName == UserComponent {
		as.BaseService.Logger.Info("Converted email to lowercase")
		updateRequest["email"] = util.ToLowerCase(util.InterfaceToString(updateRequest["email"]))
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, recordId)
	if err != nil {
		as.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	modifiedUpdatedData := common.UpdateMetaInfoFromInterface(updateRequest, ctx)
	updatingData["object_info"] = as.ComponentManager.GetUpdateRequest(modifiedUpdatedData, objectInterface.ObjectInfo, componentName)

	err = Update(as.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}

	//Update user role list when user list or role list updated in group object
	if componentName == UserGroupComponent {
		existingUserGroupInfo := GroupInfo{}
		json.Unmarshal(objectInterface.ObjectInfo, &existingUserGroupInfo)

		requestUserList := util.InterfaceToIntArray(updateRequest["users"])
		existingUserList := existingUserGroupInfo.Users
		requestRoleList := util.InterfaceToIntArray(updateRequest["roles"])
		existingRoleList := existingUserGroupInfo.Roles

		diffUserList := util.DifferenceUsers(existingUserList, requestUserList)
		diffRoleList := util.Difference(existingRoleList, requestRoleList)

		if (len(diffUserList) != 0) || (len(diffRoleList) != 0) {
			for _, userId := range diffUserList {
				var appendUserRoles []int
				err, userObjectInterface := Get(dbConnection, UserTable, userId)
				if err != nil {
					as.BaseService.Logger.Error("error in getting user by userId", zap.String("error", err.Error()))
					continue
				}
				var userInfo = make(map[string]interface{})
				json.Unmarshal(userObjectInterface.ObjectInfo, &userInfo)
				if len(existingUserList) > len(requestUserList) {
					appendUserRoles = removeValue(util.InterfaceToIntArray(userInfo["userRoles"]), requestRoleList)
				} else {
					appendUserRoles = append(util.InterfaceToIntArray(userInfo["userRoles"]), requestRoleList...)
				}

				userInfo["userRoles"] = util.RemoveDuplicateInt(appendUserRoles)

				modifiedUpdatedData := common.UpdateMetaInfoFromInterface(userInfo, ctx)
				updatingData["object_info"] = as.ComponentManager.GetUpdateRequest(modifiedUpdatedData, objectInterface.ObjectInfo, componentName)

				err = Update(as.BaseService.ReferenceDatabase, UserTable, userId, updatingData)

				if err != nil {
					as.BaseService.Logger.Error("error in updateing user info", zap.String("error", err.Error()))

					continue
				}
			}
		}
	}

	as.loadUserAccess()
	ctx.JSON(http.StatusOK, updateRequest)

}
func removeValue(slice []int, values []int) []int {
	var result []int
	valueSet := make(map[int]struct{})

	for _, v := range values {
		valueSet[v] = struct{}{}
	}

	for _, v := range slice {
		if _, found := valueSet[v]; !found {
			result = append(result, v)
		}
	}

	return result
}

func (as *AuthService) getGroupByCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := as.ComponentManager.GetTargetTable(componentName)
	projectId := util.GetProjectId(ctx)
	groupByFields := ctx.Query("groupByFields")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	searchFields := ctx.Query("search")

	dbConnection := as.BaseService.ReferenceDatabase
	var listOfObjects []component.GeneralObject
	var err error

	if searchFields != "" {

		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := as.ComponentManager.GetAbsoluteSearchQuery(componentName, searchFieldCommand)
		err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, searchQuery)
	} else {
		err, listOfObjects = GetObjects(dbConnection, targetTable)
	}

	as.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		as.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponseMap := as.ComponentManager.GetCardViewArrayOfMapInterfaceV1(listOfObjects, componentName)
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

func (as *AuthService) createNewResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	dbConnection := as.BaseService.ReferenceDatabase
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if componentName == UserComponent {
		as.BaseService.Logger.Info("Converted email to lowercase")
		createRequest["email"] = util.ToLowerCase(util.InterfaceToString(createRequest["email"]))
	}

	updatedRequest := as.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	// here we should do the validation
	err := as.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	updatedRequest["type"] = UserTypeUser

	formattedRequest := as.ComponentManager.ApplyFormatting(updatedRequest, componentName)

	rawCreateRequest, _ := json.Marshal(formattedRequest)
	modifiedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
	object := component.GeneralObject{
		ObjectInfo: modifiedRequest,
	}
	err, createRecordId := CreateFromGeneralObject(dbConnection, targetTable, object)
	if err != nil {
		response.DispatchDetailedError(ctx, common.ResourceCreationFailed,
			&response.DetailedError{
				Header:      getError(common.ErrorCreatingResource).Error(),
				Description: "Internal system error happened while creating new resource, Please try again later or raise issue ticket to system admin",
			})
		return
	}

	if componentName == UserComponent {
		err, groupObject := Get(dbConnection, UserGroupTable, as.ProfileGroupId)
		if err != nil {
			as.BaseService.Logger.Error("error getting group object", zap.String("error", err.Error()))
		} else {
			as.BaseService.Logger.Info("user profile id in creation", zap.Any("ProfileGroupId", as.ProfileGroupId), zap.Any("loaded_group_info", string(groupObject.ObjectInfo)))
			groupInfo := GroupInfo{}
			err := json.Unmarshal(groupObject.ObjectInfo, &groupInfo)
			if err != nil {
				as.BaseService.Logger.Error("error unmarshalling GroupInfo", zap.Any("GroupInfo", string(groupObject.ObjectInfo)))
			} else {
				groupInfo.Users = append(groupInfo.Users, createRecordId)
				err = Update(as.BaseService.ReferenceDatabase, UserGroupTable, groupObject.Id, groupInfo.DatabaseSerialize())
				if err != nil {
					as.BaseService.Logger.Error("Updating group info Due to:", zap.String("errors", err.Error()))
				}
			}
		}

	} else if componentName == APIAccessComponent {
		userId := common.GetUserId(ctx)
		zone := as.GetUserTimezone(userId)
		fmt.Print("modifiedRequest : ", string(modifiedRequest))
		apiAccessInfo := GetAPIAccessInfo(modifiedRequest)
		as.BaseService.Logger.Info("api access info", zap.Any("api_access_info", apiAccessInfo))
		var claims = make(map[string]interface{})
		claims["endPoints"] = apiAccessInfo.EndPoints
		claims["endDate"] = apiAccessInfo.EndDate //2024-08-21T14:00:00.000Z
		claims["id"] = apiAccessInfo.UserId
		claims["type"] = "api_access"

		token, err := auth.CreateCustomToken(ProjectID, zone, claims)
		if err == nil {
			as.BaseService.Logger.Info("Created custom token", zap.Any("token", token))
			// write the token in the object itself
			apiAccessInfo.APIKey = token
			var updateData = make(map[string]interface{})
			updateData["object_info"] = apiAccessInfo.serialised()
			err := Update(dbConnection, APIAccessTable, createRecordId, updateData)
			if err != nil {
				as.BaseService.Logger.Error("error updating the api access key info", zap.Error(err))
			}
		} else {
			as.BaseService.Logger.Error("error creating api access key info", zap.Error(err))
		}
	}
	as.loadUserAccess()
	as.BaseService.Logger.Info("generated object", zap.Any("object_details", string(modifiedRequest)))
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:    0,
		Message: "Resource has been successfully created",
	})
}

func (as *AuthService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	dbConnection := as.BaseService.ReferenceDatabase
	userId := common.GetUserId(ctx)
	zone := as.GetUserTimezone(userId)
	newRecordResponse := as.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

func (as *AuthService) getRecordFormData(ctx *gin.Context) {

	// first get the record
	// projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	recordId := ctx.Param("recordId")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := as.BaseService.ReferenceDatabase

	//userId := common.GetUserId(ctx)
	//if !as.IsAllowed(userId, projectId, componentName, "/records/{record_id}") {
	//	response.SendValidationError(ctx, http.StatusBadRequest, getError(AccessDenied), ErrorGettingObjectsInformation, AccessDeniedDescription)
	//	return
	//}

	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}
	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := as.GetUserTimezone(userId)
	response := as.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)

	ctx.JSON(http.StatusOK, response)
}

func (as *AuthService) getSearchResults(ctx *gin.Context) {

	var searchFieldCommand []component.SearchKeys
	componentName := ctx.Param("componentName")
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ReferenceDatabase
	userId := common.GetUserId(ctx)
	zone := as.GetUserTimezone(userId)
	format := ctx.Query("format")
	if format != "" {
		if format == "card_view" {
			if len(searchFieldCommand) == 0 {
				// reset the search
				ctx.Set("offset", 1)
				ctx.Set("limit", 30)
				as.getObjects(ctx)
				return
			}

			searchList := as.ComponentManager.GetSearchQueryV2(dbConnection, componentName, searchFieldCommand)
			searchQuery := " id in " + searchList
			if searchList == "()" {
				cardViewResponse := component.CardViewResponse{}
				ctx.JSON(http.StatusOK, cardViewResponse)
				return
			}
			err, listOfObjects := GetConditionalObjects(dbConnection, targetTable, searchQuery)
			// fmt.Println("listOfObjects:", len(*listOfObjects))
			if err != nil {
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
				return
			}
			cardViewResponse := as.ComponentManager.GetCardViewResponseV1(listOfObjects, componentName)
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
		as.getObjects(ctx)
		return
	}

	searchList := as.ComponentManager.GetSearchQueryV2(dbConnection, componentName, searchFieldCommand)

	listOfObjects := make([]component.GeneralObject, 0)
	var err error
	if searchList != "()" {
		searchQuery := " id in " + searchList
		err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, searchQuery)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
			return
		}
	}

	_, searchResponse := as.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, int64(len(listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)
}

func (as *AuthService) sendEmailInvitation(ctx *gin.Context) {

	dbConnection := as.BaseService.ReferenceDatabase
	recordId := util.GetRecordId(ctx)
	err, generalUserObject := Get(dbConnection, UserTable, recordId)

	if err != nil {
		as.BaseService.Logger.Error("Loading users had failed due to", zap.String("errors", err.Error()))
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid User"), ErrorGettingIndividualObjectInformation, err.Error())
		return
	}

	userInfo := GetUserInfo(generalUserObject.ObjectInfo)
	if userInfo.PlainPassword == "" {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Empty Password"), ErrorGettingIndividualObjectInformation, "Please generate new password using reset password action before sending invitation. You can not send invitation when the user is already accessed the system using his own password")
		return
	}
	as.BaseService.Logger.Info("user object ", zap.Any("user", userInfo))

	as.emailGenerator(dbConnection, WelcomeMESEmailTemplateType, recordId, UserComponent, recordId)
	// user has successfully logged in, so update the last login and last active
	userInfo.InvitationStatus = "Invited"
	userInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	Update(as.BaseService.ReferenceDatabase, UserTable, generalUserObject.Id, userInfo.DatabaseSerialize())

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Invitation link is sent to your registered email address",
	})
	return

}
func (as *AuthService) sendEmailToUser(ctx *gin.Context) {
	// projectId := ctx.Param("projectId")
	//componentName := ctx.Param("componentName")
	//recordId := ctx.Param("recordId")
	//targetTable := as.ComponentManager.GetTargetTable(componentName)
	//intRecordId, _ := strconv.Atoi(recordId)
	//dbConnection := as.BaseService.ReferenceDatabase
	//
	//err, generalObject := Get(dbConnection, targetTable, intRecordId)
	//
	//if err != nil {
	//	response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
	//	return
	//}
	//
	//// send an email now
	//emailConditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.emailType\")) = \"invitation_email\""
	//listOfSetting, err := GetConditionalObjects(as.BaseService.ReferenceDatabase, UserEmailSettingTable, emailConditionString)
	//if err != nil {
	//	response.SendDetailedError(ctx, http.StatusBadRequest, getError("System Error"), SystemError, "Unable to load internal system configuration, system admin not configured email setting")
	//}
	//if len(*listOfSetting) > 1 {
	//	response.SendDetailedError(ctx, http.StatusBadRequest, getError("System Error"), SystemError, "Invalid system configuration to process this request. Please report to admin")
	//}
	//
	//user := User{
	//	Id:         generalObject.Id,
	//	ObjectInfo: generalObject.ObjectInfo,
	//}
	//for _, emailSettingObject := range *listOfSetting {
	//	userEmailSetting := UserEmailSetting{
	//		Id:         emailSettingObject.Id,
	//		ObjectInfo: emailSettingObject.ObjectInfo,
	//	}
	//	emailTemplate := userEmailSetting.getEmailSettingInfo().Template
	//	rawUserObject := generalObject.ObjectInfo
	//	var userObjectFields map[string]interface{}
	//	json.Unmarshal(rawUserObject, &userObjectFields)
	//	for _, replaceFields := range userEmailSetting.getEmailSettingInfo().MappingFields {
	//		if val, ok := userObjectFields[replaceFields.ObjectField]; ok {
	//			emailTemplate = strings.Replace(emailTemplate, replaceFields.TemplateField, util.InterfaceToString(val), 1)
	//		}
	//
	//	}
	//	var emailMessages []common.Message
	//	emailMessage := common.Message{
	//		To:          []string{user.GetUserInfo().Email},
	//		SingleEmail: false,
	//		Subject:     userEmailSetting.getEmailSettingInfo().Subject,
	//		Body: map[string]string{
	//			"text/html": emailTemplate,
	//		},
	//		Info:          "",
	//		ReplyTo:       nil,
	//		EmbeddedFiles: nil,
	//		AttachedFiles: nil,
	//	}
	//	emailMessages = append(emailMessages, emailMessage)
	//	notificationService := services.GetService("notification").ServiceInterface.(services.NotificationInterface)
	//	notificationService.CreateMessages("906d0fd569404c59956503985b330132", emailMessages)
	//}
	//
	//// now update the status as invitation_sent
	//
	//// user has successfully logged in, so update the last login and last active
	//userInfo := user.GetUserInfo()
	//userInfo.InvitationStatus = "invitation_sent"
	//Update(as.BaseService.ReferenceDatabase, UserTable, generalObject.Id, userInfo.DatabaseSerialize())

}
