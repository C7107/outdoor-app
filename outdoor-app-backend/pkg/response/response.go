package response

import (
	"net/http"
	"outdoor-app-backend/pkg/e" // ⚠️ 记得换成你的项目名

	"github.com/gin-gonic/gin"
)

// Response 基础返回结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 请求成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: data,
	})
}

// Fail 请求失败返回
func Fail(c *gin.Context, code int, msg string) {
	// 如果 msg 为空，则根据 code 自动获取默认 msg
	if msg == "" {
		msg = e.GetMsg(code)
	}
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
