package panel

import (
	"context"
	"fmt"
	"sync"

	"gbase/src/core/cmd/explore"
	"gbase/src/core/logger"
	"gbase/src/service/panel/project"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	ctx        context.Context
	mu         sync.RWMutex
	projects   []*project.Project
	projectOps map[string]*sync.Mutex
}

func NewService() *Service {
	return &Service{
		projectOps: make(map[string]*sync.Mutex),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *Service) ListProjects() []*project.Project {
	projects, err := project.NewProjects()
	if err != nil {
		logger.Error("取得應用程式列表失敗：%v", err)
	}
	projects = uniqueProjectsByKey(projects)

	s.mu.Lock()
	s.projects = projects
	s.ensureProjectLocksLocked(projects)
	result := append([]*project.Project(nil), s.projects...)
	s.mu.Unlock()

	return result
}

func (s *Service) OpenProject(key string) error {
	p, err := s.projectByKey(key)
	if err != nil {
		logger.Error("開啟資料夾失敗：%v", err)
		return err
	}

	if err := explore.OpenDir(p.Path); err != nil {
		logger.Error("開啟資料夾失敗：%v", err)
		return err
	}

	return nil
}

func (s *Service) StartProject(key string) (bool, error) {
	unlock := s.lockProjectOperation(key)
	defer unlock()

	p, err := s.projectByKey(key)
	if err != nil {
		logger.Error("啟動專案失敗：%v", err)
		return false, err
	}

	logger.Info("%s 啟動中...", p.Config.Title)
	err = p.Start()
	if err != nil {
		logger.Error("%s 啟動失敗：%v", p.Config.Title, err)
	} else {
		logger.Info("%s 啟動成功", p.Config.Title)
	}

	return p.Running, err
}

func (s *Service) StopProject(key string) (bool, error) {
	unlock := s.lockProjectOperation(key)
	defer unlock()

	p, err := s.projectByKey(key)
	if err != nil {
		logger.Error("停止專案失敗：%v", err)
		return false, err
	}

	logger.Info("%s 停止中...", p.Config.Title)
	err = p.Stop()
	p.CheckRunning()
	if err != nil {
		if !p.Running {
			logger.Info("%s 已停止", p.Config.Title)
			return p.Running, nil
		}

		logger.Error("%s 停止失敗：%v", p.Config.Title, err)
		return p.Running, err
	}

	logger.Info("%s 已停止", p.Config.Title)
	return p.Running, nil
}

func (s *Service) InstallProject(key string) error {
	unlock := s.lockProjectOperation(key)
	defer unlock()

	p, err := s.projectByKey(key)
	if err != nil {
		logger.Error("安裝專案失敗：%v", err)
		return err
	}

	if p.Running {
		err := fmt.Errorf("%s 正在執行中，請先停止再重新安裝", p.Config.Title)
		logger.Error("%v", err)
		return err
	}

	logger.Info("%s 安裝中...", p.Config.Title)
	err = p.Install()
	if err != nil {
		logger.Error("%s 安裝失敗：%v", p.Config.Title, err)
	} else {
		logger.Info("%s 已安裝", p.Config.Title)
	}

	return err
}

func (s *Service) projectByKey(key string) (*project.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.projects {
		if p.Config.Key == key {
			return p, nil
		}
	}

	return nil, fmt.Errorf("找不到專案：%s", key)
}

func uniqueProjectsByKey(projects []*project.Project) []*project.Project {
	seen := make(map[string]struct{}, len(projects))
	unique := make([]*project.Project, 0, len(projects))

	for _, p := range projects {
		if _, ok := seen[p.Config.Key]; ok {
			logger.Error("專案 key 重複，已跳過：%s (%s)", p.Config.Key, p.Path)
			continue
		}

		seen[p.Config.Key] = struct{}{}
		unique = append(unique, p)
	}

	return unique
}

func (s *Service) ensureProjectLocksLocked(projects []*project.Project) {
	if s.projectOps == nil {
		s.projectOps = make(map[string]*sync.Mutex)
	}

	for _, p := range projects {
		if _, ok := s.projectOps[p.Config.Key]; !ok {
			s.projectOps[p.Config.Key] = &sync.Mutex{}
		}
	}
}

func (s *Service) lockProjectOperation(key string) func() {
	s.mu.Lock()
	if s.projectOps == nil {
		s.projectOps = make(map[string]*sync.Mutex)
	}
	lock, ok := s.projectOps[key]
	if !ok {
		lock = &sync.Mutex{}
		s.projectOps[key] = lock
	}
	s.mu.Unlock()

	lock.Lock()
	return lock.Unlock
}
