package service

import (
	"context"
	"fmt"

	"github.com/penonpaper/event-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsercase(repo domain.UserRepository) domain.UserService {
	return &UserUsecase{userRepo: repo}

}

func (a *UserUsecase) SignUp(ctx context.Context, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "client",
	}

	if err := a.userRepo.Create(ctx, &user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
