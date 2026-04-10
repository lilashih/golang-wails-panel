//go:build !windows

package shell

import "os/exec"

func configureBackgroundCommand(cmd *exec.Cmd) {}

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
