package eventbus

import (
	"crypto/rand"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type EventType string

type Handler func(
	corrid string,
	payload any,
)

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]Handler
	log         *zap.Logger
}

func New(log *zap.Logger) *EventBus {
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

	eb.mu.RLock()
	handlers, ok := eb.subscribers[eventType]
	eb.mu.RUnlock()

	if !ok {
		eb.log.Debug(
			"no subscribers for event",
			zap.String("event", string(eventType)),
			zap.String("corrid", corrid),
		)
		return
	}

	eb.log.Info(
		"publishing event",
		zap.String("event", string(eventType)),
		zap.String("corrid", corrid),
		zap.Int("handlers_count", len(handlers)),
	)

	for _, handler := range handlers {
		go handler(corrid, payload)
	}
}

func (eb *EventBus) generateHash() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
