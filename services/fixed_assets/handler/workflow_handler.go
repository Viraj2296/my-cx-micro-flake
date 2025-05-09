package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}

func (v *FixedAssetService) handleSubmit(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	if componentName == FixedAssetMyDisposalComponent {
		err, fixedAssetsDisposal := Get(dbConnection, FixedAssetDisposalTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetDisposal := FixedAssetDisposal{ObjectInfo: fixedAssetsDisposal.ObjectInfo}
		if fixedAssetDisposal.getDisposalRequest().CanSubmit {
			disposalRequest := fixedAssetDisposal.getDisposalRequest()
			disposalRequest.CanSubmit = false
			disposalRequest.CanHODApprove = true
			disposalRequest.CanHODReject = true
			disposalRequest.DisposalStatus = DisposalStatusSubmitted
			disposalRequest.ActionStatus = "PENDING HOD REVIEW"
			disposalRequest.WorkflowLevel = DisposalHODWorkflowLevel
			disposalRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

			err := Update(dbConnection, FixedAssetDisposalTable, recordId, disposalRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	} else if componentName == FixedAssetTransferComponent {
		err, fixedAssetsTransfer := Get(dbConnection, FixedAssetTransferTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetTransfer := FixedAssetTransfer{ObjectInfo: fixedAssetsTransfer.ObjectInfo}
		if fixedAssetTransfer.getTransferRequest().CanSubmit {
			fixedAssetTransferRequest := fixedAssetTransfer.getTransferRequest()
			fixedAssetTransferRequest.CanSubmit = false
			fixedAssetTransferRequest.CanHODApprove = true
			fixedAssetTransferRequest.CanHODReject = true
			fixedAssetTransferRequest.TransferStatus = TransferStatusSubmitted
			fixedAssetTransferRequest.ActionStatus = "PENDING HOD REVIEW"
			fixedAssetTransferRequest.WorkflowLevel = DisposalHODWorkflowLevel

			fixedAssetTransferRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

			err := Update(dbConnection, FixedAssetTransferTable, recordId, fixedAssetTransferRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	}

}

func (v *FixedAssetService) handleHODReject(ctx *gin.Context) {

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	returnRemark := util.InterfaceToString(returnFields["remark"])

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	if componentName == FixedAssetDisposalComponent {
		err, fixedAssetsDisposal := Get(dbConnection, FixedAssetDisposalTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetDisposal := FixedAssetDisposal{ObjectInfo: fixedAssetsDisposal.ObjectInfo}
		if fixedAssetDisposal.getDisposalRequest().CanHODReject {
			disposalRequest := fixedAssetDisposal.getDisposalRequest()
			disposalRequest.CanSubmit = true
			disposalRequest.CanHODApprove = false
			disposalRequest.CanHODReject = false
			disposalRequest.WorkflowLevel = DisposalUserWorkflowLevel
			disposalRequest.DisposalStatus = DisposalStatusHODRejected
			disposalRequest.ActionStatus = "USER ACTION NEEDED"
			existingActionRemarks := disposalRequest.ActionRemarks
			existingActionRemarks = append(existingActionRemarks, ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
				Status:        "REJECTED BY HOD",
				UserId:        userId,
				Remarks:       returnRemark,
				ProcessedTime: getTimeDifference(disposalRequest.CreatedAt),
			})
			disposalRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			disposalRequest.WorkflowLevel = DisposalHODWorkflowLevel
			err := Update(dbConnection, FixedAssetDisposalTable, recordId, disposalRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	} else if componentName == FixedAssetTransferComponent {
		err, fixedAssetTransferObject := Get(dbConnection, FixedAssetTransferTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetTransfer := FixedAssetTransfer{ObjectInfo: fixedAssetTransferObject.ObjectInfo}
		if fixedAssetTransfer.getTransferRequest().CanHODReject {
			transferRequest := fixedAssetTransfer.getTransferRequest()
			transferRequest.CanSubmit = true
			transferRequest.CanHODApprove = false
			transferRequest.CanHODReject = false
			transferRequest.WorkflowLevel = TransferUserWorkflowLevel
			transferRequest.TransferStatus = TransferStatusHODRejected
			transferRequest.ActionStatus = "USER ACTION NEEDED"

			existingActionRemarks := transferRequest.ActionRemarks
			existingActionRemarks = append(existingActionRemarks, ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
				Status:        "REJECTED BY HOD",
				UserId:        userId,
				Remarks:       returnRemark,
				ProcessedTime: getTimeDifference(transferRequest.CreatedAt),
			})
			transferRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			transferRequest.WorkflowLevel = TransferHODWorkflowLevel
			err := Update(dbConnection, FixedAssetTransferTable, recordId, transferRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	}

}

func (v *FixedAssetService) handleHODApprove(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	if componentName == FixedAssetDisposalComponent {
		err, fixedAssetsDisposal := Get(dbConnection, FixedAssetDisposalTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetDisposal := FixedAssetDisposal{ObjectInfo: fixedAssetsDisposal.ObjectInfo}
		if fixedAssetDisposal.getDisposalRequest().CanHODApprove {
			disposalRequest := fixedAssetDisposal.getDisposalRequest()
			disposalRequest.CanSubmit = false
			disposalRequest.CanHODApprove = false
			disposalRequest.CanHODReject = false
			disposalRequest.CanCEOApprove = true
			disposalRequest.CanCEOReject = true
			disposalRequest.DisposalStatus = DisposalStatusCEOApproved
			disposalRequest.ActionStatus = "PENDING CEO REVIEW"
			disposalRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			disposalRequest.WorkflowLevel = DisposalCEOWorkflowLevel
			err := Update(dbConnection, FixedAssetDisposalTable, recordId, disposalRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	} else if componentName == FixedAssetTransferComponent {
		err, fixedAssetsTransferObject := Get(dbConnection, FixedAssetTransferTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetTransfer := FixedAssetTransfer{ObjectInfo: fixedAssetsTransferObject.ObjectInfo}
		if fixedAssetTransfer.getTransferRequest().CanHODApprove {
			transferRequest := fixedAssetTransfer.getTransferRequest()
			transferRequest.CanSubmit = false
			transferRequest.CanHODApprove = false
			transferRequest.CanHODReject = false
			transferRequest.CanCEOApprove = true
			transferRequest.CanCEOReject = true
			transferRequest.TransferStatus = TransferStatusCEOApproved
			transferRequest.ActionStatus = "PENDING CEO REVIEW"
			transferRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			transferRequest.WorkflowLevel = TransferCEOWorkflowLevel

			transferRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			transferRequest.WorkflowLevel = TransferCEOWorkflowLevel
			err := Update(dbConnection, FixedAssetDisposalTable, recordId, transferRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	}

}

func (v *FixedAssetService) handleCEOApprove(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	if componentName == FixedAssetDisposalComponent {
		err, fixedAssetsDisposal := Get(dbConnection, FixedAssetDisposalTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetDisposal := FixedAssetDisposal{ObjectInfo: fixedAssetsDisposal.ObjectInfo}
		if fixedAssetDisposal.getDisposalRequest().CanCEOApprove {
			disposalRequest := fixedAssetDisposal.getDisposalRequest()
			disposalRequest.CanSubmit = false
			disposalRequest.CanHODApprove = false
			disposalRequest.CanHODReject = false
			disposalRequest.CanCEOApprove = false
			disposalRequest.CanCEOReject = false
			disposalRequest.CanUserAcknowledge = true
			disposalRequest.ActionStatus = "COMPLETED"
			disposalRequest.DisposalStatus = DisposalStatusCEOApproved
			disposalRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			disposalRequest.WorkflowLevel = DisposalUserWorkflowLevel
			err := Update(dbConnection, FixedAssetDisposalTable, recordId, disposalRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	} else if componentName == FixedAssetTransferComponent {
		err, fixedAssetsTransferObject := Get(dbConnection, FixedAssetTransferTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetTransfer := FixedAssetTransfer{ObjectInfo: fixedAssetsTransferObject.ObjectInfo}
		if fixedAssetTransfer.getTransferRequest().CanCEOApprove {
			transferRequest := fixedAssetTransfer.getTransferRequest()
			transferRequest.CanSubmit = false
			transferRequest.CanHODApprove = false
			transferRequest.CanHODReject = false
			transferRequest.CanCEOApprove = false
			transferRequest.CanCEOReject = false
			transferRequest.ActionStatus = "COMPLETED"
			transferRequest.TransferStatus = TransferStatusCEOApproved
			transferRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			transferRequest.WorkflowLevel = TransferUserWorkflowLevel
			err := Update(dbConnection, FixedAssetTransferTable, recordId, transferRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	}

}

func (v *FixedAssetService) handleCEOReject(ctx *gin.Context) {

	var returnFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&returnFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	returnRemark := util.InterfaceToString(returnFields["remark"])

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	if componentName == FixedAssetDisposalComponent {
		err, fixedAssetsDisposal := Get(dbConnection, FixedAssetDisposalTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetDisposal := FixedAssetDisposal{ObjectInfo: fixedAssetsDisposal.ObjectInfo}
		if fixedAssetDisposal.getDisposalRequest().CanHODReject {
			disposalRequest := fixedAssetDisposal.getDisposalRequest()
			disposalRequest.CanSubmit = true
			disposalRequest.CanHODApprove = false
			disposalRequest.CanHODReject = false
			disposalRequest.CanCEOApprove = false
			disposalRequest.CanCEOReject = false
			disposalRequest.DisposalStatus = DisposalStatusCEORejected
			disposalRequest.ActionStatus = "USER ACTION NEEDED"
			existingActionRemarks := disposalRequest.ActionRemarks
			existingActionRemarks = append(existingActionRemarks, ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
				Status:        "REJECTED BY CEO",
				UserId:        userId,
				Remarks:       returnRemark,
				ProcessedTime: getTimeDifference(disposalRequest.CreatedAt),
			})
			disposalRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			disposalRequest.WorkflowLevel = DisposalUserWorkflowLevel
			err := Update(dbConnection, FixedAssetDisposalTable, recordId, disposalRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	} else if componentName == FixedAssetTransferComponent {
		err, fixedTransferObject := Get(dbConnection, FixedAssetTransferTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Resource not found",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		userId := common.GetUserId(ctx)
		fixedAssetTransfer := FixedAssetTransfer{ObjectInfo: fixedTransferObject.ObjectInfo}
		if fixedAssetTransfer.getTransferRequest().CanHODReject {
			transferRequest := fixedAssetTransfer.getTransferRequest()
			transferRequest.CanSubmit = true
			transferRequest.CanHODApprove = false
			transferRequest.CanHODReject = false
			transferRequest.CanCEOApprove = false
			transferRequest.CanCEOReject = false
			transferRequest.TransferStatus = TransferStatusCEOApproved
			transferRequest.ActionStatus = "USER ACTION NEEDED"
			existingActionRemarks := transferRequest.ActionRemarks
			existingActionRemarks = append(existingActionRemarks, ActionRemarks{
				ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
				Status:        "REJECTED BY CEO",
				UserId:        userId,
				Remarks:       returnRemark,
				ProcessedTime: getTimeDifference(transferRequest.CreatedAt),
			})
			transferRequest.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			transferRequest.WorkflowLevel = TransferUserWorkflowLevel
			err := Update(dbConnection, FixedAssetDisposalTable, recordId, transferRequest.DatabaseSerialize(userId))
			if err != nil {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      "Server Exception",
						Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
					})
				return
			}

		} else {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Action",
					Description: "Your action can not be performed against this request due to sequence validation",
				})
		}
	}

}

func (v *FixedAssetService) handleNotifyVendor(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	_, generalDisposalObject := Get(dbConnection, targetTable, recordId)
	fixedAssetDisposalInfo := make(map[string]interface{})
	json.Unmarshal(generalDisposalObject.ObjectInfo, &fixedAssetDisposalInfo)

	assetId := util.InterfaceToInt(fixedAssetDisposalInfo["assetId"])

	_, generalObject := Get(dbConnection, FixedAssetMasterTable, assetId)
	fixedAssetMaster := FixedAssetMaster{ObjectInfo: generalObject.ObjectInfo}
	fixedAssetInfo := fixedAssetMaster.getFixedAssetMasterInfo()

	vendorEmail := fixedAssetInfo.VendorEmail

	if vendorEmail == "" {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "Vendor Email is not available",
			})
	}

	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	userRecordId := authService.EmailToUserId(vendorEmail)

	// Ssri

	v.emailGenerator(dbConnection, FixedAssetDisposalEmailTemplate, userRecordId, FixedAssetDisposalComponent, recordId)

}

func (v *FixedAssetService) handleUserAcknowledgement(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	userId := common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	_, generalDisposalObject := Get(dbConnection, targetTable, recordId)

	fixedAssetDisposal := FixedAssetDisposal{ObjectInfo: generalDisposalObject.ObjectInfo}
	fixedAssetDisposalInfo := fixedAssetDisposal.getDisposalRequest()

	if fixedAssetDisposalInfo.CanUserAcknowledge {
		fixedAssetDisposalInfo.CanSubmit = false
		fixedAssetDisposalInfo.CanUserAcknowledge = false
		fixedAssetDisposalInfo.CanHODReject = false
		fixedAssetDisposalInfo.CanHODApprove = false
		fixedAssetDisposalInfo.WorkflowLevel = DisposalUserWorkflowLevel
		fixedAssetDisposalInfo.ActionStatus = "Closed"
		// send the email about ack saying, you request is under review
		existingActionRemarks := fixedAssetDisposalInfo.ActionRemarks
		existingActionRemarks = append(existingActionRemarks, ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "ACCEPTED BY USER",
			UserId:        basicUserInfo.UserId,
			Remarks:       "Thank you for accepting your request deliver",
			ProcessedTime: getTimeDifference(fixedAssetDisposalInfo.CreatedAt),
		})
		fixedAssetDisposalInfo.ActionRemarks = existingActionRemarks
		serialisedRequestFields, _ := json.Marshal(fixedAssetDisposalInfo)
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = serialisedRequestFields
		err := Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		//listOfWorkflowUsers := v.getTargetWorkflowGroup(dbConnection, recordId, WorkFlowExecutionParty)
		//v.BaseService.Logger.Infow("execution users", "users", listOfWorkflowUsers)
		//for _, workflowUser := range listOfWorkflowUsers {
		//	v.emailGenerator(dbConnection, UserAcknowledgeEmailTemplateType, workflowUser, ITServiceMyExecutionRequestComponent, recordId)
		//}

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
