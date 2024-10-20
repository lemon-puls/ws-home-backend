package business

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"ws-home-backend/common"
	"ws-home-backend/common/jwt"
	"ws-home-backend/config"
	"ws-home-backend/dto"
	"ws-home-backend/model"
	"ws-home-backend/vo"
)

func GetUserById(userId int64) model.User {
	db := config.GetDB()
	var user model.User
	db.Where(&model.User{UserId: userId}).Find(&user)
	return user
}

func Register(dto dto.RegisterDTO) {
	db := config.GetDB()

	if isUserExists(dto.Phone) {
		panic(common.NewCustomErrorWithMsg("User already exists"))
	}

	user := model.User{
		UserId:   config.GenerateID(),
		Username: dto.Username,
		Password: common.Encode(dto.Password),
		Phone:    dto.Phone,
	}

	res := db.Create(&user)
	if res.RowsAffected == 0 {
		panic(common.NewCustomErrorWithMsg("Failed to create user"))
	}
}

func isUserExists(phone string) bool {
	db := config.GetDB()
	var user model.User
	res := db.Where(&model.User{Phone: phone}).Find(&user)
	if res.RowsAffected > 0 {
		return true
	}
	return false
}

func Login(loginDTO dto.LoginDTO, ctx *gin.Context) interface{} {
	user := GetUserByPhone(loginDTO.Phone)
	if user.UserId == 0 {
		panic(common.NewCustomError(common.CodeNotFound))
	}

	if !common.Verify(loginDTO.Password, user.Password) {
		panic(common.NewCustomErrorWithMsg("Incorrect password"))
	}

	accessToken, err := jwt.AccessToken(user.UserId)
	if err != nil {
		panic(err)
	}

	refreshToken, err := jwt.RefreshToken(user.UserId)
	if err != nil {
		panic(err)
	}
	//
	var remoteIP = ctx.RemoteIP()
	var key = common.KeyUserTokenPrefix +
		strconv.FormatInt(user.UserId, 10) + ":" + remoteIP
	config.RDB.Set(context.Background(), key, accessToken, config.Conf.JwtExpire*time.Minute)
	return vo.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// 通过手机号查询用户
func GetUserByPhone(phone string) model.User {
	db := config.GetDB()
	var user model.User
	db.Where(&model.User{Phone: phone}).Find(&user)
	return user
}

// 通过用户 ID 查询用户
func GetUserByUserId(userId int64) model.User {
	db := config.GetDB()
	var user model.User
	db.Where(&model.User{UserId: userId}).Find(&user)
	return user
}
