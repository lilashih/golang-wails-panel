package project

import (
	"encoding/json"
	"errors"
	"gbase/src/core/config"
	"gbase/src/core/logger"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Project struct {
	Config ProjectConfig

	Path    string `json:"path"`
	Running bool   `json:"running"`

	Install      func() error `json:"-"`
	Start        func() error `json:"-"`
	Stop         func() error `json:"-"`
	CheckRunning func()       `json:"-"`
}

type ProjectConfig struct {
	Title   string `json:"title"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Start   string `json:"start"`
	Stop    string `json:"stop"`
	Install string `json:"install"`
}

// 從projects目錄載入專案
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

			var cfg ProjectConfig
			if err := json.Unmarshal(data, &cfg); err != nil {
				logger.Log.Printf("解析 %s 失敗: %v", entry.Name(), err)
				continue
			}

			if cfg.Title == "" || cfg.Key == "" || cfg.Start == "" || cfg.Stop == "" {
				logger.Log.Printf("設定檔 %s 缺少必要欄位，已跳過", entry.Name())
				continue
			}

			var p *Project

			switch cfg.Type {
			case "pm2":
				p = NewPm2(cfg, filepath.Join(path, dir))
			case "exe":
				p = NewExe(cfg, filepath.Join(path, dir))
			default:
				logger.Log.Printf("不支援的專案類型: %s", cfg.Type)
				continue
			}

			p.CheckRunning()

			list = append(list, p)
		}
	}

	return list, nil
}

func readProjectDir(elem ...string) (string, []fs.DirEntry, error) {
	entries := []fs.DirEntry{}

	basePath, err := os.Getwd()
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

// 啟動或停止時不馬上檢查狀態，等待3秒後再查詢
func waitForRunningState(p *Project, expected bool) {
	deadline := time.Now().Add(3 * time.Second)

	for {
		p.CheckRunning()
		if p.Running == expected || time.Now().After(deadline) {
			return
		}

		time.Sleep(250 * time.Millisecond)
	}
}
