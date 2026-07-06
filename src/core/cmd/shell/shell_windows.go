//go:build windows

package shell

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"

	"golang.org/x/sys/windows"
)

func configureBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func configureTimeoutCommand(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_NEW_PROCESS_GROUP
}

func terminateCommandTree(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := kill.CombinedOutput(); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("taskkill failed: %w: %s", err, string(output))
	}

	return nil
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
