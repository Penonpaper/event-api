package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/penonpaper/event-api/internal/domain"
)

type Repositories struct {
	User  domain.UserRepository
	Event domain.EventRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:  NewUserRepo(db),
		Event: NewEventRepo(db),
	}
}
