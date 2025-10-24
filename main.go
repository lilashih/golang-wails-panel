package main

import (
	"context"
	"embed"
	"fmt"
	"gbase/src/core/config"
	"gbase/src/core/logger"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	wlog "github.com/wailsapp/wails/v2/pkg/logger"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed favicon.ico
var icon []byte

//─── 全域狀態 ───────────────────────────────────────────────

var (
	ctx      context.Context // Wails 的 context（給 runtime.* 用）
	ctxReady = make(chan struct{})
	hidden   = false // 追蹤視窗是否已被隱藏
)

//───────────────────────────────────────────────────────────

func main() {
	// 先跑 systray（非阻塞；實際阻塞在其內部 goroutine）
	go systray.Run(onReady, onExit)

	// 再啟動 Wails（此呼叫會阻塞直到程式結束）
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  fmt.Sprintf("%s - %s", config.App.Name, config.App.Version),
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		HideWindowOnClose: true, // 視窗關閉→縮到系統列（v2.8+）
		OnStartup: func(c context.Context) {
			ctx = c
			close(ctxReady) // 通知 systray：ctx 已可用
			app.startup(c)
		},
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

//───────────────── systray ────────────────────────────

func onReady() {

	<-ctxReady // **等 Wails 建好 ctx，再碰 runtime**

	systray.SetIcon(icon)
	systray.SetTooltip(fmt.Sprintf("%s - %s", config.App.Name, config.App.Version))

	mToggle := systray.AddMenuItem("顯示", "Show main window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "Quit the app")

	// 按鈕監聽
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				runtime.WindowShow(ctx)
			case <-mQuit.ClickedCh:
				// 先讓 Wails 關閉，再停 systray
				runtime.Quit(ctx)
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// 需要的話在此釋放資源
}
