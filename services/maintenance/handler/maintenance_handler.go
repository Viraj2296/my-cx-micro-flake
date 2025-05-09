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
func (v *MaintenanceService) loadFile(ctx *gin.Context) {
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
func (v *MaintenanceService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")

	//componentName := ctx.Param("componentName")
	//targetTable := ms.ComponentManager.GetTargetTable(componentName)
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//dbConnection := ms.BaseService.ServiceDatabases[projectId]
	//err, errorCode, totalRecords, listOfObjects := ms.ComponentManager.ImportData(dbConnection, componentName, importDataCommand)
	//if err != nil {
	//	ms.BaseService.Logger.Error("unable to import data", zap.String("error", err.Error()))
	//	response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
	//	return
	//}
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
	//
	//ctx.JSON(http.StatusOK, importDataResponse)
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
func (v *MaintenanceService) exportObjects(ctx *gin.Context) {
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
func (v *MaintenanceService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := v.ComponentManager.GetTableImportSchema(componentName)
	ctx.JSON(http.StatusOK, tableImportSchema)
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
func (v *MaintenanceService) getExportSchema(ctx *gin.Context) {
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
func (v *MaintenanceService) getObjects(ctx *gin.Context) {

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
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
	userId := common.GetUserId(ctx)
	//Have to next flag
	isNext := true

	userBasedQuery := " object_info ->>'$.assignedUserId' = " + strconv.Itoa(userId) + " "

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

		if componentName == MaintenanceWorkOrderMyTaskComponent || componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent || componentName == MyMouldMaintenancePreventiveWorkOrderTaskComponent || componentName == MyMouldMaintenanceCorrectiveWorkOrderTaskComponent {
			searchWithBaseQuery = searchWithBaseQuery + " AND " + userBasedQuery
		}
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
	} else {
		totalRecords = Count(dbConnection, targetTable)
		tableCondition := component.TableCondition(offsetValue, fields, values, condition)
		if componentName == MaintenanceWorkOrderMyTaskComponent {
			tableCondition = tableCondition + " AND " + userBasedQuery
		} else if componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent {
			tableCondition = tableCondition + " AND " + userBasedQuery
		}
		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, tableCondition)

		} else {
			var conditionString string
			limitVal, _ := strconv.Atoi(limitValue)
			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				var tableConditionDesc string
				if conditionString != "" {
					if offsetVal == -1 {
						tableConditionDesc = component.TableConditionV1(offsetValue, fields, values, condition)
					} else {
						tableConditionDesc = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
					}
					if tableConditionDesc != "" {
						conditionString = tableConditionDesc + " AND " + conditionString
					}
				} else {
					if offsetVal == -1 {
						conditionString = component.TableConditionV1(offsetValue, fields, values, condition)
					} else {
						conditionString = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
					}

				}

				if conditionString == "" {
					if componentName == MaintenanceWorkOrderMyTaskComponent || componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent || componentName == MyMouldMaintenancePreventiveWorkOrderTaskComponent || componentName == MyMouldMaintenanceCorrectiveWorkOrderTaskComponent {
						conditionString = userBasedQuery
					}
				} else {
					if componentName == MaintenanceWorkOrderMyTaskComponent || componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent || componentName == MyMouldMaintenancePreventiveWorkOrderTaskComponent || componentName == MyMouldMaintenanceCorrectiveWorkOrderTaskComponent {
						conditionString = conditionString + " AND " + userBasedQuery
					}
				}

				orderBy := "object_info ->> '$.createdAt' desc"

				listOfObjects, err = GetConditionalObjectsOrderBy(dbConnection, targetTable, conditionString, orderBy, limitVal)

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

				fmt.Println("=======================")
				fmt.Println("conditionString", conditionString)
				fmt.Println("limitVal", limitVal)
				//listOfObjects = reverseSlice(listOfObjects)
			} else {
				if componentName == MaintenanceWorkOrderMyTaskComponent || componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent || componentName == MyMouldMaintenancePreventiveWorkOrderTaskComponent || componentName == MyMouldMaintenanceCorrectiveWorkOrderTaskComponent {
					tableCondition = tableCondition + " AND " + userBasedQuery
				}
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, tableCondition, limitVal)
				currentRecordCount := len(*listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(tableCondition, "AND")
					var totalRecordObjects *[]component.GeneralObject
					if len(andClauses) > 1 {
						totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, tableCondition)

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

		currentRecordCount := len(*listOfObjects)
		limitVal, _ := strconv.Atoi(limitValue)
		if currentRecordCount < limitVal {
			isNext = false
		} else if currentRecordCount == limitVal {
			andClauses := strings.Split(tableCondition, "AND")
			var totalRecordObjects *[]component.GeneralObject
			if len(andClauses) > 1 {
				totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, tableCondition)

			} else {
				totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
			}
			lenTotalRecord := len(*totalRecordObjects)
			if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
				isNext = false
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
func (v *MaintenanceService) getCardView(ctx *gin.Context) {

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

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
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
func (v *MaintenanceService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	objectInfo := make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &objectInfo)

	err = ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
		return
	}

	if targetTable == MaintenanceWorkOrderTable {
		assetId := util.InterfaceToInt(objectInfo["assetId"])
		moduleName := util.InterfaceToString(objectInfo["moduleName"])
		objectStatusCorrectiveCondition := " object_info ->> '$.workOrderStatus' !=" + strconv.Itoa(WorkOrderDone) + " AND object_info ->> '$.objectStatus' = 'Active' AND object_info ->> '$.assetId' = " + strconv.Itoa(assetId)
		objectStatusPreventiveCondition := " object_info ->> '$.workOrderStatus' !=" + strconv.Itoa(WorkOrderDone) + " AND object_info ->> '$.objectStatus' = 'Active' AND object_info ->> '$.moduleName' = '" + moduleName + "' AND object_info ->> '$.assetId' = " + strconv.Itoa(assetId)

		totalNoOfPreventiveOrders := CountByCondition(dbConnection, MaintenanceWorkOrderTable, objectStatusPreventiveCondition)
		machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
		if moduleName == AssetClassMachines {
			totalNoOfCorrectiveOrders := CountByCondition(dbConnection, MaintenanceCorrectiveWorkOrderTable, objectStatusCorrectiveCondition)

			if totalNoOfPreventiveOrders == 0 && totalNoOfCorrectiveOrders == 0 {

				_ = machineService.MoveMachineToActive(projectId, assetId)
				_ = machineService.MoveMachineLiveStatusToActive(projectId, assetId)
			}
		} else {

			if totalNoOfPreventiveOrders == 0 {
				_ = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(assetId))
				_ = machineService.MoveMachineLiveStatusToActive(projectId, util.InterfaceToInt(assetId))
			}
		}
	} else if targetTable == MaintenanceCorrectiveWorkOrderTable {
		assetId := util.InterfaceToInt(objectInfo["assetId"])
		objectStatusCondition := " object_info ->> '$.workOrderStatus' !=" + strconv.Itoa(WorkOrderDone) + " AND object_info ->> '$.objectStatus' = 'Active' AND object_info ->> '$.assetId' = " + strconv.Itoa(assetId)
		totalNoOfCorrectiveOrders := CountByCondition(dbConnection, MaintenanceCorrectiveWorkOrderTable, objectStatusCondition)
		machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)

		if totalNoOfCorrectiveOrders == 0 {
			_ = machineService.MoveMachineToActive(projectId, assetId)
			_ = machineService.MoveMachineLiveStatusToActive(projectId, assetId)
		}
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
func (v *MaintenanceService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, recordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if targetTable == MaintenanceWorkOrderCorrectiveTaskTable {
		updateRequest["taskDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["taskDate"]))

		updateRequest["estimatedTaskEndDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["estimatedTaskEndDate"]))
		assignedUser := util.InterfaceToInt(updateRequest["assignedUserId"])

		originalObject := make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &originalObject)

		_, workOrderInterface := Get(dbConnection, MaintenanceCorrectiveWorkOrderTable, util.InterfaceToInt(originalObject["workOrderId"]))

		correctiveWorkOrderInfo := make(map[string]interface{})
		json.Unmarshal(workOrderInterface.ObjectInfo, &correctiveWorkOrderInfo)

		if !util.InterfaceToBool(originalObject["canUpdate"]) {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("can't update completed task"), ErrorUpdatingObjectInformation)
			return
		}

		if util.InterfaceToInt(correctiveWorkOrderInfo["workOrderStatus"]) > WorkOrderCreated {
			delete(updateRequest, "workOrderId")
		}

		err = v.emailGenerator(dbConnection, WorkOrderTaskAssignmentNotification, assignedUser, MaintenanceWorkOrderCorrectiveTaskTable, recordId)
		fmt.Println("Errror:", err)

	} else if targetTable == MouldMaintenanceCorrectiveWorkOrderTaskTable {
		updateRequest["taskDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["taskDate"]))

		updateRequest["estimatedTaskEndDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["estimatedTaskEndDate"]))
		assignedUser := util.InterfaceToInt(updateRequest["assignedUserId"])

		originalObject := make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &originalObject)

		_, workOrderInterface := Get(dbConnection, MouldMaintenanceCorrectiveWorkOrderTable, util.InterfaceToInt(originalObject["workOrderId"]))

		correctiveWorkOrderInfo := make(map[string]interface{})
		json.Unmarshal(workOrderInterface.ObjectInfo, &correctiveWorkOrderInfo)

		if !util.InterfaceToBool(originalObject["canUpdate"]) {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("can't update completed task"), ErrorUpdatingObjectInformation)
			return
		}

		if util.InterfaceToInt(correctiveWorkOrderInfo["workOrderStatus"]) > WorkOrderCreated {
			delete(updateRequest, "workOrderId")
		}

		err = v.emailGenerator(dbConnection, MouldWorkOrderTaskAssignmentNotification, assignedUser, MouldMaintenanceCorrectiveWorkOrderTaskTable, recordId)
		fmt.Println("Errror:", err)

	} else if targetTable == MaintenanceCorrectiveWorkOrderTable {
		originalObject := make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &originalObject)
		if !util.InterfaceToBool(originalObject["canUpdate"]) {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("can't update completed order"), ErrorUpdatingObjectInformation)
			return
		}
	} else if targetTable == MouldMaintenanceCorrectiveWorkOrderTable {
		originalObject := make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &originalObject)
		if !util.InterfaceToBool(originalObject["canUpdate"]) {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("can't update completed order"), ErrorUpdatingObjectInformation)
			return
		}
	} else if componentName == MaintenanceWorkOrderTaskComponent {
		updateRequest["taskDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["taskDate"]))

		updateRequest["estimatedTaskEndDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["estimatedTaskEndDate"]))
		assignedUser := util.InterfaceToInt(updateRequest["assignedUserId"])
		err = v.emailGenerator(dbConnection, PreventiveWorkOrderTaskAssignmentNotification, assignedUser, MaintenanceWorkOrderTaskTable, recordId)
	} else if componentName == MouldMaintenancePreventiveWorkOrderTaskComponent {
		updateRequest["taskDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["taskDate"]))

		updateRequest["estimatedTaskEndDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["estimatedTaskEndDate"]))
		assignedUser := util.InterfaceToInt(updateRequest["assignedUserId"])
		err = v.emailGenerator(dbConnection, MouldPreventiveWorkOrderTaskAssignmentNotification, assignedUser, MouldMaintenancePreventiveWorkOrderTaskTable, recordId)
	} else if componentName == MaintenanceWorkOrderComponent {
		updateRequest["workOrderScheduledEndDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["workOrderScheduledEndDate"]))
		updateRequest["workOrderScheduledStartDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["workOrderScheduledStartDate"]))
		updateRequest["remainderDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["remainderDate"]))
		updateRequest["remainderEndDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(updateRequest["remainderEndDate"]))
	}

	//Adding update preprocess request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}
	userId := common.GetUserId(ctx)

	if componentName == MaintenanceWorkOrderTaskComponent {

		if err, isDone := v.isAllTaskCompleted(dbConnection, recordId); err == nil {
			if isDone {
				// all the tasks are completed, now close the main work order
				err := v.completeWorkOrder(userId, dbConnection, recordId)
				if err != nil {
					v.BaseService.Logger.Error("Completing work order had issue, can not completed", zap.String("error", err.Error()))
				}
			}
		}
	}

	err = v.CreateUserRecordMessage(ProjectID, componentName, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	if componentName == MaintenanceWorkOrderTaskComponent {
		err = v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderMyTaskComponent, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	} else if componentName == MaintenanceWorkOrderCorrectiveTaskComponent {
		err = v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderMyCorrectiveTaskComponent, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	} else if componentName == MaintenanceWorkOrderMyTaskComponent {
		err = v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderTaskComponent, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	} else if componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent {
		err = v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderCorrectiveTaskComponent, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	} else if componentName == MouldMaintenancePreventiveWorkOrderTaskComponent {
		err = v.CreateUserRecordMessage(ProjectID, MyMouldMaintenancePreventiveWorkOrderTaskComponent, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	} else if componentName == MouldMaintenanceCorrectiveWorkOrderTaskComponent {
		err = v.CreateUserRecordMessage(ProjectID, MyMouldMaintenanceCorrectiveWorkOrderTaskComponent, "Resource got updated", recordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	}
	if err != nil {
		v.BaseService.Logger.Error("error creating record trail for updating resource", zap.Any("record_id", recordId))
	}

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
func (v *MaintenanceService) createNewResource(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createRequest["objectStatus"] = common.ObjectStatusActive

	tx := dbConnection.Begin()
	if tx == nil {
		v.BaseService.Logger.Error("Failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			v.BaseService.Logger.Error("Recovered from panic:", zap.Any("error", r))
			tx.Rollback()
			v.BaseService.Logger.Error("Transaction rolled back successfully")
		}
	}()

	var createdRecordId int
	var err error

	updatedRequest := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	// here we should do the validation
	err = v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		tx.Rollback()
		v.BaseService.Logger.Error("Transaction rollback initiated due to an error")
		v.BaseService.Logger.Error("Validation Failed:", zap.String("error", err.Error()))
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}
	userId := common.GetUserId(ctx)

	rawCreateRequest, _ := json.Marshal(updatedRequest)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}
	err, createdRecordId = Create(tx, targetTable, object)

	if err != nil {
		tx.Rollback()
		v.BaseService.Logger.Error("Transaction rollback initiated due to an error")
		v.BaseService.Logger.Error("error creating machines register information:", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error creating machines register information"), ErrorCreatingObjectInformation)
		return
	}

	switch componentName {
	case MaintenanceWorkOrderComponent:

		_, generalWorkOrder := Get(tx, MaintenanceWorkOrderTable, createdRecordId)

		workOrder := MaintenanceWorkOrder{ObjectInfo: generalWorkOrder.ObjectInfo}
		workOrderInfo := workOrder.getWorkOrderInfo()
		if createdRecordId < 10 {
			workOrderInfo.WorkOrderReferenceId = "W0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			workOrderInfo.WorkOrderReferenceId = "W000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			workOrderInfo.WorkOrderReferenceId = "W00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			workOrderInfo.WorkOrderReferenceId = "W0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			workOrderInfo.WorkOrderReferenceId = "W" + strconv.Itoa(createdRecordId)
		}
		//Added can status
		workOrderInfo.CanUpdate = true
		workOrderInfo.CanRelease = true
		workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)

		workOrderInfo.WorkOrderScheduledEndDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.WorkOrderScheduledEndDate)
		workOrderInfo.WorkOrderScheduledStartDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.WorkOrderScheduledStartDate)
		workOrderInfo.RemainderDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.RemainderDate)
		workOrderInfo.RemainderEndDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.RemainderEndDate)

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(workOrderInfo)
		updatingData["object_info"] = rawWorkOrderInfo
		Update(tx, MaintenanceWorkOrderTable, createdRecordId, updatingData)
		if workOrder.getWorkOrderInfo().ModuleName == AssetClassMoulds {
			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			if workOrder.getWorkOrderInfo().WorkOrderType == WorkOrderTypeCorrective {
				mouldInterface.PutToRepairMode(projectId, workOrder.getWorkOrderInfo().AssetId, userId)
			}
		}

		// Have to create default work order task

		if workOrder.getWorkOrderInfo().ModuleName == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			_ = machineService.MoveMachineToMaintenance(projectId, workOrderInfo.AssetId)
			_ = machineService.MoveMachineLiveStatusToMaintenance(projectId, workOrderInfo.AssetId)
		}

		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(createdRecordId)
		listOfTasksObjects, _ := GetConditionalObjects(tx, MaintenanceWorkOrderTaskComponent, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)
		var workOrderTaskId string
		if alreadyAvailableTasks < 10 {
			workOrderTaskId = "PWT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskId = "PWT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskId = "PWT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskId = "PWT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskId = "PWT" + strconv.Itoa(createdRecordId)
		}

		workOrderTaskInfo := WorkOrderTaskInfo{
			TaskName:             workOrderInfo.Name,
			ShortDescription:     workOrderInfo.Description,
			TaskDate:             util.GetCurrentTime(ISOTimeLayout),
			AssignedUserId:       userId,
			WorkOrderId:          createdRecordId,
			CanApprove:           false,
			IsOrderReleased:      false,
			CanReject:            false,
			CanCheckIn:           true,
			CanCheckOut:          false,
			TaskStatus:           1,
			WorkOrderTaskId:      workOrderTaskId,
			ObjectStatus:         common.ObjectStatusActive,
			EstimatedTaskEndDate: workOrderInfo.WorkOrderScheduledEndDate,
			LastUpdatedAt:        util.GetCurrentTime(ISOTimeLayout),
			CreatedAt:            util.GetCurrentTime(ISOTimeLayout),
		}

		rawCreateTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		workOrderTaskObject := component.GeneralObject{
			ObjectInfo: rawCreateTaskInfo,
		}

		err, taskCreatedRecordId := Create(tx, MaintenanceWorkOrderTaskTable, workOrderTaskObject)
		if err != nil {
			tx.Rollback()
			v.BaseService.Logger.Error("Transaction rollback initiated due to an error")
			v.BaseService.Logger.Error("error creating work order task:", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error creating work order task"), ErrorCreatingObjectInformation)
			return
		}

		v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderTaskTable, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderMyTaskComponent, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenancePreventiveWorkOrderStatusTable, WorkOrderCreated)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		//listOfWorkflowUsers := workOrderInfo.Supervisors
		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, WorkOrderSupervisorNotication, workflowUser, MaintenanceWorkOrderComponent, createdRecordId)
				fmt.Println("Errror:", err)
			}
		}

	case MaintenanceCorrectiveWorkOrderComponent:
		productionContinue := util.InterfaceToInt(createRequest["canContinue"])
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderCreated)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}

		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		_, generalWorkOrder := Get(tx, MaintenanceCorrectiveWorkOrderComponent, createdRecordId)

		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(generalWorkOrder.ObjectInfo, &workOrderInfo)
		if createdRecordId < 10 {
			workOrderInfo["workOrderReferenceId"] = "W0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			workOrderInfo["workOrderReferenceId"] = "W000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			workOrderInfo["workOrderReferenceId"] = "W00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			workOrderInfo["workOrderReferenceId"] = "W0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			workOrderInfo["workOrderReferenceId"] = "W" + strconv.Itoa(createdRecordId)
		}
		//Added can status
		workOrderInfo["canUpdate"] = true
		workOrderInfo["canRelease"] = true

		workOrderInfo["workOrderScheduledStartDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(workOrderInfo["workOrderScheduledStartDate"]))

		currentDate := time.Now().UTC()
		var scheduleEndDate string
		isOrderRelease := false

		lastTimeOfCurrentDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 59, 000, time.UTC)
		// convert to UTC time
		utcLastTime := lastTimeOfCurrentDate.Add(time.Hour * time.Duration(-8))
		formattedTime := utcLastTime.Format(ISOTimeLayout)

		if productionContinue == ProductionContinued {

			workOrderInfo["workOrderScheduledEndDate"] = formattedTime
			scheduleEndDate = formattedTime

			if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
				for _, workflowUser := range listOfWorkflowUsers {
					err = v.emailGenerator(dbConnection, CorrectiveWorkOrderCreateNotification, workflowUser, MaintenanceCorrectiveWorkOrderComponent, createdRecordId)
					fmt.Println("Errror:", err)
				}
			}

		} else {
			// Need to check machine has production order
			assetId := util.InterfaceToInt(createRequest["assetId"])
			machineInfo := make(map[string]interface{})
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			_, machineMasterGeneralObject := machineService.GetMachineInfoById(projectId, assetId)
			json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineInfo)

			var workOrderScheduledEndDate string
			enableProductionOrder := util.InterfaceToBool(machineInfo["enableProductionOrder"])
			if enableProductionOrder {
				_, generalObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				scheduledOrderInfo := make(map[string]interface{})
				json.Unmarshal(generalObject.ObjectInfo, &scheduledOrderInfo)
				workOrderScheduledEndDate = util.InterfaceToString(scheduledOrderInfo["endDate"])
			} else {
				workOrderScheduledEndDate = formattedTime
			}

			workOrderInfo["workOrderScheduledStartDate"] = util.GetCurrentTime(ISOTimeLayout)
			workOrderInfo["workOrderScheduledEndDate"] = workOrderScheduledEndDate

			scheduleEndDate = workOrderScheduledEndDate
			isOrderRelease = true
			workOrderInfo["canRelease"] = false
			workOrderInfo["workOrderStatus"] = WorkOrderScheduled
			workOrderInfo["canUpdate"] = false
			workOrderInfo["canRelease"] = false
			workOrderInfo["canUnRelease"] = true
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(workOrderInfo)
		updatingData["object_info"] = rawWorkOrderInfo
		Update(tx, MaintenanceCorrectiveWorkOrderComponent, createdRecordId, updatingData)
		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

			mouldInterface.PutToRepairMode(projectId, util.InterfaceToInt(workOrderInfo["assetId"]), userId)

		}

		// Have to create default work order task

		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			_ = machineService.MoveMachineToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			_ = machineService.MoveMachineLiveStatusToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		}

		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(createdRecordId)
		listOfTasksObjects, _ := GetConditionalObjects(tx, MaintenanceWorkOrderCorrectiveTaskComponent, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)
		var workOrderTaskId string
		if alreadyAvailableTasks < 10 {
			workOrderTaskId = "WT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskId = "WT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskId = "WT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskId = "WT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskId = "WT" + strconv.Itoa(createdRecordId)
		}

		workOrderTaskInfo := WorkOrderTaskInfo{
			TaskName:             util.InterfaceToString(workOrderInfo["description"]),
			ShortDescription:     util.InterfaceToString(workOrderInfo["description"]),
			TaskDate:             util.GetCurrentTime(ISOTimeLayout),
			AssignedUserId:       -1, // assigned to no one. it will be handled by front-end later
			WorkOrderId:          createdRecordId,
			CanApprove:           false,
			IsOrderReleased:      isOrderRelease,
			CanReject:            false,
			CanCheckIn:           true,
			CanCheckOut:          false,
			CanUpdate:            true,
			TaskStatus:           1,
			ObjectStatus:         "Active",
			WorkOrderTaskId:      workOrderTaskId,
			EstimatedTaskEndDate: scheduleEndDate,
			LastUpdatedAt:        util.GetCurrentTime(ISOTimeLayout),
			CreatedAt:            util.GetCurrentTime(ISOTimeLayout),
		}

		rawCreateTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		workOrderTaskObject := component.GeneralObject{
			ObjectInfo: rawCreateTaskInfo,
		}

		// var x *int
		// fmt.Println(*x)

		err, taskCreatedRecordId := Create(tx, MaintenanceWorkOrderCorrectiveTaskComponent, workOrderTaskObject)
		if err != nil {
			tx.Rollback()
			v.BaseService.Logger.Error("Transaction rollback initiated due to an error")
			v.BaseService.Logger.Error("error creating work order task:", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error creating work order task"), ErrorCreatingObjectInformation)
			return
		}

		v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderCorrectiveTaskComponent, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderMyCorrectiveTaskComponent, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		//listOfWorkflowUsers := workOrderInfo.Supervisors
		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, CorrectiveWorkOrderSupervisorNotification, workflowUser, MaintenanceCorrectiveWorkOrderComponent, createdRecordId)
				fmt.Println("Errror:", err)
			}
		}

	case MaintenanceWorkOrderTaskComponent:
		workOrderId := util.InterfaceToInt(createRequest["workOrderId"])

		// check this work order is already done, then don't allow to create new task
		if err, isTaskIsDone := v.isWorkOrderCompletedFromWorkOrderId(dbConnection, workOrderId); err == nil {
			if isTaskIsDone {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, this work order is already been done, adding any task is not allowed",
					})
				return
			}
		}
		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(workOrderId)
		listOfTasksObjects, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)

		// get the work order task just inserted
		_, generalWorkOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, createdRecordId)

		workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: generalWorkOrderTask.ObjectInfo}
		workOrderTaskInfo := workOrderTask.getWorkOrderTaskInfo()
		if alreadyAvailableTasks < 10 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT" + strconv.Itoa(createdRecordId)
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		updatingData["object_info"] = rawWorkOrderTaskInfo

		Update(dbConnection, MaintenanceWorkOrderTaskTable, createdRecordId, updatingData)
	case MaintenanceWorkOrderCorrectiveTaskComponent:
		workOrderId := util.InterfaceToInt(createRequest["workOrderId"])

		// check this work order is already done, then don't allow to create new task
		if err, isTaskIsDone := v.isWorkOrderCompletedFromWorkOrderId(dbConnection, workOrderId); err == nil {
			if isTaskIsDone {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, this work order is already been done, adding any task is not allowed",
					})
				return
			}
		}
		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(workOrderId)
		listOfTasksObjects, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)

		// get the work order task just inserted
		_, generalWorkOrderTask := Get(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, createdRecordId)

		workOrderTaskInfo := make(map[string]interface{})
		json.Unmarshal(generalWorkOrderTask.ObjectInfo, &workOrderTaskInfo)
		if alreadyAvailableTasks < 10 {
			workOrderTaskInfo["workOrderTaskId"] = "WT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskInfo["workOrderTaskId"] = "WT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskInfo["workOrderTaskId"] = "WT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskInfo["workOrderTaskId"] = "WT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskInfo["workOrderTaskId"] = "WT" + strconv.Itoa(createdRecordId)
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		updatingData["object_info"] = rawWorkOrderTaskInfo

		Update(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, createdRecordId, updatingData)

	case MouldMaintenancePreventiveWorkOrderComponent:

		_, generalWorkOrder := Get(tx, MouldMaintenancePreventiveWorkOrderTable, createdRecordId)

		workOrder := MouldMaintenancePreventiveWorkOrder{ObjectInfo: generalWorkOrder.ObjectInfo}
		workOrderInfo := workOrder.getWorkOrderInfo()
		if createdRecordId < 10 {
			workOrderInfo.WorkOrderReferenceId = "W0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			workOrderInfo.WorkOrderReferenceId = "W000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			workOrderInfo.WorkOrderReferenceId = "W00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			workOrderInfo.WorkOrderReferenceId = "W0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			workOrderInfo.WorkOrderReferenceId = "W" + strconv.Itoa(createdRecordId)
		}
		//Added can status
		workOrderInfo.CanUpdate = true
		workOrderInfo.CanRelease = true
		workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
		workOrderInfo.ModuleName = AssetClassMoulds
		workOrderInfo.WorkOrderScheduledEndDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.WorkOrderScheduledEndDate)
		workOrderInfo.WorkOrderScheduledStartDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.WorkOrderScheduledStartDate)
		workOrderInfo.RemainderDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.RemainderDate)
		workOrderInfo.RemainderEndDate = util.ConvertSingaporeTimeToUTC(workOrderInfo.RemainderEndDate)

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(workOrderInfo)
		updatingData["object_info"] = rawWorkOrderInfo
		Update(tx, MouldMaintenancePreventiveWorkOrderTable, createdRecordId, updatingData)
		// mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		// if workOrder.getWorkOrderInfo().WorkOrderType == WorkOrderTypeCorrective {
		// 	mouldInterface.PutToRepairMode(projectId, workOrder.getWorkOrderInfo().AssetId, userId)
		// }
		if util.InterfaceToString(workOrderInfo.ModuleName) == AssetClassMoulds {

			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

			mouldInterface.PutToRepairMode(projectId, util.InterfaceToInt(workOrderInfo.MouldId), userId)

		}

		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(createdRecordId)
		listOfTasksObjects, _ := GetConditionalObjects(tx, MouldMaintenancePreventiveWorkOrderTaskComponent, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)
		var workOrderTaskId string
		if alreadyAvailableTasks < 10 {
			workOrderTaskId = "PWT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskId = "PWT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskId = "PWT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskId = "PWT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskId = "PWT" + strconv.Itoa(createdRecordId)
		}

		workOrderTaskInfo := WorkOrderTaskInfo{
			TaskName:             workOrderInfo.Name,
			ShortDescription:     workOrderInfo.Description,
			TaskDate:             util.GetCurrentTime(ISOTimeLayout),
			AssignedUserId:       userId,
			WorkOrderId:          createdRecordId,
			CanApprove:           false,
			IsOrderReleased:      false,
			CanReject:            false,
			CanCheckIn:           true,
			CanCheckOut:          false,
			TaskStatus:           1,
			WorkOrderTaskId:      workOrderTaskId,
			ObjectStatus:         common.ObjectStatusActive,
			EstimatedTaskEndDate: workOrderInfo.WorkOrderScheduledEndDate,
			LastUpdatedAt:        util.GetCurrentTime(ISOTimeLayout),
			CreatedAt:            util.GetCurrentTime(ISOTimeLayout),
		}

		rawCreateTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		workOrderTaskObject := component.GeneralObject{
			ObjectInfo: rawCreateTaskInfo,
		}

		err, taskCreatedRecordId := Create(tx, MouldMaintenancePreventiveWorkOrderTaskTable, workOrderTaskObject)
		if err != nil {
			tx.Rollback()
			v.BaseService.Logger.Error("Transaction rollback initiated due to an error")
			v.BaseService.Logger.Error("error creating work order task:", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error creating work order task"), ErrorCreatingObjectInformation)
			return
		}

		v.CreateUserRecordMessage(ProjectID, MouldMaintenancePreventiveWorkOrderTaskTable, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderMyTaskComponent, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenancePreventiveWorkOrderStatusTable, WorkOrderCreated)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		//listOfWorkflowUsers := workOrderInfo.Supervisors
		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, MouldPreventiveWorkOrderAssignmentNotification, workflowUser, MouldMaintenancePreventiveWorkOrderComponent, createdRecordId)
				fmt.Println("Errror:", err)
			}
		}
	case MouldMaintenanceCorrectiveWorkOrderComponent:
		productionContinue := util.InterfaceToInt(createRequest["canContinue"])
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		listOfNotification := authService.GetUserInfoFromGroupId(v.ToolingSupervisorGroup)

		_, generalWorkOrder := Get(tx, MouldMaintenanceCorrectiveWorkOrderComponent, createdRecordId)

		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(generalWorkOrder.ObjectInfo, &workOrderInfo)
		if createdRecordId < 10 {
			workOrderInfo["workOrderReferenceId"] = "W0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			workOrderInfo["workOrderReferenceId"] = "W000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			workOrderInfo["workOrderReferenceId"] = "W00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			workOrderInfo["workOrderReferenceId"] = "W0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			workOrderInfo["workOrderReferenceId"] = "W" + strconv.Itoa(createdRecordId)
		}
		//Added can status
		workOrderInfo["canUpdate"] = true
		workOrderInfo["canRelease"] = true
		// release the orderw without checking anything
		workOrderInfo["workOrderStatus"] = WorkOrderScheduled
		workOrderInfo["workOrderScheduledStartDate"] = util.ConvertSingaporeTimeToUTC(util.InterfaceToString(workOrderInfo["workOrderScheduledStartDate"]))

		currentDate := time.Now().UTC()
		var scheduleEndDate string
		isOrderRelease := false

		lastTimeOfCurrentDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 59, 000, time.UTC)
		// convert to UTC time
		utcLastTime := lastTimeOfCurrentDate.Add(time.Hour * time.Duration(-8))
		formattedTime := utcLastTime.Format(ISOTimeLayout)

		if productionContinue == ProductionContinued {

			workOrderInfo["workOrderScheduledEndDate"] = formattedTime
			scheduleEndDate = formattedTime

			for _, configuredUserInfo := range listOfNotification {
				err = v.emailGenerator(dbConnection, MouldCorrectiveWorkOrderCreateNotification, configuredUserInfo.UserId, MouldMaintenanceCorrectiveWorkOrderComponent, createdRecordId)
				if err != nil {
					v.BaseService.Logger.Info("sending email for corrective order has failed", zap.String("error", err.Error()))
				} else {
					v.BaseService.Logger.Info("sending email for corrective order", zap.Int("user", configuredUserInfo.UserId))
				}

			}

		} else {
			// Need to check machine has production order
			assetId := util.InterfaceToInt(createRequest["assetId"])
			machineInfo := make(map[string]interface{})
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			_, machineMasterGeneralObject := machineService.GetMachineInfoById(projectId, assetId)
			json.Unmarshal(machineMasterGeneralObject.ObjectInfo, &machineInfo)

			var workOrderScheduledEndDate string
			enableProductionOrder := util.InterfaceToBool(machineInfo["enableProductionOrder"])
			if enableProductionOrder {
				_, generalObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				scheduledOrderInfo := make(map[string]interface{})
				json.Unmarshal(generalObject.ObjectInfo, &scheduledOrderInfo)
				workOrderScheduledEndDate = util.InterfaceToString(scheduledOrderInfo["endDate"])
			} else {
				workOrderScheduledEndDate = formattedTime
			}

			workOrderInfo["workOrderScheduledStartDate"] = util.GetCurrentTime(ISOTimeLayout)
			workOrderInfo["workOrderScheduledEndDate"] = workOrderScheduledEndDate

			scheduleEndDate = workOrderScheduledEndDate
			isOrderRelease = true
			workOrderInfo["canRelease"] = false
			workOrderInfo["workOrderStatus"] = WorkOrderScheduled
			workOrderInfo["canUpdate"] = false
			workOrderInfo["canRelease"] = false
			workOrderInfo["canUnRelease"] = true
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(workOrderInfo)
		updatingData["object_info"] = rawWorkOrderInfo
		Update(tx, MouldMaintenanceCorrectiveWorkOrderComponent, createdRecordId, updatingData)

		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {

			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)

			mouldInterface.PutToRepairMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)

		}

		// Have to create default work order task

		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			_ = machineService.MoveMachineToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			_ = machineService.MoveMachineLiveStatusToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		}

		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(createdRecordId)
		listOfTasksObjects, _ := GetConditionalObjects(tx, MouldMaintenanceCorrectiveWorkOrderTaskComponent, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)
		var workOrderTaskId string
		if alreadyAvailableTasks < 10 {
			workOrderTaskId = "WT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskId = "WT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskId = "WT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskId = "WT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskId = "WT" + strconv.Itoa(createdRecordId)
		}

		workOrderTaskInfo := WorkOrderTaskInfo{
			TaskName:             util.InterfaceToString(workOrderInfo["description"]),
			ShortDescription:     util.InterfaceToString(workOrderInfo["description"]),
			TaskDate:             util.GetCurrentTime(ISOTimeLayout),
			AssignedUserId:       -1, // assigned to no one. it will be handled by front-end later
			WorkOrderId:          createdRecordId,
			CanApprove:           false,
			IsOrderReleased:      isOrderRelease,
			CanReject:            false,
			CanCheckIn:           true,
			CanCheckOut:          false,
			CanUpdate:            true,
			TaskStatus:           1,
			ObjectStatus:         "Active",
			WorkOrderTaskId:      workOrderTaskId,
			EstimatedTaskEndDate: scheduleEndDate,
			LastUpdatedAt:        util.GetCurrentTime(ISOTimeLayout),
			CreatedAt:            util.GetCurrentTime(ISOTimeLayout),
		}

		rawCreateTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		workOrderTaskObject := component.GeneralObject{
			ObjectInfo: rawCreateTaskInfo,
		}

		_, taskCreatedRecordId := Create(tx, MouldMaintenanceCorrectiveWorkOrderTaskComponent, workOrderTaskObject)
		if err != nil {
			tx.Rollback()
			v.BaseService.Logger.Error("Transaction rollback initiated due to an error")
			v.BaseService.Logger.Error("error creating work order task:", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error creating work order task"), ErrorCreatingObjectInformation)
			return
		}

		v.CreateUserRecordMessage(ProjectID, MouldMaintenanceCorrectiveWorkOrderTaskComponent, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		v.CreateUserRecordMessage(ProjectID, MaintenanceWorkOrderMyCorrectiveTaskComponent, "New resource is created", taskCreatedRecordId, userId, nil, nil)

		// Send notification to configured user in creation status

		for _, configuredUserInfo := range listOfNotification {
			err = v.emailGenerator(dbConnection, MouldCorrectiveWorkOrderAssignmentNotification, configuredUserInfo.UserId, MouldMaintenanceCorrectiveWorkOrderComponent, createdRecordId)
			if err != nil {
				v.BaseService.Logger.Info("sending email for corrective order has failed", zap.String("error", err.Error()))
			} else {
				v.BaseService.Logger.Info("sending email for corrective order", zap.Int("user", configuredUserInfo.UserId))
			}

		}
	case MouldMaintenanceCorrectiveWorkOrderTaskComponent:
		workOrderId := util.InterfaceToInt(createRequest["workOrderId"])

		// check this work order is already done, then don't allow to create new task
		if err, isTaskIsDone := v.isWorkOrderCompletedFromWorkOrderId(dbConnection, workOrderId); err == nil {
			if isTaskIsDone {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, this work order is already been done, adding any task is not allowed",
					})
				return
			}
		}
		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(workOrderId)
		listOfTasksObjects, _ := GetConditionalObjects(dbConnection, MouldMaintenanceCorrectiveWorkOrderTaskComponent, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)

		// get the work order task just inserted
		_, generalWorkOrderTask := Get(dbConnection, MouldMaintenanceCorrectiveWorkOrderTaskComponent, createdRecordId)

		workOrderTaskInfo := make(map[string]interface{})
		json.Unmarshal(generalWorkOrderTask.ObjectInfo, &workOrderTaskInfo)
		if alreadyAvailableTasks < 10 {
			workOrderTaskInfo["workOrderTaskId"] = "WT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskInfo["workOrderTaskId"] = "WT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskInfo["workOrderTaskId"] = "WT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskInfo["workOrderTaskId"] = "WT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskInfo["workOrderTaskId"] = "WT" + strconv.Itoa(createdRecordId)
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		updatingData["object_info"] = rawWorkOrderTaskInfo

		Update(dbConnection, MouldMaintenanceCorrectiveWorkOrderTaskComponent, createdRecordId, updatingData)
	case MouldMaintenancePreventiveWorkOrderTaskComponent:
		workOrderId := util.InterfaceToInt(createRequest["workOrderId"])

		// check this work order is already done, then don't allow to create new task
		if err, isTaskIsDone := v.isWorkOrderCompletedFromWorkOrderId(dbConnection, workOrderId); err == nil {
			if isTaskIsDone {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, this work order is already been done, adding any task is not allowed",
					})
				return
			}
		}
		conditionString := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) =  " + strconv.Itoa(workOrderId)
		listOfTasksObjects, _ := GetConditionalObjects(dbConnection, MouldMaintenancePreventiveWorkOrderTaskTable, conditionString)
		alreadyAvailableTasks := len(*listOfTasksObjects)

		// get the work order task just inserted
		_, generalWorkOrderTask := Get(dbConnection, MouldMaintenancePreventiveWorkOrderTaskTable, createdRecordId)

		workOrderTask := MouldMaintenancePreventiveWorkOrderTask{ObjectInfo: generalWorkOrderTask.ObjectInfo}
		workOrderTaskInfo := workOrderTask.getWorkOrderTaskInfo()
		if alreadyAvailableTasks < 10 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT0000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT000" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 1000 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT00" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 10000 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT0" + strconv.Itoa(createdRecordId)
		} else if alreadyAvailableTasks < 100000 {
			workOrderTaskInfo.WorkOrderTaskId = "PWT" + strconv.Itoa(createdRecordId)
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderTaskInfo, _ := json.Marshal(workOrderTaskInfo)
		updatingData["object_info"] = rawWorkOrderTaskInfo

		Update(dbConnection, MouldMaintenancePreventiveWorkOrderTaskTable, createdRecordId, updatingData)

	}

	err = tx.Commit().Error
	if err != nil {
		v.BaseService.Logger.Error("Transaction commit failed", zap.String("error", err.Error()))
		tx.Rollback()
		v.BaseService.Logger.Error("Rollback due to commit failure")
		response.SendSimpleError(ctx, http.StatusInternalServerError, errors.New("transaction commit failed"), ErrorCreatingObjectInformation)
		return
	}
	v.BaseService.Logger.Info("Transaction committed successfully", zap.Int("createdRecordId", createdRecordId))

	url := "/project/" + projectId + "/maintenance/component/" + componentName + "/" + strconv.Itoa(createdRecordId)

	v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	ctx.Writer.Header().Set("Location", url)
	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully created",
		Error:   0,
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
func (v *MaintenanceService) getNewRecord(ctx *gin.Context) {

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
func (v *MaintenanceService) getRecordFormData(ctx *gin.Context) {

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
func (v *MaintenanceService) getSearchResults(ctx *gin.Context) {

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

	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)

}

// getDashboardSnapshot ShowAccount godoc
// @Summary Get the record form data to facilitate the update
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags Dashboards
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/machine_overview [get]
func (v *MaintenanceService) getMaintenanceOverview(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	summaryResponse := v.getSummaryResponse(projectId)
	ctx.JSON(http.StatusOK, summaryResponse)
}

func (v *MaintenanceService) removeInternalArrayReference(ctx *gin.Context) {

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
	err, objectInterface := Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	fmt.Println("objectInterface:", objectInterface)
	serializedObject := v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject
	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}
	var updatingObjectFields map[string]interface{}
	json.Unmarshal(serializedObject, &updatingObjectFields)
	ctx.JSON(http.StatusOK, updatingObjectFields)

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
func (v *MaintenanceService) getGroupBy(ctx *gin.Context) {

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
func (v *MaintenanceService) getCardViewGroupBy(ctx *gin.Context) {

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

type GroupByCardView struct {
	GroupByField string                   `json:"groupByField"`
	Cards        []map[string]interface{} `json:"cards"`
}
