package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig        `yaml:"app"`
	Database PostgreSQLConfig `yaml:"database"`
	Broker   RabbitMQConfig   `yaml:"broker"`
}

func LoadFromYaml(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("can't find config with given path=%s", path)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("can't read config file, error=%s", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		log.Fatalf("can't parse config file into config struct, error=%s", err)
	}

	return config
}
