package service

import (
	"encoding/json"
	"errors"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
)

// PublishPost 发布朋友圈动态
func PublishPost(userID uint, req *dto.PostCreateReq) error {
	// 把前端传来的字符串数组，转换成 JSON 字符串存入数据库
	imagesJSON, _ := json.Marshal(req.Images)

	post := &model.Post{
		UserID:  userID,
		Content: req.Content,
		Images:  string(imagesJSON), // 存入 JSON 字符串
	}
	return repository.CreatePost(post)
}

// AddComment 发表评论
func AddComment(userID uint, req *dto.CommentCreateReq) error {
	comment := &model.Comment{
		PostID:  req.PostID,
		UserID:  userID,
		Content: req.Content,
	}
	return repository.CreateComment(comment)
}

// ToggleLike 点赞与取消点赞
func ToggleLike(userID, postID uint) error {
	return repository.ToggleLike(postID, userID)
}

// ==========================================
// 💡 这是一个高级组装逻辑：前端不仅需要动态，还需要动态的评论和是否已赞状态
type PostRes struct {
	ID      uint       `json:"id"`
	Content string     `json:"content"`
	Images  []string   `json:"images"`
	User    model.User `json:"user"`

	IsLiked  bool            `json:"is_liked"`
	Comments []model.Comment `json:"comments"`
}

// GetPostList 获取公共动态列表
func GetPostList(userID uint, page, pageSize int) ([]dto.PostRes, error) {
	posts, err := repository.GetPostList(page, pageSize)
	if err != nil {
		return nil, err
	}
	return buildPostResList(userID, posts) // 调用底部的公共转换函数
}

// GetMyPostList 获取我发布的动态列表
func GetMyPostList(userID uint, page, pageSize int) ([]dto.PostRes, error) {
	posts, err := repository.GetPostListByUser(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return buildPostResList(userID, posts) // 调用底部的公共转换函数
}

// ==========================================
// 辅助函数：将数据库的 []model.Post 完美转换为前端需要的 []dto.PostRes
// ==========================================
func buildPostResList(currentUserID uint, posts []model.Post) ([]dto.PostRes, error) {
	var res []dto.PostRes

	// 如果没有数据，返回一个空数组而不是 null，对前端更友好
	if len(posts) == 0 {
		return []dto.PostRes{}, nil
	}

	for _, p := range posts {
		// 1. 解析图片
		var images []string
		json.Unmarshal([]byte(p.Images), &images)

		// 2. 处理评论（把评论也转换为精简格式）
		commentsData, _ := repository.GetCommentsByPostID(p.ID)
		var commentsRes []dto.CommentRes
		for _, c := range commentsData {
			commentsRes = append(commentsRes, dto.CommentRes{
				ID:        c.ID,
				Content:   c.Content,
				CreatedAt: c.CreatedAt,
				User: dto.UserBasicInfo{
					ID:       c.User.ID,
					Nickname: c.User.Nickname,
					Avatar:   c.User.Avatar,
				},
			})
		}

		// 3. 处理点赞状态
		isLiked := repository.IsPostLiked(p.ID, currentUserID)

		// 4. 组装最终的最干净的 DTO
		res = append(res, dto.PostRes{
			ID:        p.ID,
			Content:   p.Content,
			Images:    images,
			CreatedAt: p.CreatedAt, // ✅ 补充了动态的创建时间
			User: dto.UserBasicInfo{ // ✅ 彻底告别臃肿的 user 对象
				ID:       p.User.ID,
				Nickname: p.User.Nickname,
				Avatar:   p.User.Avatar,
			},
			IsLiked:  isLiked,
			Comments: commentsRes,
		})
	}
	return res, nil
}

// DeletePost 删除自己的动态 (防越权)
func DeletePost(userID, postID uint) error {
	// 1. 先查出这篇动态
	post, err := repository.GetPostByID(postID)
	if err != nil {
		return errors.New("动态不存在或已被删除")
	}

	// 2. 🔴 极度关键：水平越权校验！防止别人通过改参数删掉你的动态
	if post.UserID != userID {
		return errors.New("无权删除他人的动态")
	}

	// 3. 执行删除
	return repository.DeletePost(postID)
}
