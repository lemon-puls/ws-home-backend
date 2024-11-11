package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"ws-home-backend/config"
	_ "ws-home-backend/docs"
	"ws-home-backend/router"
)

// @title WS Home Backend API
// @version 1.0
// @description 这是 WS Home Backend 的 API 文档
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
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
	r := router.InitRouter()
	// 初始化 雪花算法 ID 生成器
	config.InitSnowflakeNode(config.Conf.SnowflakeConfig)
	// 初始化 COS Client
	config.InitCosClient(config.Conf.CosConfig)

	r.Run(fmt.Sprintf(":%d", config.Conf.ServerConfig.Port))

	//config.RDB.Set(background, "test", "Go 使用 Redis", 0)
	zap.L().Info("Server exited", zap.String("profile", config.Conf.Profile))
}
