package router

import (
	"github.com/gin-gonic/gin"
	"ws-home-backend/api"
	"ws-home-backend/middleware"
)

func RegisterAlbumRouter(router *gin.RouterGroup) {
	albumRouter := router.Group("/album", middleware.LoginRequired())
	{
		albumRouter.POST("/", api.AddAlbum)
		albumRouter.GET("/list", api.ListAlbum)
		albumRouter.POST("/img", api.AddImgToAlbum)
		albumRouter.DELETE("/img", api.RemoveImgFromAlbum)
		albumRouter.GET("/:id", api.GetAlbumById)
	}
}
