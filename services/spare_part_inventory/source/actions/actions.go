package actions

import (
	"cx-micro-flake/services/spare_part_inventory/source/repository"
	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/logging"
	"gorm.io/gorm"
)

type Actions struct {
	Logger           *logging.Logger
	Database         *gorm.DB
	Repository       repository.Repository
	ComponentManager *component_processor.ComponentManager
}

// NewActions initializes a new Actions instance
func NewActions(logger *logging.Logger, db *gorm.DB, repository repository.Repository, componentManager *component_processor.ComponentManager) *Actions {
	return &Actions{
		Logger:           logger,
		Database:         db,
		Repository:       repository,
		ComponentManager: componentManager,
	}
}
func (v *Actions) Init() {
}
