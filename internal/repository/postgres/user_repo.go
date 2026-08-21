package postgres

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
	query := `INSERT INTO users (id, email, password_hash, role, created_at, nickname)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING created_at`

	_, err := a.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Role, user.Created_at, user.Nickname)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	return nil
}

func (a *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, password_hash, role, created_at FROM users
	WHERE email = $1`

	var user domain.User
	row := a.db.QueryRow(ctx, query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.Created_at); err != nil {
		return nil, fmt.Errorf("failed to method GetbyEmail: %w", err)
	}
	return &user, nil
}
