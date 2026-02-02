package eventbus

import (
	"context"
	"crypto/rand"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"DiscordBotAgent/internal/core/zap_logger"
	"DiscordBotAgent/pkg/ctxtrace"
	"go.uber.org/zap"
)

const handlerTimeout = 15 * time.Second

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
}

func (eb *EventBus) Publish(
	eventType EventType,
	payload any,
) {
	corrid := eb.generateHash()
	baseCtx := ctxtrace.WithCorrelationID(context.Background(), corrid)

	eb.mu.RLock()
	handlers, ok := eb.subscribers[eventType]
	eb.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		h := handler
		go eb.executeHandler(baseCtx, string(eventType), h, payload)
	}
}

func (eb *EventBus) executeHandler(
	ctx context.Context,
	eventName string,
	h Handler,
	payload any,
) {
	ctx, cancel := context.WithTimeout(ctx, handlerTimeout)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			eb.log.WithCtx(ctx).Error(
				"event handler panicked",
				zap.String("event", eventName),
				zap.Any("error", r),
				zap.String("stack", string(debug.Stack())),
			)
		}
	}()

	h(ctx, payload)
}

func (eb *EventBus) generateHash() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
