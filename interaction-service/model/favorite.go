package model

import "time"

// Favorite 点赞表 /*
type Favorite struct {
	UserId     int64     `gorm:"column:user_id"`
	VideoId    int64     `gorm:"column:video_id"`
	CreateTime time.Time `gorm:"autoCreateTime"`
}
