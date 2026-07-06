package app

import (
	"gbase/src/core/logger"
	"gbase/src/service/log_viewer"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func RegisterEvents() {
	application.RegisterEvent[string]("time")
	application.RegisterEvent[logger.Event]("log")
	application.RegisterEvent[log_viewer.LogTailEvent]("logViewer:tail")
}
