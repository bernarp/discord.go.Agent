package config_manager

import (
	"context"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"

	"DiscordBotAgent/pkg/ctxtrace"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func (m *Manager) initWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	m.watcher = watcher

	go m.watchLoop()

	if err := m.watcher.Add(m.pathDefault); err != nil {
		return err
	}
	if err := m.watcher.Add(m.pathMerge); err != nil {
		return err
	}

	m.log.Debug("file system watcher initialized for config directories")
	return nil
}

func (m *Manager) watchLoop() {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				m.handleFileSystemEvent(event.Name, event.Op.String())
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
			m.log.Error("config watcher error", zap.Error(err))
		}
	}
}

func (m *Manager) handleFileSystemEvent(
	path string,
	op string,
) {
	fileName := filepath.Base(path)
	if !strings.HasSuffix(fileName, ExtensionYaml) {
		return
	}

	configName := strings.TrimSuffix(fileName, ExtensionYaml)
	configName = strings.TrimPrefix(configName, PrefixMerge)

	m.mu.RLock()
	_, registered := m.registry[configName]
	m.mu.RUnlock()

	if !registered {
		return
	}

	traceID := m.generateTraceID()
	ctx := ctxtrace.WithCorrelationID(context.Background(), traceID)

	m.log.WithCtx(ctx).Info(
		"file system event detected",
		zap.String("config", configName),
		zap.String("operation", op),
		zap.String("file", fileName),
	)

	m.mu.Lock()
	if timer, ok := m.timers[configName]; ok {
		timer.Stop()
	}

	m.timers[configName] = m.afterFunc(
		DebounceDuration, func() {
			m.reloadConfig(ctx, configName)
		},
	)
	m.mu.Unlock()

	m.log.WithCtx(ctx).Debug(
		"reload scheduled after debounce",
		zap.String("config", configName),
		zap.Duration("wait", DebounceDuration),
	)
}

func (m *Manager) reloadConfig(
	ctx context.Context,
	name string,
) {
	m.mu.Lock()
	meta, ok := m.registry[name]
	m.mu.Unlock()

	if !ok {
		return
	}

	m.log.WithCtx(ctx).Info("hot-reloading configuration", zap.String("config", name))

	newCfg, err := m.loadAndMerge(name, meta.StructType)
	if err != nil {
		m.log.WithCtx(ctx).Error(
			"hot-reload failed: load error. module will be disabled",
			zap.String("config", name),
			zap.Error(err),
		)
		if meta.OnUpdate != nil {
			go meta.OnUpdate(nil, false)
		}
		return
	}

	if err := m.validator.Validate(name, newCfg); err != nil {
		m.log.WithCtx(ctx).Error(
			"hot-reload failed: validation error. module will be disabled",
			zap.String("config", name),
		)
		if meta.OnUpdate != nil {
			go meta.OnUpdate(nil, false)
		}
		return
	}

	m.mu.Lock()
	meta.CurrentValue = newCfg
	m.mu.Unlock()

	if meta.OnUpdate != nil {
		m.log.WithCtx(ctx).Debug("triggering update callback (valid)", zap.String("config", name))
		go meta.OnUpdate(newCfg, true)
	}

	m.log.WithCtx(ctx).Info("configuration successfully reloaded", zap.String("config", name))
}

func (m *Manager) generateTraceID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("sys-%x", b)
}
