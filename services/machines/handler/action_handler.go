package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

func (v *MachineService) getProductionDashboardDisplaySetting(ctx *gin.Context) {
	condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.displayEnabled\"))= 'true'"
	projectId := util.GetProjectId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	allDisplayEnabledMachines, _ := GetConditionalObjects(dbConnection, MachineDisplaySettingTable, condition)
	var machineIds []int
	for _, machineInterface := range *allDisplayEnabledMachines {
		machineIds = append(machineIds, machineInterface.Id)
	}
	machineModuleSetting, err := GetObjects(dbConnection, MachineModuleSettingTable)
	var displayInterval = 25
	if err == nil {
		machineModuleInterface := (*machineModuleSetting)[0]
		machineModule := MachineModuleSetting{ObjectInfo: machineModuleInterface.ObjectInfo}
		displayInterval = machineModule.getMachineModuleSettingInfo().DisplayRotateInterval
	}

	var displaySettingResponse = make(map[string]interface{})
	defaultRecordInfo := component.GetDefaultRecordInfo()
	defaultRecordInfo.Value = machineIds
	defaultRecordInfo.IsEdit = false
	displaySettingResponse["displayEnabledMachines"] = defaultRecordInfo
	displaySettingResponse["displayInterval"] = component.GetRecordIntInfo(displayInterval, "int")
	ctx.JSON(http.StatusOK, displaySettingResponse)
}

func (v *MachineService) addingManualAssemblyResetMessage(ctx *gin.Context) {
	//var userId = common.GetUserId(ctx)
	//var createRequest = make(map[string]interface{})
	//if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
	//	ctx.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	//
	//var machineId = util.InterfaceToInt(createRequest["machineId"])
	//var eventId = util.InterfaceToInt(createRequest["eventId"])
	//condition := " object_info ->> '$.eventId' = " + strconv.Itoa(eventId) + " and object_info ->> '$.machineId' = " + strconv.Itoa(machineId)
	//projectId := util.GetProjectId(ctx)
	//dbConnection := v.BaseService.ServiceDatabases[projectId]
	//listOfMachineHmi, err := GetConditionalObjects(dbConnection, AssemblyMachineHmiComponent, condition)
	//
	//if err != nil {
	//	response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
	//	return
	//}
	//
	//if len(*listOfMachineHmi) < 1 {
	//	response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
	//	return
	//}
	//
	//var assemblyReset = map[string]interface{}{
	//	"eventId":   eventId,
	//	"machineId": machineId,
	//	"createdAt": util.GetCurrentTime(ISOTimeLayout),
	//	"createdBy": userId,
	//	"remark":    "reset",
	//}
	//
	//serializedRequest, _ := json.Marshal(assemblyReset)
	//object := component.GeneralObject{
	//	ObjectInfo: serializedRequest,
	//}
	//
	//err, _ = Create(dbConnection, AssemblyMachineHmiComponent, object)
	//if err != nil {
	//	response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, response.GeneralResponse{
	//	Code:    0,
	//	Message: "This machine is successfully reseted",
	//})
	var userId = common.GetUserId(ctx)
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	projectId := util.GetProjectId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var machineId = util.InterfaceToInt(createRequest["machineId"])
	var eventId = util.InterfaceToInt(createRequest["eventId"])

	err, machineHmiObject := Get(dbConnection, AssemblyMachineHmiSettingTable, machineId)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	machineHMISetting := MachineHMISetting{ObjectInfo: machineHmiObject.ObjectInfo}
	var listOfResetOperators = machineHMISetting.getHMISettingInfo().ResetOperators
	var operatorConfigured bool
	operatorConfigured = false
	for _, operatorId := range listOfResetOperators {
		if userId == operatorId {
			operatorConfigured = true
		}
	}
	if !operatorConfigured {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid HMI Operator"), ErrorGettingObjectsInformation, "You are not configured as a HMI Operator, please contact admin regarding this")
		return
	}

	err, machineMasterObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	var machineMasterObjectInfo = make(map[string]interface{})
	err = json.Unmarshal(machineMasterObject.ObjectInfo, &machineMasterObjectInfo)

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Reset is failed",
				Description: "Error in reading assembly machine data",
			})
		return
	}

	var machineConnectStatus = util.InterfaceToInt(machineMasterObjectInfo["machineConnectStatus"])
	if machineConnectStatus == machineConnectStatusWaitingForFeed {
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
		err, scheduledEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineId)

		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		}

		if scheduledEventObject.Id != 0 {
			var scheduledEventObjectInfo = make(map[string]interface{})
			json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEventObjectInfo)

			var updateScheduleEvent = make(map[string]interface{})
			updateScheduleEvent["percentDone"] = 0
			updateScheduleEvent["completedQty"] = 0

			serilizedObject, _ := json.Marshal(updateScheduleEvent)

			err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(projectId, scheduledEventObject.Id, serilizedObject)

			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Reset is failed",
						Description: "Please contact admin for further details",
					})
				return
			}
		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Validation Failed",
					Description: "Currently this machine doesn't have scheduled order.",
				})
			return
		}
	}

	condition := " object_info ->> '$.eventId' = " + strconv.Itoa(eventId) + " and object_info ->> '$.machineId' = " + strconv.Itoa(machineId)

	listOfMachineHmi, err := GetConditionalObjects(dbConnection, AssemblyMachineHmiComponent, condition)

	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	if len(*listOfMachineHmi) < 1 {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Reset is failed",
				Description: "Please start hmi before reset",
			})
		return
	}

	var assemblyReset = map[string]interface{}{
		"eventId":   eventId,
		"machineId": machineId,
		"createdAt": util.GetCurrentTime(ISOTimeLayout),
		"createdBy": userId,
		"remark":    "reset",
	}

	serializedRequest, _ := json.Marshal(assemblyReset)
	object := component.GeneralObject{
		ObjectInfo: serializedRequest,
	}

	errr, _ := Create(dbConnection, AssemblyMachineHmiComponent, object)
	if errr != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorCreatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "This machine is successfully reseted",
	})
}
