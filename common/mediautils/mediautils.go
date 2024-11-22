package mediautils

import "strings"

const (
	MediaTypeImage = iota
	MediaTypeVideo
)

func GetMediaType(url string) int8 {
	ext := strings.ToLower(url[strings.LastIndex(url, ".")+1:])
	switch ext {
	case "mp4", "avi", "mov", "wmv", "flv":
		return MediaTypeVideo
	default:
		return MediaTypeImage
	}
}
