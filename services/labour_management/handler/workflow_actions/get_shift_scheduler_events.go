package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ShiftTemplateEventRequest struct {
	ShiftTemplateId int `json:"shiftTemplateId"`
}

// GetShiftSchedulerEvents  : Getting all the shift scheduler events/*
func (v *ActionService) GetShiftSchedulerEvents(ctx *gin.Context) {
	v.Logger.Info("handle get shift scheduler events")
	shiftTemplateEventRequest := ShiftTemplateEventRequest{}

	if err := ctx.ShouldBindBodyWith(&shiftTemplateEventRequest, binding.JSON); err != nil {
		err := ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("sending error to client has failed", zap.String("error", err.Error()))
		}
		return
	}

	var shiftTemplateId = shiftTemplateEventRequest.ShiftTemplateId
	err, shiftTemplateInterface := database.Get(v.Database, const_util.LabourManagementShiftTemplateTable, shiftTemplateId)
	if err != nil {
		v.Logger.Error("error getting shift template", zap.String("error", err.Error()))
		response.SendResourceNotFound(ctx)
	}
	var shiftTemplateInfo = database.GetSShiftTemplateInfo(shiftTemplateInterface.ObjectInfo)
	err, shiftStartModifiedTime := v.GetShiftTime(shiftTemplateInfo.ShiftStartTime, 0)
	if err != nil {
		response.SendInternalSystemError(ctx)
	}
	err, shiftEndTime := v.GetShiftTime(shiftTemplateInfo.ShiftStartTime, shiftTemplateInfo.ShiftPeriod)
	if err != nil {
		response.SendInternalSystemError(ctx)
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	// get the available events
	v.Logger.Info("searching shift scheduler event start time", zap.String("start_time", shiftStartModifiedTime), zap.String("shift_end_time", shiftEndTime))
	err, listOfAvailableEvents := productionOrderInterface.GetReleasedAssemblyOrderEventsBetween(const_util.ProjectID, shiftStartModifiedTime, shiftEndTime)

	var alreadyUsedEvents = v.LoadUsedEvents()
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Error Getting Released Orders",
				Description: "Internal system error getting released order under chosen shift " + shiftTemplateInfo.Name,
			})
		return
	}

	if len(*listOfAvailableEvents) == 0 {
		v.Logger.Warn("There is no shift information found for a period", zap.String("shift_start_time", shiftStartModifiedTime), zap.String("shift_end_time", shiftEndTime))
		var shiftResponse = make([]ShiftsResponse, 0)
		ctx.JSON(http.StatusOK, shiftResponse)
		return
	}

	manufacturingModuleInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)

	err, listOfAssemblyMachineInfoInterface := machineService.GetListOfAssemblyLines(const_util.ProjectID)
	var listOfLines = make([]ListOfAssemblyLine, 0)
	if err == nil {
		for _, assemblyMachineInterface := range listOfAssemblyMachineInfoInterface {
			assemblyMachineInfo := database.GetAssemblyMachineLineInfo(assemblyMachineInterface.ObjectInfo)
			listOfLines = append(listOfLines, ListOfAssemblyLine{
				Id:   assemblyMachineInterface.Id,
				Name: assemblyMachineInfo.Name,
			})
		}
	}
	shiftResponse := ShiftsResponse{}
	// now iterate over shift, and get the schedule order details.
	for _, scheduledEventsInterface := range *listOfAvailableEvents {

		// before sending all the events, we need to make sure we have filtered out the events that are already created
		// in the previous shifts
		if util.HasInt(scheduledEventsInterface.Id, alreadyUsedEvents) {
			v.Logger.Info("event is already used, so skip sending it", zap.Int("event_id", scheduledEventsInterface.Id))
			continue
		}
		shiftSchedulerDetails := ShiftSchedulerDetails{}
		scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(scheduledEventsInterface.ObjectInfo)
		err, productionOrderData := productionOrderInterface.GetAssemblyProductionOrderById(scheduledOrderInfo.EventSourceId)
		if err != nil {
			v.Logger.Error("error getting production order information", zap.String("error", err.Error()))
			continue
		}
		var productionOrderInfo = GetAssemblyProductionOrderInfo(productionOrderData.ObjectInfo)
		_, partObject := manufacturingModuleInterface.GetAssemblyPartInfo(const_util.ProjectID, productionOrderInfo.PartNumber)
		partInfo := GetPartInfo(partObject.ObjectInfo)
		shiftSchedulerDetails.PartImage = partInfo.Image
		shiftSchedulerDetails.PartDescription = partInfo.Description
		shiftSchedulerDetails.ScheduleName = scheduledOrderInfo.Name
		// get the machine details
		err, machineGeneralComponent := machineService.GetAssemblyMachineInfoById(productionOrderInfo.MachineId)
		if err != nil {
			v.Logger.Error("error getting machine  information", zap.String("error", err.Error()))
			continue
		}
		machineInfo := GetAssemblyMachineMasterInfo(machineGeneralComponent.ObjectInfo)
		shiftSchedulerDetails.MachineName = machineInfo.NewMachineId
		shiftSchedulerDetails.Model = machineInfo.Model
		if machineInfo.IsMESDriverConfigured {
			shiftSchedulerDetails.CanEdit = false
		} else {
			shiftSchedulerDetails.CanEdit = true
		}
		shiftSchedulerDetails.ScheduledEventId = scheduledEventsInterface.Id
		shiftResponse.ShiftSchedulerDetails = append(shiftResponse.ShiftSchedulerDetails, shiftSchedulerDetails)
	}
	shiftResponse.ListOfLines = listOfLines

	v.Logger.Info("generated shift response", zap.Any("shift_response", shiftResponse))
	ctx.JSON(http.StatusOK, shiftResponse)
}
func (v *ActionService) LoadUsedEvents() []int {
	err, listOfShiftMasterInterface := database.GetObjects(v.Database, const_util.LabourManagementShiftMasterTable)
	var listOfUsedEvents []int
	if err == nil {
		for _, shiftMasterInterface := range listOfShiftMasterInterface {
			shiftMasterInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
			listOfUsedEvents = append(listOfUsedEvents, shiftMasterInfo.ScheduledOrderEvents...)
		}
	}
	return listOfUsedEvents
}
func (v *ActionService) GetShiftTime(shiftEndTime string, period int) (error, string) {
	//shiftEndTime := "13:52"

	// Get the current date
	now := time.Now()

	// Add 8 hours to the current time to adjust it for singapore, later it should be taken from system operation time or user timezone
	//TODO later move this hard coded one to system setting.
	futureTime := now.Add(8 * time.Hour)
	// Format the new time to only get the date in YYYY-MM-DD format
	currentDate := futureTime.Format("2006-01-02")
	// Combine the current date with the shift end time
	dateTimeStr := currentDate + "T" + shiftEndTime + ":00.000Z"

	// Parse the combined string into a time.Time object
	correctionTime, err := time.Parse("2006-01-02T15:04:05.000Z", dateTimeStr)
	if err != nil {
		v.Logger.Error("error parsing shift time", zap.String("error", err.Error()))
		return err, ""
	}
	// Reduce 8 hours from the parsed time
	reducedTime := correctionTime.Add(-8 * time.Hour)

	if period > 0 {
		// Add 12 hours to the parsed time
		newEndTime := reducedTime.Add(time.Duration(period) * time.Hour)
		// Format the parsed time into the desired timestamp format
		formattedTime := newEndTime.Format(time.RFC3339Nano)
		return nil, formattedTime
	}
	// Format the parsed time into the desired timestamp format
	formattedTime := reducedTime.Format(time.RFC3339Nano)
	return nil, formattedTime
}

func adjustTimestamp(input string) (string, error) {
	// Define the layout to parse the input timestamp
	const layout = "2006-01-02T15:04:05Z"

	// Parse the input timestamp
	t, err := time.Parse(layout, input)
	if err != nil {
		return "", err
	}

	// Subtract 8 hours
	t = t.Add(-8 * time.Hour)

	// Format back to the desired format without 'Z'
	const outputLayout = "2006-01-02T15:04:05"
	output := t.Format(outputLayout)

	return output, nil
}
