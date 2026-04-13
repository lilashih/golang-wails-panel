package helper

import (
	"os"
	"path/filepath"
)

const appImageEnvKey = "APPIMAGE"

func GetAppImagePath() string {
	return os.Getenv(appImageEnvKey)
}

func IsAppImage() bool {
	return GetAppImagePath() != ""
}

func getAppImageBasePath() (string, error) {
	appImagePath := GetAppImagePath()
	if appImagePath == "" {
		return "", os.ErrNotExist
	}

	return filepath.Dir(appImagePath), nil
}

// GetRuntimeBasePath 會在 release 模式優先回傳執行檔所在目錄，
// 避免桌面環境雙擊啟動時因工作目錄不同導致相對路徑失效。
func GetRuntimeBasePath() (string, error) {
	if IsRelease() {
		if appImageBasePath, err := getAppImageBasePath(); err == nil {
			return appImageBasePath, nil
		}

		executablePath, err := os.Executable()
		if err == nil {
			return filepath.Dir(executablePath), nil
		}
	}

	return os.Getwd()
}

// GetWritableRuntimeBasePath 會優先回傳可寫入的執行基準路徑。
// AppImage 掛載點通常為唯讀，因此會先改用 .AppImage 所在目錄；
// 若該目錄也不可寫，則退回使用者設定目錄。
func GetWritableRuntimeBasePath() (string, error) {
	basePath, err := GetRuntimeBasePath()
	if err == nil && isWritableDir(basePath) {
		return basePath, nil
	}

	userConfigDir, configErr := os.UserConfigDir()
	if configErr == nil {
		return filepath.Join(userConfigDir, "golang-wails-panel"), nil
	}

	if err != nil {
		return "", err
	}

	return basePath, nil
}

func isWritableDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	testFile, err := os.CreateTemp(path, ".write-test-*")
	if err != nil {
		return false
	}

	testPath := testFile.Name()
	_ = testFile.Close()
	_ = os.Remove(testPath)

	return true
}
