package service

import (
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
)

// GetUnreadMessages 获取未读消息列表
func GetUnreadMessages(uid uint) ([]*model.Message, error) {
	return repository.GetUnreadMessagesByUserID(uid)
}

// MarkMessageRead 标记消息已读
func MarkMessageRead(messageID uint) error {
	return repository.MarkRead(messageID)
}

// GetUnreadMessageCount 获取未读消息数量
func GetUnreadMessageCount(userID uint) (int64, error) {
	return repository.GetUnreadCount(userID)
}
