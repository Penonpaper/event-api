package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
)

type EventHandler struct {
	EventService domain.EventService
}
type EventInput struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	StartAt     time.Time `json:"starts_at" binding:"required"`
	EndsAt      time.Time `json:"ends_at" binding:"required"`
	TotalSeats  int       `json:"total_seats" binding:"required"`
	Status      string    `json:"status" binding:"required"`
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
	userID, exists := c.Get("userID")
	stringUserID, ok := userID.(string)
	if !ok {
		slog.Warn("Invalid to convert userID to string", "error", "Convert Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка конвертации ID юзера",
		})
		return
	}

	if !exists {
		slog.Warn("Field userID is empty", "error", "Empty field")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось обнаружить ID пользователя",
		})
	}
	err := Ev.EventService.Create(c.Request.Context(), stringUserID, eventinput.Title, eventinput.Description, eventinput.Location, eventinput.TotalSeats, eventinput.StartAt, eventinput.EndsAt, eventinput.Status)
	if err != nil {
		slog.Warn("Invalid to create event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось создать ивент",
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
	var events []domain.EventListItem

	events, err := Ev.EventService.GetEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось получить ивенты",
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (Ev *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	id, err := uuid.Parse(eventID)
	if err != nil {
		slog.Error("Invalid delete event format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось выгрузить данные",
		})
		return
	}
	userRole, exists := c.Get("userRole")
	if !exists {
		slog.Error("role not found in context ", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось определить роль пользователя",
		})
		return
	}
	role, ok := userRole.(string)
	if !ok {
		slog.Error("invalid role type in context",
			"type", fmt.Sprintf("%T", role),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Некорректная роль пользователя",
		})
		return
	}

	orgIDValue, exists := c.Get("userID")
	if !exists {
		slog.Error("userid not found in context ", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось определить пользователя",
		})
		return
	}
	orgID, ok := orgIDValue.(string)
	if !ok {
		slog.Error("invalid userID type in context",
			"type", fmt.Sprintf("%T", orgIDValue))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Некорректный id пользователя",
		})
		return
	}

	orgidUUID, err := uuid.Parse(orgID)

	if err != nil {
		slog.Error("Invalid to parse ID to uuid", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось преобразовать данные",
		})
		return
	}
	err = Ev.EventService.Delete(c.Request.Context(), id, orgidUUID, role)
	if err != nil {
		slog.Warn("Invalid to delete event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось удалить ивент",
		})
		return
	}
	slog.Info("Успешное удаление ивента",
		"title", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Успешное удаление ивента",
	})

}
