package http

import "github.com/penonpaper/event-api/internal/service"

type Handlers struct {
	User     *UserHandler
	Event    *EventHandler
	Attendee *AttendeeHandler
	Auth     *AuthHandler
}

func NewHandlers(services service.Services) *Handlers {
	return &Handlers{
		User:     NewUserHandler(services.User),
		Event:    NewEventHandler(services.Event),
		Attendee: NewAttendeeHandler(services.Attendee),
		Auth:     NewAuthHandler(services.Auth),
	}
}
