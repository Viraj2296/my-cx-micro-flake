package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) isTimeExceeded(givenTime string) bool {
	// Parse the given time string into a time.Time object
	parsedTime, err := time.Parse(time.RFC3339, givenTime)
	if err != nil {
		// Print the error and return false if there's a parsing issue
		return false
	}

	// Get the current time
	currentTime := time.Now()
	v.Logger.Info("time check", zap.String("givenTime", givenTime), zap.String("currentTime", currentTime.String()), zap.String("parsedTime", parsedTime.String()))
	// Compare current time with the parsed given time
	return currentTime.After(parsedTime)
}

// GetActiveShiftDetails This function generates the shift details with events attached to that shift
// It will return the active shifts
func (v *ActionService) GetActiveShiftDetails(ctx *gin.Context) {
	v.Logger.Info("handle get shift details")
	var condition = " object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusPending) + " OR " + "object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusActive)
	err, listOfShifts := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, condition)
	if err != nil {
		v.Logger.Error("error getting shift information", zap.String("error", err.Error()))
		var shiftResponse = make([]ShiftsResponse, 0)
		ctx.JSON(http.StatusOK, shiftResponse)
		return
	}
	if len(listOfShifts) == 0 {
		v.Logger.Warn("There is no shift information found", zap.String("condition", condition))
		var shiftResponse = make([]ShiftsResponse, 0)
		ctx.JSON(http.StatusOK, shiftResponse)
		return
	}

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	manufacturingModuleInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
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
	var listOfShiftResponse []ShiftsResponse
	for _, shiftInterface := range listOfShifts {
		shiftMasterInfo := database.GetShiftMasterInfo(shiftInterface.ObjectInfo)
		shiftScheduleStartTime := util.HumanReadableDateFormat(shiftMasterInfo.ShiftStartDate)
		shiftResponse := ShiftsResponse{
			ShiftName:             shiftMasterInfo.ShiftReferenceId,
			ShiftStatusColorCode:  v.getStatusColorCode(shiftMasterInfo.ShiftStatus),
			CanCheckIn:            shiftMasterInfo.CanCheckIn,
			CanShiftStart:         shiftMasterInfo.CanShiftStart,
			CanShiftComplete:      shiftMasterInfo.CanShiftStop,
			ShiftResourceId:       shiftInterface.Id,
			AlertMessage:          "",
			ShiftSchedulerDetails: []ShiftSchedulerDetails{},
			ListOfLines:           listOfLines,
			ShiftDate:             shiftScheduleStartTime,
		}

		err, shiftTemplateInterface := database.Get(v.Database, const_util.LabourManagementShiftTemplateTable, shiftMasterInfo.ShiftTemplateId)
		if err == nil {
			shiftTemplateInfo := database.GetSShiftTemplateInfo(shiftTemplateInterface.ObjectInfo)
			err, shiftEndTime := v.GetShiftTime(shiftTemplateInfo.ShiftStartTime, shiftTemplateInfo.ShiftPeriod)
			if err == nil && v.isTimeExceeded(shiftEndTime) {
				shiftResponse.AlertMessage = "This shift has currently overrun. Please close it"
			}
		}
		attendanceCondition := " object_info->>'$.shiftResourceId' = " + strconv.Itoa(shiftInterface.Id)
		shiftAttendanceError, shiftAttendanceRecords := database.GetConditionalObjects(v.Database, const_util.LabourManagementAttendanceTable, attendanceCondition)

		for _, scheduledOrderEventId := range shiftMasterInfo.ScheduledOrderEvents {
			shiftSchedulerDetails := ShiftSchedulerDetails{}
			err, scheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, scheduledOrderEventId)
			if err != nil {
				v.Logger.Error("error getting scheduler order information", zap.String("error", err.Error()))
				continue
			}
			scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(scheduledOrderInterface.ObjectInfo)
			err, productionOrderData := productionOrderInterface.GetAssemblyProductionOrderById(scheduledOrderInfo.EventSourceId)
			if err != nil {
				v.Logger.Error("error getting production order information", zap.String("error", err.Error()))
				continue
			}
			var productionOrderInfo = GetAssemblyProductionOrderInfo(productionOrderData.ObjectInfo)
			_, partObject := manufacturingModuleInterface.GetAssemblyPartInfo(const_util.ProjectID, productionOrderInfo.PartNumber)
			schedulerStartTime := util.HumanReadable12HoursDateTimeFormat(scheduledOrderInfo.StartDate.String())
			schedulerEndTime := util.HumanReadable12HoursDateTimeFormat(scheduledOrderInfo.EndDate.String())
			partInfo := GetPartInfo(partObject.ObjectInfo)
			shiftSchedulerDetails.PartImage = partInfo.Image
			shiftSchedulerDetails.PartDescription = partInfo.Description
			shiftSchedulerDetails.ScheduleName = scheduledOrderInfo.Name
			shiftSchedulerDetails.SchedulerStartTime = schedulerStartTime
			shiftSchedulerDetails.SchedulerEndTime = schedulerEndTime

			var shiftMasterCondition = " object_info->>'$.shiftId' = " + strconv.Itoa(shiftInterface.Id) + " AND object_info->>'$.scheduledEventId'=" + strconv.Itoa(scheduledOrderEventId)
			err, shiftMasterProductionRecords := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftMasterCondition)
			if err == nil && len(shiftMasterProductionRecords) == 1 {
				shiftProductionInfo := database.GetShiftMasterProductionInfo(shiftMasterProductionRecords[0].ObjectInfo)
				shiftSchedulerDetails.ShiftActualOutputPartTimer = shiftProductionInfo.ShiftActualOutputPartTimer
				shiftSchedulerDetails.ShiftTargetOutputPartTimer = shiftProductionInfo.ShiftTargetOutputPartTimer
				shiftSchedulerDetails.ShiftActualOutput = scheduledOrderInfo.CompletedQty
				shiftSchedulerDetails.Remarks = scheduledOrderInfo.Remarks
			}

			err, machineGeneralComponent := machineService.GetAssemblyMachineInfoById(productionOrderInfo.MachineId)
			if err == nil {
				machineInfo := GetAssemblyMachineMasterInfo(machineGeneralComponent.ObjectInfo)
				shiftSchedulerDetails.MachineName = machineInfo.NewMachineId
				shiftSchedulerDetails.Model = machineInfo.Model
				shiftSchedulerDetails.CanEdit = !machineInfo.IsMESDriverConfigured
				var assemblyLineId = machineInfo.AssemblyLineOption
				var numberOfCheckedInUsers int
				numberOfCheckedInUsers = 0
				var checkedInUserDetails = make([]LabourInfo, 0) // Initialize to store user details for this specific line
				if shiftAttendanceError == nil {
					for _, attendance := range shiftAttendanceRecords {
						var attendanceInfo = database.GetAttendanceInfo(attendance.ObjectInfo)
						if attendanceInfo.CheckOutDate == "" || attendanceInfo.CheckOutTime == "" {
							var listOfLinesCheckedIn = attendanceInfo.ManufacturingLines
							if util.HasInt(assemblyLineId, listOfLinesCheckedIn) {
								numberOfCheckedInUsers += 1
								userInfo := authService.GetUserInfoById(attendanceInfo.UserResourceId)
								var jobRoleName = authService.GetJobRoleName(userInfo.JobRole)
								checkedInUserDetails = append(checkedInUserDetails, LabourInfo{
									Name:      userInfo.FullName,
									AvatarUrl: userInfo.AvatarUrl,
									Role:      jobRoleName,
								})
							}
						}

					}
					// Set the number of checked-in users and their details for the current shiftSchedulerDetails
					shiftSchedulerDetails.NumberOfCheckedIn = numberOfCheckedInUsers
					shiftSchedulerDetails.CheckedInUserInfo = checkedInUserDetails
				} else {
					v.Logger.Error("error getting shift attendance records", zap.String("error", err.Error()))
					shiftSchedulerDetails.NumberOfCheckedIn = 0
					shiftSchedulerDetails.CheckedInUserInfo = make([]LabourInfo, 0)
				}

			} else {
				shiftSchedulerDetails.NumberOfCheckedIn = 0
				shiftSchedulerDetails.CheckedInUserInfo = make([]LabourInfo, 0)
			}
			shiftSchedulerDetails.ScheduledEventId = scheduledOrderEventId

			if scheduledOrderInfo.EventStatus == const_util.ScheduleStatusPreferenceFive {
				shiftSchedulerDetails.IsRunning = true
				shiftSchedulerDetails.IsDisabled = false
			} else {
				shiftSchedulerDetails.IsRunning = false
				shiftSchedulerDetails.IsDisabled = true
			}

			shiftResponse.ShiftSchedulerDetails = append(shiftResponse.ShiftSchedulerDetails, shiftSchedulerDetails)
		}

		listOfShiftResponse = append(listOfShiftResponse, shiftResponse)
	}
	v.Logger.Info("generated shift response", zap.Any("shift_response", listOfShiftResponse))
	ctx.JSON(http.StatusOK, listOfShiftResponse)
}
