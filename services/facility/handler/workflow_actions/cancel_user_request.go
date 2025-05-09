package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

func (v *ActionService) HandleCancelRequest(ctx *gin.Context) {
	v.Logger.Info("handling cancel request")
	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		v.Logger.Error("invalid payload, return error now", zap.String("error", err.Error()))
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canCancel"]) {
		returnRemark := util.InterfaceToString(returnFields["remark"])
		serviceRequestInfo["canEdit"] = false
		serviceRequestInfo["canUserSubmit"] = false
		serviceRequestInfo["canCancel"] = false
		serviceRequestInfo["canHODReturn"] = false
		serviceRequestInfo["canHODReject"] = false
		serviceRequestInfo["canHODApprove"] = false
		serviceRequestInfo["canSapApprove"] = false
		serviceRequestInfo["canITManagerApprove"] = false
		serviceRequestInfo["canUserAcknowledge"] = false
		serviceRequestInfo["canExecPartyComplete"] = false
		serviceRequestInfo["isAssignable"] = false
		serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
		serviceRequestInfo["actionStatus"] = const_util.ActionCanceled

		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "FACILITY REQUEST IS CANCELLED",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Facility request is canceled by " + basicUserInfo.FullName + ". Because '" + returnRemark + "'",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks

		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("exception handling request, updating the resource failed", zap.String("error", err.Error()))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		v.EmailHandler.EmailGenerator(dbConnection, const_util.FacilityCanceledRequestEmailTemplateType, util.InterfaceToInt(serviceRequestInfo["createdBy"]), const_util.FacilityServiceMyRequestComponent, recordId)
		v.Logger.Info("handle cancel request has been successfully processed", zap.Any("record_id", recordId))

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})
	} else {
		v.Logger.Error("invalid action, this action can not be performed due to canCancel flag is not set", zap.Any("record_id", recordId))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
