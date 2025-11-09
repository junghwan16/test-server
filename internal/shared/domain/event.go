package domain

import (
	"time"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}

// BaseEvent provides common event fields
type BaseEvent struct {
	occurredAt time.Time
}

// NewBaseEvent creates a new base event
func NewBaseEvent() BaseEvent {
	return BaseEvent{occurredAt: time.Now()}
}

// OccurredAt returns when the event occurred
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}
