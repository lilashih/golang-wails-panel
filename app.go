package main

import (
	"context"
	"gbase/src/core/logger"
	"gbase/src/service/log_viewer"
	"gbase/src/service/panel"
)

// App struct
type App struct {
	ctx              context.Context
	PanelService     *panel.PanelService
	LogViewerService *log_viewer.LogViewerService
}

func NewApp() *App {
	return &App{
		PanelService:     panel.NewPanelService(),
		LogViewerService: log_viewer.NewLogViewerService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	logger.RegisterCtx(ctx)

	a.PanelService.Startup(a.ctx)
	a.LogViewerService.Startup(a.ctx)
}

func (a *App) shutdown(ctx context.Context) {
	a.LogViewerService.Shutdown(a.ctx)
}
