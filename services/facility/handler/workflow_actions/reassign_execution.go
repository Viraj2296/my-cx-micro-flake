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

func (v *ActionService) HandleReassignExecution(ctx *gin.Context) {
	v.Logger.Info("handling re-assign execution")
	var reassignPayload = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&reassignPayload, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reassignUserId = util.InterfaceToInt(reassignPayload["reassignExecutionUser"])
	v.Logger.Info("reassign user id", zap.Any("user_id", reassignUserId))

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["isAssignable"]) {
		// dont' allow if ack email not configured or route email configured

		serviceRequestInfo["isAssignable"] = true
		serviceRequestInfo["assignedExecutionParty"] = reassignUserId
		serviceRequestInfo["canExecPartyComplete"] = true
		// send the email about ack saying, you request is under review
		serviceRequestInfo["canAddCustomExecution"] = true
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "ASSIGNED EXTERNALLY",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Execution party assigned from external",
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle re-assign has failed due to update resource", zap.String("error", err.Error()))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		err = v.EmailHandler.EmailGenerator(dbConnection, const_util.AssignExecutionerEmailTemplateType, reassignUserId, const_util.FacilityServiceMyExecutionRequestComponent, recordId)
		if err != nil {
			v.Logger.Error("error sending email to re-assigned users", zap.Any("error", err.Error()))
		}

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
