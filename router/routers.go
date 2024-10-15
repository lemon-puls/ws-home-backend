package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {

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

	//api := r.Group("/api")
	//{
	//	// 注册路由
	//}

	zap.L().Info("gin server start")
	return r
}
