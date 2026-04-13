//go:build windows

package exe

import (
	"encoding/csv"
	"gbase/src/core/cmd/shell"
	"io"
	"strings"
)

func isRunning(key string) bool {
	candidates := buildCandidates(key)
	if len(candidates) == 0 {
		return false
	}

	output, err := shell.RunGetOutput("tasklist /FO CSV /NH", "")
	if err != nil {
		return false
	}

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

func buildCandidates(key string) []string {
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
