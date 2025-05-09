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

func (v *Actions) HandleCheckIn(ctx *gin.Context) {
	v.Logger.Info("handle check in received")
	checkInRequest := dto.CheckInRequest{}
	if err := ctx.ShouldBindBodyWith(&checkInRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		v.Logger.Error("error sending response", zap.Error(err))
		return
	}

	for _, faultId := range checkInRequest.FaultIds {
		v.Logger.Info("checkin the machine downtime job", zap.Int("faultId", faultId))
		userId := common.GetUserId(ctx)
		err, c := orm.Get(v.Database, consts.MachineDownTimeMasterTable, faultId)
		if err == nil {
			machineDowntimeInfo := models.GetMachineDowntimeInfo(c.ObjectInfo)
			machineDowntimeInfo.CheckInDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			machineDowntimeInfo.CheckInTime = util.GetCurrentTimeInSingapore("15:04:05")
			machineDowntimeInfo.CanCheckOut = true
			machineDowntimeInfo.CanCheckIn = false
			machineDowntimeInfo.Status = consts.DowntimeStatus_Fault_Under_Investigation
			machineDowntimeInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			machineDowntimeInfo.LastUpdatedBy = userId
			machineDowntimeInfo.CheckInUserId = userId
			machineDowntimeInfo.AssignedUserId = userId
			serialisedData, _ := machineDowntimeInfo.Serialised()
			err := orm.UpdateSerialisedResourceFromId(v.Database, consts.MachineDownTimeMasterTable, faultId, userId, serialisedData)
			if err != nil {
				v.Logger.Error("error check in error", zap.Error(err))
				response.SendInternalSystemError(ctx)
				return
			} else {
				v.Logger.Info("updated the checkin job successfully", zap.Int("faultId", faultId))
			}
		} else {
			v.Logger.Error("error getting the fault id", zap.Error(err))
			response.SendInternalSystemError(ctx)
			return
		}
	}

	v.Logger.Info("machine downtime is successfully updated", zap.Any("fault_ids", checkInRequest.FaultIds))
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully CheckedIn",
	})
}
