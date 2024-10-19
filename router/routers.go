package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"time"
	"ws-home-backend/common"
	"ws-home-backend/config"
)

func InitRouter() *gin.Engine {

	if err := config.InitTranslator("zh"); err != nil {
		zap.L().Error("init translator failed", zap.Error(err))
	}

	r := gin.Default()

	// 设置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源，您可以指定特定的 URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400, // 缓存预请求的结果，单位是秒
	}))

	// 捕捉 panic 并记录日志
	r.Use(common.RecoveryWithZap(zap.L(), false))
	// 记录请求日志
	r.Use(common.LoggerWithZap(zap.L(), time.DateTime, false))

	api := r.Group("/api")
	{
		// 注册路由
		RegisterUserRouter(api)
	}

	// 配置 Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	zap.L().Info("gin server start")
	return r
}
