package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwttoken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Не удалось найти token",
			})
			return
		}
		parts := strings.Split(tokenString, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Токен не валиден",
			})
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwttoken), nil

		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Ошибка валидации токена",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok1 := claims["sub"].(string)
			if ok1 {
				c.Set("userID", userID)
			} else {
				slog.Error("Отсутствует userID: %s", "error", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Отсутствует userID",
				})
				return
			}
			userRole, ok2 := claims["role"].(string)
			if ok2 {
				c.Set("userRole", userRole)
			} else {
				slog.Error("Отсутствует role: %s", "error", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Отсутствует role",
				})
				return
			}
		}

		c.Next()

	}

}

func OrganizerMiddleware(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Роль пользователя не найдена",
			})
			return
		}

		rolestr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Неверный формат роли",
			})
		}

		for _, allowedRole := range allowedRoles {
			if rolestr == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Недостаточно прав пользователя",
		})

	}
}
