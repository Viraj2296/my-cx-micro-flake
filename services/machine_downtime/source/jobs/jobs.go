package jobs

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/models"
	"cx-micro-flake/services/machine_downtime/source/module_config"
	"cx-micro-flake/services/machine_downtime/source/repository"
	"io/ioutil"
	"os"
	"time"

	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/orm"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Jobs struct {
	JobAssigningUsers                  []int
	PushNotificationUsers              []int
	LastProcessedMachineHelpSignalTime int64
	EscalationEmailTemplate            string
	DowntimeConfig                     module_config.DowntimeConfig
	Repository                         *repository.Repository
	Database                           *gorm.DB
	Logger                             *logging.Logger
	machineDowntimeSettingInfo         *models.MachineDowntimeSettingInfo
}

// NewJobs initializes a new Jobs instance
func NewJobs(logger *logging.Logger, db *gorm.DB, JobsConfig module_config.DowntimeConfig, repository *repository.Repository) *Jobs {
	return &Jobs{
		Logger:         logger,
		Database:       db,
		DowntimeConfig: JobsConfig,
		Repository:     repository,
	}
}
func (v *Jobs) Init() {
	err, settingInterface := orm.Get(v.Database, consts.MachineDownTimeSettingTable, 1)
	if err != nil {
		v.Logger.Error("error getting labour management setting", zap.Error(err))
		os.Exit(0)
	}
	v.machineDowntimeSettingInfo = models.GetMachineDowntimeSettingInfo(settingInterface.ObjectInfo)
	var jobAssigningUsers = v.machineDowntimeSettingInfo.JobAssigningUsers
	var pushNotificationJobRoles = v.machineDowntimeSettingInfo.PushNotificationJobRoles
	v.JobAssigningUsers = jobAssigningUsers
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	v.Logger.Info("loaded job assigning users ", zap.Any("job_user", jobAssigningUsers))
	for _, roleId := range pushNotificationJobRoles {
		listOfUsers := authService.GetUsersFromJobRoleId(roleId)
		if len(listOfUsers) > 0 {
			v.PushNotificationUsers = append(v.PushNotificationUsers, listOfUsers...)
		}
	}
	var numberOfRecords = orm.Count(v.Database, consts.MachineDownTimeSignalHistoryTable)
	if numberOfRecords == 0 {
		v.Logger.Info("no signal history record found, creating a new one")
		var currentTimeEpoch = time.Now().Unix()
		v.LastProcessedMachineHelpSignalTime = currentTimeEpoch
	} else {
		var getLastSignalHistoryRecord = " id > 0 order by object_info->>'$.createdAt' desc limit 1"
		err, signalHistoryRecords := orm.GetConditionalObjects(v.Database, consts.MachineDownTimeSignalHistoryTable, getLastSignalHistoryRecord)
		if err != nil {
			v.LastProcessedMachineHelpSignalTime = time.Now().Unix()
			v.Logger.Error("error getting last signal history", zap.Error(err))
		} else {
			if len(signalHistoryRecords) == 1 {
				var signalHistoryInfo = models.GetSignalHistoryInfo(signalHistoryRecords[0].ObjectInfo)
				v.LastProcessedMachineHelpSignalTime = signalHistoryInfo.CreatedAt
				v.Logger.Info("last signal history record found,setting processed help signal time", logging.Any("LastProcessedMachineHelpSignalTime", signalHistoryInfo.CreatedAt))
			}
		}
	}

	content, err := ioutil.ReadFile(v.DowntimeConfig.JobServiceConfig.EscalationEmailTemplate)
	if err != nil {
		v.Logger.Error("error reading email template", zap.String("error", err.Error()))
		v.EscalationEmailTemplate = "" // don't send email if it is empty
	} else {
		v.EscalationEmailTemplate = string(content)
	}
	v.Logger.Info("last read signal history time is set to ", zap.Int64("ts", v.LastProcessedMachineHelpSignalTime))
}
