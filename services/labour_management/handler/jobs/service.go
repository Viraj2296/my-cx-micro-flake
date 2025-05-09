package jobs

import (
	"cx-micro-flake/services/labour_management/handler/cache"
	"io/ioutil"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobService struct {
	Database                *gorm.DB
	Logger                  *zap.Logger
	PoolingInterval         int
	EscalationEmailTemplate string
	MachineStatsCache       *cache.MachineStatsCache
}

func (v *JobService) Init() {

	content, err := ioutil.ReadFile(v.EscalationEmailTemplate)
	if err != nil {
		v.Logger.Error("error reading email template", zap.String("error", err.Error()))
		v.EscalationEmailTemplate = "" // don't send email if it is empty
	} else {
		v.EscalationEmailTemplate = string(content)
	}
}
