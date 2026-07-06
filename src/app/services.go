package app

import "github.com/wailsapp/wails/v3/pkg/application"

func (a *App) Services() []application.Service {
	return []application.Service{
		application.NewService(a.PanelService),
		application.NewService(a.AppInfoService),
		application.NewService(a.LogViewerService),
		application.NewService(a.ReadmeReaderService),
	}
}
