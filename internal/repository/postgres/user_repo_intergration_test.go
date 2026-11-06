package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/penonpaper/event-api/internal/domain"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := "postgres://test_user:test_password@localhost:5434/event_api_test?sslmode=disable"
	var err error
	testPool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		panic("Не удалось создать пул: %v" + err.Error())
	}

	if err = testPool.Ping(ctx); err != nil {
		panic("Нет коннекта к бд: " + err.Error())
	}

	code := m.Run()
	testPool.Close()
	os.Exit(code)
}

func ClearDB(ctx context.Context, t *testing.T) {
	//marks the calling function as a test helper function
	t.Helper()
	_, err := testPool.Exec(ctx, `TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("Не удалось очистить таблицу users перед тестом: %v", err)
	}

}
func TestUserRepo_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ClearDB(ctx, t)
	userRepo := NewUserRepo(testPool)

	now := time.Now()
	testUser := &domain.User{
		ID:           "d3b07384-d113-4956-a5db-e13c14c49cf2",
		Email:        "testkamimil@mail.com",
		PasswordHash: "test-hash",
		Role:         "client",
		Created_at:   now,
		Nickname:     "kamilka",
	}

	if err := userRepo.Create(ctx, testUser); err != nil {
		t.Fatalf("Не удалось создатьпользователя: %v", err)
	}

}

func TestUserRepo_Duplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	ClearDB(ctx, t)
	userRepo := NewUserRepo(testPool)

	testUser := &domain.User{
		ID:           "d3b07384-d113-4956-a5db-e13c14c49cf3",
		Email:        "testkamil@mail.com",
		PasswordHash: "test1-hash",
		Role:         "client",
		Created_at:   now,
		Nickname:     "randomname",
	}
	if err := userRepo.Create(ctx, testUser); err != nil {
		t.Fatalf("Не удалось создать первого пользователя: %v", err)
	}
	err := userRepo.Create(ctx, testUser)
	if err == nil {
		t.Fatalf("Ожидалась ошибка дубликата, но пользователь был успешно создан")
	}

}
