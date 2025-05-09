package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (v *BatchManagementService) CreateUserRecordMessage(projectId string, componentName, message string, recordId int, userId int, existingData *component.GeneralObject, updatedData *component.GeneralObject) error {
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
	recordTrail := database.BatchManagementRecordTrail{
		RecordId:      recordId,
		ComponentName: componentName,
		ObjectInfo:    rawResourceInfo,
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, recordId := database.CreateRecordTrail(dbConnection, recordTrail)
	v.BaseService.Logger.Info("creating a record trail id", zap.Any("trail_id", recordId))
	return err
}

func (v *BatchManagementService) getComponentRecordTrails(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)

	componentName := ctx.Param("componentName")
	orderByCondition := "created_at desc"
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionString = " component_name = '" + componentName + "' AND record_id ='" + strconv.Itoa(recordId) + "'"
	err, listOfObjects := database.GetConditionalObjectsOrderBy(dbConnection, const_util.BatchManagementRecordTrailTable, conditionString, orderByCondition)
	if err != nil {
		v.BaseService.Logger.Error("error loading resources", zap.Error(err))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting projects"), common.ObjectNotFound)
		return
	}
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	trailResponse := common.GetRecordTrailResponseV1(zone, listOfObjects)
	ctx.JSON(http.StatusOK, trailResponse)

}
