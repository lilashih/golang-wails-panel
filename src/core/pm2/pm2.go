package pm2

import (
	"encoding/json"
	"fmt"
	"gbase/src/core/cmd/shell"
	"strings"
)

type pm2Process struct {
	Name   string `json:"name"`
	PM2Env struct {
		Status string `json:"status"`
	} `json:"pm2_env"`
}

func IsServiceOnline(service string, cmdDir string) (bool, error) {
	commands := []string{
		"pnpm exec pm2 jlist",
		"npm exec -- pm2 jlist",
		"npx pm2 jlist",
		"pm2 jlist",
	}

	var errors []string

	for _, commandLine := range commands {
		data, err := shell.RunGetOutput(commandLine, cmdDir)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", commandLine, err))
			continue
		}

		online, err := parseServiceStatus(data, service)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", commandLine, err))
			continue
		}

		return online, nil
	}

	return false, fmt.Errorf("執行 pm2 失敗: %s", strings.Join(errors, " | "))
}

func parseServiceStatus(data string, service string) (bool, error) {
	data = cleanJSONOutput(data)

	var processes []pm2Process
	if err := json.Unmarshal([]byte(data), &processes); err != nil {
		return false, fmt.Errorf("解析 JSON 失敗: %w", err)
	}

	// 搜尋指定服務
	for _, proc := range processes {
		if proc.Name == service {
			return proc.PM2Env.Status == "online", nil
		}
	}

	// 未找到服務 -> 未啟用
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
