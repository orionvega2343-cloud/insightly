package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RateLimiter(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		id, ok := GetUserID(c)
		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "invalid user id"})
			c.Abort()
			return
		}
		key := fmt.Sprintf("rate-limit-%d", id)
		res, err := limiter.Allow(ctx, key, redis_rate.PerHour(10))
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		if res.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()

	}
}
