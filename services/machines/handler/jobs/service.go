package jobs

import (
	"cx-micro-flake/services/machines/handler/const_util"
	"cx-micro-flake/services/machines/handler/database"
	"cx-micro-flake/services/machines/handler/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobsConfig struct {
	AssemblyHelpSignalGenerationPeriod string `json:"assemblyHelpSignalGenerationPeriod"`
}
type JobService struct {
	Logger                    *zap.Logger
	Database                  *gorm.DB
	LastProcessedHelpSignalTs int
	JobConfig                 JobsConfig
}

func (v *JobService) Init() error {
	var numberOfRecords = database.Count(v.Database, const_util.AssemblyHelpSignalProcessedTimeTable)
	if numberOfRecords == 0 {
		v.Logger.Info("first time starting up, no signal processed time found, leaving empty")
		v.LastProcessedHelpSignalTs = int(time.Now().Unix())
	} else {
		err, c := database.Get(v.Database, const_util.AssemblyHelpSignalProcessedTimeTable, 1)
		if err != nil {
			v.Logger.Info("error executing query", zap.Error(err))
			v.LastProcessedHelpSignalTs = int(time.Now().Unix())
		}
		var assemblyHelpSignalProcessedTimeInfo = model.GetAssemblyHelpSignalProcessedTimeInfo(c.ObjectInfo)
		if assemblyHelpSignalProcessedTimeInfo == nil {
			v.Logger.Error("error getting assembly signal processed time")
			v.LastProcessedHelpSignalTs = int(time.Now().Unix())
		} else {
			v.LastProcessedHelpSignalTs = assemblyHelpSignalProcessedTimeInfo.ProcessedMessageTs
		}

		v.Logger.Info("found the time last processed, setting the message processed timestamp", zap.Int("message_time", v.LastProcessedHelpSignalTs))
	}

	return nil
}
