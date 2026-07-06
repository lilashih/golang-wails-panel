package runner

import (
	"fmt"
	"gbase/src/core/cmd/shell"
	"strings"
	"time"
)

var (
	timeout = 30 * time.Second
)

type Command struct {
	Config Config
	Path   string
}

func NewCommand(config Config, path string) *Command {
	return &Command{
		Config: config,
		Path:   path,
	}
}

func (c *Command) Install() error {
	if strings.TrimSpace(c.Config.Install) == "" {
		return nil
	}

	return shell.RunInNewConsole(c.Config.Install, c.Path)
}

func (c *Command) Start(waitForRunningState func(expected bool) bool) error {
	if err := shell.RunWithTimeout(c.Config.Start, c.Path, timeout); err != nil {
		return err
	}

	if !waitForRunningState(true) {
		return newStateError()
	}

	return nil
}

func (c *Command) Stop(waitForRunningState func(expected bool) bool) error {
	if err := shell.RunWithTimeout(c.Config.Stop, c.Path, timeout); err != nil {
		return err
	}

	if !waitForRunningState(false) {
		return newStateError()
	}

	return nil
}

func newStateError() error {
	return fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
}
