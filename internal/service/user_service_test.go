package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
)

type FakeUserRepository struct {
	CreateErr  error
	CreateCall bool
}

func (f *FakeUserRepository) Create(ctx context.Context, user *domain.User) error {
	f.CreateCall = true
	return f.CreateErr
}
func (f *FakeUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (a *FakeUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func TestSignUp_Success(t *testing.T) {
	tests := []struct {
		name        string
		createErr   error
		expectedErr error
	}{
		{
			name:        "success",
			createErr:   nil,
			expectedErr: nil,
		},
		{
			name:        "duplicate email",
			createErr:   domain.ErrDuplicateEmail,
			expectedErr: domain.ErrDuplicateEmail,
		},
		{
			name:        "database error",
			createErr:   domain.ErrDataBase,
			expectedErr: domain.ErrDataBase,
		},
		{
			name:        "unexpected error",
			createErr:   errors.New("Something went wrong"),
			expectedErr: domain.ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeRepo := FakeUserRepository{
				CreateErr: test.createErr,
			}
			usecase := NewUserUsercase(&fakeRepo)

			err := usecase.SignUp(context.Background(),
				"test@example.com",
				"testpass",
				"client",
				"TestUser",
			)

			if test.expectedErr == nil {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				return
			}

			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("expected error %v, got: %v",
					test.expectedErr, err)
			}

			if test.createErr == nil && !fakeRepo.CreateCall {
				t.Fatalf("expected repoitory Create to be called")
			}

		})
	}
}

func TestSignUp_CorrectPass(t *testing.T) {
	fakeRepo := FakeUserRepository{}
	usecase := NewUserUsercase(&fakeRepo)

	err := usecase.SignUp(context.Background(),
		"test123@mail.com",
		"shortshortshortshortshortshortshortshortshortshortshortshortshortshortshortshortshortshortshort",
		"client",
		"test_us77",
	)
	if !errors.Is(err, domain.ErrInvalidPassword) {
		t.Fatalf("expected %s, got %v", domain.ErrInvalidPassword, err)
	}

	if fakeRepo.CreateCall {
		t.Fatalf("Repository cannot be called when we test invalid password")
	}
}
