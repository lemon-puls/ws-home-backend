package cosutils

import (
	"net/url"
	"strings"
)

// 从完整 URL 中提取对象键名
func ExtractKeyFromUrl(url string) string {
	// 假设 URL 格式为: https://bucket.cos.region.myqcloud.com/path/to/object
	parts := strings.Split(url, ".com/")
	if len(parts) != 2 {
		// 返回原始 URL
		return url
	}
	return parts[1]
}

// ConvertObjectPath 将文件路径转换为存储至数据库中的路径
// 例如：https://www.example.com/exampleobject/1745647348066-761.jpg?q-sign-algorithm=sha1&q-ak=AKIDc6MDsKXWGm38z432-7823gGhv9D4jANM7e094m
// 转换为：exampleobject/1745647348066-761.jpg
func ConvertUrlToKey(path string) string {
	if path == "" {
		return ""
	}
	u, _ := url.Parse(path)
	return u.Path[1:]
}
