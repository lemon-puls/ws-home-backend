package vo

import "ws-home-backend/model"

type AlbumVO struct {
	model.BaseModel
	Name        string       `json:"name"`
	Description string       `json:"description"`
	UserId      int64        `json:"user_id"`
	CoverImg    string       `json:"cover_img"`
	User        UserVO       `json:"user"`
	AlbumImgs   []AlbumImgVO `json:"album_imgs"`
}

type AlbumImgVO struct {
	model.BaseModel
	AlbumId int64  `gorm:"not null;" json:"album_id"`
	Url     string `gorm:"type:varchar(255); not null;" json:"url"`
}
