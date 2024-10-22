package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"ws-home-backend/common"
	"ws-home-backend/common/jwt"
	"ws-home-backend/config"
)

const TokenKey = "Authorization"

// 登陆检验中间件
func LoginRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Token 格式： Authorization: Bearer <token>
		// 是否带了 Token
		token := ctx.Request.Header.Get(TokenKey)
		if token == "" {
			zap.L().Error("Token is missing")
			common.ErrorWithCode(ctx, common.CodeNotLogin)
			ctx.Abort()
			return
		}
		// Token 格式是否正确
		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			zap.L().Error("Token format is invalid")
			common.ErrorWithCode(ctx, common.CodeNotLogin)
			ctx.Abort()
			return
		}
		// 验证 token
		claims, err := jwt.VerifyToken(parts[1])
		if err != nil {
			zap.L().Error("Token is invalid")
			common.ErrorWithCode(ctx, common.CodeNotLogin)
			ctx.Abort()
			return
		}
		// 和 Redis 中存储的进行对比，是否一致，如果不一致，则表示已在其他地方登陆，禁止访问
		result, err := config.RDB.Get(context.Background(),
			common.GetUserTokenKey(claims.UserID, ctx.RemoteIP())).
			Result()
		if err != nil || result != parts[1] {
			if err != nil {
				zap.L().Error("Token is invalid", zap.Error(err))
			} else {
				zap.L().Error("User has logged in elsewhere",
					zap.Int64("user_id", claims.UserID), zap.String("ip", ctx.RemoteIP()))
			}
			common.ErrorWithCode(ctx, common.CodeNotLogin)
			ctx.Abort()
			return
		}
		// 验证通过，设置当前用户 ID 到上下文中
		ctx.Set("userId", claims.UserID)
		// 放行
		ctx.Next()
	}
}
