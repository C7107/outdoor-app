package repository

import (
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"
	"time"

	"gorm.io/gorm/clause"
)

// GetCityWeatherFromDB 从数据库获取某城市未来 5 天的天气
func GetCityWeatherFromDB(city string) ([]model.CityWeather, error) {
	var weathers []model.CityWeather

	// 只查今天及以后的数据 (防止查出过期的历史天气)
	today := time.Now().Format("2006-01-02")

	err := database.DB.Where("city = ? AND target_date >= ?", city, today).
		Order("target_date ASC").
		Limit(5).
		Find(&weathers).Error

	return weathers, err
}

// BatchUpsertWeather 批量更新或插入天气数据 (防数据库打爆的核心)
// 使用 GORM 的 Clauses(clause.OnConflict) 实现 MySQL 的 UPSERT
func BatchUpsertWeather(weathers []model.CityWeather) error {
	if len(weathers) == 0 {
		return nil
	}

	// 无论传进来多少条，只用一条 SQL 语句批量写入！极大降低数据库 I/O 压力
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "city"}, {Name: "target_date"}},                               // 冲突条件：城市和日期组合重复，当 city 和 target_date 的组合已存在时，视为冲突
		DoUpdates: clause.AssignmentColumns([]string{"temperature", "weather", "direct", "updated_at"}), // 重复则更新这些字段，当冲突发生时，更新这些字段
	}).Create(&weathers).Error
}
