package test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
	_ "unsafe"

	_ "gbase/src/core/logger"
)

//go:linkname compressLogFile gbase/src/core/logger.compressLogFile
func compressLogFile(path string) error

//go:linkname compressOldLogs gbase/src/core/logger.compressOldLogs
func compressOldLogs(dir, filePrefix, current string) error

func TestCompressLogFileCreatesGzipAndRemovesSource(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log-2026-04-06.log")
	content := "hello logger\n"

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write source log: %v", err)
	}

	if err := compressLogFile(path); err != nil {
		t.Fatalf("compress log file: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("source log should be removed, got err=%v", err)
	}

	gzPath := path + ".gz"
	f, err := os.Open(gzPath)
	if err != nil {
		t.Fatalf("open gzip log: %v", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("new gzip reader: %v", err)
	}
	defer gz.Close()

	data, err := io.ReadAll(gz)
	if err != nil {
		t.Fatalf("read gzip content: %v", err)
	}

	if string(data) != content {
		t.Fatalf("unexpected gzip content: got %q want %q", string(data), content)
	}
}

func TestCompressOldLogsSkipsCurrentFile(t *testing.T) {
	dir := t.TempDir()
	currentPath := filepath.Join(dir, "log-2026-04-07.log")
	oldPath := filepath.Join(dir, "log-2026-04-06.log")

	if err := os.WriteFile(currentPath, []byte("today"), 0600); err != nil {
		t.Fatalf("write current log: %v", err)
	}
	if err := os.WriteFile(oldPath, []byte("yesterday"), 0600); err != nil {
		t.Fatalf("write old log: %v", err)
	}

	if err := compressOldLogs(dir, "log-", currentPath); err != nil {
		t.Fatalf("compress old logs: %v", err)
	}

	if _, err := os.Stat(currentPath); err != nil {
		t.Fatalf("current log should remain untouched: %v", err)
	}

	if _, err := os.Stat(oldPath + ".gz"); err != nil {
		t.Fatalf("old log should be compressed: %v", err)
	}

	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("old source log should be removed, got err=%v", err)
	}
}
