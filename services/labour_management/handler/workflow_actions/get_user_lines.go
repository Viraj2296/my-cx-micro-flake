package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

type LineSelectionResponse struct {
	ListOfAssemblyLines    []ListOfAssemblyLine `json:"listOfAssemblyLines"`
	CanSelectMultipleLines bool                 `json:"canSelectMultipleLines"`
}

type UserLineRequest struct {
	EmployeeId string `json:"employeeId"`
}

// GetUserLines This function to get the list lines
func (v *ActionService) GetUserLines(ctx *gin.Context) {
	v.Logger.Info("handle GetUserLines")
	checkInRequest := UserLineRequest{}
	var shiftId = util.GetRecordId(ctx)
	if err := ctx.ShouldBindBodyWith(&checkInRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("sending error has failed", zap.Error(err))
		}
		return
	}
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
	var lineSelectionResponse = LineSelectionResponse{}
	lineSelectionResponse.ListOfAssemblyLines = make([]ListOfAssemblyLine, 0)
	lineSelectionResponse.CanSelectMultipleLines = true
	if authService.GetUserInfoById(checkInUserInfo.UserId).JobRole == v.LabourManagementSettingInfo.ShiftOperatorJobRoleId {
		lineSelectionResponse.CanSelectMultipleLines = false
	}
	lineSelectionResponse.ListOfAssemblyLines = v.getLinesForSchedulerEvents(shiftId)
	v.Logger.Info("got the list of approved lines to shift", zap.Int("shift_id", shiftId), zap.Any("assemblies", lineSelectionResponse.ListOfAssemblyLines))
	ctx.JSON(http.StatusOK, lineSelectionResponse)
}
