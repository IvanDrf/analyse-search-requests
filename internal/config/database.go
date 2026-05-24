package config

import (
	"fmt"
	"time"
)

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database int    `yaml:"db"`

	DuplicateTime time.Duration `yaml:"duplicate_time"`
	BadWordKey    string        `yaml:"bad_word_key"`
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
