package dto

type RegisterDTO struct {
	Username string `json:"userName" binding:"required,min=3"`
	Phone    string `json:"phone" binding:"required,min=11,max=11"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}

type LoginDTO struct {
	Phone    string `json:"phone" binding:"required,min=11,max=11"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}
