package vo

import "ws-home-backend/model"

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserVO       UserVO `json:"userVO"`
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
