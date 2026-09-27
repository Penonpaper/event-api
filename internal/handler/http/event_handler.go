package http

import (
	"errors"
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
	Title       string    `json:"title" example:"Воркшоп по архитектуре Go" binding:"required,min=3,max=100"`
	Description string    `json:"description" example:"Практическое занятие по проектированию чистого кода." binding:"required,max=1000"`
	Location    string    `json:"location" example:"Казань, ул. Пушкина, д. 5" binding:"required"`
	StartAt     time.Time `json:"starts_at" example:"2026-10-15T14:00:00Z" binding:"required"`
	EndsAt      time.Time `json:"ends_at" example:"2026-10-15T18:00:00Z" binding:"required"`
	TotalSeats  int       `json:"total_seats" example:"45" binding:"required"`
	Status      string    `json:"status" example:"planned" binding:"required,oneof=planned active cancelled"`
}

func NewEventHandler(service domain.EventService) *EventHandler {
	return &EventHandler{
		EventService: service,
	}
}

// @Summary Создание ивента определенной ролью
// @Description Cоздание мероприятия админом или организатором
// @Security BearerAuth
// @Tags Event
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Успешное создание ивента"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Param input body EventInput true "Описание мепроприятия"
// @Router /api/v1/event [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {

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
	err := h.EventService.Create(c.Request.Context(), stringUserID, eventinput.Title, eventinput.Description, eventinput.Location, eventinput.TotalSeats, eventinput.StartAt, eventinput.EndsAt, eventinput.Status)
	if err != nil {
		slog.Warn("Invalid to create event", "error", err)

		switch {
		case errors.Is(err, domain.ErrDataBase):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка сервера",
			})
			return
		case errors.Is(err, domain.ErrConditions):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Ошибка формата данных",
			})
			return
		default:
			slog.Error("Unexpected server error during createEvent", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return

		}
	}

	slog.Info("Успешное создание ивента",
		"title", eventinput.Title)
	c.JSON(http.StatusOK, gin.H{
		"message": "Успешное создание ивента",
	})

}

// @Summary Получение всех ивентов
// @Description Получение всех мероприятий
// @Tags Event
// @Accept json
// @Produce json
// @Success 200 {array} domain.EventListItem "Список мероприятий успешно получен"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/events [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	var events []domain.EventListItem

	events, err := h.EventService.GetEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось получить ивенты",
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

// @Summary Удаление ивента
// @Description Удаление определнного мероприятия
// @Param id path string true "ID события"
// @Security BearerAuth
// @Tags Event
// @Success 200 {object} map[string]string "Успешное удаление event"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Router /api/v1/event/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	id, err := uuid.Parse(eventID)
	fmt.Println(id)
	if err != nil {
		slog.Error("Invalid delete event format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось выгрузить данные",
		})
		return
	}
	userRole, exists := c.Get("userRole")
	if !exists {
		slog.Error("role not found in context")
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
	err = h.EventService.Delete(c.Request.Context(), id, orgidUUID, role)
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
