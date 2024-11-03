package business

import (
	"ws-home-backend/common/page"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"

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

func AddImgToAlbum(albumDTO dto.AddImgToAlbumDTO) {
	db := config.GetDB()

	albumImgs := make([]model.AlbumImg, 0)
	for _, url := range albumDTO.Urls {
		albumImg := model.AlbumImg{
			Url:     url,
			AlbumId: albumDTO.AlbumId,
		}
		albumImgs = append(albumImgs, albumImg)
	}
	if err := db.Create(&albumImgs).Error; err != nil {
		panic(err)
	}
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
	}, func(u *model.AlbumImg) interface{} {
		return &u.CreateTime
	})
	if err != nil {
		panic(err)
	}

	return result

}
