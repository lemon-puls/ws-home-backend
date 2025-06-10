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
		albumRouter.POST("/media", api.AddMediaToAlbum)
		albumRouter.DELETE("/media", api.RemoveMediaFromAlbum)
		albumRouter.GET("/:id", api.GetAlbumById)
		albumRouter.POST("/media/list", api.ListMediaByAlbumId)
		albumRouter.DELETE("/:id", api.DeleteAlbum)
		albumRouter.POST("/media/size", api.UpdateMediaSize)
		albumRouter.GET("/stats", api.GetUserAlbumStats)
		albumRouter.GET("/media/random", api.GetRandomMedia)
	}
}
