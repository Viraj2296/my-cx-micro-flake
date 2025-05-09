package views

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ViewManager struct {
	DbConn *gorm.DB
	Logger *zap.Logger
}
