// internal/core/config_manager/manager_test.go
package config_manager

import (
	"DiscordBotAgent/internal/core/zap_logger"
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Name    string `yaml:"name" validate:"required"`
	Enabled bool   `yaml:"enabled"`
	Count   int    `yaml:"count" validate:"gte=0,lte=100"`
}

func setupTestManager(t *testing.T) (*Manager, string, string) {
	t.Helper()

	tmpDir := t.TempDir()
	defaultPath := filepath.Join(tmpDir, "config_df")
	mergePath := filepath.Join(tmpDir, "config_mrg")

	logger, err := zap_logger.New()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	m, err := New(logger, defaultPath, mergePath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	return m, defaultPath, mergePath
}

func writeYamlFile(
	t *testing.T,
	path, content string,
) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func TestNew_CreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	defaultPath := filepath.Join(tmpDir, "new_default")
	mergePath := filepath.Join(tmpDir, "new_merge")

	logger, _ := zap_logger.New()
	_, err := New(logger, defaultPath, mergePath)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Error("default directory was not created")
	}
	if _, err := os.Stat(mergePath); os.IsNotExist(err) {
		t.Error("merge directory was not created")
	}
}

func TestRegister_Success(t *testing.T) {
	m, defaultPath, _ := setupTestManager(t)
	defer m.Close()

	yaml := `
name: "test"
enabled: true
count: 50
`
	writeYamlFile(t, filepath.Join(defaultPath, "testconfig.yaml"), yaml)

	var callbackCalled bool
	err := m.Register(
		"testconfig",
		TestConfig{},
		func(
			cfg any,
			isValid bool,
		) {
			callbackCalled = true
		},
	)

	if err != nil {
		t.Fatalf("Register() error: %v", err)
	}
	if !callbackCalled {
		t.Error("callback was not called on register")
	}
}

func TestRegister_ValidationError(t *testing.T) {
	m, defaultPath, _ := setupTestManager(t)
	defer m.Close()

	yaml := `
enabled: true
count: 50
`
	writeYamlFile(t, filepath.Join(defaultPath, "invalid.yaml"), yaml)

	err := m.Register("invalid", TestConfig{}, nil)
	if err == nil {
		t.Error("expected validation error, got nil")
	}
}

func TestRegister_FileNotFound(t *testing.T) {
	m, _, _ := setupTestManager(t)
	defer m.Close()

	err := m.Register("nonexistent", TestConfig{}, nil)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestGet_ReturnsRegisteredConfig(t *testing.T) {
	m, defaultPath, _ := setupTestManager(t)
	defer m.Close()

	yaml := `
name: "myconfig"
enabled: true
count: 42
`
	writeYamlFile(t, filepath.Join(defaultPath, "gettest.yaml"), yaml)

	_ = m.Register("gettest", TestConfig{}, nil)

	result := m.Get("gettest")
	if result == nil {
		t.Fatal("Get() returned nil")
	}

	cfg, ok := result.(TestConfig)
	if !ok {
		t.Fatal("Get() returned wrong type")
	}

	if cfg.Name != "myconfig" {
		t.Errorf("expected name 'myconfig', got '%s'", cfg.Name)
	}
	if cfg.Count != 42 {
		t.Errorf("expected count 42, got %d", cfg.Count)
	}
}

func TestGet_UnregisteredConfig(t *testing.T) {
	m, _, _ := setupTestManager(t)
	defer m.Close()

	result := m.Get("unknown")
	if result != nil {
		t.Error("expected nil for unregistered config")
	}
}

func TestClose_NoError(t *testing.T) {
	m, _, _ := setupTestManager(t)

	err := m.Close()
	if err != nil {
		t.Errorf("Close() error: %v", err)
	}
}
