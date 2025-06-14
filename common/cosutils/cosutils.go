package cosutils

import (
	"net/url"
	"strings"
)

// ConvertObjectPath 将文件路径转换为存储至数据库中的路径
// 例如：https://www.example.com/exampleobject/1745647348066-761.jpg?q-sign-algorithm=sha1&q-ak=AKIDc6MDsKXWGm38z432-7823gGhv9D4jANM7e094m
// 转换为：exampleobject/1745647348066-761.jpg
func ExtractKeyFromUrl(path string) string {
	if path == "" {
		return ""
	}
	// 如果路径不包含协议前缀，直接返回
	if !strings.Contains(path, "://") {
		return path
	}
	u, err := url.Parse(path)
	if err != nil {
		return path
	}
	// 如果路径为空，直接返回空字符串
	if u.Path == "" || u.Path == "/" {
		return ""
	}
	// 去掉开头的斜杠
	return u.Path[1:]
}
