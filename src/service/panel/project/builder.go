package project

import (
	"gbase/src/service/panel/project/exe"
	"gbase/src/service/panel/project/pm2"
	"gbase/src/service/panel/project/runner"
	"time"
)

type serviceRunner interface {
	Install() error
	Start(func(expected bool) bool) error
	Stop(func(expected bool) bool) error
	CheckRunning() bool
}

func buildProject(config ProjectConfig, path string) *Project {
	runner := newServiceRunner(config, path)
	if runner == nil {
		return nil
	}

	p := &Project{
		Config: config,
		Path:   path,
	}

	p.Install = runner.Install
	p.Start = func() error {
		return runner.Start(func(expected bool) bool {
			waitForRunningState(p, expected)
			return p.Running
		})
	}
	p.Stop = func() error {
		return runner.Stop(func(expected bool) bool {
			waitForRunningState(p, expected)
			return p.Running
		})
	}
	p.CheckRunning = func() {
		p.Running = runner.CheckRunning()
	}

	return p
}

func newServiceRunner(config ProjectConfig, path string) serviceRunner {
	runnerConfig := runner.Config{
		Title:   config.Title,
		Key:     config.Key,
		Install: config.Install,
		Start:   config.Start,
		Stop:    config.Stop,
	}

	switch config.Type {
	case "pm2":
		return pm2.New(runnerConfig, path)
	case "exe":
		return exe.New(runnerConfig, path)
	default:
		return nil
	}
}

// 啟動或停止時不馬上檢查狀態，等待3秒後再查詢
func waitForRunningState(p *Project, expected bool) {
	deadline := time.Now().Add(3 * time.Second)

	for {
		p.CheckRunning()
		if p.Running == expected || time.Now().After(deadline) {
			return
		}

		time.Sleep(250 * time.Millisecond)
	}
}
