package http

import "github.com/penonpaper/event-api/internal/service"

type Handlers struct {
	User  *UserHandler
	Event *EventHandler
}

func NewHandlers(services service.Services) *Handlers {
	return &Handlers{
		User:  NewUserHandler(services.User),
		Event: NewEventHandler(services.Event),
	}
}
