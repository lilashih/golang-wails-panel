//go:build !windows && !darwin

package explore

import "os/exec"

func openDirCommand(path string) *exec.Cmd {
	return exec.Command("xdg-open", path)
}
