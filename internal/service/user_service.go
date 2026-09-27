package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsercase(repo domain.UserRepository) domain.UserService {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (a *UserUsecase) SignUp(ctx context.Context, email, password, role, nickname string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w: failed to hash password: %v", domain.ErrInvalidPassword, err)
		//return fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Created_at:   time.Now(),
		Nickname:     nickname,
	}

	if err := a.userRepo.Create(ctx, &user); err != nil {

		switch {
		case errors.Is(err, domain.ErrDuplicateEmail):
			return fmt.Errorf("%w: Email already exists: %v", domain.ErrDuplicateEmail, err)
		case errors.Is(err, domain.ErrDataBase):
			return fmt.Errorf("%w: %v", domain.ErrDataBase, err)

		default:
			return fmt.Errorf("%w: failed to create user: %v", domain.ErrInternal, err)
		}

		//return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// func (a *UserUsecase) SignIn(ctx context.Context, email, password string) (*domain.TokenPair, error) {
// 	// Сравнение пароль из бд
// 	us_db, err := a.userRepo.GetByEmail(ctx, email)

// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to get data from db: %w", err)
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(us_db.PasswordHash), []byte(password))

// 	if err != nil {
// 		return nil, fmt.Errorf("Password do not match: %w", err)
// 	}

// 	accessToken, err := generateAccessToken(us_db.ID, us_db.Role)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to get access token: %w", err)
// 	}

// 	refreshToken, err := generateRefreshToken()
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to get refreshtoken: %w", err)
// 	}

// 	refreshTokenHash := generateHashRefreshToken(refreshToken)
// 	a.userRepo.RefreshTokenRepo()
// }
