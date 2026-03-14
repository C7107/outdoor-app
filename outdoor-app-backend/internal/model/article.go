package model

// Article 户外急救知识与百科
type Article struct {
	BaseModel
	AuthorID uint   `gorm:"not null;index;comment:'发布者ID(专家)'" json:"author_id"` // 👈 新增这一行
	Title    string `gorm:"type:varchar(100);not null" json:"title"`
	Category string `gorm:"type:varchar(50);index;comment:'分类:急救/常识'" json:"category"`
	Cover    string `gorm:"type:varchar(255)" json:"cover"`
	Content  string `gorm:"type:longtext;not null" json:"content"` // 存富文本或 Markdown
}
