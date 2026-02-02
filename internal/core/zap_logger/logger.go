package zap_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"DiscordBotAgent/pkg/ctxtrace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New() (*Logger, error) {
	startTime := time.Now().Format("2006-01-02_15-04-05")
	logDir := filepath.Join("logs", startTime)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, "logs.json")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	return &Logger{
		Logger: zap.New(core, zap.AddCaller()),
	}, nil
}

func (l *Logger) WithCtx(ctx context.Context) *zap.Logger {
	id := ctxtrace.Extract(ctx)
	if id == "" {
		return l.Logger
	}
	return l.Logger.With(zap.String("corrid", id))
}
