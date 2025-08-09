package main

import (
	"embed"
	"gbase/src/core/config"
	"gbase/src/core/logger"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	wlog "github.com/wailsapp/wails/v2/pkg/logger"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  config.App.Name,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
			app.PanelService,
			app.LogViewerService,
		},
		Logger:             logger.NewWailsLog(),
		LogLevel:           wlog.TRACE,
		LogLevelProduction: wlog.INFO, // 生產環境只顯示 INFO 以上的訊息
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
