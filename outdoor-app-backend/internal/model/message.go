package model

// Message 个人中心：消息通知
type Message struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	Type      string `gorm:"type:varchar(50);comment:'类型:system/activity'" json:"type"`
	Title     string `gorm:"type:varchar(100);not null" json:"title"`
	Content   string `gorm:"type:varchar(500);not null" json:"content"`
	RelatedID uint   `gorm:"default:0;comment:'关联的业务ID(如活动ID)'" json:"related_id"` // 👈 新增
	IsRead    bool   `gorm:"default:false" json:"is_read"`
}
