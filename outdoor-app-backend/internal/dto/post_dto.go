package dto

import (
	"time"
)

// PostCreateReq 发布动态请求
type PostCreateReq struct {
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images" binding:"omitempty"` // 前端传数组
}

// CommentCreateReq 发布评论请求
type CommentCreateReq struct {
	PostID  uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UserBasicInfo 专门用于列表展示的基础用户信息 (绝不包含敏感数据)
type UserBasicInfo struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CommentRes 评论的精简返回格式
type CommentRes struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"` // 评论的时间
	User      UserBasicInfo `json:"user"`       // 评论者的精简信息
}

// PostRes 动态列表的精简返回格式
type PostRes struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	Images    []string      `json:"images"`
	CreatedAt time.Time     `json:"created_at"` // 👈 新增：动态发布时间
	User      UserBasicInfo `json:"user"`       // 👈 优化：只返回发帖人的头像和昵称

	IsLiked  bool         `json:"is_liked"`
	Comments []CommentRes `json:"comments"` // 👈 优化：评论里的用户信息也精简了
}
