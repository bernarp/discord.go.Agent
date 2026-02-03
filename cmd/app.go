package main

import (
	"fmt"

	"DiscordBotAgent/internal/api"
	"DiscordBotAgent/internal/api/apierror"
	"DiscordBotAgent/internal/client"
	config "DiscordBotAgent/internal/core/config_env"
	"DiscordBotAgent/internal/core/config_manager"
	"DiscordBotAgent/internal/core/eventbus"
	"DiscordBotAgent/internal/core/module_manager"
	"DiscordBotAgent/internal/core/zap_logger"
	"DiscordBotAgent/internal/modules/template"
	"DiscordBotAgent/internal/modules/template2"

	"go.uber.org/zap"
)

type App struct {
	cfg       *config.Config
	log       *zap_logger.Logger
	eb        *eventbus.EventBus
	configMgr *config_manager.Manager
	moduleMgr *module_manager.Manager
	client    *client.Client
	api       *api.Server
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

	if err := apierror.Init(logger.Logger); err != nil {
		return nil, fmt.Errorf("api errors init: %w", err)
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

	apiServer := api.New(logger, moduleMgr)

	return &App{
		cfg:       cfg,
		log:       logger,
		eb:        eb,
		configMgr: configMgr,
		moduleMgr: moduleMgr,
		client:    discordClient,
		api:       apiServer,
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

	if a.cfg.Port != "" {
		if err := a.api.Start(a.cfg.Port); err != nil {
			a.log.Error("failed to start api server", zap.Error(err))
		}
	} else {
		a.log.Info("api server skipped (no port provided)")
	}

	a.log.Info("application started")
	a.configMgr.PrintReport()
	a.moduleMgr.PrintReport()

	a.WaitGracefulShutdown()

	return nil
}
