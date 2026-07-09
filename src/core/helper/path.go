package helper

import (
	"gbase/src/core/config"
	"os"
	"path/filepath"
	"strings"
)

// GetRuntimeBasePath 回傳程式執行時應使用的基準路徑，或執行檔所在目錄
// 這邊只能透過手動設定 APP_MODE 判斷，因為 wails, air 導致無法判斷到正確路徑
//
// release 版 (exe)：回傳執行檔所在目錄，避免使用者從其他工作目錄啟動程式時，導致設定檔、assets、projects 等相對路徑讀取錯誤。
// 非 release 版 (go run, air)：回傳目前工作目錄，方便使用 go run 或 IDE 執行時，能以專案根目錄作為基準存取檔案。
func GetRuntimeBasePath() (string, error) {
	if IsRelease() {
		if appImagePath := strings.TrimSpace(os.Getenv("APPIMAGE")); appImagePath != "" {
			return filepath.Dir(appImagePath), nil
		}

		executablePath, err := os.Executable()
		if err == nil {
			return filepath.Dir(executablePath), nil
		}
	}

	return os.Getwd()
}

// GetWritableRuntimeBasePath 會優先回傳可寫入的執行基準路徑。
// 若執行基準路徑不可寫，則退回使用者設定目錄。
func GetWritableRuntimeBasePath() (string, error) {
	basePath, err := GetRuntimeBasePath()
	if err == nil && isWritableDir(basePath) {
		return basePath, nil
	}

	userConfigDir, configErr := os.UserConfigDir()
	if configErr == nil {
		return filepath.Join(userConfigDir, appConfigDirName()), nil
	}

	if err != nil {
		return "", err
	}

	return basePath, nil
}

func appConfigDirName() string {
	name := strings.TrimSpace(config.App.Id)
	if name != "" {
		return name
	}

	name = strings.TrimSpace(config.App.Name)
	if name != "" {
		return name
	}

	return "gbase"
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
