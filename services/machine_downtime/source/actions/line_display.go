package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"net/http"
	"time"

	"go.cerex.io/transcendflow/orm"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserInfo struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}

type MachineInfo struct {
	Name string `json:"name"`
}

type DowntimeEntry struct {
	Machine          MachineInfo `json:"machine"`
	CheckInDateTime  string      `json:"checkInDateTime"`
	CheckOutDateTime string      `json:"checkOutDateTime"`
	CreatedDateTime  string      `json:"createdDateTime"`
	Status           string      `json:"status"`
	StatusColorCode  string      `json:"statusColorCode"`
	CheckInUser      UserInfo    `json:"checkInUser"`
	CheckOutUser     UserInfo    `json:"checkOutUser"`
	AssignedUser     UserInfo    `json:"assignedUser"`
}

// GetMachineDowntimeDisplay handles the endpoint for fetching machine downtime
func (v *Actions) GetMachineDowntimeDisplay(ctx *gin.Context) {
	v.Logger.Info("getting the machine downtime display...")

	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	err, machineDowntimeRecords := orm.GetConditionalObjects(v.Database, consts.MachineDownTimeMasterTable, " object_info->>'$.checkoutDate' IS NULL and (object_info->>'$.status' =1  OR object_info->>'$.status' =2) order by object_info->>'$.createdAt'")
	if err != nil {
		v.Logger.Error("error retrieving machine downtime records", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus, &response.DetailedError{
			Header:      "Internal Server Error",
			Description: "Unable to fetch machine downtime data, please contact the system administrator.",
		})
		return
	}
	v.Logger.Info("machine downtime records retrieved successfully", zap.Int("record_count", len(machineDowntimeRecords)))
	var downtimeStatusMap = make(map[int]*models.MachineDowntimeStatusInfo)
	err, i := orm.GetObjects(v.Database, consts.MachineDownTimeStatusTable)
	if err == nil {
		for _, record := range i {
			downtimeStatusMap[record.Id] = models.GetMachineDowntimeStatusInfo(record.ObjectInfo)
		}
	} else {
		v.Logger.Error("error retrieving machine downtime records", zap.String("error", err.Error()))
	}

	var downtimeEntries []DowntimeEntry
	// Process each downtime record
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	for _, record := range machineDowntimeRecords {
		downtimeInfo := models.GetMachineDowntimeInfo(record.ObjectInfo)
		if downtimeInfo == nil {
			v.Logger.Error("Error getting machine downtime info", zap.Any("record", record))
			continue
		}

		err, machineGeneralComponent := machineService.GetAssemblyMachineInfoById(downtimeInfo.MachineId)
		if err != nil {
			v.Logger.Error("error retrieving machine info", zap.Error(err))
			continue
		}
		if assemblyMachineInfo, err := dto.GetAssemblyMachineMasterInfo(machineGeneralComponent.ObjectInfo); err == nil {
			// Fetch user information for check-in and check-out users
			checkInUser := authService.GetUserInfoById(downtimeInfo.CheckInUserId)
			checkOutUser := authService.GetUserInfoById(downtimeInfo.CheckOutUserId)
			assignedUser := authService.GetUserInfoById(downtimeInfo.AssignedUserId)
			if downtimeInfo.CheckOutUserId == 0 {
				checkOutUser.FullName = "-"
			}
			if downtimeInfo.CheckInUserId == 0 {
				checkOutUser.FullName = "-"
			}
			var downtimeStatus = ""
			var defaultStatusColorCode = "#bfd3df "
			if statusValue, ok := downtimeStatusMap[downtimeInfo.Status]; ok {
				downtimeStatus = statusValue.Status
				defaultStatusColorCode = statusValue.ColorCode
			}

			checkInDateTime, _ := time.Parse("2006-01-02T15:04:05.000Z", downtimeInfo.CheckInDate)
			checkOutDateTime, _ := time.Parse("2006-01-02T15:04:05.000Z", downtimeInfo.CheckOutDate)
			createdDateTime, _ := time.Parse("2006-01-02T15:04:05.000Z", downtimeInfo.CreatedAt)

			// Adjust the time by adding 8 hours
			var updatedCheckOutDateTime time.Time
			var updatedCheckInDateTime time.Time
			if downtimeInfo.CheckOutDate == "" {
				updatedCheckOutDateTime = checkOutDateTime
			} else {
				updatedCheckOutDateTime = checkOutDateTime.Add(8 * time.Hour)
			}
			if downtimeInfo.CheckInDate == "" {
				updatedCheckInDateTime = checkInDateTime
			} else {
				updatedCheckInDateTime = checkInDateTime.Add(8 * time.Hour)
			}

			var updatedCreatedDateTime = createdDateTime.Add(8 * time.Hour)
			downtimeEntry := DowntimeEntry{
				Machine: MachineInfo{
					Name: assemblyMachineInfo.Description,
				},

				CheckInDateTime:  updatedCheckInDateTime.Format("Jan 02, 2006, 03:04:05 PM"),
				CheckOutDateTime: updatedCheckOutDateTime.Format("Jan 02, 2006, 03:04:05 PM"),
				CreatedDateTime:  updatedCreatedDateTime.Format("Jan 02, 2006, 03:04:05 PM"),
				Status:           downtimeStatus,
				StatusColorCode:  defaultStatusColorCode,

				CheckInUser: UserInfo{
					Name:      checkInUser.FullName,
					AvatarUrl: checkInUser.AvatarUrl,
				},
				CheckOutUser: UserInfo{
					Name:      checkOutUser.FullName,
					AvatarUrl: checkOutUser.AvatarUrl,
				},
				AssignedUser: UserInfo{
					Name:      assignedUser.FullName,
					AvatarUrl: assignedUser.AvatarUrl,
				},
			}
			downtimeEntries = append(downtimeEntries, downtimeEntry)
		}

	}

	if len(downtimeEntries) == 0 {
		v.Logger.Warn("no downtime entries found without checkout date")
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No downtime entries found without checkout date"})
		return
	}

	ctx.JSON(http.StatusOK, downtimeEntries)
}
