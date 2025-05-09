package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
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

type UpdateManpower struct {
	EmployeeId  string `json:"employeeId"`
	CheckInTime string `json:"checkInTime"` //2024-10-04 10:30:00”
	Remarks     string `json:"remarks"`
	Lines       []int  `json:"lines"`
}

func (v *ActionService) HandleUpdateManpower(ctx *gin.Context) {
	v.Logger.Info("update manpower request is received")
	updateManpower := UpdateManpower{}
	userId := common.GetUserId(ctx)
	if err := ctx.ShouldBindBodyWith(&updateManpower, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	v.Logger.Info("requesting user details for employee id", zap.String("employeeId", updateManpower.EmployeeId))
	checkInUserInfo := authService.GetUserInfoByEmployeeId(updateManpower.EmployeeId)

	if checkInUserInfo.UserId == 0 {
		v.Logger.Error("invalid user id", zap.Int("user_id", checkInUserInfo.UserId))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid User",
				Description: "The user is not available in the system, please contact system admin to configure the user",
			})
		return
	}

	resourceId := util.GetRecordId(ctx)

	err, shiftProductionInterface := database.Get(v.Database, const_util.LabourManagementShiftProductionTable, resourceId)
	if err != nil {
		v.Logger.Error("error getting shift production ", zap.Error(err))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Shift Production",
				Description: "Request shift production is not available in the system",
			})
		return
	}
	shiftProductionInfo := database.GetShiftMasterProductionInfo(shiftProductionInterface.ObjectInfo)
	err, shiftMasterInterface := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, shiftProductionInfo.ShiftId)
	if err == nil {
		shiftMasterInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
		if shiftMasterInfo.ShiftStatus == const_util.ShiftStatusCompleted {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Shift Status",
					Description: "You can not modify shift production data once the shift is completed",
				})
			return
		}
		shiftProductionInfo.ActualManpower = shiftProductionInfo.ActualManpower + 1
		// now update the check in records
		existingActionRemarks := shiftProductionInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "MANUAL CHECKIN",
			UserId:        userId,
			Remarks:       "Manpower is updated manually [" + updateManpower.Remarks + "]",
			ProcessedTime: getTimeDifference(util.InterfaceToString(shiftProductionInfo.CreatedAt)),
		})

		// create the attendance record
		// Parse the input time string to time.Time object
		parsedTime, err := time.Parse("2006-01-02 15:04", updateManpower.CheckInTime)
		if err != nil {
			v.Logger.Error("error parsing checkInTime, sending internal system error ", zap.Error(err))
			response.SendInternalSystemError(ctx)
			return
		}

		// Convert to desired format with "Z" for UTC
		convertedTime := parsedTime.UTC().Format("2006-01-02T15:04:05.000Z")
		// Format the time to extract only the time part (HH:mm:ss)
		timeOnly := parsedTime.Format("15:04:05")

		attendance := database.AttendanceInfo{
			ShiftResourceId:    shiftProductionInfo.ShiftId,
			CheckInDate:        convertedTime,
			CheckInTime:        timeOnly,
			CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedAt:      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			UserResourceId:     checkInUserInfo.UserId,
			CheckOutDate:       "",
			CheckOutTime:       "",
			ManufacturingLines: updateManpower.Lines,
		}
		err = v.updateShiftProductionManpower(ctx, shiftProductionInfo.ShiftId, updateManpower.Lines, checkInUserInfo.EmployeeTypeId)
		if err != nil {
			v.Logger.Error("error updating actual manpower", zap.String("error", err.Error()))
			return
		}
		generalObject := component.GeneralObject{ObjectInfo: attendance.Serialised()}
		err, attendanceRecordId := database.CreateFromGeneralObject(v.Database, const_util.LabourManagementAttendanceTable, generalObject)
		if err == nil {
			v.Logger.Info("successfully checked in", zap.String("full_name", checkInUserInfo.FullName), zap.Int("resource_id", attendanceRecordId))
			 // Create a record in ShiftProductionAttendanceInfo after a successful check-in
			shiftProductionAttendanceInfo := database.LabourManagementShiftProductionAttendance{
				CreatedBy:         checkInUserInfo.UserId,
				ShiftProductionId: shiftProductionInfo.ShiftId,
				ShiftAttendanceId:      attendanceRecordId, // Use the ID of the newly created attendance record
			}

			err = database.CreateShiftProductionAttendanceLink(v.Database, shiftProductionAttendanceInfo)
			if err != nil {
				v.Logger.Error("error creating ShiftProductionAttendanceInfo record", zap.Error(err))
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Database Error",
						Description: "Error while creating shift production attendance info.",
					})
				return
			}

			ctx.JSON(http.StatusOK, response.GeneralResponse{
				Code:    0,
				Message: "Successfully CheckedIn",
			})
		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

	}

}
