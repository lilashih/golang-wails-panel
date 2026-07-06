package logger

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Emitter interface {
	EmitLog(Event)
}

var (
	emitterMu sync.RWMutex
	emitter   Emitter
)

func SetEmitter(next Emitter) {
	emitterMu.Lock()
	defer emitterMu.Unlock()
	emitter = next
}

func Debug(format string, args ...any) {
	write("debug", format, args...)
}

func Info(format string, args ...any) {
	write("info", format, args...)
}

func Warn(format string, args ...any) {
	write("warning", format, args...)
}

func Error(format string, args ...any) {
	write("error", format, args...)
}

func write(level, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	Log.Printf("[%s] %s", level, message)
	emit(Event{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format("15:04:05"),
	})
}

func emit(event Event) {
	emitterMu.RLock()
	current := emitter
	emitterMu.RUnlock()

	if current != nil {
		current.EmitLog(event)
	}
}
