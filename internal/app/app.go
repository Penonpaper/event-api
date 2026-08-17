package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/penonpaper/event-api/config"
	"github.com/penonpaper/event-api/internal/handler/http"
	repository "github.com/penonpaper/event-api/internal/repository/postgres"
	"github.com/penonpaper/event-api/internal/service"
)

type App struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

// DSN PostgreSQL - postgres://user:password@host:port/dbname
func (a *App) Run() error {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		a.cfg.Postgres.User,
		a.cfg.Postgres.Password,
		a.cfg.Postgres.Host,
		a.cfg.Postgres.Port,
		a.cfg.Postgres.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("Критическая ошибка создания пула PostgreSQL", "error", err)
		return fmt.Errorf("failde to create pgxpool: %w", err)
	}
	a.db = pool

	var pingErr error
	for attempts := 1; attempts <= 5; attempts++ {
		pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		pingErr = a.db.Ping(pingCtx)
		cancel()

		if pingErr == nil {
			break
		}

		slog.Warn("Попытка покдлючения к PostgreSQL не удалась",
			"attempt", attempts,
			"max_attempts", 5,
		)
		time.Sleep(2 * time.Second)

	}

	if pingErr != nil {
		slog.Error("База данных не ответила после 5 попыток подключения", "error", err)
		a.db.Close()
		return fmt.Errorf("failed to connect db: %w", pingErr)
	}

	slog.Info("Успешное подключение к PostgreSQL!",
		"host", a.cfg.Postgres.Host,
		"port", a.cfg.Postgres.Port,
	)

	userRepos := repository.NewRepositories(a.db)
	userSers := service.NewServices(*userRepos, a.cfg.JWT.Secret, a.cfg.JWT.TTLMinutes)
	userHandlers := http.NewHandlers(*userSers)
	router := gin.Default()

	routes := http.NewRoutes(userHandlers.User, userHandlers.Event, a.cfg.JWT.Secret)
	routes.RegisteredRoutes(router)
	gin.SetMode(a.cfg.App.GinMode)

	slog.Info("Успешный запуск HTTP сервера: " + a.cfg.App.Port + " " + a.cfg.App.AppEnv)

	if err := router.Run(":" + a.cfg.App.Port); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func (a *App) Close() {
	if a.db != nil {
		slog.Info("Закрытие пула соединений")
		a.db.Close()
	}
}
