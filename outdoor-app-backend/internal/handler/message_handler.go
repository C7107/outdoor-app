package handler

import (
	"strconv"

	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetUnreadCount 获取未读消息数量
func GetUnreadCount(c *gin.Context) {

	userID := getUserID(c)

	count, err := service.GetUnreadMessageCount(userID)
	if err != nil {
		response.Fail(c, e.Error, "获取未读消息数量失败")
		return
	}

	response.Success(c, gin.H{
		"count": count,
	})
}

// MarkMessageRead 标记消息已读
func MarkMessageRead(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, e.InvalidParams, "消息ID无效")
		return
	}

	err = service.MarkMessageRead(uint(id))
	if err != nil {
		response.Fail(c, e.Error, "标记已读失败")
		return
	}

	response.Success(c, nil)
}
