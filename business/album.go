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
