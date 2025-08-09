package log_viewer

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"gbase/src/core/logger"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogViewerService struct {
	ctx context.Context

	mu       sync.Mutex
	curFile  string
	tailStop context.CancelFunc
	dir      string // 專案的 log 目錄
}

func NewLogViewerService() *LogViewerService {
	return &LogViewerService{
		dir: logger.GetLogDir(),
	}
}

func (s *LogViewerService) Startup(ctx context.Context) {
	s.ctx = ctx
}
func (s *LogViewerService) Shutdown(ctx context.Context) {
	s.StopTail()
}

// ===================== 檔案選取 / 清單 =====================

func (s *LogViewerService) PickFile() (string, error) {
	path, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "選擇 Log 檔案",
		Filters: []runtime.FileFilter{
			{DisplayName: "Log", Pattern: "*.log;*.txt;*.*"},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", errors.New("取消選擇")
	}
	s.curFile = filepath.Clean(path)
	return s.curFile, nil
}

type LogFileItem struct {
	Name string `json:"name"` // 檔名
	Path string `json:"path"` // 完整路徑
}

func (s *LogViewerService) ListLogFiles() ([]LogFileItem, error) {
	var files []LogFileItem
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		low := strings.ToLower(name)
		if strings.HasSuffix(low, ".log") || strings.HasSuffix(low, ".txt") {
			files = append(files, LogFileItem{
				Name: name,
				Path: filepath.Clean(filepath.Join(s.dir, name)),
			})
		}
	}
	sort.Slice(files, func(i, j int) bool {
		ni := strings.ToLower(files[i].Name)
		nj := strings.ToLower(files[j].Name)
		if ni == nj {
			return files[i].Path < files[j].Path
		}
		return ni < nj
	})
	return files, nil
}

func (s *LogViewerService) SetCurrentFile(p string) (string, error) {
	if p == "" {
		return "", errors.New("無效的檔案路徑")
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(s.dir, p)
	}
	p = filepath.Clean(p)
	if _, err := os.Stat(p); err != nil {
		return "", err
	}
	s.curFile = p
	return p, nil
}

// ===================== 讀取（由頭往後；固定行數） =====================

type ChunkResult struct {
	Lines     []string `json:"lines"`
	NextStart int64    `json:"nextStart"` // 下一次從哪個位移開始讀
	EOF       bool     `json:"eof"`
}

// ReadFirstLines：從檔案「開頭」讀，最多讀 maxLines 行
func (s *LogViewerService) ReadFirstLines(path string, maxLines int) (*ChunkResult, error) {
	if path == "" {
		path = s.curFile
	}
	if path == "" {
		return nil, errors.New("無效的檔案路徑")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines, next, eof, err := readLinesFrom(f, 0, maxLines)
	if err != nil {
		return nil, err
	}
	return &ChunkResult{Lines: lines, NextStart: next, EOF: eof}, nil
}

// LoadNext：從 afterPos 之後往後讀，最多讀 maxLines 行
func (s *LogViewerService) LoadNext(path string, afterPos int64, maxLines int) (*ChunkResult, error) {
	if path == "" {
		path = s.curFile
	}
	if path == "" {
		return nil, errors.New("無效的檔案路徑")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines, next, eof, err := readLinesFrom(f, afterPos, maxLines)
	if err != nil {
		return nil, err
	}
	return &ChunkResult{Lines: lines, NextStart: next, EOF: eof}, nil
}

// 低配置且可處理超長行：逐行讀並累計位移
func readLinesFrom(f *os.File, start int64, maxLines int) (lines []string, next int64, eof bool, err error) {
	if maxLines <= 0 {
		return []string{}, start, false, nil
	}

	if _, err = f.Seek(start, io.SeekStart); err != nil {
		return nil, start, false, err
	}

	r := bufio.NewReaderSize(f, 256*1024) // 256KB buffer
	var off = start

	for len(lines) < maxLines {
		b, rerr := r.ReadBytes('\n') // 含 '\n'，若到 EOF 且最後一行無換行，會回傳殘餘
		if len(b) > 0 {
			off += int64(len(b))
			// 去掉行尾 '\n'（若沒有則原樣）
			if b[len(b)-1] == '\n' {
				lines = append(lines, string(b[:len(b)-1]))
			} else {
				lines = append(lines, string(b))
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				eof = true
				break
			}
			return lines, off, false, rerr
		}
	}

	return lines, off, eof, nil
}

// ===================== 計算總行數（含最後一行無換行的情況） =====================

func (s *LogViewerService) CountLines(path string) (int, error) {
	if path == "" {
		path = s.curFile
	}
	if path == "" {
		return 0, errors.New("無效的檔案路徑")
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fi, _ := f.Stat()
	size := fi.Size()

	buf := make([]byte, 32*1024)
	total := 0
	var lastByte byte
	var haveByte bool

	for {
		n, rerr := f.Read(buf)
		if n > 0 {
			total += bytes.Count(buf[:n], []byte{'\n'})
			lastByte = buf[n-1]
			haveByte = true
		}
		if rerr != nil {
			if rerr == io.EOF {
				break
			}
			return total, rerr
		}
	}

	// 若檔案非空且最後一個 byte 不是 '\n'，代表最後一行沒有換行，也要算一行
	if size > 0 && haveByte && lastByte != '\n' {
		total++
	}
	return total, nil
}

// ===================== 其它 =====================

func (s *LogViewerService) StopTail() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tailStop != nil {
		s.tailStop()
		s.tailStop = nil
	}
}
