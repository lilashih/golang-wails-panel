package config

type ProjectConfig struct {
	BasePath string `env:"PROJECT_BASE_PATH" envDefault:"./"`
	Json     string `env:"PROJECT_JSON" envDefault:"project.json"`
}
