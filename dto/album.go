package dto

type AlbumAddDTO struct {
	UserId      int              `json:"user_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CoverImg    string           `json:"cover_img"`
	AlbumImgs   []AlbumImgAddDTO `json:"album_imgs"`
}

type AlbumImgAddDTO struct {
	Url string `json:"url"`
}

type AlbumQueryDTO struct {
	UserId   int64  `json:"user_id" form:"user_id"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Name     string `json:"name" form:"name"`
}
