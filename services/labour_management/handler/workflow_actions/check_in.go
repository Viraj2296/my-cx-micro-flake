package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CheckInRequest struct {
	EmployeeId string `json:"employeeId"`
	Lines      []int  `json:"lines"`
	ShiftId    int    `json:"shiftId"`
}

func (v *ActionService) HandleCheckIn(ctx *gin.Context) {
	v.Logger.Info("handle check in received")
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
	v.Logger.Info("check in is creating by", zap.Int("user_info", userInfo.UserId))
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
		// check the supervisor is already checked in
		var condition = " object_info->>'$.shiftId' = " + strconv.Itoa(recordId) + " AND object_info->>'$.userResourceId' =" + strconv.Itoa(checkInUserInfo.UserId)

		listOfAttendance := database.CountByCondition(dbConnection, const_util.LabourManagementAttendanceTable, condition)

		if shiftMasterInfo.IsSupervisorCheckedIn {
			// only allowing once it is checked in.
			if listOfAttendance == 0 {
				attendance := database.AttendanceInfo{
					ShiftResourceId:    checkInRequest.ShiftId,
					CheckInDate:        util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
					CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
					LastUpdatedAt:      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
					CheckInTime:        util.GetCurrentTimeInSingapore("15:04:05"),
					CreatedBy:          userInfo.UserId,
					UserResourceId:     checkInUserInfo.UserId,
					CheckOutDate:       "",
					CheckOutTime:       "",
					ManufacturingLines: checkInRequest.Lines,
				}
				// try to get the scheduled order event id from the line id
				err := v.updateShiftProductionManpower(ctx, checkInRequest.ShiftId, checkInRequest.Lines, checkInUserInfo.EmployeeTypeId)
				if err != nil {
					v.Logger.Error("error updating shift calculation", zap.String("error", err.Error()))
					return
				}
				generalObject := component.GeneralObject{ObjectInfo: attendance.Serialised()}
				err, attendanceRecordId := database.CreateFromGeneralObject(dbConnection, const_util.LabourManagementAttendanceTable, generalObject)
				if err == nil {
					v.Logger.Info("successfully checked in", zap.String("full-name", checkInUserInfo.FullName), zap.Int("resource_id", attendanceRecordId))

					// Prepare ShiftProductionAttendanceInfo object
					shiftProductionAttendanceInfo := database.LabourManagementShiftProductionAttendance{
						CreatedBy:         checkInUserInfo.UserId,
						ShiftProductionId: checkInRequest.ShiftId,
						ShiftAttendanceId: attendanceRecordId,
					}

					err = database.CreateShiftProductionAttendanceLink(dbConnection, shiftProductionAttendanceInfo)
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
					return
				} else {
					response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
						&response.DetailedError{
							Header:      "Server Exception",
							Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
						})
					return
				}
			} else {
				v.Logger.Error("user has already checked in", zap.String("full_name", checkInUserInfo.FullName))
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Invalid Checkin",
						Description: "The employee [" + checkInUserInfo.FullName + "] has already checked in",
					})
				return
			}
		} else {

			// check this is the supervisor
			// don't look for supervisor, if the roles configured just start it.
			if util.HasInt(checkInUserInfo.JobRole, v.LabourManagementSettingInfo.ShiftCreationRoles) {
				//if shiftMasterInfo.ShiftSupervisor == checkInUserInfo.UserId {
				// yes, supervisor checking in update the status
				var updatingData = make(map[string]interface{})
				shiftMasterInfo.IsSupervisorCheckedIn = true
				rawUserInfo, _ := json.Marshal(shiftMasterInfo)
				updatingData["object_info"] = rawUserInfo
				err := database.Update(dbConnection, const_util.LabourManagementShiftMasterTable, checkInRequest.ShiftId, updatingData)
				if err != nil {
					v.Logger.Error("error updating shift master data", zap.String("error", err.Error()))
				} else {
					attendance := database.AttendanceInfo{
						ShiftResourceId:    checkInRequest.ShiftId,
						CheckInDate:        util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
						CheckInTime:        util.GetCurrentTimeInSingapore("15:04:05"),
						CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
						LastUpdatedAt:      util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
						UserResourceId:     checkInUserInfo.UserId,
						CheckOutDate:       "",
						CheckOutTime:       "",
						ManufacturingLines: checkInRequest.Lines,
					}
					err := v.updateShiftProductionManpower(ctx, checkInRequest.ShiftId, checkInRequest.Lines, checkInUserInfo.EmployeeTypeId)
					if err != nil {
						v.Logger.Error("error updating shift calculation", zap.String("error", err.Error()))
						return
					}
					generalObject := component.GeneralObject{ObjectInfo: attendance.Serialised()}
					err, attendanceRecordId := database.CreateFromGeneralObject(dbConnection, const_util.LabourManagementAttendanceTable, generalObject)
					if err == nil {
						v.Logger.Info("successfully checked in", zap.String("fullname", checkInUserInfo.FullName), zap.Int("resource_id", attendanceRecordId))

						// Prepare ShiftProductionAttendanceInfo object
						shiftProductionAttendanceInfo := database.LabourManagementShiftProductionAttendance{
							CreatedBy:         checkInUserInfo.UserId,
							ShiftProductionId: checkInRequest.ShiftId,
							ShiftAttendanceId: attendanceRecordId,
						}

						err = database.CreateShiftProductionAttendanceLink(dbConnection, shiftProductionAttendanceInfo)
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

			} else {
				v.Logger.Error("supervisor not checked in", zap.Any("supervisor", shiftMasterInfo.ShiftSupervisor))
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Supervisor not checked in",
						Description: "The supervisor for this shift has not checked in. This typically happens when your role is not authorized to perform the initial check-in.",
					})
				return
			}
		}
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Shift",
				Description: "This shift is not valid for this operation",
			})
		return
	}
}

func (v *ActionService) updateShiftProductionManpower(ctx *gin.Context, shiftId int, lines []int, employeeTypeId int) error {
	for _, line := range lines {
		schedulerEventId := v.getSchedulerEventIdFromMachineLine(line, shiftId)
		if schedulerEventId == 0 {
			v.Logger.Error("error getting scheduler event ID")
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Line",
					Description: "The line you selected is not configured in the shift, please make sure you are selected correct line during check-in",
				})
			return errors.New("error getting scheduler event ID")
		}
		var shiftProMasterCondition = " object_info->>'$.shiftId'=" + strconv.Itoa(shiftId) + " AND object_info->>'$.scheduledEventId'=" + strconv.Itoa(schedulerEventId)
		// now update the shift production
		err, objects := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftProMasterCondition)
		if err != nil {
			v.Logger.Error("get objects error", zap.Error(err))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Shift Production Resource",
					Description: "System couldn't able to identify the shift production resource based on selection criteria",
				})
			return errors.New("error getting shift production resource based on selection criteria")
		}
		if len(objects) == 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Shift Production Found",
					Description: "System couldn't able to identify the shift production resource based on selection criteria",
				})
			return errors.New("error getting shift production resource based on selection criteria")
		}
		shiftMasterProductionInfo := database.GetShiftMasterProductionInfo((objects)[0].ObjectInfo)
		if employeeTypeId == v.LabourManagementSettingInfo.ContractEmployeeTypeId {
			// this user is the contract so update the count in parttimer counter
			//TODO if the user is already checked in, then don't increase the counter
			shiftMasterProductionInfo.ShiftActualOutputPartTimer += 1
		}
		shiftMasterProductionInfo.ActualManpower += 1
		var updatingData = make(map[string]interface{})
		updatingData["object_info"] = shiftMasterProductionInfo.Serialised()
		err = database.Update(v.Database, const_util.LabourManagementShiftProductionTable, (objects)[0].Id, updatingData)
		if err != nil {
			v.Logger.Error("update error", zap.Error(err))
			return errors.New("error updating shift production resource based on selection criteria")
		}
	}
	return nil
}
