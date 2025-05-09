package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.cerex.io/transcendflow/orm"
	"go.uber.org/zap"
)

type CheckOutInfoResponse struct {
	Id                int    `json:"id"`
	JobReferenceId    string `json:"jobReferenceId"`
	MachineName       string `json:"machineName"`
	MachineImage      string `json:"machineImage"`
	CreatedDateTime   string `json:"createdDateTime"`
	CheckedInDateTime string `json:"checkedInDateTime"`
}

func (v *Actions) HandleCheckOutFaults(ctx *gin.Context) {
	// get the list of fault type
	jobInfoRequest := dto.JobInfoRequest{}
	if err := ctx.ShouldBindBodyWith(&jobInfoRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	err, machineId := machineService.GetAssemblyMachineInfoFromEquipmentId(consts.ProjectID, jobInfoRequest.EquipmentId)
	var downtimeStatusMap = make(map[int]*models.MachineDowntimeStatusInfo)
	err, i := orm.GetObjects(v.Database, consts.MachineDownTimeStatusTable)
	if err == nil {
		for _, record := range i {
			downtimeStatusMap[record.Id] = models.GetMachineDowntimeStatusInfo(record.ObjectInfo)
		}
	} else {
		v.Logger.Error("error retrieving machine downtime records", zap.String("error", err.Error()))
	}

	// send only the job that is not checked out for that machine
	var condition = " object_info->>'$.machineId'= " + strconv.Itoa(machineId) + " AND object_info->>'$.canCheckOut'= 'true'"
	err, listOfJobs := orm.GetConditionalObjects(v.Database, consts.MachineDownTimeMasterTable, condition)
	var checkOutFaults = make([]CheckOutInfoResponse, 0)

	if err == nil {
		//get the first job now
		if len(listOfJobs) > 0 {
			for _, jobInterface := range listOfJobs {
				var downtimeInfo = models.GetMachineDowntimeInfo(jobInterface.ObjectInfo)
				checkOutInfoResponse := CheckOutInfoResponse{}
				checkOutInfoResponse.Id = jobInterface.Id
				err, c := machineService.GetAssemblyMachineInfoById(machineId)

				if err == nil {
					assemblyMachineMasterInfo, _ := dto.GetAssemblyMachineMasterInfo(c.ObjectInfo)
					checkOutInfoResponse.MachineName = assemblyMachineMasterInfo.Description
					checkOutInfoResponse.MachineImage = assemblyMachineMasterInfo.MachineImage
				}
				if downtimeInfo != nil {
					checkOutInfoResponse.CreatedDateTime = downtimeInfo.CreatedAt
					checkOutInfoResponse.JobReferenceId = downtimeInfo.JobReferenceId
					checkOutInfoResponse.CheckedInDateTime = downtimeInfo.CheckInDate + downtimeInfo.CheckInTime
				}
				checkOutFaults = append(checkOutFaults, checkOutInfoResponse)
			}

			ctx.JSON(http.StatusOK, checkOutFaults)
			v.Logger.Info("sending downtime response", zap.Any("response", checkOutFaults))
		} else {
			v.Logger.Error("no jobs found, query", zap.Any("query", condition))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "No Jobs Found",
					Description: "Sorry, no jobs found to checkout !!",
				})
		}

	} else {
		v.Logger.Error("error getting machine downtime master", zap.Error(err))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Machine",
				Description: "Internal system error getting corresponding job information",
			})
	}

}
