package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type FakeUserRepo struct {
	Users map[string]*domain.User
}

func NewFakeUserRepo() *FakeUserRepo {
	return &FakeUserRepo{
		Users: make(map[string]*domain.User),
	}
}

func (f *FakeUserRepo) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (f *FakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Пусть возвращает не пустой err
	user, exists := f.Users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil

}

func (f *FakeUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func TestSignIn(t *testing.T) {
	ctx := context.Background()

	repo := NewFakeUserRepo()

	tests := []struct {
		Name       string
		Id         string
		Email      string
		Password   string
		Role       string
		Created_at time.Time
		Nickname   string
	}{
		{
			Name:       "Success attempt",
			Id:         uuid.NewString(),
			Email:      "test@example.com",
			Password:   "testpassword",
			Role:       "client",
			Created_at: time.Now(),
			Nickname:   "testus",
		},
		{
			Name:       "The password is longer than 72 bytes(optional)",
			Id:         uuid.NewString(),
			Email:      "test2@example.com",
			Password:   "testpasswordtestpasswordtestpasswordtestpasswordtestpasswordtestpasswordtestpassword"[:15],
			Role:       "client",
			Created_at: time.Now(),
			Nickname:   "albon",
		},
	}
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	authCase := NewAuthUseCase(repo, "test-access-secret", "test-refresh-secret", 15*time.Minute, 30*time.Hour, rdb)

	for _, tt := range tests {
		hash, err := bcrypt.GenerateFromPassword([]byte(tt.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("bcrypt failed")
		}
		repo.Users[tt.Email] = &domain.User{
			ID:           tt.Id,
			Email:        tt.Email,
			PasswordHash: string(hash),
			Role:         tt.Role,
			Created_at:   time.Now(),
			Nickname:     tt.Nickname,
		}
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := authCase.SignIn(ctx, tt.Email, tt.Password)
			if err != nil {
				t.Fatalf("SignIn %s failed: %v", tt.Email, err)
			}

		})
	}

	// for email, user := range repo.Users {
	// 	_, err = authCase.SignIn(ctx, email)
	// 	if err != nil {
	// 		t.Fatalf("Expected no error, got %v", err)
	// 	}
	// }

}
