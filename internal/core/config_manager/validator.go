package config_manager

import (
	"fmt"

	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ConfigValidator struct {
	validate *validator.Validate
	log      *zap_logger.Logger
}

func NewValidator(log *zap_logger.Logger) *ConfigValidator {
	return &ConfigValidator{
		validate: validator.New(),
		log:      log,
	}
}

func (v *ConfigValidator) Validate(
	name string,
	cfg any,
) error {
	v.log.Debug("starting configuration validation", zap.String("config", name))

	err := v.validate.Struct(cfg)
	if err != nil {
		v.log.Error(
			"configuration validation failed",
			zap.String("config", name),
			zap.Error(err),
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	v.log.Debug("configuration validation passed", zap.String("config", name))
	return nil
}
