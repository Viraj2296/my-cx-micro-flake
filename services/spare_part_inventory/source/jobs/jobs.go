package jobs

import (
	"cx-micro-flake/services/spare_part_inventory/source/consts"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"cx-micro-flake/services/spare_part_inventory/source/module_config"
	"cx-micro-flake/services/spare_part_inventory/source/repository"
	"io/ioutil"

	"go.cerex.io/transcendflow/component_processor"
	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Jobs struct {
	JobConfig               module_config.JobsConfig
	Repository              repository.Repository
	Logger                  *logging.Logger
	Database                *gorm.DB
	settingInfo             *models.SparePartInventorySettingInfo
	ComponentManager        *component_processor.ComponentManager
	EscalationEmailTemplate string
}

// NewJobs initializes a new Jobs instance
func NewJobs(logger *logging.Logger, db *gorm.DB, JobsConfig module_config.JobsConfig, repository repository.Repository, componentManager *component_processor.ComponentManager) *Jobs {
	return &Jobs{
		Logger:           logger,
		Database:         db,
		JobConfig:        JobsConfig,
		Repository:       repository,
		ComponentManager: componentManager,
	}
}
func (v *Jobs) Init() {
	if v.JobConfig.InventoryLimitPollingInterval == 0 {
		v.Logger.Warn("InventoryLimitPollingInterval cannot be 0, taking it 30 seconds")
		v.JobConfig.InventoryLimitPollingInterval = 30
	}
	var settingTable = v.ComponentManager.GetTargetTable(consts.SparePartInventorySettingComponent)
	err, c := v.Repository.GetResource(settingTable, 1)
	if err != nil {
		v.Logger.Error("getting spare part inventory setting failed", logging.Error(err))
	}
	v.settingInfo = models.GetSparePartInventorySettingInfo(c.ObjectInfo)
	content, err := ioutil.ReadFile(v.JobConfig.EscalationEmailTemplate)
	if err != nil {
		v.Logger.Error("error reading email template", zap.String("error", err.Error()))
		v.EscalationEmailTemplate = "" // don't send email if it is empty
	} else {
		v.EscalationEmailTemplate = string(content)
	}
}
