package model

type Album struct {
	BaseModel
	Name        string     `gorm:"type:varchar(255); not null;" json:"name"`
	Description string     `gorm:"type:varchar(10000)" json:"description"`
	UserId      int64      `gorm:"not null;" json:"user_id"`
	CoverImg    string     `gorm:"type:varchar(255)" json:"cover_img"`
	User        User       `gorm:"references:UserId" json:"user"`
	AlbumImgs   []AlbumImg `gorm:"foreignkey:AlbumId;references:Id" json:"album_imgs"`
	PhotoCount  int64      `gorm:"-" json:"photo_count"`
}

type AlbumImg struct {
	BaseModel
	AlbumId int64  `gorm:"not null;" json:"album_id"`
	Url     string `gorm:"type:varchar(255); not null;" json:"url"`
	IsRaw   bool   `gorm:"type:tinyint(1);not null;default:0" json:"is_raw"`
	//Album   Album  `gorm:"references:Id" json:"album"`
}
