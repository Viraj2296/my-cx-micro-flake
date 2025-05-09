package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (v *ProductionOrderService) loadFile(ctx *gin.Context) {

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

func (v *ProductionOrderService) importObjects(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ReferenceDatabase
	err, errorCode, importObjects := v.ComponentManager.ImportData(dbConnection, componentName, importDataCommand)
	if err != nil {
		v.BaseService.Logger.Error("unable to import data", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}

	var failedRecords int
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	for _, object := range importObjects.InsertObjects {
		err, _ = Create(dbConnection, targetTable, object)

		if err != nil {
			v.BaseService.Logger.Error("unable to create record", zap.String("error", err.Error()))
			failedRecords = failedRecords + 1
		}
		//recordIdInString := strconv.Itoa(recordId)
		//CreateBotRecordTrail(projectId, recordIdInString, componentName, "machine master is created")
	}
	for _, object := range importObjects.UpdateObjects {
		updatingData := make(map[string]interface{})
		rawObject, _ := json.Marshal(object.ObjectInfo)
		updatingData["object_info"] = rawObject
		err = Update(dbConnection, targetTable, object.Id, updatingData)

		if err != nil {
			v.BaseService.Logger.Error("unable to create record", zap.String("error", err.Error()))
			failedRecords = failedRecords + 1
		}
		//recordIdInString := strconv.Itoa(recordId)
		//CreateBotRecordTrail(projectId, recordIdInString, componentName, "machine master is created")
	}
	var message string
	if importObjects.TotalSkippedRecords == importObjects.TotalRecords {
		// all are skipped
		message = "Data import is failed " + importObjects.SkippedRecordNames + " Records were skipped due to schema validation, please check the schema configured to import under this module. Normally records are skipped due to validation failures from master data"
	} else if importObjects.TotalSkippedRecords > 0 {
		message = "Data is successfully imported" + importObjects.SkippedRecordNames + " but some of the records were skipped due to schema validation, please check the schema configured to import under this module. Normally records are skipped due to validation failures from master data"
	} else {
		message = "Great!, all the data is successfully imported based on schema configured"
	}
	importDataResponse := component.ImportDataResponse{
		TotalRecords:        importObjects.TotalRecords,
		FailedRecords:       failedRecords,
		SkippedData:         importObjects.SkippedData,
		TotalSkippedRecords: importObjects.TotalSkippedRecords,
		Message:             message,
	}

	ctx.JSON(http.StatusOK, importDataResponse)
}

func (v *ProductionOrderService) exportObjects(ctx *gin.Context) {
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
	if len(exportCommand.Attributes) > 0 {
		// user has requested with attributes
		var conditions []string
		for key, attr := range exportCommand.Attributes {
			operator := attr.Operator
			valueType := attr.Type
			value := attr.Value

			// Handle different types
			switch valueType {
			case "integer":
				// Assuming the value is safe to be used as an integer
				condition := fmt.Sprintf("object_info->'$.%s' %s %d", key, operator, util.InterfaceToInt(value))
				conditions = append(conditions, condition)

			case "date_time":
				// Wrap date/time values in single quotes for SQL
				var convertedValue = util.InterfaceToString(value)
				if len(convertedValue) > 0 {
					condition := fmt.Sprintf("object_info->'$.%s' %s '%s'", key, operator, value)
					conditions = append(conditions, condition)
				}

			default:
				// For any other types, treat as string (adjust as needed)
				condition := fmt.Sprintf("object_info->'$.%s' %s '%s'", key, operator, util.InterfaceToString(value))
				conditions = append(conditions, condition)
			}
		}

		// Join the conditions with AND
		condition = strings.Join(conditions, " AND ")
	}
	v.BaseService.Logger.Info("passing the condition to get the objects", zap.String("condition", condition))
	if componentName == ScheduledOrderEventComponent {
		err, errorCode, exportDataResponse := v.ComponentManager.ExportDataProductionOrder(dbConnection, componentName, exportCommand, condition)
		if err != nil {
			v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
			return
		}
		ctx.JSON(http.StatusOK, exportDataResponse)
	} else {
		err, errorCode, exportDataResponse := v.ComponentManager.GeneralExportData(dbConnection, componentName, exportCommand, condition)
		if err != nil {
			v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
			return
		}
		ctx.JSON(http.StatusOK, exportDataResponse)
	}

}

func (v *ProductionOrderService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := v.ComponentManager.GetTableImportSchema(componentName)
	ctx.JSON(http.StatusOK, tableImportSchema.Fields)
}

func (v *ProductionOrderService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)
}

func (v *ProductionOrderService) getObjects(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	orderValue := ctx.Query("order")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
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

	now := time.Now()
	oneMonthBefore := now.AddDate(0, -1, 0)
	dateTimeStr := oneMonthBefore.Format(TimeLayout)
	timeCondition := " object_info ->> '$.createdAt' > '" + dateTimeStr + "'"

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
		if componentName == ProductionOrderMasterTable || componentName == AssemblyProductionOrderTable || componentName == ToolingOrderMasterTable {
			searchWithBaseQuery = searchWithBaseQuery + " AND " + timeCondition
		}
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			queryCondition := component.TableCondition(offsetValue, fields, values, condition)
			//if componentName == ProductionOrderMasterTable || componentName == AssemblyProductionOrderTable || componentName == ToolingOrderMasterTable {
			//	queryCondition = queryCondition + " AND " + timeCondition
			//}
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, queryCondition)

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			queryCondition := component.TableCondition(offsetValue, fields, values, condition)
			if orderValue == "desc" {
				//queryCondition = component.TableCondition(offsetValue, fields, values, condition)
				offsetVal, _ := strconv.Atoi(offsetValue)
				if offsetVal == -1 {
					queryCondition = component.TableConditionV1(offsetValue, fields, values, condition)
				} else {
					queryCondition = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
				}

				orderBy := "object_info ->> '$.createdAt' desc"
				listOfObjects, err = GetConditionalObjectsOrderBy(dbConnection, targetTable, queryCondition, orderBy, limitVal)

				// when we do the decending mode, if we reach to 1,then set the is next false
				for _, object := range *listOfObjects {
					if object.Id == 1 {
						isNext = false
					}
				}

			} else {
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

func (v *ProductionOrderService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	orderValue := ctx.Query("order")

	fields := ctx.Query("fields")
	values := ctx.Query("values")
	condition := ctx.Query("condition")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	if offsetValue != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)
		totalRecords := Count(dbConnection, targetTable)
		queryCondition := component.TableCondition(offsetValue, fields, values, condition)
		if orderValue == "desc" {
			offsetVal, _ := strconv.Atoi(offsetValue)
			offsetValue = strconv.Itoa(int(totalRecords) - limitVal + 1)

			limitVal = limitVal - offsetVal
			queryCondition = component.TableCondition(offsetValue, fields, values, condition)
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, queryCondition, limitVal)
			listOfObjects = reverseSlice(listOfObjects)
		} else {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, queryCondition, limitVal)
		}

		// requesting to search fields for table
		//listOfObjects, err = GetObjects(dbConnection, targetTable, limitVal)
	} else {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
	}

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

}

func (v *ProductionOrderService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	switch componentName {
	case ScheduledOrderEventComponent:
		// if this is deleted, we should add back the remianingschedueld quantiy , and delete the event from table
		productionOrderInfo := ProductionOrderInfo{}
		scheduledEvent := ScheduledOrderEvent{ObjectInfo: generalObject.ObjectInfo}

		orderStatusIdPreferenceFour := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFour)
		scheduledEventInfo := scheduledEvent.getScheduledOrderEventInfo()
		if scheduledEvent.getScheduledOrderEventInfo().EventStatus == orderStatusIdPreferenceFour {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidScheduleEventStatusError), InvalidFieldValue, ScheduleIsAlreadyReleaseDescription)
			return
		}

		err, productionOrderObject := Get(dbConnection, ProductionOrderMasterTable, scheduledEventInfo.EventSourceId)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(ObjectNotFound), ErrorGettingIndividualObjectInformation, InvalidProductionOrderDescription)
			return
		}
		// now delete the event
		err = ArchiveObject(dbConnection, ScheduledOrderEventTable, generalObject)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(DeleteEventFailed), ErrorGettingIndividualObjectInformation, DeleteScheduledEventFailedDescription)
			return
		}

		json.Unmarshal(productionOrderObject.ObjectInfo, &productionOrderInfo)
		createdOrderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceTwo)

		scheduledQty := scheduledEventInfo.ScheduledQty
		fmt.Println("scheduledQty: ", scheduledQty)
		productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty + scheduledQty
		productionOrderInfo.OrderStatus = createdOrderStatusId
		updatingData := make(map[string]interface{})
		rawProductionOrderInfo, _ := json.Marshal(productionOrderInfo)
		updatingData["object_info"] = rawProductionOrderInfo
		err = Update(dbConnection, ProductionOrderMasterTable, scheduledEventInfo.EventSourceId, updatingData)
		if err != nil {
			v.BaseService.Logger.Error("update production order master table had failed", zap.String("error", err.Error()))
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Production Order Information Failed"), ErrorUpdatingObjectInformation, "Updating production master information has failed due to ["+err.Error()+"]")
			return
		}

		// create a system notification

		systemNotification := common.SystemNotification{}
		eventName := scheduledEventInfo.Name
		systemNotification.Name = eventName + " is deleted"
		systemNotification.ColorCode = "#14F44E"
		systemNotification.IconCls = "icon-park-outline:transaction-order"
		systemNotification.RecordId = intRecordId
		systemNotification.RouteLinkComponent = "machine_timeline_event"
		systemNotification.Component = "Production Order"
		systemNotification.Description = "Order is scheduled with id [" + productionOrderInfo.ProdOrder + "] is deleted"
		systemNotification.GeneratedTime = util.GetCurrentTime(time.RFC822)
		rawSystemNotification, _ := json.Marshal(systemNotification)
		notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
		notificationService.CreateSystemNotification(projectId, rawSystemNotification)
	case ProductionOrderMasterComponent:
		err = ArchiveObject(dbConnection, targetTable, generalObject)

		if err != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		}
		//Find production order name
		var productionOrderMaster ProductionOrderInfo
		err = json.Unmarshal([]byte(generalObject.ObjectInfo), &productionOrderMaster)

		if err != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		}
		//Fetch machine time lines based on prodOrder
		timelineQuery := "select * from scheduled_order_event where object_info->>'$.eventSourceId' = '" + recordId + "'"
		var commonObject []component.GeneralObject
		dbConnection.Raw(timelineQuery).Scan(&commonObject)

		for _, timeLineObject := range commonObject {
			//Archived machine time line
			err = ArchiveObject(dbConnection, ScheduledOrderEventTable, timeLineObject)
			if err != nil {
				v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
				continue
			}

		}
	case ToolingScheduledOrderEventComponent:
		toolingOrderInfo := make(map[string]interface{})
		scheduledEvent := ToolingScheduledOrderEvent{ObjectInfo: generalObject.ObjectInfo}

		orderStatusIdPreferenceFour := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFour)
		scheduledEventInfo := scheduledEvent.getToolingScheduledOrderEventInfo()
		if scheduledEvent.getToolingScheduledOrderEventInfo().EventStatus == orderStatusIdPreferenceFour {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidScheduleEventStatusError), InvalidFieldValue, ScheduleIsAlreadyReleaseDescription)
			return
		}

		err, productionOrderObject := Get(dbConnection, ToolingOrderMasterTable, scheduledEventInfo.EventSourceId)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(ObjectNotFound), ErrorGettingIndividualObjectInformation, InvalidProductionOrderDescription)
			return
		}
		// now delete the event
		err = ArchiveObject(dbConnection, ToolingScheduledOrderEventTable, generalObject)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(DeleteEventFailed), ErrorGettingIndividualObjectInformation, DeleteScheduledEventFailedDescription)
			return
		}

		json.Unmarshal(productionOrderObject.ObjectInfo, &toolingOrderInfo)

		orderRemainingDuration := util.InterfaceToString(toolingOrderInfo["remainingDuration"])
		orderRemainingDurationSec, _ := time.ParseDuration(orderRemainingDuration)

		_, partMasterObject := Get(dbConnection, ToolingPartMasterTable, scheduledEvent.getToolingScheduledOrderEventInfo().PartId)
		partMasterInfo := make(map[string]interface{})
		json.Unmarshal(partMasterObject.ObjectInfo, &partMasterInfo)

		day := util.InterfaceToInt(partMasterInfo["day"])
		hour := util.InterfaceToInt(partMasterInfo["hour"])
		minute := util.InterfaceToInt(partMasterInfo["minute"])

		remainingDuration := day*86400 + hour*3600 + minute*60

		currentRemainingDuration := int(orderRemainingDurationSec.Seconds()) - remainingDuration

		s := fmt.Sprintf("%.2f", currentRemainingDuration) + "s"

		h, _ := time.ParseDuration(s)

		toolingOrderInfo["remainingDuration"] = fmt.Sprintf("%.2f", h.Hours()) + "h"

		Update(dbConnection, ToolingOrderMasterTable, productionOrderObject.Id, toolingOrderInfo)
	case ToolingOderMasterComponent:
		//Archieved Paroduction order master
		err = ArchiveObject(dbConnection, targetTable, generalObject)

		if err != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		}

		//Fetch machine time lines based on prodOrder
		timelineQuery := "select * from tooling_scheduled_order_event where object_info->>'$.eventSourceId' = '" + recordId + "'"
		var commonObject []component.GeneralObject
		dbConnection.Raw(timelineQuery).Scan(&commonObject)

		for _, timeLineObject := range commonObject {
			//Archived machine time line
			err = ArchiveObject(dbConnection, ToolingScheduledOrderEventTable, timeLineObject)
			if err != nil {
				v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
				continue
			}

		}
	default:
		err = ArchiveObject(dbConnection, targetTable, generalObject)
		//err = Delete(dbConnection, targetTable, generalObject)
	}

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	ctx.Status(http.StatusNoContent)

}

func (v *ProductionOrderService) updateResource(ctx *gin.Context) {

	componentName := util.GetComponentName(ctx)
	projectId := util.GetProjectId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	switch componentName {
	case ScheduledOrderEventComponent:
		v.handleUpdateSchedule(ctx)
		return
	case AssemblyScheduledOrderEventComponent:
		v.handleUpdateAssemblySchedule(ctx)
		return
	case ToolingScheduledOrderEventComponent:
		v.handleUpdateToolingSchedule(ctx)
		return
	default:
		if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		updatingData := make(map[string]interface{})
		err, objectInterface := Get(dbConnection, targetTable, recordId)
		if err != nil {
			v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
			return
		}

		if componentName == ToolingPartMasterComponent {
			if !isValidDuration(updateRequest) {
				response.DispatchDetailedError(ctx, FieldValidationFailed,
					&response.DetailedError{
						Header:      getError(InvalidScheduleEventStatusError).Error(),
						Description: "Duration can not be set as 0, you can set the duration as expected time needed to complete.",
					})
				return
			}
		}

		if !common.ValidateObjectStatus(objectInterface.ObjectInfo) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      getError(common.InvalidObjectStatusError).Error(),
					Description: "This resource is already archived, no further modifications are allowed.",
				})
			return
		}
		orderStatusPreference1Id := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceOne)
		//Can't modify balance after production order reached scheduled status
		var productionOrderNotification string
		if componentName == ProductionOrderMasterComponent {

			conditionString := " object_info ->>'$.eventSourceId' = " + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)
			if err == nil {
				var numberOfNonScheduledEvent int
				numberOfNonScheduledEvent = 0
				for _, schedulerEvent := range *listOfObjects {
					commonObject := ScheduledOrderEvent{ObjectInfo: schedulerEvent.ObjectInfo}
					if commonObject.getScheduledOrderEventInfo().EventStatus != ScheduleStatusPreferenceThree {
						numberOfNonScheduledEvent = numberOfNonScheduledEvent + 1
					}
				}

				if numberOfNonScheduledEvent > 0 {
					response.DispatchDetailedError(ctx, common.OperationNotPermitted,
						&response.DetailedError{
							Header:      getError(common.OperationNotPermittedError).Error(),
							Description: "One or more scheduled orders are already been processed, modifying order is not possible at this moment",
						})
					return
				}
			}

			var productionInfo = make(map[string]interface{})
			json.Unmarshal(objectInterface.ObjectInfo, &productionInfo)
			productionOrderNotification = updateRequest["prodOrder"].(string) + " is updated"
			if _, ok := updateRequest["balance"]; ok {
				if util.InterfaceToInt(productionInfo["balance"]) != util.InterfaceToInt(updateRequest["balance"]) {
					//Check order status of production order
					if util.InterfaceToInt(productionInfo["orderStatus"]) != orderStatusPreference1Id {
						response.DispatchDetailedError(ctx, common.OperationNotPermitted,
							&response.DetailedError{
								Header:      getError(common.OperationNotPermittedError).Error(),
								Description: "This resource is already scheduled or partially scheduled, no further modification in balance.",
							})
						return
					}
				}
			}

			// now check if the machine is getting updated,
			if util.InterfaceToInt(productionInfo["machineId"]) != util.InterfaceToInt(updateRequest["machineId"]) {
				if util.InterfaceToInt(productionInfo["orderStatus"]) > 3 {
					response.DispatchDetailedError(ctx, common.OperationNotPermitted,
						&response.DetailedError{
							Header:      getError(common.OperationNotPermittedError).Error(),
							Description: "The order is already moved in difference stage where changing machine is not possible",
						})
					return
				}

				// now update the machine to all the child events
				scheduledEventCondition := " object_info ->>'$.eventSourceId' = " + strconv.Itoa(recordId)
				listOfScheduledEvents, _ := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, scheduledEventCondition)
				for _, schedulerEvent := range *listOfScheduledEvents {
					var objectFields = make(map[string]interface{})
					json.Unmarshal(schedulerEvent.ObjectInfo, &objectFields)
					objectFields["machineId"] = util.InterfaceToInt(updateRequest["machineId"])
					serialisedObject, _ := json.Marshal(objectFields)
					scheduledUpdatingData := make(map[string]interface{})
					scheduledUpdatingData["object_info"] = serialisedObject
					err = Update(v.BaseService.ReferenceDatabase, ScheduledOrderEventTable, schedulerEvent.Id, scheduledUpdatingData)
				}
			}
		} else if componentName == ToolingOrderMasterTable {
			conditionString := " object_info ->>'$.eventSourceId' = " + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionString)
			if err == nil {
				var numberOfNonScheduledEvent int
				numberOfNonScheduledEvent = 0
				for _, schedulerEvent := range *listOfObjects {
					commonObject := make(map[string]interface{})
					json.Unmarshal(schedulerEvent.ObjectInfo, &commonObject)
					if util.InterfaceToInt(commonObject["eventStatus"]) != ScheduleStatusPreferenceThree {
						numberOfNonScheduledEvent = numberOfNonScheduledEvent + 1
					}
				}

				if numberOfNonScheduledEvent > 0 {
					response.DispatchDetailedError(ctx, common.OperationNotPermitted,
						&response.DetailedError{
							Header:      getError(common.OperationNotPermittedError).Error(),
							Description: "One or more scheduled orders are already been processed, modifying order is not possible at this moment",
						})
					return
				}
			}

			var toolingOrderInfo = make(map[string]interface{})
			json.Unmarshal(objectInterface.ObjectInfo, &toolingOrderInfo)
			// now check if the machine is getting updated,
			if util.InterfaceToInt(toolingOrderInfo["machineId"]) != util.InterfaceToInt(updateRequest["machineId"]) {
				if util.InterfaceToInt(toolingOrderInfo["orderStatus"]) > 3 {
					response.DispatchDetailedError(ctx, common.OperationNotPermitted,
						&response.DetailedError{
							Header:      getError(common.OperationNotPermittedError).Error(),
							Description: "The order is already moved in difference stage where changing machine is not possible",
						})
					return
				}

				// now update the machine to all the child events
				scheduledEventCondition := " object_info ->>'$.eventSourceId' = " + strconv.Itoa(recordId)
				listOfScheduledEvents, _ := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, scheduledEventCondition)
				for _, schedulerEvent := range *listOfScheduledEvents {
					var objectFields = make(map[string]interface{})
					json.Unmarshal(schedulerEvent.ObjectInfo, &objectFields)
					objectFields["machineId"] = util.InterfaceToInt(updateRequest["machineId"])
					serialisedObject, _ := json.Marshal(objectFields)
					scheduledUpdatingData := make(map[string]interface{})
					scheduledUpdatingData["object_info"] = serialisedObject
					err = Update(v.BaseService.ReferenceDatabase, ToolingScheduledOrderEventTable, schedulerEvent.Id, scheduledUpdatingData)
				}
			}
		}

		serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
		initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
		updatingData["object_info"] = initializedObject

		err = Update(v.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)

		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
			return
		}

		//Adding system notification for update production order
		if componentName == ProductionOrderMasterComponent {
			notificationHeader := productionOrderNotification
			notificationDescription := "The production order is modified, fields are updated. Check the production order for further details"
			v.createSystemNotification(projectId, notificationHeader, notificationDescription, recordId)
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated",
	})

}

func (v *ProductionOrderService) createNewResource(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	userId := common.GetUserId(ctx)
	switch componentName {
	case ScheduledOrderEventComponent:
		v.handleCreateSchedule(ctx)
		return
	case AssemblyScheduledOrderEventComponent:
		v.handleCreateAssemblySchedule(ctx)
		return
	case ToolingScheduledOrderEventComponent:
		v.handleCreateToolingSchedule(ctx)
		return
	case ToolingOderMasterComponent:
		orderRequest := make(map[string]interface{})
		if err := ctx.ShouldBindBodyWith(&orderRequest, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		listOfPartNo := util.InterfaceToIntArray(orderRequest["partNo"])
		var inQuery string
		inQuery = ""
		if len(listOfPartNo) == 0 {
			inQuery = "-1"
		} else {
			for index, partNo := range listOfPartNo {
				if index == len(listOfPartNo)-1 {
					inQuery += strconv.Itoa(partNo)
				} else {
					inQuery += strconv.Itoa(partNo) + ","
				}

			}
		}

		err := v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, orderRequest)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
			return
		}
		orderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceOne)
		orderRequest["orderStatus"] = orderStatusId

		rawCreateRequest, _ := json.Marshal(orderRequest)
		preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
		object := component.GeneralObject{
			ObjectInfo: preprocessedRequest,
		}
		err, createdRecordId := Create(dbConnection, targetTable, object)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
			return
		}

		v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
		ctx.JSON(http.StatusOK, component.GeneralResponse{
			RecordId: createdRecordId,
			Message:  "Successfully created the resource",
			Error:    0,
		})

		return
	case ProductionOrderMasterComponent:
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

		if _, ok := createRequest["balance"]; ok {
			if util.InterfaceToInt(createRequest["balance"]) > util.InterfaceToInt(createRequest["orderQty"]) {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
				return
			}

			if util.InterfaceToInt(createRequest["balance"]) > 0 {
				// front-end sending the balance, so means , this is how much we left, so that is remaining scheduled quantity
				createRequest["remainingScheduledQty"] = util.InterfaceToInt(createRequest["balance"])
			}
			if util.InterfaceToInt(createRequest["orderQty"]) < util.InterfaceToInt(createRequest["dailyRate"]) {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, "The daily rate can not be greater than order quantity")
				return
			}
		}

		rawCreateRequest, _ := json.Marshal(createRequest)
		preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
		object := component.GeneralObject{
			ObjectInfo: preprocessedRequest,
		}
		err, createdRecordId := Create(dbConnection, targetTable, object)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
			return
		}

		notificationHeader := createRequest["prodOrder"].(string)
		notificationDescription := "New production order is scheduled with id [" + createRequest["prodOrder"].(string) + "] click view to see more details"
		v.createSystemNotification(projectId, notificationHeader, notificationDescription, createdRecordId)
		v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
		ctx.JSON(http.StatusOK, component.GeneralResponse{
			RecordId: createdRecordId,
			Message:  "Successfully created the resource",
			Error:    0,
		})

	default:
		var createRequest = make(map[string]interface{})
		if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if componentName == ToolingPartMasterTable {
			if !isValidDuration(createRequest) {
				response.DispatchDetailedError(ctx, FieldValidationFailed,
					&response.DetailedError{
						Header:      getError(InvalidScheduleEventStatusError).Error(),
						Description: "Duration can not be set as 0, you can set the duration as expected time needed to complete.",
					})
				return
			}
		}

		// here we should do the validation
		err := v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
			return
		}

		rawCreateRequest, _ := json.Marshal(createRequest)
		preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
		object := component.GeneralObject{
			ObjectInfo: preprocessedRequest,
		}
		err, createdRecordId := Create(dbConnection, targetTable, object)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
			return
		}

		v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
		ctx.JSON(http.StatusOK, component.GeneralResponse{
			Message:  "Successfully created the resource",
			RecordId: createdRecordId,
			Error:    0,
		})
	}
}

func (v *ProductionOrderService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

func (v *ProductionOrderService) getRecordFormData(ctx *gin.Context) {

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

func (v *ProductionOrderService) splitScheduleOrders(ctx *gin.Context) {

	// first get the record
	componentName := ctx.Param("componentName")
	userId := common.GetUserId(ctx)
	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObjectOrder := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}

	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var recordMap = make(map[string]interface{})
	json.Unmarshal(generalObjectOrder.ObjectInfo, &recordMap)

	scheduleStatus := util.InterfaceToInt(recordMap["orderStatus"])

	if scheduleStatus == ScheduleStatusPreferenceOne || scheduleStatus == ScheduleStatusPreferenceTwo {
		orderStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
		if orderStatusId == -1 {
			response.DispatchDetailedError(ctx, AlreadyScheduled,
				&response.DetailedError{
					Header:      "Invalid Order Status",
					Description: "Invalid order status configured, please check the preference leaves are correctly configured before schedule it",
				})
			return
		}

		v.BaseService.Logger.Info("init scheduling")
		recordMap["id"] = intRecordId

		var scheduledOrderTableName string
		if targetTable == ProductionOrderMasterTable {
			scheduledOrderTableName = ScheduledOrderEventTable
		} else if targetTable == AssemblyProductionOrderTable {
			scheduledOrderTableName = AssemblyScheduledOrderEventTable
		} else {
			rangeObject := createRequest["range"].(map[string]interface{})
			duration := util.InterfaceToString(createRequest["duration"])
			typeOption := util.InterfaceToString(rangeObject["type"])

			if typeOption == "endDate" {
				endDate := util.InterfaceToString(rangeObject["endDate"])
				startDate := util.InterfaceToString(rangeObject["startDate"])
				startTime := util.InterfaceToString(createRequest["startTime"])
				listOfPartNo := util.InterfaceToIntArray(recordMap["partNo"])
				var inQuery string
				inQuery = ""
				if len(listOfPartNo) == 0 {
					inQuery = "-1"
				} else {
					for index, partNo := range listOfPartNo {
						if index == len(listOfPartNo)-1 {
							inQuery += strconv.Itoa(partNo)
						} else {
							inQuery += strconv.Itoa(partNo) + ","
						}

					}
				}
				conditionString := " id IN (" + inQuery + ")"
				generalPartObjects, _ := GetConditionalObjects(dbConnection, ToolingPartMasterTable, conditionString)

				if !isValidEndDate(endDate, startDate, startTime, generalPartObjects) {
					response.DispatchDetailedError(ctx, AlreadyScheduled,
						&response.DetailedError{
							Header:      "Invalid Operation",
							Description: "End date is lesser than expected date",
						})
					return
				}

			}

			listOfSplitOrders := v.SplitToolingSchedule(recordMap, dbConnection, generalObjectOrder.Id, createRequest, duration)
			orderStatusId = v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
			for _, toolingOrder := range listOfSplitOrders {
				jsonObject, _ := json.Marshal(toolingOrder)
				generalObject := component.GeneralObject{ObjectInfo: jsonObject}
				err, _ := Create(dbConnection, ToolingScheduledOrderEventTable, generalObject)
				if err != nil {
					v.BaseService.Logger.Error("error creating schedule", zap.String("error", err.Error()))
					continue
				}
			}

			updateProductionOrder := make(map[string]interface{})
			recordMap["orderStatus"] = orderStatusId

			serializedProductionOrderInfo, _ := json.Marshal(recordMap)
			updateProductionOrder["object_info"] = serializedProductionOrderInfo
			err = Update(dbConnection, targetTable, intRecordId, updateProductionOrder)

			ctx.JSON(http.StatusCreated, response.GeneralResponse{
				Code:    0,
				Message: "New schedule is successfully created",
			})
			return
		}

		errorResponse, preferenceLevel, scheduledQty := v.SplitSchedule(projectId, dbConnection, recordMap, orderStatusId, createRequest, scheduledOrderTableName)
		if errorResponse != nil {
			response.DispatchDetailedError(ctx, ErrorCreatingSchedule, errorResponse)
			return
		}
		orderStatusId = v.getOrderStatusId(dbConnection, preferenceLevel)
		if orderStatusId == -1 {
			response.DispatchDetailedError(ctx, AlreadyScheduled,
				&response.DetailedError{
					Header:      "Invalid Order Status",
					Description: "Invalid order status configured, please check the preference leaves are correctly configured before schedule it",
				})
			return
		}

		productionOrder := util.InterfaceToString(recordMap["prodOrder"])
		updateProductionOrder := make(map[string]interface{})
		recordMap["orderStatus"] = orderStatusId
		recordMap["remainingScheduledQty"] = util.InterfaceToInt(recordMap["remainingScheduledQty"]) - scheduledQty
		serializedProductionOrderInfo, _ := json.Marshal(recordMap)
		updateProductionOrder["object_info"] = serializedProductionOrderInfo
		Update(dbConnection, targetTable, intRecordId, updateProductionOrder)
		v.CreateActionRecordMessage(ProjectID, componentName, "Resource got updated", intRecordId, userId, createRequest)
		v.CreateUserRecordMessage(ProjectID, componentName, "Resource got updated", intRecordId, userId, &generalObjectOrder, &component.GeneralObject{ObjectInfo: serializedProductionOrderInfo})

		ctx.JSON(http.StatusCreated, response.GeneralResponse{
			Code:    0,
			Message: "Production order [" + productionOrder + "] has been successfully scheduled",
		})

	} else {
		//send error saying, it is already scheduled
		fmt.Println("Else condition")
		response.DispatchDetailedError(ctx, AlreadyScheduled,
			&response.DetailedError{
				Header:      "Already Scheduled",
				Description: "Invalid scheduling state, requested production order is already scheduled by the system",
			})
		return
	}

}

func (v *ProductionOrderService) resetSchedule(ctx *gin.Context) {

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
	var recordMap = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &recordMap)

	scheduleStatus := util.InterfaceToString(recordMap["status"])
	// status : 1 mean pending
	productionOrder := util.InterfaceToString(recordMap["prodOrder"])
	if scheduleStatus == "Scheduled" {
		v.BaseService.Logger.Info("scheduling it")
		condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.productionOrder\")) = \"" + productionOrder + "\""
		listOfSchedules, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, condition)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting schedules"), ErrorGettingObjectsInformation)
			return
		}
		for _, schedule := range *listOfSchedules {
			fmt.Println("schedule:", schedule)
			var recordMap = make(map[string]interface{})
			json.Unmarshal(schedule.ObjectInfo, &recordMap)
			eventStatus := util.InterfaceToString(recordMap["eventStatus"])
			if eventStatus == "confirmed" {
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("can not reset schedules now, there are schedules with confirmed stage, please make them un-confirm stage"), ErrorCreatingSchedule)
				return
			}
		}
		//for _, schedule := range *listOfSchedules{
		//	Delete(dbConnection,MachineTimelineEventTable,)
		//}

	} else {
		//send error saying, it is already scheduled
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid scheduling state, requested production order is already scheduled by the system"), ErrorCreatingSchedule)
		return
	}

}

func (v *ProductionOrderService) getSearchResults(ctx *gin.Context) {

	var searchFieldCommand []component.SearchKeys
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	if len(searchFieldCommand) == 0 {
		// reset the search
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getObjects(ctx)
		return
	}

	format := ctx.Query("format")

	searchQuery := v.ComponentManager.GetSearchQueryV2(dbConnection, componentName, searchFieldCommand)

	if searchQuery == "()" || searchQuery == "" {
		searchQuery = "(NULL)"
	}
	newSearchQuery := " id IN " + searchQuery
	listOfObjects, err := GetConditionalObjects(dbConnection, targetTable, newSearchQuery)
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

	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)
}

func (v *ProductionOrderService) deleteValidation(ctx *gin.Context) {

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
