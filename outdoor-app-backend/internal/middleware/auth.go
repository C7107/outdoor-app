package middleware

import (
	"outdoor-app-backend/internal/repository"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExpertRequired 专家权限校验中间件
// 前置条件：必须先使用 JWTAuth 中间件
func ExpertRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Context 中获取 JWTAuth 存入的 userID
		userID, exists := c.Get("userID")
		if !exists {
			response.Fail(c, e.Unauthorized, "未登录")
			c.Abort()
			return
		}

		// 2. 查数据库获取用户角色
		user, err := repository.GetUserByID(userID.(uint))
		if err != nil || user == nil {
			response.Fail(c, e.ErrorUserNotFound, "用户不存在")
			c.Abort()
			return
		}

		// 3. 判断是否为专家或管理员
		if user.Role != "expert" && user.Role != "admin" {
			response.Fail(c, e.Forbidden, "无权限：仅专家可操作")
			c.Abort()
			return
		}

		// 4. 放行
		c.Next()
	}
}
