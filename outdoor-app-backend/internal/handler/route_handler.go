package handler

import (
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRouteList 获取线路瀑布流
func GetRouteList(c *gin.Context) {
	var req dto.RouteFilterReq
	// ⚠️ 注意：GET 请求参数要用 ShouldBindQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, e.InvalidParams, "筛选参数错误")
		return
	}

	page, pageSize := getPagination(c)
	list, err := service.GetRouteList(&req, page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取线路失败")
		return
	}
	response.Success(c, list)
}

// GetRouteDetail 获取线路详情
func GetRouteDetail(c *gin.Context) {
	routeID, _ := strconv.Atoi(c.Param("id"))

	data, err := service.GetRouteDetail(getUserID(c), uint(routeID))
	if err != nil {
		response.Fail(c, e.Error, err.Error())
		return
	}
	response.Success(c, data)
}

// ToggleFavorite 切换收藏状态
func ToggleFavorite(c *gin.Context) {
	routeID, _ := strconv.Atoi(c.Param("id"))

	if err := service.ToggleFavorite(getUserID(c), uint(routeID)); err != nil {
		response.Fail(c, e.Error, "操作失败")
		return
	}
	response.Success(c, "操作成功")
}

// CreateReview 发布评价
func CreateReview(c *gin.Context) {
	routeID, _ := strconv.Atoi(c.Param("id"))

	var req dto.ReviewCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误: 评分1-5且内容必填")
		return
	}

	if err := service.CreateReview(getUserID(c), uint(routeID), &req); err != nil {
		response.Fail(c, e.Error, "发布评价失败")
		return
	}
	response.Success(c, "评价成功")
}

// PublishRoute 官方/专家发布新路线
func PublishRoute(c *gin.Context) {
	var req []dto.RouteCreateReq // 直接接收数组

	// 绑定 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数校验失败: "+err.Error())
		return
	}

	if err := service.CreateRoutes(req); err != nil {
		response.Fail(c, e.Error, "发布线路失败")
		return
	}

	response.Success(c, "发布线路成功")
}
