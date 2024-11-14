package business

import (
	"context"
	"math"
	"ws-home-backend/common/cosutils"
	"ws-home-backend/common/page"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"
	"ws-home-backend/vo"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ListAlbum(queryDto dto.AlbumQueryDTO) *page.PageResult {
	db := config.GetDB()
	albums := make([]model.Album, 0)
	query := db.Preload("User")

	if queryDto.UserId != 0 {
		query = query.Where("user_id =?", queryDto.UserId)
	}
	if queryDto.Name != "" {
		query = query.Where("name like ?", "%"+queryDto.Name+"%")
	}
	if queryDto.StartTime != nil {
		query = query.Where("start_time >= ?", queryDto.StartTime)
	}
	if queryDto.EndTime != nil {
		query = query.Where("start_time <= ?", queryDto.EndTime)
	}

	paginate, err := page.Paginate(query, queryDto.PageParam, &albums)
	if err != nil {
		panic(err)
	}

	// 获取分页结果中所有相册的ID
	albumIds := make([]int64, 0)
	for _, album := range albums {
		albumIds = append(albumIds, album.Id)
	}

	// 查询这些相册的照片数量
	var counts []struct {
		AlbumId    int64 `gorm:"column:album_id"`
		PhotoCount int64 `gorm:"column:photo_count"`
	}

	db.Model(&model.AlbumImg{}).
		Select("album_id, count(*) as photo_count").
		Where("album_id IN ?", albumIds).
		Group("album_id").
		Find(&counts)

	// 构建相册ID到照片数量的映射
	photoCountMap := make(map[int64]int64)
	for _, count := range counts {
		photoCountMap[count.AlbumId] = count.PhotoCount
	}

	// 将照片数量添加到相册对象中
	for i := range albums {
		albums[i].PhotoCount = photoCountMap[albums[i].Id]
	}

	paginate.Records = &albums
	return paginate
}

func AddImgToAlbum(albumDTO dto.AddImgToAlbumDTO) map[string]int64 {
	db := config.GetDB()

	albumImgs := make([]model.AlbumImg, 0)
	for _, img := range albumDTO.AlbumImgs {
		albumImg := model.AlbumImg{
			Url:     img.Url,
			AlbumId: albumDTO.AlbumId,
			IsRaw:   img.IsRaw,
			Size:    img.Size,
		}
		albumImgs = append(albumImgs, albumImg)
	}
	if err := db.Create(&albumImgs).Error; err != nil {
		panic(err)
	}

	// 创建 URL 到 ID 的映射
	urlToId := make(map[string]int64)
	for _, img := range albumImgs {
		urlToId[img.Url] = img.Id
	}
	return urlToId
}

func RemoveImgFromAlbum(splits []string) {
	db := config.GetDB()

	res := db.Where("id in (?)", splits).Delete(&model.AlbumImg{})
	if res.Error != nil {
		panic(res.Error)
	}
}

func GetAlbumById(id string) *model.Album {
	db := config.GetDB()
	album := &model.Album{}
	if err := db.Preload("User").Preload("AlbumImgs").Take(album, id).Error; err != nil {
		panic(err)
	}
	return album
}

func ListImgByAlbumId(queryRequest dto.CursorListAlbumImgDTO) *page.CursorPageBaseVO[model.AlbumImg] {
	db := config.GetDB()
	result, err := page.GetCursorPageByMySQL(db, queryRequest.CursorPageBaseRequest, func(db *gorm.DB) {
		if queryRequest.AlbumId != 0 {
			db.Where("album_id = ?", queryRequest.AlbumId)
		}
		if queryRequest.IsRaw != nil {
			// 使用 *queryRequest.IsRaw 获取具体的布尔值
			db.Where("is_raw = ?", *queryRequest.IsRaw)
		}
	}, func(u *model.AlbumImg) interface{} {
		return &u.CreateTime
	})
	if err != nil {
		panic(err)
	}

	return result

}

/**
 * 删除相册
 */
func DeleteAlbum(id string) {
	db := config.GetDB()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 获取相册下的所有图片
	var albumImgs []model.AlbumImg
	if err := tx.Where("album_id = ?", id).Find(&albumImgs).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	// 删除 COS 上的图片
	if len(albumImgs) > 0 {
		var keys []string
		for _, img := range albumImgs {
			// 从完整 URL 中提取对象键名
			key := cosutils.ExtractKeyFromUrl(img.Url)
			keys = append(keys, key)
		}

		// 批量删除 COS 对象
		if err := cosutils.DeleteCosObjects(keys); err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	// 删除相册图片记录
	if err := tx.Where("album_id = ?", id).Delete(&model.AlbumImg{}).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	// 删除相册
	if err := tx.Delete(&model.Album{}, id).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	// 提交事务
	tx.Commit()
}

func UpdateAllImgSize() {
	db := config.GetDB()
	cosClient := config.GetCosClient()

	// 获取所有图片记录
	var albumImgs []model.AlbumImg
	if err := db.Find(&albumImgs).Error; err != nil {
		panic(err)
	}

	// 批量更新
	for _, img := range albumImgs {
		// 从URL中提取对象键名
		key := cosutils.ExtractKeyFromUrl(img.Url)
		if key == "" {
			continue
		}

		// 获取对象属性
		resp, err := cosClient.Object.Head(context.Background(), key, nil)
		if err != nil {
			zap.L().Error("获取对象属性失败",
				zap.String("key", key),
				zap.Error(err))
			continue
		}

		// 获取Content-Length
		size := resp.ContentLength
		// 转换为MB并保留两位小数
		sizeMB := float64(size) / 1024 / 1024
		sizeMB = math.Round(sizeMB*100) / 100

		// 更新数据库
		if err := db.Model(&img).Update("size", sizeMB).Error; err != nil {
			zap.L().Error("更新图片大小失败",
				zap.String("key", key),
				zap.Error(err))
			continue
		}
	}
}

func GetUserAlbumStats(userId int64) *vo.AlbumStatsVO {
	db := config.GetDB()

	// 统计相册总数
	var totalAlbums int64
	db.Model(&model.Album{}).Where("user_id = ?", userId).Count(&totalAlbums)

	// 统计照片总数和总大小
	var result struct {
		TotalPhotos int64   `json:"total_photos"`
		TotalSize   float64 `json:"total_size"`
	}

	db.Model(&model.AlbumImg{}).
		Select("COUNT(*) as total_photos, ROUND(SUM(size), 2) as total_size").
		Joins("JOIN ws_album ON ws_album_img.album_id = ws_album.id").
		Where("ws_album.user_id = ?", userId).
		Scan(&result)

	return &vo.AlbumStatsVO{
		TotalAlbums: totalAlbums,
		TotalPhotos: result.TotalPhotos,
		TotalSize:   result.TotalSize,
	}
}
