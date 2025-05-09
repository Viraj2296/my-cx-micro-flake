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

type ShiftStatusSummary struct {
	TotalActualManPower  int                   `json:"totalActualManPower"`
	TotalPlannedManPower int                   `json:"totalPlannedManPower"`
	ShiftProductionInfo  []ShiftProductionInfo `json:"shiftProductionInfo"`
	ManpowerUtilisation  float64               `json:"manpowerUtilisation"`
}
type ShiftProductionInfo struct {
	RowId                       int     `json:"rowId"`
	PlannedLines                string  `json:"plannedLines"`
	PlannedManPower             int     `json:"plannedManPower"`
	ActualManpower              int     `json:"actualManpower"`
	ShiftActualOutput           int     `json:"shiftActualOutput"`
	ShiftTargetOutput           int     `json:"shiftTargetOutput"`
	ProductionCompletedPercent  float64 `json:"productionCompletedPercent"`
	LineManpowerUtilisation     float64 `json:"lineManpowerUtilisation"`
	ProductionCompletedQuantity int     `json:"productionCompletedQuantity"`
	LineManpower                int     `json:"lineManpower"`
}

// GetShiftProductionTvView  Note : https://cerexio.atlassian.net/browse/FUYU2-349
func (v *ActionService) GetShiftProductionTvView(ctx *gin.Context) {
	v.Logger.Info("getting the shift production tv view")
	var condition = " object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusActive) + " ORDER BY object_info->>'$.createdAt' DESC"
	var shiftStatusSummary = ShiftStatusSummary{
		TotalActualManPower:  0,
		TotalPlannedManPower: 0,
		ShiftProductionInfo:  make([]ShiftProductionInfo, 0),
	}
	err, listOfShifts := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, condition)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	if err != nil {
		// got an error, so return the empty array
		v.Logger.Error("error getting shift information", zap.String("error", err.Error()))
		ctx.JSON(http.StatusOK, shiftStatusSummary)
		return
	}
	if len(listOfShifts) == 0 {
		v.Logger.Warn("There is no shift information found", zap.String("condition", condition))
		ctx.JSON(http.StatusOK, shiftStatusSummary)
		return
	}
	var rowId int
	rowId = 1
	for _, activeShiftObject := range listOfShifts {

		v.Logger.Info("selected shift for production tv view", zap.Int("shift_master_id", activeShiftObject.Id))
		// get all the shift master production data
		var shiftProductionCondition = " object_info->>'$.shiftId' = " + strconv.Itoa(activeShiftObject.Id)
		err, listOfShiftProduction := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftProductionCondition)
		if err != nil {
			// got an error, so return the empty array
			v.Logger.Error("error getting shift production", zap.String("error", err.Error()))
			ctx.JSON(http.StatusOK, shiftStatusSummary)
			return
		}
		if len(listOfShiftProduction) == 0 {
			v.Logger.Warn("There is no shift production found", zap.String("condition", shiftProductionCondition))
			ctx.JSON(http.StatusOK, shiftStatusSummary)
			return
		}
		var attendanceCondition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(activeShiftObject.Id)
		v.Logger.Info("selected shift for production tv view, attendance condition", zap.String("attendance_condition", attendanceCondition))
		err, attendanceRecords := database.GetConditionalObjects(v.Database, const_util.LabourManagementAttendanceTable, attendanceCondition)
		var attendanceLineCache = make(map[int]int)
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		var listOfUsersCheckedIn []int
		if err == nil {
			for _, attendance := range attendanceRecords {
				attendanceInfo := database.GetAttendanceInfo(attendance.ObjectInfo)
				// don't add the checkout users
				if attendanceInfo.CheckOutTime == "" || attendanceInfo.CheckOutDate == "" {
					for _, lineId := range attendanceInfo.ManufacturingLines {
						userInfo := authService.GetUserInfoById(attendanceInfo.UserResourceId)

						if userInfo.JobRole == v.LabourManagementSettingInfo.ShiftOperatorJobRoleId {
							if !util.HasInt(userInfo.UserId, listOfUsersCheckedIn) {
								var initialCount = attendanceLineCache[lineId]
								attendanceLineCache[lineId] = initialCount + 1
								v.Logger.Info("adding line in to cache", zap.Int("lineId", lineId), zap.Int("count", attendanceLineCache[lineId]))
							}
							listOfUsersCheckedIn = append(listOfUsersCheckedIn, userInfo.UserId)
						}
					}
				}

			}
		}

		for _, shiftMasterProductionInterface := range listOfShiftProduction {
			var shiftProductionInfo = ShiftProductionInfo{}
			//select  object_info->>'$.name' as `name` from assembly_machine_lines where id  =  (select object_info->>'$.assemblyLineOption' from assembly_machine_master where id = [machine_id])
			shiftMasterProductionInfo := database.GetShiftMasterProductionInfo(shiftMasterProductionInterface.ObjectInfo)
			err, assemblyMachineObject := machineService.GetAssemblyMachineInfoById(shiftMasterProductionInfo.MachineId)
			if err != nil {
				v.Logger.Error("error getting machine details", zap.String("error", err.Error()))
				continue
			}

			assemblyMachineMasterInfo := GetAssemblyMachineMasterInfo(assemblyMachineObject.ObjectInfo)
			err, assemblyLineObject := machineService.GetAssemblyLineFromId(const_util.ProjectID, assemblyMachineMasterInfo.AssemblyLineOption)
			if err == nil {
				shiftProductionInfo.PlannedLines = database.GetAssemblyMachineLineInfo(assemblyLineObject.ObjectInfo).Name
				if shiftProductionInfo.PlannedLines == "PSU2 Manual" {
					continue
				}
				if count, ok := attendanceLineCache[assemblyLineObject.Id]; ok {
					shiftProductionInfo.ActualManpower = count
					shiftStatusSummary.TotalActualManPower = shiftStatusSummary.TotalActualManPower + count
				}
			} else {
				shiftProductionInfo.PlannedLines = "-"
			}

			err, scheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, shiftMasterProductionInfo.ScheduledEventId)
			if err == nil {
				scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(scheduledOrderInterface.ObjectInfo)
				v.Logger.Info("scheduler order info ", zap.Int("event_id", shiftMasterProductionInfo.ScheduledEventId), zap.Any("production_info", scheduledOrderInfo))

				shiftProductionInfo.ShiftTargetOutput = scheduledOrderInfo.ScheduledQty
				shiftProductionInfo.ShiftActualOutput = scheduledOrderInfo.CompletedQty
				shiftProductionInfo.PlannedManPower = scheduledOrderInfo.PlannedManPower
				shiftStatusSummary.TotalPlannedManPower = shiftStatusSummary.TotalPlannedManPower + scheduledOrderInfo.PlannedManPower

				// Calculate the new fields
				shiftProductionInfo.ProductionCompletedQuantity = shiftProductionInfo.ShiftActualOutput - shiftProductionInfo.ShiftTargetOutput
				shiftProductionInfo.LineManpower = shiftProductionInfo.ActualManpower - shiftProductionInfo.PlannedManPower
			} else {
				shiftProductionInfo.ShiftTargetOutput = 0
				shiftProductionInfo.ProductionCompletedQuantity = 0
				shiftProductionInfo.LineManpower = 0
			}

			// Existing percentage calculations
			shiftProductionInfo.ProductionCompletedPercent = util.CalculatePercentage(shiftProductionInfo.ShiftActualOutput, shiftProductionInfo.ShiftTargetOutput)
			var lineManpowerUtilisation float64
			lineManpowerUtilisation = util.CalculatePercentage(shiftProductionInfo.ActualManpower, shiftProductionInfo.PlannedManPower)
			shiftProductionInfo.LineManpowerUtilisation = lineManpowerUtilisation

			shiftProductionInfo.RowId = rowId
			shiftStatusSummary.ShiftProductionInfo = append(shiftStatusSummary.ShiftProductionInfo, shiftProductionInfo)
			rowId = rowId + 1
		}

		var manpowerUtilisation float64

		manpowerUtilisation = util.CalculatePercentage(shiftStatusSummary.TotalActualManPower, shiftStatusSummary.TotalPlannedManPower)
		shiftStatusSummary.ManpowerUtilisation = manpowerUtilisation
	}

	v.Logger.Info("generated shift production info ", zap.Any("listOfShiftProductionInfo", shiftStatusSummary))
	ctx.JSON(http.StatusOK, shiftStatusSummary)
}
