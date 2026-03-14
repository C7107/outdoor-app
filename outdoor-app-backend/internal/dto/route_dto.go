package dto

import "time"

// RouteFilterReq 线路列表筛选参数 (用于 GET 请求)
type RouteFilterReq struct {
	City        string  `form:"city"`
	Difficulty  int     `form:"difficulty"`
	MaxDuration int     `form:"max_duration"` // 最大耗时
	MinScore    float64 `form:"min_score"`    // 最低风景指数
}

// ReviewCreateReq 发布路线评价请求
type ReviewCreateReq struct {
	Score   int      `json:"score" binding:"required,min=1,max=5"`
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images" binding:"omitempty"` // 传图片数组
}

// RouteReviewRes 评价的精简返回格式
type RouteReviewRes struct {
	ID        uint          `json:"id"`
	Score     int           `json:"score"`
	Content   string        `json:"content"`
	Images    []string      `json:"images"`
	CreatedAt time.Time     `json:"created_at"`
	User      UserBasicInfo `json:"user"` // 👈 完美复用之前写的精简版 User
}

// RouteDetailRes 线路详情组合格式
type RouteDetailRes struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	City         string  `json:"city"`
	Difficulty   int     `json:"difficulty"`
	DurationDays int     `json:"duration_days"`
	SceneryScore float64 `json:"scenery_score"`
	CoverImage   string  `json:"cover_image"`
	MapTrackUrl  string  `json:"map_track_url"`
	ElevationUrl string  `json:"elevation_url"`
	Description  string  `json:"description"`

	// 聚合字段
	IsFavorited bool             `json:"is_favorited"` // 当前用户是否已收藏
	Reviews     []RouteReviewRes `json:"reviews"`      // 评价列表

	// WeatherInfo interface{} `json:"weather_info"` // TODO: 预留给天气API
}

// RouteCreateReq 发布线路请求参数
type RouteCreateReq struct {
	Title        string  `json:"title" binding:"required,max=100"`
	City         string  `json:"city" binding:"required"`
	Difficulty   int     `json:"difficulty" binding:"required,min=1,max=5"`
	DurationDays int     `json:"duration_days" binding:"required,min=1"`
	SceneryScore float64 `json:"scenery_score" binding:"required,min=1,max=5"`
	CoverImage   string  `json:"cover_image" binding:"omitempty"`   // 封面图URL
	MapTrackUrl  string  `json:"map_track_url" binding:"omitempty"` // 轨迹文件URL
	ElevationUrl string  `json:"elevation_url" binding:"omitempty"` // 海拔图URL
	Description  string  `json:"description" binding:"required"`    // 详情描述
}
