//go:build windows

package explore

import "os/exec"

func openDirCommand(path string) *exec.Cmd {
	return exec.Command("explorer", path)
}
