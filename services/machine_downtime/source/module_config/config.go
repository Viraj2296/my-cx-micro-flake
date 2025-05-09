package module_config

import (
	"go.cerex.io/transcendflow/config"
)

type JobServiceConfig struct {
	MachineHelpSignalPollingInterval int    `mapstructure:"machineHelpSignalPollingInterval"`
	EscalationPollingInterval        int    `mapstructure:"escalationPollingInterval"`
	EscalationEmailTemplate          string `mapstructure:"escalationEmailTemplate"`
	EnableEascalation                bool   `mapstructure:"enableEascalation"`
}

type DowntimeConfig struct {
	JobServiceConfig JobServiceConfig `mapstructure:"jobs"`
}

type MachineDowntimeConfig struct {
	AppConfig      config.AppConfig `mapstructure:"app"`
	DowntimeConfig DowntimeConfig   `mapstructure:"downtime"`
}
