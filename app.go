package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type App struct {
	Name    string   `json:"name"`
	Icon    string   `json:"icon"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Refresh int      `json:"refresh,omitempty"`
}

type AppService struct {
	mu       sync.RWMutex
	apps     []App
	filePath string
}

func NewAppService() *AppService {
	configDir, _ := os.UserConfigDir()
	dir := filepath.Join(configDir, "tray")
	os.MkdirAll(dir, 0755)

	svc := &AppService{
		filePath: filepath.Join(dir, "apps.json"),
	}
	svc.load()
	return svc
}

func (s *AppService) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		s.apps = []App{}
		return
	}
	json.Unmarshal(data, &s.apps)
	if s.apps == nil {
		s.apps = []App{}
	}
}

func (s *AppService) save() error {
	data, err := json.MarshalIndent(s.apps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *AppService) GetApps() []App {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apps
}

func (s *AppService) SaveApp(app App) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apps = append(s.apps, app)
	return s.save()
}

func (s *AppService) DeleteApp(index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.apps) {
		return nil
	}
	s.apps = append(s.apps[:index], s.apps[index+1:]...)
	return s.save()
}
