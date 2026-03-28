package task

import (
	"fmt"
	"log"

	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
	"outdoor-app-backend/internal/service" // 调用第一期写的智能查库方法
	"outdoor-app-backend/pkg/weatherapi"

	"github.com/robfig/cron/v3"
)

func InitCronJobs() {
	c := cron.New()
	// 每天凌晨2点、下午2点执行一次（节省额度，全自动推警报）
	c.AddFunc("0 2,14 * * *", CheckActivityWeatherAlerts)
	InitActivityJobs(c)
	c.Start()

	log.Println("✅ 定时任务调度器已启动")
	log.Println("✅ 智能天气定时巡检任务 (Cron Jobs) 已启动")
}

// CheckActivityWeatherAlerts 巡检所有进行中的活动，判断是否需要发预警
func CheckActivityWeatherAlerts() {
	log.Println("🔄 开始执行定时任务：精准气象巡检...")

	var activities []model.Activity
	// 1. 找出所有【未结束】且【填写了城市】的活动
	database.DB.Where("status IN ('enrolling', 'full', 'ongoing') AND city != ''").Find(&activities)

	// 优化：记录本轮巡检中，哪些活动被触发了报警，避免重复发消息
	for _, act := range activities {
		// 获取活动的具体日期 (只取 YYYY-MM-DD，不包含时分秒)
		gatherDate := act.GatherTime.Format("2006-01-02")

		// 2. 🌟 核心：调用 Service 智能获取该城市的 5 天天气 (它会自动判断要不要去调第三方 API)
		cityWeathers, err := service.GetAndSyncWeather(act.City)
		if err != nil || len(cityWeathers) == 0 {
			log.Printf("获取城市 %s 天气失败，跳过", act.City)
			continue
		}

		// 3. 精准匹配日期：在这 5 天的数据里，找出和活动那天同一天的数据！
		var targetWeather *model.CityWeather
		for _, cw := range cityWeathers {
			if cw.TargetDate == gatherDate {
				targetWeather = &cw
				break
			}
		}

		// 如果 API 里没有那天的数据 (比如活动在 10 天后，API 只提供 5 天)，跳过不管
		if targetWeather == nil {
			continue
		}

		// 4. 🌟 判断这天的天气是否恶劣
		hasAlert := weatherapi.IsBadWeather(targetWeather.Weather)

		// 5. 如果状态发生了改变 (比如昨天预报是晴天(false)，今天突然预报那天下雨了(true))
		if act.WeatherAlert != hasAlert {
			// 更新主表状态，前端详情页会出现那个高亮的红色警告框
			database.DB.Model(&model.Activity{}).Where("id = ?", act.ID).Update("weather_alert", hasAlert)

			// 如果变成了恶劣天气，立即给所有报名的成员发一条极其具体的系统通知！
			if hasAlert {
				notifyMembers(act, targetWeather.TargetDate, targetWeather.Weather, targetWeather.Temperature)
			}
		}
	}
	log.Println("✅ 气象巡检任务执行完毕")
}

// notifyMembers 给所有已通过的人（包括发起人）发具体的危险天气推送
func notifyMembers(act model.Activity, date, weather, temp string) {
	// 🌟 1. 动态拼装的详尽话术
	content := fmt.Sprintf("您参与的活动【%s】目的地近期气象异常。%s (%s) 预报为【%s】，气温【%s℃】。请注意防范并带好相应的应急装备！",
		act.Title, act.City, date, weather, temp)

	// 🌟 2. 构造消息实体 (基础模板)
	baseMsg := model.Message{
		Type:      "system",
		Title:     "⚠️ 恶劣天气预警",
		RelatedID: act.ID, // 点击消息直接跳到活动详情审核页
		Content:   content,
	}

	// 🌟 3. 首先：给发起人发送一条预警！
	initiatorMsg := baseMsg
	initiatorMsg.UserID = act.InitiatorID
	repository.CreateMessage(&initiatorMsg)
	// ws.SendMessageToUser(act.InitiatorID, &initiatorMsg) // 实时推送

	// 🌟 4. 其次：给所有审核通过的成员发送预警
	var members []model.ActivityMember
	database.DB.Where("activity_id = ? AND status = 'approved'", act.ID).Find(&members)

	for _, m := range members {
		// 如果成员刚好就是发起人（虽然正常逻辑发起人不在 members 表里，但为了防呆加个判断）
		if m.UserID == act.InitiatorID {
			continue
		}

		memberMsg := baseMsg
		memberMsg.UserID = m.UserID
		repository.CreateMessage(&memberMsg)
		// ws.SendMessageToUser(m.UserID, &memberMsg) // 实时推送
	}
}
