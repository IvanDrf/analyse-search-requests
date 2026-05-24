package config

import "time"

type AppConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	RequestTime    time.Duration `yaml:"request_time"`
	SearchDuration time.Duration `yaml:"search_duration"`
	SearchInterval int           `yaml:"search_interval"`

	LoggerLevel string `yaml:"logger_level"`
}
