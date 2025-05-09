package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (v *FacilityService) CreateUserRecordMessage(projectId string, componentName, message string, recordId int, userId int, existingData *component.GeneralObject, updatedData *component.GeneralObject) error {
	//TODO change this to int64  as we are using recordId (int) datatype, not a string
	resourceInfo := common.ResourceInfo{}
	resourceMeta := common.ResourceMeta{}

	resourceMeta.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	resourceMeta.UserId = userId
	resourceMeta.UpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	resourceInfo.ResourceMeta = resourceMeta
	resourceInfo.Message = message
	resourceInfo.MessageType = common.MessageTypeNotification
	resourceInfo.HasAttachment = false
	resourceInfo.AttachmentList = make([]string, 0)

	existingDataInfo := make(map[string]interface{})

	if existingData != nil {
		json.Unmarshal(existingData.ObjectInfo, &existingDataInfo)
	}

	currentTime := time.Now()

	trackingFields := common.GetTrackingFields(existingData, updatedData)
	resourceInfo.TrackingFields = trackingFields
	resourceInfo.SourceObject.Version = "version_" + currentTime.Format("2006.01.02 15:04:05")
	resourceInfo.SourceObject.ObjectInfo = existingDataInfo
	rawResourceInfo, _ := json.Marshal(resourceInfo)
	recordTrail := database.FacilityServiceRecordTrail{
		RecordId:      recordId,
		ComponentName: componentName,
		ObjectInfo:    rawResourceInfo,
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, recordId := database.CreateRecordTrail(dbConnection, recordTrail)
	v.BaseService.Logger.Info("creating a record trail id", zap.Any("trail_id", recordId))
	return err
}

// getComponentResourceRecords ShowAccount godoc
// @Summary Get all the resource related to given resource id for particular component id
// @Description based on user permission, user will get the list of records assigned and particular record
// @Tags RecordMessages
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /project/{projectId}/record_messages/component/{componentName}/record_messages/:recordId [get]
func (v *FacilityService) getComponentRecordTrails(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)

	componentName := ctx.Param("componentName")
	orderByCondition := "created_at desc"
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionString = " component_name = \"" + componentName + "\" AND record_id =\"" + strconv.Itoa(recordId) + "\""
	listOfObjects, err := database.GetConditionalObjectsOrderBy(dbConnection, const_util.FacilityServiceRecordTrailTable, conditionString, orderByCondition)
	if err != nil {
		v.BaseService.Logger.Error("error loading resources", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting projects"), common.ObjectNotFound)
		return
	}
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := common.GetRecordTrailResponse(zone, listOfObjects)
	ctx.JSON(http.StatusOK, response)

}

func (v *FacilityService) getSourceObjectVersions(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)

	componentName := ctx.Param("componentName")
	orderByCondition := "created_at desc"
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionString = " component_name = \"" + componentName + "\" AND record_id =\"" + strconv.Itoa(recordId) + "\"" + " AND object_info ->> '$.sourceObject' != 'null'"

	listOfObjects, err := database.GetConditionalObjectsOrderBy(dbConnection, const_util.FacilityServiceRecordTrailTable, conditionString, orderByCondition)
	if err != nil {
		v.BaseService.Logger.Error("error loading resources", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting projects"), common.ObjectNotFound)
		return
	}

	if len(*listOfObjects) == 0 {
		response := map[string]interface{}{"data": make([]interface{}, 0), "canRollBack": false}
		ctx.JSON(http.StatusOK, response)
		return
	}

	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := common.GetRecordVersionResponse(listOfObjects, zone)
	responseObject := make(map[string]interface{})

	responseJson, _ := json.Marshal(response)
	json.Unmarshal(responseJson, &responseObject)
	responseObject["canRollBack"] = true
	ctx.JSON(http.StatusOK, responseObject)
}

func (v *FacilityService) handleRollBack(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	var updateRequest = make(map[string]interface{})

	resourceInfo := common.ResourceInfo{}
	updatingData := make(map[string]interface{})
	err, objectInterface := database.Get(dbConnection, targetTable, recordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	versionName := util.InterfaceToString(updateRequest["version"])

	var conditionString = " component_name = \"" + componentName + "\" AND record_id =\"" + strconv.Itoa(recordId) + "\"" + " AND object_info ->> '$.sourceObject.version' =" + versionName

	listOfObjects, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRecordTrailTable, conditionString)
	if err != nil || len(*listOfObjects) == 0 {
		v.BaseService.Logger.Error("error loading resources", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting projects"), common.ObjectNotFound)
		return
	}

	json.Unmarshal((*listOfObjects)[0].ObjectInfo, &resourceInfo)
	serializedObject := v.ComponentManager.GetUpdateRequest(resourceInfo.SourceObject.ObjectInfo, objectInterface.ObjectInfo, componentName)

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Validation Failed"), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error roll back information"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	userId := common.GetUserId(ctx)
	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully updated",
		Error:   0,
	})
}
