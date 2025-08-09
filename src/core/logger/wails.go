package logger

import (
	"context"
	"time"

	wlog "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ----------------------
//
//	供 App.startup 註冊
//
// ----------------------
var ctxList = make([]context.Context, 0, 1)

func RegisterCtx(ctx context.Context) {
	ctxList = append(ctxList, ctx)
}

// 對前端送出事件
func emit(level, msg string) {
	for _, c := range ctxList {
		runtime.EventsEmit(c, "log", map[string]interface{}{
			"level":     level,
			"message":   msg,
			"timestamp": time.Now().Format("15:04:05"),
		})
	}
}

// Adapter 把 Wails 要的七個方法轉給 Log
type Adapter struct{}

func (*Adapter) Print(msg string)   { emit("info", msg); Log.Printf("[INFO]  %s", msg) }
func (*Adapter) Trace(msg string)   { emit("trace", msg); Log.Printf("[TRACE] %s", msg) }
func (*Adapter) Debug(msg string)   { emit("debug", msg); Log.Printf("[DEBUG] %s", msg) }
func (*Adapter) Info(msg string)    { emit("info", msg); Log.Printf("[INFO]  %s", msg) }
func (*Adapter) Warning(msg string) { emit("warning", msg); Log.Printf("[WARN]  %s", msg) }
func (*Adapter) Error(msg string)   { emit("error", msg); Log.Printf("[ERROR] %s", msg) }
func (*Adapter) Fatal(msg string)   { emit("fatal", msg); Log.Fatalf("[FATAL] %s", msg) }

// New 傳回符合 wails Logger 介面的實例
func NewWailsLog() wlog.Logger {
	return &Adapter{}
}
