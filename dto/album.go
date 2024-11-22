package dto

import (
	"time"
	"ws-home-backend/common/page"
)

type AlbumAddDTO struct {
	Id          int64              `json:"id"`
	UserId      int64              `json:"user_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CoverImg    string             `json:"cover_img"`
	StartTime   time.Time          `json:"start_time"`
	Medias      []AlbumMediaAddDTO `json:"medias"`
}

type AlbumMediaAddDTO struct {
	Url   string  `json:"url"`
	IsRaw bool    `json:"is_raw"`
	Size  float64 `json:"size"`
}

type AlbumQueryDTO struct {
	page.PageParam
	UserId    int64      `json:"user_id" form:"user_id"`
	Name      string     `json:"name" form:"name"`
	StartTime *time.Time `json:"start_time" form:"start_time"`
	EndTime   *time.Time `json:"end_time" form:"end_time"`
}

type AddMediaToAlbumDTO struct {
	AlbumId int64              `json:"album_id"`
	Medias  []AlbumMediaAddDTO `json:"medias"`
}

type CursorListAlbumMediaDTO struct {
	page.CursorPageBaseRequest
	AlbumId int64 `json:"album_id"`
	IsRaw   *bool `json:"is_raw"`
	Type    *int8 `json:"type"`
}
