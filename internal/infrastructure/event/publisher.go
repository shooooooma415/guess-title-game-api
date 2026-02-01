package event

import (
	"log"
	"sync"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/event"
)

// InMemoryEventPublisher is an in-memory implementation of event.Publisher
type InMemoryEventPublisher struct {
	handlers map[string][]event.EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventPublisher creates a new InMemoryEventPublisher
func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		handlers: make(map[string][]event.EventHandler),
	}
}

// Publish publishes an event to all subscribed handlers
func (p *InMemoryEventPublisher) Publish(evt event.Event) {
	p.mu.RLock()
	handlers, exists := p.handlers[evt.EventType()]
	p.mu.RUnlock()

	if !exists {
		return
	}

	// Execute handlers asynchronously
	for _, handler := range handlers {
		go func(h event.EventHandler, e event.Event) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Event handler panic: %v", r)
				}
			}()
			h(e)
		}(handler, evt)
	}
}

// Subscribe subscribes a handler to a specific event type
func (p *InMemoryEventPublisher) Subscribe(eventType string, handler event.EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers[eventType] = append(p.handlers[eventType], handler)
}
