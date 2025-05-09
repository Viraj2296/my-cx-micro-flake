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

type CheckInInfoResponse struct {
	Id              int    `json:"id"`
	JobReferenceId  string `json:"jobReferenceId"`
	MachineName     string `json:"machineName"`
	MachineImage    string `json:"machineImage"`
	CreatedDateTime string `json:"createdDateTime"`
}

func (v *Actions) HandleCheckInFaults(ctx *gin.Context) {
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
	var condition = " object_info->>'$.machineId'= " + strconv.Itoa(machineId) + " AND object_info->>'$.canCheckIn'= 'true'"
	err, listOfJobs := orm.GetConditionalObjects(v.Database, consts.MachineDownTimeMasterTable, condition)
	var checkInFaults = make([]CheckInInfoResponse, 0)

	if err == nil {
		//get the first job now
		if len(listOfJobs) > 0 {
			for _, jobInterface := range listOfJobs {
				var downtimeInfo = models.GetMachineDowntimeInfo(jobInterface.ObjectInfo)
				checkInInfoResponse := CheckInInfoResponse{}
				checkInInfoResponse.Id = jobInterface.Id
				err, c := machineService.GetAssemblyMachineInfoById(machineId)
				if err == nil {
					assemblyMachineMasterInfo, _ := dto.GetAssemblyMachineMasterInfo(c.ObjectInfo)
					checkInInfoResponse.MachineName = assemblyMachineMasterInfo.Description
					checkInInfoResponse.MachineImage = assemblyMachineMasterInfo.MachineImage
				}
				if downtimeInfo != nil {
					checkInInfoResponse.CreatedDateTime = downtimeInfo.CreatedAt
					checkInInfoResponse.JobReferenceId = downtimeInfo.JobReferenceId
				}

				checkInFaults = append(checkInFaults, checkInInfoResponse)
			}

			ctx.JSON(http.StatusOK, checkInFaults)
			v.Logger.Info("sending downtime response", zap.Any("response", checkInFaults))
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
