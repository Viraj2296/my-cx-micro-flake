package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}
type Assignment struct {
	Event    int `json:"int"`
	Id       int `json:"id"`
	Resource int `json:"resource"`
}
type Job struct {
	Id       int    `json:"id"`
	Label    string `json:"label"`
	Location int    `json:"location"`
}

type OrderScheduledEventUpdateRequest struct {
	EndDate    string      `json:"endDate"`
	StartDate  string      `json:"startDate"`
	Assignment *Assignment `json:"assignment"`
}

type PrinterConfiguration struct {
	Id          int    `json:"id"`
	NetworkIP   string `json:"networkIP"`
	NetworkPort int    `json:"networkPort"`
	LocationId  int    `json:"locationId"`
}

type ServiceRequestAction struct {
	Remark string `json:"remark"`
}

func (v *BatchManagementService) recordPOSTActionHandler(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "submit" {
		v.ActionService.PrintDocument(ctx)
		return

	} else if actionName == "print_mould" {
		v.printMouldBatch(ctx)
		return
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}

func (v *BatchManagementService) handleComponentAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "complete_batch" {
		v.completeBatch(ctx)
		return

	} else if actionName == "get_raw_material" {
		v.handleGetRawMaterial(ctx)
		return

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}
func (v *BatchManagementService) recordGetActionHandler(ctx *gin.Context) {

}

func (v *BatchManagementService) completeBatch(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//userId := common.GetUserId(ctx)
	err, completeBatchRequest := component.GetRequestFields(ctx)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// now get the mould resource Id
	var mouldBatchResourceId = util.InterfaceToInt(completeBatchRequest["mouldBatchResourceId"])
	var noOfLabels = util.InterfaceToInt(completeBatchRequest["noOfLabels"])
	v.BaseService.Logger.Info("received complete batch request", zap.Any("request", completeBatchRequest))
	err, generalObject := database.Get(dbConnection, targetTable, mouldBatchResourceId)
	var mouldBatchFields = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &mouldBatchFields)

	var updatingData = make(map[string]interface{})
	mouldBatchFields["noOfLabels"] = noOfLabels
	mouldBatchFields["labelStatus"] = const_util.PrintingStatus

	serializedObject, _ := json.Marshal(mouldBatchFields)

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, mouldBatchResourceId, updatingData)

	//var mouldBatchId = util.InterfaceToString(mouldBatchFields["mouldBatchId"])
	//var rawMaterialBatchId = util.InterfaceToInt(mouldBatchFields["rawMaterialId"])
	//qaInterface := common.GetService("qa").ServiceInterface.(common.QAInterface)
	//
	//var qaObjectFields = make(map[string]interface{})
	//qaObjectFields["objectStatus"] = component.ObjectStatusActive
	//var createdAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	//qaObjectFields["createdAt"] = createdAt
	//qaObjectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	//qaObjectFields["createdBy"] = userId
	//qaObjectFields["lastUpdatedBy"] = userId
	//qaObjectFields["mouldBatchId"] = mouldBatchId
	//qaObjectFields["rawMaterialId"] = rawMaterialBatchId
	//qaObjectFields["qualityStatus"] = 1
	//var actionRemarks = make([]interface{}, 0)
	//actionRemarks = append(actionRemarks, ActionRemarks{
	//	ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
	//	Status:        "QA EVALUATION CREATED",
	//	UserId:        userId,
	//	Remarks:       "Great, Evaluation for quality assurance is created",
	//	ProcessedTime: getTimeDifference(util.InterfaceToString(createdAt)),
	//})
	//qaObjectFields["actionRemarks"] = actionRemarks
	//
	//serialisedData, _ := json.Marshal(qaObjectFields)
	//qaInterface.CreateQAResource(serialisedData)
	//v.BaseService.Logger.Info("creating qa mould batch", zap.String("mould_batch", mouldBatchId))
	// now insert this into qa table
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully completed the batch",
	})

}

func (v *BatchManagementService) printMouldBatch(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	recordIdString := ctx.Param("recordId")
	recordId, _ := strconv.Atoi(recordIdString)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//userId := common.GetUserId(ctx)
	err, printMouldBatchPayload := component.GetRequestFields(ctx)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var noOfLabels = util.InterfaceToInt(printMouldBatchPayload["noOfLabels"])
	v.BaseService.Logger.Info("received complete batch request", zap.Any("request", printMouldBatchPayload))
	err, generalObject := database.Get(dbConnection, targetTable, recordId)

	updatingData := make(map[string]interface{})

	mouldBatchInfo := make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &mouldBatchInfo)

	mouldBatchInfo["noOfLabels"] = noOfLabels
	mouldBatchInfo["labelStatus"] = const_util.PrintingStatus

	serializedObject, _ := json.Marshal(mouldBatchInfo)

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)

	// now insert this into qa table
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully mould batch is updated",
	})

}

func (v *BatchManagementService) handleGetRawMaterial(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	// it is only for ComponentBatchManagementRawMaterial
	if componentName == const_util.ComponentBatchManagementRawMaterial {
		targetTable := v.ComponentManager.GetTargetTable(componentName)
		dbConnection := v.BaseService.ServiceDatabases[projectId]
		type RawMaterialInfo struct {
			Name          string `json:"name"`
			RawMaterialId int    `json:"rawMaterialId"`
			Location      int    `json:"location"`
		}
		err, rawMaterialRequest := component.GetRequestFields(ctx)

		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		// now get the mould resource Id
		var batchId = util.InterfaceToString(rawMaterialRequest["batchId"])

		responseObject := RawMaterialInfo{}
		var condition = "object_info->>'$.batchId' = '" + batchId + "'"
		err, generalObjects := database.GetConditionalObjects(dbConnection, targetTable, condition)

		if err == nil {
			var rawMaterialInterface = (generalObjects)[0]
			rawMaterialInfo := database.GeRawMaterialBatchInfo(rawMaterialInterface.ObjectInfo)
			responseObject.RawMaterialId = rawMaterialInterface.Id
			responseObject.Name = rawMaterialInfo.Name
			responseObject.Location = rawMaterialInfo.Location

		}
		ctx.JSON(http.StatusOK, responseObject)
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Component",
				Description: "Your action can not be performed against this request due to invalid component",
			})
	}

}
