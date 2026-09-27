package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
)

type AttendeeHandler struct {
	AttendeeHandler domain.EventAttendeesService
}

func NewAttendeeHandler(service domain.EventAttendeesService) *AttendeeHandler {
	return &AttendeeHandler{
		AttendeeHandler: service,
	}
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Неверный ввод"`
	Message string `json:"message" example:"Название продукта обязательно"`
}

// @Summary Регистрация на мероприятие
// @Description Регистрация на мероприятие
// @Security BearerAuth
// @Tags Register_attendee
// @Accept json
// @Produce json
// @Param        id       path      string  true  "UUID мероприятия"
// @Success      200      {object}  map[string]string	"Успешная регистрация на мероприятие"
// @Failure      400      {object}  map[string]string "Неверный формат запроса"
// @Failure		 500	 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/event/{id}/register [post]
func (h *AttendeeHandler) Register(c *gin.Context) {
	eventID := c.Param("id")
	fmt.Println(eventID)
	id, err := uuid.Parse(eventID)
	fmt.Println(id)
	if err != nil {
		slog.Error("Invalid register format ", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка формата данных",
		})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		slog.Error("userID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userID отсутствует",
		})
		return
	}
	stringUserID, ok := userID.(string)
	if !ok {
		slog.Error("invalid userID type in context",
			"type", fmt.Sprintf("%T", userID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный id пользователя",
		})
		return
	}
	uuidUserID, err := uuid.Parse(stringUserID)
	if err != nil {
		slog.Error("Failed to parse string: %w", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Некорректный id пользователя",
		})
		return
	}
	err = h.AttendeeHandler.RegisterAttendee(c.Request.Context(), id, uuidUserID)
	if err != nil {
		slog.Error("Failed to register attendee",
			slog.Any("attendee_id", uuidUserID),
			slog.String("error", err.Error()),
		)
		switch {
		case errors.Is(err, domain.ErrDataBase):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		default:
			slog.Error("Unexpected error during register")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}
	}

	slog.Info("Успешная регистрация пользователя",
		"userID", uuidUserID,
		"eventID", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Успешная регистрация пользователя",
	})

}

// @Summary Отмена регистрации на мероприятие
// @Description Отмена регистрации на мероприятие
// @Security BearerAuth
// @Tags UnRegister_attendee
// @Accept json
// @Produce json
// @Param        id       path      string  true  "UUID мероприятия"
// @Success      200      {object}  map[string]string       "Отмена участия в мероприятии"
// @Failure      400      {object}  map[string]string	"Ошибка формата данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/event/{id}/unregister [post]
func (h *AttendeeHandler) Unregister(c *gin.Context) {
	id := c.Param("id")

	eventID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Invalid unregister format ", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ошибка формата данных",
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		slog.Error("userID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userID отсутствует",
		})
		return
	}
	stringUserID, ok := userID.(string)
	if !ok {
		slog.Error("invalid userID type in context",
			"type", fmt.Sprintf("%T", userID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный id пользователя",
		})
		return
	}
	uuidUserID, err := uuid.Parse(stringUserID)
	if err != nil {
		slog.Error("Failed to parse string: %w", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Некорректный id пользователя",
		})
		return
	}
	err = h.AttendeeHandler.UnregisterAttendee(c.Request.Context(), eventID, uuidUserID)
	if err != nil {
		slog.Error("Failed to unregister attendee",
			slog.Any("attendee_id", uuidUserID),
			slog.String("error", err.Error()),
		)
		switch {
		case errors.Is(err, domain.ErrDataBase):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка",
			})
			return
		default:
			slog.Error("Unexpected error during unregister")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}
	}
	slog.Info("Пользователь успешно отменил регистрацию",
		"userID", uuidUserID,
		"eventID", eventID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Успешная отмена регистрации на мероприятие",
	})
}
