package domain

// EventBus publishes domain events
type EventBus interface {
	Publish(event DomainEvent) error
}

// SimpleEventBus is a simple in-memory event bus
type SimpleEventBus struct {
	handlers []EventHandler
}

// EventHandler handles domain events
type EventHandler func(event DomainEvent) error

// NewSimpleEventBus creates a new simple event bus
func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make([]EventHandler, 0),
	}
}

// Subscribe adds an event handler
func (b *SimpleEventBus) Subscribe(handler EventHandler) {
	b.handlers = append(b.handlers, handler)
}

// Publish publishes an event to all handlers
func (b *SimpleEventBus) Publish(event DomainEvent) error {
	for _, handler := range b.handlers {
		if err := handler(event); err != nil {
			// Log error but continue
			continue
		}
	}
	return nil
}
