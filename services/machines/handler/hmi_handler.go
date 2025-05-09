package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type HMIRequest struct {
	EventId          int     `json:"eventId,omitempty"`
	Remark           string  `json:"remark,omitempty"`
	MachineId        int     `json:"machineId,omitempty"`
	HMIStatus        string  `json:"hmiStatus,omitempty"`
	CreatedAt        string  `json:"createdAt,omitempty"`
	CreatedBy        int     `json:"createdBy,omitempty"`
	OperatorId       int     `json:"operatorId,omitempty"`
	ReasonId         *int    `json:"reasonId,omitempty"`
	SetupTime        *string `json:"setupTime,omitempty"`
	RejectedQuantity *int    `json:"rejectedQuantity,omitempty"`
}

func (v *MachineService) handleNewHMIResource(ctx *gin.Context) int {

	projectId := util.GetProjectId(ctx)
	componentName := util.GetComponentName(ctx)
	userId := common.GetUserId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	hmiRequest := HMIRequest{}
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return 0
	}
	serializedObject, _ := json.Marshal(createRequest)
	initializedObject := common.InitMetaInfoFromSerializedObject(serializedObject, ctx)
	json.Unmarshal(initializedObject, &hmiRequest)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	if hmiRequest.EventId == 0 {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "No Scheduler Event Found",
				Description: "There is no scheduler event is assigned to this machine. Please schedule an event first before do any operations",
			})
		return 0
	}
	//Info validation based on machine status and time line
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, machineGeneralObject := Get(dbConnection, MachineMasterTable, hmiRequest.MachineId)

	if err != nil {
		response.DispatchDetailedError(ctx, InvalidMachineId,
			&response.DetailedError{
				Header:      "Invalid Machine ID",
				Description: "HMI operation is requesting with invalid machine ID, This may happen due to invalid machine master detail loaded or request is manipulated by external systems",
			})
		return 0
	}

	if isEventAborted(v, dbConnection, hmiRequest.EventId) {
		response.DispatchDetailedError(ctx, InvalidMachineId,
			&response.DetailedError{
				Header:      "Invalid Event ID",
				Description: "This event is aborted. Can't do any further operations",
			})
		return 0
	}

	err, scheduledOrder := productionOrderInterface.GetScheduledOrderInfo(projectId, hmiRequest.EventId)
	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "Invalid Event ID",
				Description: "No scheduler event is attached to given event Id, this might be due to event id either deleted or not created yet",
			})
		return 0
	}

	var scheduledOrderInfo map[string]interface{}
	_ = json.Unmarshal(scheduledOrder.ObjectInfo, &scheduledOrderInfo)

	var machineMasterInfo MachineMasterInfo
	err = json.Unmarshal(machineGeneralObject.ObjectInfo, &machineMasterInfo)

	if err != nil {
		response.DispatchDetailedError(ctx, DecodingFailed,
			&response.DetailedError{
				Header:      "Invalid Object",
				Description: "Internal system error during object decoding, some object fields have wrong data types",
			})
		return 0
	}

	v.BaseService.Logger.Info("machine master info", zap.Any("master_info", machineMasterInfo))

	//Insert new hmi info
	if hmiRequest.HMIStatus == "start" {
		completeScheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)

		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == completeScheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Can't start scheduled order was completed",
				})
			return 0
		}

		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Machine is not Live",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		scheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceThree)

		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == scheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Can't start scheduled order before release",
				})
			return 0
		}

		if util.InterfaceToInt(scheduledOrderInfo["mouldId"]) < 1 {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Please set the mould parameter",
				})
			return 0
		}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceFive)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "started"
	}
	if hmiRequest.HMIStatus == "stop" {
		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		// when we stop the HMI, we can complete the order
		updatingData["canComplete"] = true
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "stopped"

		// if this event has the mould batch resource ID created, then use that
		if mouldBatchResourceId, ok := scheduledOrderInfo["mouldBatchResourceId"]; ok {
			v.BaseService.Logger.Info("updating the mould batch, and generating the label", zap.Any("resource_id", mouldBatchResourceId))
			// if it is available, take that update the end time
			batchManagementInterface := common.GetService("batch_management_module").ServiceInterface.(common.BatchManagementInterface)
			var resourceId = util.InterfaceToInt(mouldBatchResourceId)
			batchManagementInterface.GenerateMouldBatchLabel(projectId, resourceId)

			generalObject := batchManagementInterface.GetMouldBatch(resourceId)
			var mouldBatchFields = make(map[string]interface{})
			json.Unmarshal(generalObject.ObjectInfo, &mouldBatchFields)
			var mouldBatchId = util.InterfaceToString(mouldBatchFields["mouldBatchId"])
			var rawMaterialBatchId = util.InterfaceToInt(mouldBatchFields["rawMaterialId"])

			err = composeQaBatch(userId, mouldBatchId, rawMaterialBatchId)

			if err != nil {
				v.BaseService.Logger.Error("error creating qa batch", zap.Any("error", err.Error()))
			}
		}
	}

	if hmiRequest.HMIStatus == "abort" {
		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Invalid Machine Status",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceEight)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "aborted"
	}

	if hmiRequest.RejectedQuantity != nil {
		if machineMasterInfo.MachineConnectStatus != machineConnectStatusLive {
			response.DispatchDetailedError(ctx, InvalidMachineStatus,
				&response.DetailedError{
					Header:      "Invalid Machine Status",
					Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
				})
			return 0
		}

		completeScheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == completeScheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't add reject quantity",
					Description: "Scheduled order was completed",
				})
			return 0
		}

		// frontend sending reject quantity now
		// check the hmi started or not
		err, isMachineStarted := isMachineHMIStarted(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, MachineHMITable)

		// totalProduction := getCurrentProductionOrder(ms, dbConnection, machineTimelineInfo.ProductionOrder)
		// totalRejectedQty := getRejectQtyByProductionName(ms, dbConnection, machineTimelineInfo.ProductionOrder)

		var productionOrderInfo map[string]interface{}

		err, productionOrder := productionOrderInterface.GetMachineProductionOrderInfo(projectId, util.InterfaceToInt(scheduledOrderInfo["eventSourceId"]), hmiRequest.MachineId)
		_ = json.Unmarshal(productionOrder.ObjectInfo, &productionOrderInfo)

		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Get Production Order",
					Description: "Invalid Machine Id or Production Id",
				})
			return 0
		}

		if !IsValidRejectCount(v, dbConnection, *hmiRequest.RejectedQuantity, scheduledOrderInfo, productionOrderInfo) {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "Rejected quantity is more than production order",
				})
			return 0
		}

		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "HMI is not started yet, you can not perform this operation without starting the HMI",
				})
			return 0
		}
	}
	if hmiRequest.ReasonId != nil {
		// frontend sending some stop reason
		// we need to update the event status to production stopped

		//If we have reason id in hmi, then we should add hmi status as stopped
		hmiRequest.HMIStatus = "stopped"

		err, isMachineStarted := isMachineHMIStarted(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, MachineHMITable)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			//response.DispatchDetailedError(ctx, DecodingFailed,
			//	&response.DetailedError{
			//		Header:      "HMI Error",
			//		Description: "HMI is not started yet, you can not perform this operation without starting the HMI",
			//	})
			//return 0
		}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
	}

	serializedDate, _ := json.Marshal(hmiRequest)
	object := component.GeneralObject{
		ObjectInfo: serializedDate,
	}
	err, createdRecordId := Create(dbConnection, targetTable, object)
	if err != nil {
		response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
			&response.DetailedError{
				Header:      "Error In Creating HMI Results",
				Description: "Internal system error during create [" + err.Error() + "]",
			})
		return 0
	}
	return createdRecordId
}

func composeQaBatch(userId int, mouldBatchId string, rawMaterialBatchId int) error {
	qaInterface := common.GetService("qa").ServiceInterface.(common.QAInterface)
	var qaObjectFields = make(map[string]interface{})
	qaObjectFields["objectStatus"] = component.ObjectStatusActive
	var createdAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	qaObjectFields["createdAt"] = createdAt
	qaObjectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	qaObjectFields["createdBy"] = userId
	qaObjectFields["lastUpdatedBy"] = userId
	qaObjectFields["mouldBatchId"] = mouldBatchId
	qaObjectFields["rawMaterialId"] = rawMaterialBatchId
	qaObjectFields["qualityStatus"] = 1
	var actionRemarks = make([]interface{}, 0)
	actionRemarks = append(actionRemarks, ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "QA EVALUATION CREATED",
		UserId:        userId,
		Remarks:       "Great, Evaluation for quality assurance is created",
		ProcessedTime: getTimeDifference(util.InterfaceToString(createdAt)),
	})
	qaObjectFields["actionRemarks"] = actionRemarks

	serialisedData, _ := json.Marshal(qaObjectFields)
	err := qaInterface.CreateQAResource(serialisedData)

	return err
}

func (v *MachineService) handleNewAssemblyHMIResource(ctx *gin.Context) int {

	projectId := util.GetProjectId(ctx)
	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	hmiRequest := HMIRequest{}
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return 0
	}
	serializedObject, _ := json.Marshal(createRequest)
	initializedObject := common.InitMetaInfoFromSerializedObject(serializedObject, ctx)
	json.Unmarshal(initializedObject, &hmiRequest)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	if hmiRequest.EventId == 0 {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "No Scheduler Event Found",
				Description: "There is no scheduler event is assigned to this machine. Please schedule an event first before do any operations",
			})
		return 0
	}
	//Info validation based on machine status and time line
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, machineGeneralObject := Get(dbConnection, AssemblyMachineMasterTable, hmiRequest.MachineId)

	if err != nil {
		response.DispatchDetailedError(ctx, InvalidMachineId,
			&response.DetailedError{
				Header:      "Invalid Machine ID",
				Description: "HMI operation is requesting with invalid machine ID, This may happen due to invalid machine master detail loaded or request is manipulated by external systems",
			})
		return 0
	}

	if isAssemblyEventAborted(v, dbConnection, hmiRequest.EventId) {
		response.DispatchDetailedError(ctx, InvalidMachineId,
			&response.DetailedError{
				Header:      "Invalid Event ID",
				Description: "This event is aborted. Can't do any further operations",
			})
		return 0
	}

	err, scheduledOrder := productionOrderInterface.GetAssemblyScheduledOrderInfo(projectId, hmiRequest.EventId)
	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "Invalid Event ID",
				Description: "No scheduler event is attached to given event Id, this might be due to event id either deleted or not created yet",
			})
		return 0
	}

	var scheduledOrderInfo map[string]interface{}
	_ = json.Unmarshal(scheduledOrder.ObjectInfo, &scheduledOrderInfo)

	var machineMasterInfo MachineMasterInfo
	err = json.Unmarshal(machineGeneralObject.ObjectInfo, &machineMasterInfo)

	if err != nil {
		response.DispatchDetailedError(ctx, DecodingFailed,
			&response.DetailedError{
				Header:      "Invalid Object",
				Description: "Internal system error during object decoding, some object fields have wrong data types",
			})
		return 0
	}

	v.BaseService.Logger.Info("machine master info", zap.Any("master_info", machineMasterInfo))

	//Insert new hmi info
	if hmiRequest.HMIStatus == "start" {
		completeScheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)

		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == completeScheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Can't start scheduled order was completed",
				})
			return 0
		}

		//TODO this is removed, so that assembly machine can start without live status (Requirement from FuYu)
		/*
			if machineMasterInfo.MachineConnectStatus != "Live" {
				response.DispatchDetailedError(ctx, InvalidMachineStatus,
					&response.DetailedError{
						Header:      "Invalid Machine Status",
						Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
					})
				return 0
			}

		*/

		scheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceThree)

		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == scheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Can't start scheduled order before release",
				})
			return 0
		}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceFive)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "started"
	}
	if hmiRequest.HMIStatus == "stop" {
		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "stopped"
	}

	if hmiRequest.HMIStatus == "abort" {
		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Machine is not Live",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceEight)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "aborted"
	}

	if hmiRequest.RejectedQuantity != nil {
		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Invalid Machine Status",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		// frontend sending reject quantity now
		// check the hmi started or not
		err, isMachineStarted := isMachineHMIStarted(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, AssemblyMachineHmiTable)

		// totalProduction := getCurrentProductionOrder(ms, dbConnection, machineTimelineInfo.ProductionOrder)
		// totalRejectedQty := getRejectQtyByProductionName(ms, dbConnection, machineTimelineInfo.ProductionOrder)

		var productionOrderInfo map[string]interface{}

		err, productionOrder := productionOrderInterface.GetAssemblyProductionOrderInfo(projectId, util.InterfaceToInt(scheduledOrderInfo["eventSourceId"]), hmiRequest.MachineId)
		_ = json.Unmarshal(productionOrder.ObjectInfo, &productionOrderInfo)

		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Get Production Order",
					Description: "Invalid Machine Id or Production Id",
				})
			return 0
		}

		if !IsValidRejectCount(v, dbConnection, *hmiRequest.RejectedQuantity, scheduledOrderInfo, productionOrderInfo) {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "Rejected quantity is more than production order",
				})
			return 0
		}

		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "HMI is not started yet, you can not perform this operation without starting the HMI",
				})
			return 0
		}
	}
	if hmiRequest.ReasonId != nil {
		// frontend sending some stop reason
		// we need to update the event status to production stopped

		//If we have reason id in hmi, then we should add hmi status as stopped
		hmiRequest.HMIStatus = "stopped"

		err, isMachineStarted := isMachineHMIStarted(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, AssemblyMachineHmiTable)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			//response.DispatchDetailedError(ctx, DecodingFailed,
			//	&response.DetailedError{
			//		Header:      "HMI Error",
			//		Description: "HMI is not started yet, you can not perform this operation without starting the HMI",
			//	})
			//return 0
		}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
	}

	serializedDate, _ := json.Marshal(hmiRequest)
	object := component.GeneralObject{
		ObjectInfo: serializedDate,
	}
	err, createdRecordId := Create(dbConnection, targetTable, object)
	if err != nil {
		response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
			&response.DetailedError{
				Header:      "Error In Creating HMI Results",
				Description: "Internal system error during create [" + err.Error() + "]",
			})
		return 0
	}
	return createdRecordId
}

func (v *MachineService) handleNewToolingHMIResource(ctx *gin.Context) int {

	projectId := util.GetProjectId(ctx)
	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	hmiRequest := HMIRequest{}
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return 0
	}
	serializedObject, _ := json.Marshal(createRequest)
	initializedObject := common.InitMetaInfoFromSerializedObject(serializedObject, ctx)
	json.Unmarshal(initializedObject, &hmiRequest)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	if hmiRequest.EventId == 0 {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "No Scheduler Event Found",
				Description: "There is no scheduler event is assigned to this machine. Please schedule an event first before do any operations",
			})
		return 0
	}
	//Info validation based on machine status and time line
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, machineGeneralObject := Get(dbConnection, ToolingMachineMasterTable, hmiRequest.MachineId)

	if err != nil {
		response.DispatchDetailedError(ctx, InvalidMachineId,
			&response.DetailedError{
				Header:      "Invalid Machine ID",
				Description: "HMI operation is requesting with invalid machine ID, This may happen due to invalid machine master detail loaded or request is manipulated by external systems",
			})
		return 0
	}

	if isToolingEventAborted(v, dbConnection, hmiRequest.EventId) {
		response.DispatchDetailedError(ctx, InvalidMachineId,
			&response.DetailedError{
				Header:      "Invalid Event ID",
				Description: "This event is aborted. Can't do any further operations",
			})
		return 0
	}

	err, scheduledOrder := productionOrderInterface.GetToolingScheduledOrderInfo(projectId, hmiRequest.EventId)
	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEventId,
			&response.DetailedError{
				Header:      "Invalid Event ID",
				Description: "No scheduler event is attached to given event Id, this might be due to event id either deleted or not created yet",
			})
		return 0
	}

	var scheduledOrderInfo map[string]interface{}
	_ = json.Unmarshal(scheduledOrder.ObjectInfo, &scheduledOrderInfo)

	var machineMasterInfo ToolingMachineMasterInfo
	err = json.Unmarshal(machineGeneralObject.ObjectInfo, &machineMasterInfo)

	if err != nil {
		response.DispatchDetailedError(ctx, DecodingFailed,
			&response.DetailedError{
				Header:      "Invalid Object",
				Description: "Internal system error during object decoding, some object fields have wrong data types",
			})
		return 0
	}

	v.BaseService.Logger.Info("machine master info", zap.Any("master_info", machineMasterInfo))

	//Insert new hmi info
	if hmiRequest.HMIStatus == "start" {
		completeScheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)

		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == completeScheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Can't start scheduled order was completed",
				})
			return 0
		}

		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Machine is not Live",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		scheduledStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceThree)

		if util.InterfaceToInt(scheduledOrderInfo["eventStatus"]) == scheduledStatusId {
			response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
				&response.DetailedError{
					Header:      "Can't start HMI",
					Description: "Can't start scheduled order before release",
				})
			return 0
		}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceFive)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateToolingScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "started"
	}
	if hmiRequest.HMIStatus == "stop" {
		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateToolingScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "stopped"
	}

	if hmiRequest.HMIStatus == "abort" {
		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Invalid Machine Status",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceEight)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateToolingScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
		hmiRequest.HMIStatus = "aborted"
	}

	if hmiRequest.RejectedQuantity != nil {
		//if machineMasterInfo.MachineConnectStatus != "Live" {
		//	response.DispatchDetailedError(ctx, InvalidMachineStatus,
		//		&response.DetailedError{
		//			Header:      "Invalid Machine Status",
		//			Description: "Machine is not Live to accept any operations, Looks like machine is currently being stopped physically or error connecting real-time gateway",
		//		})
		//	return 0
		//}

		// frontend sending reject quantity now
		// check the hmi started or not
		err, isMachineStarted := isMachineHMIStarted(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, AssemblyMachineHmiTable)

		// totalProduction := getCurrentProductionOrder(ms, dbConnection, machineTimelineInfo.ProductionOrder)
		// totalRejectedQty := getRejectQtyByProductionName(ms, dbConnection, machineTimelineInfo.ProductionOrder)

		var productionOrderInfo map[string]interface{}

		err, productionOrder := productionOrderInterface.GetToolingProductionOrderInfo(projectId, util.InterfaceToInt(scheduledOrderInfo["eventSourceId"]), hmiRequest.MachineId)
		_ = json.Unmarshal(productionOrder.ObjectInfo, &productionOrderInfo)

		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Get Production Order",
					Description: "Invalid Machine Id or Production Id",
				})
			return 0
		}

		if !IsValidRejectCount(v, dbConnection, *hmiRequest.RejectedQuantity, scheduledOrderInfo, productionOrderInfo) {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "Rejected quantity is more than production order",
				})
			return 0
		}

		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "HMI is not started yet, you can not perform this operation without starting the HMI",
				})
			return 0
		}
	}
	if hmiRequest.ReasonId != nil {
		// frontend sending some stop reason
		// we need to update the event status to production stopped

		//If we have reason id in hmi, then we should add hmi status as stopped
		hmiRequest.HMIStatus = "stopped"

		err, isMachineStarted := isMachineHMIStarted(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, ToolingMachineHmiTable)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			//response.DispatchDetailedError(ctx, DecodingFailed,
			//	&response.DetailedError{
			//		Header:      "HMI Error",
			//		Description: "HMI is not started yet, you can not perform this operation without starting the HMI",
			//	})
			//return 0
		}

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSix)

		var updatingData = make(map[string]interface{})
		updatingData["eventStatus"] = orderStatusId
		serializedObject, _ = json.Marshal(updatingData)
		err = productionOrderInterface.UpdateToolingScheduledOrderFields(projectId, hmiRequest.EventId, serializedObject)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Update Event Status",
					Description: "Internal system error during update [" + err.Error() + "]",
				})
			return 0
		}
	}

	if hmiRequest.SetupTime != nil {

		err, isMachineStarted := canSetupTimeAdded(dbConnection, hmiRequest.MachineId, hmiRequest.EventId, ToolingMachineHmiTable)
		if err != nil {
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "Error Getting Event",
					Description: "Internal system error during fetching events [" + err.Error() + "]",
				})
			return 0
		}
		if !isMachineStarted {
			v.BaseService.Logger.Info("HMI is not started to machine", zap.Any("request", string(initializedObject)))
			response.DispatchDetailedError(ctx, DecodingFailed,
				&response.DetailedError{
					Header:      "HMI Error",
					Description: "HMI is not stopped yet, you can not perform this operation without stopping the HMI",
				})
			return 0
		}
	}

	serializedDate, _ := json.Marshal(hmiRequest)
	object := component.GeneralObject{
		ObjectInfo: serializedDate,
	}
	err, createdRecordId := Create(dbConnection, targetTable, object)
	if err != nil {
		response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
			&response.DetailedError{
				Header:      "Error In Creating HMI Results",
				Description: "Internal system error during create [" + err.Error() + "]",
			})
		return 0
	}
	return createdRecordId
}

func canSetupTimeAdded(dbConnection *gorm.DB, machineId, eventId int, targetTable string) (error, bool) {
	conditionString := "(JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.eventId\"))) =" + strconv.Itoa(eventId) + " AND (JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.machineId\"))) = " + strconv.Itoa(machineId)
	listOfHMIObjects, err := GetConditionalObjects(dbConnection, targetTable, conditionString)
	if err != nil {
		return err, true
	}
	if len(*listOfHMIObjects) == 0 {
		return nil, true
	}

	foundHmiStatus := false
	isHiStarted := false

	for index := len(*listOfHMIObjects) - 1; index >= 0; index-- {

		if foundHmiStatus {
			break
		}
		lastObjectHmi := (*listOfHMIObjects)[index]
		hmiInfo := make(map[string]interface{})
		json.Unmarshal(lastObjectHmi.ObjectInfo, &hmiInfo)

		hmiStatus := util.InterfaceToString(hmiInfo["hmiStatus"])

		if hmiStatus != "" {
			foundHmiStatus = true
			if hmiStatus == "stopped" {
				isHiStarted = true
			}
		}

	}

	return nil, isHiStarted
}

func isMachineHMIStarted(dbConnection *gorm.DB, machineId, eventId int, targetTable string) (error, bool) {
	conditionString := "(JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.eventId\"))) =" + strconv.Itoa(eventId) + " AND (JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.machineId\"))) = " + strconv.Itoa(machineId) + " AND (JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.hmiStatus\"))) = \"started\""
	listOfHMIObjects, err := GetConditionalObjects(dbConnection, targetTable, conditionString)

	if err != nil {
		return err, false
	}
	if len(*listOfHMIObjects) == 0 {
		return nil, false
	}
	return nil, true
}

func (v *MachineService) getHMIInfoResponse(projectId string, machineId int, operatorId int, machineMasterInfo map[string]interface{}, dbConnection *gorm.DB, targetTable string) datatypes.JSON {
	var settingTable string
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	locationList := getLocationDetails(dbConnection)
	var machineParamId = -1
	if targetTable == AssemblyMachineHmiTable {
		settingTable = AssemblyMachineHmiSettingTable
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

		operatorInfo := authService.GetUserInfoById(operatorId)
		v.BaseService.Logger.Info("operator info", zap.Any("operator_id", operatorId), zap.Any("operator_info", operatorInfo))

		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))

		err, scheduledEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		eventId := scheduledEventObject.Id
		scheduledOrderEvent := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
		startScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.StartDate)
		endScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.EndDate)
		v.BaseService.Logger.Info("current timeline event", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.eventId' =" + strconv.Itoa(eventId) + " order by  object_info->>'$.createdAt' desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, targetTable, hmiInfoConditionString)

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		var hmiStatus string
		//default is stopped, for example, if we don't have any hmi info (as it is the first time), then we need to have somthing to operate
		hmiStatus = "stopped"
		hmiId := 0

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, scheduledOrderEvent.EventStatus)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		if scheduledOrderEvent.EventStatus == orderStatusId || hmiStatus == "" {
			hmiStatus = "disabled"
		}

		stopList := v.getStopList(eventId, dbConnection, AssemblyMachineHmiTable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetAssemblyProductionOrderInfo(projectId, scheduledOrderEvent.EventSourceId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		productionOrderInfo := GetProductionOrderInfo(productionObject.ObjectInfo)
		manufacturingInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
		err, partObject := manufacturingInterface.GetAssemblyPartInfo(projectId, productionOrderInfo.PartNumber)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		partInfo := GetPartInfo(partObject.ObjectInfo)

		// Get warning message from stats table
		warningMessage := getWarningMessage(dbConnection, eventId, machineId, AssemblyMachineStatisticsTable)

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, AssemblyMachineHmiTable)
		var isItAlreadyStarted = isMachineAlreadyStarted(dbConnection, scheduledOrderEvent.MachineId, eventId)
		return buildHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, warningMessage, component.RecordInfo{}, "", "", hmiId, startedBy, stoppedBy, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList, isItAlreadyStarted, machineParamId)
	} else if targetTable == ToolingMachineHmiTable {
		settingTable = ToolingMachineHmiSettingTable
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

		operatorInfo := authService.GetUserInfoById(operatorId)
		v.BaseService.Logger.Info("operator info", zap.Any("operator_id", operatorId), zap.Any("operator_info", operatorInfo))

		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))

		err, scheduledEventObject := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineId)
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName)
		}

		eventId := scheduledEventObject.Id
		scheduledOrderEvent := make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledOrderEvent)

		startScheduledDateTime := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["startDate"]))
		endScheduledDateTime := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["endDate"]))
		eventSourceId := util.InterfaceToInt(scheduledOrderEvent["eventSourceId"])
		v.BaseService.Logger.Info("current timeline event", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.eventId' =" + strconv.Itoa(eventId) + " order by  object_info->>'$.createdAt' desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, targetTable, hmiInfoConditionString)

		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		var hmiStatus string
		//default is stopped, for example, if we don't have any hmi info (as it is the first time), then we need to have somthing to operate
		hmiStatus = "stopped"
		hmiId := 0

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		setupTimeList := getSetupTimeList(listOfHMIInfo)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, util.InterfaceToInt(scheduledOrderEvent["eventStatus"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		if util.InterfaceToInt(scheduledOrderEvent["eventStatus"]) == orderStatusId || hmiStatus == "" {
			hmiStatus = "disabled"
		}

		stopList := v.getStopList(eventId, dbConnection, ToolingMachineHmiTable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetToolingProductionOrderInfo(projectId, util.InterfaceToInt(scheduledOrderEvent["eventSourceId"]), machineId)
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}

		productionOrderInfo := make(map[string]interface{})
		json.Unmarshal(productionObject.ObjectInfo, &productionOrderInfo)

		//totalDuration := util.InterfaceToInt(productionOrderInfo["day"])*24 + util.InterfaceToInt(productionOrderInfo["hour"]) + util.InterfaceToInt(productionOrderInfo["minute"])/60

		err, partObject := productionOrderInterface.GetToolingPartById(projectId, util.InterfaceToInt(scheduledOrderEvent["partId"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		partInfo := make(map[string]interface{})
		json.Unmarshal(partObject.ObjectInfo, &partInfo)

		//totalDuration := util.InterfaceToFloat(partInfo["day"])*24 + util.InterfaceToFloat(partInfo["hour"]) + util.InterfaceToFloat(partInfo["minute"])/60
		totalDuration := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["endDate"])).DateTimeEpoch - util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["startDate"])).DateTimeEpoch

		// Convert seconds to a duration
		duration := time.Second * time.Duration(totalDuration)

		// Format the duration as a string representing hours
		hoursString := fmt.Sprintf("%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60)
		durationString := hoursString + " hours"

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, ToolingMachineHmiTable)

		_, bomGeneralObject := productionOrderInterface.GetBomInfo(projectId, eventSourceId)
		var bomObjectInfo = make(map[string]interface{})
		json.Unmarshal(bomGeneralObject.ObjectInfo, &bomObjectInfo)
		bomName := util.InterfaceToString(bomObjectInfo["name"])

		return buildToolingHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, "", durationString, bomName, hmiId, startedBy, stoppedBy, setupTimeList, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList)
	} else {
		settingTable = MachineHMISettingSettingTable

		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

		operatorInfo := authService.GetUserInfoById(operatorId)
		v.BaseService.Logger.Info("operator info", zap.Any("operator_id", operatorId), zap.Any("operator_info", operatorInfo))
		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))

		err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineId)

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		eventId := scheduledEventObject.Id

		scheduledOrderEvent := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
		startScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.StartDate)
		endScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.EndDate)
		v.BaseService.Logger.Info("current timeline event", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " object_info ->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.eventId' =" + strconv.Itoa(eventId) + " order by  object_info->>'$.createdAt' desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, targetTable, hmiInfoConditionString)

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		var hmiStatus string
		//default is stopped, for example, if we don't have any hmi info (as it is the first time), then we need to have somthing to operate
		hmiStatus = "stopped"
		hmiId := 0
		var mouldId int

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first

		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, scheduledOrderEvent.EventStatus)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		if scheduledOrderEvent.EventStatus == orderStatusId || hmiStatus == "" {
			hmiStatus = "disabled"
		}

		stopList := v.getStopList(eventId, dbConnection, MachineHMITable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, scheduledOrderEvent.EventSourceId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		productionOrderInfo := GetProductionOrderInfo(productionObject.ObjectInfo)
		mouldModuleInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		err, partObject := mouldModuleInterface.GetPartInfo(projectId, productionOrderInfo.PartNumber)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		partInfo := GetPartInfo(partObject.ObjectInfo)

		// Get mould ids from part number
		partId := strconv.Itoa(productionOrderInfo.PartNumber)
		_, mouldList := mouldModuleInterface.GetMouldsByPartNo(projectId, partId)
		recordInfo := component.RecordInfo{}
		var dropDownArray []component.OrderedData
		index := 0
		mouldId = scheduledOrderEvent.MouldId

		mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		_, mouldGeneralObj := mouldService.GetMouldInfoById(projectId, scheduledOrderEvent.MouldId)
		mouldInfo := make(map[string]interface{})
		json.Unmarshal(mouldGeneralObj.ObjectInfo, &mouldInfo)

		var mouldDescription string
		mouldToolNo := util.InterfaceToString(mouldInfo["toolNo"])
		if description, ok := mouldInfo["description"]; ok {
			mouldDescription = util.InterfaceToString(description)
		}

		for _, mould := range mouldList {
			id := mould.Id
			mouldObjectInfo := make(map[string]interface{})
			json.Unmarshal(mould.ObjectInfo, &mouldObjectInfo)
			dropdownValue := util.InterfaceToString(mouldObjectInfo["toolNo"])
			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:    id,
				Value: dropdownValue,
			})
			if index == 0 {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}

			if mouldId == id {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}
			index = index + 1
		}
		recordInfo.Data = dropDownArray

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, MachineHMITable)

		// Get warning message from stats table
		warningMessage := getWarningMessage(dbConnection, eventId, machineId, MachineStatisticsTable)
		var isItAlreadyStarted = isMachineAlreadyStarted(dbConnection, scheduledOrderEvent.MachineId, eventId)
		machineParamId = v.GetMouldTestMachineParam(machineId, eventId)
		v.BaseService.Logger.Info("setting the machine param ID for current one ", zap.Int("machine_param_id", machineParamId))

		return buildHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, warningMessage, recordInfo, mouldDescription, mouldToolNo, hmiId, startedBy, stoppedBy, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList, isItAlreadyStarted, machineParamId)

	}

}

func getLocationDetails(dbconnection *gorm.DB) component.RecordInfo {
	recordInfo := component.RecordInfo{}
	var dropDownArray []component.OrderedData

	factoryInterface := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
	generalObjects := factoryInterface.GetBuildingInfo()

	if generalObjects != nil {
		for index, buildingObjects := range *generalObjects {
			resultInfo := make(map[string]interface{})
			json.Unmarshal(buildingObjects.ObjectInfo, &resultInfo)

			id := buildingObjects.Id

			dropdownValue := util.InterfaceToString(resultInfo["name"])
			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:    id,
				Value: dropdownValue,
			})
			if index == 0 {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}

			index = index + 1
		}
		recordInfo.Data = dropDownArray
	}
	return recordInfo
}

func getLastHmiOperateUser(listOfHMIInfo *[]component.GeneralObject) (string, string, string, string) {
	// Hmi info in decending order
	lastStartUser := "-"
	lastStopUser := "-"

	lastStartAvatarUser := ""
	lastStopAvatarUser := ""

	isLastStartPersonFound := false
	isLastStopPersonFound := false

	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	for _, hmiObject := range *listOfHMIInfo {
		machineHMIInfo := make(map[string]interface{})
		json.Unmarshal(hmiObject.ObjectInfo, &machineHMIInfo)

		if isLastStartPersonFound && isLastStopPersonFound {
			break
		}

		if util.InterfaceToString(machineHMIInfo["hmiStatus"]) == "started" {
			isLastStartPersonFound = true
			userInfo := authService.GetUserInfoById(util.InterfaceToInt(machineHMIInfo["createdBy"]))
			if userInfo.FullName != "" {
				lastStartUser = userInfo.FullName

				upstreamResponse, err := util.UpstreamGet(userInfo.AvatarUrl + "/action/meta_info")
				if err == nil {
					var upstreamResponseFields = make(map[string]interface{})
					json.Unmarshal(upstreamResponse, &upstreamResponseFields)
					lastStartAvatarUser = util.InterfaceToString(upstreamResponseFields["url"])
				}

			}

		}

		if util.InterfaceToString(machineHMIInfo["hmiStatus"]) == "stopped" {
			isLastStopPersonFound = true
			userInfo := authService.GetUserInfoById(util.InterfaceToInt(machineHMIInfo["createdBy"]))
			if userInfo.FullName != "" {
				lastStopUser = userInfo.FullName

				upstreamResponse, err := util.UpstreamGet(userInfo.AvatarUrl + "/action/meta_info")

				if err == nil {
					var upstreamResponseFields = make(map[string]interface{})
					json.Unmarshal(upstreamResponse, &upstreamResponseFields)

					lastStopAvatarUser = util.InterfaceToString(upstreamResponseFields["url"])
				}
			}

		}
	}

	return lastStartUser, lastStopUser, lastStartAvatarUser, lastStopAvatarUser

}

func getWarningMessage(dbConnection *gorm.DB, eventId int, machineId int, statsTable string) string {
	warningMessage := ""
	statisticsQuery := "select * from " + statsTable + " where stats_info ->> '$.eventId' = " + strconv.Itoa(eventId) + " and machine_id=" + strconv.Itoa(machineId) + " order by ts desc limit 1"
	var machineStatics MachineStatistics
	var machineStatsInfo MachineStatisticsInfo

	dbConnection.Raw(statisticsQuery).Scan(&machineStatics)

	_ = json.Unmarshal(machineStatics.StatsInfo, &machineStatsInfo)

	listOfWarningMessage := machineStatsInfo.WarningMessage

	for _, msg := range listOfWarningMessage {
		warningMessage = msg
	}

	if warningMessage != "" {
		warningMessageList := strings.Split(warningMessage, " ")
		warningMessage = strings.Join(warningMessageList[:len(warningMessageList)-1], " ")
	}

	return warningMessage
}

func (v *MachineService) getNextHMIResponse(projectId string, eventId, machineId int, operatorId int, machineMasterInfo map[string]interface{}, dbConnection *gorm.DB, targetTable string) datatypes.JSON {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	operatorInfo := authService.GetUserInfoById(operatorId)
	var settingTable string
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	locationList := getLocationDetails(dbConnection)
	var machineParamId = -1
	if targetTable == AssemblyMachineHmiTable {
		settingTable = AssemblyMachineHmiSettingTable
		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		err, scheduledEventObject := productionOrderInterface.GetNextAssemblyMachineScheduledOrderEvent(projectId, machineId, eventId)
		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		if len((*scheduledEventObject)) < 1 {
			err, scheduledEventObject = productionOrderInterface.GetFirstAssemblyScheduledOrderEvent(projectId, machineId)
		}

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		scheduledOrderEvent := GetScheduledOrderEventInfo((*scheduledEventObject)[0].ObjectInfo)
		eventId = (*scheduledEventObject)[0].Id
		startScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.StartDate)
		endScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.EndDate)
		v.BaseService.Logger.Info("current timeline event:", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " JSON_EXTRACT(object_info, \"$.machineId\") =" + strconv.Itoa(machineId) + " AND JSON_EXTRACT(object_info, \"$.eventId\") =" + strconv.Itoa(eventId) + " order by  JSON_EXTRACT(object_info, \"$.createdAt\") desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, AssemblyMachineHmiTable, hmiInfoConditionString)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		var hmiStatus string
		var hmiId int

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, scheduledOrderEvent.EventStatus)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		stopList := v.getStopList(eventId, dbConnection, AssemblyMachineHmiTable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, scheduledOrderEvent.EventSourceId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		productionOrderInfo := GetProductionOrderInfo(productionObject.ObjectInfo)
		manufacturingInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
		err, partObject := manufacturingInterface.GetAssemblyPartInfo(projectId, productionOrderInfo.PartNumber)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		partInfo := GetPartInfo(partObject.ObjectInfo)

		//orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		//if scheduledOrderEvent.EventStatus == orderStatusId || hmiStatus == "" {
		//	err, currentEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineId)
		//	if err == nil && eventId != currentEventObject.Id {
		//		hmiStatus = "disabled"
		//	}
		//}
		_, currentEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineId)
		if eventId != currentEventObject.Id {
			hmiStatus = "disabled"
		}

		if hmiStatus == "" {
			hmiStatus = "stopped"
		}

		// Get mould ids from part number

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, AssemblyMachineHmiTable)

		// Get warning message from stats table
		warningMessage := getWarningMessage(dbConnection, eventId, machineId, AssemblyMachineStatisticsTable)
		var isItAlreadyStarted = isMachineAlreadyStarted(dbConnection, scheduledOrderEvent.MachineId, eventId)
		return buildHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, warningMessage, component.RecordInfo{}, "", "", hmiId, startedBy, stoppedBy, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList, isItAlreadyStarted, machineParamId)

	} else if targetTable == ToolingMachineHmiTable {
		settingTable = ToolingMachineHmiSettingTable
		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		err, scheduledEventObject := productionOrderInterface.GetNextToolingMachineScheduledOrderEvent(projectId, machineId, eventId)
		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName)
		}

		if len((*scheduledEventObject)) < 1 {
			err, scheduledEventObject = productionOrderInterface.GetFirstToolingScheduledOrderEvent(projectId, machineId)
		}

		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName)
		}

		scheduledOrderEvent := make(map[string]interface{})
		json.Unmarshal((*scheduledEventObject)[0].ObjectInfo, &scheduledOrderEvent)
		eventId = (*scheduledEventObject)[0].Id
		startScheduledDateTime := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["startDate"]))
		endScheduledDateTime := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["endDate"]))
		v.BaseService.Logger.Info("current timeline event:", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " JSON_EXTRACT(object_info, \"$.machineId\") =" + strconv.Itoa(machineId) + " AND JSON_EXTRACT(object_info, \"$.eventId\") =" + strconv.Itoa(eventId) + " order by  JSON_EXTRACT(object_info, \"$.createdAt\") desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, ToolingMachineHmiTable, hmiInfoConditionString)
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		var hmiStatus string
		var hmiId int

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		eventSourceId := util.InterfaceToInt(scheduledOrderEvent["eventSourceId"])
		rejectList := getRejectList(listOfHMIInfo, eventId)
		setupTimeList := getSetupTimeList(listOfHMIInfo)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, util.InterfaceToInt(scheduledOrderEvent["eventStatus"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		if util.InterfaceToInt(scheduledOrderEvent["eventStatus"]) == orderStatusId || hmiStatus == "" {
			hmiStatus = "disabled"
		}

		stopList := v.getStopList(eventId, dbConnection, ToolingMachineHmiTable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetToolingProductionOrderInfo(projectId, util.InterfaceToInt(scheduledOrderEvent["eventSourceId"]), machineId)
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}

		productionOrderInfo := make(map[string]interface{})
		json.Unmarshal(productionObject.ObjectInfo, &productionOrderInfo)

		//totalDuration := util.InterfaceToInt(productionOrderInfo["day"])*24 + util.InterfaceToInt(productionOrderInfo["hour"]) + util.InterfaceToInt(productionOrderInfo["minute"])/60

		err, partObject := productionOrderInterface.GetToolingPartById(projectId, util.InterfaceToInt(scheduledOrderEvent["partId"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		partInfo := make(map[string]interface{})
		json.Unmarshal(partObject.ObjectInfo, &partInfo)

		totalDuration := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["endDate"])).DateTimeEpoch - util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["startDate"])).DateTimeEpoch

		// Convert seconds to a duration
		duration := time.Second * time.Duration(totalDuration)

		// Format the duration as a string representing hours
		hoursString := fmt.Sprintf("%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60)
		durationString := hoursString + " hours"

		_, currentEventObject := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineId)
		if eventId != currentEventObject.Id {
			hmiStatus = "disabled"
		}

		if hmiStatus == "" {
			hmiStatus = "stopped"
		}

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, ToolingMachineHmiTable)

		_, bomGeneralObject := productionOrderInterface.GetBomInfo(projectId, eventSourceId)
		var bomObjectInfo = make(map[string]interface{})
		json.Unmarshal(bomGeneralObject.ObjectInfo, &bomObjectInfo)
		bomName := util.InterfaceToString(bomObjectInfo["name"])

		return buildToolingHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, "", durationString, bomName, hmiId, startedBy, stoppedBy, setupTimeList, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList)
	} else {
		settingTable = MachineHMISettingSettingTable

		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		err, scheduledEventObject := productionOrderInterface.GetNextMachineScheduledOrderEvent(projectId, machineId, eventId)
		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		if len((*scheduledEventObject)) < 1 {
			err, scheduledEventObject = productionOrderInterface.GetFirstScheduledOrderEvent(projectId, machineId)
		}

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		scheduledOrderEvent := GetScheduledOrderEventInfo((*scheduledEventObject)[0].ObjectInfo)
		eventId = (*scheduledEventObject)[0].Id
		startScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.StartDate)
		endScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.EndDate)
		v.BaseService.Logger.Info("current timeline event:", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " JSON_EXTRACT(object_info, \"$.machineId\") =" + strconv.Itoa(machineId) + " AND JSON_EXTRACT(object_info, \"$.eventId\") =" + strconv.Itoa(eventId) + " order by  JSON_EXTRACT(object_info, \"$.createdAt\") desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, MachineHMITable, hmiInfoConditionString)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		var hmiStatus string
		var mouldId int
		var hmiId int

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)
		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, scheduledOrderEvent.EventStatus)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		stopList := v.getStopList(eventId, dbConnection, MachineHMITable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, scheduledOrderEvent.EventSourceId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		productionOrderInfo := GetProductionOrderInfo(productionObject.ObjectInfo)
		mouldModuleInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		err, partObject := mouldModuleInterface.GetPartInfo(projectId, productionOrderInfo.PartNumber)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		partInfo := GetPartInfo(partObject.ObjectInfo)

		//orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		//if scheduledOrderEvent.EventStatus == orderStatusId || hmiStatus == "" {
		//	err, currentEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineId)
		//	fmt.Println("=======================================")
		//	fmt.Println(currentEventObject.ObjectInfo)
		//	if err == nil && eventId != currentEventObject.Id {
		//		hmiStatus = "disabled"
		//	} else {
		//		fmt.Println(hmiStatus)
		//		if hmiStatus == "" {
		//			hmiStatus = "stopped"
		//		}
		//	}
		//}

		_, currentEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineId)

		if eventId != currentEventObject.Id {
			hmiStatus = "disabled"
		}

		if hmiStatus == "" {
			hmiStatus = "stopped"
		}
		// Get mould ids from part number
		partId := strconv.Itoa(productionOrderInfo.PartNumber)
		_, mouldList := mouldModuleInterface.GetMouldsByPartNo(projectId, partId)
		recordInfo := component.RecordInfo{}
		var dropDownArray []component.OrderedData
		index := 0
		mouldId = scheduledOrderEvent.MouldId

		mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		_, mouldGeneralObj := mouldService.GetMouldInfoById(projectId, mouldId)
		mouldInfo := make(map[string]interface{})
		json.Unmarshal(mouldGeneralObj.ObjectInfo, &mouldInfo)

		var mouldDescription string
		mouldToolNo := util.InterfaceToString(mouldInfo["toolNo"])
		if description, ok := mouldInfo["description"]; ok {
			mouldDescription = util.InterfaceToString(description)
		}

		for _, mould := range mouldList {
			id := mould.Id
			mouldObjectInfo := make(map[string]interface{})
			json.Unmarshal(mould.ObjectInfo, &mouldObjectInfo)
			dropdownValue := util.InterfaceToString(mouldObjectInfo["toolNo"])
			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:    id,
				Value: dropdownValue,
			})
			if index == 0 {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}

			if mouldId == id {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}
			index = index + 1
		}
		recordInfo.Data = dropDownArray

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, MachineHMITable)
		// Get warning message from stats table
		warningMessage := getWarningMessage(dbConnection, eventId, machineId, AssemblyMachineStatisticsTable)
		var isItAlreadyStarted = isMachineAlreadyStarted(dbConnection, scheduledOrderEvent.MachineId, eventId)
		machineParamId = v.GetMouldTestMachineParam(machineId, eventId)
		return buildHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, warningMessage, recordInfo, mouldDescription, mouldToolNo, hmiId, startedBy, stoppedBy, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList, isItAlreadyStarted, machineParamId)

	}

}

func (v *MachineService) getPreviousHMIResponse(projectId string, eventId, machineId int, operatorId int, machineMasterInfo map[string]interface{}, dbConnection *gorm.DB, targetTable string) datatypes.JSON {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	operatorInfo := authService.GetUserInfoById(operatorId)
	var settingTable string
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	locationList := getLocationDetails(dbConnection)
	var machineParamId = -1
	if targetTable == AssemblyMachineHmiTable {
		settingTable = AssemblyMachineHmiSettingTable
		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		err, scheduledEventObject := productionOrderInterface.GetPreviousAssemblyMachineScheduledOrderEvent(projectId, machineId, eventId)
		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		if len((*scheduledEventObject)) < 1 {
			err, scheduledEventObject = productionOrderInterface.GetLastAssemblyScheduledOrderEvent(projectId, machineId)
		}

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		scheduledOrderEvent := GetScheduledOrderEventInfo((*scheduledEventObject)[0].ObjectInfo)
		eventId = (*scheduledEventObject)[0].Id
		startScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.StartDate)
		endScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.EndDate)
		v.BaseService.Logger.Info("current timeline event:", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " JSON_EXTRACT(object_info, \"$.machineId\") =" + strconv.Itoa(machineId) + " AND JSON_EXTRACT(object_info, \"$.eventId\") =" + strconv.Itoa(eventId) + " order by  JSON_EXTRACT(object_info, \"$.createdAt\") desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, AssemblyMachineHmiTable, hmiInfoConditionString)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		var hmiStatus string
		var hmiID int

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiID = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, scheduledOrderEvent.EventStatus)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		stopList := v.getStopList(eventId, dbConnection, AssemblyMachineHmiTable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, scheduledOrderEvent.EventSourceId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		productionOrderInfo := GetProductionOrderInfo(productionObject.ObjectInfo)
		manufacturingInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
		err, partObject := manufacturingInterface.GetAssemblyPartInfo(projectId, productionOrderInfo.PartNumber)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		partInfo := GetPartInfo(partObject.ObjectInfo)

		err, currentEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineId)
		if eventId != currentEventObject.Id {
			hmiStatus = "disabled"
		}

		if hmiStatus == "" {
			hmiStatus = "stopped"
		}

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, AssemblyMachineHmiTable)

		// Get warning message from stats table
		warningMessage := getWarningMessage(dbConnection, eventId, machineId, AssemblyMachineStatisticsTable)
		var isItAlreadyStarted = isMachineAlreadyStarted(dbConnection, scheduledOrderEvent.MachineId, eventId)
		return buildHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, warningMessage, component.RecordInfo{}, "", "", hmiID, startedBy, stoppedBy, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList, isItAlreadyStarted, machineParamId)

	} else if targetTable == ToolingMachineHmiTable {
		settingTable = ToolingMachineHmiSettingTable
		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		err, scheduledEventObject := productionOrderInterface.GetPreviousToolingMachineScheduledOrderEvent(projectId, machineId, eventId)
		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName)
		}

		if len((*scheduledEventObject)) < 1 {
			err, scheduledEventObject = productionOrderInterface.GetLastToolingScheduledOrderEvent(projectId, machineId)
		}

		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName)
		}

		scheduledOrderEvent := make(map[string]interface{})
		json.Unmarshal((*scheduledEventObject)[0].ObjectInfo, &scheduledOrderEvent)

		eventId = (*scheduledEventObject)[0].Id
		startScheduledDateTime := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["startDate"]))
		endScheduledDateTime := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["endDate"]))
		v.BaseService.Logger.Info("current timeline event:", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " JSON_EXTRACT(object_info, \"$.machineId\") =" + strconv.Itoa(machineId) + " AND JSON_EXTRACT(object_info, \"$.eventId\") =" + strconv.Itoa(eventId) + " order by  JSON_EXTRACT(object_info, \"$.createdAt\") desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, ToolingMachineHmiTable, hmiInfoConditionString)
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		var hmiStatus string
		var hmiId int

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}
		eventSourceId := util.InterfaceToInt(scheduledOrderEvent["eventSourceId"])
		rejectList := getRejectList(listOfHMIInfo, eventId)
		setupTimeList := getSetupTimeList(listOfHMIInfo)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, util.InterfaceToInt(scheduledOrderEvent["eventStatus"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		if util.InterfaceToInt(scheduledOrderEvent["eventStatus"]) == orderStatusId || hmiStatus == "" {
			hmiStatus = "disabled"
		}

		stopList := v.getStopList(eventId, dbConnection, ToolingMachineHmiTable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetToolingProductionOrderInfo(projectId, util.InterfaceToInt(scheduledOrderEvent["eventSourceId"]), machineId)
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}

		productionOrderInfo := make(map[string]interface{})
		json.Unmarshal(productionObject.ObjectInfo, &productionOrderInfo)

		//totalDuration := util.InterfaceToInt(productionOrderInfo["day"])*24 + util.InterfaceToInt(productionOrderInfo["hour"]) + util.InterfaceToInt(productionOrderInfo["minute"])/60

		err, partObject := productionOrderInterface.GetToolingPartById(projectId, util.InterfaceToInt(scheduledOrderEvent["partId"]))
		if err != nil {
			return getUnassignedToolingHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName)
		}
		partInfo := make(map[string]interface{})
		json.Unmarshal(partObject.ObjectInfo, &partInfo)

		totalDuration := util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["endDate"])).DateTimeEpoch - util.ConvertStringToDateTime(util.InterfaceToString(scheduledOrderEvent["startDate"])).DateTimeEpoch

		// Convert seconds to a duration
		duration := time.Second * time.Duration(totalDuration)

		// Format the duration as a string representing hours
		hoursString := fmt.Sprintf("%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60)
		durationString := hoursString + " hours"

		_, currentEventObject := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineId)
		if eventId != currentEventObject.Id {
			hmiStatus = "disabled"
		}

		if hmiStatus == "" {
			hmiStatus = "stopped"
		}

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, ToolingMachineHmiTable)

		_, bomGeneralObject := productionOrderInterface.GetBomInfo(projectId, eventSourceId)
		var bomObjectInfo = make(map[string]interface{})
		json.Unmarshal(bomGeneralObject.ObjectInfo, &bomObjectInfo)
		bomName := util.InterfaceToString(bomObjectInfo["name"])

		return buildToolingHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, "", durationString, bomName, hmiId, startedBy, stoppedBy, setupTimeList, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList)
	} else {
		settingTable = MachineHMISettingSettingTable
		hmiStopReasonList := v.getMachineHMIStopReasonsList(dbConnection, machineId, settingTable)
		v.BaseService.Logger.Info("HMI stop reason list", zap.Any("stop reasons", hmiStopReasonList))

		err, scheduledEventObject := productionOrderInterface.GetPreviousMachineScheduledOrderEvent(projectId, machineId, eventId)
		connectStatusColor, statusName := getMachineConnectStatusColorCode(dbConnection, util.InterfaceToInt(machineMasterInfo["machineConnectStatus"]))
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		if len((*scheduledEventObject)) < 1 {
			err, scheduledEventObject = productionOrderInterface.GetLastScheduledOrderEvent(projectId, machineId)
		}

		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, nil, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		scheduledOrderEvent := GetScheduledOrderEventInfo((*scheduledEventObject)[0].ObjectInfo)
		eventId = (*scheduledEventObject)[0].Id
		startScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.StartDate)
		endScheduledDateTime := util.ConvertStringToDateTime(scheduledOrderEvent.EndDate)
		v.BaseService.Logger.Info("current timeline event:", zap.Any("start_time", startScheduledDateTime), zap.Any("end_time", endScheduledDateTime))

		hmiInfoConditionString := " JSON_EXTRACT(object_info, \"$.machineId\") =" + strconv.Itoa(machineId) + " AND JSON_EXTRACT(object_info, \"$.eventId\") =" + strconv.Itoa(eventId) + " order by  JSON_EXTRACT(object_info, \"$.createdAt\") desc"
		listOfHMIInfo, err := GetConditionalObjects(dbConnection, MachineHMITable, hmiInfoConditionString)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}

		startedBy, stoppedBy, startUrl, stopUrl := getLastHmiOperateUser(listOfHMIInfo)

		var hmiStatus string
		var mouldId int
		var hmiId int
		// we need to send what the last status from the list, since it is descending, we will be having the last one first
		for _, hmiResult := range *listOfHMIInfo {
			machineHMIInfo := MachineHMIInfo{}
			json.Unmarshal(hmiResult.ObjectInfo, &machineHMIInfo)
			if machineHMIInfo.HMIStatus != "" {
				hmiStatus = machineHMIInfo.HMIStatus
				hmiId = hmiResult.Id
				break
			}
		}

		rejectList := getRejectList(listOfHMIInfo, eventId)
		v.BaseService.Logger.Info("rejected list", zap.Any("rejected", rejectList))
		err, eventStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, scheduledOrderEvent.EventStatus)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		eventStatusInfo := GetProductionOrderStatusInfo(eventStatusObject.ObjectInfo)

		stopList := v.getStopList(eventId, dbConnection, MachineHMITable)
		v.BaseService.Logger.Info("stop list", zap.Any("stopped", stopList))

		err, productionObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, scheduledOrderEvent.EventSourceId, machineId)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		productionOrderInfo := GetProductionOrderInfo(productionObject.ObjectInfo)
		mouldModuleInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		err, partObject := mouldModuleInterface.GetPartInfo(projectId, productionOrderInfo.PartNumber)
		if err != nil {
			return getUnassignedHMIResponse(machineId, machineMasterInfo, operatorInfo, hmiStopReasonList, scheduledOrderEvent, nil, nil, "", "", connectStatusColor, statusName, locationList)
		}
		partInfo := GetPartInfo(partObject.ObjectInfo)

		//orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceSeven)
		//if scheduledOrderEvent.EventStatus == orderStatusId || hmiStatus == "" {
		//	err, currentEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineId)
		//	if err == nil && eventId != currentEventObject.Id {
		//		hmiStatus = "disabled"
		//	} else {
		//		if hmiStatus == "" {
		//			hmiStatus = "stopped"
		//		}
		//	}
		//
		//}

		err, currentEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineId)
		if eventId != currentEventObject.Id {
			hmiStatus = "disabled"
		}
		if hmiStatus == "" {
			hmiStatus = "stopped"
		}

		// Get mould ids from part number
		partId := strconv.Itoa(productionOrderInfo.PartNumber)
		_, mouldList := mouldModuleInterface.GetMouldsByPartNo(projectId, partId)
		recordInfo := component.RecordInfo{}
		var dropDownArray []component.OrderedData
		index := 0
		mouldId = scheduledOrderEvent.MouldId

		mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		_, mouldGeneralObj := mouldService.GetMouldInfoById(projectId, mouldId)
		mouldInfo := make(map[string]interface{})
		json.Unmarshal(mouldGeneralObj.ObjectInfo, &mouldInfo)

		var mouldDescription string
		mouldToolNo := util.InterfaceToString(mouldInfo["toolNo"])
		if description, ok := mouldInfo["description"]; ok {
			mouldDescription = util.InterfaceToString(description)
		}

		for _, mould := range mouldList {
			id := mould.Id
			mouldObjectInfo := make(map[string]interface{})
			json.Unmarshal(mould.ObjectInfo, &mouldObjectInfo)
			dropdownValue := util.InterfaceToString(mouldObjectInfo["toolNo"])
			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:    id,
				Value: dropdownValue,
			})
			if index == 0 {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}

			if mouldId == id {
				recordInfo.Index = id
				recordInfo.Value = dropdownValue
			}
			index = index + 1
		}
		recordInfo.Data = dropDownArray

		actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId, MachineHMITable)

		// Get warning message from stats table
		warningMessage := getWarningMessage(dbConnection, eventId, machineId, MachineStatisticsTable)
		var isItAlreadyStarted = isMachineAlreadyStarted(dbConnection, scheduledOrderEvent.MachineId, eventId)
		machineParamId = v.GetMouldTestMachineParam(machineId, eventId)
		return buildHMIResponse(eventId, machineMasterInfo, scheduledOrderEvent, partInfo, eventStatusInfo, hmiStatus, operatorInfo, actualStartTime, actualEndTime, hmiStopReasonList, stopList, rejectList, warningMessage, recordInfo, mouldDescription, mouldToolNo, hmiId, startedBy, stoppedBy, connectStatusColor, statusName, startUrl, stopUrl, operatorId, locationList, isItAlreadyStarted, machineParamId)

	}

}

func (v *MachineService) getMachineHMIStopReasonsList(dbConnection *gorm.DB, machineId int, settingTable string) component.RecordInfo {
	recordInfo := component.RecordInfo{}
	err, hmiSettingGeneralObject := Get(dbConnection, settingTable, machineId)
	if err != nil {
		return recordInfo
	}
	machineHMISetting := MachineHMISetting{ObjectInfo: hmiSettingGeneralObject.ObjectInfo}
	stopReasonsArray := machineHMISetting.getHMISettingInfo().HmiStopReasons
	var stopReasonList []interface{}
	recordInfo.Data = stopReasonList
	recordInfo.IsEdit = false
	recordInfo.Type = "object_array"
	for _, stopReasonId := range stopReasonsArray {
		err, hmiStopReasonGeneralObject := Get(dbConnection, MachineHMIStopReasonTable, stopReasonId)
		if err != nil {
			return recordInfo
		}
		machineHMIStopReasons := MachineHMIStopReason{ObjectInfo: hmiStopReasonGeneralObject.ObjectInfo}
		hmiStopReason := machineHMIStopReasons.getMachineHMIStopReasonInfo()
		fmt.Println("hmiStopReason", hmiStopReason)
		hmiStopReason.Id = hmiStopReasonGeneralObject.Id
		stopReasonList = append(stopReasonList, hmiStopReason)
		recordInfo.Data = stopReasonList
	}
	return recordInfo
}

func IsValidRejectCount(ms *MachineService, dbConnection *gorm.DB, rejectQty int, scheduledOrderInfo map[string]interface{}, productionOrderInfo map[string]interface{}) bool {
	isValid := false
	//Can't exceed actual quantity
	//Solution get complted quantity from time line and check it
	if util.InterfaceToInt(scheduledOrderInfo["completedQty"]) >= rejectQty {
		isValid = true
	}
	//total reject can't exceed completed qty
	//Solution get production order and check prodQty it
	totalRejectedQty := getOverallRejectedQuantity(ms, dbConnection, util.InterfaceToInt(scheduledOrderInfo["eventSourceId"]))

	if util.InterfaceToInt(productionOrderInfo["prodQty"]) >= (totalRejectedQty + rejectQty) {
		isValid = true
	}
	return isValid
}

func isEventAborted(ms *MachineService, dbConnection *gorm.DB, eventId int) bool {
	isAbort := false

	hmiConditionString := "JSON_EXTRACT(object_info, \"$.eventId\") = " + strconv.Itoa(eventId)
	listOfHmiObjects, _ := GetConditionalObjects(dbConnection, MachineHMITable, hmiConditionString)

	if len(*listOfHmiObjects) == 0 {
		return isAbort
	}

	for _, hmiObject := range *listOfHmiObjects {
		machineHMI := MachineHMI{ObjectInfo: hmiObject.ObjectInfo}
		hmiStatus := machineHMI.getMachineHMIInfo().HMIStatus

		if hmiStatus == "aborted" {
			return true
		}
	}

	return isAbort

}

func isAssemblyEventAborted(ms *MachineService, dbConnection *gorm.DB, eventId int) bool {
	isAbort := false

	hmiConditionString := "JSON_EXTRACT(object_info, \"$.eventId\") = " + strconv.Itoa(eventId)
	listOfHmiObjects, _ := GetConditionalObjects(dbConnection, AssemblyMachineHmiTable, hmiConditionString)

	if len(*listOfHmiObjects) == 0 {
		return isAbort
	}

	for _, hmiObject := range *listOfHmiObjects {
		machineHMI := MachineHMI{ObjectInfo: hmiObject.ObjectInfo}
		hmiStatus := machineHMI.getMachineHMIInfo().HMIStatus

		if hmiStatus == "aborted" {
			return true
		}
	}

	return isAbort

}

func reverseSlice(generalObjects *[]component.GeneralObject) *[]component.GeneralObject {
	if generalObjects != nil {
		for i, j := 0, len(*generalObjects)-1; i < j; i, j = i+1, j-1 {
			(*generalObjects)[i], (*generalObjects)[j] = (*generalObjects)[j], (*generalObjects)[i]
		}
	}

	return generalObjects
}

func isToolingEventAborted(ms *MachineService, dbConnection *gorm.DB, eventId int) bool {
	isAbort := false

	hmiConditionString := "JSON_EXTRACT(object_info, \"$.eventId\") = " + strconv.Itoa(eventId)
	listOfHmiObjects, _ := GetConditionalObjects(dbConnection, ToolingMachineHmiTable, hmiConditionString)

	if len(*listOfHmiObjects) == 0 {
		return isAbort
	}

	for _, hmiObject := range *listOfHmiObjects {
		machineHMI := MachineHMI{ObjectInfo: hmiObject.ObjectInfo}
		hmiStatus := machineHMI.getMachineHMIInfo().HMIStatus

		if hmiStatus == "aborted" {
			return true
		}
	}

	return isAbort

}

func (v *MachineService) updateHmiInfo(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err, objectInterface := Get(dbConnection, targetTable, recordId)
	currentData := make(map[string]interface{})
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.Any("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	updatingData := make(map[string]interface{})

	json.Unmarshal(objectInterface.ObjectInfo, &currentData)
	currentData["reasonId"] = updateRequest["reasonId"]
	currentData["remark"] = updateRequest["remark"]
	serializedObject, _ := json.Marshal(currentData)
	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Updating Resource Failed",
				Description: "Error updating resource information due to internal system error. Please report this error to system administrator",
			})
		return
	}
}

func getMachineConnectStatusColorCode(dbConnection *gorm.DB, statusId int) (string, string) {
	var machineConnectStatus map[string]interface{}
	_, connectStatusObject := Get(dbConnection, MachineConnectStatusTable, statusId)
	json.Unmarshal(connectStatusObject.ObjectInfo, &machineConnectStatus)

	searchingColor := util.InterfaceToString(machineConnectStatus["colorCode"])
	searchingName := util.InterfaceToString(machineConnectStatus["status"])

	return searchingColor, searchingName
}
