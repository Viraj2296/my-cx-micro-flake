package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *ActionService) UnReleaseMouldTestRequest(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)

	err, eventObject := database.Get(v.Database, const_util.MouldTestRequestTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	mouldTestRequest := database.MouldTestRequest{ObjectInfo: eventObject.ObjectInfo}
	mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()
	mouldTestRequestInfo.MouldTestStatus = const_util.MouldTestWorkFlowPlanner
	mouldTestRequestInfo.ActionStatus = const_util.MouldTestRequestActionCreated
	mouldTestRequestInfo.CanUnRelease = false
	mouldTestRequestInfo.CanRelease = true
	err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, mouldTestRequestInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), const_util.ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})
}
