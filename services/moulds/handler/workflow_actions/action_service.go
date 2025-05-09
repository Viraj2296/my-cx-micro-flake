package workflow_actions

import (
	"cx-micro-flake/services/moulds/handler/notification"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ActionService struct {
	Logger       *zap.Logger
	Database     *gorm.DB
	EmailHandler *notification.EmailHandler
}

func (v *ActionService) Init(logger *zap.Logger) {
	v.Logger = logger
}
