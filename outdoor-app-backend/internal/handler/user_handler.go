package handler

import (
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// Register 用户注册接口
func Register(c *gin.Context) {
	var req dto.UserRegisterReq

	// 1. 绑定并校验前端传来的 JSON 参数 (根据 DTO 里的 binding 标签拦截错误)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数校验失败: 请检查邮箱格式或密码长度")
		return
	}

	// 2. 调用 Service 层处理核心逻辑
	if err := service.Register(&req); err != nil {
		// 如果业务逻辑报错（比如邮箱已存在），把错误信息返回给前端
		response.Fail(c, e.Error, err.Error())
		return
	}

	// 3. 成功返回
	response.Success(c, "注册成功，请登录")
}

// Login 用户登录接口
func Login(c *gin.Context) {
	var req dto.UserLoginReq

	// 1. 绑定并校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数校验失败: 邮箱或密码不能为空")
		return
	}

	// 2. 调用 Service 层执行登录
	res, err := service.Login(&req)
	if err != nil {
		// 登录失败（密码错误或找不到人）
		response.Fail(c, e.ErrorPassword, "用户不存在或密码错误")
		return
	}

	// 3. 登录成功，把 Token 和用户信息返回给前端
	response.Success(c, res)
}
