package model

import "time"

// Activity 活动约伴表
type Activity struct {
	BaseModel
	InitiatorID uint      `gorm:"not null;index" json:"initiator_id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	City        string    `gorm:"type:varchar(50);index;comment:'活动所属城市'" json:"city"` // 👈 新增：用于筛选
	Destination string    `gorm:"type:varchar(100);not null" json:"destination"`
	GatherTime  time.Time `gorm:"not null;index;comment:'集合时间'" json:"gather_time"`
	GatherPlace string    `gorm:"type:varchar(100);not null" json:"gather_place"`
	FeeType     string    `gorm:"type:varchar(20);comment:'AA/免费'" json:"fee_type"`
	GroupQrCode string    `gorm:"type:varchar(255);comment:'微信群二维码'" json:"group_qr_code"`

	// 👈 新增：图文描述与封面
	CoverImage  string `gorm:"type:varchar(255);comment:'活动封面'" json:"cover_image"`
	Description string `gorm:"type:text;comment:'活动具体描述'" json:"description"`
	Images      string `gorm:"type:json;comment:'活动补充图片数组'" json:"images"`

	// 门槛与状态
	MinFitness     int    `gorm:"type:tinyint;default:1" json:"min_fitness"`
	AgeLimit       string `gorm:"type:varchar(50)" json:"age_limit"`
	MaxMembers     int    `gorm:"type:int;not null" json:"max_members"`
	CurrentMembers int    `gorm:"type:int;default:1" json:"current_members"`

	// enrolling(报名中), full(满员), ongoing(进行中), ended(已结束)
	Status       string `gorm:"type:varchar(20);default:'enrolling';index" json:"status"`
	WeatherAlert bool   `gorm:"default:false;comment:'是否有恶劣天气'" json:"weather_alert"`

	Initiator User `gorm:"foreignKey:InitiatorID" json:"initiator"`
}

// ActivityMember 报名记录 (保持不变)
type ActivityMember struct {
	BaseModel
	ActivityID    uint   `gorm:"not null;uniqueIndex:idx_act_user" json:"activity_id"`
	UserID        uint   `gorm:"not null;uniqueIndex:idx_act_user" json:"user_id"`
	EquipmentNote string `gorm:"type:varchar(255);comment:'装备备注'" json:"equipment_note"`
	Status        string `gorm:"type:varchar(20);default:'pending'" json:"status"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Activity Activity `gorm:"foreignKey:ActivityID" json:"activity"`
}
