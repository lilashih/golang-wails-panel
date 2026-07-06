package app

import (
	"gbase/src/core/logger"
	"gbase/src/service/app_info"
	"gbase/src/service/log_viewer"
	"gbase/src/service/panel"
	"gbase/src/service/readme_reader"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	PanelService        *panel.Service
	AppInfoService      *app_info.Service
	LogViewerService    *log_viewer.Service
	ReadmeReaderService *readme_reader.Service
}

func New() *App {
	return &App{
		PanelService:        panel.NewService(),
		AppInfoService:      app_info.NewService(),
		LogViewerService:    log_viewer.NewService(),
		ReadmeReaderService: readme_reader.NewService(),
	}
}

func (a *App) ConnectWails(wailsApp *application.App) {
	emitter := &eventEmitter{app: wailsApp}
	logger.SetEmitter(emitter)
	log_viewer.SetEmitter(emitter)
}

type eventEmitter struct {
	app *application.App
}

func (e *eventEmitter) EmitLog(event logger.Event) {
	e.app.Event.Emit("log", event)
}

func (e *eventEmitter) EmitLogTail(event log_viewer.LogTailEvent) {
	e.app.Event.Emit("logViewer:tail", event)
}
