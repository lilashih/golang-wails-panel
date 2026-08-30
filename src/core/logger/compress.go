package logger

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gbase/src/core/config"
	"gbase/src/core/def"
)

func isLogFile(name, filePrefix string) bool {
	base := strings.TrimSuffix(name, ".gz")

	if !strings.HasPrefix(base, filePrefix) || !strings.HasSuffix(base, logFileExt) {
		return false
	}

	datePart := strings.TrimSuffix(strings.TrimPrefix(base, filePrefix), logFileExt)
	_, err := time.Parse(def.YYYY_MM_DD, datePart)
	return err == nil
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
		if path == current || strings.HasSuffix(name, ".gz") || !isLogFile(name, filePrefix) {
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
