package logger

import (
	"context"
	"log/slog"
)

func NewSlogLogger(level slog.Leveler) *slog.Logger {
	if level == nil {
		level = slog.LevelInfo
	}

	return slog.New(&slogHandler{level: level})
}

type slogHandler struct {
	level slog.Leveler
}

func (h *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *slogHandler) Handle(_ context.Context, record slog.Record) error {
	Log.Printf("[WAILS:%s] %s", record.Level.String(), record.Message)
	return nil
}

func (h *slogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *slogHandler) WithGroup(_ string) slog.Handler {
	return h
}
