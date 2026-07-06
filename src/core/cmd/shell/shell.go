package shell

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"gbase/src/core/logger"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const maxLogLineSize = 1024 * 1024

// 執行指令，並回傳標準輸出結果
func RunGetOutput(commandLine string, cmdDir string) (string, error) {
	cmd := buildShellCommand(commandLine, cmdDir)
	cmd.Env = buildCommandEnv()
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
	cmd.Env = buildCommandEnv()
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

	logWG := streamCommandOutput(stdout, stderr)
	err = cmd.Wait()
	logWG.Wait()

	return err
}

// 執行指令，並回加上timeout
// 使用方式: err := shell.RunWithTimeout(c.Config.Start, c.Path, 5*time.Second)
func RunWithTimeout(commandLine string, cmdDir string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := buildShellCommand(commandLine, cmdDir)
	cmd.Env = buildCommandEnv()
	configureBackgroundCommand(cmd)
	configureTimeoutCommand(cmd)

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

	logWG := streamCommandOutput(stdout, stderr)

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	select {
	case err := <-waitDone:
		logWG.Wait()
		if err != nil {
			logger.Error("%v", err)
		}
		return err
	case <-ctx.Done():
		if killErr := terminateCommandTree(cmd); killErr != nil {
			logger.Error("終止逾時指令失敗：%v", killErr)
		}
		<-waitDone
		logWG.Wait()
		return fmt.Errorf("指令執行逾時，常駐指令需自行背景化並結束，如 start /B ...、nohup ...：%s", commandLine)
	}
}

func streamCommandOutput(stdout io.ReadCloser, stderr io.ReadCloser) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(2)
	go streamLogger(stdout, &wg)
	go streamLogger(stderr, &wg)
	return &wg
}

// 即時 log 輸出到 log
func streamLogger(pipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLogLineSize)
	for scanner.Scan() {
		logger.Log.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil && !isClosedPipeError(err) {
		logger.Error("讀取指令輸出失敗：%v", err)
	}
}

// cmd.Wait() 結束時 os/exec 會關閉 stdout/stderr pipe，log reader 的 Scanner 偶爾會拿到：  read |0: file already closed
// 這不是指令輸出讀取真正失敗，也不是業務錯誤，所以不應該打成 error。
func isClosedPipeError(err error) bool {
	return errors.Is(err, os.ErrClosed) || strings.Contains(err.Error(), "file already closed")
}

// 在新視窗中執行指令(會彈出小黑窗)，並等待結束
func RunInNewConsole(commandLine string, dir string) error {
	cmd := buildConsoleCommand(commandLine, dir)
	cmd.Env = buildCommandEnv()
	return cmd.Run() // 會同步等待；exit code ≠ 0 時回傳 *exec.ExitError
}
