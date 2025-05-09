package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LabourInfo struct {
	Role           string `json:"role"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatarUrl"`
	HierarchyLevel int    `json:"hierarchyLevel"`
}

type SummaryResponse struct {
	ShiftTemplateName        string        `json:"shiftTemplateName"`
	ShiftTemplateDescription string        `json:"shiftTemplateDescription"`
	ShiftStartDateTime       string        `json:"shiftStartDateTime"`
	ShiftEndDateTime         string        `json:"shiftEndDateTime"`
	ShiftActualStartDateTime string        `json:"shiftActualStartDateTime"`
	LabourInfo               []RoleLabours `json:"labourInfo"`
}

type RoleLabours struct {
	RoleName       string       `json:"roleName"`
	Labours        []LabourInfo `json:"labours"`
	HierarchyLevel int          `json:"hierarchyLevel"`
}

func (v *ActionService) GetaAssemblyShiftSummaryDisplay(ctx *gin.Context) {
	v.Logger.Info("getting the assembly labour display")

	var conditionString = " object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusActive)
	err, listOfActiveShifts := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, conditionString)
	if err != nil {
		// it is the error send the error description
		v.Logger.Error("error getting shift master", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "Sorry, server couldn't able to complete this action due to internal error, please report this error to system admin",
			})
		return
	}
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	// identify what shift is running on requested line
	// it is better to init the default values
	// it should send the empty response instead of null.
	var labourInfo = make([]LabourInfo, 0)
	for _, activeShiftInterface := range listOfActiveShifts {

		var attendanceCondition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(activeShiftInterface.Id)
		v.Logger.Error("line found, getting the attendance", zap.Any("attendanceCondition", attendanceCondition))
		err, listOfAttendance := database.GetConditionalObjects(v.Database, const_util.LabourManagementAttendanceTable, attendanceCondition)
		if err != nil {
			// add the empty array
			v.Logger.Error("error getting operators", zap.String("error", err.Error()))
			labourInfo = make([]LabourInfo, 0)
		}
		var uniqueUserList []int
		for _, info := range listOfAttendance {
			attendanceInfo := database.GetAttendanceInfo(info.ObjectInfo)
			uniqueUserList = append(uniqueUserList, attendanceInfo.UserResourceId)
		}
		uniqueUserList = util.RemoveDuplicateInt(uniqueUserList)
		for _, userId := range uniqueUserList {
			userInfo := authService.GetUserInfoById(userId)
			// only configured based on configuration
			if util.HasInt(userInfo.JobRole, v.LabourManagementSettingInfo.LabourManagementSummaryJobRoles) {
				var jobRoleName = authService.GetJobRoleName(userInfo.JobRole)
				labourInfo = append(labourInfo, LabourInfo{
					AvatarUrl: userInfo.AvatarUrl,
					Name:      userInfo.FullName,
					Role:      jobRoleName,
				})
			} else {
				v.Logger.Warn("skipping adding user as it is not configured as summary generator job role", zap.Int("job_role", userInfo.JobRole))
			}

		}
	}
	ctx.JSON(http.StatusOK, labourInfo)
}

// GetaAssemblyShiftSummaryDisplayV2 Note : https://cerexio.atlassian.net/browse/FUYU2-352
// This is to support new way of showing the labour management summary
func (v *ActionService) GetaAssemblyShiftSummaryDisplayV2(ctx *gin.Context) {
	v.Logger.Info("getting the assembly labour display")

	// Query to get the active shift
	var conditionString = " object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusActive) + " ORDER BY object_info->>'$.createdAt' DESC"
	err, activeShiftLists := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, conditionString)
	if err != nil {
		// Log the error and return a detailed error response
		v.Logger.Error("error getting shift master", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "Sorry, server couldn't able to complete this action due to internal error, please report this error to system admin",
			})
		return
	}

	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	// Initialize the response
	summaryResponse := SummaryResponse{
		ShiftTemplateName:        "",
		ShiftTemplateDescription: "",
		ShiftStartDateTime:       "",
		ShiftEndDateTime:         "",
		ShiftActualStartDateTime: "",
		LabourInfo:               []RoleLabours{},
	}

	// If no active shifts, return the empty response
	if len(activeShiftLists) == 0 {
		ctx.JSON(http.StatusOK, summaryResponse)
		return
	}
	var uniqueUserList []int

	for _, activeShiftObject := range activeShiftLists {
		// Extract shift master info
		shiftMasterInfo := database.GetShiftMasterInfo(activeShiftObject.ObjectInfo)

		// Get shift template info
		err, shiftTemplateInterface := database.Get(v.Database, const_util.LabourManagementShiftTemplateTable, shiftMasterInfo.ShiftTemplateId)
		if err == nil {
			shiftTemplateInfo := database.GetSShiftTemplateInfo(shiftTemplateInterface.ObjectInfo)
			summaryResponse.ShiftTemplateName = shiftTemplateInfo.Name
			summaryResponse.ShiftTemplateDescription = shiftTemplateInfo.Description

			// Combine date and time
			startDateTime, _ := CombineDateTimeToISO8601(shiftMasterInfo.ShiftStartDate, shiftMasterInfo.ShiftStartTime)
			actualStartDateTime, _ := ConvertToISO8601NoMilliseconds(shiftMasterInfo.CreatedAt)
			summaryResponse.ShiftStartDateTime = startDateTime
			summaryResponse.ShiftActualStartDateTime = actualStartDateTime

			// Calculate shift end time
			err, endTime := v.GetShiftEndTime(shiftMasterInfo.ShiftStartDate, shiftTemplateInfo.ShiftStartTime, shiftTemplateInfo.ShiftPeriod)
			if err != nil {
				v.Logger.Error("error getting shift end time", zap.String("error", err.Error()))
				summaryResponse.ShiftEndDateTime = "-"
			} else {
				summaryResponse.ShiftEndDateTime = endTime
			}

			// Get attendance info
			var attendanceCondition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(activeShiftObject.Id)
			v.Logger.Info("line found, getting the attendance", zap.Any("attendanceCondition", attendanceCondition))
			err, listOfAttendance := database.GetConditionalObjects(v.Database, const_util.LabourManagementAttendanceTable, attendanceCondition)
			if err != nil {
				v.Logger.Error("error getting operators", zap.String("error", err.Error()))
			}

			for _, info := range listOfAttendance {
				attendanceInfo := database.GetAttendanceInfo(info.ObjectInfo)
				if attendanceInfo.CheckOutDate == "" || attendanceInfo.CheckOutTime == "" {
					uniqueUserList = append(uniqueUserList, attendanceInfo.UserResourceId)
				} else {
					v.Logger.Warn("user has already checked out, so removing it from showing in the shift summary display", zap.Any("attendanceInfo", attendanceInfo))
				}
			}

		}
	}
	uniqueUserList = util.RemoveDuplicateInt(uniqueUserList)

	// Group users by their roles
	roleLaboursMap := make(map[string]RoleLabours)
	for _, userId := range uniqueUserList {
		userInfo := authService.GetUserInfoById(userId)

		// Only include configured roles
		if util.HasInt(userInfo.JobRole, v.LabourManagementSettingInfo.LabourManagementSummaryJobRoles) {
			jobRoleName := authService.GetJobRoleName(userInfo.JobRole)

			// Retrieve the hierarchy level for the job role (implement this method)
			hierarchyLevel := authService.GetJobRoleHierarchy(userInfo.JobRole) // Assuming this returns int

			// Check if the role is already present in the map
			if roleLabour, exists := roleLaboursMap[jobRoleName]; exists {
				// Append the labour info to the existing role
				roleLabour.Labours = append(roleLabour.Labours, LabourInfo{
					AvatarUrl:      userInfo.AvatarUrl,
					Name:           userInfo.FullName,
					Role:           jobRoleName,
					HierarchyLevel: hierarchyLevel,
				})
				roleLaboursMap[jobRoleName] = roleLabour // Update the map with modified roleLabour
			} else {
				// Create a new entry in the map
				roleLaboursMap[jobRoleName] = RoleLabours{
					RoleName:       jobRoleName,
					HierarchyLevel: hierarchyLevel,
					Labours: []LabourInfo{
						{
							AvatarUrl:      userInfo.AvatarUrl,
							Name:           userInfo.FullName,
							Role:           jobRoleName,
							HierarchyLevel: hierarchyLevel,
						},
					},
				}
			}
		} else {
			v.Logger.Warn("skipping user not in summary job role", zap.Int("job_role", userInfo.JobRole))
		}
	}

	// Convert the map to the desired structure
	for _, roleLabour := range roleLaboursMap {
		summaryResponse.LabourInfo = append(summaryResponse.LabourInfo, roleLabour)
	}

	// Sort the LabourInfo by HierarchyLevel
	sort.Slice(summaryResponse.LabourInfo, func(i, j int) bool {
		return summaryResponse.LabourInfo[i].HierarchyLevel < summaryResponse.LabourInfo[j].HierarchyLevel
	})

	// Return the summary response
	ctx.JSON(http.StatusOK, summaryResponse)
}
