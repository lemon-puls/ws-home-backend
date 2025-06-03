package dto

// GetPresignedURLReq 获取预签名URL请求
type GetPresignedURLReq struct {
	Type string `json:"type" binding:"required,oneof=upload download"` // 类型：upload-上传，download-下载
	Key  string `json:"key" binding:"required"`                        // 对象键
}

// BatchDeleteObjectsReq 批量删除对象请求
type BatchDeleteObjectsReq struct {
	Keys []string `json:"keys" binding:"required,min=1"` // 要删除的对象键列表
}
