package config

type AppConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	LoggerLevel string `yaml:"string"`
}
