package config

type PostgreSQLConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`

	Database string `yaml:"database"`
}
