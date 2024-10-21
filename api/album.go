package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"ws-home-backend/business"
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

// ListAlbum : 获取相册列表
// @Summary 获取相册列表
// @Description 获取相册列表
// @Tags 相册功能
// @Param user_id query string false "用户ID"
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param name query string false "相册名称"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=[]model.Album} "成功响应"
// @Router /album/list [get]
func ListAlbum(ctx *gin.Context) {
	var albumQueryDto dto.AlbumQueryDTO
	if err := ctx.ShouldBindQuery(&albumQueryDto); err != nil {
		common.ValidateError(ctx, err)
		return
	}
	albumList := business.ListAlbum(albumQueryDto)
	common.OkWithData(ctx, albumList)
}
