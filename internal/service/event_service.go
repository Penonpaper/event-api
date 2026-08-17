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

func (Ev *EventUseCase) Create(ctx context.Context, title, desctiption string, total_seats int) error {

	if total_seats <= 0 {
		return fmt.Errorf("Seats must be positive")
	}

	Event := &domain.Event{
		ID:          uuid.New().String(),
		Title:       title,
		Description: desctiption,
		Total_seats: total_seats,
		CreatedAt:   time.Now(),
	}

	if err := Ev.EventRepo.CreateEvent(ctx, Event); err != nil {
		return fmt.Errorf("Failed to create event (service): %w", err)
	}

	return nil
}

func (Ev *EventUseCase) GetEvents(ctx context.Context) ([]domain.Event, error) {

	var events []domain.Event
	events, err := Ev.EventRepo.GetAllEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get events(service): %w", err)
	}

	return events, nil

}
