package module_manager

import (
	"sync"
	"time"
)

type moduleState struct {
	module       Module
	status       ModuleStatus
	configValid  bool
	currentCfg   any
	dependencies []string
	errorMessage string
	lastUpdated  time.Time
	mu           sync.RWMutex
}

func newModuleState(
	m Module,
	deps []string,
) *moduleState {
	return &moduleState{
		module:       m,
		status:       StatusDisabled,
		configValid:  false,
		dependencies: deps,
		lastUpdated:  time.Now(),
	}
}

func (s *moduleState) setEnabled(cfg any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = StatusEnabled
	s.configValid = true
	s.currentCfg = cfg
	s.errorMessage = ""
	s.lastUpdated = time.Now()
}

func (s *moduleState) setDisabled(reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = StatusDisabled
	s.errorMessage = reason
	s.lastUpdated = time.Now()
}

func (s *moduleState) setDepDisabled(depName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = StatusDepDisabled
	s.errorMessage = "dependency disabled: " + depName
	s.lastUpdated = time.Now()
}

func (s *moduleState) setError(err string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = StatusError
	s.configValid = false
	s.errorMessage = err
	s.lastUpdated = time.Now()
}

func (s *moduleState) updateConfig(cfg any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentCfg = cfg
	s.configValid = true
	s.lastUpdated = time.Now()
}

func (s *moduleState) getStatus() ModuleStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *moduleState) getInfo(dependents []string) ModuleInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return ModuleInfo{
		Name:         s.module.Name(),
		Status:       s.status,
		ConfigKey:    s.module.ConfigKey(),
		ConfigValid:  s.configValid,
		Dependencies: s.dependencies,
		Dependents:   dependents,
		ErrorMessage: s.errorMessage,
		LastUpdated:  s.lastUpdated,
	}
}

func (s *moduleState) getConfig() (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentCfg, s.configValid
}

func (s *moduleState) isEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status == StatusEnabled
}
