package model

type User struct {
	BaseModel
	UserId   int64   `gorm:"unique; not null" json:"userId"`
	Username string  `gorm:"unique; varchar(255); not null" json:"userName"`
	Password string  `gorm:"type:varchar(255); not null" json:"-"`
	Email    string  `gorm:"type:varchar(255); not null" json:"email"`
	Phone    string  `gorm:"type:varchar(255); not null" json:"phone"`
	Gender   int8    `gorm:"type:tinyint" json:"gender"`
	Age      int8    `gorm:"type:tinyint" json:"age"`
	Avatar   string  `gorm:"type:varchar(255)" json:"avatar"`
	Albums   []Album `gorm:"foreignkey:UserId;references:UserId" json:"albums"`
}
