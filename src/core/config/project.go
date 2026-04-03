package config

type ProjectConfig struct {
	BasePath string `env:"PROJECT_BASE_PATH" envDefault:"./projects"`
	Json     string `env:"PROJECT_JSON" envDefault:"project.json"`
}
