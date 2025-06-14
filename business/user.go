package business

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
	"ws-home-backend/common"
	"ws-home-backend/common/jwt"
	"ws-home-backend/config"
	"ws-home-backend/config/db"
	"ws-home-backend/dto"
	"ws-home-backend/model"
	"ws-home-backend/vo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetUserById(userId int64) model.User {
	db := db.GetDB()
	var user model.User
	db.Where(&model.User{UserId: userId}).Find(&user)
	return user
}

func Register(dto dto.RegisterDTO) {
	db := db.GetDB()

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
	db := db.GetDB()
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

	userVO.Avatar, _ = config.GetCosClient().GenerateDownloadPresignedURL(user.Avatar)

	return vo.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserVO:       userVO,
	}
}

// WxLogin 微信小程序登录
func WxLogin(loginDTO dto.LoginDTO, ctx *gin.Context) interface{} {
	// 调用微信接口获取openid和session_key
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.Conf.WxConfig.AppID,
		config.Conf.WxConfig.AppSecret,
		loginDTO.Code)

	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error("调用微信接口查询 openid 失败", zap.Error(err))
		panic(common.NewCustomErrorWithMsg("调用微信接口失败"))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取微信接口响应失败", zap.Error(err))
		panic(common.NewCustomErrorWithMsg("读取微信接口响应失败"))
	}

	var wxResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &wxResp); err != nil {
		zap.L().Error("解析微信接口响应失败", zap.Error(err))
		panic(common.NewCustomErrorWithMsg("解析微信接口响应失败"))
	}

	if wxResp.ErrCode != 0 {
		zap.L().Error("微信登录失败", zap.Error(err))
		panic(common.NewCustomErrorWithMsg(fmt.Sprintf("微信登录失败: %s", wxResp.ErrMsg)))
	}

	// 根据openid查找用户是否存在
	db := db.GetDB()
	var user model.User
	res := db.Where(&model.User{OpenID: wxResp.OpenID}).First(&user)

	// 如果用户不存在，则创建新用户
	if res.RowsAffected == 0 {
		zap.L().Info("用户不存在，无法登录")
		panic(common.NewCustomErrorWithMsg("用户不存在，无法登录"))
	}

	// 生成token
	accessToken, err := jwt.AccessToken(user.UserId)
	if err != nil {
		panic(err)
	}

	refreshToken, err := jwt.RefreshToken(user.UserId)
	if err != nil {
		panic(err)
	}

	// 存储token到Redis
	var remoteIP = ctx.RemoteIP()
	var key = common.GetUserTokenKey(user.UserId, remoteIP)
	config.RDB.Set(context.Background(), key, accessToken, config.Conf.JwtExpire*time.Minute)

	var userVO vo.UserVO
	copier.Copy(&userVO, user)

	userVO.Avatar, _ = config.GetCosClient().GenerateDownloadPresignedURL(user.Avatar)

	return vo.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserVO:       userVO,
	}
}

// 通过手机号查询用户
func GetUserByPhone(phone string) model.User {
	db := db.GetDB()
	var user model.User
	db.Where(&model.User{Phone: phone}).Find(&user)
	return user
}

// 通过用户 ID 查询用户
func GetUserByUserId(userId int64) model.User {
	db := db.GetDB()
	var user model.User
	db.Where(&model.User{UserId: userId}).Find(&user)
	return user
}

func UpdateUser(userId int64, dto dto.UpdateUserDTO) {
	db := db.GetDB()
	user := GetUserByUserId(userId)
	if user.UserId == 0 {
		panic(common.NewCustomError(common.CodeNotFound))
	}

	// 检查手机号是否已被使用
	if dto.Phone != "" && dto.Phone != user.Phone {
		if isUserExists(dto.Phone) {
			panic(common.NewCustomErrorWithMsg("手机号已被使用"))
		}
	}

	// 如果要修改密码
	if dto.OldPassword != "" && dto.NewPassword != "" {
		if !common.Verify(dto.OldPassword, user.Password) {
			panic(common.NewCustomErrorWithMsg("旧密码错误"))
		}
		user.Password = common.Encode(dto.NewPassword)
	}

	// 使用 copier 复制非空字段
	if err := copier.CopyWithOption(&user, &dto, copier.Option{
		IgnoreEmpty: true,
		// 忽略密码字段
		Converters: []copier.TypeConverter{
			{
				SrcType: string(""),
				DstType: string(""),
				Fn: func(src interface{}) (interface{}, error) {
					if src == nil {
						return nil, nil
					}
					return src, nil
				},
			},
		},
	}); err != nil {
		panic(common.NewCustomErrorWithMsg("更新用户信息失败"))
	}

	if err := db.Save(&user).Error; err != nil {
		panic(common.NewCustomErrorWithMsg("更新用户信息失败"))
	}
}

func RefreshToken(refreshToken string, ctx *gin.Context) interface{} {
	// 验证刷新令牌
	claims, err := jwt.VerifyToken(refreshToken)
	if err != nil {
		panic(common.NewCustomError(common.CodeNotLogin))
	}

	// 获取用户信息
	user := GetUserByUserId(claims.UserID)
	if user.UserId == 0 {
		panic(common.NewCustomError(common.CodeNotFound))
	}

	// 生成新的访问令牌
	accessToken, err := jwt.AccessToken(user.UserId)
	if err != nil {
		panic(err)
	}

	// 生成新的刷新令牌
	newRefreshToken, err := jwt.RefreshToken(user.UserId)
	if err != nil {
		panic(err)
	}

	// 更新 Redis 中的访问令牌
	var remoteIP = ctx.RemoteIP()
	var key = common.GetUserTokenKey(user.UserId, remoteIP)
	config.RDB.Set(context.Background(), key, accessToken, config.Conf.JwtExpire*time.Minute)

	var userVO vo.UserVO
	copier.Copy(&userVO, user)

	return vo.Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		UserVO:       userVO,
	}
}
