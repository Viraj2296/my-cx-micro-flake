package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"net/http"
	"strconv"
	"strings"
)

func (v *NotificationService) getObjects(ctx *gin.Context) {
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
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
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
		searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
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
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (v *NotificationService) getNotificationList(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	notificationListResponse := NotificationListResponse{}
	userId := common.GetUserId(ctx)

	var notificationList []datatypes.JSON
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	userSpecificNotificationIds := authService.GetNotificationList(userId)
	userSpecificViewNotificationIds := authService.GetViewNotificationList(userId)
	//Get notification limit
	userInfo := authService.GetUserInfoById(userId)
	notificationLimit := 20
	if userInfo.NotificationLimit != 0 {
		notificationLimit = userInfo.NotificationLimit
	}
	conditionString := "JSON_CONTAINS(object_info->'$.targetUsers', CAST(" + strconv.Itoa(userId) + " AS JSON), '$')"
	orderByCondition := "id desc"
	listOfObjects, err := GetConditionalObjectsOrderBy(dbConnection, targetTable, conditionString, orderByCondition, notificationLimit)
	var notViewCount int
	var isAnyNotificationAdded bool
	if err != nil {
		notificationListResponse.ViewCount = 0
		var emptyArray []datatypes.JSON
		notificationListResponse.Notifications = emptyArray
	} else {

		for _, systemNotificationObject := range *listOfObjects {
			var notificationObjectFields = make(map[string]interface{})
			json.Unmarshal(systemNotificationObject.ObjectInfo, &notificationObjectFields)
			if IsNotificationIdExist(userSpecificNotificationIds, systemNotificationObject.Id) {
				if IsNotificationIdExist(userSpecificViewNotificationIds, systemNotificationObject.Id) {
					notificationObjectFields["isView"] = true
				} else {
					notificationObjectFields["isView"] = false
					notViewCount = notViewCount + 1
				}
				notificationObjectFields["id"] = systemNotificationObject.Id
				rawNotificationObject, _ := json.Marshal(notificationObjectFields)
				notificationList = append(notificationList, rawNotificationObject)
				notificationListResponse.Notifications = notificationList
				isAnyNotificationAdded = true
			}

		}
	}
	// if the list is empty, then send the empty array
	if !isAnyNotificationAdded {
		var emptyArray = make([]datatypes.JSON, 0)
		notificationListResponse.Notifications = emptyArray
	}
	notificationListResponse.ViewCount = notViewCount
	ctx.JSON(http.StatusOK, notificationListResponse)
}

func (v *NotificationService) updateResource(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	switch componentName {
	case SystemNotificationComponent:
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userId := common.GetUserId(ctx)
		v.BaseService.Logger.Info("adding view notification id to user", zap.Any("user_id", userId), zap.Any("notification id", intRecordId))
		err := authService.AddViewNotificationIds(userId, intRecordId)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("unable to update notification seen status"), ErrorGettingObjectInformation)
			return
		}

	default:
		err, objectInterface := Get(dbConnection, targetTable, intRecordId)
		if err != nil {
			v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingObjectInformation)
			return
		}

		var updateRequest = make(map[string]interface{})

		if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		updatingData := make(map[string]interface{})

		updatingData["object_info"] = v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
		err = Update(dbConnection, targetTable, intRecordId, updatingData)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
			return
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated",
	})

}

func (v *NotificationService) handleComponentAction(ctx *gin.Context) {
	actionName := ctx.Param("actionName")
	if actionName == "send_test_email" {
		v.sendTestEmail(ctx)
	} else if actionName == ActionMarkAllRead {

		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userId := common.GetUserId(ctx)
		existingNotificationList := authService.GetNotificationList(userId)

		//Loop through the not notify ids
		for _, existingNotificationId := range existingNotificationList {
			err := authService.AddViewNotificationIds(userId, existingNotificationId.Id)
			if err != nil {
				v.BaseService.Logger.Error("error in Adding notification id through auth service", zap.String("error", err.Error()))
				continue
			}
		}
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Successfully updated",
		})
	}

}
