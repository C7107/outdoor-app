package model

// Route 户外路线
type Route struct {
	BaseModel
	Title        string  `gorm:"type:varchar(100);index" json:"title"`
	City         string  `gorm:"type:varchar(50);index" json:"city"`
	Difficulty   int     `gorm:"type:tinyint;index;comment:'难度1-5'" json:"difficulty"`
	DurationDays int     `gorm:"type:int;comment:'耗时(天)'" json:"duration_days"`
	SceneryScore float64 `gorm:"type:decimal(3,1);comment:'风景指数'" json:"scenery_score"`
	CoverImage   string  `gorm:"type:varchar(255)" json:"cover_image"`
	MapTrackUrl  string  `gorm:"type:varchar(255);comment:'地图轨迹'" json:"map_track_url"`
	ElevationUrl string  `gorm:"type:varchar(255);comment:'海拔图'" json:"elevation_url"`
	Description  string  `gorm:"type:text;comment:'路况描述'" json:"description"`
}

// RouteReview 路线评价 (一对多)
type RouteReview struct {
	BaseModel
	RouteID uint   `gorm:"not null;index" json:"route_id"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	Score   int    `gorm:"type:tinyint;not null;comment:'评分1-5'" json:"score"`
	Content string `gorm:"type:text" json:"content"`
	Images  string `gorm:"type:json;comment:'图文评价的图片数组'" json:"images"` // 存 JSON 数组

	User User `gorm:"foreignKey:UserID" json:"user"` // 关联用户信息
}

// FavoriteRoute 用户的“路线收藏”中间表 (多对多)
type FavoriteRoute struct {
	UserID  uint `gorm:"primaryKey" json:"user_id"`
	RouteID uint `gorm:"primaryKey" json:"route_id"`
}
