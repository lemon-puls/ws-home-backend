package main

import (
	"go.uber.org/zap"
	"ws-home-backend/config"
)

func main() {
	// 初始化配置
	config.InitConfig("./config/config-dev.yaml")
	// 初始化日志
	config.InitLogger(config.Conf.LogConfig, config.Conf.Profile)
	zap.L().Info("Config initialized", zap.Any("config", config.Conf))
	// 连接数据库
	config.InitDB(config.Conf.MysqlConfig)

	zap.L().Info("Server started", zap.String("profile", config.Conf.Profile))
}
