package repository

import (
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"

	"gorm.io/gorm"
)

func CreatePost(post *model.Post) error {
	return database.DB.Create(post).Error
}

func GetPostList(page, pageSize int) ([]model.Post, error) {
	var posts []model.Post
	err := database.DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("User"). // 预加载发帖人信息
		Find(&posts).Error
	return posts, err
}

func CreateComment(comment *model.Comment) error {
	return database.DB.Create(comment).Error
}

// GetCommentsByPostID 获取某条动态的所有评论
func GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := database.DB.Where("post_id = ?", postID).Preload("User").Find(&comments).Error
	return comments, err
}

// 🌟 修复后的 ToggleLike：点赞操作（使用事务保证总数准确）
func ToggleLike(postID, userID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error { //database.DB.Transaction 是 GORM 提供的一个事务处理函数
		var like model.PostLike
		err := tx.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error

		if err == nil {
			// 1. 如果记录存在，说明已赞 -> 执行取消赞
			if err := tx.Delete(&like).Error; err != nil {
				return err
			}
			// 将总赞数 - 1
			return tx.Model(&model.Post{}).Where("id = ?", postID).Update("like_count", gorm.Expr("like_count - 1")).Error
		}

		// 2. 如果记录不存在 -> 执行点赞
		if err := tx.Create(&model.PostLike{PostID: postID, UserID: userID}).Error; err != nil {
			return err
		}
		// 将总赞数 + 1
		return tx.Model(&model.Post{}).Where("id = ?", postID).Update("like_count", gorm.Expr("like_count + 1")).Error
	})
}

// IsPostLiked 判断当前用户是否对该动态点过赞 (用于前端高亮心形图标)
func IsPostLiked(postID, userID uint) bool {
	var count int64
	database.DB.Model(&model.PostLike{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count)
	return count > 0
}

// GetPostListByUser 获取某个用户发布的全部动态（分页）
func GetPostListByUser(userID uint, page, pageSize int) ([]model.Post, error) {
	var posts []model.Post
	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("User"). // 同样预加载发布者信息
		Find(&posts).Error
	return posts, err
}

// GetPostByID 根据 ID 查找单条动态 (用于删除前的权限校验)
func GetPostByID(postID uint) (*model.Post, error) {
	var post model.Post
	err := database.DB.First(&post, postID).Error
	return &post, err
}

// DeletePost 删除单条动态 (由于使用了 gorm.DeletedAt，这里是软删除)
// 如果希望彻底删除，可以用 database.DB.Unscoped().Delete(&post)
func DeletePost(postID uint) error {
	return database.DB.Delete(&model.Post{}, postID).Error
}
