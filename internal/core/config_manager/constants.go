package config_manager

import "time"

const (
	// ExtensionYaml расширение файлов конфигурации
	ExtensionYaml = ".yaml"

	// PrefixMerge префикс для файлов переопределения
	PrefixMerge = "MERGE."

	// DebounceDuration время ожидания после последнего события ФС перед перезагрузкой
	DebounceDuration = 200 * time.Millisecond
)
