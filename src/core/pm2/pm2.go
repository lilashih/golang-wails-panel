package pm2

import (
	"encoding/json"
	"fmt"
	"gbase/src/core/cmd"
	"strings"
)

type pm2Process struct {
	Name   string `json:"name"`
	PM2Env struct {
		Status string `json:"status"`
	} `json:"pm2_env"`
}

func IsServiceOnline(service string) (bool, error) {
	data, err := cmd.RunGetOutput("cmd", []string{"pm2", "jlist"}, "")
	if err != nil {
		return false, fmt.Errorf("執行 pm2 失敗: %w", err)
	}

	data = cleanJSONOutput(data) // 只取 JSON 區塊

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

func cleanJSONOutput(data string) string {
	start := strings.Index(data, "[")
	end := strings.LastIndex(data, "]")
	if start >= 0 && end > start {
		return data[start : end+1]
	}
	return data
}
