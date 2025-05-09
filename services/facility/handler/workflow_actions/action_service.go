package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"cx-micro-flake/services/facility/handler/notification"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type ActionService struct {
	Logger       *zap.Logger
	Database     *gorm.DB
	EmailHandler *notification.EmailHandler
}

func (v *ActionService) Init(logger *zap.Logger) {
	v.Logger = logger
}
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
func (v *ActionService) getBasicInfo(ctx *gin.Context) (int, *gorm.DB, common.UserBasicInfo, map[string]interface{}) {
	recordId := util.GetRecordId(ctx)
	userId := common.GetUserId(ctx)
	_, serviceRequestGeneralObject := database.Get(v.Database, const_util.FacilityServiceRequestTable, recordId)
	serviceRequestInfo := make(map[string]interface{})
	json.Unmarshal(serviceRequestGeneralObject.ObjectInfo, &serviceRequestInfo)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)

	return recordId, v.Database, basicUserInfo, serviceRequestInfo
}

func (v *ActionService) getStatusName(dbConnection *gorm.DB, serviceStatusId int) string {
	_, requestStatus := database.Get(dbConnection, const_util.FacilityServiceRequestStatusTable, serviceStatusId)
	var requestStatusInfo database.FacilityServiceRequestStatusInfo
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
