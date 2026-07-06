package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		cut := strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(cut, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("email", claims.Email)
		c.Next()
	}
}
