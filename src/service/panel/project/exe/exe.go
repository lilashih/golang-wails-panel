// Package exe 實作 project.json 中 type="exe" 的專案 runner。
//
// exe runner 適合用一般 shell 指令啟動、停止，並透過作業系統 process 清單判斷執行狀態的專案。
// 它內嵌 runner.Command，因此 Install、Start、Stop 都沿用共用的 command 行為：
//   - Install 會在新 console 視窗執行 config.install。
//   - Start 會執行 config.start，等待該指令結束後，再檢查目標 process 是否已經執行。
//   - Stop 會執行 config.stop，等待該指令結束後，再檢查目標 process 是否已經停止。
//
// CheckRunning 是 exe runner 自己負責的部分。它會呼叫 isRunning(config.key)，
// 並由各 OS 平台的實作檢查目前的 process 清單。
// Linux 目前主要比對執行檔名稱；Windows 則檢查 task list。
// 因此 config.key 應該填入可辨識的 process 名稱，不應只是任意顯示用名稱。
//
// exe runner 可以用來啟動常駐的開發指令，但 config.start 必須自己讓指令背景化並結束， 否則 Panel 會一直等待 Start 返回。
// 例如 Windows 可搭配 "start /B air"，Linux 可搭配 "nohup air > air.log 2>&1 &"。
package exe

import (
	"gbase/src/service/panel/project/runner"
	"strings"
)

type Runner struct {
	*runner.Command
}

func New(config runner.Config, path string) *Runner {
	return &Runner{
		Command: runner.NewCommand(config, path),
	}
}

func (r *Runner) CheckRunning() bool {
	return isRunning(r.Config.Key)
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))

	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		result = append(result, value)
	}

	return result
}
