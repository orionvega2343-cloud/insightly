package middlewares

import (
	"insightly/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		cut := strings.TrimPrefix(tokenString, "Bearer ")

		//Валидация JWT токена
		token, err := jwt.ParseWithClaims(cut, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		//Приведение к type assertion
		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		//Установка значений и запуск мидлвара
		c.Set("userId", claims.UserId)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// GetUserID достаёт userId, установленный AuthMiddleware, из контекста запроса.
func GetUserID(c *gin.Context) (int, bool) {
	value, ok := c.Get("userId")
	if !ok {
		return 0, false
	}
	id, ok := value.(int)
	return id, ok
}
