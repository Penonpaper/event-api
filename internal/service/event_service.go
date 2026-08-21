package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
)

type EventUseCase struct {
	EventRepo domain.EventRepository
}

func NewEventService(repo domain.EventRepository) domain.EventService {
	return &EventUseCase{
		EventRepo: repo,
	}
}

func (Ev *EventUseCase) Create(ctx context.Context, userID string, title, desctiption string, location string,
	total_seats int, startAt, endsAt time.Time, status string) error {

	if total_seats <= 0 {
		return fmt.Errorf("Seats must be positive")
	}

	if title == " " {
		return fmt.Errorf("Title cant be empty")
	}

	Event := &domain.Event{
		ID:          uuid.New().String(),
		OrganizerID: userID,
		Title:       title,
		Description: desctiption,
		Location:    location,
		Total_seats: total_seats,
		StartAt:     startAt,
		EndsAt:      endsAt,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	if err := Ev.EventRepo.CreateEvent(ctx, Event); err != nil {
		return fmt.Errorf("Failed to create event (service): %w", err)
	}

	return nil
}

func (Ev *EventUseCase) GetEvents(ctx context.Context) ([]domain.EventListItem, error) {

	var events []domain.EventListItem
	events, err := Ev.EventRepo.GetAllEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get events(service): %w", err)
	}

	return events, nil

}
func (Ev *EventUseCase) Delete(ctx context.Context, id, organizer_id uuid.UUID, role string) error {

	if role == "manager" {
		err := Ev.EventRepo.DeleteEvent(ctx, id, organizer_id)
		if err != nil {
			return fmt.Errorf("Failed to delete from events (service): %w", err)
		}
	} else {
		err := Ev.EventRepo.DeleteEventMN(ctx, id)
		if err != nil {
			return fmt.Errorf("Failed to delete from events (service): %w", err)
		}
	}

	return nil
}
