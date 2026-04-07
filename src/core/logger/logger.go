package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gbase/src/core/config"
	"gbase/src/core/def"
	"gbase/src/core/helper"
)

const (
	logDirName = "log"
	logFileExt = ".log"
)

type Logger struct {
	once sync.Once
	log  *log.Logger
	err  error
}

func (l *Logger) Init(factory func() (*log.Logger, error)) (*log.Logger, error) {
	l.once.Do(func() {
		l.log, l.err = factory()
	})

	return l.log, l.err
}

type dailyWriter struct {
	mu              sync.Mutex
	cur             *os.File
	currentDate     string
	currentPath     string
	dir             string
	filePrefix      string
	compressEnabled bool
	cleanupRunning  bool
	lastCleanup     string
}

func newDailyWriter(filePrefix string, dirs ...string) (*dailyWriter, error) {
	dir, err := getLogDir(dirs...)
	if err != nil {
		return nil, err
	}

	w := &dailyWriter{
		dir:             dir,
		filePrefix:      filePrefix,
		compressEnabled: isCompressEnabled(),
	}

	if err := w.rotate(time.Now()); err != nil {
		return nil, err
	}

	if w.compressEnabled {
		if err := w.cleanupOldLogs(w.currentPath); err != nil {
			log.Printf("初始化時壓縮日誌失敗: %v", err)
		} else {
			w.lastCleanup = w.currentDate
		}
	}

	return w, nil
}

func (w *dailyWriter) Write(p []byte) (int, error) {
	now := time.Now()
	today := now.Format(def.YYYY_MM_DD)

	w.mu.Lock()
	if w.currentDate != today {
		if err := w.rotate(now); err != nil {
			w.mu.Unlock()
			return 0, err
		}
	}

	if w.cur == nil {
		w.mu.Unlock()
		return 0, fmt.Errorf("初始化日誌寫入器失敗")
	}

	n, err := w.cur.Write(p)

	shouldCleanup := w.compressEnabled && !w.cleanupRunning && w.lastCleanup != today
	currentFile := w.currentPath
	if shouldCleanup {
		w.cleanupRunning = true
	}
	w.mu.Unlock() // 檢查並設定完畢後再釋放鎖，避免同時看到 !w.cleanupRunning 並同時啟動清理

	if shouldCleanup {
		go w.runCleanup(today, currentFile)
	}

	return n, err
}

func (w *dailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cur != nil {
		err := w.cur.Close()
		w.cur = nil
		w.currentPath = ""
		return err
	}

	return nil
}

func (w *dailyWriter) rotate(t time.Time) error {
	today := t.Format(def.YYYY_MM_DD)
	filename := w.logFilePath(today)

	newFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	old := w.cur
	w.cur = newFile
	w.currentDate = today
	w.currentPath = filename

	if old != nil {
		if err := old.Close(); err != nil {
			log.Printf("%s 檔案無法關閉: %v", old.Name(), err)
		}
	}

	return nil
}

func (w *dailyWriter) runCleanup(today, currentFile string) {
	err := w.cleanupOldLogs(currentFile)
	if err != nil {
		log.Printf("壓縮過期日誌失敗: %v", err)
	}

	w.mu.Lock()
	w.cleanupRunning = false
	if err == nil {
		w.lastCleanup = today
	}
	w.mu.Unlock()
}

func (w *dailyWriter) cleanupOldLogs(currentFile string) error {
	if !w.compressEnabled {
		return nil
	}

	return compressOldLogs(w.dir, w.filePrefix, currentFile)
}

func (w *dailyWriter) logFilePath(date string) string {
	return filepath.Join(w.dir, fmt.Sprintf("%s%s%s", w.filePrefix, date, logFileExt))
}

func getLogDir(dirs ...string) (string, error) {
	parts := append([]string{logDirName}, dirs...)
	return helper.GetStorageDirOrCreate(parts...)
}

func isLogFile(name, filePrefix string) bool {
	base := strings.TrimSuffix(name, ".gz")

	if !strings.HasPrefix(base, filePrefix) || !strings.HasSuffix(base, logFileExt) {
		return false
	}

	datePart := strings.TrimSuffix(strings.TrimPrefix(base, filePrefix), logFileExt)
	_, err := time.Parse(def.YYYY_MM_DD, datePart)
	if err != nil {
		return false
	}

	return true
}

func compressOldLogs(dir, filePrefix, current string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(dir, name)

		if path == current {
			continue
		}
		if strings.HasSuffix(name, ".gz") {
			continue
		}
		if !isLogFile(name, filePrefix) {
			continue
		}

		if err := compressLogFile(path); err != nil {
			log.Printf("壓縮日誌失敗: %s, err=%v", path, err)
		}
	}

	return nil
}

func compressLogFile(path string) (err error) {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if src != nil {
			_ = src.Close()
		}
	}()

	info, err := src.Stat()
	if err != nil {
		return err
	}

	gzPath := path + ".gz"
	if _, err := os.Stat(gzPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	dst, err := os.Create(gzPath)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if dst != nil {
			_ = dst.Close()
		}
		if !success {
			_ = os.Remove(gzPath)
		}
	}()

	gz := gzip.NewWriter(dst)
	gz.Name = filepath.Base(path)
	gz.ModTime = info.ModTime()

	if _, err = io.Copy(gz, src); err != nil {
		_ = gz.Close()
		return err
	}

	if err = gz.Close(); err != nil {
		return err
	}

	if err = dst.Close(); err != nil {
		return err
	}
	dst = nil

	if err = src.Close(); err != nil {
		return err
	}
	src = nil

	if err = os.Chtimes(gzPath, info.ModTime(), info.ModTime()); err != nil {
		return err
	}

	if err = os.Remove(path); err != nil {
		return err
	}

	success = true
	return nil
}

func isCompressEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(config.Logger.Compress), "true")
}

func GetLogDir() string {
	dir, _ := helper.GetStorageDirOrCreate(logDirName)

	return dir
}
