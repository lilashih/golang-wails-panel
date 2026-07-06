package config

type AppConfig struct {
	Id       string `env:"APP_ID" envDefault:"Panel"`
	Version  string `envDefault:"2.0.0"`
	Mode     string `env:"APP_MODE" envDefault:"release"`
	BasePath string `env:"APP_BASE_PATH" envDefault:"./"`
	Name     string `env:"APP_NAME" envDefault:"Panel"`
}
