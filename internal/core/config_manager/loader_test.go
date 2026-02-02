package config_manager

import (
	"path/filepath"
	"reflect"
	"testing"

	"DiscordBotAgent/internal/core/zap_logger"
)

func TestDeepMerge_OverwritesPrimitives(t *testing.T) {
	logger, _ := zap_logger.New()
	m := &Manager{log: logger}

	base := map[string]any{
		"name":  "original",
		"count": 10,
	}
	merge := map[string]any{
		"name": "overwritten",
	}

	m.deepMerge(base, merge)

	if base["name"] != "overwritten" {
		t.Errorf("expected 'overwritten', got '%v'", base["name"])
	}
	if base["count"] != 10 {
		t.Errorf("expected count to remain 10, got %v", base["count"])
	}
}

func TestDeepMerge_MergesNestedMaps(t *testing.T) {
	logger, _ := zap_logger.New()
	m := &Manager{log: logger}

	base := map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
	}
	merge := map[string]any{
		"server": map[string]any{
			"port": 9090,
		},
	}

	m.deepMerge(base, merge)

	server := base["server"].(map[string]any)
	if server["host"] != "localhost" {
		t.Errorf("expected host 'localhost', got '%v'", server["host"])
	}
	if server["port"] != 9090 {
		t.Errorf("expected port 9090, got %v", server["port"])
	}
}

func TestDeepMerge_AddsNewKeys(t *testing.T) {
	logger, _ := zap_logger.New()
	m := &Manager{log: logger}

	base := map[string]any{
		"existing": "value",
	}
	merge := map[string]any{
		"newkey": "newvalue",
	}

	m.deepMerge(base, merge)

	if base["newkey"] != "newvalue" {
		t.Errorf("expected new key to be added")
	}
}

func TestLoadAndMerge_BaseOnly(t *testing.T) {
	m, defaultPath, _ := setupTestManager(t)
	defer m.Close()

	yaml := `
name: "baseonly"
enabled: true
count: 25
`
	writeYamlFile(t, filepath.Join(defaultPath, "baseonly.yaml"), yaml)

	result, err := m.loadAndMerge("baseonly", reflect.TypeOf(TestConfig{}))
	if err != nil {
		t.Fatalf("loadAndMerge() error: %v", err)
	}

	cfg := result.(TestConfig)
	if cfg.Name != "baseonly" {
		t.Errorf("expected 'baseonly', got '%s'", cfg.Name)
	}
}

func TestLoadAndMerge_WithMergeFile(t *testing.T) {
	m, defaultPath, mergePath := setupTestManager(t)
	defer m.Close()

	baseYaml := `
name: "original"
enabled: false
count: 10
`
	mergeYaml := `
name: "merged"
count: 99
`
	writeYamlFile(t, filepath.Join(defaultPath, "mergetest.yaml"), baseYaml)
	writeYamlFile(t, filepath.Join(mergePath, "MERGE.mergetest.yaml"), mergeYaml)

	result, err := m.loadAndMerge("mergetest", reflect.TypeOf(TestConfig{}))
	if err != nil {
		t.Fatalf("loadAndMerge() error: %v", err)
	}

	cfg := result.(TestConfig)
	if cfg.Name != "merged" {
		t.Errorf("expected name 'merged', got '%s'", cfg.Name)
	}
	if cfg.Enabled != false {
		t.Error("expected enabled to remain false")
	}
	if cfg.Count != 99 {
		t.Errorf("expected count 99, got %d", cfg.Count)
	}
}

func TestScanUsage_DetectsUnusedFiles(t *testing.T) {
	m, defaultPath, _ := setupTestManager(t)
	defer m.Close()

	writeYamlFile(t, filepath.Join(defaultPath, "unused.yaml"), "name: test")
	writeYamlFile(t, filepath.Join(defaultPath, "used.yaml"), "name: used\ncount: 1")
	_ = m.Register("used", TestConfig{}, nil)

	report := m.ScanUsage()

	if len(report.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(report.Files))
	}

	unusedFound := false
	usedFound := false
	for _, f := range report.Files {
		if filepath.Base(f.Path) == "unused.yaml" && !f.IsUsed {
			unusedFound = true
		}
		if filepath.Base(f.Path) == "used.yaml" && f.IsUsed {
			usedFound = true
		}
	}

	if !unusedFound {
		t.Error("unused file not detected")
	}
	if !usedFound {
		t.Error("used file not marked as used")
	}
}
