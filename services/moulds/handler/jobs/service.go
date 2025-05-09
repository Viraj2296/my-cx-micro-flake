package jobs

import (
	"cx-micro-flake/services/moulds/handler/notification"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobService struct {
	Database        *gorm.DB
	Logger          *zap.Logger
	PoolingInterval int
	EmailHandler    *notification.EmailHandler
}

func (v *JobService) Init() {

}
