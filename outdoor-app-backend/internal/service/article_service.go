package service

import (
	"errors"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
)

func CreateArticle(authorID uint, req *dto.ArticleCreateReq) error {
	article := &model.Article{
		AuthorID: authorID,
		Title:    req.Title,
		Category: req.Category,
		Cover:    req.Cover,
		Content:  req.Content,
	}
	return repository.CreateArticle(article)
}

// 获取某人的文章列表
func GetMyArticles(authorID uint, page, pageSize int) ([]model.Article, error) {
	return repository.GetArticlesByAuthor(authorID, page, pageSize)
}

// GetArticleList 分页获取公共百科列表
func GetArticleList(page, pageSize int) ([]model.Article, error) {
	return repository.GetArticleList(page, pageSize)
}

// UpdateArticle 修改百科 (核心校验：只能改自己的)
func UpdateArticle(authorID, articleID uint, req *dto.ArticleUpdateReq) error {
	// 1. 先把文章查出来
	article, err := repository.GetArticleByID(articleID)
	if err != nil {
		return errors.New("文章不存在")
	}

	// 2. 🔴 安全校验：判断这篇文章是不是当前登录专家发的
	if article.AuthorID != authorID {
		return errors.New("无权修改他人的文章")
	}

	// 3. 组装要更新的字段 (Map 过滤零值)
	updateData := map[string]interface{}{}
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Category != "" {
		updateData["category"] = req.Category
	}
	if req.Cover != "" {
		updateData["cover"] = req.Cover
	}
	if req.Content != "" {
		updateData["content"] = req.Content
	}

	// 4. 执行更新
	return repository.UpdateArticlePartial(articleID, updateData)
}

// DeleteArticle 删除百科 (核心校验：只能删自己的)
func DeleteArticle(authorID, articleID uint) error {
	article, err := repository.GetArticleByID(articleID)
	if err != nil {
		return errors.New("文章不存在")
	}

	if article.AuthorID != authorID {
		return errors.New("无权删除他人的文章")
	}

	return repository.DeleteArticle(articleID)
}
