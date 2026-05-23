package config

import "fmt"

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database int    `yaml:"db"`
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
