package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromId   int64
	ToUserId int64
	Content  string `gorm:"not null"`
}

type Friend struct {
	gorm.Model
	UserId   int64
	FriendId int64
}
