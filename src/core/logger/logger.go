package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"gbase/src/core/helper"
	"gbase/src/def"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Log *log.Logger

func init() {
	var writers []io.Writer

	// 包了 lumberjack
	daily := newDailyWriter()
	writers = append(writers, newAnsiStripWriter(daily))

	if !helper.IsRelease() {
		writers = append(writers, os.Stdout) // dev 模式同時打到 console
	}

	Log = log.New(io.MultiWriter(writers...), "", log.LstdFlags|log.Lshortfile)
}

// --------------------------------------------------------------------------------
// 建立 dailyWriter：包一層 lumberjack.Logger，跨日就自動換檔
// --------------------------------------------------------------------------------
type dailyWriter struct {
	cur  *lumberjack.Logger // 目前使用中的 lumberjack
	date atomic.Value       // yyyyMMdd 字串
}

func newDailyWriter() io.Writer {
	w := &dailyWriter{}
	w.rotate() // 啟動時先建今天的檔案
	return w
}

func (w *dailyWriter) Write(p []byte) (int, error) {
	// 每次寫入前確認是否跨日
	if w.needRotate() {
		w.rotate()
	}
	return w.cur.Write(p)
}

func (w *dailyWriter) needRotate() bool {
	today := time.Now().Format(def.YYYY_MM_DD)
	if d, ok := w.date.Load().(string); ok {
		return d != today
	}
	return true
}

func (w *dailyWriter) rotate() {
	today := time.Now().Format(def.YYYY_MM_DD)

	dir, _ := helper.GetStorageDirOrCreate("log")
	filename := filepath.Join(dir, "log-"+today+".log")

	w.cur = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100, // 每個檔案上限 100 MB，再大就額外切一支
		MaxBackups: 30,  // 最多保留 30 支歷史檔
		MaxAge:     30,  // 超過 30 天自動刪
		Compress:   false,
	}
	w.date.Store(today)
}

func GetLogDir() string {
	dir, _ := helper.GetStorageDirOrCreate("log")

	return dir
}
