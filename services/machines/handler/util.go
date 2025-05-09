package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func isMachineAlreadyStarted(dbConnection *gorm.DB, machineId int, eventId int) bool {
	// check whether this machine is already started
	conditionString := " object_info->>'$.machineId' =" + strconv.Itoa(machineId) + " AND object_info->>'$.eventId' = " + strconv.Itoa(eventId) + " AND object_info->>'$.hmiStatus' = 'started'"
	var startedCount = CountByCondition(dbConnection, MachineHMITable, conditionString)
	var isItAlreadyStarted bool
	if startedCount == 0 {
		isItAlreadyStarted = false
	} else {
		isItAlreadyStarted = true
	}
	return isItAlreadyStarted
}

func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}

func getMouldingMachineOrderDetails(machineId int) *ScheduledOrderEventInfo {
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEvent(ProjectID, machineId)
	if err == nil {
		schedulerOrderEventInfo := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
		return schedulerOrderEventInfo
	}

	return nil

}

func getMachineOrderDetails(machineId int) *ScheduledOrderEventInfo {
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, scheduledEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(ProjectID, machineId)
	if err == nil {
		schedulerOrderEventInfo := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
		return schedulerOrderEventInfo
	}

	return nil

}

func (v *MachineService) GetMouldTestMachineParam(machineId int, eventId int) int {
	var machineParam = -1
	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, c := productionOrderInterface.GetScheduledOrderInfo(ProjectID, eventId)
	if err == nil {
		var scheduledOrderEvent = make(map[string]interface{})
		err := json.Unmarshal(c.ObjectInfo, &scheduledOrderEvent)
		if err != nil {
			return machineParam
		}
		if mouldIdInterface, ok := scheduledOrderEvent["mouldId"]; ok {
			var mouldId = util.InterfaceToInt(mouldIdInterface)
			machineParam = mouldService.GetMouldMachineTestParam(machineId, mouldId)
		}
	}
	return machineParam
}

func formatTime(input string) (string, error) {
	// Parse the input time string
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return "", err
	}

	// Format the output in the required format
	formatted := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dT",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return formatted, nil
}
