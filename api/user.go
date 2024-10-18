package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"ws-home-backend/business"
	"ws-home-backend/dto"
	"ws-home-backend/utils"
)

// GetUserInfoById : 获取用户详情
// @Summary 获取用户详情
// @Description 获取用户详情
// @Tags 用户模块
// @Produce json
// @Accept json
// @Param userId query string true "用户ID"
// @Success 0 {object} utils.Response{data=model.User} "成功响应"
// @Failure 3 {object} utils.Response "失败响应"
// @Router /user/one [get]
func GetUserInfoById(ctx *gin.Context) {

	value := ctx.Query("userId")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)

	user := business.GetUserById(int64(userId))

	if user.UserId == 0 {
		// user not found
		zap.L().Error("User not found", zap.Int32("userId", int32(userId)))
		utils.ErrorWithCode(ctx, utils.CodeNotFound)
		return
	}

	zap.L().Info("Get user info by id", zap.Any("user", user))
	utils.OkWithData(ctx, user)
}

// Register : 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户模块
// @Produce json
// @Accept json
// @Param body body dto.RegisterDTO true "用户注册信息"
// @Success 0 {object} utils.Response{data=string} "成功响应"
// @Router /user/register [post]
func Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	if err := ctx.ShouldBind(&registerDTO); err != nil {
		// 参数校验失败
		utils.ValidateError(ctx, err)
		return
	}
	// TODO 注册用户逻辑
	utils.OkWithMsg(ctx, "注册成功")
}
