package middleware

import (
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/jwt"
	"outdoor-app-backend/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization 字段
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Fail(c, e.Unauthorized, "请求头中缺少 Token")
			c.Abort() // 拦截请求
			return
		}

		// 2. 按空格分割 (标准的 Token 格式为 "Bearer eyJhbG...")
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Fail(c, e.Unauthorized, "Token 格式错误")
			c.Abort()
			return
		}

		// 3. 解析 Token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, e.Unauthorized, "Token 无效或已过期，请重新登录")
			c.Abort()
			return
		}

		// 4. 解析成功，将 userID 存入 Gin 的上下文 (Context) 中
		// 这样后续的 Handler 就可以直接通过 c.Get("userID") 知道是哪个用户了！
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		c.Next() // 放行，继续执行后续的 Handler
	}
}
