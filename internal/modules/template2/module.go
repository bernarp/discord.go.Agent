package template2

import (
	"context"

	"DiscordBotAgent/internal/core/config_manager"
	"DiscordBotAgent/internal/core/eventbus"
	"DiscordBotAgent/internal/core/module_manager"
	"DiscordBotAgent/internal/core/zap_logger"
	"DiscordBotAgent/internal/modules/template"
)

const ModuleName = "template2"

type Module struct {
	log           *zap_logger.Logger
	eb            *eventbus.EventBus
	mm            *module_manager.Manager
	template      *template.Module
	handler       *Handler
	cfg           Config
	subscriptions []eventbus.SubscriptionID
}

func New(
	log *zap_logger.Logger,
	eb *eventbus.EventBus,
	mm *module_manager.Manager,
	template *template.Module,
) *Module {
	m := &Module{
		log:           log,
		eb:            eb,
		mm:            mm,
		template:      template,
		subscriptions: make([]eventbus.SubscriptionID, 0),
	}
	m.handler = NewHandler(NewService(log), m)
	return m
}

func (m *Module) Name() string {
	return ModuleName
}

func (m *Module) ConfigKey() string {
	return config_manager.Contract.System.Discord.Template2
}

func (m *Module) ConfigTemplate() any {
	return Config{
		Prefix:  "!",
		Enabled: true,
		MaxLogs: 100,
	}
}

func (m *Module) OnEnable(
	ctx context.Context,
	cfg any,
) {
	m.cfg = cfg.(Config)

	id := m.eb.Subscribe(eventbus.MessageCreate, m.handler.OnMessageCreate)
	m.subscriptions = append(m.subscriptions, id)

	m.log.Info("template2 module: enabled")
}

func (m *Module) OnDisable(ctx context.Context) {
	for _, id := range m.subscriptions {
		m.eb.Unsubscribe(id)
	}
	m.subscriptions = make([]eventbus.SubscriptionID, 0)

	m.log.Info("template2 module: disabled")
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
