package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorId      int64
	Title         string `gorm:"not null"`
	PlayUrl       string `gorm:"not null"`
	CoverUrl      string `gorm:"not null"`
	FavoriteCount int64
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
	UserId  int64
	VideoId int64
}
