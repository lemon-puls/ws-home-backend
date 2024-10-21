package vo

import "ws-home-backend/model"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserVO struct {
	model.BaseModel
	UserId   int64  `json:"userId"`
	Username string `json:"userName"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   int8   `json:"gender"`
	Age      int8   `json:"age"`
}
