package service

import (
	"log"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
	"outdoor-app-backend/pkg/weatherapi"
	"time"
)

// GetAndSyncWeather 智能获取天气 (带自动按需同步机制)
func GetAndSyncWeather(city string) ([]model.CityWeather, error) {
	// 1. 先去数据库查
	dbWeathers, err := repository.GetCityWeatherFromDB(city)

	needSync := false

	// 2. 核心判断：需要去调第三方 API 吗？
	// 条件 A: 数据库里连今天的数据都没有（可能这个城市是第一次被搜）
	if len(dbWeathers) == 0 || dbWeathers[0].TargetDate != time.Now().Format("2006-01-02") {
		needSync = true
	} else {
		// 条件 B: 数据库里有数据，但是上次更新时间距今超过了 4 小时 (缓存过期了)
		lastUpdate := dbWeathers[0].UpdatedAt
		if time.Since(lastUpdate).Hours() > 4 {
			needSync = true
		}
	}

	// 3. 如果不需要同步，直接把查到的库里数据返回给前端！(耗时 < 5ms，0 费用)
	if !needSync {
		return dbWeathers, nil
	}

	// ==========================================
	// 4. 执行按需同步 (去调第三方 API)
	// ==========================================
	log.Printf("☁️ 城市 %s 的天气数据已过期或不存在，正在请求第三方 API...", city)

	apiData, err := weatherapi.FetchFromJuhe(city) // 我们稍后在 pkg 里实现这个真实的请求
	if err != nil {
		// 如果 API 挂了或者欠费了，为了不让前端白屏，我们容错降级，返回旧的数据库数据
		log.Printf("⚠️ 调用天气 API 失败: %v, 降级使用旧数据", err)
		return dbWeathers, nil
	}

	// 5. 把第三方拿到的数据转为数据库 Model
	var newWeathers []model.CityWeather
	now := time.Now()

	for _, f := range apiData.Result.Future {
		newWeathers = append(newWeathers, model.CityWeather{
			City:        city,
			TargetDate:  f.Date,
			Temperature: f.Temperature,
			Weather:     f.Weather,
			Direct:      f.Direct,
			UpdatedAt:   now,
		})
	}

	// 6. 异步批量入库！(不阻塞当前前端的请求，极速返回)
	go func(ws []model.CityWeather) {
		if err := repository.BatchUpsertWeather(ws); err != nil {
			log.Printf("❌ 天气数据入库失败: %v", err)
		}
	}(newWeathers)

	// 7. 直接把最新的数据返回给前端
	return newWeathers, nil
}
