package ratelimit

import (
	"fmt"
	"net/http"
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// 脚本，用Lua将INCR和EXPIRE打包成原子操作，解决了代码中可能出现的key永不过期的问题。
var luaScript = `
local current = redis.call("INCR", KEYS[1])
if tonumber(current) == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`

// RateLimit 限流中间件
// maxRequests: 时间窗口内允许的最大请求数
// windowSeconds: 时间窗口大小(秒)
func RateLimit(maxRequests int64, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {

		var key string

		// 优先使用 userID，否则 IP + API
		userID, exists := c.Get("userID")
		if exists {
			key = fmt.Sprintf("rate_limit:user:%v:%s", userID, c.FullPath())
		} else {
			ip := c.ClientIP()
			key = fmt.Sprintf("rate_limit:ip:%s:%s", ip, c.FullPath())
		}

		res, err := database.RedisClient.Eval(
			database.Ctx,
			luaScript,
			[]string{key},
			windowSeconds,
		).Result()

		if err != nil {
			// Redis挂了直接放行
			c.Next()
			return
		}

		count := res.(int64)

		if count > maxRequests {
			response.Fail(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
