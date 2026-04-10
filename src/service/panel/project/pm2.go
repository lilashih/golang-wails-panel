package project

import (
	"fmt"
	"gbase/src/core/cmd/shell"
	"gbase/src/core/logger"
	core "gbase/src/core/pm2"
	"strings"
)

func NewPm2(config ProjectConfig, path string) *Project {
	p := &Project{
		Config: config,
		Path:   path,
	}
	p.Install = func() error {
		if strings.TrimSpace(p.Config.Install) == "" {
			return nil
		}

		return shell.RunInNewConsole(p.Config.Install, p.Path)
	}
	p.Start = func() error {
		err := shell.Run(p.Config.Start, p.Path)
		waitForRunningState(p, true)

		if !p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}
	p.Stop = func() error {
		err := shell.Run(p.Config.Stop, p.Path)
		waitForRunningState(p, false)

		if p.Running {
			err = fmt.Errorf("未知錯誤 (請檢查.env檔或相關設定)")
		}

		return err
	}
	p.CheckRunning = func() {
		online, err := core.IsServiceOnline(p.Config.Key, p.Path)
		if err != nil {
			logger.Log.Printf("檢查 pm2 服務 %s 是否啟用失敗：%v", p.Config.Title, err)
		}
		p.Running = online
	}

	return p
}
