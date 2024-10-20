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
