package vo

import (
	"time"
	"ws-home-backend/model"
)

type AlbumVO struct {
	model.BaseModel
	Name        string       `json:"name"`
	Description string       `json:"description"`
	UserId      int64        `json:"user_id"`
	CoverImg    string       `json:"cover_img"`
	StartTime   time.Time    `json:"start_time"`
	User        UserVO       `json:"user"`
	AlbumImgs   []AlbumImgVO `json:"album_imgs"`
	PhotoCount  int64        `json:"photo_count"`
}

type AlbumImgVO struct {
	model.BaseModel
	AlbumId int64   `json:"album_id"`
	Url     string  `json:"url"`
	IsRaw   bool    `json:"is_raw"`
	Size    float64 `json:"size"`
}
