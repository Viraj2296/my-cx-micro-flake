package actions

import (
	"cx-micro-flake/services/machine_downtime/source/repository"
	"go.cerex.io/transcendflow/logging"
	"gorm.io/gorm"
)

type Actions struct {
	Logger     *logging.Logger
	Database   *gorm.DB
	Repository *repository.Repository
}

// NewActions initializes a new Actions instance
func NewActions(logger *logging.Logger, db *gorm.DB, repository *repository.Repository) *Actions {
	return &Actions{
		Logger:     logger,
		Database:   db,
		Repository: repository,
	}
}
func (v *Actions) Init() {
}
