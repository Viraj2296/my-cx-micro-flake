package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShiftsResponse struct {
	ShiftName             string                  `json:"shiftName"`
	ShiftStatusColorCode  string                  `json:"shiftStatus"` // better to send the color code to front-end, instead of status ID, it is no use.
	ShiftSchedulerDetails []ShiftSchedulerDetails `json:"shiftSchedulerDetails"`
	ListOfLines           []ListOfAssemblyLine    `json:"listOfLines"`
	CanCheckIn            bool                    `json:"canCheckIn"`
	ShiftResourceId       int                     `json:"shiftResourceId"`
	CanShiftStart         bool                    `json:"canShiftStart"`
	CanShiftComplete      bool                    `json:"canShiftComplete"`
	AlertMessage          string                  `json:"alertMessage"`
	ShiftDate             string                  `json:"shiftDate"`
}
type ShiftSchedulerDetails struct {
	MachineName                string       `json:"machineName"`
	Model                      string       `json:"model"`
	MaterialName               string       `json:"materialName"`
	PartImage                  string       `json:"partImage"`
	PartDescription            string       `json:"partDescription"`
	CanEdit                    bool         `json:"canEdit"`
	ScheduledEventId           int          `json:"scheduledEventId"`
	ScheduleName               string       `json:"scheduleName"`
	ShiftActualOutputPartTimer int          `json:"shiftActualOutputPartTimer"`
	ShiftTargetOutputPartTimer int          `json:"shiftTargetOutputPartTimer"`
	ShiftActualOutput          int          `json:"shiftActualOutput"`
	Remarks                    string       `json:"remarks"`
	IsRunning                  bool         `json:"isRunning"`
	IsDisabled                 bool         `json:"isDisabled"`
	NumberOfCheckedIn          int          `json:"numberOfCheckedIn"`
	CheckedInUserInfo          []LabourInfo `json:"checkedInUserInfo"`
	SchedulerStartTime         string       `json:"schedulerStartTime"`
	SchedulerEndTime           string       `json:"schedulerEndTime"`
}

type ListOfAssemblyLine struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// GetHistoryShiftDetails This function generates the shift details with events attached to that shift
// It will return the history of shifts
func (v *ActionService) GetHistoryShiftDetails(ctx *gin.Context) {
	v.Logger.Info("handle get shift details")
	var condition = " DATEDIFF(CURDATE(), STR_TO_DATE(object_info->>'$.shiftStartDate', '%Y-%m-%d')) < 7 AND object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusCompleted)
	err, listOfShifts := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, condition)
	if err != nil {
		// got an error, so return the empty array
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
	// now iterate over shift, and get the schedule order details.
	for _, shiftInterface := range listOfShifts {
		// each shift will have multiple scheduled order events, so it should be the child one
		shiftMasterInfo := database.GetShiftMasterInfo(shiftInterface.ObjectInfo)
		shiftScheduleStartTime := util.HumanReadableDateFormat(shiftMasterInfo.ShiftStartDate)
		shiftResponse := ShiftsResponse{}
		shiftResponse.ShiftName = shiftMasterInfo.ShiftReferenceId
		shiftResponse.ShiftStatusColorCode = v.getStatusColorCode(shiftMasterInfo.ShiftStatus)
		shiftResponse.CanCheckIn = shiftMasterInfo.CanCheckIn
		shiftResponse.CanShiftStart = shiftMasterInfo.CanShiftStart
		shiftResponse.CanShiftComplete = shiftMasterInfo.CanShiftStop
		shiftResponse.ShiftResourceId = shiftInterface.Id
		shiftResponse.ShiftDate = shiftScheduleStartTime

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
			partInfo := GetPartInfo(partObject.ObjectInfo)
			schedulerStartTime := util.HumanReadable12HoursDateTimeFormat(scheduledOrderInfo.StartDate.String())
			schedulerEndTime := util.HumanReadable12HoursDateTimeFormat(scheduledOrderInfo.EndDate.String())
			shiftSchedulerDetails.PartImage = partInfo.Image
			shiftSchedulerDetails.PartDescription = partInfo.Description
			shiftSchedulerDetails.ScheduleName = scheduledOrderInfo.Name
			shiftSchedulerDetails.SchedulerStartTime = schedulerStartTime
			shiftSchedulerDetails.SchedulerEndTime = schedulerEndTime
			// each shift scheduler event id, get the shift production information
			var shiftMasterCondition = " object_info->>'$.shiftId' = " + strconv.Itoa(shiftInterface.Id) + " AND object_info->>'$.scheduledEventId'=" + strconv.Itoa(scheduledOrderEventId)
			err, shiftMasterProductionRecords := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftMasterCondition)
			if err != nil {
				v.Logger.Error("error getting shift information", zap.String("error", err.Error()))
			} else {
				if len(shiftMasterProductionRecords) == 1 {
					shiftProductionInfo := database.GetShiftMasterProductionInfo((shiftMasterProductionRecords)[0].ObjectInfo)
					shiftSchedulerDetails.ShiftActualOutputPartTimer = shiftProductionInfo.ShiftActualOutputPartTimer
					shiftSchedulerDetails.ShiftTargetOutputPartTimer = shiftProductionInfo.ShiftTargetOutputPartTimer
				}
			}
			// get the machine details
			err, machineGeneralComponent := machineService.GetAssemblyMachineInfoById(productionOrderInfo.MachineId)
			if err != nil {
				v.Logger.Error("error getting machine  information", zap.String("error", err.Error()))
				continue
			}
			machineInfo := GetAssemblyMachineMasterInfo(machineGeneralComponent.ObjectInfo)
			shiftSchedulerDetails.MachineName = machineInfo.NewMachineId
			shiftSchedulerDetails.Model = machineInfo.Model

			var assemblyLineId = machineInfo.AssemblyLineOption
			var numberOfCheckedInUsers int
			numberOfCheckedInUsers = 0
			var checkedInUserDetails = make([]LabourInfo, 0) // Initialize to store user details for this specific line
			if shiftAttendanceError == nil {
				for _, attendance := range shiftAttendanceRecords {
					var attendanceInfo = database.GetAttendanceInfo(attendance.ObjectInfo)
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
				// Set the number of checked-in users and their details for the current shiftSchedulerDetails
				shiftSchedulerDetails.NumberOfCheckedIn = numberOfCheckedInUsers
				shiftSchedulerDetails.CheckedInUserInfo = checkedInUserDetails
			} else {
				v.Logger.Error("error getting shift attendance records", zap.String("error", err.Error()))
				shiftSchedulerDetails.NumberOfCheckedIn = 0
				shiftSchedulerDetails.CheckedInUserInfo = make([]LabourInfo, 0)
			}

			if machineInfo.IsMESDriverConfigured {
				shiftSchedulerDetails.CanEdit = false
			} else {
				shiftSchedulerDetails.CanEdit = true
			}
			shiftSchedulerDetails.ScheduledEventId = scheduledOrderEventId
			shiftSchedulerDetails.IsRunning = false
			shiftSchedulerDetails.IsDisabled = true
			shiftResponse.ShiftSchedulerDetails = append(shiftResponse.ShiftSchedulerDetails, shiftSchedulerDetails)

		}
		shiftResponse.ListOfLines = listOfLines
		listOfShiftResponse = append(listOfShiftResponse, shiftResponse)
	}
	v.Logger.Info("generated shift response", zap.Any("shift_response", listOfShiftResponse))
	ctx.JSON(http.StatusOK, listOfShiftResponse)
}

func (v *ActionService) getStatusColorCode(statusId int) string {
	err, statusInterface := database.Get(v.Database, const_util.LabourManagementShiftStatusTable, statusId)
	if err == nil {
		return database.GetShiftStatusInfo(statusInterface.ObjectInfo).ColorCode
	} else {
		v.Logger.Error("error getting color code", zap.String("error", err.Error()))
	}
	return "#000000"
}
