package config_manager

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func (m *Manager) loadAndMerge(
	name string,
	targetType reflect.Type,
) (any, error) {
	basePath := filepath.Clean(filepath.Join(m.pathDefault, name+ExtensionYaml))
	mergePath := filepath.Clean(filepath.Join(m.pathMerge, PrefixMerge+name+ExtensionYaml))

	m.log.Debug("reading base configuration file", zap.String("path", basePath))
	baseData, err := m.readYamlToMap(basePath)
	if err != nil {
		return nil, fmt.Errorf("read default config: %w", err)
	}

	if _, err := os.Stat(mergePath); err == nil {
		m.log.Debug("reading merge configuration file", zap.String("path", mergePath))
		mergeData, err := m.readYamlToMap(mergePath)
		if err != nil {
			m.log.Warn(
				"failed to read merge file, skipping overrides",
				zap.String("path", mergePath),
				zap.Error(err),
			)
		} else {
			m.log.Debug("applying deep merge for configuration", zap.String("config", name))
			m.deepMerge(baseData, mergeData)
		}
	} else {
		m.log.Debug("no merge file found, using defaults", zap.String("config", name))
	}

	finalYaml, err := yaml.Marshal(baseData)
	if err != nil {
		return nil, fmt.Errorf("marshal merged data: %w", err)
	}

	resultStruct := reflect.New(targetType).Interface()

	decoder := yaml.NewDecoder(strings.NewReader(string(finalYaml)))
	decoder.KnownFields(true)

	if err := decoder.Decode(resultStruct); err != nil {
		return nil, fmt.Errorf("strict unmarshal error (check for unknown fields): %w", err)
	}

	m.log.Debug(
		"configuration successfully unmarshaled into struct",
		zap.String("config", name),
		zap.String("type", targetType.String()),
	)

	return reflect.ValueOf(resultStruct).Elem().Interface(), nil
}

func (m *Manager) readYamlToMap(path string) (map[string]any, error) {
	cleanPath := filepath.Clean(path)

	absDefault, _ := filepath.Abs(m.pathDefault)
	absMerge, _ := filepath.Abs(m.pathMerge)
	absPath, _ := filepath.Abs(cleanPath)

	if !strings.HasPrefix(absPath, absDefault) && !strings.HasPrefix(absPath, absMerge) {
		return nil, fmt.Errorf("path outside config directories: %s", path)
	}

	file, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	data := make(map[string]any)
	if err := yaml.Unmarshal(file, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *Manager) deepMerge(base, merge map[string]any) {
	for k, v := range merge {
		if baseVal, ok := base[k]; ok {
			baseMap, isBaseMap := baseVal.(map[string]any)
			mergeMap, isMergeMap := v.(map[string]any)
			if isBaseMap && isMergeMap {
				m.deepMerge(baseMap, mergeMap)
				continue
			}
		}
		base[k] = v
	}
}

func (m *Manager) ScanUsage() ScanReport {
	m.mu.RLock()
	defer m.mu.RUnlock()

	report := ScanReport{}
	files, err := os.ReadDir(m.pathDefault)
	if err != nil {
		m.log.Error("failed to scan default config directory", zap.Error(err))
		return report
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ExtensionYaml) {
			continue
		}

		name := strings.TrimSuffix(f.Name(), ExtensionYaml)
		_, isUsed := m.registry[name]

		report.Files = append(
			report.Files, FileStatus{
				Path:   filepath.Join(m.pathDefault, f.Name()),
				IsUsed: isUsed,
			},
		)
	}

	return report
}
