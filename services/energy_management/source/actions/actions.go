package actions

import (
	"cx-micro-flake/services/energy_management/source/repository"

	"go.cerex.io/transcendflow/logging"
	time_series "go.cerex.io/transcendflow/time-series"
	"gorm.io/gorm"
)

type Actions struct {
	Logger            *logging.Logger
	Database          *gorm.DB
	Repository        *repository.Repository
	RealtimeDBManager *time_series.RealtimeDBManager
	EnergyStartDate   string
	IdlePowerCache    map[int][]CachedIdleData
}

// NewActions initializes a new Actions instance
func NewActions(logger *logging.Logger, db *gorm.DB, repository *repository.Repository, RealtimeDBManager *time_series.RealtimeDBManager, energyStartDate string) *Actions {
	var defaultIdlePowerCache = make(map[int][]CachedIdleData)
	return &Actions{
		Logger:            logger,
		Database:          db,
		Repository:        repository,
		RealtimeDBManager: RealtimeDBManager,
		EnergyStartDate:   energyStartDate,
		IdlePowerCache:    defaultIdlePowerCache,
	}
}
func (v *Actions) Init() {
}
