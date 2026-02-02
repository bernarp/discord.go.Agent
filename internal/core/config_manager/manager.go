package config_manager

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type Manager struct {
	mu          sync.RWMutex
	log         *zap_logger.Logger
	validator   *ConfigValidator
	registry    map[string]*ConfigMeta
	timers      map[string]*time.Timer
	watcher     *fsnotify.Watcher
	pathDefault string
	pathMerge   string
	afterFunc   func(
		time.Duration,
		func(),
	) *time.Timer
}

func New(
	log *zap_logger.Logger,
	pathDefault, pathMerge string,
) (*Manager, error) {
	for _, p := range []string{pathDefault, pathMerge} {
		if err := os.MkdirAll(p, 0750); err != nil {
			return nil, fmt.Errorf("failed to create config directory %s: %w", p, err)
		}
	}

	m := &Manager{
		log:         log,
		validator:   NewValidator(log),
		registry:    make(map[string]*ConfigMeta),
		timers:      make(map[string]*time.Timer),
		pathDefault: pathDefault,
		pathMerge:   pathMerge,
		afterFunc:   time.AfterFunc,
	}

	if err := m.initWatcher(); err != nil {
		return nil, fmt.Errorf("failed to init config watcher: %w", err)
	}

	return m, nil
}

func (m *Manager) Register(
	name string,
	template any,
	callback UpdateCallback,
) error {
	m.log.Info("registering configuration", zap.String("config", name))

	m.mu.Lock()
	defer m.mu.Unlock()

	t := reflect.TypeOf(template)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	cfg, err := m.loadAndMerge(name, t)
	if err != nil {
		return err
	}

	if err := m.validator.Validate(name, cfg); err != nil {
		return err
	}

	meta := &ConfigMeta{
		Name:         name,
		StructType:   t,
		OnUpdate:     callback,
		IsUsed:       true,
		CurrentValue: cfg,
	}

	m.registry[name] = meta

	if callback != nil {
		m.log.Debug("executing initial configuration callback", zap.String("config", name))
		callback(cfg, true)
	}

	m.log.Info("configuration registered and active", zap.String("config", name))
	return nil
}

func (m *Manager) Get(name string) any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	meta, ok := m.registry[name]
	if !ok {
		m.log.Warn("attempted to get unregistered configuration", zap.String("config", name))
		return nil
	}

	return meta.CurrentValue
}

func (m *Manager) PrintReport() {
	m.log.Info("starting configuration usage audit")
	report := m.ScanUsage()

	unusedCount := 0
	for _, f := range report.Files {
		if !f.IsUsed {
			m.log.Warn("unused configuration file detected", zap.String("path", f.Path))
			unusedCount++
		} else {
			m.log.Debug("active configuration file", zap.String("path", f.Path))
		}
	}

	m.log.Info(
		"configuration audit finished",
		zap.Int("total_files", len(report.Files)),
		zap.Int("unused_files", unusedCount),
	)
}

func (m *Manager) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}
