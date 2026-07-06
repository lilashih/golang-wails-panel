//go:build !windows

package shell

import (
	"os/exec"
	"syscall"
)

func configureBackgroundCommand(cmd *exec.Cmd) {}

func configureTimeoutCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func terminateCommandTree(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

func buildShellCommand(commandLine string, dir string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", commandLine)
	cmd.Dir = dir

	return cmd
}

func buildConsoleCommand(commandLine string, dir string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", commandLine)
	cmd.Dir = dir

	return cmd
}
