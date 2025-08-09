package project

import (
	"fmt"
	"gbase/src/core/cmd"
	"strings"
)

func NewExe(config ProjectConfig, path string) *Project {
	p := &Project{
		Config: config,
		Path:   path,
	}
	p.Install = func() error {
		return nil
	}
	p.Start = func() error {
		err := cmd.Run("cmd", []string{p.Config.Start}, p.Path)
		p.CheckRunning()

		if !p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}
	p.Stop = func() error {
		err := cmd.Run("cmd", []string{p.Config.Stop}, p.Path)
		p.CheckRunning()

		if p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}

	p.CheckRunning = func() {
		output, err := cmd.RunGetOutput("cmd", []string{"tasklist"}, "")
		if err != nil {
			p.Running = false
		} else {
			p.Running = strings.Contains(string(output), p.Config.Key)
		}
	}

	return p
}
