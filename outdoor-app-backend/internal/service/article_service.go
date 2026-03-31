package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
	"time"
)

func CreateArticle(authorID uint, req *dto.ArticleCreateReq) error {
	article := &model.Article{
		AuthorID: authorID,
		Title:    req.Title,
		Category: req.Category,
		Cover:    req.Cover,
		Content:  req.Content,
	}

	// 1. 入库
	if err := repository.CreateArticle(article); err != nil {
		return err
	}

	// 2. 🌟 删除所有文章分页缓存（关键！！）
	// 模糊删除 articles:page:*
	pattern := "articles:page:*"
	keys, err := database.RedisClient.Keys(database.Ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		database.RedisClient.Del(database.Ctx, keys...)
		fmt.Println("🧹 已清除文章列表缓存")
	}

	return nil
}

// 获取某人的文章列表
func GetMyArticles(authorID uint, page, pageSize int) ([]model.Article, error) {
	return repository.GetArticlesByAuthor(authorID, page, pageSize)
}

// GetArticleList 分页获取公共百科列表 (带 Redis 缓存)
func GetArticleList(page, pageSize int) ([]model.Article, error) {
	// 1. 定义 Redis 中的 Key (根据页码不同，Key也不同，如: articles:page:1)
	cacheKey := fmt.Sprintf("articles:page:%d:size:%d", page, pageSize)

	// 2. 尝试从 Redis 中获取数据
	cachedData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		// 🌟 缓存命中！直接将 JSON 字符串反序列化并返回
		var articles []model.Article
		json.Unmarshal([]byte(cachedData), &articles)
		fmt.Println("🚀 从 Redis 缓存中极速读取了数据！")
		return articles, nil
	}

	// 3. 缓存未命中：老老实实去 MySQL 查数据库
	fmt.Println("🐌 缓存未命中，从 MySQL 数据库读取...")
	articles, err := repository.GetArticleList(page, pageSize)
	if err != nil {
		return nil, err
	}

	// 4. 查到数据后，序列化成 JSON 字符串，存入 Redis
	// 🌟 设置过期时间为 5 分钟，5分钟后自动失效，下次请求又会去查 MySQL 保证数据更新
	if jsonData, err := json.Marshal(articles); err == nil {
		database.RedisClient.Set(database.Ctx, cacheKey, jsonData, 5*time.Minute)
	}

	return articles, nil
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
	if err := repository.UpdateArticlePartial(articleID, updateData); err != nil {
		return err
	}
	// 5. 🌟 删除所有文章分页缓存（关键！！）

	pattern := "articles:page:*"
	keys, err := database.RedisClient.Keys(database.Ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		database.RedisClient.Del(database.Ctx, keys...)
		fmt.Println("🧹 已清除文章列表缓存")
	}

	return nil
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

	if err := repository.DeleteArticle(articleID); err != nil {
		return err
	}

	// 5. 🌟 删除所有文章分页缓存（关键！！）
	pattern := "articles:page:*"
	keys, err := database.RedisClient.Keys(database.Ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		database.RedisClient.Del(database.Ctx, keys...)
		fmt.Println("🧹 已清除文章列表缓存")
	}

	return nil
}
