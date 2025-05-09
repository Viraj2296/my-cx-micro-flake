package repository

import (
	"go.cerex.io/transcendflow/logging"
	"gorm.io/gorm"
)

type Repository struct {
	Database *gorm.DB
	Logger   *logging.Logger
}

// NewRepository initializes a new Repository instance
func NewRepository(logger *logging.Logger, db *gorm.DB) *Repository {
	return &Repository{
		Logger:   logger,
		Database: db,
	}
}
