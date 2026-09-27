package postgres

import (
	"context"
	"fmt"

	uuid "github.com/google/uuid"

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
	query := `INSERT INTO events(id, organizer_id, title, description, location, start_at, 
								ends_at, status, total_seats)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)		
	`

	_, err := Ev.db.Exec(ctx, query, Event.ID, Event.OrganizerID, Event.Title, Event.Description,
		Event.Location, Event.StartAt, Event.EndsAt, Event.Status, Event.TotalSeats)

	if err != nil {
		return fmt.Errorf("Failed to create event (repo): %w", err)
	}

	return nil
}

func (Ev *EventRepo) GetAllEvents(ctx context.Context) ([]domain.EventListItem, error) {
	query := `SELECT e.id,
	e.title,
	e.description,
	e.location,
	e.start_at,
	e.ends_at,
	e.organizer_id,
	u.nickname
	FROM events e
	JOIN users u
	ON u.id = e.organizer_id`
	rows, err := Ev.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Failed to get events: %w", err)
	}

	defer rows.Close()

	var events []domain.EventListItem

	for rows.Next() {
		var event domain.EventListItem
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.StartAt,
			&event.EndsAt,
			&event.OrganizerID,
			&event.OrganizerName,
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
func (Ev *EventRepo) DeleteEvent(ctx context.Context, id uuid.UUID, organizer_id uuid.UUID) error {
	query := `DELETE FROM events
			WHERE id = $1
				AND organizer_id = $2`

	_, err := Ev.db.Exec(ctx, query, id, organizer_id)
	if err != nil {
		return fmt.Errorf("Failed to delete from events: %w", err)
	}
	return nil
}

func (Ev *EventRepo) DeleteEventMN(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events
			WHERE id = $1`

	_, err := Ev.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete from events: %w", err)
	}
	return nil
}
