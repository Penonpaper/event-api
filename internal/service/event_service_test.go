package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
)

type FakeEventRepository struct {
	CreateErr  error
	CreateCall bool
}

func (f *FakeEventRepository) CreateEvent(ctx context.Context, Event *domain.Event) error {

	f.CreateCall = true
	return f.CreateErr
}

func (f *FakeEventRepository) GetAllEvents(ctx context.Context) ([]domain.EventListItem, error) {
	return nil, nil
}

func (f *FakeEventRepository) DeleteEvent(ctx context.Context, id uuid.UUID, organizer_id uuid.UUID) error {
	return nil
}
func (f *FakeEventRepository) DeleteEventMN(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (f *FakeEventRepository) GetByTags(ctx context.Context, filt domain.Filter) ([]domain.EventListItem, error) {
	return nil, nil
}
func TestCreateEvent_Success(t *testing.T) {
	tests := []struct {
		Name        string
		CreateErr   error
		ExpectedErr error
	}{
		{
			Name:        "Валидные данные",
			CreateErr:   nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Ошибка в repository",
			CreateErr:   domain.ErrDataBase,
			ExpectedErr: domain.ErrDataBase,
		},
		// {
		// 	Name:        "Ошибки в credentials",
		// 	CreateErr:   domain.ErrConditions,
		// 	ExpectedErr: domain.ErrConditions,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			FakeEventRepository := FakeEventRepository{
				CreateErr: tt.CreateErr,
			}
			usecase := NewEventService(&FakeEventRepository)
			now := time.Now()
			err := usecase.Create(context.Background(),
				"f123",
				"Go meetup",
				"Обсуждаем goшку на среднем уровне",
				"telemost/7667",
				41,
				now,
				now.Add(4*time.Hour),
				"planned",
			)
			if tt.ExpectedErr == nil {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			}
			if !errors.Is(err, tt.ExpectedErr) {
				t.Fatalf("Expected %v, got %v", tt.ExpectedErr, err)
			}
			if tt.CreateErr == nil && !FakeEventRepository.CreateCall {
				t.Fatalf("expected repoitory Create to be called")
			}

		})
	}
}
