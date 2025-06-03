package router

import (
	"ws-home-backend/api"
	"ws-home-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCosRouter(router *gin.RouterGroup) {
	cosRouter := router.Group("/cos", middleware.LoginRequired())
	{
		cosRouter.POST("/presigned-url", api.GetPresignedURL)
		cosRouter.POST("/batch-delete", api.BatchDeleteObjects)
	}
}
