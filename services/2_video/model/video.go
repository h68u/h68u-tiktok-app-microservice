package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorId      int64
	Title         string `gorm:"not null"`
	PlayUrl       string `gorm:"not null"`
	CoverUrl      string `gorm:"not null"`
	FavoriteCount int64  `gorm:"column:favorite_count;"`
	CommentCount  int64

	// has many
	Comments  []Comment
	Favorites []Favorite
}

type Comment struct {
	gorm.Model
	UserId  int64
	VideoId int64
	Content string `gorm:"not null"`
}

type Favorite struct {
	gorm.Model
	UserId  int64 `gorm:"column:user_id;"`
	VideoId int64 `gorm:"column:video_id;"`

	// belongs to
	Video Video
}
