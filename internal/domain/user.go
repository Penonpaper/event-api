package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Created_at   time.Time `json:"created_at"`
	Nickname     string    `json:"nickname"`
}
type TokenClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type UserService interface {
	SignUp(ctx context.Context, email, password, role, nickname string) error
	SignIn(ctx context.Context, email, password string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}
