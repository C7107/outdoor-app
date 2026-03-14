package repository

import (
	"errors"
	"outdoor-app-backend/internal/database" // ⚠️ 替换为你的项目模块名
	"outdoor-app-backend/internal/model"

	"gorm.io/gorm"
)

// CreateUser 在数据库中创建新用户 (注册用)
func CreateUser(user *model.User) error {
	// gorm 的 Create 方法会自动把生成的主键 ID 塞回 user 变量里
	return database.DB.Create(user).Error
}

// GetUserByEmail 根据邮箱查找用户 (登录、防重复注册用)
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	// 使用 Where 查询，First 表示只拿第一条记录
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没找到用户不算代码错误，返回 nil 让业务层判断
		}
		return nil, err // 真正的数据库查询错误
	}
	return &user, nil
}

// GetUserByID 根据 ID 查找用户 (获取个人信息用)
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息 (修改个人档案用)
func UpdateUser(userID uint, data map[string]interface{}) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(data).Error
}
