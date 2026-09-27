package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          string    `json:"id" example:"a4f21d3e-90ab-4cde-8f12-34567890abcdef"`
	OrganizerID string    `json:"organizer_id" example:"d3b07384-d113-4956-a5db-e13c14c48cf2"`
	Title       string    `json:"title" example:"Конференция по Go" binding:"required"`
	Description string    `json:"description" example:"Ежегодная встреча Go-разработчиков"`
	Location    string    `json:"location" example:"Москва, Технопарк"`
	StartAt     time.Time `json:"start_at" example:"2026-09-10T10:00:00Z" binding:"required"`
	EndsAt      time.Time `json:"ends_at" example:"2026-09-10T18:00:00Z" binding:"required"`
	Status      string    `json:"status" example:"planned"`
	TotalSeats  int       `json:"total_seats" example:"150" binding:"required"`
	CreatedAt   time.Time `json:"created_at" example:"2026-08-29T15:04:05Z"`
}

type EventListItem struct {
	ID            string    `json:"id" example:"a4f21d3e-90ab-4cde-8f12-34567890abcdef"`
	Title         string    `json:"title" example:"Конференция по Go"`
	Description   string    `json:"description" example:"Ежегодная встреча Go-разработчиков"`
	Location      string    `json:"location" example:"Москва, Технопарк"`
	StartAt       time.Time `json:"start_at" example:"2026-09-10T10:00:00Z"`
	EndsAt        time.Time `json:"ends_at" example:"2026-09-10T18:00:00Z"`
	OrganizerID   string    `json:"organizer_id" example:"d3b07384-d113-4956-a5db-e13c14c48cf2"`
	OrganizerName string    `json:"organizer_name" example:"Иван Иванов"`
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
