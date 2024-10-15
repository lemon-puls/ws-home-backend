package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var RDB *redis.Client

func InitRedis(cfg *RedisConfig, ctx context.Context) {

	// 实例化 redis 客户端
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// 验证连接
	if _, err := RDB.Ping(ctx).Result(); err != nil {
		zap.L().Error("redis connect failed", zap.Error(err))
		panic(err)
	}
	zap.L().Info("redis connect success")
	return
}
