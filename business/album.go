package business

import (
	"ws-home-backend/common/page"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"
)

func ListAlbum(queryDto dto.AlbumQueryDTO) *page.PageResult {
	db := config.GetDB()
	albums := make([]model.Album, 0)
	query := db.Preload("User").
		Preload("AlbumImgs")

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
