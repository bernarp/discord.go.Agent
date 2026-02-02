package module_manager

import (
	"context"
)

type Module interface {
	Name() string
	ConfigKey() string
	ConfigTemplate() any
	OnEnable(
		ctx context.Context,
		cfg any,
	)
	OnDisable(ctx context.Context)
	OnConfigUpdate(
		ctx context.Context,
		cfg any,
	)
}
