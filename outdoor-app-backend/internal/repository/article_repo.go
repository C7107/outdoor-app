package repository

import (
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"
)

// GetArticleList 获取所有百科列表（公共接口，所有人可看，需分页）
func GetArticleList(page, pageSize int) ([]model.Article, error) {
	var articles []model.Article
	err := database.DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error
	return articles, err
}

// GetArticlesByAuthor 获取某个专家发布的全部文章（需分页）
func GetArticlesByAuthor(authorID uint, page, pageSize int) ([]model.Article, error) {
	var articles []model.Article
	err := database.DB.Where("author_id = ?", authorID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error
	return articles, err
}

// CreateArticle 创建一篇新的百科文章 (供 Service 层调用)
func CreateArticle(article *model.Article) error {
	return database.DB.Create(article).Error
}

// UpdateArticle 更新整篇文章 (仅作者本人可操作，需在 Service 层校验权限)
func UpdateArticle(article *model.Article) error {
	return database.DB.Save(article).Error
}

// DeleteArticle 删除文章 (仅作者本人可操作，需在 Service 层校验权限)
func DeleteArticle(id uint) error {
	return database.DB.Delete(&model.Article{}, id).Error
}

// GetArticleByID 根据 ID 获取单篇文章 (用于校验是不是作者本人)
func GetArticleByID(id uint) (*model.Article, error) {
	var article model.Article
	err := database.DB.First(&article, id).Error
	return &article, err
}

// UpdateArticlePartial 局部更新文章内容 (借鉴你之前改用户资料的优秀写法)
func UpdateArticlePartial(id uint, data map[string]interface{}) error {
	return database.DB.Model(&model.Article{}).Where("id = ?", id).Updates(data).Error
}
