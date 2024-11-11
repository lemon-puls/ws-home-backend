package cosutils

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"strings"
	"ws-home-backend/config"
)

/**
 * 批量删除 COS 对象
 */
func DeleteCosObjects(keys []string) error {

	client := config.GetCosClient()
	// 构建删除对象请求
	obs := []cos.Object{}
	for _, key := range keys {
		obs = append(obs, cos.Object{Key: key})
	}

	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
		// 布尔值，这个值决定了是否启动 Quiet 模式
		// 值为 true 启动 Quiet 模式，值为 false 则启动 Verbose 模式，默认值为 false
		Quiet: true,
	}

	// 执行批量删除
	_, _, err := client.Object.DeleteMulti(context.Background(), opt)
	return err
}

// 从完整 URL 中提取对象键名
func ExtractKeyFromUrl(url string) string {
	// 假设 URL 格式为: https://bucket.cos.region.myqcloud.com/path/to/object
	parts := strings.Split(url, ".com/")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
