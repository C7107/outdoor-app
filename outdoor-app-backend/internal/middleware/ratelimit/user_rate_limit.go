package ratelimit

import (
	"fmt"
	"net/http"
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func UserGlobalRateLimit(maxRequests int64, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {

		var key string

		// 优先使用 userID
		userID, exists := c.Get("userID")
		if exists {
			key = fmt.Sprintf("rate_limit:user_global:%v", userID)
		} else {
			ip := c.ClientIP()
			key = fmt.Sprintf("rate_limit:ip_global:%s", ip)
		}

		count, err := database.RedisClient.Incr(database.Ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		// 第一次访问设置过期
		if count == 1 {
			database.RedisClient.Expire(
				database.Ctx,
				key,
				time.Duration(windowSeconds)*time.Second,
			)
		}

		if count > maxRequests {
			response.Fail(c, http.StatusTooManyRequests, "请求次数过多，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
