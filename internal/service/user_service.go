package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/penonpaper/event-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo      domain.UserRepository
	jwtsecret     string
	jwtttlminutes int
}

func NewUserUsercase(repo domain.UserRepository, jwtsecret string, jwtttlminutes int) domain.UserService {
	return &UserUsecase{
		userRepo:      repo,
		jwtsecret:     jwtsecret,
		jwtttlminutes: jwtttlminutes,
	}
}

func (a *UserUsecase) SignUp(ctx context.Context, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "client",
		Created_at:   time.Now(),
	}

	if err := a.userRepo.Create(ctx, &user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (a *UserUsecase) SignIn(ctx context.Context, email, password string) (string, error) {
	// Сравнение пароль из бд
	us_db, err := a.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return " ", fmt.Errorf("Failed to get data from db: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(us_db.PasswordHash), []byte(password))

	if err != nil {
		return " ", fmt.Errorf("Password do not match: %w", err)
	}

	expirationTime := time.Now().Add(time.Duration(a.jwtttlminutes) * time.Minute)

	claims := &domain.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   us_db.ID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: us_db.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString([]byte(a.jwtsecret))
	if err != nil {
		return " ", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenstring, nil
}
