package model

// Post 户外圈子动态 (晒图/轨迹)
type Post struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Images    string `gorm:"type:json;comment:'动态图片数组'" json:"images"`
	LikeCount int    `gorm:"type:int;default:0" json:"like_count"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

// Comment 动态评论
type Comment struct {
	BaseModel
	PostID  uint   `gorm:"not null;index" json:"post_id"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	Content string `gorm:"type:varchar(500);not null" json:"content"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

// PostLike 点赞记录 (多对多)
type PostLike struct {
	PostID uint `gorm:"primaryKey" json:"post_id"`
	UserID uint `gorm:"primaryKey" json:"user_id"`
}
