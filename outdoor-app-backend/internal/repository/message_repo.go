package repository

import (
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"
)

// CreateMessage 创建一条新的系统消息 (供 Service 层调用)
func CreateMessage(msg *model.Message) error {
	return database.DB.Create(msg).Error
}

// GetUserMessages 分页获取用户的消息列表
// 增加了 page (当前页) 和 pageSize (每页条数) 参数
func GetUserMessages(userID uint, page, pageSize int) ([]model.Message, error) {
	var msgs []model.Message

	// 计算偏移量
	offset := (page - 1) * pageSize

	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&msgs).Error

	return msgs, err
}

// MarkRead 将消息标记为已读
func MarkRead(messageID uint) error {
	return database.DB.Model(&model.Message{}).
		Where("id = ?", messageID).
		Update("is_read", true).Error
}

// GetUnreadCount 获取未读消息数量 (前端右上角小红点常用)
func GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Message{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// GetUnreadMessagesByUserID 获取用户的未读消息列表
func GetUnreadMessagesByUserID(uid uint) ([]*model.Message, error) {
	var msgs []*model.Message
	err := database.DB.Where("user_id = ? AND is_read = ?", uid, false).
		Order("created_at DESC").
		Find(&msgs).Error
	return msgs, err
}
