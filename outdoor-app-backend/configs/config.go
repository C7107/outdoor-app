package configs

import (
	"log"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// AppConfig 暴露给外部访问的全局配置变量
var AppConfig *Config

// InitConfig 初始化并加载配置文件
func InitConfig() {
	viper.SetConfigName("config")     // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")       // 配置文件类型
	viper.AddConfigPath("./configs/") // 查找路径：相对于 main.go 运行的目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ 读取配置文件失败: %v", err)
	}

	// 将配置反序列化到结构体中
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("❌ 解析配置文件失败: %v", err)
	}

	log.Println("✅ 配置文件加载成功！")
}
