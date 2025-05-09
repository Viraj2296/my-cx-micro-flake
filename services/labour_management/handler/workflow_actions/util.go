package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

func (v *ActionService) getSchedulerEventIdFromMachineLine(machineLineId, shiftId int) int {
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	err, shiftMasterInterface := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, shiftId)
	if err != nil {
		v.Logger.Error("error getting shift information", zap.String("error", err.Error()))
		return 0
	}
	shiftMasterInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
	var listOfShiftEvents = shiftMasterInfo.ScheduledOrderEvents
	for _, event := range listOfShiftEvents {
		err, assemblyScheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, event)
		if err != nil {
			v.Logger.Error("error getting assembly scheduled order information", zap.String("error", err.Error()))
			return 0
		}
		scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(assemblyScheduledOrderInterface.ObjectInfo)
		err, machineInterface := machineService.GetAssemblyMachineInfoById(scheduledOrderInfo.MachineId)
		if err != nil {
			v.Logger.Error("error getting assembly machine master information", zap.String("error", err.Error()))
			return 0
		}
		assemblyMasterInfo := GetAssemblyMachineMasterInfo(machineInterface.ObjectInfo)
		if assemblyMasterInfo.AssemblyLineOption == machineLineId {
			// yes found the machine equal to line selected from the user
			return event
		}

	}
	return 0
}

func (v *ActionService) getLinesForSchedulerEvents(shiftId int) []ListOfAssemblyLine {
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	var listOfPermittedLines = make([]ListOfAssemblyLine, 0)
	err, shiftMasterInterface := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, shiftId)
	if err != nil {
		v.Logger.Error("error getting shift information", zap.String("error", err.Error()))
		return listOfPermittedLines
	}
	shiftMasterInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
	var listOfShiftEvents = shiftMasterInfo.ScheduledOrderEvents

	err, listOfAssemblyMachineInfoInterface := machineService.GetListOfAssemblyLines(const_util.ProjectID)
	if err != nil {
		return listOfPermittedLines
	}

	for _, event := range listOfShiftEvents {
		err, assemblyScheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, event)
		if err != nil {
			v.Logger.Error("error getting assembly scheduled order information", zap.String("error", err.Error()))
			return listOfPermittedLines
		}
		scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(assemblyScheduledOrderInterface.ObjectInfo)
		err, machineInterface := machineService.GetAssemblyMachineInfoById(scheduledOrderInfo.MachineId)
		if err != nil {
			v.Logger.Error("error getting assembly machine master information", zap.String("error", err.Error()))
			return listOfPermittedLines
		}
		assemblyMasterInfo := GetAssemblyMachineMasterInfo(machineInterface.ObjectInfo)
		for _, assemblyMachineInterface := range listOfAssemblyMachineInfoInterface {
			assemblyMachineInfo := database.GetAssemblyMachineLineInfo(assemblyMachineInterface.ObjectInfo)
			if assemblyMasterInfo.AssemblyLineOption == assemblyMachineInterface.Id {
				listOfPermittedLines = append(listOfPermittedLines, ListOfAssemblyLine{
					Id:   assemblyMachineInterface.Id,
					Name: assemblyMachineInfo.Name,
				})
			}

		}

	}
	return listOfPermittedLines
}

func (v *ActionService) createAssemblyHmiEntry(eventId, userId int, hmiStatus string) error {
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	err, assemblyScheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, eventId)
	if err != nil {
		return err
	}
	var scheduleOrderInfo = make(map[string]interface{})
	json.Unmarshal(assemblyScheduledOrderInterface.ObjectInfo, &scheduleOrderInfo)
	var machineId = util.InterfaceToInt(scheduleOrderInfo["machineId"])
	var timeNowStr = util.GetCurrentTime(const_util.ISOTimeLayout)
	var machineHmiInfo = map[string]interface{}{
		"eventId":   eventId,
		"createdAt": timeNowStr,
		"createdBy": userId,
		"hmiStatus": hmiStatus,
		"machineId": machineId,
	}
	err = machineService.CreateAssemblyHmiEntry(const_util.ProjectID, machineHmiInfo)
	return err
}

func (v *ActionService) GetShiftEndTime(shiftStartDate string, shiftStartTime string, period int) (error, string) {
	// Combine the shift start date with shift start time
	dateTimeStr := shiftStartDate + "T" + shiftStartTime + ":00.000Z"

	// Parse the combined string into a time.Time object
	startTime, err := time.Parse("2006-01-02T15:04:05.000Z", dateTimeStr)
	if err != nil {
		v.Logger.Error("error parsing shift start time", zap.String("error", err.Error()))
		return err, ""
	}

	// Add the period hours to get end time
	endTime := startTime.Add(time.Duration(period) * time.Hour)
	
	// Format the end time into RFC3339 format
	formattedEndTime := endTime.Format(time.RFC3339)
	
	return nil, formattedEndTime
}
