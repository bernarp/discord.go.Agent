package eventbus

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"

	"DiscordBotAgent/internal/core/zap_logger"
	"DiscordBotAgent/pkg/ctxtrace"
	"go.uber.org/zap"
)

type EventType string

type Handler func(
	ctx context.Context,
	payload any,
)

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]Handler
	log         *zap_logger.Logger
}

func New(log *zap_logger.Logger) *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]Handler),
		log:         log,
	}
}

func (eb *EventBus) Subscribe(
	eventType EventType,
	handler Handler,
) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
	eb.log.Debug("new subscription", zap.String("event", string(eventType)))
}

func (eb *EventBus) Publish(
	eventType EventType,
	payload any,
) {
	corrid := eb.generateHash()
	ctx := ctxtrace.WithCorrelationID(context.Background(), corrid)

	eb.mu.RLock()
	handlers, ok := eb.subscribers[eventType]
	eb.mu.RUnlock()

	if !ok {
		eb.log.WithCtx(ctx).Debug(
			"no subscribers for event",
			zap.String("event", string(eventType)),
		)
		return
	}

	eb.log.WithCtx(ctx).Info(
		"publishing event",
		zap.String("event", string(eventType)),
		zap.Int("handlers_count", len(handlers)),
	)

	for _, handler := range handlers {
		h := handler
		go h(ctx, payload)
	}
}

func (eb *EventBus) generateHash() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
