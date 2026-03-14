package main

import (
	"flag"
	"log"
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/router" // 导入路由模块
)

func main() {
	seed := flag.Bool("seed", false, "是否执行 init.sql 导入初始数据")
	flag.Parse()

	log.Println("🏃‍♂️ 正在启动 Outdoor App 后端服务...")

	// 1. 初始化数据库
	database.InitMySQL(*seed)

	// 2. 初始化路由
	r := router.InitRouter()

	// 3. 启动服务 (Gin 会接管程序运行，不需要 select {})
	log.Println("🌐 后端服务启动成功，监听端口 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ 服务启动失败: %v", err)
	}
}
