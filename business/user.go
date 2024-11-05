package business

import (
	"context"
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
		// 默认头像 登陆后可以更换
		Avatar: "https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/3ee5f13fb09879ecb5185e440cef6eb9.png~tplv-uwbnlip3yd-webp.webp",
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
	var key = common.GetUserTokenKey(user.UserId, remoteIP)
	config.RDB.Set(context.Background(), key, accessToken, config.Conf.JwtExpire*time.Minute)
	var userVO vo.UserVO
	copier.Copy(&userVO, user)
	return vo.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserVO:       userVO,
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
