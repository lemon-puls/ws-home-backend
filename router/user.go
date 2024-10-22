package router

import (
	"github.com/gin-gonic/gin"
	"ws-home-backend/api"
	"ws-home-backend/middleware"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/one", middleware.LoginRequired(), api.GetUserInfoById)
		userRouter.POST("/register", api.Register)
		userRouter.POST("/login", api.Login)
	}
}
