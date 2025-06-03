package api

import (
	"net/http"
	"time"
	"ws-home-backend/common"
	"ws-home-backend/common/cosutils"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/vo"

	"github.com/gin-gonic/gin"
)

// @Summary 获取预签名URL
// @Description 获取文件上传或下载的预签名URL
// @Tags 对象存储相关
// @Accept json
// @Produce json
// @Param data body dto.GetPresignedURLReq true "请求参数"
// @Success 200 {object} common.Response{data=vo.GetPresignedURLVO} "成功"
// @Router /cos/presigned-url [post]
func GetPresignedURL(ctx *gin.Context) {
	var req dto.GetPresignedURLReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ValidateError(ctx, err)
		return
	}

	cosClient := config.GetCosClient()

	// 默认过期时间 1 小时
	expire := 3600

	// 生成预签名URL
	var (
		url string
		err error
	)

	switch req.Type {
	case "upload":
		url, err = cosClient.GenerateUploadPresignedURL(req.Key)
	case "download":
		url, err = cosClient.GenerateDownloadPresignedURL(req.Key)
	default:
		common.ErrorWithCode(ctx, http.StatusBadRequest)
		return
	}

	if err != nil {
		common.ErrorWithCode(ctx, http.StatusInternalServerError)
		return
	}

	// 计算过期时间戳
	expireAt := time.Now().Add(time.Duration(expire) * time.Second).Unix()

	common.OkWithData(ctx, vo.GetPresignedURLVO{
		URL:      url,
		Key:      req.Key,
		ExpireAt: expireAt,
	})
}

// @Summary 批量删除对象
// @Description 批量删除指定的对象
// @Tags 对象存储相关
// @Accept json
// @Produce json
// @Param data body dto.BatchDeleteObjectsReq true "请求参数"
// @Success 200 {object} common.Response "成功"
// @Router /cos/batch-delete [post]
// TODO 该接口不够安全严谨，后续进行优化：对象的 key 路径中包含用户 id，用户只能删除自己的对象
func BatchDeleteObjects(ctx *gin.Context) {
	var req dto.BatchDeleteObjectsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ValidateError(ctx, err)
		return
	}

	// 处理每个 key，去除域名和协议（如果有）
	processedKeys := make([]string, len(req.Keys))
	for i, key := range req.Keys {
		processedKeys[i] = cosutils.ExtractKeyFromUrl(key)
	}

	cosClient := config.GetCosClient()
	if err := cosClient.DeleteObjects(processedKeys); err != nil {
		common.ErrorWithCode(ctx, http.StatusInternalServerError)
		return
	}

	common.Ok(ctx)
}
