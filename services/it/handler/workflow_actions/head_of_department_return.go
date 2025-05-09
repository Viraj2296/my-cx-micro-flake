package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

func (v *ActionService) HandleDepartmentHeadReturn(ctx *gin.Context) {
	v.Logger.Info("handle head of department return request is received")
	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		v.Logger.Error("error getting HOD return request payload", zap.String("error", err.Error()))
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	returnRemark := util.InterfaceToString(returnFields["remark"])
	recordId, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canHODReturn"]) {
		if !v.EmailHandler.IsEmailTemplateExist(v.Database, const_util.HeadOfDepartmentReturnEmailTemplateType) {
			v.Logger.Error("error getting HOD return email template type", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}

		if existingLevelCounter, ok := serviceRequestInfo["levelCounter"]; ok {
			var intLevelCounter = util.InterfaceToInt(existingLevelCounter)
			serviceRequestInfo["levelCounter"] = intLevelCounter - 1
		}
		serviceRequestInfo["canEdit"] = true
		serviceRequestInfo["canHODReturn"] = false
		serviceRequestInfo["canUserSubmit"] = true
		serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
		serviceRequestInfo["actionStatus"] = const_util.ActionReturnedByHoD
		// send the email about ack saying, you request is under review
		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "RETURNED BY HOD",
			UserId:        basicUserInfo.UserId,
			Remarks:       returnRemark,
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := database.Update(v.Database, const_util.ITServiceRequestTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}

		v.EmailHandler.EmailGenerator(v.Database, const_util.HeadOfDepartmentReturnEmailTemplateType, util.InterfaceToInt(serviceRequestInfo["createdBy"]), const_util.ITServiceMyRequestComponent, recordId)
		v.Logger.Info("handle HOD return request is success", zap.Any("record_id", recordId))
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your request has been processed successfully",
		})
	} else {
		v.Logger.Error("handle HOD return request has failed due to flag canHODReturn is not true", zap.Any("record_id", recordId))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Your action can not be performed against this request due to sequence validation",
			})
	}

}
