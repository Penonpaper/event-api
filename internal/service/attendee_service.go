package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
)

type AttendeeUseCase struct {
	AttendeeRepo domain.EventAttendeesRepository
}

func NewAttendeeUseCase(repo domain.EventAttendeesRepository) domain.EventAttendeesService {
	return &AttendeeUseCase{
		AttendeeRepo: repo,
	}
}

func (A *AttendeeUseCase) RegisterAttendee(ctx context.Context, eventID, userID uuid.UUID) error {

	if err := A.AttendeeRepo.Register(ctx, eventID, userID); err != nil {
		return fmt.Errorf("%w: Failed method: register(repo): %v", domain.ErrDataBase, err)
	}
	return nil
}

func (A *AttendeeUseCase) UnregisterAttendee(ctx context.Context, eventID, userID uuid.UUID) error {
	if err := A.AttendeeRepo.Unregister(ctx, eventID, userID); err != nil {
		return fmt.Errorf("%w: Failed method: unregister(repo): %v", domain.ErrDataBase, err)
	}
	return nil
}
