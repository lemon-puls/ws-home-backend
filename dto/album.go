package dto

import "ws-home-backend/common/page"

type AlbumAddDTO struct {
	UserId      int              `json:"user_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CoverImg    string           `json:"cover_img"`
	AlbumImgs   []AlbumImgAddDTO `json:"album_imgs"`
}

type AlbumImgAddDTO struct {
	Url string `json:"url"`
}

type AlbumQueryDTO struct {
	page.PageParam
	UserId int64  `json:"user_id" form:"user_id"`
	Name   string `json:"name" form:"name"`
}
