package config

type LoggerConfig struct {
	Compress string `env:"LOGGER_COMPRESS" envDefault:"true"`
}
