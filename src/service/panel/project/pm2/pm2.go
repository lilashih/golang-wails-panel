// Package pm2 實作 project.json 中 type="pm2" 的專案 runner。
//
// pm2 runner 適合交由 PM2 管理的 Node.js、前端開發伺服器，或其他可透過 pm2 start / stop 控制的常駐服務。
// 它內嵌 runner.Command，因此 Install、Start、Stop 都沿用共用的 command 行為：
//   - Install 會在新 console 視窗執行 config.install。
//   - Start 會執行 config.start，等待該指令結束後，再透過 PM2 檢查服務是否已經 online。
//   - Stop 會執行 config.stop，等待該指令結束後，再透過 PM2 檢查服務是否已經停止。
//
// CheckRunning 是 pm2 runner 自己負責的部分。它會執行 "pm2 jlist"，取得 PM2 process JSON 後，
// 用 config.key 對應 process name，並以 pm2_env.status == "online" 判斷是否執行中。
// 因此 config.key 必須和 PM2 裡的 process name 一致，也就是 pm2 start 時 --name 指定的名稱。
// 若 PM2 不在 PATH，CheckRunning 會回傳 false 並記錄錯誤。
package pm2

import (
	"encoding/json"
	"fmt"
	"gbase/src/core/cmd/shell"
	"gbase/src/core/logger"
	"gbase/src/service/panel/project/runner"
	"strings"
)

type Runner struct {
	*runner.Command
}

type process struct {
	Name   string `json:"name"`
	PM2Env struct {
		Status string `json:"status"`
	} `json:"pm2_env"`
}

func New(config runner.Config, path string) *Runner {
	return &Runner{
		Command: runner.NewCommand(config, path),
	}
}

func (r *Runner) CheckRunning() bool {
	online, err := isServiceOnline(r.Config.Key, r.Path)
	if err != nil {
		logger.Log.Printf("檢查 pm2 服務 %s 是否啟用失敗：%v", r.Config.Title, err)
	}

	return online
}

func isServiceOnline(service string, cmdDir string) (bool, error) {
	const commandLine = "pm2 jlist"

	data, err := shell.RunGetOutput(commandLine, cmdDir)
	if err != nil {
		return false, fmt.Errorf("%s: %v", commandLine, err)
	}

	online, err := parseServiceStatus(data, service)
	if err != nil {
		return false, fmt.Errorf("%s: %v", commandLine, err)
	}

	return online, nil
}

func parseServiceStatus(data string, service string) (bool, error) {
	data = cleanJSONOutput(data)

	var processes []process
	if err := json.Unmarshal([]byte(data), &processes); err != nil {
		return false, fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	for _, proc := range processes {
		if proc.Name == service {
			return proc.PM2Env.Status == "online", nil
		}
	}

	return false, nil
}

// pm2 jlist 的輸出有時會夾帶額外文字或日誌，把其中的 JSON 陣列（[...]）擷取出來
func cleanJSONOutput(data string) string {
	start := strings.Index(data, "[")
	end := strings.LastIndex(data, "]")
	if start >= 0 && end > start {
		return data[start : end+1]
	}

	return data
}
