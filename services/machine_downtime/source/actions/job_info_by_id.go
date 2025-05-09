package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/orm"
	"go.uber.org/zap"
	"net/http"
)

type JobDetail struct {
	Id              int    `json:"id"`
	JobReferenceId  string `json:"jobReferenceId"`
	CheckInTime     string `json:"checkInTime"`
	CheckInDate     string `json:"checkInDate"`
	CheckOutDate    string `json:"checkOutDate"`
	CheckOutTime    string `json:"checkOutTime"`
	CheckInUserId   string `json:"checkInUserId"`
	CheckOutUserId  string `json:"checkOutUserId"`
	MachineImage    string `json:"machineImage"`
	MachineName     string `json:"machineName"`
	CreatedAt       string `json:"createdAt"`
	LastUpdatedAt   string `json:"lastUpdatedAt"`
	CreatedBy       int    `json:"createdBy"`
	LastUpdatedBy   int    `json:"lastUpdatedBy"`
	CanCheckOut     bool   `json:"canCheckOut"`
	FaultType       int    `json:"faultType"`
	FaultCode       int    `json:"faultCode"`
	Status          string `json:"status"`
	StatusColorCode string `json:"statusColorCode"`
	Remarks         string `json:"remarks"`
}

func (v *Actions) HandleJobInfoById(ctx *gin.Context) {
	recordId := util.GetRecordId(ctx)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	err, c := orm.Get(v.Database, consts.MachineDownTimeMasterTable, recordId)
	if err != nil {
		v.Logger.Error("error getting machine downtime data", zap.Error(err))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Job ID",
				Description: "Internal system error getting corresponding job information",
			})
	}
	var downtimeStatusMap = make(map[int]*models.MachineDowntimeStatusInfo)
	err, i := orm.GetObjects(v.Database, consts.MachineDownTimeStatusTable)
	if err == nil {
		for _, record := range i {
			downtimeStatusMap[record.Id] = models.GetMachineDowntimeStatusInfo(record.ObjectInfo)
		}
	} else {
		v.Logger.Error("error retrieving machine downtime records", zap.String("error", err.Error()))
	}

	downtimeInfo := models.GetMachineDowntimeInfo(c.ObjectInfo)
	err, machineObject := machineService.GetAssemblyMachineInfoById(downtimeInfo.MachineId)
	assemblyMachineMasterInfo, _ := dto.GetAssemblyMachineMasterInfo(machineObject.ObjectInfo)
	if err == nil {
		//get the first job now
		jobDetail := JobDetail{}
		jobDetail.Id = c.Id
		jobDetail.MachineImage = assemblyMachineMasterInfo.MachineImage
		jobDetail.MachineName = assemblyMachineMasterInfo.Description
		jobDetail.FaultCode = downtimeInfo.FaultCode
		jobDetail.JobReferenceId = downtimeInfo.JobReferenceId
		jobDetail.FaultType = downtimeInfo.FaultType
		jobDetail.CanCheckOut = downtimeInfo.CanCheckOut
		jobDetail.CheckInDate = downtimeInfo.CheckInDate
		jobDetail.CheckOutDate = downtimeInfo.CheckOutDate
		jobDetail.LastUpdatedBy = downtimeInfo.LastUpdatedBy
		jobDetail.LastUpdatedAt = downtimeInfo.LastUpdatedAt
		jobDetail.CreatedAt = downtimeInfo.CreatedAt
		var downtimeStatus = ""
		var defaultStatusColorCode = "#bfd3df "
		if statusValue, ok := downtimeStatusMap[downtimeInfo.Status]; ok {
			downtimeStatus = statusValue.Status
			defaultStatusColorCode = statusValue.ColorCode
		}
		jobDetail.Status = downtimeStatus
		jobDetail.StatusColorCode = defaultStatusColorCode
		jobDetail.Remarks = downtimeInfo.Remarks
		jobDetail.CheckInUserId = authService.GetUserInfoById(downtimeInfo.CheckInUserId).FullName
		jobDetail.CheckOutUserId = authService.GetUserInfoById(downtimeInfo.CheckOutUserId).FullName
		ctx.JSON(http.StatusOK, jobDetail)
		v.Logger.Info("sending job detail response", zap.Any("response", jobDetail))

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Machine",
				Description: "Internal system error getting corresponding job information",
			})
	}

}
