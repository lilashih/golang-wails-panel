package panel

import (
	"context"
	"gbase/src/core/cmd"
	"gbase/src/service/panel/project"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PanelService struct {
	ctx      context.Context
	Projects []*project.Project
}

func NewPanelService() *PanelService {
	return &PanelService{}
}

func (s *PanelService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// ListProjects 回傳子目錄名稱切片
func (s *PanelService) ListProjects() []*project.Project {
	projects, err := project.NewProjects()
	s.Projects = projects

	if err != nil {
		runtime.LogErrorf(s.ctx, "取得專案列表失敗：%v", err)
	}

	return s.Projects
}

// OpenProject 範例：在作業系統上「開啟」該子目錄（自己決定要做什麼）
func (s *PanelService) OpenProject(index int) {
	p := s.Projects[index]

	if err := cmd.OpenDir(p.Path); err != nil {
		runtime.LogErrorf(s.ctx, "開啟資料夾失敗：%v", err)
	}
}

func (s *PanelService) StartProject(index int) bool {
	p := s.Projects[index]

	runtime.LogInfof(s.ctx, "%s 啟動中...", p.Config.Title)
	err := p.Start()

	if err != nil {
		runtime.LogErrorf(s.ctx, "%s 啟動失敗：%v", p.Config.Title, err)
	} else {
		runtime.LogInfof(s.ctx, "%s 啟動成功", p.Config.Title)
	}

	return p.Running
}

func (s *PanelService) StopProject(index int) bool {
	p := s.Projects[index]

	runtime.LogInfof(s.ctx, "%s 停止中...", p.Config.Title)
	err := p.Stop()

	if err != nil {
		runtime.LogErrorf(s.ctx, "%s 停止失敗：%v", p.Config.Title, err)
	} else {
		runtime.LogInfof(s.ctx, "%s 已停止", p.Config.Title)
	}

	return p.Running
}

func (s *PanelService) InstallProject(index int) {
	p := s.Projects[index]

	if p.Running {
		runtime.LogErrorf(s.ctx, "%s 正在執行中，請先停止再重新安裝", p.Config.Title)
		return
	}

	runtime.LogInfof(s.ctx, "%s 安裝中...", p.Config.Title)
	err := p.Install()

	if err != nil {
		runtime.LogErrorf(s.ctx, "%s 停止失敗：%v", p.Config.Title, err)
	} else {
		runtime.LogInfof(s.ctx, "%s 已安裝", p.Config.Title)
	}
}
