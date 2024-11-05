package router

import (
	"ws-home-backend/api"
	"ws-home-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/one", middleware.LoginRequired(), api.GetUserInfoById)
		userRouter.POST("/register", api.Register)
		userRouter.POST("/login", api.Login)
		userRouter.PUT("", middleware.LoginRequired(), api.UpdateUser)
		userRouter.GET("/current", middleware.LoginRequired(), api.GetCurrentUserInfo)
	}
}
