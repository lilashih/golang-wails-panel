package explore

import (
	"fmt"
	"os"
)

func OpenDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("開啟資料夾失敗，路徑不存在：%s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("開啟資料夾失敗，路徑不是資料夾：%s", path)
	}

	cmd := openDirCommand(path)
	return cmd.Start()
}
