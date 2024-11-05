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

type UpdateUserDTO struct {
	Username    string `json:"userName" binding:"omitempty,min=3"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone" binding:"omitempty,min=11,max=11"`
	Gender      *int8  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Age         *int8  `json:"age" binding:"omitempty,gte=0,lte=150"`
	Avatar      string `json:"avatar" binding:"omitempty,url"`
	OldPassword string `json:"oldPassword" binding:"omitempty,required_with=NewPassword,min=6,max=15"`
	NewPassword string `json:"newPassword" binding:"omitempty,required_with=OldPassword,min=6,max=15"`
}
