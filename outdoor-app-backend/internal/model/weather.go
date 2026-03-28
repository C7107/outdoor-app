package model

import "time"

// CityWeather 城市天气缓存表
type CityWeather struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	City        string    `gorm:"type:varchar(50);index;uniqueIndex:idx_city_date;not null;comment:'城市名'" json:"city"`
	TargetDate  string    `gorm:"type:varchar(20);index;uniqueIndex:idx_city_date;not null;comment:'预报日期(YYYY-MM-DD)'" json:"target_date"`
	Temperature string    `gorm:"type:varchar(20);comment:'温度范围'" json:"temperature"`
	Weather     string    `gorm:"type:varchar(50);comment:'天气情况(如:雷阵雨)'" json:"weather"`
	Direct      string    `gorm:"type:varchar(50);comment:'风向'" json:"direct"`
	UpdatedAt   time.Time `json:"updated_at"` // 记录最后一次从 API 拉取的时间
}
