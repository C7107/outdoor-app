package handler

import (
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/service"
	"outdoor-app-backend/pkg/e"
	"outdoor-app-backend/pkg/response"
	"outdoor-app-backend/pkg/upload"
	"outdoor-app-backend/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 辅助函数：获取上下文中存储的 userID
func getUserID(c *gin.Context) uint {
	id, _ := c.Get("userID")
	return id.(uint)
}

// 辅助函数：获取分页参数
func getPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	return page, pageSize
}

// GetProfile 获取个人资料
func GetProfile(c *gin.Context) {
	res, err := service.GetUserProfile(getUserID(c))
	if err != nil {
		response.Fail(c, e.Error, "获取用户信息失败")
		return
	}
	response.Success(c, res)
}

// UpdateProfile 修改个人资料
func UpdateProfile(c *gin.Context) {
	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误")
		return
	}
	if err := service.UpdateUserProfile(getUserID(c), &req); err != nil {
		response.Fail(c, e.Error, "更新失败")
		return
	}
	response.Success(c, "更新资料成功")
}

// 以下是各种列表获取接口...
func GetMyPublishedActivities(c *gin.Context) {
	page, pageSize := getPagination(c)
	data, err := service.GetMyPublished(getUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取失败")
		return
	}
	response.Success(c, data)
}

func GetMyJoinedActivities(c *gin.Context) {
	page, pageSize := getPagination(c)
	data, err := service.GetMyJoined(getUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取失败")
		return
	}
	response.Success(c, data)
}

func GetMyFavoriteRoutes(c *gin.Context) {
	data, err := service.GetMyFavorites(getUserID(c))
	if err != nil {
		response.Fail(c, e.Error, "获取失败")
		return
	}
	response.Success(c, data)
}

func GetMyMessages(c *gin.Context) {
	page, pageSize := getPagination(c)
	data, err := service.GetMyMessages(getUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取失败")
		return
	}
	response.Success(c, data)
}

// UploadImage 通用图片上传接口
func UploadImage(c *gin.Context) {
	// 1. 从前端请求的 form-data 中获取名为 "file" 的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, e.InvalidParams, "读取图片文件失败")
		return
	}

	// 2. 获取前端传来的业务类型 (比如 "avatar", "post", "route")
	// 如果前端没传type，默认放到 "temp" 文件夹里
	folderName := c.DefaultPostForm("type", "temp")

	// 3. 调用工具包把文件存到本地并生成 URL
	// 参数: gin上下文, 文件句柄, 要保存的文件夹名
	imageUrl, err := upload.SaveImageToLocal(c, file, folderName)
	if err != nil {
		response.Fail(c, e.Error, "保存图片失败: "+err.Error())
		return
	}

	// 4. 返回拼接好的网络 URL 给前端
	// 前端拿到这个 URL 后，就可以拿着它去请求你的 UpdateProfile 接口了
	response.Success(c, imageUrl)
}

// 在 profile_handler.go 追加以下代码：

// GetMyArticles 获取我发布的百科 (专家专属)
func GetMyArticles(c *gin.Context) {
	page, pageSize := getPagination(c)
	data, err := service.GetMyArticles(getUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取失败")
		return
	}
	response.Success(c, data)
}

// GetMyPosts 获取我发布的动态列表 (个人中心)
func GetMyPosts(c *gin.Context) {
	page, pageSize := getPagination(c)
	data, err := service.GetMyPostList(getUserID(c), page, pageSize)
	if err != nil {
		response.Fail(c, e.Error, "获取动态列表失败")
		return
	}
	response.Success(c, data)
}

// Logout 退出登录
func Logout(c *gin.Context) {
	// 在基于 JWT 的无状态鉴权中，后端其实不需要做任何数据库操作。
	// 真正的“退出登录”是前端主动把本地存储 (AsyncStorage) 里的 token 删掉。
	// 但为了语义完整，后端返回成功即可。
	// 进阶玩法：后端可以把这个 token 加入 Redis 的黑名单，防止其在过期前继续被使用。
	response.Success(c, "已成功退出登录")
}

func ChangePassword(c *gin.Context) {
	var req dto.ChangePwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, e.InvalidParams, "参数错误")
		return
	}

	userID := getUserID(c)
	// 1. 查出原密码
	user, _ := service.GetUserPasswordbyID(userID)

	// 2. 校验旧密码
	if !utils.CheckPassword(req.OldPassword, user) {
		response.Fail(c, e.ErrorPassword, "原密码错误")
		return
	}

	// 3. 加密新密码并保存
	newHash, _ := utils.HashPassword(req.NewPassword)
	service.UpdateUserPasswordbyID(userID, newHash)

	response.Success(c, "密码修改成功")
}
