package handler

import (
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PublishArticle 专家发布百科
func PublishArticle(c *gin.Context) {
	var req dto.ArticleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误: "+err.Error())
		return
	}

	authorID := getUserID(c) // 借助之前写好的辅助函数
	if err := service.CreateArticle(authorID, &req); err != nil {
		response.Fail(c, e.Error, "发布失败")
		return
	}
	response.Success(c, "发布成功")
}

// GetArticleList 获取百科列表 (公共接口)
func GetArticleList(c *gin.Context) {
	page, pageSize := getPagination(c) // 调用之前写的辅助函数
	list, err := service.GetArticleList(page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取百科列表失败")
		return
	}
	response.Success(c, list)
}

// UpdateArticle 专家修改百科
func UpdateArticle(c *gin.Context) {
	// 1. 获取 URL 中的 /update/:id 参数
	articleIDStr := c.Param("id")
	articleID, _ := strconv.Atoi(articleIDStr)

	// 2. 解析请求 Body
	var req dto.ArticleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误")
		return
	}

	// 3. 执行修改逻辑
	authorID := getUserID(c)
	if err := service.UpdateArticle(authorID, uint(articleID), &req); err != nil {
		response.Fail(c, e.Forbidden, err.Error())
		return
	}
	response.Success(c, "文章更新成功")
}

// DeleteArticle 专家删除百科
func DeleteArticle(c *gin.Context) {
	// 1. 获取 URL 中的 /delete/:id 参数
	articleIDStr := c.Param("id")
	articleID, _ := strconv.Atoi(articleIDStr)

	// 2. 执行删除逻辑
	authorID := getUserID(c)
	if err := service.DeleteArticle(authorID, uint(articleID)); err != nil {
		response.Fail(c, e.Forbidden, err.Error())
		return
	}
	response.Success(c, "文章已删除")
}
