package repository

import (
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"
)

// GetRouteList 多维度筛选线路
func GetRouteList(city string, difficulty, maxDuration int, minScore float64, page, pageSize int) ([]model.Route, error) {
	var routes []model.Route
	query := database.DB.Model(&model.Route{})

	// 动态拼接查询条件 (如果前端传了对应的值，才加上这个条件)
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if difficulty > 0 {
		query = query.Where("difficulty = ?", difficulty)
	}
	if maxDuration > 0 {
		query = query.Where("duration_days <= ?", maxDuration) // 筛选耗时小于等于 N 天的
	}
	if minScore > 0 {
		query = query.Where("scenery_score >= ?", minScore) // 筛选风景评分大于等于 N 的
	}

	// 瀑布流展示，通常按照推荐度（风景评分）或最新创建排序
	err := query.Order("scenery_score DESC, created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&routes).Error

	return routes, err
}

// GetRouteByID 根据 ID 查询单条路线详情。
func GetRouteByID(id uint) (*model.Route, error) {
	var route model.Route
	err := database.DB.First(&route, id).Error
	return &route, err
}

// CreateReview 创建一条新的路线评价记录。
func CreateReview(review *model.RouteReview) error {
	return database.DB.Create(review).Error
}

// GetReviewsByRoute 获取指定路线的所有评价，并预加载每个评价的用户信息。
func GetReviewsByRoute(routeID uint) ([]model.RouteReview, error) {
	var reviews []model.RouteReview
	err := database.DB.Preload("User").Where("route_id = ?", routeID).Find(&reviews).Error
	return reviews, err
}

//.Preload("User")
//作用：告诉 GORM 在查询评价的同时，预加载每条评价关联的用户信息。这里 "User" 对应 RouteReview 模型中定义的 User 字段（类型为 model.User）。
//原理：GORM 会先执行主查询（查询评价），然后根据主查询结果中的 UserID 字段，自动发起额外的查询（或使用 JOIN）来获取对应的用户信息，并将结果填充到 reviews 中每个元素的 User 字段中。
//好处：如果没有 Preload，要获取评价及其用户信息，你可能需要先查询评价列表，然后遍历每个评价再单独查询用户，导致 N+1 次查询。Preload 将查询次数减少到 2 次（或 1 次，如果使用 JOIN 策略），大幅提升性能。

func AddFavorite(userID, routeID uint) error {
	favorite := model.FavoriteRoute{
		UserID:  userID,
		RouteID: routeID,
	}
	return database.DB.Create(&favorite).Error
}

// RemoveFavorite 取消收藏
func RemoveFavorite(userID, routeID uint) error {
	return database.DB.Where("user_id = ? AND route_id = ?", userID, routeID).
		Delete(&model.FavoriteRoute{}).Error
}

// IsFavorited 检查用户是否已收藏该线路 (前端展示“心形”图标点亮状态时用)
func IsFavorited(userID, routeID uint) bool {
	var count int64
	database.DB.Model(&model.FavoriteRoute{}).
		Where("user_id = ? AND route_id = ?", userID, routeID).
		Count(&count)
	return count > 0
}

// GetUserFavorites 获取用户收藏的所有线路列表
// 这里使用了 Preload，实现“通过中间表查询出具体的路线信息”
func GetUserFavorites(userID uint) ([]model.Route, error) {
	var routes []model.Route

	// 这里通过 JOIN 关联查询，查出 UserID 对应的所有 Route 记录
	err := database.DB.Table("routes").
		Joins("JOIN favorite_routes ON favorite_routes.route_id = routes.id").
		Where("favorite_routes.user_id = ?", userID).
		Find(&routes).Error

	return routes, err
}

// CreateRoute 创建一条新的公共路线 (后台或专家用)
func CreateRoutes(routes []model.Route) error {
	return database.DB.Create(&routes).Error
}
