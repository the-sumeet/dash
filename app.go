package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type App struct {
	Name    string   `json:"name"`
	Icon    string   `json:"icon"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Refresh int      `json:"refresh,omitempty"`
}

type cachedOutput struct {
	output string
	time   time.Time
}

type AppService struct {
	mu       sync.RWMutex
	apps     []App
	filePath string
	cache    map[int]cachedOutput
}

func NewAppService() *AppService {
	configDir, _ := os.UserConfigDir()
	dir := filepath.Join(configDir, "tray")
	os.MkdirAll(dir, 0755)

	svc := &AppService{
		filePath: filepath.Join(dir, "apps.json"),
		cache:    make(map[int]cachedOutput),
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
	s.cache = make(map[int]cachedOutput)
	return s.save()
}

func (s *AppService) RunApp(index int) (string, error) {
	s.mu.RLock()
	if index < 0 || index >= len(s.apps) {
		s.mu.RUnlock()
		return "", fmt.Errorf("invalid app index: %d", index)
	}
	app := s.apps[index]
	if cached, ok := s.cache[index]; ok {
		if app.Refresh <= 0 || time.Since(cached.time) < time.Duration(app.Refresh)*time.Second {
			s.mu.RUnlock()
			return cached.output, nil
		}
	}
	s.mu.RUnlock()

	cmd := exec.Command(app.Command, app.Args...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return string(output), fmt.Errorf("%s: %w", string(output), err)
	}

	s.mu.Lock()
	s.cache[index] = cachedOutput{output: string(output), time: time.Now()}
	s.mu.Unlock()

	return string(output), nil
}

func (s *AppService) RefreshApp(index int) (string, error) {
	s.mu.Lock()
	delete(s.cache, index)
	s.mu.Unlock()
	return s.RunApp(index)
}

func (s *AppService) DeleteApp(index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.apps) {
		return nil
	}
	s.apps = append(s.apps[:index], s.apps[index+1:]...)
	s.cache = make(map[int]cachedOutput)
	return s.save()
}
