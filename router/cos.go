package router

import (
	"github.com/gin-gonic/gin"
	"ws-home-backend/api"
	"ws-home-backend/middleware"
)

func RegisterCosRouter(router *gin.RouterGroup) {
	cosRouter := router.Group("/cos", middleware.LoginRequired())
	{
		cosRouter.GET("/credentials", api.GetTempCredentials)
	}
}
