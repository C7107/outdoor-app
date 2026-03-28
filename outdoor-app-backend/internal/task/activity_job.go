package task

import (
	"log"
	"outdoor-app-backend/internal/repository"

	"github.com/robfig/cron/v3"
)

// InitActivityJobs 挂载到你的 InitCronJobs 中
func InitActivityJobs(c *cron.Cron) {
	// 每小时的第 0 分钟执行一次 (0 * * * *)
	// 你也可以设置每 10 分钟执行一次 (*/10 * * * *)
	c.AddFunc("0 * * * *", UpdateExpiredActivities)
	log.Println("✅ 活动状态自动流转任务已启动")
}

func UpdateExpiredActivities() {
	log.Println("🔄 开始执行定时任务：清理过期活动状态...")

	rows, err := repository.AutoUpdateActivityStatus()
	if err != nil {
		log.Printf("❌ 状态流转任务失败: %v", err)
		return
	}

	if rows > 0 {
		log.Printf("✅ 状态流转完毕！共将 %d 个过期活动标记为 '已结束(ended)'", rows)
	} else {
		log.Println("⚡ 没有需要过期的活动")
	}
}
