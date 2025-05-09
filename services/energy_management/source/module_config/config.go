package module_config

import (
	"go.cerex.io/transcendflow/config"
	time_series "go.cerex.io/transcendflow/time-series"
)

type Configuration struct {
	AppConfig       config.AppConfig             `mapstructure:"app"`
	InfluxConfig    time_series.RealtimeDBConfig `mapstructure:"dbManager"`
	EnergyStartDate string                       `mapstructure:"energyStartDate"`
}
