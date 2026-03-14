package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400
	Unauthorized  = 401
	Forbidden     = 403

	// 可以根据业务继续追加，例如：
	ErrorUserExist    = 10001
	ErrorUserNotFound = 10002
	ErrorPassword     = 10003
)

// GetMsg 根据错误码获取错误信息
func GetMsg(code int) string {
	msgFlags := map[int]string{
		Success:           "操作成功",
		Error:             "服务器内部错误",
		InvalidParams:     "请求参数错误",
		Unauthorized:      "Token验证失败或已失效",
		Forbidden:         "无权限访问",
		ErrorUserExist:    "邮箱已被注册",
		ErrorUserNotFound: "用户不存在",
		ErrorPassword:     "密码错误",
	}
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}
	return msgFlags[Error]
}
