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
	"strconv"
)

type CheckOutRequest struct {
	EmployeeId string `json:"employeeId"`
	ShiftId    int    `json:"shiftId"`
}

func (v *ActionService) HandleCheckOut(ctx *gin.Context) {
	v.Logger.Info("handle check out received")
	checkOutRequest := CheckOutRequest{}
	if err := ctx.ShouldBindBodyWith(&checkOutRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	v.Logger.Info("checkout payload", zap.Any("checkOutRequest", checkOutRequest))
	recordId := checkOutRequest.ShiftId
	// Call getBasicInfo with recordId
	_, dbConnection, userInfo, shiftMasterInfo := v.getShiftMasterInfo(ctx, recordId)
	v.Logger.Info("checkout is requested from", zap.Int("user_id", userInfo.UserId))
	if shiftMasterInfo.ShiftStatus == const_util.ShiftStatusActive {
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		v.Logger.Info("requesting user details for employee id", zap.String("employeeId", checkOutRequest.EmployeeId))
		checkOutUserInfo := authService.GetUserInfoByEmployeeId(checkOutRequest.EmployeeId)
		v.Logger.Info("checkout user info", zap.String("username", checkOutUserInfo.Username), zap.Int("user_id", checkOutUserInfo.UserId))
		if checkOutUserInfo.UserId == 0 {
			v.Logger.Error("invalid user id", zap.Int("user_id", checkOutUserInfo.UserId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid User",
					Description: "The user is not available in the system, please contact system admin to configure the user",
				})
			return
		}
		// check the supervisor is already checked in
		var condition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(recordId) + " AND object_info->>'$.userResourceId' =" + strconv.Itoa(checkOutUserInfo.UserId)
		err, listOfAttendance := database.GetConditionalObjects(dbConnection, const_util.LabourManagementAttendanceTable, condition)
		if err != nil {
			v.Logger.Error("error getting attendance", zap.String("error", err.Error()), zap.String("condition", condition))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		var listOfLinesCheckedIn []int
		for _, attendanceInterface := range listOfAttendance {
			attendanceInfo := database.GetAttendanceInfo(attendanceInterface.ObjectInfo)
			attendanceInfo.CheckOutDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			attendanceInfo.CheckOutTime = util.GetCurrentTimeInSingapore("15:04:05")
			attendanceInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			attendanceInfo.LastUpdatedBy = userInfo.UserId
			updateObject := make(map[string]interface{})
			updateObject["object_info"] = attendanceInfo.Serialised()
			err := database.Update(dbConnection, const_util.LabourManagementAttendanceTable, attendanceInterface.Id, updateObject)
			if err != nil {
				// ideally all the shift employees should check out in complete shift, but log the error
				v.Logger.Error("error updating checkout time", zap.String("error", err.Error()))
			} else {
				v.Logger.Info("successfully checked to user id", zap.Int("user_id", attendanceInterface.Id))
			}
			listOfLinesCheckedIn = append(listOfLinesCheckedIn, attendanceInfo.ManufacturingLines...)
		}

		// remove duplicates
		listOfLinesCheckedIn = util.RemoveDuplicateInt(listOfLinesCheckedIn)
		v.Logger.Info("list of lines checked in for checked out operations", zap.Any("lines", listOfLinesCheckedIn))
		// now reduce the actual manpower
		var shiftProMasterCondition = " object_info->>'$.shiftId'=" + strconv.Itoa(recordId)
		err, listOfShiftMasterProduction := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftProMasterCondition)
		machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
		// now go through each shift master production, and see which line operator checked in, based online we can reduce the man power to that shift master production
		for _, shiftMasterInterface := range listOfShiftMasterProduction {
			shiftMasterProductionInfo := database.GetShiftMasterProductionInfo(shiftMasterInterface.ObjectInfo)
			err, assemblymanInfo := machineService.GetAssemblyMachineInfoById(shiftMasterProductionInfo.MachineId)
			if err == nil {
				assemblyLineId := GetAssemblyMachineMasterInfo(assemblymanInfo.ObjectInfo).AssemblyLineOption
				if util.HasInt(assemblyLineId, listOfLinesCheckedIn) {
					v.Logger.Info("line found, so checking out the user actual man power...")
					shiftMasterProductionInfo.ActualManpower = shiftMasterProductionInfo.ActualManpower - 1
					var updatingData = make(map[string]interface{})
					updatingData["object_info"] = shiftMasterProductionInfo.Serialised()
					err := database.Update(v.Database, const_util.LabourManagementShiftProductionTable, shiftMasterInterface.Id, updatingData)
					if err != nil {
						v.Logger.Error("error updating actual man power", zap.String("error", err.Error()))
					} else {
						v.Logger.Info("successfully updated actual man power", zap.Int("actual", shiftMasterProductionInfo.ActualManpower))
					}
				}
			} else {
				v.Logger.Error("error getting assembly info", zap.String("error", err.Error()))
			}
		}
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Successfully Checkout",
		})
	} else {
		v.Logger.Warn("invalid operation, this shift is not active at this moment to checkout", zap.Int("shift_id", checkOutRequest.ShiftId))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid operation",
				Description: "This shift is not active at this moment to checkout, please contact system admin",
			})
		return
	}

}
