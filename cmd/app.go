package main

import (
	"fmt"

	"DiscordBotAgent/internal/client"
	config "DiscordBotAgent/internal/core/config_env"
	"DiscordBotAgent/internal/core/zap_logger"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	log    *zap.Logger
	client *client.Client
}

func New() (*App, error) {
	logger, err := zap_logger.New()
	if err != nil {
		return nil, fmt.Errorf("app logger: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("app config: %w", err)
	}

	discordClient, err := client.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("app client: %w", err)
	}

	return &App{
		cfg:    cfg,
		log:    logger,
		client: discordClient,
	}, nil
}

func (a *App) Run() error {
	defer a.log.Sync()

	if err := a.client.Connect(); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	a.log.Info("application started")

	a.WaitGracefulShutdown()

	return nil
}
