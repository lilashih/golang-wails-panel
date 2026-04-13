//go:build !windows

package exe

import (
	"gbase/src/core/cmd/shell"
	"path"
	"strings"
)

func isRunning(key string) bool {
	candidates := buildCandidates(key)
	if len(candidates) == 0 {
		return false
	}

	output, err := shell.RunGetOutput("ps -eo args=", "")
	if err != nil {
		return false
	}

	for _, line := range strings.Split(output, "\n") {
		processName := extractUnixProcessName(line)
		if processName == "" {
			continue
		}

		for _, candidate := range candidates {
			if processName == candidate {
				return true
			}
		}
	}

	return false
}

func extractUnixProcessName(line string) string {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return ""
	}

	executable := strings.ReplaceAll(fields[0], "\\", "/")
	name := strings.TrimSpace(path.Base(executable))
	if name == "." {
		return ""
	}

	return name
}

func buildCandidates(key string) []string {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}

	normalizedKey := strings.ReplaceAll(key, "\\", "/")
	baseName := strings.TrimSpace(path.Base(normalizedKey))
	if baseName == "" || baseName == "." {
		return nil
	}

	candidates := []string{baseName}
	if strings.HasSuffix(strings.ToLower(baseName), ".exe") {
		candidates = append(candidates, baseName[:len(baseName)-4])
	}

	return uniqueStrings(candidates)
}
