package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penonpaper/event-api/internal/domain"
)

type AuthHandler struct {
	AuthService domain.AuthService
}

func NewAuthHandler(service domain.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: service,
	}
}

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Авторизация в приложение
// @Description В авторизации происходит генерация двух токенов refresh и access
// @Tags Signin
// @Accept json
// @Produce json
// @Param input body SignInInput true "Данные для входа"
// @Success 200 {object} map[string]string "Успешная авторизация"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Failure 401 {object} map[string]string "Не удалось авторизовать пользователя"
// @Router /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var input SignInInput

	if err := c.ShouldBind(&input); err != nil {
		slog.Warn("Invalid signin format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса",
		})
		return
	}

	tokenpair, err := h.AuthService.SignIn(c.Request.Context(), input.Email, input.Password)

	if err != nil {
		slog.Warn("Failed signin attempt",
			"email", input.Email,
			"error", err)
		switch {
		case errors.Is(err, domain.ErrDataBase):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Не удалось авторизовать пользователя",
			})
			return
		case errors.Is(err, domain.ErrConditions):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Не удалось авторизовать пользователя",
			})
			return

		case errors.Is(err, domain.ErrInternal):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутрення ошибка сервера",
			})
			return
		case errors.Is(err, domain.ErrRefreshToken):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Не удалось авторизовать пользователя",
			})
		case errors.Is(err, domain.ErrAccessToken):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Не удалось авторизовать пользователя",
			})
		default:
			slog.Error("Unexpected error during signin",
				"email", input.Email,
				"error", err,
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}
	}

	maxAge := 30 * 24 * 60 * 60 // 30 дней
	c.SetCookie("refresh_token", tokenpair.RefreshToken, maxAge, "/api/v1/auth", "", true, true)

	slog.Info("successful SignIn attemp", "email", input.Email)
	c.JSON(http.StatusCreated, gin.H{
		"access_token": tokenpair.AccessToken,
	})

}

// @Summary Обновление токен пары
// @Description Обновляет токен пару клиента
// @Security BearerAuth
// @Tags Refresh
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Успешное обновление токена"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 401 {object} map[string]string "Токен обновления пуст"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("Missing refresh token cookie")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Токен обновления пуст",
		})
	}

	tokenpair, err := h.AuthService.Refresh(c.Request.Context(), refreshToken)
	if err != nil {

		switch {
		case errors.Is(err, domain.ErrRefreshToken):
			//slog.Error("Invalid refresh token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Ошибка обновления токена",
			})
			return
		case errors.Is(err, domain.ErrInternal):

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		case errors.Is(err, domain.ErrConditions):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Ошибка проверки данных",
			})
			return
		case errors.Is(err, domain.ErrRedis):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка сервера, попробуйте позже",
			})
			return
		case errors.Is(err, domain.ErrAccessToken):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Ошибка обновления токена",
			})
			return
		default:
			slog.Error("Unexpected error during refresh tokens")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}
	}

	maxAgeSeconds := 30 * 24 * 60 * 60
	c.SetCookie("refresh_token", tokenpair.RefreshToken, maxAgeSeconds,
		"/api/v1/auth", "", true, true)

	slog.Info("successful refresh attemp")
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenpair.AccessToken,
	})
}

// @Summary Удаление токена из Redis
// @Description Оставляет токен клиента, но удаляет полностью refresh token
// @Security BearerAuth
// @Tags Logout
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Успешная деавторизация"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Failure 401 {object} map[string]string "Токен обновления пуст"

// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("Missing refresh token cookie")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Токен обновления пуст",
		})
	}

	err = h.AuthService.Logout(c.Request.Context(), refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRefreshToken):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка обновления токена",
			})
			return
		case errors.Is(err, domain.ErrInternal):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		case errors.Is(err, domain.ErrRedis):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка сервера, попробуйте позже",
			})
			return
		default:
			slog.Error("Unexpected error during logout")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}

	}

	// Удаление куки
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", true, true)

	slog.Info("successful logout")
	c.JSON(http.StatusOK, gin.H{
		"message": "Успешная деавторизация",
	})
}
