package module_manager

import "time"

type ModuleStatus string

const (
	StatusDisabled    ModuleStatus = "disabled"
	StatusEnabled     ModuleStatus = "enabled"
	StatusError       ModuleStatus = "error"
	StatusDepDisabled ModuleStatus = "dependency_disabled"
)

type ModuleInfo struct {
	Name         string
	Status       ModuleStatus
	ConfigKey    string
	ConfigValid  bool
	Dependencies []string
	Dependents   []string
	ErrorMessage string
	LastUpdated  time.Time
}
