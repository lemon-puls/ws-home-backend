package router

import (
	"github.com/gin-gonic/gin"
	"ws-home-backend/api"
)

func RegisterCosRouter(router *gin.RouterGroup) {
	cosRouter := router.Group("/cos")
	{
		cosRouter.GET("/credentials", api.GetTempCredentials)
	}
}
