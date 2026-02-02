// internal/core/eventbus/eventbus.go
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

type SubscriptionID string

type subscriber struct {
	id      SubscriptionID
	handler Handler
}

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]subscriber
	log         *zap_logger.Logger
	idCounter   int64
}

func New(log *zap_logger.Logger) *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]subscriber),
		log:         log,
	}
}

func (eb *EventBus) Subscribe(
	eventType EventType,
	handler Handler,
) SubscriptionID {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.idCounter++
	id := SubscriptionID(fmt.Sprintf("%s_%d", eventType, eb.idCounter))

	eb.subscribers[eventType] = append(
		eb.subscribers[eventType], subscriber{
			id:      id,
			handler: handler,
		},
	)

	eb.log.Debug(
		"handler subscribed",
		zap.String("event", string(eventType)),
		zap.String("id", string(id)),
	)

	return id
}

func (eb *EventBus) Unsubscribe(id SubscriptionID) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for eventType, subs := range eb.subscribers {
		for i, sub := range subs {
			if sub.id == id {
				eb.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
				eb.log.Debug(
					"handler unsubscribed",
					zap.String("event", string(eventType)),
					zap.String("id", string(id)),
				)
				return
			}
		}
	}
}

func (eb *EventBus) Publish(
	eventType EventType,
	payload any,
) {
	corrid, err := eb.generateHash()
	if err != nil {
		eb.log.Error("failed to generate correlation id", zap.Error(err))
		corrid = "unknown"
	}

	baseCtx := ctxtrace.WithCorrelationID(context.Background(), corrid)

	eb.mu.RLock()
	subs := make([]subscriber, len(eb.subscribers[eventType]))
	copy(subs, eb.subscribers[eventType])
	eb.mu.RUnlock()

	if len(subs) == 0 {
		eb.log.WithCtx(baseCtx).Debug(
			"event published but no subscribers found",
			zap.String("event", string(eventType)),
		)
		return
	}

	eb.log.WithCtx(baseCtx).Info(
		"publishing event",
		zap.String("event", string(eventType)),
		zap.Int("handlers_count", len(subs)),
	)

	for _, sub := range subs {
		go eb.executeHandler(baseCtx, string(eventType), sub.handler, payload)
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

func (eb *EventBus) generateHash() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand read: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
