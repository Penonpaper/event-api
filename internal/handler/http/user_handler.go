package http

import (
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

func (h *UserHandler) SignUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса" + err.Error(),
		})
	}
	if err := h.userService.SignUp(c.Request.Context(), input.Email, input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось зарегистрировать пользователя" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}
