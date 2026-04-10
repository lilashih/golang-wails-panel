package shell

import (
	"bufio"
	"bytes"
	"gbase/src/core/logger"
	"io"
	"os"
	"strings"
)

// 執行指令，並回傳標準輸出結果
func RunGetOutput(commandLine string, cmdDir string) (string, error) {
	cmd := buildShellCommand(commandLine, cmdDir)
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	configureBackgroundCommand(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // 也可以合併錯誤輸出

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// 背景執行執行專案指令，將命令提示自元的訊息存到log
func Run(commandLine string, cmdDir string) error {
	cmd := buildShellCommand(commandLine, cmdDir)
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	configureBackgroundCommand(cmd)

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
func RunInNewConsole(commandLine string, dir string) error {
	cmd := buildConsoleCommand(commandLine, dir)
	return cmd.Run() // 會同步等待；exit code ≠ 0 時回傳 *exec.ExitError
}
