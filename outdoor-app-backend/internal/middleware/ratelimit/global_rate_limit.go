package ratelimit

import (
	"net/http"
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func GlobalRateLimit(maxRequests int64, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {

		key := "rate_limit:global"

		count, err := database.RedisClient.Incr(database.Ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		// 第一次请求设置过期时间
		if count == 1 {
			database.RedisClient.Expire(database.Ctx, key, time.Duration(windowSeconds)*time.Second)
		}

		if count > maxRequests {
			response.Fail(c, http.StatusTooManyRequests, "服务器繁忙，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
