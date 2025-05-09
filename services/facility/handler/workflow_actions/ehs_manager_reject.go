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

func (v *ActionService) HandleEHSManagerReject(ctx *gin.Context) {
	updateRequest := make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recordId, dbConnection, basicUserInfo, serviceRequestInfo := v.getBasicInfo(ctx)
	if util.InterfaceToBool(serviceRequestInfo["canEHSManagerReject"]) {
		// dont' allow if ack email not configured or route email configured
		if !v.EmailHandler.IsEmailTemplateExist(dbConnection, const_util.EHSManagerRejectEmailTemplateType) {
			v.Logger.Error("handle EHS manager reject has failed due to invalid email template type", zap.Any("record_id", recordId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring user acknowledgement and routing email templates to proceed further",
				})
			return
		}
		categoryId := util.InterfaceToInt(serviceRequestInfo["categoryId"])
		err, categoryObject := database.Get(dbConnection, const_util.FacilityServiceRequestCategoryTable, categoryId)

		if err != nil {
			v.Logger.Error("handle EHS manager reject has failed due to getting category has failed", zap.Any("category_id", categoryId))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Actions Allowed",
					Description: "Please consider configuring category in the it service request",
				})
			return
		}

		categoryInfo := make(map[string]interface{})
		json.Unmarshal(categoryObject.ObjectInfo, &categoryInfo)
		categoryTemplateIdId := util.InterfaceToInt(categoryInfo["categoryTemplate"])

		err, categoryTemplateObject := database.Get(dbConnection, const_util.FacilityServiceCategoryTemplateTable, categoryTemplateIdId)
		categoryTemplateInfo := make(map[string]interface{})
		json.Unmarshal(categoryTemplateObject.ObjectInfo, &categoryTemplateInfo)

		v.EmailHandler.EmailGenerator(dbConnection, const_util.EHSManagerRejectEmailTemplateType, util.InterfaceToInt(serviceRequestInfo["createdBy"]), const_util.FacilityServiceMyRequestComponent, recordId)

		serviceRequestInfo["canEHSManagerReject"] = false
		serviceRequestInfo["canEHSManagerApprove"] = false
		serviceRequestInfo["canExecPartyComplete"] = false
		serviceRequestInfo["isAssignable"] = false
		serviceRequestInfo["canEdit"] = false
		serviceRequestInfo["canCancel"] = false
		serviceRequestInfo["serviceStatus"] = const_util.WorkFlowUser
		serviceRequestInfo["actionStatus"] = const_util.ActionRejected
		// send the email about ack saying, you request is under review

		remark := util.InterfaceToString(updateRequest["remark"])

		existingActionRemarks := serviceRequestInfo["actionRemarks"].([]interface{})
		existingActionRemarks = append(existingActionRemarks, database.ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(const_util.ISOTimeLayout),
			Status:        "REJECTED BY EHS MANAGER",
			UserId:        basicUserInfo.UserId,
			Remarks:       remark,
			ProcessedTime: getTimeDifference(util.InterfaceToString(serviceRequestInfo["createdAt"])),
		})
		serviceRequestInfo["actionRemarks"] = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(serviceRequestInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err = database.Update(dbConnection, const_util.FacilityServiceRequestTable, recordId, updateObject)
		if err != nil {
			v.Logger.Error("handle IT manager rejection has failed due to update resource failed", zap.String("error", err.Error()))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
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
