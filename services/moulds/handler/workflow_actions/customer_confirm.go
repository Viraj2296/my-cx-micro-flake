package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CustomerConfirmRequest struct {
	IsExport bool `json:"isExport"`
}

func (v *ActionService) CustomerConfirm(ctx *gin.Context) {
	recordId := util.GetRecordId(ctx)
	customerConfirmRequest := CustomerConfirmRequest{}
	if err := ctx.ShouldBindBodyWith(&customerConfirmRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	v.Logger.Info("customer confirm request is received ", zap.Any("request", customerConfirmRequest))
	err, eventObject := database.Get(v.Database, const_util.MouldMasterTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	var mouldMasterFields = make(map[string]interface{})
	json.Unmarshal(eventObject.ObjectInfo, &mouldMasterFields)
	var mouldStatus int
	if mouldMasterStatusInterface, ok := mouldMasterFields["mouldStatus"]; !ok {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Mould Status",
				Description: "Internal system error, system couldn't able to determine the mould status",
			})
		return
	} else {
		mouldStatus = util.InterfaceToInt(mouldMasterStatusInterface)
	}

	if mouldStatus == const_util.MouldStatusCustomerApproval {
		updatingData := make(map[string]interface{})
		if customerConfirmRequest.IsExport == true {
			mouldStatus = const_util.MouldStatusExport
		} else {
			mouldStatus = const_util.MouldStatusActive
		}

		mouldMasterFields["lastUpdatedBy"] = userId
		mouldMasterFields["canCreateWorkOrder"] = true
		mouldMasterFields["canSubmitTestRequest"] = false
		mouldMasterFields["canCustomerApprove"] = false
		mouldMasterFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		serialisedData, _ := json.Marshal(mouldMasterFields)
		updatingData["object_info"] = serialisedData

		err = database.Update(v.Database, const_util.MouldMasterTable, recordId, updatingData)

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Status",
				Description: "Make sure the mould is get fully qualified before moving into customer confirmation",
			})
		return
	}

	// now make any mould test request related to this mould mark it active
	var condition = " object_info->>'$.mouldId' = " + strconv.Itoa(recordId)
	listOfMouldTestRequests, err := database.GetConditionalObjects(v.Database, const_util.MouldMasterTable, condition)
	if err != nil {
		v.Logger.Error("failed to fetch mould test request", zap.Error(err))
	} else {
		for _, mouldTestRequestInterface := range *listOfMouldTestRequests {
			mouldTestRequest := database.MouldTestRequest{ObjectInfo: mouldTestRequestInterface.ObjectInfo}
			mouldTestRequestInfo := mouldTestRequest.GetMouldTestRequestInfo()
			mouldTestRequestInfo.CanApprove = true
			var serialisedData = mouldTestRequestInfo.DatabaseSerialize(userId)
			v.Logger.Info("updating mould test request status canApprove true", zap.Any("test_id", mouldTestRequestInterface.Id))
			err = database.Update(v.Database, const_util.MouldTestRequestTable, recordId, serialisedData)
			if err != nil {
				v.Logger.Error("failed to update mould test request status canApprove true", zap.Error(err))
			}
		}
	}
	if customerConfirmRequest.IsExport == true {
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Message: "Your mould has now become an export",
			Code:    0,
		})
	} else {
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Message: "Your mould has now become an active, it is now ready for production order",
			Code:    0,
		})
	}

}
