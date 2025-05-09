package api_access

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type APIService struct {
	DBConnection *gorm.DB
	Logger       *zap.Logger
}

func (v *APIService) Init(logger *zap.Logger) {
	v.Logger = logger
}
