package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

type OrderScheduledEventUpdateRequest struct {
	EndDate    string      `json:"endDate"`
	StartDate  string      `json:"startDate"`
	Assignment *Assignment `json:"assignment"`
}

type ServiceRequestAction struct {
	Remark string `json:"remark"`
}

func (v *QAService) recordPOSTActionHandler(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)

	if actionName == "approve" {
		v.handleQAApproved(ctx)
		return

	} else if actionName == "reject" {
		v.handleQARejected(ctx)
		return

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}
}

func (v *QAService) handleQAApproved(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	err, qaObject := Get(dbConnection, QABatchTable, intRecordId)
	if err == nil {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(qaObject.ObjectInfo, &objectFields)
		objectFields["qualityStatus"] = 2
		var lastUpdatedTime = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		objectFields["lastUpdatedAt"] = lastUpdatedTime
		objectFields["lastUpdatedBy"] = userId
		objectFields["canApprove"] = false
		objectFields["canReject"] = false
		existingActionRemarks := objectFields["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "APPROVED BY QA",
			UserId:        userId,
			Remarks:       "Great, the batch has been approved for further processing",
			ProcessedTime: getTimeDifference(util.InterfaceToString(objectFields["createdAt"])),
		})
		objectFields["actionRemarks"] = existingActionRemarks

		serialisedData, _ := json.Marshal(objectFields)

		var updatingObject = make(map[string]interface{})
		updatingObject["object_info"] = serialisedData

		Update(dbConnection, QABatchTable, intRecordId, updatingObject)

		// now create the traceablity records under traceablity module
		traceabilityInterface := common.GetService("traceability_module").ServiceInterface.(common.TraceabilityInterface)
		var mouldBatchId = util.InterfaceToString(objectFields["mouldBatchId"])
		v.BaseService.Logger.Info("creating the traceable records for ", zap.String("mould_batch_id", mouldBatchId))
		var qaStartTime = util.InterfaceToString(objectFields["createdAt"])
		traceabilityInterface.CreateTraceabilityResource(mouldBatchId, qaStartTime, lastUpdatedTime, userId)

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})
		return
	}
	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Internal System Error",
			Description: "The system not allows to make the QA batch to approve",
		})
	return
}

func (v *QAService) handleQARejected(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	var rejectionRemarkFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&rejectionRemarkFields, binding.JSON); err != nil {
		v.BaseService.Logger.Error("invalid payload on handle head of department reject request", zap.String("error", err.Error()))
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	rejectionRemarks := util.InterfaceToString(rejectionRemarkFields["remark"])
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	err, qaObject := Get(dbConnection, QABatchTable, intRecordId)
	if err == nil {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(qaObject.ObjectInfo, &objectFields)
		objectFields["qualityStatus"] = 3
		objectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		objectFields["lastUpdatedBy"] = userId
		objectFields["canApprove"] = false
		objectFields["canReject"] = false
		existingActionRemarks := objectFields["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "REJECTED BY QA",
			UserId:        userId,
			Remarks:       "Sorry, this batch was not approved to move forward [" + rejectionRemarks + "]",
			ProcessedTime: getTimeDifference(util.InterfaceToString(objectFields["createdAt"])),
		})
		objectFields["actionRemarks"] = existingActionRemarks

		serialisedData, _ := json.Marshal(objectFields)

		var updatingObject = make(map[string]interface{})
		updatingObject["object_info"] = serialisedData

		Update(dbConnection, QABatchTable, intRecordId, updatingObject)

		traceabilityInterface := common.GetService("traceability_module").ServiceInterface.(common.TraceabilityInterface)
		var mouldBatchId = util.InterfaceToString(objectFields["mouldBatchId"])
		var lastUpdatedTime = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		v.BaseService.Logger.Info("creating the traceable records for ", zap.String("mould_batch_id", mouldBatchId))
		var qaStartTime = util.InterfaceToString(objectFields["createdAt"])
		traceabilityInterface.CreateTraceabilityResource(mouldBatchId, qaStartTime, lastUpdatedTime, userId)

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})
		return
	} else {
		v.BaseService.Logger.Error("can not able to get the qa batch record", zap.String("error", err.Error()))
	}
	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Internal System Error",
			Description: "The system not allows to make the QA batch to approve",
		})
	return
}
