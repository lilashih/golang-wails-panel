package config

type ProjectConfig struct {
	BasePath string `env:"PROJECT_BASE_PATH" envDefault:"./projects"` // 專案目錄路徑以外的要用絕對路徑
	Json     string `env:"PROJECT_JSON" envDefault:"project.json"`
}
