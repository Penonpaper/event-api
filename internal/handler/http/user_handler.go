package http

import (
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

type signUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса" + err.Error(),
		})
		return
	}
	if err := h.userService.SignUp(c.Request.Context(), input.Email, input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось зарегистрировать пользователя" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}
func (h *UserHandler) SignIn(c *gin.Context) {
	var input signInInput

	if err := c.ShouldBind(&input); err != nil {
		slog.Warn("Invalid signin format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса",
		})
		return
	}

	tokenstr, err := h.userService.SignIn(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		slog.Warn("Failed signin attempt",
			"email", input.Email,
			"error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Не удалось авторизовать пользователя",
		})
		return
	}

	slog.Info("successful SignIn attemp",
		"email", input.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenstr,
	})

}
