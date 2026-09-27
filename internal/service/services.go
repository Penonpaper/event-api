package service

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/penonpaper/event-api/internal/domain"
	"github.com/penonpaper/event-api/internal/repository/postgres"
)

type Services struct {
	User     domain.UserService
	Event    domain.EventService
	Attendee domain.EventAttendeesService
	Auth     domain.AuthService
}

func NewServices(repos postgres.Repositories, accessTokenSecret, refreshTokenSecret string, accessTokenTTL, refreshTokenTTL time.Duration, rdb *redis.Client) *Services {
	return &Services{
		User:     NewUserUsercase(repos.User),
		Event:    NewEventService(repos.Event),
		Attendee: NewAttendeeUseCase(repos.Attendee),
		Auth:     NewAuthUseCase(repos.User, accessTokenSecret, refreshTokenSecret, accessTokenTTL, refreshTokenTTL, rdb),
	}
}
