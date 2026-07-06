package log_viewer

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"gbase/src/core/cmd/explore"
	"gbase/src/core/logger"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	ctx context.Context

	mu       sync.Mutex
	curFile  string
	tailStop context.CancelFunc
	dir      string
}

type TailEmitter interface {
	EmitLogTail(LogTailEvent)
}

var (
	emitterMu sync.RWMutex
	emitter   TailEmitter
)

func SetEmitter(next TailEmitter) {
	emitterMu.Lock()
	defer emitterMu.Unlock()
	emitter = next
}

type LogTailEvent struct {
	Path      string   `json:"path"`
	Lines     []string `json:"lines"`
	Timestamp string   `json:"timestamp"`
}

type LogFileItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ChunkResult struct {
	Lines     []string `json:"lines"`
	NextStart int64    `json:"nextStart"`
	EOF       bool     `json:"eof"`
}

func NewService() *Service {
	return &Service{
		dir: logger.GetLogDir(),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *Service) ServiceShutdown() error {
	s.StopTail()
	return nil
}

func (s *Service) ListLogFiles() ([]LogFileItem, error) {
	files := []LogFileItem{}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
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

func (s *Service) SetCurrentFile(path string) (string, error) {
	path, err := s.resolveFile(path)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.curFile = path
	s.mu.Unlock()

	return path, nil
}

func (s *Service) ReadFirstLines(path string, maxLines int) (*ChunkResult, error) {
	path, err := s.pathOrCurrent(path)
	if err != nil {
		return nil, err
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

func (s *Service) LoadNext(path string, afterPos int64, maxLines int) (*ChunkResult, error) {
	path, err := s.pathOrCurrent(path)
	if err != nil {
		return nil, err
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

func (s *Service) OpenLogFileLocation(path string) error {
	path, err := s.pathOrCurrent(path)
	if err != nil {
		return err
	}

	return explore.OpenDir(filepath.Dir(path))
}

func (s *Service) CountLines(path string) (int, error) {
	path, err := s.pathOrCurrent(path)
	if err != nil {
		return 0, err
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 32*1024)
	total := 0
	var lastByte byte
	var haveByte bool

	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			total += bytes.Count(buf[:n], []byte{'\n'})
			lastByte = buf[n-1]
			haveByte = true
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return total, readErr
		}
	}

	if info.Size() > 0 && haveByte && lastByte != '\n' {
		total++
	}

	return total, nil
}

func (s *Service) StartTail(path string) error {
	s.StopTail()

	path, err := s.SetCurrentFile(path)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()

	ctx := s.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	tailCtx, cancel := context.WithCancel(ctx)

	s.mu.Lock()
	s.tailStop = cancel
	s.mu.Unlock()

	go s.tailFile(tailCtx, path, offset)
	return nil
}

func (s *Service) StopTail() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tailStop != nil {
		s.tailStop()
		s.tailStop = nil
	}
}

func (s *Service) tailFile(ctx context.Context, path string, offset int64) {
	f, err := os.Open(path)
	if err != nil {
		logger.Error("開啟 log tail 失敗：%v", err)
		return
	}
	defer f.Close()

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		logger.Error("設定 log tail 位置失敗：%v", err)
		return
	}

	reader := bufio.NewReaderSize(f, 256*1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, readErr := reader.ReadString('\n')
		if line != "" {
			s.emitTail(path, []string{strings.TrimRight(line, "\r\n")})
		}

		if readErr != nil {
			if readErr == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			logger.Error("讀取 log tail 失敗：%v", readErr)
			return
		}
	}
}

func (s *Service) emitTail(path string, lines []string) {
	if len(lines) == 0 {
		return
	}

	emitterMu.RLock()
	current := emitter
	emitterMu.RUnlock()

	if current == nil {
		return
	}

	current.EmitLogTail(LogTailEvent{
		Path:      path,
		Lines:     lines,
		Timestamp: time.Now().Format("15:04:05"),
	})
}

func (s *Service) pathOrCurrent(path string) (string, error) {
	if path == "" {
		s.mu.Lock()
		path = s.curFile
		s.mu.Unlock()
	}
	return s.resolveFile(path)
}

func (s *Service) resolveFile(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("無效的檔案路徑")
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(s.dir, path)
	}
	path = filepath.Clean(path)

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}

func readLinesFrom(f *os.File, start int64, maxLines int) (lines []string, next int64, eof bool, err error) {
	if maxLines <= 0 {
		return []string{}, start, false, nil
	}

	if _, err = f.Seek(start, io.SeekStart); err != nil {
		return nil, start, false, err
	}

	reader := bufio.NewReaderSize(f, 256*1024)
	offset := start

	for len(lines) < maxLines {
		line, readErr := reader.ReadBytes('\n')
		if len(line) > 0 {
			offset += int64(len(line))
			if line[len(line)-1] == '\n' {
				lines = append(lines, string(line[:len(line)-1]))
			} else {
				lines = append(lines, string(line))
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				eof = true
				break
			}
			return lines, offset, false, readErr
		}
	}

	return lines, offset, eof, nil
}
