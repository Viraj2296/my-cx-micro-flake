package repository

import (
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/models"
	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/orm"
	"gorm.io/gorm"
	"strconv"
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

func (v *Repository) GetEmailEscalationCount(downtimeId int) int {
	var conditionalString = "object_info ->> '$.downtimeId' = " + strconv.Itoa(downtimeId)
	var escalationInterfaceCount = orm.CountByCondition(v.Database, consts.MachineDownTimeEmailEscalationTable, conditionalString)
	return int(escalationInterfaceCount)
}

func (v *Repository) GetEmailEscalationInfo(downtimeId int) *models.MachineDowntimeEmailEscalationInfo {
	var conditionalString = "object_info ->> '$.downtimeId' = " + strconv.Itoa(downtimeId)
	err, escalationInterface := orm.GetConditionalObjects(v.Database, consts.MachineDownTimeEmailEscalationTable, conditionalString)
	if err != nil {
		return nil
	}
	return models.GetMachineDowntimeEmailEscalationInfo(escalationInterface[len(escalationInterface)-1].ObjectInfo)
}
