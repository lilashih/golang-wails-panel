package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

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

func GetLogDir() string {
	dir, _ := helper.GetStorageDirOrCreate(logDirName)

	return dir
}
