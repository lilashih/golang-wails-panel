//go:build !windows

package shell

import (
	"os"
	"path/filepath"
	"strings"
)

func applyCommandEnvDefaults(env []string) []string {
	paths := splitEnvPath(getEnvValue(env, "PATH"))
	preferred := []string{}

	preferred = append(preferred, splitEnvPath(getEnvValue(env, "COMMAND_EXTRA_PATH"))...)

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		preferred = append(preferred,
			filepath.Join(home, ".local", "bin"),
			filepath.Join(home, "go", "bin"),
			filepath.Join(home, ".pnpm-live"),
			filepath.Join(home, ".local", "share", "pnpm"),
			filepath.Join(home, ".npm-global", "bin"),
		)

		// 嘗試加入 NVM 路徑 (自動尋找最新的 Node 版本 bin 目錄)
		nvmDir := filepath.Join(home, ".nvm", "versions", "node")
		if entries, err := os.ReadDir(nvmDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					preferred = append(preferred, filepath.Join(nvmDir, entry.Name(), "bin"))
				}
			}
		}
	}

	preferred = append(preferred,
		"/opt/homebrew/bin",
		"/opt/homebrew/sbin",
		"/usr/local/go/bin",
		"/usr/local/bin",
		"/usr/local/sbin",
		"/usr/bin",
		"/bin",
		"/usr/sbin",
		"/sbin",
	)

	for i := len(preferred) - 1; i >= 0; i-- {
		paths = prependMissingPath(paths, preferred[i])
	}

	return setEnvValue(env, "PATH", strings.Join(paths, string(os.PathListSeparator)))
}

func envKeyEqual(a, b string) bool {
	return a == b
}
