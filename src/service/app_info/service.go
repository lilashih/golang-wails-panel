package app_info

import "gbase/src/core/config"

type Service struct{}

type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Mode        string `json:"mode"`
	ProjectPath string `json:"project_path"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetAppInfo() AppInfo {
	return AppInfo{
		Name:        config.App.Name,
		Version:     config.App.Version,
		Mode:        config.App.Mode,
		ProjectPath: config.Project.BasePath,
	}
}
