package database

import (
	"log"
	"os"

	"outdoor-app-backend/internal/model" // ⚠️ 替换为你的真实项目模块名

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接对象，供 Repo 层调用
var DB *gorm.DB

// InitMySQL 初始化数据库连接并执行迁移
// 参数 isSeed: 是否执行 init.sql 导入初始数据
func InitMySQL(isSeed bool) {
	// 1. 配置 DSN (Data Source Name)
	// 格式: 用户名:密码@tcp(IP:端口)/数据库名?参数...
	// ⚠️ 请确保你已经在 MySQL 里提前建好了 outdoor_db 这个空数据库！
	dsn := "root:123456@tcp(127.0.0.1:3306)/outdoor_db?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"

	// 2. 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 在终端打印生成的 SQL 语句，方便调试，logger.Info	打印所有 SQL，logger.Warn 只打印慢，SQLlogger.Error	只打印错误，logger.Silent 完全不打印 SQL
	})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	DB = db
	log.Println("✅ MySQL 数据库连接成功！")

	// 3. 自动迁移表结构 (AutoMigrate)
	log.Println("🔄 正在同步数据库表结构...")
	err = db.AutoMigrate(
		&model.User{},
		&model.Route{}, &model.RouteReview{}, &model.FavoriteRoute{},
		&model.Activity{}, &model.ActivityMember{},
		&model.Post{}, &model.Comment{}, &model.PostLike{},
		&model.Article{},
		&model.Message{},
	)
	if err != nil {
		log.Fatalf("❌ 表结构同步失败: %v", err)
	}
	log.Println("✅ 数据库表结构同步完成！")

	// 4. 判断是否需要导入初始测试数据
	if isSeed {
		runInitSQL(db)
	}
}

// runInitSQL 读取并执行 migrations/init.sql
func runInitSQL(db *gorm.DB) {
	log.Println("🚀 接收到 -seed 指令，准备导入初始测试数据...")

	// 读取 sql 文件
	sqlBytes, err := os.ReadFile("./migrations/init.sql")
	if err != nil {
		log.Printf("⚠️ 读取 init.sql 失败 (请检查 ./migrations/init.sql 文件是否存在): %v", err)
		return
	}

	// 执行 sql 语句
	err = db.Exec(string(sqlBytes)).Error
	if err != nil {
		log.Printf("❌ 导入初始数据失败: %v", err)
	} else {
		log.Println("🎉 初始数据导入成功！测试环境已准备就绪。")
	}
}
