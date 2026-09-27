package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `json:"id" example:"d3b07384-d113-4956-a5db-e13c14c48cf2"`
	Email        string    `json:"email" example:"user@example.com" binding:"required,email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role" example:"client"`
	Created_at   time.Time `json:"created_at" example:"2026-08-29T15:04:05Z"`
	Nickname     string    `json:"nickname" example:"john_doe"`
}

type UserService interface {
	SignUp(ctx context.Context, email, password, role, nickname string) error
	//SignIn(ctx context.Context, email, password string) (*TokenPair, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}
