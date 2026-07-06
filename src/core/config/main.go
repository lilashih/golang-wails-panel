package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var App AppConfig
var Project ProjectConfig
var Logger LoggerConfig
var LoadErrors []error // 初始化的錯誤，在 log 檔建立前就出錯的錯誤

type configs struct {
	*AppConfig
	*ProjectConfig
	*LoggerConfig
}

func init() {
	envFile := ".env"

	set(&configs{&App, &Project, &Logger}) // 先載入預設值

	// 即使 .env 檔案不存在，程式也能正常使用預設值啟動
	if file, err := openEnvFile(envFile); err == nil {
		defer file.Close()
		setFromEnvFile(file, &configs{&App, &Project, &Logger})
	}
}

func openEnvFile(envFile string) (*os.File, error) {
	file, err := os.Open(envFile)
	if err == nil {
		return file, nil
	}

	executablePath, executableErr := os.Executable()
	if executableErr != nil {
		return nil, err
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, err
	}

	executableDir := filepath.Dir(executablePath)
	// 處理 windows 捷徑路徑問題
	if !shouldFallbackToExecutableEnv(cwd, executableDir) {
		return nil, err
	}

	return os.Open(filepath.Join(executableDir, envFile))
}

func shouldFallbackToExecutableEnv(cwd string, executableDir string) bool {
	if App.Mode != "release" {
		return false
	}

	// air 等開發工具通常會在專案工作目錄底下執行暫存二進制檔，例如 ./bin/app.exe
	// 這種情況下如果 cwd/.env 不存在，應該使用預設值，不應讀取可能過期的 bin/.env
	return !isPathWithinOrEqual(cwd, executableDir)
}

func isPathWithinOrEqual(parent string, child string) bool {
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return false
	}

	childAbs, err := filepath.Abs(child)
	if err != nil {
		return false
	}

	rel, err := filepath.Rel(parentAbs, childAbs)
	if err != nil {
		return false
	}

	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && !filepath.IsAbs(rel))
}

func set(structure interface{}) {
	if err := env.Parse(structure); err != nil {
		recordLoadError(err)
	}
}

func setFromEnvFile(file *os.File, structure interface{}) {
	vars, err := readEnvFile(file)
	if err != nil {
		recordLoadError(err)
		return
	}

	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			recordLoadError(fmt.Errorf("dotenv: cannot set %s; err: %w", key, err))
		}
	}

	set(structure)
}

func readEnvFile(file *os.File) (map[string]string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("dotenv: read %s; err: %w", file.Name(), err)
	}

	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	vars, err := godotenv.UnmarshalBytes(content)
	if err != nil {
		return nil, fmt.Errorf("dotenv: parse %s; err: %w", file.Name(), err)
	}

	return vars, nil
}

func recordLoadError(err error) {
	LoadErrors = append(LoadErrors, err)
	log.Printf("%+v", err)
}
