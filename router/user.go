package router

import (
	"github.com/gin-gonic/gin"
	"ws-home-backend/api"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/one", api.GetUserInfoById)
		userRouter.POST("/register", api.Register)
	}
}
