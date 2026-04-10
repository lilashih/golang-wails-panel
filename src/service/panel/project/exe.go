package project

import (
	"encoding/csv"
	"fmt"
	"gbase/src/core/cmd/shell"
	"io"
	"runtime"
	"strings"
)

func NewExe(config ProjectConfig, path string) *Project {
	p := &Project{
		Config: config,
		Path:   path,
	}
	p.Install = func() error {
		if strings.TrimSpace(p.Config.Install) == "" {
			return nil
		}

		return shell.RunInNewConsole(p.Config.Install, p.Path)
	}
	p.Start = func() error {
		err := shell.Run(p.Config.Start, p.Path)
		waitForRunningState(p, true)

		if !p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}
	p.Stop = func() error {
		err := shell.Run(p.Config.Stop, p.Path)
		waitForRunningState(p, false)

		if p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}

	p.CheckRunning = func() {
		p.Running = isExeRunning(p.Config.Key)
	}

	return p
}

func isExeRunning(key string) bool {
	switch runtime.GOOS {
	case "windows":
		return isExeRunningOnWindows(key)
	default:
		return isExeRunningOnUnix(key)
	}
}

func isExeRunningOnWindows(key string) bool {
	output, err := shell.RunGetOutput("tasklist /FO CSV /NH", "")
	if err != nil {
		return false
	}

	return isExeRunningOnWindowsOutput(output, key)
}

func isExeRunningOnWindowsOutput(output string, key string) bool {
	candidates := buildExeCandidates(key)
	reader := csv.NewReader(strings.NewReader(output))

	for {
		record, err := reader.Read()
		if err == io.EOF {
			return false
		}

		if err != nil || len(record) == 0 {
			continue
		}

		imageName := strings.TrimSpace(record[0])
		for _, candidate := range candidates {
			if strings.EqualFold(imageName, candidate) {
				return true
			}
		}
	}
}

func buildExeCandidates(key string) []string {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}

	candidates := []string{key}
	if len(key) < 4 || !strings.EqualFold(key[len(key)-4:], ".exe") {
		candidates = append(candidates, key+".exe")
	}

	return uniqueStrings(candidates)
}

func isExeRunningOnUnix(key string) bool {
	for _, candidate := range buildExeCandidates(key) {
		commandLine := fmt.Sprintf("pgrep -f %q", candidate)
		if _, err := shell.RunGetOutput(commandLine, ""); err == nil {
			return true
		}
	}

	return false
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
