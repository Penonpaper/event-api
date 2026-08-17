package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/penonpaper/event-api/internal/domain"
)

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) domain.EventRepository {
	return &EventRepo{
		db: db,
	}
}

func (Ev *EventRepo) CreateEvent(ctx context.Context, Event *domain.Event) error {
	query := `INSERT INTO events(id, title, description, total_seats, created_at)
			VALUES ($1, $2, $3, $4, $5)		
	`

	_, err := Ev.db.Exec(ctx, query, Event.ID, Event.Title, Event.Description, Event.Total_seats, Event.CreatedAt)

	if err != nil {
		return fmt.Errorf("Failed to create event (repo): %w", err)
	}

	return nil
}

func (Ev *EventRepo) GetAllEvents(ctx context.Context) ([]domain.Event, error) {
	query := `SELECT id, title, description, total_seats, created_at FROM events`

	rows, err := Ev.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Failed to get events: %w", err)
	}

	defer rows.Close()

	var events []domain.Event

	for rows.Next() {
		var event domain.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Total_seats,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row (get all event): %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return events, nil

}
