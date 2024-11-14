package router

import (
	"ws-home-backend/api"
	"ws-home-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRouter(router *gin.RouterGroup) {
	albumRouter := router.Group("/album", middleware.LoginRequired())
	{
		albumRouter.POST("", api.AddOrUpdateAlbum)
		albumRouter.GET("/list", api.ListAlbum)
		albumRouter.POST("/img", api.AddImgToAlbum)
		albumRouter.DELETE("/img", api.RemoveImgFromAlbum)
		albumRouter.GET("/:id", api.GetAlbumById)
		albumRouter.POST("/img/list", api.ListImgByAlbumId)
		albumRouter.DELETE("/:id", api.DeleteAlbum)
		albumRouter.POST("/img/size", api.UpdateImgSize)
		albumRouter.GET("/stats", api.GetUserAlbumStats)
	}
}
