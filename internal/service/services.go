package service

import (
	"github.com/penonpaper/event-api/internal/domain"
	"github.com/penonpaper/event-api/internal/repository/postgres"
)

type Services struct {
	User  domain.UserService
	Event domain.EventService
}

func NewServices(repos postgres.Repositories, jwtSecret string, ttl int) *Services {
	return &Services{
		User:  NewUserUsercase(repos.User, jwtSecret, ttl),
		Event: NewEventService(repos.Event),
	}
}
