package config

type RabbitMQConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`

	Queue string `yaml:"queue"`
}
