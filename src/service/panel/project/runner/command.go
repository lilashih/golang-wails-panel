package runner

import (
	"fmt"
	"gbase/src/core/cmd/shell"
	"strings"
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
	err := shell.Run(c.Config.Start, c.Path)

	if !waitForRunningState(true) {
		err = newStateError()
	}

	return err
}

func (c *Command) Stop(waitForRunningState func(expected bool) bool) error {
	err := shell.Run(c.Config.Stop, c.Path)

	if waitForRunningState(false) {
		err = newStateError()
	}

	return err
}

func newStateError() error {
	return fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
}
