package module_manager

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"DiscordBotAgent/internal/core/config_manager"
	"DiscordBotAgent/internal/core/zap_logger"

	"go.uber.org/zap"
)

type Manager struct {
	log        *zap_logger.Logger
	cm         *config_manager.Manager
	modules    map[string]*moduleState
	dependents map[string][]string
	mu         sync.RWMutex
}

func New(
	log *zap_logger.Logger,
	cm *config_manager.Manager,
) *Manager {
	return &Manager{
		log:        log,
		cm:         cm,
		modules:    make(map[string]*moduleState),
		dependents: make(map[string][]string),
	}
}

func (m *Manager) Register(mod Module) error {
	name := mod.Name()

	m.mu.RLock()
	_, exists := m.modules[name]
	m.mu.RUnlock()

	if exists {
		return fmt.Errorf("module %s already registered", name)
	}

	deps := scanDependencies(mod)
	state := newModuleState(mod, deps)

	m.mu.Lock()
	m.modules[name] = state
	for _, dep := range deps {
		m.dependents[dep] = append(m.dependents[dep], name)
	}
	m.mu.Unlock()

	if configKey := mod.ConfigKey(); configKey != "" {
		template := mod.ConfigTemplate()

		err := m.cm.Register(
			configKey,
			template,
			func(
				cfg any,
				isValid bool,
			) {
				m.onConfigUpdate(name, cfg, isValid)
			},
		)

		if err != nil {
			if errors.Is(err, config_manager.ErrPlaceholderCreated) {
				m.log.Warn(
					"MODULE DISABLED: configuration file was missing",
					zap.String("module", name),
					zap.String("config_file", configKey+config_manager.ExtensionYaml),
					zap.String("action", "Please check config_df folder and fill the generated file"),
				)
				state.setDisabled("missing configuration (placeholder created)")
				return nil
			}

			state.setError(fmt.Sprintf("config registration failed: %v", err))
			return fmt.Errorf("module %s config registration: %w", name, err)
		}
	} else {
		m.tryEnable(name, nil)
	}

	return nil
}

func (m *Manager) onConfigUpdate(
	moduleName string,
	cfg any,
	isValid bool,
) {
	m.mu.RLock()
	state, exists := m.modules[moduleName]
	m.mu.RUnlock()

	if !exists {
		return
	}

	ctx := context.Background()
	wasEnabled := state.isEnabled()

	if !isValid {
		state.setDisabled("invalid configuration")
		if wasEnabled {
			state.module.OnDisable(ctx)
			m.log.Warn(
				"module disabled due to invalid config",
				zap.String("module", moduleName),
			)
			m.disableDependents(moduleName)
		}
		return
	}

	state.updateConfig(cfg)

	if wasEnabled {
		state.module.OnConfigUpdate(ctx, cfg)
		m.log.Info("module config updated", zap.String("module", moduleName))
	} else {
		m.tryEnable(moduleName, cfg)
	}
}

func (m *Manager) tryEnable(
	moduleName string,
	cfg any,
) {
	m.mu.RLock()
	state, exists := m.modules[moduleName]
	m.mu.RUnlock()

	if !exists || state.isEnabled() {
		return
	}

	for _, depName := range state.dependencies {
		m.mu.RLock()
		depState, depExists := m.modules[depName]
		m.mu.RUnlock()

		if !depExists {
			state.setError(fmt.Sprintf("dependency %s not registered", depName))
			m.log.Error(
				"module dependency not found",
				zap.String("module", moduleName),
				zap.String("dependency", depName),
			)
			return
		}

		if !depState.isEnabled() {
			state.setDepDisabled(depName)
			m.log.Warn(
				"module waiting for dependency",
				zap.String("module", moduleName),
				zap.String("dependency", depName),
			)
			return
		}
	}

	ctx := context.Background()

	if cfg == nil {
		cfg, _ = state.getConfig()
	}

	state.setEnabled(cfg)
	state.module.OnEnable(ctx, cfg)
	m.log.Info("module enabled", zap.String("module", moduleName))

	m.tryEnableDependents(moduleName)
}

func (m *Manager) tryEnableDependents(moduleName string) {
	m.mu.RLock()
	deps := m.dependents[moduleName]
	m.mu.RUnlock()

	for _, depName := range deps {
		m.mu.RLock()
		state, exists := m.modules[depName]
		m.mu.RUnlock()

		if exists && state.getStatus() == StatusDepDisabled {
			cfg, _ := state.getConfig()
			m.tryEnable(depName, cfg)
		}
	}
}

func (m *Manager) disableDependents(moduleName string) {
	m.mu.RLock()
	deps := m.dependents[moduleName]
	m.mu.RUnlock()

	ctx := context.Background()

	for _, depName := range deps {
		m.mu.RLock()
		state, exists := m.modules[depName]
		m.mu.RUnlock()

		if exists && state.isEnabled() {
			state.setDepDisabled(moduleName)
			state.module.OnDisable(ctx)
			m.log.Warn(
				"module disabled due to dependency",
				zap.String("module", depName),
				zap.String("dependency", moduleName),
			)

			m.disableDependents(depName)
		}
	}
}

func (m *Manager) Disable(moduleName string) error {
	m.mu.RLock()
	state, exists := m.modules[moduleName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("module %s not found", moduleName)
	}

	if !state.isEnabled() {
		return nil
	}

	ctx := context.Background()
	state.setDisabled("manually disabled")
	state.module.OnDisable(ctx)

	m.log.Info("module manually disabled", zap.String("module", moduleName))

	m.disableDependents(moduleName)

	return nil
}

func (m *Manager) Enable(moduleName string) error {
	m.mu.RLock()
	state, exists := m.modules[moduleName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("module %s not found", moduleName)
	}

	if state.isEnabled() {
		return nil
	}

	cfg, valid := state.getConfig()
	if !valid && state.module.ConfigKey() != "" {
		return fmt.Errorf("module %s has invalid config", moduleName)
	}

	m.tryEnable(moduleName, cfg)
	return nil
}

func (m *Manager) GetConfig(moduleName string) (any, bool) {
	m.mu.RLock()
	state, exists := m.modules[moduleName]
	m.mu.RUnlock()

	if !exists {
		return nil, false
	}

	return state.getConfig()
}

func (m *Manager) GetModuleInfo(moduleName string) (ModuleInfo, bool) {
	m.mu.RLock()
	state, exists := m.modules[moduleName]
	deps := m.dependents[moduleName]
	m.mu.RUnlock()

	if !exists {
		return ModuleInfo{}, false
	}

	return state.getInfo(deps), true
}

func (m *Manager) GetAllModules() []ModuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]ModuleInfo, 0, len(m.modules))
	for name, state := range m.modules {
		result = append(result, state.getInfo(m.dependents[name]))
	}
	return result
}

func (m *Manager) PrintReport() {
	for _, info := range m.GetAllModules() {
		fields := []zap.Field{
			zap.String("module", info.Name),
			zap.String("status", string(info.Status)),
		}
		if len(info.Dependencies) > 0 {
			fields = append(fields, zap.Strings("depends_on", info.Dependencies))
		}
		if len(info.Dependents) > 0 {
			fields = append(fields, zap.Strings("required_by", info.Dependents))
		}
		if info.ErrorMessage != "" {
			fields = append(fields, zap.String("error", info.ErrorMessage))
		}
		m.log.Info("module status", fields...)
	}
}

func (m *Manager) IsModuleEnabled(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if state, ok := m.modules[name]; ok {
		return state.isEnabled()
	}
	return false
}
