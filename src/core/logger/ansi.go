package logger

import (
	"io"
	"regexp"
)

// 編譯一次即可重用
var ansiRegexp = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

// 移除命令提示字元的文字顏色
// ansiStripWriter 會把傳入資料中的 ANSI 轉義序列移除後再寫入底層 writer
type ansiStripWriter struct{ dst io.Writer }

func newAnsiStripWriter(w io.Writer) io.Writer { return &ansiStripWriter{dst: w} }

func (w *ansiStripWriter) Write(p []byte) (int, error) {
	clean := ansiRegexp.ReplaceAll(p, nil)
	return w.dst.Write(clean)
}
