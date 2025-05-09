package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (v *ActionService) CheckInRequest(ctx *gin.Context) {
	v.Logger.Info("processing check in request")
	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	err, mouldTestRequestGeneralObject := database.Get(v.Database, const_util.MouldTestRequestTable, recordId)
	if err != nil {
		v.Logger.Error("check in request failed due to resource not found", zap.Error(err))
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	mouldTestRequest := database.MouldTestRequest{ObjectInfo: mouldTestRequestGeneralObject.ObjectInfo}
	mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()

	if mouldTestRequestInfo.ActionStatus == const_util.MouldTestRequestActionScheduled {
		mouldTestRequestInfo.ActionStatus = const_util.MouldTestRequestActionTestInProgress
		err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, mouldTestRequestInfo.DatabaseSerialize(userId))
		if err != nil {
			v.Logger.Error("update mould test request has failed", zap.Error(err))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Internal system error",
					Description: "Internal system error starting test request, please summit error code to system admin",
				})
			return
		}

		// now create the test machine teset parameter
		machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
		_, machineParamId := machineService.CreateMachineParam(projectId, mouldTestRequestGeneralObject.Id, userId)
		mouldTestRequestInfo.MachineParamId = machineParamId
		mouldTestRequestInfo.CanCheckIn = false
		mouldTestRequestInfo.CanCheckOut = true
		mouldTestRequestInfo.CanContinueTest = true
		mouldTestRequestInfo.CanUnRelease = false

		mouldTestRequestInfo.TestedBy = userId
		mouldTestRequestInfo.CanView = true

		existingActionRemarks := mouldTestRequestInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "TAKEN FOR TEST",
			UserId:        userId,
			Remarks:       "Mould test request is taken for test",
			ProcessedTime: GetTimeDifference(mouldTestRequestInfo.CreatedAt),
		})
		mouldTestRequestInfo.ActionRemarks = existingActionRemarks

		err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, mouldTestRequestInfo.DatabaseSerialize(userId))

		if err != nil {
			v.Logger.Error("update mould test request failed after set all fields", zap.Error(err))
		}
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Message: "Mould test request is successfully started",
			Code:    0,
		})
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal system error",
				Description: "Please wait until your test request get scheduled. Testing is only possible after successfully scheduling confirmed by planner",
			})
		return
	}

}
