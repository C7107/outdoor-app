package handler

import (
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PublishPost 发布动态
func PublishPost(c *gin.Context) {
	var req dto.PostCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误")
		return
	}

	if err := service.PublishPost(getUserID(c), &req); err != nil {
		response.Fail(c, e.Error, "发布失败")
		return
	}
	response.Success(c, "发布成功")
}

// AddComment 发表评论
func AddComment(c *gin.Context) {
	var req dto.CommentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误")
		return
	}

	if err := service.AddComment(getUserID(c), &req); err != nil {
		response.Fail(c, e.Error, "评论失败")
		return
	}
	response.Success(c, "评论成功")
}

// ToggleLike 点赞/取消点赞
func ToggleLike(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.Atoi(postIDStr)

	if err := service.ToggleLike(getUserID(c), uint(postID)); err != nil {
		response.Fail(c, e.Error, "操作失败")
		return
	}
	response.Success(c, "操作成功")
}

// GetPostList 获取动态列表
func GetPostList(c *gin.Context) {
	page, pageSize := getPagination(c)

	list, err := service.GetPostList(getUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取动态失败")
		return
	}
	response.Success(c, list)
}

// DeletePost 删除自己的动态
func DeletePost(c *gin.Context) {
	// 获取 URL 参数 /delete/:id
	postIDStr := c.Param("id")
	postID, _ := strconv.Atoi(postIDStr)

	// 调用 Service 执行删除
	if err := service.DeletePost(getUserID(c), uint(postID)); err != nil {
		response.Fail(c, e.Forbidden, err.Error())
		return
	}
	response.Success(c, "动态删除成功")
}
