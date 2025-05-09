package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.cerex.io/transcendflow/orm"
	"go.uber.org/zap"
)

func (v *Actions) HandleCheckout(ctx *gin.Context) {
	v.Logger.Info("handle check out received")
	recordId := util.GetRecordId(ctx)
	// get the job
	err, c := orm.Get(v.Database, consts.MachineDownTimeMasterTable, recordId)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Job",
				Description: "The job is not available in the system, job might have archived or removed from system, please contact system administrator",
			})
	}
	checkOutRequest := dto.CheckOutRequest{}
	if err := ctx.ShouldBindBodyWith(&checkOutRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	userId := common.GetUserId(ctx)
	machineDowntimeInfo := models.GetMachineDowntimeInfo(c.ObjectInfo)
	machineDowntimeInfo.CheckOutDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	// machineDowntimeInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	machineDowntimeInfo.CheckOutTime = util.GetCurrentTimeInSingapore("15:04:05")
	machineDowntimeInfo.CanCheckOut = false
	machineDowntimeInfo.CheckOutUserId = userId
	machineDowntimeInfo.LastUpdatedBy = userId
	machineDowntimeInfo.CanCancel = false
	machineDowntimeInfo.Status = consts.DowntimeStatus_Fault_Repaired
	machineDowntimeInfo.FaultType = checkOutRequest.FaultType
	machineDowntimeInfo.FaultCode = checkOutRequest.FaultCode
	machineDowntimeInfo.Remarks = checkOutRequest.Remarks
	if serialisedData, err := machineDowntimeInfo.Serialised(); err == nil {
		err := orm.UpdateSerialisedResourceFromId(v.Database, consts.MachineDownTimeMasterTable, recordId, userId, serialisedData)
		if err != nil {
			v.Logger.Error("error updating downtime master", zap.Error(err))
		} else {
			v.Logger.Info("updated the checkout job successfully", zap.Int("faultId", recordId))

			v.Logger.Info("successfully check-out")
			// make the machine into live mode.
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveAssemblyMachineToActive(consts.ProjectID, machineDowntimeInfo.MachineId)
			if err != nil {
				v.Logger.Error("move assembly machine to live error", zap.Error(err))
			}
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
