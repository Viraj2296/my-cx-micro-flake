package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (as *AuthService) CreateUserRecordMessage(projectId string, componentName, message string, recordId int, userId int, existingData *component.GeneralObject, updatedData *component.GeneralObject) error {
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

	trackingFields := common.GetTrackingFields(existingData, updatedData)
	resourceInfo.TrackingFields = trackingFields
	rawResourceInfo, _ := json.Marshal(resourceInfo)
	recordTrail := AuthRecordTrail{
		RecordId:      recordId,
		ComponentName: componentName,
		ObjectInfo:    rawResourceInfo,
	}
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	err, recordId := CreateRecordTrail(dbConnection, recordTrail)
	as.BaseService.Logger.Info("creating a record trail id", zap.Any("trail_id", recordId))
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
func (as *AuthService) getComponentRecordTrails(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)

	componentName := ctx.Param("componentName")
	orderByCondition := "created_at desc"
	dbConnection := as.BaseService.ServiceDatabases[projectId]
	var conditionString = " component_name = \"" + componentName + "\" AND record_id =\"" + strconv.Itoa(recordId) + "\""
	listOfObjects, err := GetConditionalObjectsOrderBy(dbConnection, AuthRecordTrailTable, conditionString, orderByCondition)
	if err != nil {
		as.BaseService.Logger.Error("error loading resources", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting projects"), common.ObjectNotFound)
		return
	}
	userId := common.GetUserId(ctx)
	zone := as.GetUserTimezone(userId)
	response := common.GetRecordTrailResponse(zone, listOfObjects)
	ctx.JSON(http.StatusOK, response)

}
