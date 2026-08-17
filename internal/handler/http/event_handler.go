package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penonpaper/event-api/internal/domain"
)

type EventHandler struct {
	EventService domain.EventService
}
type EventInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	TotalSeats  int    `json:"total_seats" binding:"required"`
}

func NewEventHandler(service domain.EventService) *EventHandler {
	return &EventHandler{
		EventService: service,
	}
}

func (Ev *EventHandler) CreateEvent(c *gin.Context) {

	// Здесь через контекст можно брать и использовать userID авторизированного пользователя
	var eventinput EventInput

	if err := c.ShouldBind(&eventinput); err != nil {
		slog.Warn("Invalid create event format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось выгрузить данные",
		})
		return
	}

	err := Ev.EventService.Create(c.Request.Context(), eventinput.Title, eventinput.Description, eventinput.TotalSeats)
	if err != nil {
		slog.Warn("Invalid to create event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось создать ивент",
		})
		return
	}

	slog.Info("Успешное создание ивента",
		"title", eventinput.Title)
	c.JSON(http.StatusOK, gin.H{
		"message": "Успешное создание ивента",
	})

}

func (Ev *EventHandler) GetEvents(c *gin.Context) {
	var events []domain.Event

	events, err := Ev.EventService.GetEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось получить ивенты",
		})
		return
	}
	c.JSON(http.StatusOK, events)
}
