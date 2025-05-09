package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"encoding/json"
	"strconv"

	"gorm.io/gorm"
)

func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}

func (v *ITService) getStatusName(dbConnection *gorm.DB, serviceStatusId int) string {
	_, requestStatus := database.Get(dbConnection, const_util.ITServiceRequestStatusTable, serviceStatusId)
	var requestStatusInfo database.ITServiceRequestStatusInfo
	json.Unmarshal(requestStatus.ObjectInfo, &requestStatusInfo)
	return requestStatusInfo.Status
}

func isDoHConfigured(userId int) bool {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	listOfHeads := authService.GetHeadOfDepartments(userId)
	if len(listOfHeads) == 0 {
		return false
	}
	return true
}
