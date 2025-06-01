package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"ws-home-backend/config"
	"ws-home-backend/config/db"
	"ws-home-backend/config/logging"
	_ "ws-home-backend/docs"
	"ws-home-backend/email"
	"ws-home-backend/router"

	"go.uber.org/zap"
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
	logging.InitLogger(config.Conf.LogConfig, config.Conf.Profile)
	zap.L().Info("Config initialized", zap.Any("config", config.Conf))
	// 连接 Mysql 数据库
	db.InitDB(config.Conf.MysqlConfig)
	// 初始化 Redis 连接
	config.InitRedis(config.Conf.RedisConfig, background)
	// 初始化 Gin Router
	r := router.InitRouter()
	// 初始化 雪花算法 ID 生成器
	config.InitSnowflakeNode(config.Conf.SnowflakeConfig)
	// 初始化 COS Client
	config.InitCosClient(config.Conf.CosConfig)

	// 初始化并启动定时问候邮件任务
	morningGreeting := email.NewMorningGreeting(config.Conf.EmailConfig)
	morningGreeting.StartScheduler()
	zap.L().Info("Morning greeting scheduler started")

	// 创建一个通道来接收系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", config.Conf.ServerConfig.Port)); err != nil {
			zap.L().Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// 等待中断信号
	<-quit
	zap.L().Info("Shutting down server...")

	// 停止定时任务
	morningGreeting.StopScheduler()

	zap.L().Info("Server exited", zap.String("profile", config.Conf.Profile))
}
