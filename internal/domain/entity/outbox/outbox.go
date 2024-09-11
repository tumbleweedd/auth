package outbox

import (
	"github.com/google/uuid"

	"time"
)

type EventType string

const (
	AggregateTypeUser    EventType = "User"
	EventTypeUserCreated EventType = "UserCreated"
	EventTypeUserDeleted EventType = "UserDeleted"
)

type Event struct {
	ID        uuid.UUID `json:"id"`
	Type      EventType `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
	Payload   []byte    `json:"payload"`
}

func NewOutbox(id uuid.UUID, eventType EventType, payload []byte) *Event {
	return &Event{
		ID:        id,
		Type:      eventType,
		CreatedAt: time.Now(),
		Payload:   payload,
	}
}
