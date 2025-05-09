package workflow_actions

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ActionService struct {
	Logger   *zap.Logger
	Database *gorm.DB
}

func (v *ActionService) Init(logger *zap.Logger) {
	v.Logger = logger
}
