package project

import (
	"encoding/csv"
	"fmt"
	"gbase/src/core/cmd"
	"io"
	"strings"
)

func NewExe(config ProjectConfig, path string) *Project {
	p := &Project{
		Config: config,
		Path:   path,
	}
	p.Install = func() error {
		return nil
	}
	p.Start = func() error {
		err := cmd.Run("cmd", []string{p.Config.Start}, p.Path)
		waitForRunningState(p, true)

		if !p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}
	p.Stop = func() error {
		err := cmd.Run("cmd", []string{p.Config.Stop}, p.Path)
		waitForRunningState(p, false)

		if p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}

	p.CheckRunning = func() {
		output, err := cmd.RunGetOutput("cmd", []string{"tasklist", "/FO", "CSV", "/NH"}, "")
		if err != nil {
			p.Running = false
			return
		}

		p.Running = isExeRunning(output, p.Config.Key)
	}

	return p
}

func isExeRunning(taskListOutput string, key string) bool {
	candidates := buildExeCandidates(key)
	reader := csv.NewReader(strings.NewReader(taskListOutput))

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
