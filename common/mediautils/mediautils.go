package mediautils

import "strings"

const (
	MediaTypeImage = iota
	MediaTypeVideo
)

func GetMediaType(url string) int8 {
	// 移除URL中的查询参数
	url = strings.Split(url, "?")[0]
	// 获取文件扩展名
	ext := strings.ToLower(url[strings.LastIndex(url, ".")+1:])
	switch ext {
	case "mp4", "avi", "mov", "wmv", "flv":
		return MediaTypeVideo
	default:
		return MediaTypeImage
	}
}
