package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

type ValidateCheckInResponse struct {
	CanCheckIn   bool   `json:"canCheckIn"`
	AlertMessage string `json:"alertMessage"`
}

func (v *ActionService) ValidateCheckIn(ctx *gin.Context) {
	v.Logger.Info("handle ValidateCheckIn received")
	checkInRequest := CheckInRequest{}
	if err := ctx.ShouldBindBodyWith(&checkInRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	recordId := checkInRequest.ShiftId
	// Call getBasicInfo with recordId
	_, dbConnection, userInfo, shiftMasterInfo := v.getShiftMasterInfo(ctx, recordId)
	v.Logger.Info("ValidateCheckIn is creating by", zap.Int("user_info", userInfo.UserId))
	if shiftMasterInfo.ShiftStatus == const_util.ShiftStatusActive {
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		v.Logger.Info("requesting user details for employee id", zap.String("employeeId", checkInRequest.EmployeeId))
		checkInUserInfo := authService.GetUserInfoByEmployeeId(checkInRequest.EmployeeId)
		v.Logger.Info("checkin user info", zap.String("username", checkInUserInfo.Username))
		if checkInUserInfo.UserId == 0 {
			v.Logger.Error("invalid user id", zap.Int("user_id", checkInUserInfo.UserId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid User",
					Description: "The user is not available in the system, please contact system admin to configure the user",
				})
			return
		}
		// check the user has already checked in for line
		var condition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(recordId) + " AND object_info->>'$.userResourceId' =" + strconv.Itoa(checkInUserInfo.UserId) + " AND object_info->>'$.checkOutDate' = ''"

		err, listOfAttendance := database.GetConditionalObjects(dbConnection, const_util.LabourManagementAttendanceTable, condition)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Internal System Error",
					Description: "System couldn't able get the meta information, please report to system admin",
				})
			return
		}
		if len(listOfAttendance) == 0 {
			v.Logger.Info("no users checked in.... so we can check in the user now")
			ctx.JSON(http.StatusOK, ValidateCheckInResponse{CanCheckIn: true, AlertMessage: ""})
		} else {
			attendanceInfo := database.GetAttendanceInfo(listOfAttendance[0].ObjectInfo)
			if len(attendanceInfo.ManufacturingLines) > 1 {
				// this should be the supervisor
				v.Logger.Info("no users checked in.... so we can check in the user now")
				ctx.JSON(http.StatusOK, ValidateCheckInResponse{CanCheckIn: true, AlertMessage: ""})
			} else {
				var checkedInLine = attendanceInfo.ManufacturingLines[0]
				v.Logger.Info("check in line ", zap.Int("checked_in_line", checkedInLine), zap.Int("request_line", checkInRequest.Lines[0]))
				if checkInRequest.Lines[0] == checkedInLine {
					ctx.JSON(http.StatusOK, ValidateCheckInResponse{CanCheckIn: false, AlertMessage: "User has already checked-in for this line, Please select the different line"})
				} else {
					ctx.JSON(http.StatusOK, ValidateCheckInResponse{CanCheckIn: true, AlertMessage: "User has already checked-in for another line, proceeding with this line would make the user checked-out to previous line, Would you like to proceed ?"})
				}
			}
		}
	}
}
