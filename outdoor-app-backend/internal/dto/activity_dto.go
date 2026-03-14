package dto

import "time"

// ActivityFilterReq 瀑布流筛选条件 (GET 参数)
type ActivityFilterReq struct {
	City       string `form:"city"`
	MinFitness int    `form:"min_fitness"` // 体能筛选
	Status     string `form:"status"`      // 状态筛选(如：只看 enrolling)
	AgeLimit   string `form:"age_limit"`   // 年龄限制筛选
}

// ActivityCreateReq 发布活动请求
type ActivityCreateReq struct {
	Title       string    `json:"title" binding:"required,max=100"`
	City        string    `json:"city" binding:"required"`
	Destination string    `json:"destination" binding:"required"`
	GatherTime  time.Time `json:"gather_time" binding:"required"`
	GatherPlace string    `json:"gather_place" binding:"required"`
	FeeType     string    `json:"fee_type" binding:"required"`
	GroupQrCode string    `json:"group_qr_code" binding:"required"` // 必填，但只给通过的人看
	CoverImage  string    `json:"cover_image" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Images      []string  `json:"images" binding:"omitempty"`
	MinFitness  int       `json:"min_fitness" binding:"required,min=1,max=5"`
	AgeLimit    string    `json:"age_limit" binding:"omitempty"`
	MaxMembers  int       `json:"max_members" binding:"required,min=2"`
}

// ActivityDetailRes 活动详情精简返回格式
type ActivityDetailRes struct {
	ID             uint          `json:"id"`
	Title          string        `json:"title"`
	City           string        `json:"city"`
	Destination    string        `json:"destination"`
	GatherTime     time.Time     `json:"gather_time"`
	GatherPlace    string        `json:"gather_place"`
	FeeType        string        `json:"fee_type"`
	GroupQrCode    string        `json:"group_qr_code"` // ⚠️ 后端在返回时动态控制是否有值
	CoverImage     string        `json:"cover_image"`
	Description    string        `json:"description"`
	Images         []string      `json:"images"`
	MinFitness     int           `json:"min_fitness"`
	AgeLimit       string        `json:"age_limit"`
	MaxMembers     int           `json:"max_members"`
	CurrentMembers int           `json:"current_members"`
	Status         string        `json:"status"`
	WeatherAlert   bool          `json:"weather_alert"`
	Initiator      UserBasicInfo `json:"initiator"` // 发起人信息

	// 当前用户的报名状态聚合
	HasApplied  bool   `json:"has_applied"`  // 是否已报名
	ApplyStatus string `json:"apply_status"` // pending/approved/rejected (若未报名为空)
}

// ActivityApplyReq 报名与审核请求 (共用简单结构)
type ActivityApplyReq struct {
	EquipmentNote string `json:"equipment_note" binding:"omitempty,max=255"`
}

type ActivityAuditReq struct {
	MemberID uint   `json:"member_id" binding:"required"`
	Status   string `json:"status" binding:"required,oneof=approved rejected"`
}

// ActivityMemberRes 报名用户精简列表
type ActivityMemberRes struct {
	MemberID      uint          `json:"member_id"` // 👈 用于调用审核接口的主键
	EquipmentNote string        `json:"equipment_note"`
	Status        string        `json:"status"` // pending/approved/rejected
	User          UserBasicInfo `json:"user"`   // 报名者的头像昵称
}
