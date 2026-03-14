package service

import (
	"errors"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
	"outdoor-app-backend/pkg/jwt"
	"outdoor-app-backend/pkg/utils"
)

// Register 处理用户注册逻辑
func Register(req *dto.UserRegisterReq) error {
	// 1. (模拟) 校验验证码，前期为了方便测试，先固定写死 123456
	if req.Code != "123456" {
		return errors.New("验证码错误")
	}

	defaultAvatar := "http://127.0.0.1:8080/uploads/avatar/123456789.jpg"
	// 2. 检查邮箱是否已被注册
	existUser, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return err // 数据库查询出错
	}
	if existUser != nil {
		return errors.New("该邮箱已被注册")
	}

	// 3. 对密码进行 bcrypt 加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 4. 组装 Model 并存入数据库
	newUser := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Nickname:     "户外新驴_" + req.Email[:4], // 默认给个昵称
		Role:         "user",                  // 默认普通用户
		Avatar:       defaultAvatar,           // 👈 在这里赋值默认头像 URL
	}

	return repository.CreateUser(newUser)
}

// Login 处理用户登录逻辑
func Login(req *dto.UserLoginReq) (*dto.UserLoginRes, error) {
	// 1. 根据邮箱去数据库找人
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("服务器错误")
	}
	if user == nil {
		return nil, errors.New("该邮箱尚未注册")
	}

	// 2. 校验密码 (明文 vs 数据库里的哈希值)
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 3. 密码正确，生成 JWT Token 派发给前端
	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("生成身份令牌失败")
	}

	// 4. 组装 Response DTO 返回
	res := &dto.UserLoginRes{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     user.Role,
	}

	return res, nil
}
