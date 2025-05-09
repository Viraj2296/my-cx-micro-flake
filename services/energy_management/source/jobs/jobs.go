package jobs

import (
	"cx-micro-flake/services/energy_management/source/repository"
	"go.cerex.io/transcendflow/logging"
	"gorm.io/gorm"
)

type Jobs struct {
	Repository *repository.Repository
	Database   *gorm.DB
	Logger     *logging.Logger
}

// NewJobs initializes a new Jobs instance
func NewJobs(logger *logging.Logger, db *gorm.DB, repository *repository.Repository) *Jobs {
	return &Jobs{
		Logger:     logger,
		Database:   db,
		Repository: repository,
	}
}
func (v *Jobs) Init() {
}
