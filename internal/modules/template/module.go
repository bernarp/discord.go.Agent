package template

import (
	"context"

	"DiscordBotAgent/internal/core/config_manager"
	"DiscordBotAgent/internal/core/eventbus"
	"DiscordBotAgent/internal/core/module_manager"
	"DiscordBotAgent/internal/core/zap_logger"
)

const ModuleName = "template"

type Module struct {
	log           *zap_logger.Logger
	eb            *eventbus.EventBus
	mm            *module_manager.Manager
	handler       *Handler
	cfg           Config
	subscriptions []eventbus.SubscriptionID
}

func New(
	log *zap_logger.Logger,
	eb *eventbus.EventBus,
	mm *module_manager.Manager,
) *Module {
	m := &Module{
		log:           log,
		eb:            eb,
		mm:            mm,
		subscriptions: make([]eventbus.SubscriptionID, 0),
	}
	m.handler = NewHandler(NewService(log), m)
	return m
}

func (m *Module) Name() string {
	return ModuleName
}

func (m *Module) ConfigKey() string {
	return config_manager.Contract.System.Discord.Template
}

func (m *Module) ConfigTemplate() any {
	defaultEnabled := true
	return Config{
		Enabled: &defaultEnabled,
		LogDetails: LogDetails{
			Guild:   true,
			Channel: true,
			Author:  true,
			Content: true,
		},
	}
}

func (m *Module) OnEnable(
	ctx context.Context,
	cfg any,
) {
	m.cfg = cfg.(Config)

	id := m.eb.Subscribe(eventbus.MessageCreate, m.handler.OnMessageCreate)
	m.subscriptions = append(m.subscriptions, id)

	m.log.Info("template module: enabled")
}

func (m *Module) OnDisable(ctx context.Context) {
	for _, id := range m.subscriptions {
		m.eb.Unsubscribe(id)
	}
	m.subscriptions = make([]eventbus.SubscriptionID, 0)

	m.log.Info("template module: disabled")
}

func (m *Module) OnConfigUpdate(
	ctx context.Context,
	cfg any,
) {
	m.cfg = cfg.(Config)
}

func (m *Module) GetConfig() Config {
	return m.cfg
}
