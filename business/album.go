package business

import (
	"context"
	"math"
	"sort"
	"ws-home-backend/common/cosutils"
	"ws-home-backend/common/maputils"
	"ws-home-backend/common/mediautils"
	"ws-home-backend/common/page"
	"ws-home-backend/config"
	"ws-home-backend/config/db"
	"ws-home-backend/dto"
	"ws-home-backend/model"
	"ws-home-backend/vo"

	"github.com/goccy/go-json"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ListAlbum(queryDto dto.AlbumQueryDTO) *page.PageResult {
	db := db.GetDB()
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

	// 查询这些相册的照片和视频数量
	var counts []struct {
		AlbumId    int64 `gorm:"column:album_id"`
		Type       int8  `gorm:"column:type"`
		MediaCount int64 `gorm:"column:media_count"`
	}

	db.Model(&model.AlbumMedia{}).
		Select("album_id, type, count(*) as media_count").
		Where("album_id IN ?", albumIds).
		Group("album_id, type").
		Find(&counts)

	// 构建相册ID到照片和视频数量的映射
	photoCountMap := make(map[int64]int64)
	videoCountMap := make(map[int64]int64)
	for _, count := range counts {
		if count.Type == mediautils.MediaTypeImage {
			photoCountMap[count.AlbumId] = count.MediaCount
		} else {
			videoCountMap[count.AlbumId] = count.MediaCount
		}
	}

	// 将照片和视频数量添加到相册对象中
	for i := range albums {
		albums[i].PhotoCount = photoCountMap[albums[i].Id]
		albums[i].VideoCount = videoCountMap[albums[i].Id]
	}

	paginate.Records = &albums
	return paginate
}

func AddMediaToAlbum(albumDTO dto.AddMediaToAlbumDTO) map[string]int64 {
	db := db.GetDB()

	medias := make([]model.AlbumMedia, 0)
	for _, media := range albumDTO.Medias {
		// 通过经纬度获取相应地址信息
		addressInfo, err := maputils.GetAddressFromLocation(media.Meta.Longitude, media.Meta.Latitude)
		if err != nil {
			zap.L().Error("Get address from location error", zap.Error(err))
		} else {
			media.Meta.Address = addressInfo.FormattedAddress
		}
		// 元数据转为 JSON 字符串落库
		metaJson, err := json.Marshal(media.Meta)
		if err != nil {
			zap.L().Error("Media meta json marshal error", zap.Int64("album_id", albumDTO.AlbumId),
				zap.String("url", media.Url), zap.Error(err))
		}

		mediaType := mediautils.GetMediaType(media.Url)

		albumMedia := model.AlbumMedia{
			Url:     cosutils.ConvertUrlToKey(media.Url),
			AlbumId: albumDTO.AlbumId,
			Type:    mediaType,
			IsRaw:   media.IsRaw,
			Size:    media.Size,
			Meta:    string(metaJson),
		}
		medias = append(medias, albumMedia)
	}

	if err := db.Create(&medias).Error; err != nil {
		panic(err)
	}

	urlToId := make(map[string]int64)
	for _, media := range medias {
		urlToId[media.Url] = media.Id
	}
	return urlToId
}

func RemoveMediaFromAlbum(splits []string) {
	db := db.GetDB()

	res := db.Where("id in (?)", splits).Delete(&model.AlbumMedia{})
	if res.Error != nil {
		panic(res.Error)
	}
}

func GetAlbumById(id string) *model.Album {
	db := db.GetDB()
	album := &model.Album{}
	if err := db.Preload("User").Take(album, id).Error; err != nil {
		panic(err)
	}

	// 单独查询相册的图片总大小
	var totalSize float64
	db.Model(&model.AlbumMedia{}).
		Select("ROUND(SUM(size), 2) as total_size").
		Where("album_id = ?", id).
		Scan(&totalSize)

	album.TotalSize = totalSize

	// 单独查询相册的图片数和视频数
	var photoCount, videoCount int64
	db.Model(&model.AlbumMedia{}).
		Where("album_id = ? AND type = ?", id, mediautils.MediaTypeImage).
		Count(&photoCount)

	db.Model(&model.AlbumMedia{}).
		Where("album_id = ? AND type = ?", id, mediautils.MediaTypeVideo).
		Count(&videoCount)

	album.PhotoCount = photoCount
	album.VideoCount = videoCount

	return album
}

func ListMediaByAlbumId(queryRequest dto.CursorListAlbumMediaDTO) *page.CursorPageBaseVO[model.AlbumMedia] {
	db := db.GetDB()
	result, err := page.GetCursorPageByMySQL(db, queryRequest.CursorPageBaseRequest, func(db *gorm.DB) {
		if queryRequest.AlbumId != 0 {
			db.Where("album_id = ?", queryRequest.AlbumId)
		}
		if queryRequest.IsRaw != nil {
			// 使用 *queryRequest.IsRaw 获取具体的布尔值
			db.Where("is_raw = ?", *queryRequest.IsRaw)
		}
		if queryRequest.Type != nil {
			db.Where("type = ?", *queryRequest.Type)
		}
	}, func(u *model.AlbumMedia) interface{} {
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
	db := db.GetDB()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 获取相册下的所有图片
	var albumImgs []model.AlbumMedia
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
		if err := config.GetCosClient().DeleteObjects(keys); err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	// 删除相册图片记录
	if err := tx.Where("album_id = ?", id).Delete(&model.AlbumMedia{}).Error; err != nil {
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

func UpdateAllMediaSize() {
	db := db.GetDB()
	originalCosClient := config.GetCosClient().GetOriginalClient()

	// 获取所有图片记录
	var albumImgs []model.AlbumMedia
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
		resp, err := originalCosClient.Object.Head(context.Background(), key, nil)
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
	db := db.GetDB()

	// 统计相册总数
	var totalAlbums int64
	if err := db.Model(&model.Album{}).Where("user_id = ?", userId).Count(&totalAlbums).Error; err != nil {
		zap.L().Error("统计相册总数失败", zap.Error(err))
		return &vo.AlbumStatsVO{}
	}

	// 统计照片总数和总大小
	type result struct {
		Type      int8    `json:"type"`
		Count     int64   `json:"count"`
		TotalSize float64 `json:"total_size"`
	}

	var results []result
	err := db.Model(&model.AlbumMedia{}).
		Select("ws_album_media.type as type, COUNT(*) as count, ROUND(SUM(size), 2) as total_size").
		Joins("JOIN ws_album ON ws_album_media.album_id = ws_album.id").
		Where("ws_album.user_id = ?", userId).
		Group("ws_album_media.type").
		Scan(&results).Error

	if err != nil {
		zap.L().Error("统计媒体数据失败", zap.Error(err))
		return &vo.AlbumStatsVO{TotalAlbums: totalAlbums}
	}

	// 如果没有任何媒体数据,直接返回
	if len(results) == 0 {
		return &vo.AlbumStatsVO{TotalAlbums: totalAlbums}
	}

	var totalPhotos, totalVideos int64
	var totalSize, photoSize, videoSize float64

	// 遍历结果并累加数据
	for _, r := range results {
		if r.Type == mediautils.MediaTypeImage {
			totalPhotos = r.Count
			photoSize = r.TotalSize
		} else {
			totalVideos = r.Count
			videoSize = r.TotalSize
		}
		totalSize += r.TotalSize
	}

	// 获取所有相册的统计信息
	type albumStats struct {
		Id        int64   `gorm:"column:id"`
		Name      string  `gorm:"column:name"`
		Type      int8    `gorm:"column:type"`
		Count     int64   `gorm:"column:count"`
		TotalSize float64 `gorm:"column:total_size"`
	}

	var albumStatsList []albumStats
	err = db.Model(&model.AlbumMedia{}).
		Select("ws_album.id, ws_album.name, ws_album_media.type, COUNT(*) as count, ROUND(SUM(size), 2) as total_size").
		Joins("JOIN ws_album ON ws_album_media.album_id = ws_album.id").
		Where("ws_album.user_id = ?", userId).
		Group("ws_album.id, ws_album.name, ws_album_media.type").
		Scan(&albumStatsList).Error

	if err != nil {
		zap.L().Error("统计相册数据失败", zap.Error(err))
		return &vo.AlbumStatsVO{
			TotalAlbums: totalAlbums,
			TotalPhotos: totalPhotos,
			TotalVideos: totalVideos,
			TotalSize:   totalSize,
			PhotoSize:   photoSize,
			VideoSize:   videoSize,
		}
	}

	// 按相册ID分组统计
	albumStatsMap := make(map[int64]*vo.AlbumStatsItemVO)
	for _, stat := range albumStatsList {
		if _, exists := albumStatsMap[stat.Id]; !exists {
			albumStatsMap[stat.Id] = &vo.AlbumStatsItemVO{
				Id:   stat.Id,
				Name: stat.Name,
			}
		}
		item := albumStatsMap[stat.Id]
		if stat.Type == mediautils.MediaTypeImage {
			item.PhotoCount = stat.Count
			item.PhotoSize = stat.TotalSize
		} else {
			item.VideoCount = stat.Count
			item.VideoSize = stat.TotalSize
		}
		item.TotalSize = item.PhotoSize + item.VideoSize
	}

	// 转换为切片并按总大小排序
	albums := make([]vo.AlbumStatsItemVO, 0, len(albumStatsMap))
	for _, item := range albumStatsMap {
		albums = append(albums, *item)
	}
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].TotalSize > albums[j].TotalSize
	})

	return &vo.AlbumStatsVO{
		TotalAlbums: totalAlbums,
		TotalPhotos: totalPhotos,
		TotalVideos: totalVideos,
		TotalSize:   totalSize,
		PhotoSize:   photoSize,
		VideoSize:   videoSize,
		Albums:      albums,
	}
}
