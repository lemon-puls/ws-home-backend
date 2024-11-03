package api

import (
	"strings"
	"ws-home-backend/business"
	"ws-home-backend/common"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"
	"ws-home-backend/vo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// AddOrUpdateAlbum : 添加相册
// @Summary 添加相册
// @Description 添加相册
// @Tags 相册功能
// @Param body body dto.AlbumAddDTO true "相册信息"
// @Accept  json
// @Produce  json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album [post]
func AddOrUpdateAlbum(ctx *gin.Context) {
	DB := config.GetDB()
	var albumDto dto.AlbumAddDTO
	if err := ctx.ShouldBindJSON(&albumDto); err != nil {
		common.ErrorWithMsg(ctx, err.Error())
		return
	}
	var album model.Album
	// 新建
	if albumDto.Id == 0 {
		var user model.User
		res := DB.Take(&user, "user_id = ?", albumDto.UserId)
		if res.RowsAffected == 0 {
			common.ErrorWithMsg(ctx, "User not found")
			return
		}
		album.User = user
		err := copier.Copy(&album, &albumDto)
		if err != nil {
			common.ErrorWithMsg(ctx, err.Error())
			return
		}
	} else {
		// 更新
		DB.Take(&album, "id = ?", albumDto.Id)
		err := copier.CopyWithOption(&album, &albumDto, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		})
		if err != nil {
			common.ErrorWithMsg(ctx, err.Error())
			return
		}
		album.AlbumImgs = nil
	}
	res1 := DB.Save(&album)
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
// @Param limit query int true "每页数量"
// @Param order_by query string false "排序字段"
// @Param order query string false "排序方式"
// @Param name query string false "相册名称"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=[]vo.AlbumVO} "成功响应"
// @Router /album/list [get]
func ListAlbum(ctx *gin.Context) {
	var albumQueryDto dto.AlbumQueryDTO
	if err := ctx.ShouldBindQuery(&albumQueryDto); err != nil {
		common.ValidateError(ctx, err)
		return
	}
	pageRes := business.ListAlbum(albumQueryDto)
	albums, _ := pageRes.Records.(*[]model.Album)
	// 封装为 vo
	var albumVos []vo.AlbumVO
	for _, album := range *albums {
		var albumVo vo.AlbumVO
		copier.Copy(&albumVo, &album)
		albumVos = append(albumVos, albumVo)
	}

	pageRes.Records = &albumVos

	common.OkWithData(ctx, pageRes)
}

// AddImgToAlbum : 添加图片到相册
// @Summary 添加图片到相册
// @Description 添加图片到相册
// @Tags 相册功能
// @Param body body dto.AddImgToAlbumDTO true "图片信息"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album/img [post]
func AddImgToAlbum(ctx *gin.Context) {

	var addImgToAlbumDTO dto.AddImgToAlbumDTO
	if err := ctx.ShouldBindJSON(&addImgToAlbumDTO); err != nil {
		common.ErrorWithMsg(ctx, err.Error())
		return
	}

	business.AddImgToAlbum(addImgToAlbumDTO)

}

// RemoveImgFromAlbum : 从相册中移除图片
// @Summary 从相册中移除图片
// @Description 从相册中移除图片
// @Tags 相册功能
// @Param ids query string true "相册ID"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album/img [delete]
func RemoveImgFromAlbum(ctx *gin.Context) {
	ids := ctx.Query("ids")
	splits := strings.Split(ids, ",")
	business.RemoveImgFromAlbum(splits)
	common.OkWithMsg(ctx, "success")
}

// GetAlbumById : 获取相册详情
// @Summary 获取相册详情
// @Description 获取相册详情
// @Tags 相册功能
// @Param id path string true "相册ID"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=vo.AlbumVO} "成功响应"
// @Router /album/{id} [get]
func GetAlbumById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		common.ErrorWithCodeAndMsg(ctx, common.CodeInvalidParams, "id is required")
		return
	}
	album := business.GetAlbumById(id)
	var albumVo vo.AlbumVO
	copier.Copy(&albumVo, &album)
	common.OkWithData(ctx, albumVo)
}

// ListImgByAlbumId : 获取相册图片列表
// @Summary 获取相册图片列表
// @Description 获取相册图片列表
// @Tags 相册功能
// @Param body body dto.CursorListAlbumImgDTO true "查询条件"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=[]vo.AlbumImgVO} "成功响应"
// @Router /album/img/list [post]
func ListImgByAlbumId(ctx *gin.Context) {
	var queryRequest dto.CursorListAlbumImgDTO
	if err := ctx.ShouldBindJSON(&queryRequest); err != nil {
		common.ErrorWithCode(ctx, common.CodeInvalidParams)
		return
	}
	albumImgs := business.ListImgByAlbumId(queryRequest)
	common.OkWithData(ctx, albumImgs)
}
