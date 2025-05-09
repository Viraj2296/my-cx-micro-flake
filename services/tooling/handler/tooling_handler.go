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
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

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
func (v *ToolingService) getObjects(ctx *gin.Context) {

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
	orderValue := ctx.Query("order")
	searchFields := ctx.Query("search")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	userBasedQuery := " JSON_EXTRACT(object_info, \"$.assignedUserId\")=" + strconv.Itoa(userId) + " "
	var listOfObjects *[]component.GeneralObject

	//Have to next flag
	isNext := true

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

		if componentName == ToolingProjectMyTaskComponent {
			searchWithBaseQuery = searchWithBaseQuery + " AND " + userBasedQuery
		}
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		if componentName == ToolingProjectMyTaskComponent {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, userBasedQuery)
			totalRecords = int64(len(*listOfObjects))
		} else {
			listOfObjects, err = GetObjects(dbConnection, targetTable)
			totalRecords = int64(len(*listOfObjects))
		}

	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			if componentName == ToolingProjectMyTaskComponent {
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition)+" AND "+userBasedQuery)
			} else {
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))
			}

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			var conditionString string

			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				offsetValue = strconv.Itoa(int(totalRecords) - limitVal + 1)

				limitVal = limitVal - offsetVal
				conditionString = component.TableCondition(offsetValue, fields, values, condition)
				if componentName == ToolingProjectMyTaskComponent {
					conditionString = component.TableCondition(offsetValue, fields, values, condition) + " AND " + userBasedQuery
				}
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)
				listOfObjects = reverseSlice(listOfObjects)

				currentRecordCount := len(*listOfObjects)

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

					if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[0].Id {
						isNext = false
					}
				}

			} else {
				if componentName == ToolingProjectMyTaskComponent {
					conditionString = component.TableCondition(offsetValue, fields, values, condition) + " AND " + userBasedQuery
				}

				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)

				currentRecordCount := len(*listOfObjects)

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

			//listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)

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
func (v *ToolingService) getCardView(ctx *gin.Context) {

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

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

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
func (v *ToolingService) getNewRecord(ctx *gin.Context) {

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
func (v *ToolingService) getRecordFormData(ctx *gin.Context) {

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
func (v *ToolingService) createNewResource(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")

	userId, _ := auth.ExtractRefreshTokenID(ctx.Request)

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// before going further check we have notification email is configured
	if componentName == ToolingProjectTaskComponent {
		if !isTaskAssignNotificationEmailConfigured(dbConnection) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Task Creation Failed",
					Description: "Please consider configuring notification email setting acknowledgement and routing email templates to proceed further",
				})
			return
		}
	}

	var createdRecordId int
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

	err, createdRecordId = Create(dbConnection, targetTable, object)
	if err != nil {
		v.BaseService.Logger.Error("error creating resource", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
		return
	}

	switch componentName {
	case ToolingProjectComponent:

		_, generalWorkOrder := Get(dbConnection, ToolingProjectTable, createdRecordId)

		//toolingProject := ToolingProject{ObjectInfo: generalWorkOrder.ObjectInfo}
		toolingProjectInfo := make(map[string]interface{})
		json.Unmarshal(generalWorkOrder.ObjectInfo, &toolingProjectInfo)
		//projectPreferenceId := util.InterfaceToString(toolingProjectInfo["projectReferenceId"])
		if createdRecordId < 10 {
			toolingProjectInfo["projectReferenceId"] = "TP0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			toolingProjectInfo["projectReferenceId"] = "TP000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			toolingProjectInfo["projectReferenceId"] = "TP00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			toolingProjectInfo["projectReferenceId"] = "TP0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			toolingProjectInfo["projectReferenceId"] = "TP" + strconv.Itoa(createdRecordId)
		}

		toolingProjectInfo["createdAt"] = util.GetCurrentTime(ISOTimeLayout)
		toolingProjectInfo["lastUpdatedAt"] = util.GetCurrentTime(ISOTimeLayout)

		existingActionRemarks := make([]ActionRemarks, 0)
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "PROJECT IS CREATED",
			UserId:        userId,
			Remarks:       "The project is successfully created",
			ProcessedTime: getTimeDifference(util.InterfaceToString(toolingProjectInfo["createdAt"])),
		})
		toolingProjectInfo["actionRemarks"] = existingActionRemarks

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(toolingProjectInfo)
		updatingData["object_info"] = rawWorkOrderInfo

		Update(dbConnection, ToolingProjectTable, createdRecordId, updatingData)
	case ToolingProjectTaskComponent:

		_, projectTask := Get(dbConnection, ToolingProjectTaskTable, createdRecordId)

		toolingProjectTask := ToolingProjectTask{ObjectInfo: projectTask.ObjectInfo}
		toolingProjectTaskInfo := toolingProjectTask.getToolingTaskInfo()
		if createdRecordId < 10 {
			toolingProjectTaskInfo.TaskReferenceId = "TPT0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			toolingProjectTaskInfo.TaskReferenceId = "TPT000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			toolingProjectTaskInfo.TaskReferenceId = "TPT00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			toolingProjectTaskInfo.TaskReferenceId = "TPT0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			toolingProjectTaskInfo.TaskReferenceId = "TPT" + strconv.Itoa(createdRecordId)
		}
		toolingProjectTaskInfo.Status = ProjectTaskToDO

		existingActionRemarks := make([]ActionRemarks, 0)
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "TASK ASSIGNED BY SUPERVISOR",
			UserId:        userId,
			Remarks:       "The task is successfully created and assigned",
			ProcessedTime: getTimeDifference(toolingProjectTaskInfo.CreatedAt),
		})
		toolingProjectTaskInfo.ActionRemarks = existingActionRemarks

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(toolingProjectTaskInfo)
		updatingData["object_info"] = rawWorkOrderInfo

		Update(dbConnection, ToolingProjectTaskTable, createdRecordId, updatingData)

	case ToolingProjectSprintComponent:
		_, projectTaskSprint := Get(dbConnection, ToolingProjectSprintTable, createdRecordId)

		toolingProjectSpring := ToolingProjectSprint{ObjectInfo: projectTaskSprint.ObjectInfo}
		toolingProjectSprintInfo := toolingProjectSpring.getToolingProjectSprintInfo()
		toolingProjectSprintInfo.ListOfAssignedTasks = make([]int, 0)
		toolingProjectSprintInfo.ListOfCompletedTasks = make([]int, 0)
		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(toolingProjectSprintInfo)
		updatingData["object_info"] = rawWorkOrderInfo

		Update(dbConnection, ToolingProjectSprintTable, createdRecordId, updatingData)
	}

	v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Resource is successfully created",
		Error:   0,
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
func (v *ToolingService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	err = ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
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
func (v *ToolingService) updateResource(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

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

	var updateRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updateRequest["lastUpdatedAt"] = util.GetCurrentTime(ISOTimeLayout)

	updatingData := make(map[string]interface{})

	//Adding update process request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the resource",
	})

}

func (v *ToolingService) checkInTask(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectTaskTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	toolingProjectTask := ToolingProjectTask{ObjectInfo: taskGeneralObject.ObjectInfo}
	toolingTaskInfo := toolingProjectTask.getToolingTaskInfo()

	if toolingTaskInfo.Status == ProjectTaskInProgress || toolingTaskInfo.Status == ProjectTaskDone {
		// already checked in , then why?
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Task Status",
				Description: "Task is already checked in, can not be proceed",
			})
		return
	}
	toolingTaskInfo.Status = ProjectTaskInProgress

	toolingTaskInfo.ActualStartDate = util.GetCurrentTime(ISOTimeLayout)
	toolingTaskInfo.CanCheckIn = false
	toolingTaskInfo.CanCheckOut = true

	checkInActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS CHECKED IN",
		UserId:        userId,
		Remarks:       "The task is successfully checked in",
		ProcessedTime: getTimeDifference(toolingTaskInfo.CreatedAt),
	}
	toolingTaskInfo.ActionRemarks = append(toolingTaskInfo.ActionRemarks, checkInActionRemarks)

	err = Update(dbConnection, ToolingProjectTaskTable, recordId, toolingTaskInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
		return
	}

	_, projectGeneralObject := Get(dbConnection, ToolingProjectTable, toolingTaskInfo.ProjectId)
	toolingProject := ToolingProject{ObjectInfo: projectGeneralObject.ObjectInfo}
	projectInfo := toolingProject.getToolingProjectInfo()

	checkInTaskActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS CHECKED IN",
		UserId:        userId,
		Remarks:       toolingTaskInfo.Name + " is successfully checked in",
		ProcessedTime: getTimeDifference(toolingTaskInfo.CreatedAt),
	}
	projectInfo.ActionRemarks = append(projectInfo.ActionRemarks, checkInTaskActionRemarks)

	if projectInfo.Status == ProjectCreated {
		projectInfo.ActualStartDate = util.GetCurrentTime(ISOTimeLayout)

	}

	_ = Update(dbConnection, ToolingProjectTable, toolingTaskInfo.ProjectId, projectInfo.DatabaseSerialize())
	v.CreateUserRecordMessage(ProjectID, ToolingTaskComponent, toolingTaskInfo.Name+" is checked in", projectGeneralObject.Id, userId, nil, nil)
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully checked in",
		Code:    0,
	})

	//TODO generate the email to project owner using ToolingProjectTaskStatusEmailTemplateType

}

func (v *ToolingService) approveTask(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectTaskTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	toolingProjectTask := ToolingProjectTask{ObjectInfo: taskGeneralObject.ObjectInfo}
	toolingTaskInfo := toolingProjectTask.getToolingTaskInfo()

	if toolingTaskInfo.Status != ProjectTaskDone {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Task Status",
				Description: "Wait until task get complete",
			})
		return
	}
	toolingTaskInfo.Status = ProjectTaskApproved

	toolingTaskInfo.CanCheckIn = false
	toolingTaskInfo.CanCheckOut = false

	toolingTaskInfo.CanApprove = false
	toolingTaskInfo.CanReject = false

	checkInActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS APPROVED",
		UserId:        userId,
		Remarks:       "Great, this task has been approved",
		ProcessedTime: getTimeDifference(toolingTaskInfo.CreatedAt),
	}
	toolingTaskInfo.ActionRemarks = append(toolingTaskInfo.ActionRemarks, checkInActionRemarks)

	err = Update(dbConnection, ToolingProjectTaskTable, recordId, toolingTaskInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating task information"), ErrorUpdatingObjectInformation)
		return
	}

	_, projectGeneralObject := Get(dbConnection, ToolingProjectTable, toolingTaskInfo.ProjectId)
	toolingProject := ToolingProject{ObjectInfo: projectGeneralObject.ObjectInfo}
	projectInfo := toolingProject.getToolingProjectInfo()

	checkInTaskActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS APPROVED",
		UserId:        userId,
		Remarks:       toolingTaskInfo.Name + " is approved as complete",
		ProcessedTime: getTimeDifference(toolingTaskInfo.CreatedAt),
	}
	projectInfo.ActionRemarks = append(projectInfo.ActionRemarks, checkInTaskActionRemarks)

	_ = Update(dbConnection, ToolingProjectTable, toolingTaskInfo.ProjectId, projectInfo.DatabaseSerialize())

	v.notifyTaskCompletion(dbConnection, toolingTaskInfo.AssignedUserId, taskGeneralObject.Id)
	v.CreateUserRecordMessage(ProjectID, ToolingTaskComponent, toolingTaskInfo.Name+" is approved", projectGeneralObject.Id, userId, nil, nil)
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully checked in",
		Code:    0,
	})

}

func (v *ToolingService) rejectTask(ctx *gin.Context) {

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	returnRemark := util.InterfaceToString(returnFields["remark"])

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectTaskTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	toolingProjectTask := ToolingProjectTask{ObjectInfo: taskGeneralObject.ObjectInfo}
	toolingTaskInfo := toolingProjectTask.getToolingTaskInfo()

	if toolingTaskInfo.Status != ProjectTaskDone {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Task Status",
				Description: "Wait until task get complete",
			})
		return
	}
	toolingTaskInfo.Status = ProjectTaskReDO

	toolingTaskInfo.CanCheckIn = true
	toolingTaskInfo.CanCheckOut = false

	toolingTaskInfo.CanApprove = false
	toolingTaskInfo.CanReject = false

	checkInActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS RE-ASSIGNED",
		UserId:        userId,
		Remarks:       returnRemark,
		ProcessedTime: getTimeDifference(toolingTaskInfo.CreatedAt),
	}
	toolingTaskInfo.ActionRemarks = append(toolingTaskInfo.ActionRemarks, checkInActionRemarks)

	err = Update(dbConnection, ToolingProjectTaskTable, recordId, toolingTaskInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating task information"), ErrorUpdatingObjectInformation)
		return
	}

	_, projectGeneralObject := Get(dbConnection, ToolingProjectTable, toolingTaskInfo.ProjectId)
	toolingProject := ToolingProject{ObjectInfo: projectGeneralObject.ObjectInfo}
	projectInfo := toolingProject.getToolingProjectInfo()

	checkInTaskActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS REJECTED",
		UserId:        userId,
		Remarks:       toolingTaskInfo.Name + "[" + returnRemark + "]",
		ProcessedTime: getTimeDifference(toolingTaskInfo.CreatedAt),
	}
	projectInfo.ActionRemarks = append(projectInfo.ActionRemarks, checkInTaskActionRemarks)

	_ = Update(dbConnection, ToolingProjectTable, toolingTaskInfo.ProjectId, projectInfo.DatabaseSerialize())

	v.notifyTaskCompletion(dbConnection, toolingTaskInfo.AssignedUserId, taskGeneralObject.Id)
	v.CreateUserRecordMessage(ProjectID, ToolingTaskComponent, toolingTaskInfo.Name+" is rejected", projectGeneralObject.Id, userId, nil, nil)

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Your action has been successfully applied",
		Code:    0,
	})

}

func (v *ToolingService) completeTask(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectTaskTable, recordId)

	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	toolingProjectTask := ToolingProjectTask{ObjectInfo: taskGeneralObject.ObjectInfo}
	taskInfo := toolingProjectTask.getToolingTaskInfo()
	taskInfo.Status = ProjectTaskDone

	taskInfo.ActualEndDate = util.GetCurrentTime(ISOTimeLayout)
	taskInfo.CanCheckIn = false
	taskInfo.CanCheckOut = false

	taskInfo.CanApprove = true
	taskInfo.CanReject = true

	checkOutActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS COMPLETED",
		UserId:        userId,
		Remarks:       "The task is successfully completed",
		ProcessedTime: getTimeDifference(taskInfo.CreatedAt),
	}
	taskInfo.ActionRemarks = append(taskInfo.ActionRemarks, checkOutActionRemarks)

	err = Update(dbConnection, ToolingProjectTaskTable, recordId, taskInfo.DatabaseSerialize(userId))

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
		return
	}

	_, projectGeneralObject := Get(dbConnection, ToolingProjectTable, taskInfo.ProjectId)
	toolingProject := ToolingProject{ObjectInfo: projectGeneralObject.ObjectInfo}
	projectInfo := toolingProject.getToolingProjectInfo()

	checkOutTaskActionRemarks := ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK IS COMPLETED",
		UserId:        userId,
		Remarks:       taskInfo.Name + " is successfully completed",
		ProcessedTime: getTimeDifference(taskInfo.CreatedAt),
	}
	projectInfo.ActionRemarks = append(projectInfo.ActionRemarks, checkOutTaskActionRemarks)

	if isAllTaskCompleted(projectGeneralObject.Id, dbConnection) {
		projectInfo.ActualEndDate = util.GetCurrentTime(ISOTimeLayout)

	}

	_ = Update(dbConnection, ToolingProjectTable, taskInfo.ProjectId, projectInfo.DatabaseSerialize())

	updateToolingSprint(recordId, toolingProjectTask, dbConnection)
	fmt.Println("Supervisor id:", projectInfo.Supervisor)
	v.notifyTaskCompletion(dbConnection, projectInfo.Supervisor, taskGeneralObject.Id)

	v.CreateUserRecordMessage(ProjectID, ToolingTaskComponent, taskInfo.Name+" is completed", projectGeneralObject.Id, userId, nil, nil)

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully completed",
		Code:    0,
	})

	//TODO generate the email to project owner using ToolingProjectTaskStatusEmailTemplateType

}

func (v *ToolingService) kickOffProject(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	userId := common.GetUserId(ctx)

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectTable, recordId)

	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}

	checkListGeneralObject, err := GetObjects(dbConnection, ToolingProjectTaskListTable)

	for _, checkListObj := range *checkListGeneralObject {
		checkListInfo := make(map[string]interface{})
		json.Unmarshal(checkListObj.ObjectInfo, &checkListInfo)

		if util.InterfaceToString(checkListInfo["objectStatus"]) != "Active" {
			continue
		}

		taskObject := ToolingTaskInfo{
			ProjectId:    recordId,
			Description:  util.InterfaceToString(checkListInfo["description"]),
			Name:         util.InterfaceToString(checkListInfo["name"]),
			TargetDate:   util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			Status:       1,
			CreatedAt:    util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:    userId,
			CanCheckIn:   true,
			CanCheckOut:  false,
			CanApprove:   false,
			CanReject:    false,
			ObjectStatus: "Active",
		}

		updatingData := component.GeneralObject{}
		updatingData.ObjectInfo, _ = json.Marshal(taskObject)

		err, _ = Create(dbConnection, ToolingProjectTaskTable, updatingData)

		if err != nil {
			v.BaseService.Logger.Error("error inserting check list into task", zap.String("error", err.Error()))
		}
	}

	toolingProject := make(map[string]interface{})
	json.Unmarshal(taskGeneralObject.ObjectInfo, &toolingProject)
	toolingProject["isKickOff"] = false
	toolingProject["status"] = ProjectCREATED

	updatingData := make(map[string]interface{})
	updatingData["object_info"], _ = json.Marshal(toolingProject)

	err = Update(dbConnection, ToolingProjectTable, recordId, updatingData)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating project information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully completed",
		Code:    0,
	})

	// send the email to all the task users about sprint is actived, and you need to check in regarding task
	projectTaskCondition := "JSON_UNQUOTE(JSON_EXTRACT(stats_info, \"$.projectId\")) = " + projectId
	listOfTasks, _ := GetConditionalObjects(dbConnection, ToolingProjectTaskTable, projectTaskCondition)
	for _, task := range *listOfTasks {
		toolingProjectTask := ToolingProjectTask{ObjectInfo: task.ObjectInfo}
		v.notifyProjectTaskAssignment(dbConnection, toolingProjectTask.getToolingTaskInfo().AssignedUserId, task.Id)
	}
	//TODO generate the email to project owner using ToolingProjectTaskStatusEmailTemplateType

}

func (v *ToolingService) activateSprint(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectSprintTable, recordId)

	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}

	toolingProjectSprint := ToolingProjectSprint{ObjectInfo: taskGeneralObject.ObjectInfo}
	sprintInfo := toolingProjectSprint.getToolingProjectSprintInfo()
	sprintInfo.Status = ProjectSprintInProgress

	if len(sprintInfo.ListOfAssignedTasks) == 0 || sprintInfo.ListOfAssignedTasks == nil {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "No Assigned Task Found",
				Description: "There is no assigned tesk in sprint. Please add an task first before do any operations",
			})
		return
	}

	sprintInfo.CanActivate = false
	sprintInfo.CanComplete = true
	err = Update(dbConnection, ToolingProjectSprintTable, recordId, sprintInfo.DatabaseSerialize())

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating Sprint information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully completed",
		Code:    0,
	})

	// send the email to all the task users about sprint is actived, and you need to check in regarding task
	projectTaskCondition := "JSON_UNQUOTE(JSON_EXTRACT(stats_info, \"$.projectId\")) = " + projectId
	listOfTasks, _ := GetConditionalObjects(dbConnection, ToolingProjectTaskTable, projectTaskCondition)
	for _, task := range *listOfTasks {
		toolingProjectTask := ToolingProjectTask{ObjectInfo: task.ObjectInfo}
		v.notifyProjectTaskAssignment(dbConnection, toolingProjectTask.getToolingTaskInfo().AssignedUserId, task.Id)
	}
	//TODO generate the email to project owner using ToolingProjectTaskStatusEmailTemplateType

}

func (v *ToolingService) completeSprint(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, ToolingProjectSprintTable, recordId)

	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}

	toolingProjectSprint := ToolingProjectSprint{ObjectInfo: taskGeneralObject.ObjectInfo}
	sprintInfo := toolingProjectSprint.getToolingProjectSprintInfo()
	sprintInfo.Status = ProjectSprintCompleted

	sprintInfo.CanActivate = false
	sprintInfo.CanComplete = false
	err = Update(dbConnection, ToolingProjectSprintTable, recordId, sprintInfo.DatabaseSerialize())

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully completed",
		Code:    0,
	})

	//TODO generate the email to project owner using ToolingProjectTaskStatusEmailTemplateType

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
func (v *ToolingService) getSearchResults(ctx *gin.Context) {

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
	listOfObjects, err := GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), ErrorGettingObjectsInformation)
			return

		}
	}
	// searchResponse := ts.ComponentManager.GetSearchResponse(componentName, listOfObjects)
	fmt.Println("zone:", zone)
	err, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)
}

func (v *ToolingService) getGroupBy(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")

	groupByAction := GroupByAction{}
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
		results := getGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range results {
			groupByChildren := GroupByChildren{}

			groupByChildren.Data = level1Value
			groupByChildren.Type = "json"

			tableGroupResponse := TableGroupByResponse{}
			tableGroupResponse.Label = level1Key
			tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Data = append(finalResponse.Data, rawData)
		}
	} else if len(groupByColumns) == 2 {
		level1Results := getGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range level1Results {
			tableGroupResponse := TableGroupByResponse{}
			if len(level1Value) > 1 {
				var internalGroupResponse []interface{}
				// here we need to group again

				level2Results := getGroupByResultsFromInterface(groupByColumns[1], level1Value)
				for level2Key, level2Value := range level2Results {
					level2Children := GroupByChildren{}
					level2Children.Data = level2Value
					level2Children.Type = "json"

					internalTableGroupResponse := TableGroupByResponse{}
					internalTableGroupResponse.Label = level2Key
					internalTableGroupResponse.Children = append(internalTableGroupResponse.Children, level2Children)
					internalGroupResponse = append(internalGroupResponse, internalTableGroupResponse)
				}
				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = internalGroupResponse
			} else {
				groupByChildren := GroupByChildren{}
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
func (v *ToolingService) getCardViewGroupBy(ctx *gin.Context) {

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

// updateResource ShowAccount godoc
// @Summary update given resource based on resource id
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags UserManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   resourceId     path    string     true        "Resource Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId}/remove_internal_array_reference [put]
func (v *ToolingService) removeInternalArrayReference(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	//projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	projectId := ctx.Param("projectId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

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
	updatingData["object_info"] = v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating assigned user information"), ErrorUpdatingObjectInformation)
		return
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully updated",
		Code:    0,
	})

}

func isAllTaskCompleted(toolingProjectId int, dbConnection *gorm.DB) bool {
	isCompleted := false
	conditionString := " JSON_EXTRACT(object_info, \"$.projectId\") "
	listOfObjects, err := GetConditionalObjects(dbConnection, ToolingProjectTable, conditionString)

	if err != nil {
		return isCompleted
	}

	totalNoTask := len(*listOfObjects)

	for _, object := range *listOfObjects {
		toolingProjectTask := ToolingProjectTask{ObjectInfo: object.ObjectInfo}
		taskInfo := toolingProjectTask.getToolingTaskInfo()

		if taskInfo.Status == ProjectTaskDone {
			totalNoTask -= 1
		}
	}

	if totalNoTask == 0 {
		isCompleted = true
	}

	return isCompleted
}
