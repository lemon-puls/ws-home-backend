package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"ws-home-backend/common"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"
)

// AddAlbum : 添加相册
// @Summary 添加相册
// @Description 添加相册
// @Tags 相册功能
// @Param body body dto.AlbumAddDTO true "相册信息"
// @Accept  json
// @Produce  json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album [post]
func AddAlbum(ctx *gin.Context) {
	DB := config.GetDB()
	var albumDto dto.AlbumAddDTO
	if err := ctx.ShouldBindJSON(&albumDto); err != nil {
		common.ErrorWithMsg(ctx, err.Error())
		return
	}
	var user model.User
	res := DB.Take(&user, "user_id = ?", albumDto.UserId)
	if res.RowsAffected == 0 {
		common.ErrorWithMsg(ctx, "User not found")
		return
	}
	var album model.Album
	err := copier.Copy(&album, &albumDto)
	if err != nil {
		common.ErrorWithMsg(ctx, err.Error())
		return
	}
	album.User = user
	res1 := DB.Create(&album)
	if res1.RowsAffected == 0 {
		common.ErrorWithMsg(ctx, "Failed to create album")
		return
	}
	common.OkWithData(ctx, album.Id)
}
