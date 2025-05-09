package http

import (
	"cx-micro-flake/services/spare_part_inventory/source/repository"
	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/logging"
)

type Service struct {
	Logger           *logging.Logger
	Repository       repository.Repository
	ComponentManager *component_processor.ComponentManager
}

// NewHTTPService NewRPCService initializes a new Actions instance
func NewHTTPService(logger *logging.Logger, repository repository.Repository, componentManager *component_processor.ComponentManager) *Service {
	return &Service{
		Logger:           logger,
		Repository:       repository,
		ComponentManager: componentManager,
	}
}
func (v *Service) Init() {
}
