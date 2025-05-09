package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Assignment struct {
	Event    int `json:"int"`
	Id       int `json:"id"`
	Resource int `json:"resource"`
}

type OrderScheduledEventUpdateRequest struct {
	EndDate    string      `json:"endDate"`
	StartDate  string      `json:"startDate"`
	Assignment *Assignment `json:"assignment"`
}

func (v *ActionService) UpdateSchedulerEvent(ctx *gin.Context) {
	var updateScheduleEventRequest = OrderScheduledEventUpdateRequest{}
	if err := ctx.ShouldBindBodyWith(&updateScheduleEventRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	recordId := util.GetRecordId(ctx)

	err, eventObject := database.Get(v.Database, const_util.MouldTestRequestTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}

	userId := common.GetUserId(ctx)
	mouldTestRequest := database.MouldTestRequest{ObjectInfo: eventObject.ObjectInfo}
	mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()

	orderInfoStartDate := mouldTestRequestInfo.RequestTestStartDate
	orderInfoEndDate := mouldTestRequestInfo.RequestTestEndDate
	mouldTestRequestInfo.RequestTestStartDate = updateScheduleEventRequest.StartDate
	mouldTestRequestInfo.RequestTestEndDate = updateScheduleEventRequest.EndDate

	currentTime := time.Now().Unix()
	reqStart, _ := time.Parse(time.RFC3339, updateScheduleEventRequest.StartDate)
	reqEnd, _ := time.Parse(time.RFC3339, updateScheduleEventRequest.EndDate)
	dbStart, _ := time.Parse(time.RFC3339, util.InterfaceToString(orderInfoStartDate))
	dbEnd, _ := time.Parse(time.RFC3339, util.InterfaceToString(orderInfoEndDate))

	if updateScheduleEventRequest.Assignment != nil {
		if updateScheduleEventRequest.Assignment.Resource != mouldTestRequestInfo.MachineId {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError(const_util.InvalidSourceError), const_util.InvalidScheduleStatus, "Changing resource which already assigned is not permitted, consider moving along with the resources")
			return
		}
	}

	if !reqStart.Equal(dbStart) {

		if currentTime > dbStart.Unix() {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError(const_util.InvalidSourceError), const_util.InvalidScheduleStatus, "You are trying to modify the event which is already passed, This modification is rejected")
			return
		}
	} else if !reqEnd.Equal(dbEnd) {

		if currentTime > dbEnd.Unix() {

			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError(const_util.InvalidSourceError), const_util.InvalidScheduleStatus, "You are trying to modify the event which is already passed, This modification is rejected")
			return
		}
	}

	err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, mouldTestRequestInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, const_util.GetError("error updating machines event information"), const_util.ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})
}
