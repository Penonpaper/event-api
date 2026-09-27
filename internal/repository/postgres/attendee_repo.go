package postgres

import (
	"context"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/penonpaper/event-api/internal/domain"
)

type AttendeeRepo struct {
	db *pgxpool.Pool
}

func NewAttendeeRepo(db *pgxpool.Pool) domain.EventAttendeesRepository {
	return &AttendeeRepo{
		db: db,
	}
}

func (r *AttendeeRepo) Register(ctx context.Context, eventID, userID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var totalSeats int32

	row := tx.QueryRow(ctx, `
		SELECT total_seats
		FROM events
		WHERE id = $1
		FOR UPDATE
	`, eventID).Scan(&totalSeats)

	if row != nil {
		return fmt.Errorf("Failed to get total_seats: %w", row)
	}

	var attendeescount int32

	row = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM event_attendees
		WHERE event_id = $1
	`, eventID).Scan(&attendeescount)

	fmt.Printf("АУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУУ : %d\n", attendeescount)
	if row != nil {
		return fmt.Errorf("Failed to get attendees count: %w", row)
	}

	if attendeescount >= totalSeats {
		return fmt.Errorf("There are no more seats left: ")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO event_attendees(event_id, user_id) VALUES ($1, $2)
	`, eventID, userID)
	if err != nil {
		return fmt.Errorf("Failed to insert attendee: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (r *AttendeeRepo) Unregister(ctx context.Context, eventID, userID uuid.UUID) error {

	query := `DELETE FROM event_attendees WHERE event_id = $1 AND user_id = $2`

	rows, err := r.db.Exec(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("Failed to delete from event_attendees: %w", err)
	}

	if rows.RowsAffected() <= 0 {
		return fmt.Errorf("No rows affected (DELETE form)")
	}
	return nil
}
