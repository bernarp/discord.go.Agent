package template

import (
	"DiscordBotAgent/internal/core/eventbus"
	"DiscordBotAgent/internal/core/zap_logger"
)

type Module struct {
	handler *Handler
	eb      *eventbus.EventBus
}

func New(
	log *zap_logger.Logger,
	eb *eventbus.EventBus,
) *Module {
	service := NewService(log)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
		eb:      eb,
	}
}

func (m *Module) Init() {
	m.eb.Subscribe(eventbus.MessageCreate, m.handler.OnMessageCreate)
}
