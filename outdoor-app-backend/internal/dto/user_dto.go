package dto

// ==========================================
// 接收请求 (Request DTO)
// ==========================================

// UserRegisterReq 用户注册请求参数
type UserRegisterReq struct {
	Email    string `json:"email" binding:"required,email"`           // 必填，邮箱格式
	Password string `json:"password" binding:"required,min=6,max=20"` // 必填，6-20位
	Code     string `json:"code" binding:"required,len=6"`            // 必填，6位验证码 (目前阶段我们可以在逻辑里先写死或跳过真实校验)
}

// UserLoginReq 用户登录请求参数
type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ==========================================
// 返回响应 (Response DTO)
// ==========================================

// UserLoginRes 登录成功后返回给前端的数据
type UserLoginRes struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// UserProfileRes 个人主页返回数据
type UserProfileRes struct {
	ID               uint   `json:"id"`
	Email            string `json:"email"`
	Nickname         string `json:"nickname"`
	Avatar           string `json:"avatar"`
	Signature        string `json:"signature"`
	FitnessLevel     int    `json:"fitness_level"`
	EmergencyContact string `json:"emergency_contact"`
	Role             string `json:"role"`
}

// UserUpdateReq 修改个人资料请求
type UserUpdateReq struct {
	Nickname         string `json:"nickname" binding:"omitempty,max=20"`
	Avatar           string `json:"avatar" binding:"omitempty"` // 前端先调上传接口拿到URL，再传给这个接口
	Signature        string `json:"signature" binding:"omitempty,max=100"`
	FitnessLevel     int    `json:"fitness_level" binding:"omitempty,oneof=1 2 3 4 5"`
	EmergencyContact string `json:"emergency_contact" binding:"omitempty,len=11"`
}

type ChangePwdReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
