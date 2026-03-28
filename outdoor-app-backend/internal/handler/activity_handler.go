package handler

import (
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetActivityList 获取活动瀑布流
func GetActivityList(c *gin.Context) {
	var req dto.ActivityFilterReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, e.InvalidParams, "筛选参数错误")
		return
	}

	page, pageSize := getPagination(c)
	list, err := service.GetActivityList(&req, page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取活动列表失败")
		return
	}
	response.Success(c, list)
}

// GetActivityDetail 获取活动详情
func GetActivityDetail(c *gin.Context) {
	activityIDStr := c.Param("id")
	activityID, _ := strconv.Atoi(activityIDStr)

	data, err := service.GetActivityDetail(getUserID(c), uint(activityID))
	if err != nil {
		response.Fail(c, e.Error, err.Error())
		return
	}
	response.Success(c, data)
}

// CreateActivity 发起活动
func CreateActivity(c *gin.Context) {
	var req dto.ActivityCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误: "+err.Error())
		return
	}

	if err := service.CreateActivity(getUserID(c), &req); err != nil {
		response.Fail(c, e.Error, err.Error())
		return
	}
	response.Success(c, "发起活动成功")
}

// DeleteActivity 删除活动
func DeleteActivity(c *gin.Context) {
	activityIDStr := c.Param("id")
	activityID, _ := strconv.Atoi(activityIDStr)

	if err := service.DeleteActivity(getUserID(c), uint(activityID)); err != nil {
		response.Fail(c, e.Forbidden, err.Error())
		return
	}
	response.Success(c, "活动已删除")
}

// ApplyActivity 用户报名
func ApplyActivity(c *gin.Context) {
	activityIDStr := c.Param("id")
	activityID, _ := strconv.Atoi(activityIDStr)

	var req dto.ActivityApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误")
		return
	}

	if err := service.ApplyActivity(getUserID(c), uint(activityID), &req); err != nil {
		response.Fail(c, e.Forbidden, err.Error()) // 抛出前置拦截的错误
		return
	}
	response.Success(c, "报名申请已提交，等待审核")
}

// AuditMember 发起人审核报名
func AuditMember(c *gin.Context) {
	var req dto.ActivityAuditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误: 只能填 approved 或 rejected")
		return
	}

	if err := service.AuditActivityMember(getUserID(c), &req); err != nil {
		response.Fail(c, e.Forbidden, err.Error())
		return
	}
	response.Success(c, "审核操作成功")
}

// GetActivityMembers 获取报名人员列表
func GetActivityMembers(c *gin.Context) {
	activityIDStr := c.Param("id")
	activityID, _ := strconv.Atoi(activityIDStr)

	// 前端可以传 ?status=pending 来只看不看待审核的人
	status := c.Query("status")
	page, pageSize := getPagination(c)

	list, err := service.GetActivityMembers(getUserID(c), uint(activityID), status, page, pageSize)
	if err != nil {
		response.Fail(c, e.Forbidden, err.Error())
		return
	}
	response.Success(c, list)
}

// SearchActivityByLocation 模糊搜索目的地
func SearchActivityByLocation(c *gin.Context) {
	// 1. 获取搜索关键词 (query 参数)
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Fail(c, e.InvalidParams, "搜索关键词不能为空")
		return
	}

	// 2. 获取分页参数 (复用之前写好的 getPagination 辅助函数)
	page, pageSize := getPagination(c)

	// 3. 调用 Service
	list, err := service.SearchActivityByLocation(keyword, page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "搜索失败")
		return
	}

	// 4. 返回成功响应
	response.Success(c, list)
}

// GetCityWeather 获取指定城市的未来 5 天天气 (供前端详情页展示用)
func GetCityWeather(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		response.Fail(c, e.InvalidParams, "城市名不能为空")
		return
	}

	// 🌟 调用我们写好的智能缓存服务，不耗费额度极速返回
	weathers, err := service.GetAndSyncWeather(city)
	if err != nil {
		response.Fail(c, e.Error, "获取天气失败")
		return
	}

	response.Success(c, weathers)
}
