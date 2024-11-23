package vo

import (
	"time"
	"ws-home-backend/model"
)

type AlbumVO struct {
	model.BaseModel
	Name        string         `json:"name"`
	Description string         `json:"description"`
	UserId      int64          `json:"user_id"`
	CoverImg    string         `json:"cover_img"`
	StartTime   time.Time      `json:"start_time"`
	User        UserVO         `json:"user"`
	Medias      []AlbumMediaVO `json:"medias"`
	MediaCount  int64          `json:"media_count"`
	TotalSize   float64        `json:"total_size"`
}

type AlbumMediaVO struct {
	model.BaseModel
	AlbumId int64   `json:"album_id"`
	Url     string  `json:"url"`
	Type    int8    `json:"type"`
	IsRaw   bool    `json:"is_raw"`
	Size    float64 `json:"size"`
}

// AlbumStatsVO 相册统计信息
type AlbumStatsVO struct {
	TotalAlbums int64   `json:"totalAlbums"` // 总相册数
	TotalPhotos int64   `json:"totalPhotos"` // 总照片数
	TotalSize   float64 `json:"totalSize"`   // 总大小(MB)
	TotalVideos int64   `json:"totalVideos"` // 总视频数
}
