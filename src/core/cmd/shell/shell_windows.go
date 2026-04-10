//go:build windows

package shell

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func configureBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func buildShellCommand(commandLine string, dir string) *exec.Cmd {
	cmd := exec.Command("cmd", "/C", commandLine)
	cmd.Dir = dir

	return cmd
}

func buildConsoleCommand(commandLine string, dir string) *exec.Cmd {
	/*
	   start "" /WAIT cmd /C <commandLine>
	   - start ""      → 真正開新黑窗
	   - /WAIT         → 等 <commandLine> 結束才結束 start
	   - cmd /C        → 指令跑完自動關視窗
	*/
	cmd := exec.Command("cmd", "/C", "start", "", "/WAIT", "cmd", "/C", commandLine)
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	return cmd
}
