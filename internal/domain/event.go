package domain

import (
	"context"
	"time"
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Total_seats int       `json:"total_seats"`
	CreatedAt   time.Time `json:"created_at"`
}
type EventRepository interface {
	CreateEvent(ctx context.Context, Event *Event) error
	GetAllEvents(ctx context.Context) ([]Event, error)
}

type EventService interface {
	Create(ctx context.Context, title, desctiption string, total_seats int) error
	GetEvents(ctx context.Context) ([]Event, error)
}
