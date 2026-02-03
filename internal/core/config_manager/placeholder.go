package config_manager

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func (m *Manager) createPlaceholder(
	path string,
	template any,
) error {
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("marshal template: %w", err)
	}

	content := []byte(PlaceholderHeader)
	content = append(content, data...)

	if err := os.WriteFile(path, content, 0600); err != nil {
		return fmt.Errorf("write placeholder: %w", err)
	}

	return nil
}
