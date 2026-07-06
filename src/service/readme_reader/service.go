package readme_reader

import (
	"fmt"
	"gbase/src/core/helper"
	"os"
	"path/filepath"
	"strings"
)

type Service struct{}

type DocumentOption struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
}

type DocumentContent struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

var documents = []DocumentOption{
	{Key: "zh-TW", Title: "繁體中文", Filename: "README.md"},
	{Key: "en", Title: "English", Filename: "README.en.md"},
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListDocuments() ([]DocumentOption, error) {
	if helper.IsRelease() {
		return []DocumentOption{}, nil
	}
	basePath, err := helper.GetRuntimeBasePath()
	if err != nil {
		return nil, fmt.Errorf("取得執行路徑失敗：%w", err)
	}

	result := make([]DocumentOption, 0, len(documents))
	for _, doc := range documents {
		path, err := resolveDocumentPath(basePath, doc.Filename)
		if err != nil {
			continue
		}
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			result = append(result, doc)
		}
	}

	return result, nil
}

func (s *Service) GetDocument(key string) (*DocumentContent, error) {
	if helper.IsRelease() {
		return nil, fmt.Errorf("release 模式不提供 README 閱讀功能")
	}
	basePath, err := helper.GetRuntimeBasePath()
	if err != nil {
		return nil, fmt.Errorf("取得執行路徑失敗：%w", err)
	}

	doc, ok := findDocument(key)
	if !ok {
		return nil, fmt.Errorf("找不到文件：%s", key)
	}

	path, err := resolveDocumentPath(basePath, doc.Filename)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("讀取文件失敗：%w", err)
	}

	return &DocumentContent{
		Key:      doc.Key,
		Title:    doc.Title,
		Filename: doc.Filename,
		Content:  string(data),
	}, nil
}

func findDocument(key string) (DocumentOption, bool) {
	for _, doc := range documents {
		if doc.Key == key {
			return doc, true
		}
	}

	return DocumentOption{}, false
}

func resolveDocumentPath(basePath string, filename string) (string, error) {
	baseAbs, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("解析執行路徑失敗：%w", err)
	}

	path := filepath.Clean(filepath.Join(baseAbs, filename))
	rel, err := filepath.Rel(baseAbs, path)
	if err != nil {
		return "", fmt.Errorf("解析文件路徑失敗：%w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("文件路徑超出允許範圍：%s", filename)
	}

	return path, nil
}
