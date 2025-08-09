package project

import (
	"fmt"
	"gbase/src/core/cmd"
	"gbase/src/core/logger"
	core "gbase/src/core/pm2"
)

func NewPm2(config ProjectConfig, path string) *Project {
	p := &Project{
		Config: config,
		Path:   path,
	}
	p.Install = func() error {
		return cmd.RunInNewConsole([]string{p.Config.Install}, p.Path)
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
		online, err := core.IsServiceOnline(p.Config.Key)
		if err != nil {
			logger.Log.Printf("檢查 pm2 服務 %s 是否啟用失敗：%v", p.Config.Title, err)
		}
		p.Running = online
	}

	return p
}
