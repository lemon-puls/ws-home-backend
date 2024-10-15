package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"ws-home-backend/config"
	"ws-home-backend/router"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "cfg", "./config/config.yaml", "配置文件路径")
	flag.Parse()

	background := context.Background()

	// 初始化配置
	config.InitConfig("./config/config-dev.yaml")
	// 初始化日志
	config.InitLogger(config.Conf.LogConfig, config.Conf.Profile)
	zap.L().Info("Config initialized", zap.Any("config", config.Conf))
	// 连接 Mysql 数据库
	config.InitDB(config.Conf.MysqlConfig)
	// 初始化 Redis 连接
	config.InitRedis(config.Conf.RedisConfig, background)
	// 初始化 Gin Router
	router.InitRouter()

	//config.RDB.Set(background, "test", "Go 使用 Redis", 0)
	zap.L().Info("Server started", zap.String("profile", config.Conf.Profile))
}
