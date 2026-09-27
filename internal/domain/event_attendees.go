package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventAttendees struct {
	EventID uuid.UUID `json:"event_id" example:"a4f21d3e-90ab-4cde-8f12-34567890abcdef" binding:"required"`
	UserID  uuid.UUID `json:"user_id" example:"d3b07384-d113-4956-a5db-e13c14c48cf2" binding:"required"`
}

type Attendee struct {
	ID       uuid.UUID `json:"id" example:"d3b07384-d113-4956-a5db-e13c14c48cf2"`
	Email    string    `json:"email" example:"attendee@example.com" binding:"required,email"`
	Nickname string    `json:"nickname" example:"go_developer"`
}

type Events struct {
	ID       uuid.UUID `json:"id" example:"a4f21d3e-90ab-4cde-8f12-34567890abcdef"`
	Title    string    `json:"title" example:"Интенсив по базам данных" binding:"required"`
	Location string    `json:"location" example:"Казань, ИТ-Парк"`
	StartAt  time.Time `json:"start_at" example:"2026-09-15T12:00:00Z" binding:"required"`
}

type EventAttendeesRepository interface {
	Register(ctx context.Context, eventID, userID uuid.UUID) error
	Unregister(ctx context.Context, eventID, userID uuid.UUID) error

	//IsRegistered(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (bool, error)
	//GetEventAttendees(ctx context.Context, eventID uuid.UUID) ([]Attendee, error)

	//GetUserEvents(ctx context.Context, userID uuid.UUID) ([]Events, error)

	//RemoveAttendee(ctx context.Context, userID uuid.UUID) error

}

type EventAttendeesService interface {
	RegisterAttendee(ctx context.Context, eventID, userID uuid.UUID) error
	UnregisterAttendee(ctx context.Context, eventID, userID uuid.UUID) error
}
