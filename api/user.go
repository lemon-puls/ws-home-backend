package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"ws-home-backend/business"
	"ws-home-backend/utils"
)

func GetUserInfoById(ctx *gin.Context) {

	value := ctx.Query("userId")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)

	user := business.GetUserById(int32(userId))

	if user.UserId == 0 {
		// user not found
		zap.L().Error("User not found", zap.Int32("userId", int32(userId)))
		utils.ErrorWithCode(ctx, utils.CodeNotFound)
		return
	}

	zap.L().Info("Get user info by id", zap.Any("user", user))
	utils.OkWithData(ctx, user)
}
