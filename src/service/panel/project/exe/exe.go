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
