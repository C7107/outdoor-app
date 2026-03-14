package model

type User struct {
	BaseModel
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"` // 密码绝对不能返回给前端
	Nickname     string `gorm:"type:varchar(50)" json:"nickname"`
	Avatar       string `gorm:"type:varchar(255)" json:"avatar"`
	Signature    string `gorm:"type:varchar(255)" json:"signature"`

	// 毕设需求：体能与安全
	Role             string `gorm:"type:varchar(20);default:'user';comment:'user/expert'" json:"role"`
	FitnessLevel     int    `gorm:"type:tinyint;default:1;comment:'体能1-5'" json:"fitness_level"`
	EmergencyContact string `gorm:"type:varchar(20);comment:'紧急联系人'" json:"emergency_contact"`
}
