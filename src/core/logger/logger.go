package logger

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"gbase/src/core/config"
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
	if isCompressEnabled() {
		today := time.Now().Format(def.YYYY_MM_DD)
		dir, _ := helper.GetStorageDirOrCreate("log")
		compressOldLogs(dir, filepath.Join(dir, "log-"+today+".log"))
	}

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
		Compress:   isCompressEnabled(),
	}
	w.date.Store(today)
}

func GetLogDir() string {
	dir, _ := helper.GetStorageDirOrCreate("log")

	return dir
}

func isCompressEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(config.Logger.Compress), "true")
}

func compressOldLogs(dir, current string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if path == current || filepath.Ext(path) != ".log" {
			continue
		}

		if err = compressLogFile(path); err != nil {
			log.Printf("compress log failed: %s, err=%v", path, err)
		}
	}
}

func compressLogFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}

	info, err := src.Stat()
	if err != nil {
		_ = src.Close()
		return err
	}

	gzPath := path + ".gz"
	dst, err := os.Create(gzPath)
	if err != nil {
		return err
	}

	ok := false
	defer func() {
		dst.Close()
		if !ok {
			_ = os.Remove(gzPath)
		}
	}()

	gz := gzip.NewWriter(dst)
	gz.Name = filepath.Base(path)
	gz.ModTime = info.ModTime()

	if _, err = io.Copy(gz, src); err != nil {
		_ = src.Close()
		_ = gz.Close()
		return err
	}
	if err = gz.Close(); err != nil {
		_ = src.Close()
		return err
	}
	if err = dst.Close(); err != nil {
		_ = src.Close()
		return err
	}
	if err = os.Chtimes(gzPath, info.ModTime(), info.ModTime()); err != nil {
		_ = src.Close()
		return err
	}
	if err = src.Close(); err != nil {
		return err
	}
	if err = os.Remove(path); err != nil {
		return err
	}

	ok = true
	return nil
}
