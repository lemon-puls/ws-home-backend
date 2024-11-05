package api

import (
	"fmt"
	"strconv"
	"ws-home-backend/business"
	"ws-home-backend/common"
	"ws-home-backend/dto"
	"ws-home-backend/vo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

// GetUserInfoById : 获取用户详情
// @Summary 获取用户详情
// @Description 获取用户详情
// @Tags 用户模块
// @Produce json
// @Accept json
// @Param userId query string true "用户ID"
// @Success 0 {object} common.Response{data=model.User} "成功响应"
// @Failure 3 {object} common.Response "失败响应"
// @Router /user/one [get]
func GetUserInfoById(ctx *gin.Context) {

	value := ctx.Query("userId")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)

	user := business.GetUserById(int64(userId))

	if user.UserId == 0 {
		// user not found
		zap.L().Error("User not found", zap.Int32("userId", int32(userId)))
		common.ErrorWithCode(ctx, common.CodeNotFound)
		return
	}

	zap.L().Info("Get user info by id", zap.Any("user", user))
	common.OkWithData(ctx, user)
}

// Register : 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户模块
// @Produce json
// @Accept json
// @Param body body dto.RegisterDTO true "用户注册信息"
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /user/register [post]
func Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	if err := ctx.ShouldBind(&registerDTO); err != nil {
		// 参数校验失败
		common.ValidateError(ctx, err)
		return
	}
	// 注册用户逻辑
	business.Register(registerDTO)
	common.OkWithMsg(ctx, "注册成功")
}

// Login : 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户模块
// @Produce json
// @Accept json
// @Param body body dto.LoginDTO true "用户登录信息"
// @Success 0 {object} common.Response{data=vo.Tokens} "成功响应"
// @Router /user/login [post]
func Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	if err := ctx.ShouldBindJSON(&loginDTO); err != nil {
		common.ValidateError(ctx, err)
		return
	}
	token := business.Login(loginDTO, ctx)
	common.OkWithData(ctx, token)
}

// UpdateUser : 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户信息(包括密码)
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param body body dto.UpdateUserDTO true "用户信息"
// @Success 0 {object} common.Response{data=string} "成功响应"
// @Router /user [put]
func UpdateUser(ctx *gin.Context) {
	var updateUserDTO dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&updateUserDTO); err != nil {
		common.ValidateError(ctx, err)
		return
	}

	// 从上下文获取当前用户ID
	userId := ctx.GetInt64("userId")
	business.UpdateUser(userId, updateUserDTO)
	common.OkWithMsg(ctx, "更新成功")
}

// GetCurrentUserInfo : 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户模块
// @Produce json
// @Success 0 {object} common.Response{data=vo.UserVO} "成功响应"
// @Failure 3 {object} common.Response "失败响应"
// @Router /user/current [get]
func GetCurrentUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userId := ctx.GetInt64("userId")

	// 通过用户ID获取用户信息
	user := business.GetUserById(userId)

	if user.UserId == 0 {
		// 用户未找到
		zap.L().Error("User not found", zap.Int64("userId", userId))
		common.ErrorWithCode(ctx, common.CodeNotFound)
		return
	}

	var userVO vo.UserVO
	copier.Copy(&userVO, &user)

	zap.L().Info("Get current user info", zap.Any("user", userVO))
	common.OkWithData(ctx, userVO)
}
