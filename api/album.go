package api

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"ws-home-backend/business"
	"ws-home-backend/common"
	"ws-home-backend/common/cosutils"
	"ws-home-backend/common/page"
	"ws-home-backend/config"
	"ws-home-backend/config/db"
	"ws-home-backend/dto"
	"ws-home-backend/model"
	"ws-home-backend/vo"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"go.uber.org/zap"

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
	DB := db.GetDB()
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
		// 从上下文获取当前用户ID
		userId := ctx.GetInt64("userId")
		// 检查相册所有者是否为当前用户
		if album.UserId != userId {
			common.ErrorWithMsg(ctx, "您没有权限修改此相册")
			return
		}

		err := copier.CopyWithOption(&album, &albumDto, copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		})
		if err != nil {
			common.ErrorWithMsg(ctx, err.Error())
			return
		}
		album.Medias = nil
	}

	album.CoverImg = cosutils.ConvertUrlToKey(album.CoverImg)

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

	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	albumQueryDto.UserId = userId

	pageRes := business.ListAlbum(albumQueryDto)
	albums, _ := pageRes.Records.(*[]model.Album)

	cosClient := config.GetCosClient()

	// 封装为 vo
	var albumVos []vo.AlbumVO
	for _, album := range *albums {
		var albumVo vo.AlbumVO
		copier.Copy(&albumVo, &album)

		albumVo.CoverImg, _ = cosClient.GenerateDownloadPresignedURL(album.CoverImg)

		albumVos = append(albumVos, albumVo)
	}

	pageRes.Records = &albumVos

	common.OkWithData(ctx, pageRes)
}

// AddMediaToAlbum : 添加图片到相册
// @Summary 添加图片到相册
// @Description 添加图片到相册
// @Tags 相册功能
// @Param body body dto.AddMediaToAlbumDTO true "图片信息(包含url、大小等)"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=map[string]int64} "成功响应"
// @Router /album/media [post]
func AddMediaToAlbum(ctx *gin.Context) {
	var addMediaToAlbumDTO dto.AddMediaToAlbumDTO
	if err := ctx.ShouldBindJSON(&addMediaToAlbumDTO); err != nil {
		common.ErrorWithMsg(ctx, err.Error())
		return
	}
	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	// 检查相册所有者是否为当前用户
	album := business.GetAlbumById(strconv.FormatInt(addMediaToAlbumDTO.AlbumId, 10))
	if album.UserId != userId {
		common.ErrorWithMsg(ctx, "您没有权限修改此相册")
		return
	}

	urlToId := business.AddMediaToAlbum(addMediaToAlbumDTO)
	common.OkWithData(ctx, urlToId)
}

// RemoveMediaFromAlbum : 从相册中移除图片
// @Summary 从相册中移除图片
// @Description 从相册中移除图片
// @Tags 相册功能
// @Param ids query string true "相册ID"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album/media [delete]
func RemoveMediaFromAlbum(ctx *gin.Context) {
	ids := ctx.Query("ids")
	if ids == "" {
		common.ErrorWithMsg(ctx, "ids不能为空")
		return
	}

	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")

	splits := strings.Split(ids, ",")

	// 获取第一张图片所属的相册信息
	// TODO 这里假设所有图片都来自同一个相册，后续需要更严谨的鉴权再优化
	db := db.GetDB()
	var albumMedia model.AlbumMedia
	if err := db.Where("id = ?", splits[0]).First(&albumMedia).Error; err != nil {
		common.ErrorWithMsg(ctx, "图片不存在")
		return
	}

	// 查询相册信息
	var album model.Album
	if err := db.Where("id = ?", albumMedia.AlbumId).First(&album).Error; err != nil {
		common.ErrorWithMsg(ctx, "相册不存在")
		return
	}

	// 检查相册是否属于当前用户
	if album.UserId != userId {
		common.ErrorWithMsg(ctx, "您没有权限删除此图片")
		return
	}

	business.RemoveMediaFromAlbum(splits)
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
	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	album := business.GetAlbumById(id)
	// 检查相册所有者是否为当前用户
	if album.UserId != userId {
		common.ErrorWithMsg(ctx, "您没有权限查看此相册")
		return
	}
	var albumVo vo.AlbumVO
	copier.Copy(&albumVo, &album)

	albumVo.User.Avatar, _ = config.GetCosClient().GenerateDownloadPresignedURL(album.User.Avatar)

	albumVo.CoverImg, _ = config.GetCosClient().GenerateDownloadPresignedURL(album.CoverImg)

	common.OkWithData(ctx, albumVo)
}

// ListMediaByAlbumId : 获取相册图片列表
// @Summary 获取相册图片列表
// @Description 获取相册图片列表
// @Tags 相册功能
// @Param body body dto.CursorListAlbumMediaDTO true "查询条件"
// @Produce  json
// @Accept  json
// @Success 0 {object} common.Response{data=[]vo.AlbumMediaVO} "成功响应"
// @Router /album/media/list [post]
func ListMediaByAlbumId(ctx *gin.Context) {
	var queryRequest dto.CursorListAlbumMediaDTO
	if err := ctx.ShouldBindJSON(&queryRequest); err != nil {
		common.ErrorWithCode(ctx, common.CodeInvalidParams)
		return
	}
	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	// 检查相册所有者是否为当前用户
	album := business.GetAlbumById(strconv.FormatInt(queryRequest.AlbumId, 10))
	if album.UserId != userId {
		common.ErrorWithMsg(ctx, "您没有权限查看此相册")
		return
	}
	cursorPageVo := business.ListMediaByAlbumId(queryRequest)

	// 封装为 vo
	convertedCursorPageVo, err := page.ConvertCursorPageVO[model.AlbumMedia, vo.AlbumMediaVO](cursorPageVo)
	if err != nil {
		common.ErrorWithMsg(ctx, err.Error())
		return
	}
	for i, media := range cursorPageVo.Data {
		meta := media.Meta
		var metaDto *dto.MediaMetaDTO
		if meta != "" {
			err := json.Unmarshal([]byte(meta), &metaDto)
			if err != nil {
				zap.L().Error("media meta json unmarshal failed",
					zap.Int64("media_id", media.Id), zap.String("meta", meta), zap.Error(err))
			} else {
				convertedCursorPageVo.Data[i].Meta = metaDto
			}
		} else {
			convertedCursorPageVo.Data[i].Meta = nil
		}
		convertedCursorPageVo.Data[i].Url, _ = config.GetCosClient().GenerateDownloadPresignedURL(media.Url)
	}

	common.OkWithData(ctx, convertedCursorPageVo)
}

// DeleteAlbum : 删除相册
// @Summary 删除相册
// @Description 删除相册及其所有照片
// @Tags 相册功能
// @Param id path string true "相册ID"
// @Produce json
// @Accept json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album/{id} [delete]
func DeleteAlbum(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		common.ErrorWithCodeAndMsg(ctx, common.CodeInvalidParams, "id is required")
		return
	}

	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	album := business.GetAlbumById(id)

	if album.UserId == 0 {
		common.ErrorWithMsg(ctx, "相册不存在")
		return
	}

	// 检查相册所有者是否为当前用户
	if album.UserId != userId {
		common.ErrorWithMsg(ctx, "您没有权限删除此相册")
		return
	}

	business.DeleteAlbum(id)
	common.OkWithMsg(ctx, "删除成功")
}

// UpdateMediaSize : 更新所有图片大小
// @Summary 更新所有图片大小
// @Description 从COS获取并更新所有图片的实际大小(MB)
// @Tags 相册功能
// @Produce json
// @Accept json
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /album/media/size [post]
func UpdateMediaSize(ctx *gin.Context) {
	business.UpdateAllMediaSize()
	common.OkWithMsg(ctx, "更新成功")
}

// GetUserAlbumStats : 获取用户相册统计信息
// @Summary 获取用户相册统计信息
// @Description 获取用户的总相册数、总照片数、总照片大小
// @Tags 相册功能
// @Produce json
// @Accept json
// @Success 0 {object} common.Response{data=vo.AlbumStatsVO} "成功响应"
// @Router /album/stats [get]
func GetUserAlbumStats(ctx *gin.Context) {
	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	stats := business.GetUserAlbumStats(userId)
	common.OkWithData(ctx, stats)
}

// GetRandomMedia : 获取随机图片
// @Summary 获取随机图片
// @Description 从用户所有相册中随机选择5张图片
// @Tags 相册功能
// @Produce json
// @Accept json
// @Success 0 {object} common.Response{data=[]string} "成功响应"
// @Router /album/media/random [get]
func GetRandomMedia(ctx *gin.Context) {
	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")

	// 获取用户的所有相册
	DB := db.GetDB()
	var albums []model.Album
	if err := DB.Where("user_id = ?", userId).Find(&albums).Error; err != nil {
		common.ErrorWithMsg(ctx, "获取相册失败")
		return
	}

	// 获取所有相册ID
	albumIds := lo.Map(albums, func(album model.Album, _ int) int64 {
		return album.Id
	})

	// 批量查询所有图片
	var allMedia []model.AlbumMedia
	if err := DB.Where("album_id IN ?", albumIds).Find(&allMedia).Error; err != nil {
		common.ErrorWithMsg(ctx, "获取图片失败")
		return
	}

	// 如果图片数量不足5张，返回所有图片
	if len(allMedia) <= 5 {
		cosClient := config.GetCosClient()
		urls := lo.Map(allMedia, func(media model.AlbumMedia, _ int) string {
			url, _ := cosClient.GenerateDownloadPresignedURL(media.Url)
			return url
		})
		common.OkWithData(ctx, urls)
		return
	}

	// 随机选择5张图片
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allMedia), func(i, j int) {
		allMedia[i], allMedia[j] = allMedia[j], allMedia[i]
	})

	// 获取前5张图片的URL
	cosClient := config.GetCosClient()
	urls := lo.Map(allMedia[:5], func(media model.AlbumMedia, _ int) string {
		url, _ := cosClient.GenerateDownloadPresignedURL(media.Url)
		return url
	})

	common.OkWithData(ctx, urls)
}
