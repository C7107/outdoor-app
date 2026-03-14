package dto

// ArticleCreateReq 发布百科请求参数
type ArticleCreateReq struct {
	Title    string `json:"title" binding:"required,max=100"`
	Category string `json:"category" binding:"required"`
	Cover    string `json:"cover" binding:"omitempty"`
	Content  string `json:"content" binding:"required"`
}

// ArticleUpdateReq 修改百科请求参数
type ArticleUpdateReq struct {
	Title    string `json:"title" binding:"omitempty,max=100"`
	Category string `json:"category" binding:"omitempty"`
	Cover    string `json:"cover" binding:"omitempty"`
	Content  string `json:"content" binding:"omitempty"`
}
