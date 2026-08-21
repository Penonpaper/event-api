package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          string    `json:"id"`
	OrganizerID string    `json:"organizer_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartAt     time.Time `json:"start_at"`
	EndsAt      time.Time `json:"ends_at"`
	Status      string    `json:"status"`
	Total_seats int       `json:"total_seats"`
	CreatedAt   time.Time `json:"created_at"`
}
type EventListItem struct {
	ID            string
	Title         string
	Description   string
	Location      string
	StartAt       time.Time
	EndsAt        time.Time
	OrganizerID   string
	OrganizerName string
}

type EventRepository interface {
	CreateEvent(ctx context.Context, Event *Event) error
	GetAllEvents(ctx context.Context) ([]EventListItem, error)
	DeleteEvent(ctx context.Context, id, organizer_id uuid.UUID) error
	DeleteEventMN(ctx context.Context, id uuid.UUID) error // Удаление админом
}

type EventService interface {
	Create(ctx context.Context, userID string, title, desctiption string, location string,
		total_seats int, startAt, endsAt time.Time, status string) error
	GetEvents(ctx context.Context) ([]EventListItem, error)
	Delete(ctx context.Context, id, organizer_id uuid.UUID, role string) error
}
