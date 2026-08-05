package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/penonpaper/event-api/internal/domain"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) domain.UserRepository {
	return &UserRepo{
		db: db,
	}
}

func (a *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password_hash, role)
			VALUES ($1, $2, $3)
			RETURNING id, created_at`

	row := a.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Role)
	if err := row.Scan(&user.ID, &user.Created_at); err != nil {
		return fmt.Errorf("failed to execute insert user query: %w", err)
	}
	return nil
}

func (a *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
