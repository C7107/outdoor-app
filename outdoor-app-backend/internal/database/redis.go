package database

import (
	"context"
	"log"
	"outdoor-app-backend/configs"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background() // Redis v9 必须传 context

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     configs.AppConfig.Redis.Addr, // Redis 默认端口
		Password: "",                           // 默认没有密码
		DB:       0,                            // 默认使用 0 号数据库
	})

	// 测试连接
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis 连接失败: %v", err)
	}
	log.Println("✅ Redis 缓存服务连接成功！")
}
