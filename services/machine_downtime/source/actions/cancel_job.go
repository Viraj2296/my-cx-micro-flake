package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.cerex.io/transcendflow/component"
	"go.cerex.io/transcendflow/header_parser"
	"go.cerex.io/transcendflow/orm"
	"go.cerex.io/transcendflow/util"
	"go.uber.org/zap"
)

func (v *Actions) HandleCancelJob(ctx *gin.Context) {
	v.Logger.Info("handle cancel job")
	recordId := header_parser.GetRecordId(ctx)
	var userId = common.GetUserId(ctx)

	err, c := orm.Get(v.Database, consts.MachineDownTimeMasterTable, recordId)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Job",
				Description: "The job is not available in the system, job might have archived or removed from system, please contact system administrator",
			})
	}
	cancelJob := dto.CancelJobRequest{}
	if err := ctx.ShouldBindBodyWith(&cancelJob, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	machineDowntimeInfo := models.GetMachineDowntimeInfo(c.ObjectInfo)
	machineDowntimeInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	machineDowntimeInfo.LastUpdatedBy = userId
	machineDowntimeInfo.CanCheckOut = false
	machineDowntimeInfo.CanCheckIn = false
	machineDowntimeInfo.CanCancel = false
	machineDowntimeInfo.Remarks = cancelJob.Remarks
	machineDowntimeInfo.Status = consts.DowntimeStatus_Fault_Cacelled

	machineDowntimeInfo.ActionRemarks = append(machineDowntimeInfo.ActionRemarks, component.ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(util.ISOTimeLayout),
		Status:        "Fault is forcefully cancelled",
		UserId:        1,
		Remarks:       "The fault is cancelled now",
		ProcessedTime: util.GetTimeDifference(util.InterfaceToString(machineDowntimeInfo.CreatedAt)),
	})

	// var updatingObject = make(map[string]interface{})
	if serialisedData, err := machineDowntimeInfo.Serialised(); err == nil {
		// updatingObject["object_info"] = serialisedData
		err := orm.UpdateSerialisedResourceFromId(v.Database, consts.MachineDownTimeMasterTable, recordId, userId, serialisedData)
		if err != nil {
			v.Logger.Error("error updating downtime master", zap.Error(err))
		} else {
			v.Logger.Info("successfully check-out")
			ctx.JSON(http.StatusOK, response.GeneralResponse{
				Code:    0,
				Message: "Successfully Checkout",
			})

			return
		}
	} else {
		v.Logger.Error("error getting machine downtime", zap.Error(err))
	}

	response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
		&response.DetailedError{
			Header:      "Invalid operation",
			Description: "This operation is not able complete due to internal system error",
		})

}
