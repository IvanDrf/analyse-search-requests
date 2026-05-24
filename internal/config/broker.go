package config

import "fmt"

type RabbitMQConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`

	Queue   string `yaml:"queue"`
	Workers int    `yaml:"workers"`
}

func (r *RabbitMQConfig) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
}
