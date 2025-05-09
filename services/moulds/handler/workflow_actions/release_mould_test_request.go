package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) ReleaseMouldTestRequest(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)

	// send to all in the group
	err, c := database.Get(v.Database, const_util.MouldSettingTable, 1)
	if err != nil {
		v.Logger.Error("error getting setting", zap.Error(err))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("no testing group configured to proceed the testing, please configure"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	err, eventObject := database.Get(v.Database, const_util.MouldTestRequestTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	mouldTestRequest := database.MouldTestRequest{ObjectInfo: eventObject.ObjectInfo}
	mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEventByMachineId(const_util.ProjectID, mouldTestRequestInfo.MachineId, mouldTestRequestInfo.RequestTestStartDate, mouldTestRequestInfo.RequestTestEndDate)

	if scheduledEventObject != nil {
		if len(*scheduledEventObject) > 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      const_util.GetError(const_util.DuplicateRecordFound).Error(),
					Description: "The Mould Test Request cannot be released as the requested machine ID is already scheduled for a production order.",
				})
			return
		}
	}

	mouldTestRequestInfo.CanCheckIn = true
	mouldTestRequestInfo.CanRelease = false
	mouldTestRequestInfo.CanUnRelease = true
	mouldTestRequestInfo.CanView = false
	mouldTestRequestInfo.MouldTestStatus = const_util.MouldTestWorkFlowProcessDepartment
	mouldTestRequestInfo.ActionStatus = const_util.MouldTestRequestActionScheduled
	err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, mouldTestRequestInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	mouldSettingInfo := database.GetMouldSettingInfo(c.ObjectInfo)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	listOfUsers := authService.GetUserInfoFromGroupId(mouldSettingInfo.MouldTestingGroup)
	for _, user := range listOfUsers {
		v.Logger.Info("email is sending for mould testing", zap.Int("user_id", user.UserId))
		v.EmailHandler.EmailGenerator(v.Database, const_util.MouldTestRequestScheduledTemplateType, user.UserId, const_util.MouldTestRequestComponent, recordId)
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Mould Test Request has been successfully released",
		Code:    0,
	})

}
