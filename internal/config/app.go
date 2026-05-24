package config

import "time"

type AppConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	RequestTime    time.Duration `yaml:"request_time"`
	SearchInterval time.Duration `yaml:"search_duration"`

	LoggerLevel string `yaml:"logger_level"`
}
