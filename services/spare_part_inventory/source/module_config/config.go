package module_config

import (
	"go.cerex.io/transcendflow/config"
)

type JobsConfig struct {
	InventoryLimitPollingInterval int    `mapstructure:"inventoryLimitPollingInterval"`
	EscalationEmailTemplate       string `mapstructure:"escalationEmailTemplate"`
}

type SparePartConfig struct {
	JobServiceConfig JobsConfig `mapstructure:"jobs"`
}
type SparePartInventoryConfig struct {
	AppConfig       config.AppConfig `mapstructure:"app"`
	SparePartConfig SparePartConfig  `mapstructure:"sparePartInventory"`
}
