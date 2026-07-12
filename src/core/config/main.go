package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

type envFileOpenResult struct {
	File     *os.File
	Path     string
	Errors   []error
	NotFound bool
}

func init() {
	envFile := ".env"

	set(&configs{&App, &Project, &Logger}) // 先載入預設值

	// 即使 .env 檔案不存在，程式也能正常使用預設值啟動
	result := openEnvFile(envFile)
	if result.File == nil {
		if !result.NotFound {
			fatalLoadError(formatEnvFileOpenError(envFile, result))
		}

		return
	}

	defer result.File.Close()
	if err := setFromEnvFile(result.File, &configs{&App, &Project, &Logger}); err != nil {
		fatalLoadError(err)
	}
}

func openEnvFile(envFile string) envFileOpenResult {
	result := envFileOpenResult{NotFound: true}

	for _, path := range envFileCandidates(envFile) {
		file, err := os.Open(path)
		if err == nil {
			result.File = file
			result.Path = path
			result.NotFound = false
			return result
		}

		result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
		if !errors.Is(err, os.ErrNotExist) {
			result.NotFound = false
		}
	}

	return result
}

func formatEnvFileOpenError(envFile string, result envFileOpenResult) error {
	if len(result.Errors) == 0 {
		return fmt.Errorf(".env: 找不到 %s，且沒有可檢查的候選路徑", envFile)
	}

	messages := make([]string, 0, len(result.Errors))
	for _, err := range result.Errors {
		messages = append(messages, err.Error())
	}

	if result.NotFound {
		return fmt.Errorf(".env: 找不到 %s，已檢查候選路徑: %s", envFile, strings.Join(messages, "; "))
	}

	return fmt.Errorf(".env: 無法開啟 %s，候選路徑錯誤: %s", envFile, strings.Join(messages, "; "))
}

func envFileCandidates(envFile string) []string {
	candidates := []string{envFile}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		recordLoadError(fmt.Errorf(".env: 取得工作目錄失敗: %w", cwdErr))
		return candidates
	}

	if appImagePath := strings.TrimSpace(os.Getenv("APPIMAGE")); appImagePath != "" {
		appImageDir := filepath.Dir(appImagePath)
		if shouldFallbackToExecutableEnv(cwd, appImageDir) {
			candidates = append(candidates, filepath.Join(appImageDir, envFile))
		}
	}

	executablePath, executableErr := os.Executable()
	if executableErr != nil {
		recordLoadError(fmt.Errorf(".env: 取得執行檔路徑失敗: %w", executableErr))
		return candidates
	}

	executableDir := filepath.Dir(executablePath)
	// 處理 Windows 捷徑路徑問題，以及 AppImage 掛載目錄內的設定檔。
	if shouldFallbackToExecutableEnv(cwd, executableDir) {
		candidates = append(candidates, filepath.Join(executableDir, envFile))
	}

	return candidates
}

func shouldFallbackToExecutableEnv(cwd string, executableDir string) bool {
	// air 等開發工具通常會在專案工作目錄底下執行暫存二進制檔，例如 ./bin/app.exe。
	// 這種情況下如果 cwd/.env 不存在，不應讀取可能過期的 bin/.env。
	within, err := isPathWithinOrEqual(cwd, executableDir)
	if err != nil {
		recordLoadError(fmt.Errorf(".env: 判斷 executable .env fallback 路徑失敗: %w", err))
		return false
	}

	return !within
}

func isPathWithinOrEqual(parent string, child string) (bool, error) {
	parentPath, err := normalizeComparablePath(parent)
	if err != nil {
		return false, fmt.Errorf("parent path %q: %w", parent, err)
	}

	childPath, err := normalizeComparablePath(child)
	if err != nil {
		return false, fmt.Errorf("child path %q: %w", child, err)
	}

	if parentPath == childPath {
		return true, nil
	}

	if runtime.GOOS == "windows" && strings.EqualFold(parentPath, childPath) {
		return true, nil
	}

	rel, err := filepath.Rel(parentPath, childPath)
	if err != nil {
		return false, err
	}

	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && !filepath.IsAbs(rel), nil
}

func normalizeComparablePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	cleaned := filepath.Clean(abs)
	evaluated, err := filepath.EvalSymlinks(cleaned)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return normalizePathCase(cleaned), nil
		}

		return "", err
	}

	return normalizePathCase(filepath.Clean(evaluated)), nil
}

func normalizePathCase(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ToLower(path)
	}

	return path
}

func set(structure interface{}) {
	if err := env.Parse(structure); err != nil {
		recordLoadError(err)
	}
}

func setFromEnvFile(file *os.File, structure interface{}) error {
	vars, err := readEnvFile(file)
	if err != nil {
		recordLoadError(err)
		return err
	}

	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			loadErr := fmt.Errorf(".env: cannot set %s; err: %w", key, err)
			recordLoadError(loadErr)
			return loadErr
		}
	}

	set(structure)
	return nil
}

func readEnvFile(file *os.File) (map[string]string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(".env: read %s; err: %w", file.Name(), err)
	}

	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	vars, err := godotenv.UnmarshalBytes(content)
	if err != nil {
		return nil, fmt.Errorf(".env: parse %s; err: %w", file.Name(), err)
	}

	return vars, nil
}

func recordLoadError(err error) {
	LoadErrors = append(LoadErrors, err)
	log.Printf("%+v", err)
}

func fatalLoadError(err error) {
	recordLoadError(err)
	log.Fatalf("config 載入失敗: %v", err)
}
