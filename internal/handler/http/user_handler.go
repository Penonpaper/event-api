package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penonpaper/event-api/internal/domain"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

type SignUpInput struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"secret123" binding:"required,min=6,max=32"`
	// Ограничиваем список ролей через oneof, чтобы пользователь не мог зарегистрироваться как admin
	Role     string `json:"role" example:"client" binding:"required,oneof=client manager"`
	Nickname string `json:"nickname" example:"john_doe" binding:"required,min=3,max=20"`
}

// @Summary Регистрация пользователя
// @Description Стандартная регистрация пользователя в бд
// @Tags SignUp
// @Accept json
// @Produce json
// @Param input body SignUpInput true "Поле регистрации"
// @Success 200 {object} map[string]string "Возрващаем токены"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 409 {object} map[string]string "Такой эмаил уже есть"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"ы
// @Router /api/v1/auth/signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var input SignUpInput

	if err := c.ShouldBind(&input); err != nil {
		slog.Warn("Invalid signup format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса",
		})
		return
	}

	err := h.userService.SignUp(c.Request.Context(), input.Email, input.Password, input.Role, input.Nickname)

	if err != nil {
		slog.Warn("Failed signup attempt",
			"email", input.Email,
			"error", err,
		)

		switch {
		case errors.Is(err, domain.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{
				"error": "Пользователь с таким email уже существует",
			})
			return
		case errors.Is(err, domain.ErrDataBase):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка сервера, попробуйте позже",
			})
			return
		case errors.Is(err, domain.ErrInternal):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		default:
			// Неожиданная ошибка
			slog.Error("Unexpected error during signup",
				"email", input.Email,
				"error", err,
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}
	}
	slog.Info("Successful signup",
		"email", input.Email,
		"role", input.Role,
	)
	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

// func (h *UserHandler) SignIn(c *gin.Context) {
// 	var input signInInput

// 	if err := c.ShouldBind(&input); err != nil {
// 		slog.Warn("Invalid signin format", "error", err)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Неверный формат запроса",
// 		})
// 		return
// 	}

// 	tokenstr, err := h.userService.SignIn(c.Request.Context(), input.Email, input.Password)
// 	if err != nil {
// 		slog.Warn("Failed signin attempt",
// 			"email", input.Email,
// 			"error", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "Не удалось авторизовать пользователя",
// 		})
// 		return
// 	}

// 	slog.Info("successful SignIn attemp",
// 		"email", input.Email)
// 	c.JSON(http.StatusOK, gin.H{
// 		"token": tokenstr,
// 	})

// }
