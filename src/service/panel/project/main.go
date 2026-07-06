package project

import (
	"encoding/json"
	"errors"
	"gbase/src/core/config"
	"gbase/src/core/helper"
	"gbase/src/core/logger"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// 從 projects 目錄載入專案。
//
// projects 支援以下三種路徑格式：
//   - 絕對路徑
//   - 相對路徑：若與目前專案同層可直接使用；若需多層 ../../ 導覽，建議改用絕對路徑以提高穩定性
//   - 捷徑：支援 Windows、Linux 捷徑
func NewProjects() ([]*Project, error) {
	var list []*Project

	path, entries, err := readProjectDir()
	if err != nil {
		return list, err
	}

	for _, entry := range entries {
		symlink, err := os.Readlink(filepath.Join(path, entry.Name())) // 檢查是否為捷徑

		if entry.IsDir() || err == nil {
			dir := entry.Name()
			if err == nil {
				logger.Log.Printf("讀取 %s 捷徑 %s ", dir, symlink)
				dir = symlink
			}

			cfgPath := filepath.Join(path, dir, config.Project.Json)
			data, err := os.ReadFile(cfgPath)
			if err != nil {
				logger.Log.Printf("讀取 %s 失敗: %v", entry.Name(), err)
				continue
			}

			var configs []ProjectConfig
			if err := json.Unmarshal(data, &configs); err != nil {
				logger.Log.Printf("解析 %s 失敗: project.json 必須為陣列格式: %v", entry.Name(), err)
				continue
			}

			cfg, ok := findConfigForCurrentOS(configs)
			if !ok {
				logger.Log.Printf("專案 %s 找不到 %s 對應的設定，已跳過", entry.Name(), runtime.GOOS)
				continue
			}

			if cfg.Title == "" || cfg.Key == "" || cfg.Type == "" || cfg.Start == "" || cfg.Stop == "" {
				logger.Log.Printf("設定檔 %s 缺少必要欄位，已跳過", entry.Name())
				continue
			}

			p := buildProject(cfg, filepath.Join(path, dir))
			if p == nil {
				logger.Log.Printf("不支援的專案類型: %s", cfg.Type)
				continue
			}

			p.CheckRunning()

			list = append(list, p)
		}
	}

	return list, nil
}

func findConfigForCurrentOS(configs []ProjectConfig) (ProjectConfig, bool) {
	normalizedTargetOS := normalizeProjectOS(runtime.GOOS)

	for _, cfg := range configs {
		if normalizeProjectOS(cfg.OS) == normalizedTargetOS {
			return cfg, true
		}
	}

	return ProjectConfig{}, false
}

func normalizeProjectOS(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func readProjectDir(elem ...string) (string, []fs.DirEntry, error) {
	entries := []fs.DirEntry{}

	basePath, err := helper.GetRuntimeBasePath()
	if err != nil {
		message := "開啟 projects 目錄失敗"
		logger.Log.Printf("%s：%v", message, err)
		return "", entries, errors.New(message)
	}

	path, err := resolveProjectDir(basePath, elem...)
	if err != nil {
		message := "解析 projects 目錄路徑失敗"
		logger.Log.Printf("%s：%v", message, err)
		return "", entries, errors.New(message)
	}

	entries, err = os.ReadDir(path)
	if err != nil {
		message := "讀取 projects 目錄失敗"
		logger.Log.Printf("%s：%v", message, err)
		return "", entries, errors.New(message)
	}

	return path, entries, nil
}

// 解析絕對路徑、相對路徑，直接使用 PROJECT_BASE_PATH 指定的目錄
// 1. PROJECT_BASE_PATH 如果是相對路徑，就用 helper.GetRuntimeBasePath() 補成絕對路徑。
// 2. 如果 PROJECT_BASE_PATH 本來就是絕對路徑，就直接使用。
// 3. Project.Path = filepath.Join(projectsBasePath, projectDir)。
//
// 唯一要注意的是 symlink：如果 entry 是 symlink，這段會把 dir = symlink：
// 因此不論 projects 目錄的捷徑 symlink 是相對路徑(Linux)或是絕對路徑(Windows)，最後都會轉成絕對路徑
func resolveProjectDir(basePath string, elem ...string) (string, error) {
	cfgBase := strings.TrimSpace(strings.Trim(config.Project.BasePath, "\"")) // 移掉前後空白、雙引號
	if cfgBase == "" {
		cfgBase = "."
	}

	if !filepath.IsAbs(cfgBase) {
		cfgBase = filepath.Join(basePath, cfgBase)
	}

	cfgBase = filepath.Clean(cfgBase)
	path := filepath.Join(append([]string{cfgBase}, elem...)...)

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", errors.New("設定的 PROJECT_BASE_PATH 不是資料夾")
	}

	return path, nil
}
