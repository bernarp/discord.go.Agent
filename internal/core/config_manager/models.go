package config_manager

import (
	"reflect"
)

type UpdateCallback func(
	cfg any,
	isValid bool,
)

type ConfigMeta struct {
	Name         string
	StructType   reflect.Type
	OnUpdate     UpdateCallback
	CurrentValue any
	IsUsed       bool
}

type FileStatus struct {
	Path   string
	IsUsed bool
}

type ScanReport struct {
	Files []FileStatus
}
