package api

import (
	"github.com/gin-gonic/gin"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"go.uber.org/zap"
	"time"
	"ws-home-backend/common"
	"ws-home-backend/config"
)

// GetTempCredentials : 获取临时密钥
// @Summary 获取临时密钥
// @Description 获取临时密钥
// @Tags COS 对象存储
// @Produce json
// @Accept json
// @SUCCESS 0 {object} common.Response "成功响应"
// @Router /cos/credentials [get]
func GetTempCredentials(ctx *gin.Context) {
	c := sts.NewClient(
		// 通过环境变量获取密钥, os.Getenv 方法表示获取环境变量
		config.Conf.AccessKey, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		config.Conf.SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		nil,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          config.Conf.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						// 这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						// 存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						//"qcs::cos:" + region + ":uid/" + appid + ":" + bucket + "/exampleobject",
						"*",
					},
					// 开始构建生效条件 condition
					// 关于 condition 的详细设置规则和COS支持的condition类型可以参考https://cloud.tencent.com/document/product/436/71306
					//Condition: map[string]map[string]interface{}{
					//	"ip_equal": map[string]interface{}{
					//		"qcs:ip": []string{
					//			"*",
					//		},
					//	},
					//},
				},
			},
		},
	}

	// 请求临时密钥
	res, err := c.GetCredential(opt)
	if err != nil {
		zap.L().Error("GetCredential failed", zap.Error(err))
		common.ErrorWithMsg(ctx, err.Error())
		return
	}
	// 返回临时密钥
	common.OkWithData(ctx, res)
}
