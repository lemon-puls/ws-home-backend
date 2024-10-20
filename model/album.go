package model

type Album struct {
	BaseModel
	Name        string     `gorm:"type:varchar(255); not null;" json:"name"`
	Description string     `gorm:"type:varchar(10000)" json:"description"`
	UserId      int64      `gorm:"not null;" json:"user_id"`
	CoverImg    string     `gorm:"type:varchar(255)" json:"cover_img"`
	User        User       `gorm:"references:UserId" json:"user"`
	AlbumImgs   []AlbumImg `gorm:"foreignkey:AlbumId" json:"album_imgs"`
}

type AlbumImg struct {
	BaseModel
	AlbumId int64  `gorm:"not null;" json:"album_id"`
	Url     string `gorm:"type:varchar(255); not null;" json:"url"`
	Album   Album  `gorm:"foreignkey:AlbumId" json:"album"`
}
