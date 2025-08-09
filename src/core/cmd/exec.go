package cmd

import (
	"bufio"
	"bytes"
	"gbase/src/core/logger"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// 執行指令，並回傳標準輸出結果
func RunGetOutput(name string, args []string, cmdDir string) (string, error) {
	if name == "cmd" {
		if runtime.GOOS == "windows" {
			args = append([]string{"/C"}, args...)
		} else {
			name = "sh"
			args = append([]string{"-c"}, args...)
		}
	}

	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Dir = cmdDir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // 也可以合併錯誤輸出

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// 背景執行執行專案指令，將命令提示自元的訊息存到log
func Run(name string, args []string, cmdDir string) error {
	if name == "cmd" {
		if runtime.GOOS == "windows" {
			args = append([]string{"/C"}, args...)
		} else {
			name = "sh"
			args = append([]string{"-c"}, args...)
		}
	}

	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Dir = cmdDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go streamLogger(stdout)
	go streamLogger(stderr)

	return cmd.Wait()
}

// 即時 log 輸出到 log
func streamLogger(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		logger.Log.Println(scanner.Text())
	}
}

// 在新視窗中執行指令(會彈出小黑窗)，並等待結束
func RunInNewConsole(args []string, dir string) error {
	commandLine := strings.Join(args, " ")

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
	return cmd.Run() // 會同步等待；exit code ≠ 0 時回傳 *exec.ExitError
}
