package models

type Song struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Group       string `json:"group" binding:"required"`
	Title       string `json:"title" binding:"required"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongFilter struct {
	Group       string `form:"group"`
	Title       string `form:"song"`
	ReleaseDate string `form:"releaseDate"`
	Page        int    `form:"page" binding:"min=1"`
	PageSize    int    `form:"pageSize" binding:"min=1"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type Lyrics struct {
	SongID        uint     `json:"songId"`
	Title         string   `json:"title"`
	Group         string   `json:"group"`
	Verses        []string `json:"verses"`
	CurrentPage   int      `json:"currentPage"`
	TotalPages    int      `json:"totalPages"`
	VersesPerPage int      `json:"versesPerPage"`
}
