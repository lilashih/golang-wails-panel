package shell

import (
	"os"
	"strings"
)

func buildCommandEnv() []string {
	env := setEnvValue(os.Environ(), "NO_COLOR", "1")
	return applyCommandEnvDefaults(env)
}

func getEnvValue(env []string, key string) string {
	for _, item := range env {
		name, value, ok := strings.Cut(item, "=")
		if ok && envKeyEqual(name, key) {
			return value
		}
	}

	return ""
}

func setEnvValue(env []string, key string, value string) []string {
	for i, item := range env {
		name, _, ok := strings.Cut(item, "=")
		if ok && envKeyEqual(name, key) {
			env[i] = name + "=" + value
			return env
		}
	}

	return append(env, key+"="+value)
}

func splitEnvPath(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return strings.Split(value, string(os.PathListSeparator))
}

func prependMissingPath(paths []string, next string) []string {
	next = strings.TrimSpace(next)
	if next == "" {
		return paths
	}

	for _, path := range paths {
		if path == next {
			return paths
		}
	}

	return append([]string{next}, paths...)
}
