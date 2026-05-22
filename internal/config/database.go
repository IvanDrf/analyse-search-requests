package config

import "fmt"

type PostgreSQLConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`

	Database string `yaml:"database"`
}

func (d *PostgreSQLConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Database,
	)
}
