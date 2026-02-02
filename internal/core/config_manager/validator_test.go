package config_manager

import (
	"testing"

	"DiscordBotAgent/internal/core/zap_logger"
)

type ValidConfig struct {
	Name  string `validate:"required"`
	Count int    `validate:"gte=0,lte=100"`
}

type InvalidConfig struct {
	Name  string `validate:"required"`
	Count int    `validate:"gte=0,lte=10"`
}

func TestValidate_Success(t *testing.T) {
	logger, _ := zap_logger.New()
	v := NewValidator(logger)

	cfg := ValidConfig{
		Name:  "test",
		Count: 50,
	}

	err := v.Validate("testconfig", cfg)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	logger, _ := zap_logger.New()
	v := NewValidator(logger)

	cfg := ValidConfig{
		Name:  "",
		Count: 50,
	}

	err := v.Validate("testconfig", cfg)
	if err == nil {
		t.Error("expected validation error for missing required field")
	}
}

func TestValidate_OutOfRange(t *testing.T) {
	logger, _ := zap_logger.New()
	v := NewValidator(logger)

	cfg := InvalidConfig{
		Name:  "test",
		Count: 999,
	}

	err := v.Validate("testconfig", cfg)
	if err == nil {
		t.Error("expected validation error for out of range value")
	}
}

func TestValidate_NegativeValue(t *testing.T) {
	logger, _ := zap_logger.New()
	v := NewValidator(logger)

	cfg := ValidConfig{
		Name:  "test",
		Count: -5,
	}

	err := v.Validate("testconfig", cfg)
	if err == nil {
		t.Error("expected validation error for negative value")
	}
}
