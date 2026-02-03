package apierror

import (
	_ "embed"
	"fmt"
	"reflect"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

//go:embed errors_code.yaml
var defaultErrorsConfig []byte

type errorRegistry struct {
	INTERNAL_ERROR    *AppError
	INVALID_REQUEST   *AppError
	VALIDATION_FAILED *AppError

	UNAUTHORIZED        *AppError
	INVALID_API_KEY     *AppError
	PERMISSION_DENIED   *AppError
	RATE_LIMIT_EXCEEDED *AppError

	MODULE_NOT_FOUND          *AppError
	MODULE_ALREADY_ENABLED    *AppError
	MODULE_ALREADY_DISABLED   *AppError
	MODULE_DEPENDENCY_MISSING *AppError
	MODULE_HAS_DEPENDENTS     *AppError

	CONFIG_NOT_FOUND   *AppError
	CONFIG_INVALID     *AppError
	CONFIG_PARSE_ERROR *AppError

	DISCORD_NOT_CONNECTED     *AppError
	DISCORD_GUILD_NOT_FOUND   *AppError
	DISCORD_CHANNEL_NOT_FOUND *AppError
	DISCORD_API_ERROR         *AppError

	EVENT_HANDLER_TIMEOUT *AppError
}

var Errors = &errorRegistry{
	INTERNAL_ERROR:            &AppError{Code: "INTERNAL_ERROR", Status: 500},
	INVALID_REQUEST:           &AppError{Code: "INVALID_REQUEST", Status: 400},
	VALIDATION_FAILED:         &AppError{Code: "VALIDATION_FAILED", Status: 400},
	UNAUTHORIZED:              &AppError{Code: "UNAUTHORIZED", Status: 401},
	INVALID_API_KEY:           &AppError{Code: "INVALID_API_KEY", Status: 401},
	PERMISSION_DENIED:         &AppError{Code: "PERMISSION_DENIED", Status: 403},
	RATE_LIMIT_EXCEEDED:       &AppError{Code: "RATE_LIMIT_EXCEEDED", Status: 429},
	MODULE_NOT_FOUND:          &AppError{Code: "MODULE_NOT_FOUND", Status: 404},
	MODULE_ALREADY_ENABLED:    &AppError{Code: "MODULE_ALREADY_ENABLED", Status: 409},
	MODULE_ALREADY_DISABLED:   &AppError{Code: "MODULE_ALREADY_DISABLED", Status: 409},
	MODULE_DEPENDENCY_MISSING: &AppError{Code: "MODULE_DEPENDENCY_MISSING", Status: 424},
	MODULE_HAS_DEPENDENTS:     &AppError{Code: "MODULE_HAS_DEPENDENTS", Status: 409},
	CONFIG_NOT_FOUND:          &AppError{Code: "CONFIG_NOT_FOUND", Status: 404},
	CONFIG_INVALID:            &AppError{Code: "CONFIG_INVALID", Status: 400},
	CONFIG_PARSE_ERROR:        &AppError{Code: "CONFIG_PARSE_ERROR", Status: 400},
	DISCORD_NOT_CONNECTED:     &AppError{Code: "DISCORD_NOT_CONNECTED", Status: 503},
	DISCORD_GUILD_NOT_FOUND:   &AppError{Code: "DISCORD_GUILD_NOT_FOUND", Status: 404},
	DISCORD_CHANNEL_NOT_FOUND: &AppError{Code: "DISCORD_CHANNEL_NOT_FOUND", Status: 404},
	DISCORD_API_ERROR:         &AppError{Code: "DISCORD_API_ERROR", Status: 502},
	EVENT_HANDLER_TIMEOUT:     &AppError{Code: "EVENT_HANDLER_TIMEOUT", Status: 504},
}

var (
	log      *zap.Logger
	initOnce sync.Once
	initErr  error
)

type yamlConfig struct {
	Errors map[string]struct {
		Status  int    `yaml:"status"`
		Message string `yaml:"message"`
	} `yaml:"errors"`
}

func Init(logger *zap.Logger) error {
	initOnce.Do(
		func() {
			log = logger
			initErr = loadFromYAML(defaultErrorsConfig)
		},
	)
	return initErr
}

func loadFromYAML(data []byte) error {
	var config yamlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parse error config: %w", err)
	}

	val := reflect.ValueOf(Errors).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		field := val.Field(i)
		appErr := field.Interface().(*AppError)

		if entry, ok := config.Errors[fieldName]; ok {
			appErr.Message = entry.Message
			appErr.Status = entry.Status
		} else if log != nil {
			log.Warn("error code missing in YAML", zap.String("code", fieldName))
		}
	}

	return nil
}
