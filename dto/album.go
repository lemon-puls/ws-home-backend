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
	Url   string       `json:"url"`
	IsRaw bool         `json:"is_raw"`
	Size  float64      `json:"size"`
	Meta  MediaMetaDTO `json:"meta"`
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

// MediaMetaDTO 媒体元信息
type MediaMetaDTO struct {
	// 公共字段
	TakeTime  string `json:"takeTime,omitempty"`  // 拍摄时间
	Latitude  string `json:"latitude,omitempty"`  // 位置信息 纬度
	Longitude string `json:"longitude,omitempty"` // 位置信息 经度
	Address   string `json:"address,omitempty"`   // 位置信息 地址(通过经纬度获取)
	// 图片信息
	Make         string `json:"make,omitempty"`         // 相机品牌
	Model        string `json:"model,omitempty"`        // 相机型号
	ISO          int32  `json:"iso,omitempty"`          // ISO
	FNumber      string `json:"fNumber,omitempty"`      // 光圈
	ExposureTime string `json:"exposureTime,omitempty"` // 快门速度
	FocalLength  string `json:"focalLength,omitempty"`  // 焦距
	// 视频信息
	Duration   float64 `json:"duration,omitempty"`   // 视频时长（秒）
	Resolution string  `json:"resolution,omitempty"` // 分辨率
	Codec      string  `json:"codec,omitempty"`      // 编码
	Bitrate    float64 `json:"bitrate,omitempty"`    // 比特率
	FPS        float64 `json:"fps,omitempty"`        // 帧率
}
