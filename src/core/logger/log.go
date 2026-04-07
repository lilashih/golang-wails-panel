package logger

import (
	"gbase/src/core/helper"
	"io"
	"log"
	"os"
)

var Log *log.Logger

const logFilePrefix = "log-"

func init() {
	var newLogger = &Logger{}

	std, err := newLogger.Init(newLogFactory)
	if err == nil {
		Log = std
		return
	}

	log.Printf("初始化log失敗，改用標準錯誤輸出: %v", err)
	Log = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
}

func newLogFactory() (*log.Logger, error) {
	daily, err := newDailyWriter(logFilePrefix)
	if err != nil {
		return nil, err
	}

	var writers []io.Writer
	writers = append(writers, newAnsiStripWriter(daily))

	if !helper.IsRelease() {
		writers = append(writers, os.Stdout)
	}

	return log.New(io.MultiWriter(writers...), "", log.LstdFlags|log.Lshortfile), nil
}
