package router

import (
	"github.com/gin-gonic/gin"
	"ws-home-backend/api"
)

func RegisterAlbumRouter(router *gin.RouterGroup) {
	albumRouter := router.Group("/album")
	{
		albumRouter.POST("/", api.AddAlbum)
		albumRouter.GET("/list", api.ListAlbum)
	}
}
