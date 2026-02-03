package main

import (
	"fmt"

	"DiscordBotAgent/internal/client"
	config "DiscordBotAgent/internal/core/config_env"
	"DiscordBotAgent/internal/core/config_manager"
	"DiscordBotAgent/internal/core/eventbus"
	"DiscordBotAgent/internal/core/module_manager"
	"DiscordBotAgent/internal/core/zap_logger"
	"DiscordBotAgent/internal/modules/template"
	"DiscordBotAgent/internal/modules/template2"
)

type App struct {
	cfg       *config.Config
	log       *zap_logger.Logger
	eb        *eventbus.EventBus
	configMgr *config_manager.Manager
	moduleMgr *module_manager.Manager
	client    *client.Client
}

func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("app config: %w", err)
	}

	logger, err := zap_logger.New()
	if err != nil {
		return nil, fmt.Errorf("app logger: %w", err)
	}

	configMgr, err := config_manager.New(logger, "config_df", "config_mrg")
	if err != nil {
		return nil, fmt.Errorf("config manager: %w", err)
	}

	eb := eventbus.New(logger)
	moduleMgr := module_manager.New(logger, configMgr)

	templateMod := template.New(logger, eb, moduleMgr)
	if err := moduleMgr.Register(templateMod); err != nil {
		return nil, fmt.Errorf("template module: %w", err)
	}

	template2Mod := template2.New(logger, eb, moduleMgr, templateMod)
	if err := moduleMgr.Register(template2Mod); err != nil {
		return nil, fmt.Errorf("template2 module: %w", err)
	}

	discordClient, err := client.New(cfg, logger, eb)
	if err != nil {
		return nil, fmt.Errorf("app client: %w", err)
	}

	return &App{
		cfg:       cfg,
		log:       logger,
		eb:        eb,
		configMgr: configMgr,
		moduleMgr: moduleMgr,
		client:    discordClient,
	}, nil
}

func (a *App) Run() error {
	defer func() {
		_ = a.configMgr.Close()
		_ = a.log.Sync()
	}()

	if err := a.client.Connect(); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	a.log.Info("application started")
	a.configMgr.PrintReport()
	a.moduleMgr.PrintReport()

	a.WaitGracefulShutdown()

	return nil
}
