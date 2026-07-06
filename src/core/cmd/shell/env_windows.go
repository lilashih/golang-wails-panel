//go:build windows

package shell

import "strings"

func applyCommandEnvDefaults(env []string) []string {
	return env
}

// Windows 環境變數名稱大小寫不敏感。
// 實務上 PATH 常會儲存成 "Path"，因此在 Windows 上不能用大小寫敏感的字串比對。
func envKeyEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
