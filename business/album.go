package business

import (
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"
)

func ListAlbum(queryDto dto.AlbumQueryDTO) []model.Album {
	db := config.GetDB()
	albums := make([]model.Album, 0)
	query := db.Preload("User").
		Preload("AlbumImgs").
		Scopes(config.Paginate(queryDto.Page, queryDto.PageSize))

	if queryDto.UserId != 0 {
		query = query.Where("user_id =?", queryDto.UserId)
	}
	if queryDto.Name != "" {
		query = query.Where("name like ?", "%"+queryDto.Name+"%")
	}

	query.Find(&albums)
	return albums
}
