package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"time"

	backendapp "gbase/src/app"
	"gbase/src/core/config"
	"gbase/src/core/logger"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	backendapp.RegisterEvents()
}

func main() {
	config.Initialize()
	logger.Initialize()
	logger.Log.Print("系統已啟動")

	backend := backendapp.New()

	app := application.New(application.Options{
		Name:        config.App.Name,
		Description: "",
		Services:    backend.Services(),
		Logger:      logger.NewSlogLogger(slog.LevelInfo),
		LogLevel:    slog.LevelInfo,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	backend.ConnectWails(app)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  fmt.Sprintf("%s - %s", config.App.Name, config.App.Version),
		Width:  1400,
		Height: 660,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
