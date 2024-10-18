package business

import (
	"ws-home-backend/config"
	"ws-home-backend/model"
)

func GetUserById(userId int64) model.User {
	db := config.GetDB()
	var user model.User
	db.Where(&model.User{UserId: userId}).Find(&user)
	return user
}
