package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/penonpaper/event-api/config"
	"github.com/penonpaper/event-api/internal/app"
)

func main() {
	// Загрузка конфигурации
	config, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Критическая ошибка конфигурации: %v", err))
	}

	// Инициализация логгера
	logger := setupLogger(config.App.AppEnv)
	slog.SetDefault(logger)

	slog.Info(
		"Конфигурация успешно загружена",
		"port", config.App.Port,
		"db_host", config.Postgres.Host,
	)
	// Создание экземпляра приложения
	application := app.NewApp(config)

	defer application.Close()

	if err := application.Run(); err != nil {
		slog.Error("Ошибка при запуске приложения", "error", err)
		os.Exit(1)
	}

}

func setupLogger(env string) *slog.Logger {
	switch env {
	case "local":
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
		return logger
	case "prod":
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		slog.SetDefault(logger)
		return logger

	default:

		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		slog.SetDefault(logger)
		return logger
	}

}
