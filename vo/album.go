package vo

import (
	"time"
	"ws-home-backend/dto"
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
	PhotoCount  int64          `json:"photo_count"`
	VideoCount  int64          `json:"video_count"`
	TotalSize   float64        `json:"total_size"`
}

type AlbumMediaVO struct {
	model.BaseModel
	AlbumId int64             `json:"album_id"`
	Url     string            `json:"url"`
	Type    int8              `json:"type"`
	IsRaw   bool              `json:"is_raw"`
	Size    float64           `json:"size"`
	Meta    *dto.MediaMetaDTO `json:"meta"`
}

// AlbumStatsVO 相册统计信息
type AlbumStatsVO struct {
	TotalAlbums int64              `json:"totalAlbums"` // 总相册数
	TotalPhotos int64              `json:"totalPhotos"` // 总照片数
	TotalVideos int64              `json:"totalVideos"` // 总视频数
	TotalSize   float64            `json:"totalSize"`   // 总大小(MB)
	PhotoSize   float64            `json:"photoSize"`   // 图片总大小(MB)
	VideoSize   float64            `json:"videoSize"`   // 视频总大小(MB)
	Albums      []AlbumStatsItemVO `json:"albums"`      // 相册统计列表
}

// AlbumStatsItemVO 相册统计项
type AlbumStatsItemVO struct {
	Id         int64   `json:"id"`         // 相册ID
	Name       string  `json:"name"`       // 相册名称
	PhotoCount int64   `json:"photoCount"` // 图片数量
	VideoCount int64   `json:"videoCount"` // 视频数量
	PhotoSize  float64 `json:"photoSize"`  // 图片总大小(MB)
	VideoSize  float64 `json:"videoSize"`  // 视频总大小(MB)
	TotalSize  float64 `json:"totalSize"`  // 总大小(MB)
}
