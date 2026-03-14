package service

import (
	"encoding/json"
	"errors"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
)

// GetRouteList 获取线路列表
func GetRouteList(req *dto.RouteFilterReq, page, pageSize int) ([]model.Route, error) {
	return repository.GetRouteList(req.City, req.Difficulty, req.MaxDuration, req.MinScore, page, pageSize)
}

// GetRouteDetail 获取聚合的单条线路详情
func GetRouteDetail(userID, routeID uint) (*dto.RouteDetailRes, error) {
	// 1. 查线路基本信息
	route, err := repository.GetRouteByID(routeID)
	if err != nil {
		return nil, errors.New("线路不存在")
	}

	// 2. 查是否被当前用户收藏
	isFavorited := repository.IsFavorited(userID, routeID)

	// 3. 查评价列表，并把 model 转换为 dto.RouteReviewRes
	reviewsData, _ := repository.GetReviewsByRoute(routeID)
	var reviewsRes []dto.RouteReviewRes

	for _, r := range reviewsData {
		var images []string
		json.Unmarshal([]byte(r.Images), &images) // 解析图片数组

		reviewsRes = append(reviewsRes, dto.RouteReviewRes{
			ID:        r.ID,
			Score:     r.Score,
			Content:   r.Content,
			Images:    images,
			CreatedAt: r.CreatedAt,
			User: dto.UserBasicInfo{
				ID:       r.User.ID,
				Nickname: r.User.Nickname,
				Avatar:   r.User.Avatar,
			},
		})
	}

	// 4. TODO: 调用第三方天气 API (先留空，后续再接入)
	// weatherData := weatherapi.GetWeather(route.City)

	// 5. 组装返回
	res := &dto.RouteDetailRes{
		ID:           route.ID,
		Title:        route.Title,
		City:         route.City,
		Difficulty:   route.Difficulty,
		DurationDays: route.DurationDays,
		SceneryScore: route.SceneryScore,
		CoverImage:   route.CoverImage,
		MapTrackUrl:  route.MapTrackUrl,
		ElevationUrl: route.ElevationUrl,
		Description:  route.Description,
		IsFavorited:  isFavorited,
		Reviews:      reviewsRes,
	}

	return res, nil
}

// ToggleFavorite 收藏/取消收藏
func ToggleFavorite(userID, routeID uint) error {
	if repository.IsFavorited(userID, routeID) {
		return repository.RemoveFavorite(userID, routeID)
	}
	return repository.AddFavorite(userID, routeID)
}

// CreateReview 发布评价
func CreateReview(userID, routeID uint, req *dto.ReviewCreateReq) error {
	imagesJSON, _ := json.Marshal(req.Images)

	review := &model.RouteReview{
		RouteID: routeID,
		UserID:  userID,
		Score:   req.Score,
		Content: req.Content,
		Images:  string(imagesJSON),
	}
	return repository.CreateReview(review)
}

// CreateRoute 发布一条新线路（给后台用）
func CreateRoutes(req []dto.RouteCreateReq) error {

	var routes []model.Route

	for _, r := range req {
		route := model.Route{
			Title:        r.Title,
			City:         r.City,
			Difficulty:   r.Difficulty,
			DurationDays: r.DurationDays,
			SceneryScore: r.SceneryScore,
			CoverImage:   r.CoverImage,
			MapTrackUrl:  r.MapTrackUrl,
			ElevationUrl: r.ElevationUrl,
			Description:  r.Description,
		}

		routes = append(routes, route)
	}

	return repository.CreateRoutes(routes)
}
